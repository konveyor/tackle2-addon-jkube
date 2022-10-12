#!/bin/bash

set -o errexit
set -o nounset
set -o xtrace

host="${HOST:-localhost:8080/hub}"

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

# Create a Review
curl -X POST ${host}/reviews -d \
'{
    "businessCriticality": 4,
    "effortEstimate": "extra_large",
    "proposedAction": "repurchase",
    "workPriority": 1,
    "comments": "This is hard.",
    "application": {"id":1}
}' | jq -M .


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