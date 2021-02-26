# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
#GOTEST=$(GOCMD) test
#GOGET=$(GOCMD) get
BINARY_NAME=soup
SOURCE_NAME=cmd/soup/main.go

all: build

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v $(SOURCE_NAME)

#test: 
#	$(GOTEST) -v ./...

run: build
	./$(BINARY_NAME)

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

#deps:
#	$(GOGET) github.com/markbates/goth
#	$(GOGET) github.com/markbates/pop
