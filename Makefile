GO_BIN := go
GOFMT_BIN := gofmt

BINARY_NAME := changelog

build:
	$(GO_BIN) build -o $(BINARY_NAME)

test: test-unit test-format

test-unit:
	@echo "Unit Test"
	$(GO_BIN) test ./...

test-format:
	@echo "Test Format"
	@result=$$($(GOFMT_BIN) -l .);\
	if [[ ! -z $${result} ]]; then \
		echo $${result}; \
		exit 1;\
	fi

clean:
	rm $(BINARY_NAME)

