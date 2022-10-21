/*
Copyright 2021 Flant JSC

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

package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"text/template"
)

var registryTemplate = `// Code generated by "register.go" DO NOT EDIT.
package main

import (
	_ "github.com/flant/addon-operator/sdk"
{{ range $value := . }}
	_ "{{ $value }}"
{{- end }}
)
`

func cwd() string {
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		panic("cannot get caller")
	}

	dir, err := filepath.Abs(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ { // ../../
		dir = filepath.Dir(dir)
	}

	return dir
}

func searchHooks(hookModules *[]string, dir, workDir string) error {
	files := make(map[string]interface{})

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f != nil && f.IsDir() {
			if f.Name() == "internal" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		trimDir := workDir
		moduleName := filepath.Join(
			deckhouseModuleName,
			filepath.Dir(
				strings.TrimPrefix(path, trimDir),
			),
		)
		files[moduleName] = struct{}{}
		return nil
	})

	for hook := range files {
		*hookModules = append(*hookModules, hook)
	}

	return err
}

const deckhouseModuleName = "github.com/deckhouse/deckhouse/"

func main() {
	workDir := cwd()

	var (
		output string
		stream = os.Stdout
	)
	flag.StringVar(&output, "output", "", "output file for generated code")
	flag.Parse()

	if output != "" {
		var err error
		stream, err = os.Create(output)
		if err != nil {
			panic(err)
		}

		defer stream.Close()
	}

	var hookModules []string
	if err := searchHooks(&hookModules, filepath.Join(workDir, "global-hooks"), workDir); err != nil {
		panic(err)
	}

	moduleDirs, err := filepath.Glob(filepath.Join(workDir, "modules/*/hooks"))
	if err != nil {
		panic(err)
	}

	additionalModuleDirs, err := filepath.Glob(filepath.Join(workDir, "ee/modules/*/hooks"))
	if err != nil {
		panic(err)
	}
	moduleDirs = append(moduleDirs, additionalModuleDirs...)

	additionalModuleDirs, err = filepath.Glob(filepath.Join(workDir, "ee/fe/modules/*/hooks"))
	if err != nil {
		panic(err)
	}
	moduleDirs = append(moduleDirs, additionalModuleDirs...)

	moduleDirs = append(moduleDirs, requirementCheckDirs(workDir)...)

	for _, dir := range moduleDirs {
		if err := searchHooks(&hookModules, dir, workDir); err != nil {
			panic(err)
		}
	}

	sort.Strings(hookModules)

	t := template.New("registry")
	t, err = t.Parse(registryTemplate)
	if err != nil {
		panic(err)
	}

	err = t.Execute(stream, hookModules)
	if err != nil {
		panic(err)
	}
}

func requirementCheckDirs(workDir string) []string {
	moduleDirs, err := filepath.Glob(filepath.Join(workDir, "modules/*/requirements"))
	if err != nil {
		panic(err)
	}

	additionalModuleDirs, err := filepath.Glob(filepath.Join(workDir, "ee/modules/*/requirements"))
	if err != nil {
		panic(err)
	}
	moduleDirs = append(moduleDirs, additionalModuleDirs...)

	additionalModuleDirs, err = filepath.Glob(filepath.Join(workDir, "ee/fe/modules/*/requirements"))
	if err != nil {
		panic(err)
	}
	moduleDirs = append(moduleDirs, additionalModuleDirs...)

	return moduleDirs
}
