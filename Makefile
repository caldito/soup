# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
#GOTEST=$(GOCMD) test
BINARY_NAME=bin/soup
SOURCE_NAME=cmd/soup/main.go

all: build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v $(SOURCE_NAME)

#test: 
#	$(GOTEST) -v ./...

run: build
	./$(BINARY_NAME) -repo https://github.com/caldito/soup-test.git

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

fmt:
	gofmt -w $(SOURCE_NAME)

build-podman: build
	podman build . -t soup

build-docker: build
	docker build . -t soup

run-podman: build-podman
	podman run -it --entrypoint /bin/soup soup -repo https://github.com/caldito/soup-test.git

run-docker: build-docker
	docker run -it --entrypoint /bin/soup soup -repo https://github.com/caldito/soup-test.git
