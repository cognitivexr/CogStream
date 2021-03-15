.PHONY: mediator engines clean all

BUILD_DIR=dist/

all: mediator engines

mediator:
	$(MAKE) -C mediator all BUILD_DIR=../$(BUILD_DIR)

engines:
	$(MAKE) -C engines/engines-go all BUILD_DIR=../../$(BUILD_DIR)engines/

clean:
	rm -rf $(BUILD_DIR)
