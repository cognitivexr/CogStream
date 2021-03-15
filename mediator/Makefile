GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BUILD_DIR=bin/

all: mediator

mediator:
	$(GOBUILD) -o $(BUILD_DIR) ./cmd/mediator

clean:
	rm -rf $(BUILD_DIR)
