BINARY_NAME=prompter
BIN_DIR=bin
SCRATCH_DIR=scratch

.PHONY: all build setup test run clean

all: build test
  
setup:
	mkdir -p ${BIN_DIR} ${SCRATCH_DIR}
  
build: setup
	go build -o ./${BIN_DIR}/${BINARY_NAME} main.go
  
test: build
	go test ./... > /dev/null 2>&1
	go test -coverprofile=coverage.out ./internal/tools >> ./${SCRATCH_DIR}/test_report.json 2>&1 || true
	cat ./${SCRATCH_DIR}/test_report.json | jq '. | select(.Action == "fail")' 2>/dev/null || true
	tests/test_server.sh
  
run: build test
	./${BIN_DIR}/${BINARY_NAME}
  
clean:
	go clean
	rm -rf ./${BIN_DIR} ./${SCRATCH_DIR}
