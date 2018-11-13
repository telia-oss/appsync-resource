DOCKER_REPO  = teliaoss/appsync-resource
TARGET      ?= linux
ARCH        ?= amd64
SRC          = $(shell find . -type f -name '*.go' -not -path "./vendor/*")


default: test

install: 
	curl -L -s https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 -o $(GOPATH)/bin/dep
	chmod +x $(GOPATH)/bin/dep
	dep ensure -v

build: test install
	@echo "== Build =="
	env GOOS=$(TARGET) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/in in/in.go
	env GOOS=$(TARGET) GOARCH=$(ARCH) CGO_ENABLED=0  go build -ldflags="-s -w" -o bin/check check/check.go
	env GOOS=$(TARGET) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/out main.go

test:
	@echo "== Test =="
	gofmt -s -l -w $(SRC)

docker:
	@echo "== Docker build =="
	docker build -t $(DOCKER_REPO):dev .

.PHONY: default generate build test docker
