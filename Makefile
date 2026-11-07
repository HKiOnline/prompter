BINARY_NAME=prompter
BIN_DIR=bin
SCRATCH_DIR=scratch

# Build-time version stamps (UTC): date as ISO 8601, build ID as
# lb.{YYYYMMDDHHMMSS} for local builds.
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILD_ID=lb.$(shell date -u +%Y%m%d%H%M%S)
VERSION_PKG=github.com/hkionline/prompter/internal/version
LDFLAGS=-X ${VERSION_PKG}.BuildDate=${BUILD_DATE} -X ${VERSION_PKG}.BuildID=${BUILD_ID}

.PHONY: all build setup test run clean coverage

all: build test
  
setup:
	@mkdir -p ${BIN_DIR} ${SCRATCH_DIR}
   
build: setup
	@go build -tags sqlite -ldflags "${LDFLAGS}" -o ./${BIN_DIR}/${BINARY_NAME} main.go
  
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
