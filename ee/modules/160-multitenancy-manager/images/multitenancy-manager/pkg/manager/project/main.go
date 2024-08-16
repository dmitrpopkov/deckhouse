/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package project

import (
	"context"
	"fmt"
	"sync"
	"time"

	"controller/pkg/apis/deckhouse.io/v1alpha2"
	"controller/pkg/helm"
	"controller/pkg/validate"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/util/wait"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
)

const (
	Finalizer = "projects.deckhouse.io/project-exists"
)

type Interface interface {
	Init(ctx context.Context, checker healthz.Checker, init *sync.WaitGroup) error
	Handle(ctx context.Context, project *v1alpha2.Project) (ctrl.Result, error)
	Delete(ctx context.Context, project *v1alpha2.Project) (ctrl.Result, error)
}
type manager struct {
	log        logr.Logger
	client     client.Client
	helmClient helm.Interface
}

func New(client client.Client, helmClient helm.Interface, log logr.Logger) Interface {
	return &manager{
		log:        log.WithName("project-manager"),
		client:     client,
		helmClient: helmClient,
	}
}

func (m *manager) Init(ctx context.Context, checker healthz.Checker, init *sync.WaitGroup) error {
	defer init.Done()

	m.log.Info("waiting for webhook server starting")
	check := func(ctx context.Context) (bool, error) {
		if err := checker(nil); err != nil {
			m.log.Info("webhook server not startup yet")
			return false, nil
		}
		return true, nil
	}

	if err := wait.PollUntilContextTimeout(ctx, time.Second, 10*time.Second, true, check); err != nil {
		m.log.Error(err, "webhook server failed to start")
		return fmt.Errorf("webhook server failed to start: %w", err)
	}

	// to make sure that the server is started, without working server reconcile is failed
	if err := wait.PollUntilContextTimeout(ctx, time.Second, 10*time.Second, false, check); err != nil {
		m.log.Error(err, "webhook server failed to start")
		return fmt.Errorf("webhook server failed to start: %w", err)
	}

	m.log.Info("webhook server started")
	return nil
}

// Handle ensures project`s resources
func (m *manager) Handle(ctx context.Context, project *v1alpha2.Project) (ctrl.Result, error) {
	// set deploying status
	if err := m.updateProjectStatus(ctx, project, v1alpha2.ProjectStateDeploying, nil); err != nil {
		m.log.Error(err, "failed to set project status")
		return ctrl.Result{Requeue: true}, nil
	}

	// set template label and delete sync require annotation
	m.log.Info("preparing the project", "project", project.Name, "projectTemplate", project.Spec.ProjectTemplateName)
	if err := m.prepareProject(ctx, project); err != nil {
		m.log.Error(err, "failed to prepare project")
		return ctrl.Result{Requeue: true}, nil
	}

	// get a project template for the project
	m.log.Info("getting project template for project", "project", project.Name, "projectTemplate", project.Spec.ProjectTemplateName)
	projectTemplate, err := m.projectTemplateByName(ctx, project.Spec.ProjectTemplateName)
	if err != nil {
		m.log.Error(err, "failed to get project template", "project", project.Name, "projectTemplate", project.Spec.ProjectTemplateName)
		cond := m.makeCondition(v1alpha2.ConditionTypeProjectTemplateFound, v1alpha2.ConditionTypeFalse, err.Error())
		if statusErr := m.updateProjectStatus(ctx, project, v1alpha2.ProjectStateError, cond); statusErr != nil {
			m.log.Error(statusErr, "failed to set project status", "project", project.Name, "projectTemplate", project.Spec.ProjectTemplateName)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}
	// check if the project template exists
	if projectTemplate == nil {
		m.log.Info("the project template not found for the project", "project", project.Name, "projectTemplate", project.Spec.ProjectTemplateName)
		cond := m.makeCondition(v1alpha2.ConditionTypeProjectTemplateFound, v1alpha2.ConditionTypeFalse, "The project template not found")
		if statusErr := m.updateProjectStatus(ctx, project, v1alpha2.ProjectStateError, cond); statusErr != nil {
			m.log.Error(statusErr, "failed to set the project status", "project", project.Name, "projectTemplate", project.Spec.ProjectTemplateName)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}

	// update conditions
	cond := m.makeCondition(v1alpha2.ConditionTypeProjectTemplateFound, v1alpha2.ConditionTypeTrue, "")
	if statusErr := m.updateProjectStatus(ctx, project, "", cond); statusErr != nil {
		m.log.Error(statusErr, "failed to update the project status", "project", project.Name, "projectTemplate", project.Spec.ProjectTemplateName)
		return ctrl.Result{Requeue: true}, nil
	}

	// validate the project against the project template
	m.log.Info("validating the project spec", "project", project.Name, "projectTemplate", projectTemplate.Name)
	if err = validate.Project(project, projectTemplate); err != nil {
		m.log.Error(err, "failed to validate the project spec", "project", project.Name, "projectTemplate", projectTemplate.Name)
		m.makeCondition(v1alpha2.ConditionTypeProjectValidated, v1alpha2.ConditionTypeFalse, err.Error())
		if statusErr := m.updateProjectStatus(ctx, project, v1alpha2.ProjectStateError, cond); statusErr != nil {
			m.log.Error(statusErr, "failed to set the project status", "project", project.Name, "projectTemplate", projectTemplate.Name)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}

	// update conditions
	cond = m.makeCondition(v1alpha2.ConditionTypeProjectValidated, v1alpha2.ConditionTypeTrue, "")
	if statusErr := m.updateProjectStatus(ctx, project, "", cond); statusErr != nil {
		m.log.Error(statusErr, "failed to update the project status", "project", project.Name, "projectTemplate", project.Spec.ProjectTemplateName)
		return ctrl.Result{Requeue: true}, nil
	}

	// upgrade project`s resources
	m.log.Info("upgrading resources for the project", "project", project.Name, "projectTemplate", projectTemplate.Name)
	if err = m.helmClient.Upgrade(ctx, project, projectTemplate); err != nil {
		m.log.Error(err, "failed to upgrade resources for the project", "project", project.Name, "projectTemplate", projectTemplate.Name)
		cond = m.makeCondition(v1alpha2.ConditionTypeProjectResourcesUpgraded, v1alpha2.ConditionTypeFalse, err.Error())
		if statusErr := m.updateProjectStatus(ctx, project, v1alpha2.ProjectStateError, cond); statusErr != nil {
			m.log.Error(statusErr, "failed to set the project status", "project", project.Name, "projectTemplate", projectTemplate.Name)
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, nil
	}

	// update conditions
	cond = m.makeCondition(v1alpha2.ConditionTypeProjectResourcesUpgraded, v1alpha2.ConditionTypeTrue, "")

	// set deployed status
	m.log.Info("setting deployed status for the project", "project", project.Name, "projectTemplate", projectTemplate.Name)
	if err = m.updateProjectStatus(ctx, project, v1alpha2.ProjectStateDeployed, cond); err != nil {
		m.log.Error(err, "failed to set the project status", "project", project.Name, "projectTemplate", projectTemplate.Name)
		return ctrl.Result{Requeue: true}, nil
	}

	// set finalizer
	m.log.Info("setting finalizer for the project", "project", project.Name, "projectTemplate", projectTemplate.Name)
	if err = m.setFinalizer(ctx, project); err != nil {
		m.log.Error(err, "failed to set the project finalizer")
		return ctrl.Result{Requeue: true}, nil
	}

	m.log.Info("the project reconciled", "project", project.Name, "projectTemplate", projectTemplate.Name)
	return ctrl.Result{}, nil
}

// Delete deletes project`s resources
func (m *manager) Delete(ctx context.Context, project *v1alpha2.Project) (ctrl.Result, error) {
	// delete resources
	if err := m.helmClient.Delete(ctx, project.Name); err != nil {
		m.log.Error(err, "failed to delete project", "project", project.Name)
		return ctrl.Result{Requeue: true}, nil
	}

	// remove finalizer
	if err := m.removeFinalizer(ctx, project); err != nil {
		m.log.Error(err, "failed to remove finalizer from project", "project", project.Name)
		return ctrl.Result{Requeue: true}, nil
	}

	m.log.Info("successfully deleted project", "project", project.Name)
	return ctrl.Result{}, nil
}
