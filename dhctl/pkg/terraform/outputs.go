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

package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

func GetMasterIPAddressForSSH(statePath string) (string, error) {
	args := []string{
		"output",
		"-no-color",
		"-json",
		fmt.Sprintf("-state=%s", statePath),
		"master_ip_address_for_ssh",
	}

	executor := &CMDExecutor{}

	result, err := executor.Output(args...)
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			err = fmt.Errorf("%s\n%v", string(ee.Stderr), err)
		}

		return "", fmt.Errorf("failed to get terraform output for 'master_ip_address_for_ssh'\n%v", err)
	}

	var output string

	err = json.Unmarshal(result, &output)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal terraform output for 'master_ip_address_for_ssh'\n%v", err)
	}

	return output, nil
}
