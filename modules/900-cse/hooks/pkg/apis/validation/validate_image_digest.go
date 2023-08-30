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

package validation

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	log "github.com/sirupsen/logrus"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	"github.com/slok/kubewebhook/v2/pkg/model"
	kwhvalidating "github.com/slok/kubewebhook/v2/pkg/webhook/validating"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type validationHandler struct {
	logger            *log.Entry
	registryTransport *http.Transport
}

func NewValidationHandler(skipVerify bool) *validationHandler {
	logger := log.WithField("prefix", "image-digest-validation")
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if skipVerify {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &validationHandler{
		logger:            logger,
		registryTransport: customTransport,
	}
}

func (vh *validationHandler) imageDigestValidationHandler() http.Handler {
	vf := kwhvalidating.ValidatorFunc(func(ctx context.Context, review *model.AdmissionReview, obj metav1.Object) (result *kwhvalidating.ValidatorResult, err error) {
		pod, ok := obj.(*corev1.Pod)
		if !ok {
			return rejectResult("incorrect pod data")
		}
		for _, image := range vh.GetImagesFromPod(pod) {
			err := vh.CheckImageDigest(image)
			if err != nil {
				return rejectResult(err.Error())
			}
		}
		return allowResult("all images is correct")
	})

	// Create webhook.
	wh, _ := kwhvalidating.NewWebhook(kwhvalidating.WebhookConfig{
		ID:        "image-digest-validation",
		Validator: vf,
		Logger:    validationLogger,
		Obj:       &corev1.Pod{},
	})

	return kwhhttp.MustHandlerFor(kwhhttp.HandlerConfig{Webhook: wh, Logger: validationLogger})
}

func (vh *validationHandler) GetImagesFromPod(pod *corev1.Pod) []string {
	images := []string{}
	for _, container := range pod.Spec.Containers {
		images = append(images, container.Image)
	}
	return images
}

func (vh *validationHandler) CheckImageDigest(imageName string) error {
	ref, err := vh.ParseImageName(imageName)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	image, err := remote.Image(
		ref,
		remote.WithTransport(vh.registryTransport),
		remote.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	imageDigest, err := image.Digest()
	if err != nil {
		return err
	}
	vh.logger.WithField("imageDigest", imageDigest.String()).WithField("imageName", imageName).Info("image from remote")
	return nil
}

func (vh *validationHandler) ParseImageName(image string) (name.Reference, error) {
	return name.ParseReference(image)
}
