GOBIN ?= $(shell go env GOPATH)/bin
PKG = github.com/dotzero/git-profile
BIN := git-profile

HASH := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%FT%T%z)

LDFLAGS := "-s -w \
	$(if $(strip $(VERSION)),-X main.Version=$(VERSION)) \
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

require-version:
	@test -n "$(strip $(VERSION))" || { echo "VERSION is required (e.g. make dist-check VERSION=1.9.2)" >&2; exit 1; }

dist-check: require-version
	GORELEASER_CURRENT_TAG=v$(VERSION) goreleaser release --snapshot --clean

dist-release: require-version
	goreleaser release --clean

.PHONY: all build test lint clean require-version dist-check dist-release
