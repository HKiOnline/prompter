BINARY_NAME=prompter
BIN_DIR=bin
SCRATCH_DIR=scratch

.PHONY: all build setup test run clean coverage

all: build test
  
setup:
	@mkdir -p ${BIN_DIR} ${SCRATCH_DIR}
   
build: setup
	@go build -o ./${BIN_DIR}/${BINARY_NAME} main.go
  
test: build
	@go test ./... > /dev/null 2>&1
	@go test ./... > ./${SCRATCH_DIR}/test_report.txt 2>&1 || true
	@cat ./${SCRATCH_DIR}/test_report.txt
	@grep -q "FAIL" ./${SCRATCH_DIR}/test_report.txt || true
	@tests/test_server.sh
  
run: build test
	./${BIN_DIR}/${BINARY_NAME}
  
clean:
	go clean
	rm -rf ./${BIN_DIR} ./${SCRATCH_DIR}

coverage:
	@echo "Generating coverage report..."
	@go test ./... -coverprofile=./${SCRATCH_DIR}/coverage.out 2>&1 | grep "coverage:" || true
	@if [ -f ./${SCRATCH_DIR}/coverage.out ]; then echo "Overall coverage summary:" && go tool cover -func=./${SCRATCH_DIR}/coverage.out | tail -1; else echo "No coverage data generated"; fi
