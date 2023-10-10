#!/bin/bash

# Extra environment variables:
#  FE_SVC: FE service name, required.
#  FE_QUERY_PORT: FE service query port, optional, default: 9030
#  ACC_USER: account name to execute sql, optional, default: k8sopr
#  ACC_PWD: account password to execute sql, optional.

source entrypoint_helper.sh

BE_MIRROR_CONF_DIR=/etc/apache-doris/be/
BE_CONF_DIR=${DORIS_HOME}/be/conf/
BE_CONF_FILE=${DORIS_HOME}/be/conf/be.conf

# self fqdn host
declare SELF_HOST
# doris be heartbeat port
declare HEARTBEAT_PORT
# doris fe query port
declare FE_QUERY_PORT

# be probe interval: 2 seconds
BE_PROBE_INTERVAL=${BE_PROBE_INTERVAL:-2}
# timeout for probing cn: 60 seconds
BE_PROBE_TIMEOUT=${BE_PROBE_TIMEOUT:-60}

# collect env info from container
collect_env() {
  SELF_HOST=$(myself_host)
  HEARTBEAT_PORT=$(get_value_from_conf_file "$BE_CONF_FILE" 'heartbeat_service_port' 9050)
  if [[ -z $FE_QUERY_PORT ]]; then
    FE_QUERY_PORT=9030
  fi
}

override_be_conf() {
  if [[ -d "$BE_MIRROR_CONF_DIR" ]]; then
    cp -f "$BE_MIRROR_CONF_DIR"* "$BE_CONF_DIR"
  fi
  # make sure node role is mix
  inject_item_into_conf_file "$BE_CONF_FILE" 'be_node_role' 'mix'
}

show_backends() {
  timeout 15 mysql --connect-timeout 2 -h "$FE_SVC" -P "$FE_QUERY_PORT" -u"$ACC_USER" -p"$ACC_PWD" --skip-column-names --batch -e 'SHOW BACKENDS;'
}

# add self to cluster
add_self() {
  set +e
  local start
  local expire
  local now
  start=$(date +%s)
  expire=$((start + BE_PROBE_TIMEOUT))

  while true; do
    doris_note "Try to add myself($SELF_HOST:$HEARTBEAT_PORT) to cluster as BACKEND..."
    # check if it has been added to the cluster
    if show_backends | grep -q -w "$SELF_HOST" &>/dev/null; then
      doris_note "Myself($SELF_HOST:$HEARTBEAT_PORT) already exists in cluster."
      break
    fi

    timeout 15 mysql --connect-timeout 2 -h "$FE_SVC" -P "$FE_QUERY_PORT" -u"$ACC_USER" -p"$ACC_PWD" --skip-column-names --batch -e "ALTER SYSTEM ADD BACKEND \"$SELF_HOST:$HEARTBEAT_PORT\";"

    # check if it was added successfully
    if show_backends | grep -q -w "$SELF_HOST" &>/dev/null; then
      doris_note "Add myself to cluster successfully."
      break
    fi
    # check probe process timeout
    now=$(date +%s)
    if [[ $expire -le $now ]]; then
      doris_error "Add myself to cluster timed out."
    fi
    sleep $BE_PROBE_INTERVAL
  done
}

# main process
if [[ -z $FE_SVC ]]; then
  doris_error "Missing environment variable FE_SVC for the FE service name"
fi

collect_env
override_be_conf
add_self
doris_note "Ready to start BE!"
start_be.sh --console
