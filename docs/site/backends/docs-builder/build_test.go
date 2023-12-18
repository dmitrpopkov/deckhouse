// Copyright 2023 Flant JSC
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

package main

import "testing"

func TestAssembleErrorRegexp(t *testing.T) {
	input := "error building site: assemble: \"/app/hugo/content/modules/moduleName/BROKEN.md:1:1\": EOF looking for end YAML front matter delimiter"

	path, ok := getAssembleErrorPath(input)
	if !ok || path != "/app/hugo/content/modules/moduleName/BROKEN.md" {
		t.Fatalf("unedxpcted path %q", path)
	}
}

func TestAssembleErrorWithColorRegexp(t *testing.T) {
	input := "error building site: assemble: \x1b[1;36m\"/app/hugo/content/modules/moduleName/BROKEN.md:1:1\"\x1b[0m: EOF looking for end YAML front matter delimiter"

	path, ok := getAssembleErrorPath(input)
	if !ok || path != "/app/hugo/content/modules/moduleName/BROKEN.md" {
		t.Fatalf("unedxpcted path %q", path)
	}
}

func TestGetModulePath(t *testing.T) {
	var tests = []struct {
		filePath string
		expected string
	}{
		{
			filePath: "/app/hugo/content/modules/moduleName/BROKEN.md",
			expected: "/app/hugo/content/modules/moduleName",
		},
		{
			filePath: "/app/hugo/content/modules/moduleName/alpha/BROKEN.md",
			expected: "/app/hugo/content/modules/moduleName/alpha",
		},
	}

	for _, test := range tests {
		t.Run(test.filePath, func(t *testing.T) {
			got := getModulePath(test.filePath)
			if got != test.expected {
				t.Error("unexpected result", got)
			}
		})
	}
}
