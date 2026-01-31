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
	go test -json ./... > ./${SCRATCH_DIR}/test_report.json && cat ./${SCRATCH_DIR}/test_report.json | jq '. | select(.Action == "fail")'
	tests/test_server.sh
  
run: build test
	./${BIN_DIR}/${BINARY_NAME}
  
clean:
	go clean
	rm -rf ./${BIN_DIR} ./${SCRATCH_DIR}
