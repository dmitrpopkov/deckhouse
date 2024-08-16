/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package template

import (
	"context"
	"controller/pkg/validate"
	"fmt"
	"net/http"

	"controller/pkg/apis/deckhouse.io/v1alpha1"
	"controller/pkg/apis/deckhouse.io/v1alpha2"
	"controller/pkg/helm"

	admissionv1 "k8s.io/api/admission/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/yaml"
)

const DeckhouseHeritage = "deckhouse"

func Register(runtimeManager manager.Manager, serviceAccount string) {
	hook := &webhook.Admission{Handler: &validator{client: runtimeManager.GetClient(), serviceAccount: serviceAccount}}
	runtimeManager.GetWebhookServer().Register("/validate/v1alpha1/templates", hook)
}

type validator struct {
	serviceAccount string
	client         client.Client
}

func (v *validator) Handle(_ context.Context, req admission.Request) admission.Response {
	template := new(v1alpha1.ProjectTemplate)
	if req.Operation == admissionv1.Create || req.Operation == admissionv1.Update {
		if err := yaml.Unmarshal(req.Object.Raw, template); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		// cannot create/update template with deckhouse heritage
		if msg := v.validateHeritage(req, template); msg != "" {
			return admission.Denied(msg)
		}

		// cannot create/update invalid template
		if err := validate.ProjectTemplate(template); err != nil {
			return admission.Errored(http.StatusBadRequest, fmt.Errorf("project template validation failed: %v", err))
		}
	}
	if req.Operation == admissionv1.Delete {
		if err := yaml.Unmarshal(req.OldObject.Raw, template); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}

		// cannot delete template with deckhouse heritage
		if msg := v.validateHeritage(req, template); msg != "" {
			return admission.Denied(msg)
		}

		// cannot delete template if it is used
		projects := new(v1alpha2.ProjectList)
		if err := v.client.List(context.Background(), projects, client.MatchingLabels{helm.ProjectTemplateLabel: template.Name}); err != nil {
			return admission.Errored(http.StatusInternalServerError, err)
		}
		if len(projects.Items) > 0 {
			msg := fmt.Sprintf("The '%s' project template cannot be deleted, it is used in the '%s' project", template.Name, projects.Items[0].Name)
			return admission.Denied(msg)
		}
	}
	return admission.Allowed("")
}

func (v *validator) validateHeritage(req admission.Request, template *v1alpha1.ProjectTemplate) string {
	if template.Labels != nil {
		heritage, ok := template.Labels[helm.HeritageLabel]
		if ok && heritage == DeckhouseHeritage && req.UserInfo.Username != v.serviceAccount {
			return fmt.Sprintf("The '%s' project template has the 'heritage' label with forbidden value: 'deckhouse'", template.Name)
		}
	}
	return ""
}
