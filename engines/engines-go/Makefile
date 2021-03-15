.PHONY: plugins/specs

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOBUILD_PLUGIN=$(GOBUILD) -buildmode=plugin

BUILD_DIR=bin/

all: plugins

plugins: plugins/bin plugins/specs

plugins/bin:
	$(GOBUILD_PLUGIN) -o $(BUILD_DIR) ./plugin/...

plugins/specs:
	find ./plugin/ -name "*.spec.json" -type f -exec cp {} $(BUILD_DIR) \;

clean:
	rm -rf $(BUILD_DIR)
