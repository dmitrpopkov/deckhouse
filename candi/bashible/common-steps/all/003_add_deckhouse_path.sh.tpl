# Copyright 2021 Flant JSC
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

bb-sync-file /etc/profile.d/02-deckhouse-path.sh - << "EOF"
export PATH="/opt/deckhouse/bin:$PATH"
EOF

# On Alt Linux there is condition
# for f in /etc/profile.d/*.sh; do
#        if [ -f "$f" -a -r "$f" -a -x "$f" -a -s "$f" -a ! -L "$f" ]; then
#                . "$f"
#        fi
# done
# so /etc/profile.d/02-deckhouse-path.sh should be executable
chmod a+x /etc/profile.d/02-deckhouse-path.sh

# also need to sure that path to deckhouse bin in root bashrc
#if grep -q "^PATH=" /root/.bashrc 2>/dev/null; then
#  echo "export PATH=/opt/deckhouse/bin:$PATH" >> /root/.bashrc
#fi
if grep -q "^PATH=" /root/.bashrc 2>/dev/null; then
  sed 's/^PATH=.*//' -i /root/.bashrc
fi