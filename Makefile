# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GODEPS=$(GOCMD) get
#GOTEST=$(GOCMD) test
SERVER_PROGRAM_NAME=mimcas-server
SERVER_BINARY_NAME=bin/$(SERVER_PROGRAM_NAME)
SERVER_SOURCE_NAME=cmd/$(SERVER_PROGRAM_NAME)/main.go
CLI_PROGRAM_NAME=mimcas-cli
CLI_BINARY_NAME=bin/$(CLI_PROGRAM_NAME)
CLI_SOURCE_NAME=cmd/$(CLI_PROGRAM_NAME)/main.go
VERSION=v0.1.0

all: build

build: 
	CGO_ENABLED=0 $(GOBUILD) -o $(SERVER_BINARY_NAME) -v $(SERVER_SOURCE_NAME)
	CGO_ENABLED=0 $(GOBUILD) -o $(CLI_BINARY_NAME) -v $(CLI_SOURCE_NAME)

run: build
	./$(SERVER_BINARY_NAME)

clean: 
	$(GOCLEAN)
	rm -f $(SERVER_BINARY_NAME)

fmt:
	gofmt -w .

deps:
	$(GODEPS) -d ./...

build-docker: build
	docker build . -f Dockerfile-server -t pablogcaldito/mimcas-server:$(VERSION)
	docker build . -f Dockerfile-cli -t pablogcaldito/mimcas-cli:$(VERSION)
