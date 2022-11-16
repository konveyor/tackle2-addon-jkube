Tackle2 Addon Jkube
===================
This project exposes [JKube's](https://www.eclipse.org/jkube) `kubernetes-maven-plugin` as an addon for [Tackle2](https://github.com/konveyor?q=tackle2&type=all&language=&sort=), enabling Java applications with Maven integration to have container image and Kubernetes manifests built for them.

This addon focus on two tasks:
- Generates container images with flexible and powerful configuration.
- Generates vanilla Kubernetes descriptors (YAML files).

It leverages the [Kubernetes Maven plugin](https://www.eclipse.org/jkube/docs/kubernetes-maven-plugin)  to generate container image and Kubernetes manifests.

## Development
To browse code - [Open in VSCode](https://open.vscode.dev/konveyor/tackle2-addon-jkube)

We use [tackle2-hub addon](https://github.com/konveyor/tackle2-hub/tree/main/addon) package to integrate addon with [tackle2-hub](https://github.com/konveyor/tackle2-hub). Here is a good starter template to build an addon - [Test Addon](https://github.com/konveyor/tackle2-hub/tree/main/hack/cmd/addon).
Tackle Hub handles all requests to an addon via [Task APIs](https://github.com/konveyor/tackle2-hub/blob/main/api/task.go) and forward these requests to desired addon.


A simple request for an addon would look like this:
```bash
host="${HOST:-localhost:8080/hub}"

# Make a request to hub
request_cmd="$(curl -i -o - -X POST ${host}/tasks -d \
'{
    "name":"Jkube",
    "state": "Ready",
    "locator": "jkube",
    "addon": "jkube",
    "application": {"id": 2},
    "data": {}
}')"
```

In above request we mention `addon: jkube` as an addon identifier and `application: {id: 2}` states the application against which this addon will be executed. The `data` field here is empty because we are not accepting any data for this addon, but we can pass data during api request to addons which accepts data, and we use [this](https://github.com/konveyor/tackle2-addon-windup/blob/84174448d3d7cd2abc7ba6ab27e66a55890b9061/cmd/main.go#L33-L47) struct to receive data for the addon.

Once a request is received - addon will fetch the application mentioned in the request and perform the operations defined in the addon. For this addon we perform following operations:
- Fetch the application (clone GitHub/subversion repository)
- Parse the maven config file and add kubernetes-maven-plugin to `pom.xml`
- Build the application using maven build
- Use kubernetes-maven-plugin to generate container image and k8s manifests
- Commit the generated resources to repo

## Development Environment Setup

Steps to successfully execute Jkube addon:
1. Clone this repo
    ```bash
    https://github.com/konveyor/tackle2-addon-jkube.git
    ```
2. Start a minikube environment and install tackle
    ```bash
    make start-minikube install-tackle
    ```
3. Add Jkube addon to tackle - [jkube addon](https://github.com/konveyor/tackle2-addon-jkube/blob/main/hack/addon.yml)
    ```bash
    kubectl apply -f hack/addon.yml
    ```
4. Create an **application** in Tackle using [this](https://github.com/konveyor/tackle2-addon-jkube/blob/a90ea44e8a4dadbcbdca556e52dce71ebb1b78b9/hack/task.sh#L39-L54) command
    ```bash
    host="${HOST:-localhost:8080/hub}"
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
        "businessService": {"id":1}
    }' | jq -M .
    ```

5. Get an application id (required in step 7), list all applications using
    ```bash
    kubectl port-forward service/tackle-ui 8080:8080 -n konveyor-tackle
    curl -X GET localhost:8080/hub/applications | jq .
    ```

6. Populate required data in tackle, including **application** against which addon will be executed - [reference](https://github.com/konveyor/tackle2-hub/tree/main/hack/add), [jkube-data](https://github.com/konveyor/tackle2-addon-jkube/blob/main/hack/task.sh)
    ```bash
    bash hack/task.sh
    ```

7. Create task to run JKube addon against our application - [task](https://github.com/konveyor/tackle2-addon-jkube/blob/a90ea44e8a4dadbcbdca556e52dce71ebb1b78b9/hack/task_ready.sh#L10-L18)
    ```bash
    bash hack/task_ready.sh
    ```

Follow the [Jkube Addon Demonstration](#jkube-addon-demonstration) for end-to-end setup and execution.

## Reference - Windup Addon
For [Windup addon](https://github.com/konveyor/tackle2-addon-windup), we have build an end-to-end [pipeline](https://github.com/konveyor/tackle2-addon-windup/blob/main/.github/workflows/test-windup.yml) to test the addon.

Here you will find everything you need to successfully execute an addon, we did it using GitHub actions. This pipeline is very loosely coupled, and you can utilize its pieces to create your own workflow, to test your addon quickly.

## Output
* Deployment artifacts
  * Dockerfile
  * Kubernetes Manifests

# Jkube Addon Demonstration

https://user-images.githubusercontent.com/20452032/201997511-3b2a86cb-dfd0-4a0a-9c90-bbc0fdbd987e.mov

📽️ [Watch on YouTube](https://youtu.be/fJM10cq7txg)
