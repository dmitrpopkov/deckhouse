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

if [[ "$(getenforce)" != "Enforcing" ]]; then
  exit 0
fi

bb-event-on 'selinux_deckhouse_policy_changed' '_on_selinux_deckhouse_policy_changed'
_on_selinux_deckhouse_policy_changed() {
  checkmodule -M -m -o /var/lib/bashible/policies/deckhouse.mod /var/lib/bashible/policies/deckhouse.te
  semodule_package -o /var/lib/bashible/policies/deckhouse.pp -m /var/lib/bashible/policies/deckhouse.mod
  semodule -i /var/lib/bashible/policies/deckhouse.pp
}

mkdir -p /var/lib/bashible/policies
bb-sync-file /var/lib/bashible/policies/deckhouse.te - selinux_deckhouse_policy_changed << "EOF"
module deckhouse 1.1;
require {
	type var_lib_t;
	type init_t;
	class file execute;
}

#============= init_t ==============

#!!!! This avc is allowed in the current policy
allow init_t var_lib_t:file execute;
EOF
