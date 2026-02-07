#!/bin/bash
# Test script to communicate with the prompter MCP server using test data files

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DATA_DIR="$SCRIPT_DIR/data"
CWD=$(pwd)

# Ensure describe_tampere.md exists in prompts directory
PROMPTS_DIR="$HOME/.config/prompter/prompts"
mkdir -p "$PROMPTS_DIR"
if [ ! -f "$PROMPTS_DIR/describe_tampere.md" ]; then
    echo "Copying describe_tampere.md to prompts directory..."
    cp "$DATA_DIR/tampere_prompt.md" "$PROMPTS_DIR/describe_tampere.md"
fi

# Create scratch directory if it doesn't exist
SCRATCH_DIR="$CWD/scratch"
mkdir -p "$SCRATCH_DIR"

# Use bin directory like Makefile
BIN_DIR="$CWD/bin"
BINARY_NAME="prompter"

echo "Starting prompter MCP server test..."
echo "====================================="

# Send all requests as a single stream to the server
echo ""
echo "Sending test sequence..."
echo "-------------------------"
(
    cat "$DATA_DIR/client_initialize_call.json" &&
    sleep 0.5 &&
    cat "$DATA_DIR/client_initialized_call.json" &&
    sleep 0.5 &&
    cat "$DATA_DIR/client_prompts_list_call.json" &&
    sleep 0.5 &&
    cat "$DATA_DIR/client_prompts_get_tampere.json"
) | "$BIN_DIR/$BINARY_NAME" > "$SCRATCH_DIR/full_response.txt" 2>&1

# Wait for responses
echo ""
echo "Test Results:"
echo "-------------"
cat "$SCRATCH_DIR/full_response.txt"

# Clean up
rm -f "$SCRATCH_DIR/prompter_input.txt" "$SCRATCH_DIR/prompter_output.txt" "$SCRATCH_DIR/full_response.txt"

echo ""
echo "Test completed."
