.PHONY: default build image test stop clean-image clean

BINARY = kuberbs


GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =

default: build test

build:
	"$(GOCMD)" build ${GOFLAGS} ${LDFLAGS} -o "${BINARY}"

image:
	@docker build -t "${BINARY}" -f Dockerfile .

test:
	"$(GOCMD)" test -race -v $(shell go list ./... | grep -v '/vendor/')

stop:
	@docker stop "${BINARY}" || true # Do not fail if container does not exist

clean-image: stop
	@docker rmi "${BINARY}" || true # Do not fail if image does not exist

clean:
	"$(GOCMD)" clean -i
