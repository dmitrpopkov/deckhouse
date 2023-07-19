# Copyright 2023 Flant JSC
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

function check_hostname(){
    local a=`hostname`
    bb-log-info "Check hostname: $a"
    if hostname|grep -P '^[a-z0-9]{1}(([a-z0-9\-\.]{0,61}[a-z0-9]{1})|[a-z0-9]{0,62})$' > /dev/null 2>&1; then
        bb-log-info "SUCCESS"
        exit 0
    fi
    bb-log-error "FAIL Hostname '$a' should be contain no more than 63 characters, contain only lowercase alphanumeric characters, '-' or '.', start with an alphanumeric character, end with an alphanumeric character"
    exit 1
}

check_hostname
exit $?

