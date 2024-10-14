// Copyright 2022 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/utils"
)

var semVerRegex = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`)

// Now:
//
// list releases
// list releases --all
// get release alpha
// get release alpha --all
// list sources
// list modules deckhouse-prod
// list module-release deckhouse-prod console
// list module-release deckhouse-prod console --all
// get module-release deckhouse-prod console alpha
// get module-release deckhouse-prod console alpha --all

// Proposal:
//
// get releases (test alias release or releases)
// get releases --all
// get releases -c auto (default channel, print name, if empty - stable)
// get releases -c alpha --all
// get sources
// get modules deckhouse-prod
// get modules deckhouse-prod console
// get modules deckhouse-prod console --all
// get modules deckhouse-prod console -c alpha
// get modules deckhouse-prod console -c alpha --all

func DefineRegistryCommand(kpApp *kingpin.Application) {
	registryCmd := kpApp.Command("registry", "Deckhouse repository work.").
		PreAction(func(_ *kingpin.ParseContext) error {
			kpApp.UsageTemplate(kingpin.DefaultUsageTemplate)
			return nil
		})

	registryGetCmd := registryCmd.Command("get", "get from registry").
		PreAction(func(_ *kingpin.ParseContext) error {
			kpApp.UsageTemplate(kingpin.DefaultUsageTemplate)
			return nil
		})

	registerReleaseCommand(registryGetCmd)
	registerSourceCommand(registryGetCmd)
	registerModuleCommand(registryGetCmd)
}

func registerReleaseCommand(parentCMD *kingpin.CmdClause) {
	releasesCmd := parentCMD.Command("releases", "Release resource. Aliases: 'release','rel'").
		Alias("release").Alias("rel")

	releaseChannel := releasesCmd.Flag("channel", "Release channel."+
		" If release is 'auto' - using default channel from configuration."+
		" If there is not default channel in configuration - use 'stable'").Short('c').String()
	allFlag := releasesCmd.Flag("all", "Output without restrictions.").Bool()
	releasesCmd.Action(func(_ *kingpin.ParseContext) error {
		ctx := context.TODO()

		registry, channel, rconf, err := getDeckhouseRegistry(ctx)
		if err != nil {
			return fmt.Errorf("get deckhouse registry: %w", err)
		}

		svc := NewDeckhouseService(registry, rconf)

		if *releaseChannel != "" {
			if *releaseChannel != "auto" {
				channel = *releaseChannel
			}

			return handleGetDeckhouseRelease(svc, channel, *allFlag)
		}

		return handleListDeckhouseReleases(ctx, svc, *allFlag)
	})
}

func handleListDeckhouseReleases(ctx context.Context, svc *DeckhouseService, all bool) error {
	ls, err := svc.ListDeckhouseReleases(ctx)
	if err != nil {
		return fmt.Errorf("list deckhouse releases: %w", err)
	}

	// if we need full tags list, not only semVer
	if !all {
		res := make([]string, 0, 1)
		for _, v := range ls {
			if semVerRegex.MatchString(v) {
				res = append(res, v)
			}
		}

		ls = res
	}

	if len(ls) == 0 {
		fmt.Println()

		if all {
			fmt.Println("Releases not found")
		} else {
			fmt.Println("Releases with semVer not found. Use --all argument to watch all releases in the registry")
		}

		return nil
	}

	fmt.Println()
	fmt.Println(strings.Join(ls, "\n"))

	return nil
}

func handleGetDeckhouseRelease(svc *DeckhouseService, channel string, all bool) error {
	meta, err := svc.GetDeckhouseRelease(channel)
	if err != nil && !errors.Is(err, ErrChannelIsNotFound) {
		return fmt.Errorf("get deckhouse release: %w", err)
	}

	if err != nil {
		return fmt.Errorf("deckhouse release with channel '%s' is not found", channel)
	}

	if !all {
		fmt.Println()
		fmt.Printf("Deckhouse version in channel '%s': %s\n", channel, meta.Version)
		fmt.Println()

		return nil
	}

	b, err := json.MarshalIndent(meta, "", "    ")
	if err != nil {
		return fmt.Errorf("marshall indent: %w", err)
	}

	fmt.Println()
	fmt.Printf("%s", b)
	fmt.Println()

	return nil
}

func registerSourceCommand(parentCMD *kingpin.CmdClause) {
	sourcesCmd := parentCMD.Command("sources", "Source resources. Aliases: 'source','src'").
		Alias("source").Alias("src")
	sourcesCmd.Action(func(_ *kingpin.ParseContext) error {
		ctx := context.TODO()

		k8sClient, err := newKubernetesClient()
		if err != nil {
			panic(err)
		}

		msl := new(v1alpha1.ModuleSourceList)
		if err := k8sClient.List(ctx, msl); err != nil {
			return fmt.Errorf("list ModuleSource: %w", err)
		}

		srcs := make([]string, 0, len(msl.Items))
		for _, ms := range msl.Items {
			srcs = append(srcs, ms.GetName())
		}

		fmt.Println()
		fmt.Printf("Module sources found (%d):\n\n", len(srcs))

		for _, src := range srcs {
			fmt.Printf("%s\n", src)
		}

		return nil
	})
}

func registerModuleCommand(parentCMD *kingpin.CmdClause) {
	// deckhouse-controller registry list modules <module-source>
	modulesCmd := parentCMD.Command("modules", "Show modules list. Aliases: 'module','mod'").
		Alias("module").Alias("mod")
	moduleSource := modulesCmd.Arg("module-source", "Module source name.").Required().String()
	moduleName := modulesCmd.Arg("module-name", "Module name.").String()
	moduleChannel := modulesCmd.Flag("channel", "Module name.").Short('c').String()
	allFlag := modulesCmd.Flag("all", "Complete list of tags.").Bool()

	modulesCmd.Action(func(_ *kingpin.ParseContext) error {
		ctx := context.TODO()

		registry, rconf, err := getModuleRegistry(ctx, *moduleSource)
		if err != nil {
			return fmt.Errorf("get module registry: %w", err)
		}

		svc := NewModuleService(registry, rconf)

		if *moduleName != "" {
			if *moduleChannel != "" {
				return handleGetModuleInfoInChannel(svc, *moduleName, *moduleChannel, *allFlag)
			}

			return handleListModulesVersions(ctx, svc, *moduleName, *allFlag)
		}

		return handleListModulesNames(ctx, svc, *allFlag)
	})
}

func handleGetModuleInfoInChannel(svc *ModuleService, name string, channel string, all bool) error {
	meta, err := svc.GetModuleRelease(name, channel)
	if err != nil && !errors.Is(err, ErrChannelIsNotFound) {
		return fmt.Errorf("get module release %s: %w", name, err)
	}

	if err != nil {
		return fmt.Errorf("module release with name '%s' and channel '%s' is not found", name, channel)
	}

	if !all {
		fmt.Println()
		fmt.Printf("Module version in channel '%s': %s", channel, meta.Version)
		fmt.Println()

		return nil
	}

	b, err := json.MarshalIndent(meta, "", "    ")
	if err != nil {
		return fmt.Errorf("marshall indent: %w", err)
	}

	fmt.Println()
	fmt.Printf("%s", b)
	fmt.Println()

	return nil
}

func handleListModulesVersions(ctx context.Context, svc *ModuleService, name string, all bool) error {
	ls, err := svc.ListModuleTags(ctx, name)
	if err != nil && !errors.Is(err, ErrModuleIsNotFound) {
		return fmt.Errorf("list module tags: %w", err)
	}

	if err != nil {
		return fmt.Errorf("module release with name '%s' is not found", name)
	}

	// if we need full tags list, not only semVer
	if !all {
		res := make([]string, 0, 1)
		for _, v := range ls {
			if semVerRegex.MatchString(v) {
				res = append(res, v)
			}
		}

		ls = res
	}

	if len(ls) == 0 {
		fmt.Println()

		if all {
			fmt.Println("Module releases not found")
		} else {
			fmt.Println("Module releases with semVer not found. Use --all argument to watch all releases in the registry")
		}

		return nil
	}

	fmt.Println()
	fmt.Println(strings.Join(ls, "\n"))

	return nil
}

func handleListModulesNames(ctx context.Context, svc *ModuleService, all bool) error {
	modules, err := svc.ListModules(ctx)
	if err != nil {
		return fmt.Errorf("list modules: %w", err)
	}

	if len(modules) == 0 {
		fmt.Println()

		if all {
			fmt.Println("Modules not found")
		} else {
			fmt.Println("Modules with semVer not found. Use --all argument to watch all releases in the registry")
		}

		fmt.Println()

		return nil
	}

	fmt.Println()
	fmt.Printf("Modules found (%d):", len(modules))
	fmt.Println()

	fmt.Println()
	fmt.Println(strings.Join(modules, "\n"))

	return nil
}

func newKubernetesClient() (client.Client, error) {
	scheme := runtime.NewScheme()

	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))

	restConfig := ctrl.GetConfigOrDie()
	opts := client.Options{
		Scheme: scheme,
	}

	k8sClient, err := client.New(restConfig, opts)
	if err != nil {
		return nil, fmt.Errorf("create kubernetes client: %w", err)
	}

	return k8sClient, nil
}

func getDeckhouseRegistry(ctx context.Context) (string, string, *utils.RegistryConfig, error) {
	k8sClient, err := newKubernetesClient()
	if err != nil {
		panic(err)
	}

	secret := new(corev1.Secret)
	if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: "d8-system", Name: "deckhouse-registry"}, secret); err != nil {
		return "", "", nil, fmt.Errorf("list ModuleSource got an error: %w", err)
	}

	drs, _ := utils.ParseDeckhouseRegistrySecret(secret.Data)

	var discoverySecret corev1.Secret
	key := types.NamespacedName{Namespace: "d8-system", Name: "deckhouse-discovery"}
	if err := k8sClient.Get(ctx, key, &discoverySecret); err != nil {
		return "", "", nil, fmt.Errorf("get deckhouse discovery sectret got an error: %w", err)
	}

	clusterUUID, ok := discoverySecret.Data["clusterUUID"]
	if !ok {
		return "", "", nil, fmt.Errorf("not found clusterUUID in discovery secret: %w", err)
	}

	releaseChannel := string(discoverySecret.Data["releaseChannel"])

	rconf := &utils.RegistryConfig{
		DockerConfig: drs.DockerConfig,
		Scheme:       drs.Scheme,
		UserAgent:    string(clusterUUID),
	}

	return drs.ImageRegistry, releaseChannel, rconf, nil
}

func getModuleRegistry(ctx context.Context, moduleSource string) (string, *utils.RegistryConfig, error) {
	k8sClient, err := newKubernetesClient()
	if err != nil {
		panic(err)
	}

	ms := new(v1alpha1.ModuleSource)
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: moduleSource}, ms); err != nil {
		return "", nil, fmt.Errorf("get ModuleSource %s got an error: %w", moduleSource, err)
	}

	rconf := &utils.RegistryConfig{
		DockerConfig: ms.Spec.Registry.DockerCFG,
		Scheme:       ms.Spec.Registry.Scheme,
		CA:           ms.Spec.Registry.CA,
		UserAgent:    "deckhouse-controller/ModuleControllers",
	}

	return ms.Spec.Registry.Repo, rconf, nil
}