#!/bin/bash

# Copyright 2021 Flant CJSC
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

set -Eeo pipefail

# For local run on macos.
grepcmd=grep
if [[ "$OSTYPE" == "darwin"* ]]; then
  grepcmd=ggrep
fi

EE_LICENSE='(?s)Copyright 2021 Flant CJSC.*Licensed under the Deckhouse Platform Enterprise Edition \(EE\) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE'
CE_LICENSE="(?s)\
[/#{-]*(\s)*Copyright 2021 Flant CJSC[-}\n#/]*\
[/#{-]*(\s)*Licensed under the Apache License, Version 2.0 \(the \"License\"\);[-}\n]*\
[/#{-]*(\s)*you may not use this file except in compliance with the License.[-}\n]*\
[/#{-]*(\s)*You may obtain a copy of the License at[-}\n#/]*\
[/#{-]*(\s)*http://www.apache.org/licenses/LICENSE-2.0[-}\n#/]*\
[/#{-]*(\s)*Unless required by applicable law or agreed to in writing, software[-}\n]*\
[/#{-]*(\s)*distributed under the License is distributed on an \"AS IS\" BASIS,[-}\n]*\
[/#{-]*(\s)*WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.[-}\n]*\
[/#{-]*(\s)*See the License for the specific language governing permissions and[-}\n]*\
[/#{-]*(\s)*limitations under the License.[-}\n]*"

function check_flant_license() {
  echo -n " check flant license "
  filename=${1}
  if [[ $filename =~ ^/?ee/ ]]; then
    if ! $grepcmd -Pzq "$EE_LICENSE" $filename; then
      echo -e "ERROR\n  $filename doesn't contain EE license"
      return 1
    fi
  else
    if ! $grepcmd -Pzq "$CE_LICENSE" $filename; then
      echo -e "ERROR\n  $filename doesn't contain CE license"
      return 1
    fi
  fi
}

function check_file_copyright() {
  filename=${1}

  # Skip if is not a file or is a symlink.
  if [[ ! -f $filename ]] || [[ -L $filename ]]; then return 0; fi

  # Skip autogenerated file or file already has other than Flant copyright
  if [[ -n $(head -n 10 $filename | grep -E 'Copyright The|autogenerated|DO NOT EDIT') ]]; then return 0; fi

  if [[ -n $(head -n 10 $filename | grep -E 'Flant|Deckhouse') ]]; then
    if ! check_flant_license "${filename}"; then
      return 1
    fi
  elif [[ -n $(head -n 10 $filename | grep -E 'Copyright') ]]; then
    # Skip file with some other copyright
    return 0
  else
    echo -e "ERROR\n  $filename doesn't contain any copyright information."
    return 1
  fi

}

DIFF_DATA="$(git diff origin/main --name-only -w --ignore-blank-lines --diff-filter=A)"

if [[ -z $DIFF_DATA ]]; then
  echo "Empty diff data"
  exit 0
fi

hasErrors=0
for item in ${DIFF_DATA}; do
  echo -n "Check ${item} ... "
  pattern="\.go$|/[^/.]+$|\.sh$|\.lua$|\.py$"
  skip_pattern="geohash.lua$|Dockerfile$|Makefile$|/docs/documentation/|/docs/site/|bashrc$|inputrc$"
  if [[ "$item" =~ $pattern ]] && ! [[ "$item" =~ $skip_pattern ]]; then
    if ! check_file_copyright "${item}"; then
      hasErrors=1
    else
      echo "OK"
    fi
  else
    echo "Skip"
  fi
done

if [[ $hasErrors == 0 ]]; then
  echo "Validation successful."
fi

exit $hasErrors
