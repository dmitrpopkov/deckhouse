/*
Copyright 2022 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"context"
	"strings"

	. "github.com/deckhouse/deckhouse/testing/helm"
	. "github.com/deckhouse/deckhouse/testing/hooks"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/google/go-containerregistry/pkg/authn"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlSrlzr "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"sigs.k8s.io/yaml"
)

var _ = Describe("Modules :: delivery :: hooks :: werf_sources ::", func() {
	decUnstructured := yamlSrlzr.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	Context("parsing of WerfSources resource into inner formet", func() {
		table.DescribeTable("Parsing werf_sources", func(wsyaml string, expected werfSource) {
			// Setup
			obj := &unstructured.Unstructured{}
			_, _, err := decUnstructured.Decode([]byte(wsyaml), nil, obj)
			Expect(err).ToNot(HaveOccurred())

			// Action
			ws, err := filterWerfSource(obj)

			// Assert
			Expect(err).ToNot(HaveOccurred())
			Expect(ws).To(Equal(expected))
		},
			table.Entry("Minimal: only image repo", `
apiVersion: deckhouse.io/v1alpha1
kind: WerfSource
metadata:
  name: minimal
spec:
  imageRepo: cr.example.com/the/path
`,
				werfSource{
					name:   "minimal",
					repo:   "cr.example.com/the/path",
					apiURL: "https://cr.example.com",
					argocdRepo: &argocdRepoConfig{
						project: "default",
					},
				}),

			table.Entry("Full", `
apiVersion: deckhouse.io/v1alpha1
kind: WerfSource
metadata:
  name: full-object
spec:
  imageRepo: cr.example.com/the/path
  apiURL: https://different.example.com
  pullSecretName: registry-credentials
  argocdRepoEnabled: true
  argocdRepo:
    project: ecommerce

`,
				werfSource{
					name:   "full-object",
					repo:   "cr.example.com/the/path",
					apiURL: "https://different.example.com",

					pullSecretName: "registry-credentials",
					argocdRepo: &argocdRepoConfig{
						project: "ecommerce",
					},
				}),

			table.Entry("argocdRepoEnabled=false omits the repo config for Argo", `
apiVersion: deckhouse.io/v1alpha1
kind: WerfSource
metadata:
  name: repo-off
spec:
  imageRepo: cr.example.com/the/path
  argocdRepoEnabled: false
`,
				werfSource{
					name:   "repo-off",
					repo:   "cr.example.com/the/path",
					apiURL: "https://cr.example.com",
				}),

			table.Entry("argocdRepoEnabled=false omits the repo config for Argo even when repo options are specified ", `
apiVersion: deckhouse.io/v1alpha1
kind: WerfSource
metadata:
  name: repo-off-yet-specified
spec:
  imageRepo: cr.example.com/the/path
  argocdRepoEnabled: false
  argocdRepo:
    project: actually-skipped
`,
				werfSource{
					name:   "repo-off-yet-specified",
					repo:   "cr.example.com/the/path",
					apiURL: "https://cr.example.com",
				}),

			table.Entry("Argo CD non-defaul project", `
apiVersion: deckhouse.io/v1alpha1
kind: WerfSource
metadata:
  name: not-default-project
spec:
  imageRepo: cr.example.com/the/path
  argocdRepo:
    project: greater-good
`,
				werfSource{
					name:   "not-default-project",
					repo:   "cr.example.com/the/path",
					apiURL: "https://cr.example.com",
					argocdRepo: &argocdRepoConfig{
						project: "greater-good",
					},
				}),
		)
	})

	Context("Converting werf sources to configs ", func() {
		ws1 := werfSource{
			name:           "ws1",
			repo:           "cr-1.example.com/the/path",
			apiURL:         "https://cr.example.com",
			pullSecretName: "registry-credentials-1",
			argocdRepo: &argocdRepoConfig{
				project: "default",
			},
		}

		ws2 := werfSource{
			name:           "ws2",
			repo:           "cr-2.example.com/the/path",
			apiURL:         "https://registry-api.other.com",
			pullSecretName: "registry-credentials-2",
			argocdRepo: &argocdRepoConfig{
				project: "top-secret",
			},
		}

		ws3 := werfSource{
			name: "ws3-no-creds",
			repo: "open.example.com/the/path",
			argocdRepo: &argocdRepoConfig{
				project: "default",
			},
		}

		ws4 := werfSource{
			name:           "ws4-no-repo",
			repo:           "cr-4.example.com/the/path",
			pullSecretName: "registry-credentials-4",
		}

		configsBySecretName := map[string]dockerFileConfig{
			"registry-credentials-1":      {Auths: map[string]authn.AuthConfig{"cr-1.example.com": {Username: "n-1", Password: "pwd-1"}}},
			"registry-credentials-2":      {Auths: map[string]authn.AuthConfig{"cr-2.example.com": {Username: "n-2", Password: "pwd-2"}}},
			"unused-registry-credentials": {Auths: map[string]authn.AuthConfig{"noop.example.com": {Username: "n-3", Password: "pwd-3"}}},
			"registry-credentials-4":      {Auths: map[string]authn.AuthConfig{"cr-4.example.com": {Username: "n-4", Password: "pwd-4"}}},
		}

		vals, err := mapWerfSources([]werfSource{ws1, ws2, ws3, ws4}, configsBySecretName)

		It("returns no errors", func() {
			Expect(err).ToNot(HaveOccurred())
		})
		It("parses to argo cd repositories as expected", func() {
			Expect(vals.ArgoCD.Repositories).To(ConsistOf(
				argocdHelmOCIRepository{
					Name:     "ws1",
					URL:      "cr-1.example.com/the/path",
					Username: "n-1",
					Password: "pwd-1",
					Project:  "default",
				},
				argocdHelmOCIRepository{
					Name:     "ws2",
					URL:      "cr-2.example.com/the/path",
					Username: "n-2",
					Password: "pwd-2",
					Project:  "top-secret",
				},
				argocdHelmOCIRepository{
					Name:    "ws3-no-creds",
					URL:     "open.example.com/the/path",
					Project: "default",
				},
			))
		})

		It("parses to argo cd image updater registries as expected", func() {
			Expect(vals.ArgoCDImageUpdater.Registries).To(ConsistOf(
				imageUpdaterRegistry{
					Name:        "ws1",
					Prefix:      "cr-1.example.com",
					APIURL:      "https://cr.example.com",
					Credentials: "pullsecret:d8-delivery/registry-credentials-1",
					Default:     false,
				},
				imageUpdaterRegistry{
					Name:        "ws2",
					Prefix:      "cr-2.example.com",
					APIURL:      "https://registry-api.other.com",
					Credentials: "pullsecret:d8-delivery/registry-credentials-2",
					Default:     false,
				},
				imageUpdaterRegistry{
					Name:    "ws3-no-creds",
					Prefix:  "open.example.com",
					APIURL:  "https://open.example.com",
					Default: false,
				},
				imageUpdaterRegistry{
					Name:        "ws4-no-repo",
					Prefix:      "cr-4.example.com",
					APIURL:      "https://cr-4.example.com",
					Credentials: "pullsecret:d8-delivery/registry-credentials-4",
					Default:     false,
				},
			))
		})
	})

	Context("YAML rendering of Argo CD repo", func() {
		It("renders full struct", func() {
			b, err := yaml.Marshal(argocdHelmOCIRepository{
				Name:     "ws1",
				URL:      "cr-1.example.com/the/path",
				Username: "n-1",
				Password: "pwd-1",
				Project:  "default",
			})

			expected := `
name: ws1
password: pwd-1
project: default
url: cr-1.example.com/the/path
username: n-1
`
			Expect(err).ToNot(HaveOccurred())
			Expect("\n" + string(b)).To(Equal(expected))
		})
		It("omits optional fields", func() {
			b, err := yaml.Marshal(argocdHelmOCIRepository{
				Name:     "ws1",
				URL:      "cr-1.example.com/the/path",
				Username: "",
				Password: "",
				Project:  "default",
			})

			expected := `
name: ws1
project: default
url: cr-1.example.com/the/path
`
			Expect(err).ToNot(HaveOccurred())
			Expect("\n" + string(b)).To(Equal(expected))
		})
	})

	Context("YAML rendering of Argo CD Image Updater registry", func() {
		It("renders full struct", func() {
			b, err := yaml.Marshal(imageUpdaterRegistry{
				Name:        "ws1",
				Prefix:      "cr-1.example.com",
				APIURL:      "https://cr.example.com",
				Credentials: "pullsecret:d8-delivery/registry-credentials-1",
				Default:     false,
			})
			expected := `
api_url: https://cr.example.com
credentials: pullsecret:d8-delivery/registry-credentials-1
default: false
name: ws1
prefix: cr-1.example.com
`
			Expect(err).ToNot(HaveOccurred())
			Expect("\n" + string(b)).To(Equal(expected))
		})

		It("omits optional fields", func() {
			b, err := yaml.Marshal(imageUpdaterRegistry{
				Name:    "ws1",
				Prefix:  "cr-1.example.com",
				APIURL:  "https://cr.example.com",
				Default: false,
			})
			expected := `
api_url: https://cr.example.com
default: false
name: ws1
prefix: cr-1.example.com
`
			Expect(err).ToNot(HaveOccurred())
			Expect("\n" + string(b)).To(Equal(expected))
		})
	})

	Context("Parsing registry secrets to the map", func() {
		state := `
data:
  # 				 {"auths":{"cr-1.example.com":{"username":"n-1","password":"pwd-1"}}}
  .dockerconfigjson: eyJhdXRocyI6eyJjci0xLmV4YW1wbGUuY29tIjp7InVzZXJuYW1lIjoibi0xIiwicGFzc3dvcmQiOiJwd2QtMSJ9fX0=
apiVersion: v1
kind: Secret
metadata:
  name: registry-credentials-1
  namespace: d8-delivery
type: kubernetes.io/dockerconfigjson
---
data:
  # 				 {"auths":{"cr-2.example.com":{"username":"n-2","password":"pwd-2"}}}
  .dockerconfigjson: eyJhdXRocyI6eyJjci0yLmV4YW1wbGUuY29tIjp7InVzZXJuYW1lIjoibi0yIiwicGFzc3dvcmQiOiJwd2QtMiJ9fX0=
apiVersion: v1
kind: Secret
metadata:
  name: registry-credentials-2
  namespace: d8-delivery
type: kubernetes.io/dockerconfigjson
---
data:
  # 				 {"auths":{"cr-1.example.com":{"username":"n-1","password":"pwd-1"},"cr-2.example.com":{"username":"n-2other","password":"pwd-2other"}}}
  .dockerconfigjson: eyJhdXRocyI6eyJjci0xLmV4YW1wbGUuY29tIjp7InVzZXJuYW1lIjoibi0xIiwicGFzc3dvcmQiOiJwd2QtMSJ9LCJjci0yLmV4YW1wbGUuY29tIjp7InVzZXJuYW1lIjoibi0yb3RoZXIiLCJwYXNzd29yZCI6InB3ZC0yb3RoZXIifX19
apiVersion: v1
kind: Secret
metadata:
  name: mixed-registries-1-2
  namespace: d8-delivery
type: kubernetes.io/dockerconfigjson
---
data:
  # 				 {"auths":{"cr-1.example.com":{"username":"n-1other","password":"pwd-1other"}}}
  .dockerconfigjson: eyJhdXRocyI6eyJjci0xLmV4YW1wbGUuY29tIjp7InVzZXJuYW1lIjoibi0xb3RoZXIiLCJwYXNzd29yZCI6InB3ZC0xb3RoZXIifX19
apiVersion: v1
kind: Secret
metadata:
  name: registry-credentials-1-other
  namespace: d8-delivery
type: kubernetes.io/dockerconfigjson
---
data:
  # 				 {"auths":{"cr-3.example.com":{"auth":"bi0zOnB3ZC0z"}}}
  .dockerconfigjson: eyJhdXRocyI6eyJjci0zLmV4YW1wbGUuY29tIjp7ImF1dGgiOiJiaTB6T25CM1pDMHoifX19
apiVersion: v1
kind: Secret
metadata:
  name: registry-credentials-3-auth
  namespace: d8-delivery
type: kubernetes.io/dockerconfigjson
---
data: {}
apiVersion: v1
kind: Secret
metadata:
  name: ignored-no-data
  namespace: d8-delivery
type: kubernetes.io/dockerconfigjson
---
data:
  .dockerconfigjson: eyJhdXRocyI6eyJjci0zLmV4YW1wbGUuY29tIjp7ImF1dGgiOiJiaTB6T25CM1pDMHoifX19
apiVersion: v1
kind: Secret
metadata:
  name: ignored-type-opaque
  namespace: d8-delivery
type: Opaque
---
data:
  ca.crt: LS0tLS0tLS0t
  namespace: LS0tLS0tLS0t
  token: LS0tLS0tLS0t
apiVersion: v1
kind: Secret
metadata:
  name: ignored-type-service-account-token
  namespace: d8-delivery
  type: kubernetes.io/service-account-token
`
		manifests := strings.Split(state, "---")
		filterResults := make([]go_hook.FilterResult, len(manifests))
		for i, manifest := range manifests {
			obj := &unstructured.Unstructured{}
			_, _, err := decUnstructured.Decode([]byte(manifest), nil, obj)
			if err != nil {
				panic(err)
			}
			filterResult, err := filterDockerConfigJSON(obj)
			if err != nil {
				panic(err)
			}
			filterResults[i] = filterResult
		}

		want := map[string]dockerFileConfig{
			"registry-credentials-1":       {Auths: map[string]authn.AuthConfig{"cr-1.example.com": {Username: "n-1", Password: "pwd-1"}}},
			"registry-credentials-2":       {Auths: map[string]authn.AuthConfig{"cr-2.example.com": {Username: "n-2", Password: "pwd-2"}}},
			"mixed-registries-1-2":         {Auths: map[string]authn.AuthConfig{"cr-1.example.com": {Username: "n-1", Password: "pwd-1"}, "cr-2.example.com": {Username: "n-2other", Password: "pwd-2other"}}},
			"registry-credentials-1-other": {Auths: map[string]authn.AuthConfig{"cr-1.example.com": {Username: "n-1other", Password: "pwd-1other"}}},
			"registry-credentials-3-auth":  {Auths: map[string]authn.AuthConfig{"cr-3.example.com": {Auth: "bi0zOnB3ZC0z"}}},
		}

		It("parses secrets with .dockerconfigjson", func() {
			got, err := parseDockerConfigsBySecretName(filterResults)
			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(Equal(want))
		})

		table.DescribeTable("extracting credentials from the map by secret name and registry", func(secret, registry, wantUsername, wantPassword string) {
			u, p := extractCredentials(want, secret, registry)
			Expect(u).To(Equal(wantUsername))
			Expect(p).To(Equal(wantPassword))
		},
			table.Entry("registry & secret 1",
				"registry-credentials-1", "cr-1.example.com",
				"n-1", "pwd-1"),
			table.Entry("registry & secret 2",
				"registry-credentials-2", "cr-2.example.com",
				"n-2", "pwd-2"),
			table.Entry("shared registry/secret with same creds",
				"mixed-registries-1-2", "cr-1.example.com",
				"n-1", "pwd-1"),
			table.Entry("shared registry/secret with different creds",
				"mixed-registries-1-2", "cr-2.example.com",
				"n-2other", "pwd-2other"),
			table.Entry("shared registry with different creds",
				"registry-credentials-1-other", "cr-1.example.com",
				"n-1other", "pwd-1other"),
			table.Entry("auth field",
				"registry-credentials-3-auth", "cr-3.example.com",
				"n-3", "pwd-3"),
			table.Entry("registry and secret unknown",
				"registry-unknown", "secret-unknown",
				"", ""),
			table.Entry("registry unknown",
				"registry-unknown", "cr-1.example.com",
				"", ""),
			table.Entry("secret unknown",
				"registry-credentials-1", "cr-0000.example.com",
				"", ""),
		)
	})

	XContext("Hook flow", func() {
		f := HookExecutionConfigInit(`{}`, `{}`)
		f.RegisterCRD("deckhouse.io", "v1alpha1", "WerfSource", false)

		// Docker config JSON for cr-1.example.com
		// {"auths":{"cr-1.example.com":{"username":"n-1","password":"pwd-1","auth":"test-auth"}}}
		// ↓↓↓
		state := `
---
data:
  .dockerconfigjson: eyJhdXRocyI6eyJjci0xLmV4YW1wbGUuY29tIjp7InVzZXJuYW1lIjoibi0xIiwicGFzc3dvcmQiOiJwd2QtMSIsImF1dGgiOiJ0ZXN0LWF1dGgifX19
apiVersion: v1
kind: Secret
metadata:
  name: registry-credentials-1
  namespace: d8-delivery
type: kubernetes.io/dockerconfigjson
---
apiVersion: deckhouse.io/v1alpha1
kind: WerfSource
metadata:
  name: ws1
spec:
  imageRepo: cr-1.example.com/the/path
  apiURL: https://cr.example.com
  pullSecretName: registry-credentials-1
---
apiVersion: deckhouse.io/v1alpha1
kind: WerfSource
metadata:
  name: ws3-no-creds
spec:
  imageRepo: open.example.com/the/path
`

		BeforeEach(func() {
			f.BindingContexts.Set(f.KubeStateSet(state))
			f.RunHook()
		})

		It("Runs successfully", func() {
			Expect(f).To(ExecuteSuccessfully())
		})

		It("creates repo secrets for ArgoCD", func() {
			repo1, err := f.KubeClient().CoreV1().Secrets("d8-delivery").Get(context.Background(), "repo-ws1", metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(repo1.Data["username"]).To(Equal([]byte("n-1")))
		})

		It("creates configmap for Argo CD image updater", func() {
			cm, err := f.KubeClient().CoreV1().ConfigMaps("d8-delivery").Get(context.Background(), "argocd-image-updater-config", metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect("\n" + cm.Data["registries.yaml"]).To(Equal(`
registries:
- api_url: https://cr.example.com
  credentials: pullsecret:d8-delivery/registry-credentials-1
  default: false
  name: ws1
  prefix: cr-1.example.com
- api_url: https://cr.example.com
  default: false
  name: ws3-no-creds
  prefix: open.example.com
`))
		})
	})

	Context("templates", func() {
		Context("Repo and registry configurations", func() {
			f := SetupHelmConfig(``)

			values := internalValues{
				ArgoCD: internalArgoCDValues{
					Repositories: []argocdHelmOCIRepository{
						{
							Name:     "ws1",
							URL:      "cr-1.example.com/the/path",
							Username: "n-1",
							Password: "pwd-1",
							Project:  "default",
						},
						{
							Name:    "ws3-no-creds",
							URL:     "open.example.com/the/path",
							Project: "default",
						},
					},
				},
				ArgoCDImageUpdater: internalUpdaterValues{
					Registries: []imageUpdaterRegistry{
						{
							Name:        "ws1",
							Prefix:      "cr-1.example.com",
							APIURL:      "https://cr.example.com",
							Credentials: "pullsecret:d8-delivery/registry-credentials-1",
							Default:     false,
						},
						{
							Name:    "ws3-no-creds",
							Prefix:  "open.example.com",
							APIURL:  "https://open.example.com",
							Default: false,
						},
					},
				},
			}

			BeforeEach(func() {
				f.ValuesSetFromYaml("global", globalValues)
				f.ValuesSet("global.modulesImages", GetModulesImages())
				f.ValuesSetFromYaml("delivery", moduleValues)
				f.ValuesSet("delivery.internal", values)
				f.HelmRender()
			})

			It("rendered without an error", func() {
				Expect(f.RenderError).ShouldNot(HaveOccurred())
			})

			It("creates repo secrets for ArgoCD", func() {
				repo1 := f.KubernetesResource("Secret", "d8-delivery", "repo-ws1")
				Expect(repo1.Exists()).To(BeTrue())
				Expect(repo1.Field("stringData").String()).Should(MatchYAML(`{
					"type": "helm",
					"enableOCI": "true",
					"name": "ws1",
					"username": "n-1",
					"password": "pwd-1",
					"project": "default",
					"url": "cr-1.example.com/the/path"
				}`))

				repo3 := f.KubernetesResource("Secret", "d8-delivery", "repo-ws3-no-creds")
				Expect(repo3.Exists()).To(BeTrue())
				Expect(repo3.Field("stringData").String()).Should(MatchYAML(`{
					"type": "helm",
					"enableOCI": "true",
					"name": "ws3-no-creds",
					"project": "default",
					"url": "open.example.com/the/path"
				}`))
			})

			It("creates configmap for Argo CD image updater", func() {
				updaterConfig := f.KubernetesResource("ConfigMap", "d8-delivery", "argocd-image-updater-config")
				Expect(updaterConfig.Exists()).To(BeTrue())
				Expect(updaterConfig.Field("data").Map()["registries.conf"].String()).Should(MatchYAML(`
registries:
- api_url: https://cr.example.com
  credentials: pullsecret:d8-delivery/registry-credentials-1
  default: false
  name: ws1
  prefix: cr-1.example.com
- api_url: https://open.example.com
  default: false
  name: ws3-no-creds
  prefix: open.example.com
`))
			})
		})
	})
})

type mockCredGetter map[string][]byte

func (cg mockCredGetter) Get(context.Context) (map[string][]byte, error) {
	return cg, nil
}

const globalValues = `
clusterConfiguration:
  apiVersion: deckhouse.io/v1
  cloud:
    prefix: myprefix
    provider: OpenStack
  clusterDomain: cluster.local
  clusterType: "Cloud"
  defaultCRI: Docker
  kind: ClusterConfiguration
  kubernetesVersion: "1.21"
  podSubnetCIDR: 10.111.0.0/16
  podSubnetNodeCIDRPrefix: "24"
  serviceSubnetCIDR: 10.222.0.0/16
enabledModules: ["vertical-pod-autoscaler-crd", "upmeter"]
modules:
  https:
    mode: CustomCertificate
  publicDomainTemplate: "%s.example.com"
  placement: {}
discovery:
  d8SpecificNodeCountByRole:
    system: 1
    master: 1
  kubernetesVersion: 1.24.2
`

const moduleValues = `
auth: {}
argocd:
  admin:
    enabled: false
https:
  mode: CustomCertificate
internal:
  argocd:
    repositories: []
  argocdImageUpdater:
    registries: []
`
