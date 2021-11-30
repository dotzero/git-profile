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
	go build -ldflags=$(LDFLAGS) -o $(GOBIN)/$(BIN)

install:
	go install -ldflags=$(LDFLAGS)

test:
	go test -v ./...

clean:
	if [ -f $(GOBIN)/$(BIN) ] ; then rm -f $(GOBIN)/$(BIN) ; fi

lint:
	golangci-lint run

dist-check:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: build install test clean lint dist-check
