# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GODEPS=$(GOCMD) get
#GOTEST=$(GOCMD) test
BINARY_NAME=bin/soup
SOURCE_NAME=cmd/soup/main.go
VERSION=0.2.0

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

build-podman: build
	podman build . -t pablogcaldito/soup:$(VERSION)

build-docker: build
	docker build . -t pablogcaldito/soup:$(VERSION)
