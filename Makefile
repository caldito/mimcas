# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GODEPS=$(GOCMD) get
#GOTEST=$(GOCMD) test
BINARY_NAME=bin/go-memcached
SOURCE_NAME=cmd/go-memcached/main.go
VERSION=v0.1.0

all: build

build: 
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -v $(SOURCE_NAME)

run: build
	./bin/go-memcached

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

fmt:
	gofmt -w .

deps:
	$(GODEPS) -d ./...

build-docker: build
	docker build . -t pablogcaldito/go-memcached:$(VERSION)
