DOCKER_REPO  = mhd999/appsync-resource
TARGET      ?= linux
ARCH        ?= amd64
SRC          = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

default: test

build: test
	@echo "== Build =="
	env GOOS=$(TARGET) GOARCH=$(ARCH) go build -ldflags="-s -w" -o out main.go

test:
	@echo "== Test =="
	gofmt -s -l -w $(SRC)

docker:
	@echo "== Docker build =="
	docker build -t $(DOCKER_REPO) .

.PHONY: default generate build test docker