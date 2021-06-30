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

run: build
	./$(BINARY_NAME) -repo https://github.com/caldito/soup-test.git

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

fmt:
	gofmt -w $(SOURCE_NAME)

deps:
	$(GODEPS) -d ./...

build-podman: build
	podman build . -t pablogcaldito/soup:$(VERSION)

build-docker: build
	docker build . -t pablogcaldito/soup:$(VERSION)

test-podman: build-podman
	podman run -it --entrypoint /bin/soup pablogcaldito/soup:$(VERSION) -repo https://github.com/caldito/soup-test.git

test-docker: build-docker
	docker run -it --entrypoint /bin/soup pablogcaldito/soup:$(VERSION) -repo https://github.com/caldito/soup-test.git
