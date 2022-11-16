#!/bin/bash

#set -o errexit
#set -o nounset
#set -o xtrace

host="${HOST:-localhost:8080/hub}"

# Port Forwarding
kubectl port-forward service/tackle-ui 8080:8080 -n konveyor-tackle > /dev/null 2>&1 &
pid=$!

# kill the port-forward regardless of how this script exits
trap '{
    kill $pid
}' EXIT

# wait for port to become available
while ! nc -vz localhost 8080 > /dev/null 2>&1 ; do
    sleep 0.1
done

# Make a request to hub
request_cmd="$(curl -i -o - -X POST ${host}/tasks -d \
'{
    "name":"Jkube",
    "state": "Ready",
    "locator": "jkube",
    "addon": "jkube",
    "application": {"id": 1},
    "data": {}
}')"

# Get status code from the curl request
status_code="$(echo "$request_cmd" | grep HTTP | awk '{print $2}')"

# Get output from the curl request
output_response=$(echo "$request_cmd")
echo "Output response: $output_response"

# Check if status_code starts with 2
if [[ "${status_code}" != 2* ]]; then
    echo "Failed to create jkube task"
    echo "Got Response Status: ${status_code}"
    exit 1
fi