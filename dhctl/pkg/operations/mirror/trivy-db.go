// Copyright 2023 Flant JSC
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

package mirror

import (
	"fmt"
	"path"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/layout"
)

func PullTrivyVulnerabilityDatabaseImageToLayout(
	sourceRepo string,
	authProvider authn.Authenticator,
	targetLayout layout.Path,
	insecure, skipVerifyTLS bool,
) error {
	nameOpts, _ := MakeRemoteRegistryRequestOptions(authProvider, insecure, skipVerifyTLS)

	trivyDBPath := path.Join(sourceRepo, "security", "trivy-db:2")
	ref, err := name.ParseReference(trivyDBPath, nameOpts...)
	if err != nil {
		return fmt.Errorf("parse trivy-db reference: %w", err)
	}

	images := map[string]struct{}{ref.String(): {}}
	if err = PullImageSet(authProvider, targetLayout, images, insecure, skipVerifyTLS, false); err != nil {
		return fmt.Errorf("pull vulnerability database: %w", err)
	}

	return nil
}
