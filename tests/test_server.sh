#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
export PROMPTER_TEST_BINARY="${PROMPTER_TEST_BINARY:-$ROOT/bin/prompter}"
exec go test ./tests -run '^TestReleaseBinaryMCP$' -count=1 -v
