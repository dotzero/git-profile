GOBIN ?= $(shell go env GOPATH)/bin
PKG = github.com/dotzero/git-profile
BIN := git-profile

VERSION := 1.3.0
HASH := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%FT%T%z)

LDFLAGS := "-s -w \
	-X main.Version=$(VERSION) \
	-X main.CommitHash=$(HASH) \
	-X main.CompileDate=$(DATE)"

all: build

build:
	GO111MODULE=on go build -ldflags=$(LDFLAGS) -o $(GOBIN)/$(BIN)

install:
	GO111MODULE=on go install -ldflags=$(LDFLAGS)

test:
	GO111MODULE=on go test -v ./...

clean:
	if [ -f $(GOBIN)/$(BIN) ] ; then rm -f $(GOBIN)/$(BIN) ; fi

lint:
	GO111MODULE=on golangci-lint run

.PHONY: build install test clean lint
