GOBIN ?= $(shell go env GOPATH)/bin
PKG = github.com/dotzero/git-profile
BIN := git-profile

VERSION := 1.8.1
HASH := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%FT%T%z)

LDFLAGS := "-s -w \
	-X main.Version=$(VERSION) \
	-X main.CommitHash=$(HASH) \
	-X main.CompileDate=$(DATE)"

all: build

build:
	go build -ldflags=$(LDFLAGS) -o $(GOBIN)/$(BIN)

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	if [ -f $(GOBIN)/$(BIN) ] ; then rm -f $(GOBIN)/$(BIN) ; fi

dist-check:
	goreleaser release --snapshot --clean

dist-release:
	goreleaser release --clean

.PHONY: all build test lint clean dist-check dist-release
