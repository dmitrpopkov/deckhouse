#!/bin/bash

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

while getopts ":h:i:o:" option; do
   case $option in
      h) # display Help
         Help
         exit;;
      i)
        FILE_TEMPLATE_BINS=$OPTARG
        ;;
      o)
        RDIR=$OPTARG
        ;;
      \?)
        echo "Error: Invalid option"
        exit;;
   esac
done

Help()
{
   # Display Help
   echo "Add description of the script functions here."
   echo
   echo "Syntax: scriptTemplate [-h|i|o]"
   echo "options:"
   echo "i     Files with paths to binaryes; Support mask like /sbin/m*"
   echo "o     Output directory (Default value: '/relocate')"
   echo "h     Print this help"
   echo
}

if [[ -z $RDIR ]];then
  RDIR="/relocate"
fi

mkdir -p "${RDIR}"

function relocate() {
  local binary=$1
  relocate_item ${binary}

  for lib in $(ldd ${binary} 2>/dev/null | awk '{if ($2=="=>") print $3; else print $1}'); do
    # don't try to relocate linux-vdso.so lib due to this lib is virtual
    if [[ "${lib}" =~ "linux-vdso" ]]; then
      continue
    fi
    relocate_item ${lib}
  done
}

function relocate_item() {
  local file=$1

  local new_place="${RDIR}$(dirname ${file})"
  mkdir -p ${new_place}

  cp -a ${file} ${new_place}

  # if symlink, copy original file too
  local orig_file="$(readlink -f ${file})"
  if [[ "${file}" != "${orig_file}" ]]; then
    cp -a ${orig_file} ${new_place}
  fi
}

function get_binnary_path () {
  local bin
  BINARY_LIST=()
  
  for bin in "$@"; do
    if [[ -z $(ls -la $bin 2>/dev/null) ]]; then
      continue
    fi
    BINARY_LIST+=$(ls -la $bin 2>/dev/null | awk '{print $9}')" "
  done

  if [[ -z $BINARY_LIST ]]; then echo "No binaryes for replace"; exit 1; fi;
}

if [[ -n $FILE_TEMPLATE_BINS ]]; then
  BIN_TEMPLATE=$(cat $FILE_TEMPLATE_BINS)
  get_binnary_path ${BIN_TEMPLATE}
else
  get_binnary_path ${@}
fi

for binary in ${BINARY_LIST[@]}; do
  relocate ${binary}
done
