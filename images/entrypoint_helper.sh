#!/bin/bash
# Common helper functions for entrypoint scripts

set -eo pipefail
shopt -s nullglob

DORIS_HOME="/opt/apache-doris"

# account name and password to execute sql
declare ACC_USER
declare ACC_PWD

ACC_USER_DEFAULT="k8sopr"
ACC_PWD_DEFAULT="NM7Cr4k9SfH6f5w0pdEJ4A=="

if [[ -z $ACC_USER ]]; then
  ACC_USER=$ACC_USER_DEFAULT
fi
if [[ -z $ACC_PWD ]]; then
  ACC_PWD=$ACC_PWD_DEFAULT
fi

# Logging functions
doris_log() {
  local type=$1
  shift
  # accept argument string or stdin
  local text="$*"
  if [ "$#" -eq 0 ]; then text="$(cat)"; fi
  local dt
  dt="$(date -Iseconds)"
  printf '%s [%s] [Entrypoint]: %s\n' "$dt" "$type" "$text"
}
doris_note() {
  doris_log NOTE "$@"
}
doris_warn() {
  doris_log WARN "$@" >&2
}
doris_error() {
  doris_log ERROR "$@" >&2
  exit 1
}

# Get the FQDN DNS of the current container.
myself_host() {
  hostname -f
#  if [[ -n $POD_NAME ]]; then
#    if [[ -n $POD_NAMESPACE ]]; then
#      echo "${POD_NAME}.${POD_NAMESPACE}"
#    else
#      echo "${POD_NAME}"
#    fi
#  else
#    hostname -f
#  fi
}

# Injects an entry in "key=value" format into the specified file when it does not exist.
inject_item_into_conf_file() {
  local conf_file=$1
  local key=$2
  local value=$3
  if ! grep -qE "^${key}\s*\=\s*${value}" "$conf_file"; then
    echo "" >>"$conf_file"
    echo "${key}=${value}" >>"$conf_file"
    doris_note "Inject '${key}=${value}' into ${conf_file}"
  fi
}

# Get the value corresponding to the key of the specified doris config file
# with optional default value.
get_value_from_conf_file() {
  local conf_file=$1
  local key=$2
  local default_value=$3
  local value
  value=$(grep "\<$key\>" "$conf_file" | grep -v '^\s*#' | sed 's|^\s*'"$key"'\s*=\s*\(.*\)\s*$|\1|g')
  if [[ -z $value && -n $default_value ]]; then
    value=$default_value
  fi
  echo "$value"
}
