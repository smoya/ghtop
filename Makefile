SOURCES          ?= $(shell find . -name "*.go" | grep -v vendor/)
PACKAGES         ?= $(shell go list ./...)
GOTOOLS          ?= github.com/kardianos/govendor \
                    github.com/GeertJohan/fgt \
					golang.org/x/tools/cmd/goimports \
					github.com/golang/lint/golint \
					github.com/kisielk/errcheck \
					honnef.co/go/tools/cmd/staticcheck

deps: tools
	govendor sync

test: deps
	go test $(PACKAGES)

test-e2e: deps
	go test github.com/smoya/ghtop/pkg/e2e -v

lint: tools
	fgt go fmt $(PACKAGES)
	fgt goimports -w $(SOURCES)
	echo $(PACKAGES) | xargs -L1 fgt golint
	fgt go vet $(PACKAGES)
	fgt errcheck -ignore Close $(PACKAGES)
	staticcheck $(PACKAGES)
.SILENT: lint

check: lint test test-e2e

tools:
	go get $(GOTOOLS)
.SILENT: tools

tools-update:
	go get -u $(GOTOOLS)

build: deps
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/ghtop
.PHONY: build
