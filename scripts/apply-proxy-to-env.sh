#!/bin/bash

set -exu -o pipefail

. "$(dirname "$0")/lib/environment.sh"

export BOSH_ENV="${KUBO_ENVIRONMENT_DIR}"
export BOSH_NAME=$(basename ${BOSH_ENV})
export DEBUG=1

cp "$PWD/kubo-lock/metadata" "${KUBO_ENVIRONMENT_DIR}/director.yml"
cp "$PWD/gcs-bosh-creds/creds.yml" "${KUBO_ENVIRONMENT_DIR}/"

. "$PWD/git-kubo-deployment/bin/lib/deploy_utils"
. "$PWD/git-kubo-deployment/bin/set_bosh_environment"

tinyproxy_ip=$(BOSH_CLIENT=bosh_admin BOSH_CLIENT_SECRET="$(get_bosh_secret)" bosh-cli -n -e "${BOSH_ENVIRONMENT}" -d tinyproxy vms --json | bosh-cli int - --path='/Tables/Content=vms/Rows/0/ips')
proxy_setting="$tinyproxy_ip:8888"

cp -r kubo-lock/* kubo-lock-with-proxy/

echo >> kubo-lock-with-proxy/metadata
echo "http_proxy: $proxy_setting" >> kubo-lock-with-proxy/metadata
echo "https_proxy: $proxy_setting" >> kubo-lock-with-proxy/metadata
