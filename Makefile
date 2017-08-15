GOBIN ?= $(shell go env GOPATH)/bin
PKG = github.com/dotzero/git-profile
BIN := git-profile

VERSION := 1.0.0
HASH := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%FT%T%z)

LDFLAGS := "-s -w \
	-X $(PKG)/cmd.Version=$(VERSION) \
	-X $(PKG)/cmd.CommitHash=$(HASH) \
	-X $(PKG)/cmd.CompileDate=$(DATE)"

all: build

build: fmt vet
	go build -ldflags=$(LDFLAGS) -o $(GOBIN)/$(BIN)

install:
	go install -ldflags=$(LDFLAGS)

test:
	go test -v $(go list ./... | grep -v /vendor/)

clean:
	if [ -f $(GOBIN)/$(BIN) ] ; then rm -f $(GOBIN)/$(BIN) ; fi

fmt:
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

vet:
	go vet ./...

.PHONY: build install test clean fmt vet
