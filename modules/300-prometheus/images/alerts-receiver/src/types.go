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

package main

import (
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	clusterAlertFiring       = "Firing"
	clusterAlertFiringStaled = "Firing (staled)"
)

type ClusterAlert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Alert  ClusterAlertSpec   `json:"alert,omitempty"`
	Status ClusterAlertStatus `json:"status,omitempty"`
}

type ClusterAlertStatus struct {
	AlertStatus    string `json:"alertStatus,omitempty"`
	StartsAt       string `json:"startsAt,omitempty"`
	LastUpdateTime string `json:"lastUpdateTime,omitempty"`
}

type ClusterAlertSpec struct {
	Name          string         `json:"name"`
	SeverityLevel string         `json:"severity_level,omitempty"`
	Summary       string         `json:"summary,omitempty"`
	Description   string         `json:"description,omitempty"`
	Annotations   model.LabelSet `json:"annotations,omitempty"`
	Labels        model.LabelSet `json:"labels"`
}
