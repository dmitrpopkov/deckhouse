// Copyright 2024 Flant JSC
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

package operations

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/deckhouse/deckhouse/dhctl/pkg/app"
	"github.com/deckhouse/deckhouse/dhctl/pkg/kubernetes/client"
	"github.com/deckhouse/deckhouse/dhctl/pkg/log"
	"github.com/deckhouse/deckhouse/dhctl/pkg/util/retry"
	"github.com/deckhouse/deckhouse/dhctl/pkg/util/tomb"
)

const allowUnsafeAnnotation = "deckhouse.io/allow-unsafe"

var abstractEditing = Edit

func SecretEdit(kubeCl *client.KubernetesClient, name string, Namespace string, secret string, dataKey string) error {
	config, err := kubeCl.CoreV1().Secrets(Namespace).Get(context.TODO(), secret, metav1.GetOptions{})
	if err != nil {
		return err
	}

	configData := config.Data[dataKey]

	var modifiedData []byte
	tomb.WithoutInterruptions(func() { modifiedData, err = abstractEditing(configData) })
	if err != nil {
		return err
	}

	// This flag is validating by webhooks to allow editing unsafe resource's fields.
	if app.SanityCheck {
		addUnsafeAnnotation(config)
	}

	return log.Process(
		"common",
		fmt.Sprintf("Save %s back to the Kubernetes cluster", name), func() error {
			if string(configData) == string(modifiedData) {
				log.InfoLn("Configurations are equal. Nothing to update.")
				return nil
			}

			config.Data[dataKey] = modifiedData

			return retry.
				NewLoop(fmt.Sprintf("Update %s secret", name), 5, 5*time.Second).
				Run(func() error {
					_, err = kubeCl.CoreV1().
						Secrets(Namespace).
						Update(context.TODO(), config, metav1.UpdateOptions{})
					if err != nil {
						return err
					}

					if app.SanityCheck {
						log.InfoLn("Remove allow-unsafe annotation")
						removeUnsafeAnnotation(config)

						_, err = kubeCl.CoreV1().
							Secrets(Namespace).
							Update(context.TODO(), config, metav1.UpdateOptions{})
					}

					return err
				})
		})

}

func addUnsafeAnnotation(doc *v1.Secret) {
	if doc.Annotations == nil {
		doc.Annotations = make(map[string]string)
	}
	doc.Annotations[allowUnsafeAnnotation] = "true"
}

func removeUnsafeAnnotation(doc *v1.Secret) {
	delete(doc.Annotations, allowUnsafeAnnotation)
}