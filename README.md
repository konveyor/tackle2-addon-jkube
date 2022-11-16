Tackle2 Addon Jkube
===================
The addon is part of the [Tackle2](https://github.com/konveyor?q=tackle2&type=all&language=&sort=) projects.

This is an implementation of [Jkube's](https://www.eclipse.org/jkube) `kubernetes-maven-plugin` as a konveyor addon. This addon is a small step towards Application Modernization. It brings your <strong>Java</strong> applications on to Kubernetes by leveraging the tasks required to make your application cloud-native. It provides a tight integration into Maven and benefits from the build configuration already provided. 

This addon focus on two tasks: 
- Generates container images with flexible and powerful configuration.
- Generates vanilla Kubernetes descriptors (YAML files).

It leverage the [Kubernetes Maven plugin](https://www.eclipse.org/jkube/docs/kubernetes-maven-plugin)  to generate conatiner image and Kubernetes resources.

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

In above request we mention `addon: jkube` as an addon identifier and `application: {id: 2}` states the application against which this addon will be executed. The `data` field here is empty but we can pass data to addon during api request and we use following struct to receive data for the addon.

```go
// Data Addon data passed in the secret.
type Data struct {
	// Output directory within application bucket.
	Output string `json:"output" binding:"required"`
	// Mode options.
	Mode Mode `json:"mode"`
	// Sources list.
	Sources Sources `json:"sources"`
	// Targets list.
	Targets Targets `json:"targets"`
	// Scope options.
	Scope Scope `json:"scope"`
	// Rules options.
	Rules *Rules `json:"rules"`
}
```

Once a request is received - addon will fetch the application mentioned in the request and perform the operations defined in the addon. For this addon we perform following operations:
- Fetch the application (clone github/subversion repository)
- Parse the maven config file and add kubernetes-maven-plugin to `pom.xml`
- Build the application using maven build
- Use kubernetes-maven-plugin to generate container image and k8s manifests
- Commit the generated resources to repo

## Development Environment Setup

Steps to successfully execute Jkube addon:
- Setup Minikube and install tackle - [start-minikube action](https://github.com/konveyor/tackle2-operator/tree/main/.github/actions/start-minikube), [install-tackle action](https://github.com/konveyor/tackle2-operator/tree/main/.github/actions/install-tackle), [install-tackle](https://github.com/konveyor/tackle2-operator/blob/main/hack/install-tackle.sh)

- Forward port for tackle-ui service
  ```bash
  kubectl port-forward service/tackle-ui 8080:8080 -n konveyor-tackle
  ```
- Populate required data in tackle, including application against which addon will be executed - [reference](https://github.com/konveyor/tackle2-hub/tree/main/hack/add), [jkube-data](https://github.com/konveyor/tackle2-addon-jkube/blob/main/hack/task.sh)

- Add Jkube addon to tackle - [jkube addon](https://github.com/konveyor/tackle2-addon-jkube/blob/main/hack/addon.yml)
  ```bash
  kubectly apply -f addon.yml
  ```
- Make curl request to Jkube addon - [request]()
  ```bash
  host="${HOST:-localhost:8080/hub}"

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
  ```
  
<strong>Reference - </strong> End-to-End setup to run an addon - [Test Windup Addon](https://github.com/konveyor/tackle2-addon-windup/blob/84174448d3d7cd2abc7ba6ab27e66a55890b9061/.github/workflows/test-windup.yml#L36-L53)  
  
 ## Output
* Deployment artifacts
  * Dockerfile
  * Kubernetes Yamls
  
# Jkube Addon Demonstration

https://user-images.githubusercontent.com/20452032/201997511-3b2a86cb-dfd0-4a0a-9c90-bbc0fdbd987e.mov
