/*
Copyright 2024 Flant JSC

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

package conversions

import (
	"testing"

	"github.com/deckhouse/deckhouse/go_lib/deckhouse-config/conversion"
)

func Test(t *testing.T) {
	pathToConversions := "."
	cases := []struct {
		name            string
		pathToSettings  string
		pathToExpected  string
		currentVersion  int
		expectedVersion int
	}{
		{
			name:            "should convert from 1 to 2 version",
			pathToSettings:  "testdata/v1_settings.yaml",
			pathToExpected:  "testdata/v2_settings.yaml",
			currentVersion:  1,
			expectedVersion: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			converted, expected, equal, err := conversion.TestConvert(
				c.pathToSettings,
				c.pathToExpected,
				pathToConversions,
				c.currentVersion,
				c.expectedVersion)
			if err != nil {
				t.Error(err)
				return
			}
			if !equal {
				t.Errorf("Expected %v, got %v", converted, expected)
			}
		})
	}
}
