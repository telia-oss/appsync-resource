DOCKER_REPO  = mhd999/appsync-resource
TARGET      ?= linux
ARCH        ?= amd64
SRC          = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: test

install: 
	curl -L -s https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 -o $(GOPATH)/bin/dep
	chmod +x $(GOPATH)/bin/dep
	dep ensure

build: test install	
	@echo "== Build =="
	env GOOS=$(TARGET) GOARCH=$(ARCH) go build -ldflags="-s -w" -o in main.go
	env GOOS=$(TARGET) GOARCH=$(ARCH) go build -ldflags="-s -w" -o check main.go
	env GOOS=$(TARGET) GOARCH=$(ARCH) go build -ldflags="-s -w" -o out main.go

test:
	@echo "== Test =="
	gofmt -s -l -w $(SRC)

docker:
	@echo "== Docker build =="
	docker build -t $(DOCKER_REPO) .

.PHONY: default generate build test docker
