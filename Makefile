# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GODEPS=$(GOCMD) get
#GOTEST=$(GOCMD) test
BINARY_NAME=bin/kv-store
SOURCE_NAME=cmd/kv-store/main.go
VERSION=v0.1.0

all: build

build: 
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -v $(SOURCE_NAME)

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

fmt:
	gofmt -w .

deps:
	$(GODEPS) -d ./...

build-docker: build
	docker build . -t pablogcaldito/kv-store:$(VERSION)
