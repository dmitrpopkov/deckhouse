/*
Copyright 2023 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hooks

import (
	"os"
	"path"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	deckhouse_config "github.com/deckhouse/deckhouse/go_lib/deckhouse-config"
	module_manager "github.com/deckhouse/deckhouse/go_lib/deckhouse-config/module-manager"
	"github.com/deckhouse/deckhouse/modules/005-external-module-manager/hooks/internal/apis/v1alpha1"
	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = Describe("Modules :: external module manager :: hooks :: apply release ::", func() {
	var tmpDir string

	f := HookExecutionConfigInit(`
global:
  deckhouseVersion: "12345"
  modulesImages:
    registry:
      base: registry.deckhouse.io/deckhouse/fe
external-module-manager:
  internal: {}
`, `{}`)
	f.RegisterCRD("deckhouse.io", "v1alpha1", "ModuleRelease", false)
	f.RegisterCRD("deckhouse.io", "v1alpha1", "ModuleUpdatePolicy", false)

	Context("Cluster has Deployed, Superseded and Pending ModuleReleases", func() {
		BeforeEach(func() {
			tmpDir, _ = os.MkdirTemp(os.TempDir(), "exrelease-*")
			_ = os.Mkdir(tmpDir+"/modules", 0777)
			_ = os.Setenv("EXTERNAL_MODULES_DIR", tmpDir)
			testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")

			f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
status:
  phase: Deployed
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.2
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.2
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.3
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.3
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleUpdatePolicy
metadata:
  name: echoserver-policy
spec:
  moduleReleaseSelector:
    labelSelector:
      matchLabels:
        source: echoserver
        module: echoserver
  releaseChannel: Stable
  update:
    mode: Manual
`)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		It("module should be in Pending state", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhaseDeployed))
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.2").Field("status.phase").String()).To(Equal(v1alpha1.PhaseSuperseded))
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.3").Field("status.phase").String()).To(Equal(v1alpha1.PhasePending))
		})
	})

	Context("Cluster has ModuleRelease not associated with any policy", func() {
		BeforeEach(func() {
			tmpDir, _ = os.MkdirTemp(os.TempDir(), "exrelease-*")
			_ = os.Mkdir(tmpDir+"/modules", 0777)
			_ = os.Setenv("EXTERNAL_MODULES_DIR", tmpDir)
			testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")

			f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
status:
  phase: Pending
`)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		It("module should be in PolicyUndefined state", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhasePolicyUndefined))
		})
	})

	Context("Cluster has pending ModuleRelease and policy Auto with window and in time", func() {
		BeforeEach(func() {
			tmpDir, _ = os.MkdirTemp(os.TempDir(), "exrelease-*")
			_ = os.Mkdir(tmpDir+"/modules", 0777)
			_ = os.Setenv("EXTERNAL_MODULES_DIR", tmpDir)
			testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")

			f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
status:
  phase: Pending
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleUpdatePolicy
metadata:
  name: echoserver-policy
spec:
  moduleReleaseSelector:
    labelSelector:
      matchLabels:
        source: echoserver
        module: echoserver
  releaseChannel: Stable
  update:
    mode: Auto
    windows:
    - days:
      - Mon
      - Tue
      - Wed
      - Thu
      - Fri
      - Sat
      - Sun
      from: "00:00"
      to: "23:59"
`)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		It("module symlink should be created", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhaseDeployed))
			moduleLinks, err := os.ReadDir(tmpDir + "/modules")
			if err != nil {
				Fail(err.Error())
			}
			Expect(moduleLinks).To(HaveLen(1))
			Expect(moduleLinks[0].Name()).To(Equal("900-echoserver"))
		})
	})

	Context("Cluster has pending ModuleRelease and policy Auto with window and not in time", func() {
		BeforeEach(func() {
			tmpDir, _ = os.MkdirTemp(os.TempDir(), "exrelease-*")
			_ = os.Mkdir(tmpDir+"/modules", 0777)
			_ = os.Setenv("EXTERNAL_MODULES_DIR", tmpDir)
			testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")

			f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
status:
  phase: Pending
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleUpdatePolicy
metadata:
  name: echoserver-policy
spec:
  moduleReleaseSelector:
    labelSelector:
      matchLabels:
        source: echoserver
        module: echoserver
  releaseChannel: Stable
  update:
    mode: Auto
    windows:
    - days:
      - Sun
      from: "00:05"
      to: "00:06"
`)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		It("module should be waiting for update window", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhasePending))
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.message").String()).To(ContainSubstring(waitingForWindow[:len(waitingForWindow)-2]))
		})
	})

	Context("Cluster has pending ModuleRelease and policy Auto without windows", func() {
		BeforeEach(func() {
			tmpDir, _ = os.MkdirTemp(os.TempDir(), "exrelease-*")
			_ = os.Mkdir(tmpDir+"/modules", 0777)
			_ = os.Setenv("EXTERNAL_MODULES_DIR", tmpDir)
			testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")

			f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
status:
  phase: Pending
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleUpdatePolicy
metadata:
  name: echoserver-policy
spec:
  moduleReleaseSelector:
    labelSelector:
      matchLabels:
        source: echoserver
        module: echoserver
  releaseChannel: Stable
  update:
    mode: Auto
`)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		It("module symlink should be created", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhaseDeployed))
			moduleLinks, err := os.ReadDir(tmpDir + "/modules")
			if err != nil {
				Fail(err.Error())
			}
			Expect(moduleLinks).To(HaveLen(1))
			Expect(moduleLinks[0].Name()).To(Equal("900-echoserver"))
		})

		Context("ModuleRelease was deleted", func() {
			BeforeEach(func() {
				st := f.KubeStateSet(``)
				f.BindingContexts.Set(st)
				fsSynchronized = false
				f.RunHook()
			})

			It("should delete module from FS", func() {
				Expect(f).To(ExecuteSuccessfully())
				moduleLinks, err := os.ReadDir(tmpDir + "/modules")
				if err != nil {
					Fail(err.Error())
				}
				Expect(moduleLinks).To(HaveLen(0))
			})
		})
	})

	Context("Cluster has unapproved pending ModuleRelease and policy Manual", func() {
		BeforeEach(func() {
			tmpDir, _ = os.MkdirTemp(os.TempDir(), "exrelease-*")
			_ = os.Mkdir(tmpDir+"/modules", 0777)
			_ = os.Setenv("EXTERNAL_MODULES_DIR", tmpDir)
			testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")

			f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
status:
  phase: Pending
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleUpdatePolicy
metadata:
  name: echoserver-policy
spec:
  moduleReleaseSelector:
    labelSelector:
      matchLabels:
        source: echoserver
        module: echoserver
  releaseChannel: Stable
  update:
    mode: Manual
`)
			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		It("module should be waiting for manual approval", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhasePending))
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.message").String()).To(Equal(manualApproval))
		})
	})

	Context("Cluster has approved pending ModuleRelease", func() {
		BeforeEach(func() {
			tmpDir, _ = os.MkdirTemp(os.TempDir(), "exrelease-*")
			_ = os.Mkdir(tmpDir+"/modules", 0777)
			_ = os.Setenv("EXTERNAL_MODULES_DIR", tmpDir)
			testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")

			f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  annotations:
    release.deckhouse.io/approved: "true"
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
status:
  phase: Pending
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleUpdatePolicy
metadata:
  name: echoserver-policy
spec:
  moduleReleaseSelector:
    labelSelector:
      matchLabels:
        source: echoserver
        module: echoserver
  releaseChannel: Stable
  update:
    mode: Manual
`)

			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		It("module symlink should be created", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhaseDeployed))
			moduleLinks, err := os.ReadDir(tmpDir + "/modules")
			if err != nil {
				Fail(err.Error())
			}
			Expect(moduleLinks).To(HaveLen(1))
			Expect(moduleLinks[0].Name()).To(Equal("900-echoserver"))
		})

		Context("ModuleRelease was deleted", func() {
			BeforeEach(func() {
				st := f.KubeStateSet(``)
				f.BindingContexts.Set(st)
				fsSynchronized = false
				f.RunHook()
			})

			It("should delete module from FS", func() {
				Expect(f).To(ExecuteSuccessfully())
				moduleLinks, err := os.ReadDir(tmpDir + "/modules")
				if err != nil {
					Fail(err.Error())
				}
				Expect(moduleLinks).To(HaveLen(0))
			})
		})
	})

	Context("Cluster has ModuleRelease with custom weight", func() {
		BeforeEach(func() {
			tmpDir, _ = os.MkdirTemp(os.TempDir(), "exrelease-*")
			_ = os.Mkdir(tmpDir+"/modules", 0777)
			_ = os.Setenv("EXTERNAL_MODULES_DIR", tmpDir)
			testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")

			f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
  weight: 987
status:
  phase: Pending
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleUpdatePolicy
metadata:
  name: echoserver-policy
spec:
  moduleReleaseSelector:
    labelSelector:
      matchLabels:
        source: echoserver
        module: echoserver
  releaseChannel: Stable
  update:
    mode: Auto
`)

			f.BindingContexts.Set(f.GenerateBeforeHelmContext())
			f.RunHook()
		})

		AfterEach(func() {
			_ = os.RemoveAll(tmpDir)
		})

		It("module symlink should be created with custom weight", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhaseDeployed))
			moduleLinks, err := os.ReadDir(tmpDir + "/modules")
			if err != nil {
				Fail(err.Error())
			}
			Expect(moduleLinks).To(HaveLen(1))
			Expect(moduleLinks[0].Name()).To(Equal("987-echoserver"))
		})

		Context("ModuleRelease was changed with another weight", func() {
			BeforeEach(func() {
				testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.1")
				testCreateModuleOnFS(tmpDir, "echoserver", "v0.0.2")
				f.KubeStateSet(``) // Empty cluster
				st := f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.1
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.1
  weight: 987
status:
  phase: Deployed
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: echoserver-v0.0.2
  labels:
    source: echoserver
    module: echoserver
    module-update-policy: echoserver-policy
spec:
  moduleName: echoserver
  version: 0.0.2
  weight: 913
status:
  phase: Pending
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleUpdatePolicy
metadata:
  name: echoserver-policy
spec:
  moduleReleaseSelector:
    labelSelector:
      matchLabels:
        source: echoserver
        module: echoserver
  releaseChannel: Stable
  update:
    mode: Auto
`)
				f.BindingContexts.Set(st)
				fsSynchronized = false
				f.RunHook()
			})
			AfterEach(func() {
				_ = os.RemoveAll(tmpDir)
			})

			It("should change module symlink", func() {
				Expect(f).To(ExecuteSuccessfully())
				Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhaseSuperseded))
				Expect(f.KubernetesGlobalResource("ModuleRelease", "echoserver-v0.0.2").Field("status.phase").String()).To(Equal(v1alpha1.PhaseDeployed))
				moduleLinks, err := os.ReadDir(tmpDir + "/modules")
				if err != nil {
					Fail(err.Error())
				}
				Expect(moduleLinks).To(HaveLen(1))
				Expect(moduleLinks[0].Name()).To(Equal("913-echoserver"))
			})
		})

		Context("Target module does not exist on fs", func() {
			BeforeEach(func() {
				mm, _ := module_manager.InitBasic("", "")
				_ = mm.RegisterModules()
				deckhouse_config.InitService(mm)
				st := f.KubeStateSet(`
---
apiVersion: deckhouse.io/v1alpha1
kind: ModuleRelease
metadata:
  name: absent-v0.0.1
spec:
  moduleName: absent
  version: 0.0.1
  weight: 987
status:
  phase: Deployed
`)
				f.BindingContexts.Set(st)
				fsSynchronized = true
				f.RunHook()
			})
			AfterEach(func() {
				_ = os.RemoveAll(tmpDir)
			})

			It("Should suspend the release", func() {
				Expect(f).To(ExecuteSuccessfully())
				Expect(f.KubernetesGlobalResource("ModuleRelease", "absent-v0.0.1").Field("status.phase").String()).To(Equal(v1alpha1.PhaseSuspended))
				Expect(f.KubernetesGlobalResource("ModuleRelease", "absent-v0.0.1").Field("status.message").String()).To(Equal("Desired version of the module met problems: not found"))
			})
		})
	})
})

// nolint: unparam
func testCreateModuleOnFS(tmpDir, moduleName, moduleVersion string) {
	modulePath := path.Join(tmpDir, moduleName, moduleVersion)
	_ = os.MkdirAll(modulePath, 0666)
	_, _ = os.Create(path.Join(modulePath, "Chart.yaml"))
	_, _ = os.Create(path.Join(modulePath, "values.yaml"))
}

func TestSymlinkFinder(t *testing.T) {
	mt, err := os.MkdirTemp("", "target-*")
	require.NoError(t, err)
	defer os.RemoveAll(mt)

	tmp, err := os.MkdirTemp("", "modules-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	_ = os.Symlink(mt, path.Join(tmp, "100-module1"))
	_ = os.Symlink(mt, path.Join(tmp, "200-module2"))
	_ = os.Symlink(mt, path.Join(tmp, "300-module3"))
	_, _ = os.Create(path.Join(tmp, "333-module2"))

	res1, err := findExistingModuleSymlink(tmp, "module2")
	require.NoError(t, err)
	assert.True(t, strings.HasSuffix(res1, path.Join(tmp, "200-module2")))

	res2, err := findExistingModuleSymlink(tmp, "module5")
	require.NoError(t, err)
	assert.Empty(t, res2)
}
