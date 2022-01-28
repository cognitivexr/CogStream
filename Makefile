.PHONY: mediator engines engines-py engines-go clean all 

BUILD_DIR=dist/

all: mediator engines

mediator:
	$(MAKE) -C mediator all BUILD_DIR=../$(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(MAKE) -C mediator all BUILD_DIR=../$(BUILD_DIR)

engines-go:
	GOARCH=amd64 $(MAKE) -C engines/engines-go all BUILD_DIR=../../$(BUILD_DIR)engines/

engines-py:
	cd engines/engines-py; $(MAKE) install

engines: engines-go engines-py

clean:
	rm -rf $(BUILD_DIR)
	cd engines/engines-py; $(MAKE) clean 

