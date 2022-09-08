/*
Copyright 2022 Flant JSC

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

package modules

import (
	"crypto/sha256"
	"os"
	"os/exec"
	"sync"

	"sigs.k8s.io/yaml"

	"github.com/deckhouse/deckhouse/testing/matrix/linter/rules/errors"
	"github.com/deckhouse/deckhouse/testing/matrix/linter/storage"
	"github.com/deckhouse/deckhouse/testing/matrix/linter/utils"
)

type checkResult struct {
	success bool
	errMsg  string
}

const promtoolPath = "/deckhouse/bin/promtool"

var (
	rulesCache   = map[string]checkResult{}
	rulesCacheMu = sync.RWMutex{}
)

func promtoolAvailable() bool {
	info, _ := os.Stat(promtoolPath)
	return info != nil && (info.Mode().Perm()&0111 != 0)
}

func marshalChartYaml(object storage.StoreObject) (marshal []byte, hash string, err error) {
	marshal, err = yaml.Marshal(object.Unstructured.Object["spec"])
	if err != nil {
		return
	}
	hash = string(newSHA256(marshal))
	return
}

func writeTempRuleFileFromObject(m utils.Module, marshalledYaml []byte) (path string, err error) {
	renderedFile, err := os.CreateTemp("", m.Name+".*.yml")
	if err != nil {
		return "", err
	}

	_, err = renderedFile.Write(marshalledYaml)
	if err != nil {
		return "", err
	}

	err = renderedFile.Close()
	if err != nil {
		return renderedFile.Name(), err
	}

	return renderedFile.Name(), nil
}

func checkRuleFile(path string) error {
	promtoolComand := exec.Command(promtoolPath, "check", "rules", path)
	_, err := promtoolComand.Output()
	if err != nil {
		return err
	}
	return nil
}

func newSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func createPromtoolError(m utils.Module, errMsg string) errors.LintRuleError {
	return errors.NewLintRuleError(
		"MODULE060",
		moduleLabel(m.Name),
		m.Path,
		"Promtool check failed for Helm chart:\n%s",
		errMsg,
	)
}

func PromtoolRuleCheck(m utils.Module, object storage.StoreObject) errors.LintRuleError {
	if object.Unstructured.GetKind() == "PrometheusRule" {
		if !promtoolAvailable() {
			return errors.NewLintRuleError(
				"MODULE060",
				m.Name,
				m.Path,
				"Promtool is not available. Execute `make bin/promtool` prior to starting matrix tests.",
			)
		}

		marshal, hash, err := marshalChartYaml(object)
		if err != nil {
			return errors.NewLintRuleError(
				"MODULE060",
				m.Name,
				m.Path,
				"Error marshalling Helm chart to yaml",
			)
		}

		rulesCacheMu.RLock()
		res, ok := rulesCache[hash]
		rulesCacheMu.RUnlock()
		if ok {
			if !res.success {
				return createPromtoolError(m, res.errMsg)
			}
			return errors.EmptyRuleError
		}

		path, err := writeTempRuleFileFromObject(m, marshal)
		if err != nil {
			return errors.NewLintRuleError(
				"MODULE060",
				m.Name,
				m.Path,
				"Error creating temporary rule file from Helm chart:\n%s",
				err.Error(),
			)
		}
		defer func(name string) {
			_ = os.Remove(name)
		}(path)

		err = checkRuleFile(path)
		if err != nil {
			errorMessage := string(err.(*exec.ExitError).Stderr)
			rulesCacheMu.Lock()
			rulesCache[hash] = checkResult{
				success: false,
				errMsg:  errorMessage,
			}
			rulesCacheMu.Unlock()
			return createPromtoolError(m, errorMessage)
		}
		rulesCacheMu.Lock()
		rulesCache[hash] = checkResult{success: true}
		rulesCacheMu.Unlock()
	}
	return errors.EmptyRuleError
}
