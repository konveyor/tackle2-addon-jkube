#!/bin/bash

set -o errexit
set -o nounset
set -o xtrace

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

# Create a Stake Holder Group
curl -X POST ${host}/stakeholdergroups -d \
'{
    "name": "Big Dogs",
    "description": "Group of big dogs."
}' | jq -M .

# Create a Stake Holder
curl -X POST ${host}/stakeholders -d \
'{
    "name": "tackle",
    "displayName":"Elmer",
    "email": "tackle@konveyor.org",
    "role": "Administrator",
    "stakeholderGroups": [{"id": 1}],
    "jobFunction" : {"id": 1}
}' | jq -M .

# Create a Business Service
curl -X POST ${host}/businessservices -d \
'{
    "createUser": "tackle",
    "name": "Marketing",
    "Description": "Marketing Dept.",
    "owner": {
      "id": 1
    }
}' | jq -M .


# Create an Application
curl -X POST ${host}/applications -d \
'{
    "createUser": "tackle",
    "name":"JkubeDemo",
    "description": "Spring Boot demo application.",
    "repository": {
      "name": "jkube_demo",
      "url": "https://github.com/mundra-ankur/jkube_demo.git",
      "branch": "main"
    },
    "facts": {
      "analysed": true
    },
    "businessService": {"id":1}
}' | jq -M .