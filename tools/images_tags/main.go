/*
Copyright 2022 Flant JSC

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
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

var imagesTagsTemplate = `// Code generated by "tools/images_tags.go" DO NOT EDIT.
// To generate run 'make generate'
package library

var	DefaultImagesTags = map[string]map[string]string{
{{- range $module, $images := . }}
	"{{ $module }}": {
  {{- range $image, $tag := $images }}
	  "{{ $image }}": "{{ $tag }}",
  {{- end }}
	},
{{- end }}
}
`

func main() {
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

	tags := make(map[string]map[string]string)

	cmd := exec.Command("werf", "config", "render")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "CI_COMMIT_REF_NAME=", "CI_COMMIT_TAG=", "WERF_ENV=FE")
	cmd.Dir = path.Join("..")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	a := strings.NewReader(string(out))
	d := yaml.NewDecoder(a)
	for {
		spec := new(Spec)
		err := d.Decode(&spec)
		if spec == nil {
			continue
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		if spec.Image == "images-tags" {
			for _, dependency := range spec.Dependencies {
				for _, i := range dependency.Imports {
					s := strings.Split(i.TargetEnv, "_")
					module := s[3]
					tag := s[4]
					if tags[module] == nil {
						tags[module] = make(map[string]string)
					}
					tags[module][tag] = fmt.Sprintf("imageHash-%s-%s", module, tag)
				}
			}
			break
		}
	}

	var buf bytes.Buffer
	t := template.New("imagesTags")
	t, err = t.Parse(imagesTagsTemplate)
	if err != nil {
		panic(err)
	}
	err = t.Execute(&buf, tags)
	if err != nil {
		panic(err)
	}
	p, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	stream.Write(p)

}

type Spec struct {
	Image        string       `yaml:"image"`
	Dependencies []Dependency `yaml:"dependencies"`
}

type Dependency struct {
	Imports []Import `yaml:"imports"`
}

type Import struct {
	TargetEnv string `yaml:"targetEnv"`
}
