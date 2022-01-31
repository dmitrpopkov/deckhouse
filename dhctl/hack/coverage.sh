#!/bin/sh

# Copyright 2022 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -Eeuo pipefail
shopt -s failglob

tmpfile=$(mktemp /tmp/coverage-report.XXXXXX)

go test -cover -coverprofile=${tmpfile} -vet=off ./pkg/...
coverage=$(go tool cover -func  ${tmpfile} | grep total | awk '{print $3}')

echo "Coverage: ${coverage}"
echo "Success!"
