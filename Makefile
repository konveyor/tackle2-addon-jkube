GOBIN ?= ${GOPATH}/bin

all: cmd

fmt:
	go fmt ./...

vet:
	go vet ./...

cmd: fmt vet
	go build -ldflags="-w -s" -o bin/addon github.com/mundra-ankur/tackle2-addon-jkube/cmd