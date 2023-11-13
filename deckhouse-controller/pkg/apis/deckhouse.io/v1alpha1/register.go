// Copyright 2021 Flant JSC
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

package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	deckhouse_io "github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: deckhouse_io.GroupName, Version: "v1alpha1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// localSchemeBuilder and AddToScheme will stay in k8s.io/kubernetes.
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addKnownTypes)
}

// Adds the list of known types to api.Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ModuleConfig{},
		&ModuleConfigList{},
		&Module{},
		&ModuleList{},
		&ModuleSource{},
		&ModuleSourceList{},
		&ModuleRelease{},
		&ModuleReleaseList{},
	)
	err := scheme.AddFieldLabelConversionFunc(ModuleReleaseGVK, func(label, value string) (internalLabel, internalValue string, err error) {
		switch label {
		case "status.phase":
			return label, value, nil

		default:
			return "", "", fmt.Errorf("field label not supported: %s", label)
		}
	})
	if err != nil {
		return err
	}

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
