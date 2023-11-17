#!/bin/bash

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

export TIMEOUT_SEC=15
export DISKLESS_STORAGE_POOL="DfltDisklessStorPool"
export LINSTOR_NAMESPACE="d8-linstor"

command -v jq >/dev/null 2>&1 || { echo "jq is required but it's not installed.  Aborting." >&2; exit 1; }

exec_linstor_with_exit_code_check() {
  execute_command kubectl -n ${LINSTOR_NAMESPACE} exec -ti deploy/linstor-controller -c linstor-controller -- linstor "$@"
}

linstor() {
  kubectl -n ${LINSTOR_NAMESPACE} exec -ti deploy/linstor-controller -c linstor-controller -- linstor "$@"
}

execute_command() {
  while true; do
    eval "$@"
    local exit_code=$?
    if [[ $exit_code -eq 0 ]]; then
        break
    else
      echo "Command \"$@\" failed with exit code \"${exit_code}\""
      if get_user_confirmation "Would you like to retry?" "y" "n"; then
        sleep 2
        continue
      else
        if get_user_confirmation "Ignore the error and continue?" "y" "n"; then  
          break
        else
          exit_function
        fi
      fi
    fi
  done
}

linstor_check_faulty() {
  while true; do
    echo "Checking for LINSTOR controller online"
    local count=0
    local max_attempts=30
    until linstor node list > /dev/null 2>&1 || [ $count -eq $max_attempts ]; do
      echo "LINSTOR controller is not online. Waiting $TIMEOUT_SEC seconds and rechecking for LINSTOR controller online. Attempt $((count+1))/$max_attempts."
      sleep $TIMEOUT_SEC
      ((count++))
    done

    if [ $count -eq $max_attempts ]; then
      echo "Timeout reached. LINSTOR controller is not online."
      if get_user_confirmation "Perform recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
        echo "Waiting $TIMEOUT_SEC seconds and rechecking for LINSTOR controller online"
        sleep $TIMEOUT_SEC
      else
        exit_function
      fi
      continue
    fi

    echo "LINSTOR controller is online"
    echo "Checking for faulty resources"
    if (( $(linstor resource list --faulty | tee /dev/tty | grep -v -i sync | grep  "[a-zA-Z0-9]" | wc -l) > 1)); then
      echo "Faulty resources found."
      if get_user_confirmation "Perform recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
        echo "Waiting $TIMEOUT_SEC seconds and rechecking for faulty resources"
        sleep $TIMEOUT_SEC
        continue
      else
        exit_function
      fi
    else
      echo "No faulty resources found"
      return
    fi
  done
}

linstor_check_advise() {
  while true; do
    echo "Checking for advise"
    if (( $(linstor advise r | tee /dev/tty | grep  "[a-zA-Z0-9]" | wc -l) > 1)); then
      echo "Advise found."
      echo "It is recommended to perform the advised actions manually."
      if get_user_confirmation "Do you perform the advised actions manually?" "y" "n"; then
        linstor_check_faulty
      else
        exit_function
      fi
    else
      echo "No advise found"
      break
    fi
  done
}

linstor_check_connection() {
  while true; do
    echo "Checking connection of LINSTOR controller to its satellites."

    SATELLITES_ONLINE=$(linstor -m --output-version=v1 node list | jq -r '.[][] | select(.type == "SATELLITE" and .connection_status == "ONLINE").name')
    if [[ -z $SATELLITES_ONLINE ]]; then
      echo "No satellites are online. This is usually a sign of issues with the LINSTOR controller operation. It is recommended to restart the controller and satellites."
      echo "List of satellites:"
      exec_linstor_with_exit_code_check node list
      if get_user_confirmation "Perform connection recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
        echo "Waiting $TIMEOUT_SEC seconds and rechecking the connection of LINSTOR controller to its satellites"
        sleep $TIMEOUT_SEC
      else
        exit_function
      fi
    else
      if [ $(linstor -m --output-version=v1 storage-pool list -s DfltDisklessStorPool -n $SATELLITES_ONLINE | jq '.[][].reports[]?.message' | grep 'No active connection to satellite' | wc -l) -ne 0 ]; then
        echo "Some satellites are not connected, even though they are online. This is usually a sign of issues with the LINSTOR controller operation. It is recommended to restart the controller and satellites."
        echo "List of satellites:"
        exec_linstor_with_exit_code_check node list
        echo "List of storage pools:"
        exec_linstor_with_exit_code_check storage-pool list
        if get_user_confirmation "Perform connection recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
          echo "Waiting $TIMEOUT_SEC seconds and rechecking the connection of LINSTOR controller to its satellites"
          sleep $TIMEOUT_SEC
        else
          exit_function
        fi
      else
        echo "LINSTOR controller has connection with all satellites that are online."
        break
      fi
    fi
  done
}

linstor_wait_sync() {
  local max_number_of_parallel_syncs=$1
  while true; do
    echo "Checking the number of replicas currently syncing"
    export SYNC_TARGET_RESOURCES=$(linstor -m --output-version=v1 resource list-volumes | jq -r '.[][] | select(.volumes[].state.disk_state | contains("SyncTarget")).name') 
    if [[ -n "${SYNC_TARGET_RESOURCES}" ]]; then
      echo "Resources found to be syncing at the moment. List of such resources:"
      exec_linstor_with_exit_code_check resource list-volumes -r ${SYNC_TARGET_RESOURCES}
      local number_of_parallel_syncs=$(echo ${SYNC_TARGET_RESOURCES} | wc -w)
      if (( ${number_of_parallel_syncs} > ${max_number_of_parallel_syncs} )); then
        echo "Number of sync operations at the moment=${number_of_parallel_syncs}. This is more than the maximum allowed number of simultaneously performed sync operations (${max_number_of_parallel_syncs}). Waiting for synchronization to complete."
        sleep 15
      else
        echo "Number of sync operations at the moment=${number_of_parallel_syncs}. This is less than or equal to the maximum allowed number of simultaneously performed sync operations (${max_number_of_parallel_syncs}). Ending synchronization wait."
        break
      fi
    else
      echo "No resources found to be syncing at the moment. Ending synchronization wait."
      break
    fi
  done
}


is_linstor_satellite_online() {
  node_connection_status=$(linstor -m --output-version=v1 node list -n ${NODE_FOR_EVICT} | jq -r --arg nodeName "${NODE_FOR_EVICT}" '.[][] | select(.name == $nodeName).connection_status')
  if [[ ${node_connection_status^^} == "ONLINE" ]]; then
    return 0
  else
    echo "Node ${NODE_FOR_EVICT} is not ONLINE in LINSTOR."
    echo "List of satellites:"
    exec_linstor_with_exit_code_check node list 
    return 1
  fi
}

is_linstor_satellite_does_not_exist() {
node_info=$(linstor -m --output-version=v1 node list -n ${NODE_FOR_EVICT} | jq -r --arg nodeName "${NODE_FOR_EVICT}" '.[][] | select(.name == $nodeName).connection_status')
if [[ -z $node_info ]]; then
  echo "Node ${NODE_FOR_EVICT} does not exist in LINSTOR."
  echo "List of satellites:"
  exec_linstor_with_exit_code_check node list
  return 0
else
  return 1
fi
}

linstor_change_replicas_count() {
  local replicas_to_add=$1
  shift
  local timeout=$1
  shift
  local resource_and_group_names=$1
  shift
  local resource_groups=$1
  shift
  local resource_names_list=("$@")

  
  local changed_resources=0

  for resource_name in $resource_names_list; do
    echo "Beginning the process of changing the count of diskfull replicas for resource ${resource_name}, which has replicas on the node being evicted ${NODE_FOR_EVICT}"

    local is_tiebreaker_needed=false
    resource_group=$(echo $resource_and_group_names | jq -r --arg resource_name "${resource_name}" '. | select(.resource == $resource_name).resource_group')
    place_count=$(echo $resource_groups | jq -r --arg resource_group "${resource_group}" '. | select(.resource_group == $resource_group).place_count')

    desired_diskful_replicas_count=$((place_count + replicas_to_add))
    # if (( $desired_diskful_replicas_count % 2 == 0 )); then # auto tiebreaker is not worked with more than 2 replicas (4, 6, 8, etc.)
    if (( $desired_diskful_replicas_count == 2 )); then
      is_tiebreaker_needed=true
    fi

    while true; do
      current_diskful_replicas_count=$(linstor -m --output-version=v1 resource list -r ${resource_name} | jq -r  --arg disklessStorPoolName "${DISKLESS_STORAGE_POOL}" '[.[][] | select(.props.StorPoolName != $disklessStorPoolName).name] | length')    
      if [[ -z $current_diskful_replicas_count ]]; then
        echo "Warning! Can't get the total number of diskfull replicas for resource ${resource_name}."
        echo "Resource status:"
        exec_linstor_with_exit_code_check resource list-volumes -r $resource_name
        if get_user_confirmation "Perform recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
          echo "Waiting $TIMEOUT_SEC seconds and rechecking the total number of diskfull replicas for resource ${resource_name}"
          sleep $TIMEOUT_SEC
        else
          exit_function
        fi
      else
        break
      fi
    done


    echo "Current status of the resource:"
    exec_linstor_with_exit_code_check resource list-volumes -r $resource_name
    sleep 2
    if (( ${current_diskful_replicas_count} != ${desired_diskful_replicas_count} )); then
      echo "The total number of diskfull replicas for resource ${resource_name} (${current_diskful_replicas_count}) does not equal the desired number of replicas(${desired_diskful_replicas_count})"
      echo "Performing checks before changing replica count"
      linstor_check_faulty
      linstor_wait_sync 3
      echo -n "Setting replica count to ${desired_diskful_replicas_count} for resource ${resource_name}"
      exec_linstor_with_exit_code_check resource-definition auto-place --place-count ${desired_diskful_replicas_count} ${resource_name}
      sleep ${timeout}
      echo "Resource status after changing replica count:"
      exec_linstor_with_exit_code_check resource list-volumes -r $resource_name

      while true; do
        current_diskful_replicas_count_new=$(linstor -m --output-version=v1 resource list -r ${resource_name} | jq -r  --arg disklessStorPoolName "${DISKLESS_STORAGE_POOL}" '[.[][] | select(.props.StorPoolName != $disklessStorPoolName).name] | length')
        if [[ -z $current_diskful_replicas_count_new ]]; then
          echo "Warning! Can't get the total number of diskfull replicas for resource ${resource_name}."
          echo "Resource status:"
          exec_linstor_with_exit_code_check resource list-volumes -r $resource_name
          if get_user_confirmation "Perform recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
            echo "Waiting $TIMEOUT_SEC seconds and rechecking the total number of diskfull replicas for resource ${resource_name}"
            sleep $TIMEOUT_SEC
          else
            exit_function
          fi
        else
          break
        fi
      done
      
      if (( ${current_diskful_replicas_count_new} != ${desired_diskful_replicas_count} )); then
        echo "Warning! The total number of diskfull replicas for resource ${resource_name} (${current_diskful_replicas_count}) does not equal the desired number of replicas(${desired_diskful_replicas_count}) even after changing the replica count. The following command was executed: \"linstor resource-definition auto-place --place-count ${desired_diskful_replicas_count} $resource_name\". Manual action required."
        exit_function
      fi
      changed_resources=$((changed_resources + 1))
    else
      echo "The total number of diskfull replicas for resource ${resource_name} (${current_diskful_replicas_count}) equals the desired number of diskfull replicas(${desired_diskful_replicas_count}). Skipping changing the count of diskful replicas for this resource"
    fi
    if [[ ${is_tiebreaker_needed} == true ]]; then
        echo "Resource ${resource_name} has an even number of diskful replicas. Checking for the presence of a TieBreaker for this resource."
        create_tiebreaker ${resource_name}
    fi
          
    echo "Processing of resource $resource_name completed."
    sleep 2
  done
  return $changed_resources
}


create_tiebreaker() {
  resource_name=$1

  while true; do
    diskless_replicas_count=$(linstor -m --output-version=v1 resource list -r ${resource_name} | jq -r  --arg disklessStorPoolName "${DISKLESS_STORAGE_POOL}" '[.[][] | select(.props.StorPoolName == $disklessStorPoolName).name] | length')
    if [[ -z $diskless_replicas_count ]]; then
      echo "Warning! Can't get the total number of diskless replicas for resource ${resource_name}."
      echo "Resource status:"
      exec_linstor_with_exit_code_check resource list-volumes -r $resource_name
      if get_user_confirmation "Perform recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
        echo "Waiting $TIMEOUT_SEC seconds and rechecking the total number of diskless replicas for resource ${resource_name}"
        sleep $TIMEOUT_SEC
      else
        exit_function
      fi
    else
      break
    fi
  done



  if (( ${diskless_replicas_count} < 1 )); then
    echo "The count of diskless replicas is ${diskless_replicas_count}. Creating a TieBreaker for resource ${resource_name}"
    exec_linstor_with_exit_code_check resource-definition set-property $resource_name DrbdOptions/auto-add-quorum-tiebreaker true
    sleep 5

    while true; do
      diskless_replicas_count_new=$(linstor -m --output-version=v1 resource list -r ${resource_name} | jq -r  --arg disklessStorPoolName "${DISKLESS_STORAGE_POOL}" '[.[][] | select(.props.StorPoolName == $disklessStorPoolName).name] | length')
      if [[ -z $diskless_replicas_count_new ]]; then
        echo "Warning! Can't get the total number of diskless replicas for resource ${resource_name}."
        echo "Resource status:"
        exec_linstor_with_exit_code_check resource list-volumes -r $resource_name
        if get_user_confirmation "Perform recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
          echo "Waiting $TIMEOUT_SEC seconds and rechecking the total number of diskless replicas for resource ${resource_name}"
          sleep $TIMEOUT_SEC
        else
          exit_function
        fi
      else
        break
      fi
    done
    if (( ${diskless_replicas_count_new} < 1 )); then
      echo "Warning! TieBreaker for the resource did not create. Resource status:"
      exec_linstor_with_exit_code_check resource list-volumes -r $resource_name
      exit_function
    fi
  else
    echo "The count of diskless replicas is ${diskless_replicas_count}. TieBreaker for resource ${resource_name} is not needed."
  fi
}

linstor_delete_resources_from_node() {
  local resource_names_list=("$@")

  for resource_name in $resource_names_list; do
    echo "Beginning the process of deleting resource ${resource_name}, which has replicas on the node being evicted ${NODE_FOR_EVICT}"
    linstor_check_connection
    linstor_check_faulty
    linstor_wait_sync 0
    echo "Current status of the resource:"
    exec_linstor_with_exit_code_check resource list-volumes -r ${resource_name}
    sleep 2
    echo "Deleting resource ${resource_name}"
    exec_linstor_with_exit_code_check resource delete ${NODE_FOR_EVICT} ${resource_name}
    echo "Resource status after deleting:"
    sleep 2
    exec_linstor_with_exit_code_check resource list-volumes -r ${resource_name}
    sleep 2

    while true; do
      current_diskful_replicas_count=$(linstor -m --output-version=v1 resource list -r ${resource_name} | jq -r  --arg disklessStorPoolName "${DISKLESS_STORAGE_POOL}" '[.[][] | select(.props.StorPoolName != $disklessStorPoolName).name] | length')    
      if [[ -z $current_diskful_replicas_count ]]; then
        echo "Warning! Can't get the total number of diskfull replicas for resource ${resource_name}."
        echo "Resource status:"
        exec_linstor_with_exit_code_check resource list-volumes -r $resource_name
        if get_user_confirmation "Perform recheck in $TIMEOUT_SEC seconds?" "y" "n"; then
          echo "Waiting $TIMEOUT_SEC seconds and rechecking the total number of diskfull replicas for resource ${resource_name}"
          sleep $TIMEOUT_SEC
        else
          exit_function
        fi
      else
        break
      fi
    done
    
    if (( $current_diskful_replicas_count == 2 )); then
      echo "Resource ${resource_name} has an even number of diskfull replicas. Creating a TieBreaker for this resource if it does not exist."
      create_tiebreaker ${resource_name}
    fi
    
    echo "Processing of resource $resource_name completed."
  done

}

get_user_confirmation() {
  local prompt="$1"
  local positive_case="$2"
  local negative_case="$3"

  while true; do
    echo -n "$prompt ($positive_case/$negative_case): "
    read user_input
    
    case "$user_input" in
      "$positive_case")
        return 0  # true
        ;;
      "$negative_case")
        return 1  # false
        ;;
      *)
        echo "Invalid input. Try again."
        sleep 1
        ;;
    esac
  done

}

exit_function(){
  if get_user_confirmation "Terminate the script?" "y" "n"; then
    if get_user_confirmation "Return node ${NODE_FOR_EVICT} to LINSTOR scheduler?" "y" "n"; then
      echo "Returning node ${NODE_FOR_EVICT} to LINSTOR scheduler"
      exec_linstor_with_exit_code_check node set-property ${NODE_FOR_EVICT} AutoplaceTarget
    fi
    if get_user_confirmation "Perform uncordon on node ${NODE_FOR_EVICT}?" "y" "n"; then
      echo "Performing uncordon on node ${NODE_FOR_EVICT}"
      execute_command "kubectl uncordon ${NODE_FOR_EVICT}"
    fi
    echo "Terminating the script"
    exit 0
  fi
  echo "The script operation will be continued."
}

kubernetes_check_node() {
  export DISKFUL_RESOURCES_TO_EVICT=$(linstor -m --output-version=v1 resource list -n ${NODE_FOR_EVICT} | jq -r --arg nodeName "${NODE_FOR_EVICT}" --arg disklessStorPoolName "${DISKLESS_STORAGE_POOL}" '.[][] | select(.node_name == $nodeName and .props.StorPoolName != $disklessStorPoolName).name')
  if [[ -z "${DISKFUL_RESOURCES_TO_EVICT}" ]]; then
      echo
      echo "List of resources to evict is empty. Please enter the values again."
      echo "List of storage pools and nodes in LINSTOR:"
      exec_linstor_with_exit_code_check storage-pool list
      sleep 2
      return 1
  fi

  if execute_command "kubectl get nodes ${NODE_FOR_EVICT}" | grep -q "SchedulingDisabled"; then
    return 0
  else
    echo "The cordon command has not been executed for node ${NODE_FOR_EVICT}."
    if get_user_confirmation "Perform node drain? (if confirmed, the command \"kubectl drain ${NODE_FOR_EVICT} --delete-emptydir-data  --ignore-daemonsets\" will be executed)" "y" "n"; then
      execute_command "kubectl drain ${NODE_FOR_EVICT} --delete-emptydir-data  --ignore-daemonsets"
      echo "Drain of node ${NODE_FOR_EVICT} completed"
      return 0
    else
      echo "WARNING! Performing node drain is mandatory before evicting LINSTOR resources from it"
      exit_function
      return 1
    fi  
  fi
  
  return 0
}

wait_for_deployment_scale_down() {
  local DEPLOYMENT_NAME=$1
  local NAMESPACE=$2

  local count=0
  local max_attempts=60

  until [[ $(kubectl get pods -n "$NAMESPACE" -l app="$DEPLOYMENT_NAME" --no-headers 2>/dev/null | wc -l) -eq 0 ]] || [[ $count -eq $max_attempts ]]; do
    echo "Waiting for pods to be deleted... Attempt $((count+1))/$max_attempts."
    sleep 5
    ((count++))
  done
  
  if [[ $count -eq $max_attempts ]]; then
    echo "Timeout reached. Pods were not deleted."
    exit_function
    return
  fi
  echo "Pods were deleted."

}

linstor_backup_database() {
  echo "Performing database backup"
  
  while true; do
    echo "Checking the number of replicas for LINSTOR controller"
    LINSTOR_CONTROLLER_REPLICAS_COUNT=$(kubectl -n ${LINSTOR_NAMESPACE} get deployment linstor-controller -o jsonpath='{.spec.replicas}')
    if [[ -n "${LINSTOR_CONTROLLER_REPLICAS_COUNT}" ]]; then
      echo "Scale down LINSTOR controller"
      execute_command "kubectl -n ${LINSTOR_NAMESPACE} scale deployment linstor-controller --replicas=0"
      echo "Waiting for LINSTOR controller to scale down"
      wait_for_deployment_scale_down "linstor-controller" "${LINSTOR_NAMESPACE}"
      break
    else
      echo "Can't get the number of replicas for LINSTOR controller."
      if get_user_confirmation "Should we recheck the number of replicas for the LINSTOR controller after $TIMEOUT_SEC seconds? (Note that the database backup will not be performed if this is not done.)" "y" "n"; then
        echo "Waiting $TIMEOUT_SEC seconds and rechecking the number of replicas for LINSTOR controller"
        sleep $TIMEOUT_SEC
        continue
      else
        break
      fi
    fi
  done

  echo "Creating a backup of the LINSTOR database"
  export current_datetime=$(date +%Y-%m-%d_%H-%M-%S)
  mkdir linstor_db_backup_before_evict_${current_datetime}
  kubectl get crds | grep -o ".*.internal.linstor.linbit.com" | xargs kubectl get crds -oyaml > ./linstor_db_backup_before_evict_${current_datetime}/crds.yaml
  kubectl get crds | grep -o ".*.internal.linstor.linbit.com" | xargs -i{} sh -xc "kubectl get {} -oyaml > ./linstor_db_backup_before_evict_${current_datetime}/{}.yaml"
  echo "Database backup completed"
  echo "Scale up LINSTOR controller"
  execute_command "kubectl -n ${LINSTOR_NAMESPACE} scale deployment linstor-controller --replicas=${LINSTOR_CONTROLLER_REPLICAS_COUNT}"
  echo "Waiting for LINSTOR controller to scale up"
  sleep 15
  linstor_check_faulty
}

delete_node_from_kubernetes_and_linstor() {
  while true; do
    deckhouse_current_replicas=$(kubectl -n d8-system get deployment deckhouse -o jsonpath='{.spec.replicas}')
    piraeus_current_replicas=$(kubectl -n ${LINSTOR_NAMESPACE} get deployment piraeus-operator -o jsonpath='{.spec.replicas}')

    if [[ -z "${deckhouse_current_replicas}" || -z "${piraeus_current_replicas}" ]]; then
      echo "Can't get the number of replicas for Deckhouse or Piraeus operator."
      if get_user_confirmation "Should we recheck the number of replicas for Deckhouse and Piraeus operator after $TIMEOUT_SEC seconds? (Note that the node will not be deleted from Kubernetes and LINSTOR if this is not done.)" "y" "n"; then
        echo "Waiting $TIMEOUT_SEC seconds and rechecking the number of replicas for Deckhouse and Piraeus operator"
        sleep $TIMEOUT_SEC
        continue
      else
        echo "Warning! The node will not be deleted from Kubernetes and LINSTOR."
        return
      fi
    fi
    break
  done

  while true; do
    if is_linstor_satellite_online; then
      linstor_satellite_online="true"
    else
      linstor_satellite_online="false"
    fi
      
    if [[ $linstor_satellite_online == "true" ]]; then
      echo "Node $NODE_FOR_EVICT is ONLINE in LINSTOR. Performing standard node deletion procedure from LINSTOR"
    else
      echo "Warning! Node ${NODE_FOR_EVICT} is not ONLINE in LINSTOR. It is impossible to perform standard node deletion procedure from LINSTOR. The command \"linstor node lost ${NODE_FOR_EVICT}\" needs to be executed."
      if get_user_confirmation "Perform a re-check of the connection with node ${NODE_FOR_EVICT} after $TIMEOUT_SEC seconds?" "y" "n"; then
        echo "Waiting for $TIMEOUT_SEC seconds and re-checking the connection with node ${NODE_FOR_EVICT}"
        sleep $TIMEOUT_SEC
        continue
      fi
    fi
  done

  echo "The procedure for deleting a node from LINSTOR will consist of the following actions:" 
  echo "1. Shutting down Deckhouse and Piraeus operator"
  if [[ $linstor_satellite_online == "true" ]]; then
    echo "2. Standard node deletion from LINSTOR"
  else
    echo "2. Executing the command \"linstor node lost ${NODE_FOR_EVICT}\""
  fi
  echo "3. Deleting the node from Kubernetes"
  echo "4. Turning Deckhouse and Piraeus operator back on"

  if get_user_confirmation "Perform the actions listed above?" "yes-i-am-sane-and-i-understand-what-i-am-doing" "n"; then
    execute_command "kubectl -n d8-system scale deployment deckhouse --replicas=0"
    echo "Waiting for Deckhouse to scale down"
    wait_for_deployment_scale_down "deckhouse" "d8-system"

    execute_command "kubectl -n ${LINSTOR_NAMESPACE} scale deployment piraeus-operator --replicas=0"
    echo "Waiting for piraeus-operator to scale down"
    wait_for_deployment_scale_down "piraeus-operator" "${LINSTOR_NAMESPACE}"

    if [[ $linstor_satellite_online == "true" ]]; then
      echo "Performing standard node deletion procedure from LINSTOR"
      exec_linstor_with_exit_code_check node delete ${NODE_FOR_EVICT}
    else
      echo "Executing the command \"linstor node lost ${NODE_FOR_EVICT}\""
      exec_linstor_with_exit_code_check node lost ${NODE_FOR_EVICT}
    fi

    if is_linstor_satellite_does_not_exist; then
      execute_command "kubectl delete node ${NODE_FOR_EVICT}"
      execute_command "kubectl -n d8-system scale deployment deckhouse --replicas=${deckhouse_current_replicas}"
      execute_command "kubectl -n ${LINSTOR_NAMESPACE} scale deployment piraeus-operator --replicas=${piraeus_current_replicas}"
      return
    else
      echo "Warning! Node ${NODE_FOR_EVICT} has not been deleted from LINSTOR. Turning Deckhouse and Piraeus operator back on."
      execute_command "kubectl -n d8-system scale deployment deckhouse --replicas=${deckhouse_current_replicas}"
      execute_command "kubectl -n ${LINSTOR_NAMESPACE} scale deployment piraeus-operator --replicas=${piraeus_current_replicas}"
      echo "It is recommended to terminate the script and investigate the cause of the error."
      exit_function
      return
    fi
  else
    exit_function
    return
  fi

}


echo "The script for evicting LINSTOR resources has been launched. Performing necessary checks before starting."
linstor_check_faulty
linstor_check_connection
linstor_wait_sync 0

echo "List of storage pools and nodes in LINSTOR:"
exec_linstor_with_exit_code_check storage-pool list

while true; do
  while true; do
    echo -n "Enter the name of the node from which LINSTOR resources need to be evicted: "
    read NODE_FOR_EVICT
    if [[ -z "$NODE_FOR_EVICT" ]]; then
        echo "Name cannot be empty. Please enter the value again."
    else
        break
    fi
  done

  echo
  
  if kubernetes_check_node; then
    break
  fi
  
done

if get_user_confirmation "Perform database backup before evicting resources from node ${NODE_FOR_EVICT}?" "y" "n"; then
  linstor_backup_database
fi

echo "Excluding node ${NODE_FOR_EVICT} from LINSTOR scheduler"
exec_linstor_with_exit_code_check node set-property ${NODE_FOR_EVICT} AutoplaceTarget false

RESOURCE_AND_GROUP_NAMES=$(linstor -m --output-version=v1 resource-definition list -r ${DISKFUL_RESOURCES_TO_EVICT} | jq -r '.[][] | {resource: .name, resource_group: .resource_group_name}')
RESOURCE_GROUPS=$(linstor -m --output-version=v1 resource-group list | jq -r '.[][] | {resource_group: .name, place_count: .select_filter.place_count}')

echo "List of resources to be evicted from node ${NODE_FOR_EVICT}:"
exec_linstor_with_exit_code_check resource list-volumes -r ${DISKFUL_RESOURCES_TO_EVICT}

linstor_change_replicas_count 1 10 "${RESOURCE_AND_GROUP_NAMES}" "${RESOURCE_GROUPS}" "${DISKFUL_RESOURCES_TO_EVICT[@]}" 
echo "Increase in replica count for movable resources completed"


# If any replica here got stuck in Inconsistent, then the following can be done:
## Deactivate and subsequently activate the resource on the node where synchronization got stuck
## If that didn't help, then remove the replica on the problematic node and manually create it on another node with the command linstor r create <node name> <resource name>

echo "Status of processed resources"
exec_linstor_with_exit_code_check resource list-volumes -r ${DISKFUL_RESOURCES_TO_EVICT}
sleep 2

echo "Performing checks before deleting the node"
linstor_check_faulty
linstor_wait_sync 0

echo "Getting the list of resources, replicas of which are on the evicted node ${NODE_FOR_EVICT} again, to check if the new resources have appeared"
export DISKFUL_RESOURCES_TO_EVICT_NEW=$(linstor -m --output-version=v1 resource list -n ${NODE_FOR_EVICT} | jq -r --arg nodeName "${NODE_FOR_EVICT}" --arg disklessStorPoolName "${DISKLESS_STORAGE_POOL}" '.[][] | select(.node_name == $nodeName and .props.StorPoolName != $disklessStorPoolName).name')
RESOURCE_AND_GROUP_NAMES_NEW=$(linstor -m --output-version=v1 resource-definition list -r ${DISKFUL_RESOURCES_TO_EVICT_NEW} | jq -r '.[][] | {resource: .name, resource_group: .resource_group_name}')
RESOURCE_GROUPS_NEW=$(linstor -m --output-version=v1 resource-group list | jq -r '.[][] | {resource_group: .name, place_count: .select_filter.place_count}')
added_resources=$(comm -13 <(echo "$DISKFUL_RESOURCES_TO_EVICT" | sort) <(echo "$DISKFUL_RESOURCES_TO_EVICT_NEW" | sort))

echo "Old resource list: ${DISKFUL_RESOURCES_TO_EVICT}"
echo "New resource list: ${DISKFUL_RESOURCES_TO_EVICT_NEW}"
if [[ -n "${added_resources}" ]]; then
  echo "The following resources have appeared on the node ${NODE_FOR_EVICT}:"
  exec_linstor_with_exit_code_check resource list-volumes -r ${added_resources}
  echo "Start increasing the number of replicas for these resources"
  added_resource_and_group_names=$(linstor -m --output-version=v1 resource-definition list -r ${added_resources} | jq -r '.[][] | {resource: .name, resource_group: .resource_group_name}')
  linstor_change_replicas_count 1 2 "${added_resource_and_group_names}" "${RESOURCE_GROUPS_NEW}" "${added_resources[@]}"
  linstor_check_faulty
  linstor_wait_sync 0
else
  echo No new resources have appeared on the node ${NODE_FOR_EVICT}
fi

echo "Get resources that have TieBreaker on evicted node ${NODE_FOR_EVICT}"
RESOURCES_WITH_TIEBREAKER_ON_EVICTED_NODE=$(linstor -m --output-version=v1 resource list -n ${NODE_FOR_EVICT} | jq -r --arg nodeName "${NODE_FOR_EVICT}" --arg disklessStorPoolName "${DISKLESS_STORAGE_POOL}" '.[][] | select(.node_name == $nodeName and .props.StorPoolName == $disklessStorPoolName and .state.in_use == false).name')

echo "Attention! Before evacuate resources from node ${NODE_FOR_EVICT}, make sure that all resources in LINSTOR are in UpToDate state."

if get_user_confirmation "Do you want to delete node ${NODE_FOR_EVICT} from Kubernetes and LINSTOR? If not, then only the deletion of all resources from ${NODE_FOR_EVICT} will be executed." "y" "n"; then  
  delete_node_from_kubernetes_and_linstor
else
  if get_user_confirmation "Do you want to delete all resources from node ${NODE_FOR_EVICT}?" "y" "n"; then
    echo "Only the deletion of all resources from ${NODE_FOR_EVICT} will be executed. The node will NOT be deleted from Kubernetes and LINSTOR."
    linstor_delete_resources_from_node "${RESOURCES_WITH_TIEBREAKER_ON_EVICTED_NODE[@]}"
    linstor_delete_resources_from_node "${DISKFUL_RESOURCES_TO_EVICT_NEW[@]}"
  else
    echo "Deletion of resources from node ${NODE_FOR_EVICT} will not be performed."
  fi
fi


linstor_check_faulty
echo "State of all resources after evacuate resources from node ${NODE_FOR_EVICT}:"
exec_linstor_with_exit_code_check resource list-volumes
sleep 2

while true; do
  if get_user_confirmation "Should we reduce the number of replicas to their previous values for movable resources and also check for TieBreakers?" "y" "n"; then
    linstor_change_replicas_count 0 2 "${RESOURCE_AND_GROUP_NAMES_NEW}" "${RESOURCE_GROUPS_NEW}" "${DISKFUL_RESOURCES_TO_EVICT_NEW[@]}"
    break
  fi
  exit_function
done

linstor_check_faulty
echo "Reduction in the number of replicas for movable resources completed"

if [[ -n $RESOURCES_WITH_TIEBREAKER_ON_EVICTED_NODE ]]; then
  echo "Create TieBreaker if needed for resources that had TieBreaker on evicted node \"${NODE_FOR_EVICT}\""
  for resource in $RESOURCES_WITH_TIEBREAKER_ON_EVICTED_NODE; do
    create_tiebreaker $resource
  done
fi

linstor_check_faulty
linstor_check_advise

echo "Script operation completed"
