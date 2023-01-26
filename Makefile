PROJECTNAME	= wfmon

SHELL := /bin/bash

GOBASE		= $(shell pwd)
GOBIN		= $(GOBASE)/bin
GOMAIN		= ${GOBASE}/cmd/main.go
GOPKG		= ${GOBASE}/pkg

## build: Build binary files.
build: clean
# MacOS
	GOOS="darwin" go build -race -o $(GOBIN)/$(PROJECTNAME) ${GOMAIN}

## run: Run
run:
	$(eval -include .env)
	$(eval export)
	env
	go run ${GOMAIN}

## clean: Clean build files.
clean:
	go clean
	rm -rf ${GOBIN}

## dep: Downloads modules dependencies.
dep:
	go mod tidy -v
	go mod download -x

## lint: Runs `golangci-lint` internally.
lint:
	golangci-lint run

## test: Runs tests.
test:
	go fmt $(shell go list ./... | grep -v /vendor/)
	go vet $(shell go list ./... | grep -v /vendor/)
	go test -race $(shell go list ./... | grep -v /vendor/)

test-json:
	go test -race $(shell go list ./... | grep -v /vendor/) -json > ./test-results.json

all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo