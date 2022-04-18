# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GODEPS=$(GOCMD) get
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
BINARY_NAME=bin/soup
SOURCE_NAME=cmd/soup/main.go

all: build

build:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -v $(SOURCE_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test: build
	$(GOTEST) ./...

fmt:
	$(GOFMT) ./...

deps:
	$(GODEPS) -d ./...

build-docker: build
	docker build . -t pablogcaldito/soup:$(VERSION)
