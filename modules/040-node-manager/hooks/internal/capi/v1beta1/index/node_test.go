/*
Copyright 2021 The Kubernetes Authors.

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

package index

import (
	"testing"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestIndexNodeByProviderID(t *testing.T) {
	validProviderID := "aws://region/zone/id"

	testCases := []struct {
		name     string
		object   client.Object
		expected []string
	}{
		{
			name:     "Node has no providerID",
			object:   &corev1.Node{},
			expected: nil,
		},
		{
			name: "Node has invalid providerID",
			object: &corev1.Node{
				Spec: corev1.NodeSpec{
					ProviderID: "",
				},
			},
			expected: nil,
		},
		{
			name: "Node has valid providerID",
			object: &corev1.Node{
				Spec: corev1.NodeSpec{
					ProviderID: validProviderID,
				},
			},
			expected: []string{validProviderID},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			got := NodeByProviderID(tc.object)
			g.Expect(got).To(BeEquivalentTo(tc.expected))
		})
	}
}
