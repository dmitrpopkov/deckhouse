/*
Copyright 2021 Flant JSC

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

package hooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/flant/shell-operator/pkg/metric_storage/operation"
	"github.com/sirupsen/logrus"

	d8http "github.com/deckhouse/deckhouse/go_lib/dependency/http"
)

// This hook is needed to fill the gaps between Deckhouse restarts and avoid alerts flapping.
// it takes latest metrics from prometheus and duplicate them on Deckhouse startup

// var _ = sdk.RegisterFunc(&go_hook.HookConfig{
// 	Queue: "/modules/helm/helm_releases",
// 	OnStartup: &go_hook.OrderedConfig{
// 		Order: 1,
// 	},
// }, dependency.WithExternalDependencies(handleStartupMetrics))

func init() {
	cl := d8http.NewClient(d8http.WithInsecureSkipVerify())
	logger := logrus.New()
	logEntry := logrus.NewEntry(logger)

	// curl -s --connect-timeout 10 --max-time 10 -k -XGET -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" "https://prometheus.d8-monitoring:9090/api/v1/query?query=resource_versions_compatibility"
	promURL := "https://prometheus.d8-monitoring:9090/api/v1/query?query=resource_versions_compatibility"
	req, err := http.NewRequest("GET", promURL, nil)
	if err != nil {
		logEntry.Errorf("Prometheus request build failed: %s", err)
		return
	}
	err = d8http.SetKubeAuthToken(req)
	if err != nil {
		logEntry.Errorf("Set auth token failed: %s", err)
		return
	}

	res, err := cl.Do(req)
	if err != nil {
		logEntry.Errorf("Prometheus request failed: %s", err)
		return
	}
	defer res.Body.Close()

	d, _ := ioutil.ReadAll(res.Body)
	fmt.Println("BODY", string(d))
	var response promMetrics

	_ = json.Unmarshal(d, &response)
	// err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		logEntry.Errorf("Unmarshal failed: %s", err)
		return
	}

	fmt.Println("GOT METRICS", response)

	metricsFile := os.Getenv("METRICS_PATH")
	fmt.Println("FILE", metricsFile)

	f, err := os.OpenFile(metricsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logEntry.Errorf("Open file failed: %s", err)
		return
	}
	defer f.Close()

	for _, metricRecord := range response.Data.Result {
		if len(metricRecord.Value) < 2 {
			logEntry.Warnf("Broken metric value from prometheus: %s. Skipping", metricRecord.Value)
			continue
		}
		value, err := strconv.ParseFloat(metricRecord.Value[1].(string), 64)
		if err != nil {
			logEntry.Warnf("Failed metric convert: %s. Skipping", metricRecord.Value[1])
			continue
		}

		op := operation.MetricOperation{
			Name:  "resource_versions_compatibility",
			Value: &value,
			Labels: map[string]string{
				"helm_release_name":      metricRecord.Metric.HelmReleaseName,
				"helm_release_namespace": metricRecord.Metric.HelmReleaseNamespace,
				"k8s_version":            metricRecord.Metric.K8sVersion,
				"resource_name":          metricRecord.Metric.ResourceName,
				"resource_namespace":     metricRecord.Metric.ResourceNamespace,
				"kind":                   metricRecord.Metric.Kind,
				"api_version":            metricRecord.Metric.APIVersion,
			},
			Group:  "/modules/helm/metrics",
			Action: "set",
		}

		data, _ := json.Marshal(op)

		_, _ = f.Write(data)
	}
}

//
// func handleStartupMetrics(input *go_hook.HookInput, dc dependency.Container) error {
// 	cl := dc.GetHTTPClient(d8http.WithInsecureSkipVerify())
//
// 	// curl -s --connect-timeout 10 --max-time 10 -k -XGET -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" "https://prometheus.d8-monitoring:9090/api/v1/query?query=resource_versions_compatibility"
// 	promURL := "https://prometheus.d8-monitoring:9090/api/v1/query?query=resource_versions_compatibility"
// 	req, err := http.NewRequest("GET", promURL, nil)
// 	if err != nil {
// 		return err
// 	}
// 	err = d8http.SetKubeAuthToken(req)
// 	if err != nil {
// 		return err
// 	}
//
// 	res, err := cl.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()
//
// 	d, _ := ioutil.ReadAll(res.Body)
// 	fmt.Println("BODY", string(d))
// 	var response promMetrics
//
// 	_ = json.Unmarshal(d, &response)
// 	// err = json.NewDecoder(res.Body).Decode(&response)
// 	if err != nil {
// 		return err
// 	}
//
// 	fmt.Println("GOT METRICS", response)
//
// 	for _, metricRecord := range response.Data.Result {
// 		if len(metricRecord.Value) < 2 {
// 			input.LogEntry.Warnf("Broken metric value from prometheus: %s. Skipping", metricRecord.Value)
// 			continue
// 		}
// 		value, err := strconv.ParseFloat(metricRecord.Value[1].(string), 64)
// 		if err != nil {
// 			input.LogEntry.Warnf("Failed metric convert: %s. Skipping", metricRecord.Value[1])
// 			continue
// 		}
// 		input.MetricsCollector.Set("resource_versions_compatibility", value, map[string]string{
// 			"helm_release_name":      metricRecord.Metric.HelmReleaseName,
// 			"helm_release_namespace": metricRecord.Metric.HelmReleaseNamespace,
// 			"k8s_version":            metricRecord.Metric.K8sVersion,
// 			"resource_name":          metricRecord.Metric.ResourceName,
// 			"resource_namespace":     metricRecord.Metric.ResourceNamespace,
// 			"kind":                   metricRecord.Metric.Kind,
// 			"api_version":            metricRecord.Metric.APIVersion,
// 		})
// 	}
//
// 	return nil
// }

type promMetrics struct {
	Data struct {
		Result []struct {
			Metric struct {
				APIVersion           string `json:"api_version"`
				Kind                 string `json:"kind"`
				HelmReleaseName      string `json:"helm_release_name"`
				HelmReleaseNamespace string `json:"helm_release_namespace"`
				K8sVersion           string `json:"k8s_version"`
				ResourceName         string `json:"resource_name"`
				ResourceNamespace    string `json:"resource_namespace"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
