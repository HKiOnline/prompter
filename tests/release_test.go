package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestReleaseBinaryMCP exercises the actual Makefile or GoReleaser binary.
// Requests are sequenced by responses, not sleeps. A deadline and cleanup
// guarantee that missing responses and failures cannot leave a server running.
func TestReleaseBinaryMCP(t *testing.T) {
	bin := os.Getenv("PROMPTER_TEST_BINARY")
	if bin == "" {
		t.Skip("set PROMPTER_TEST_BINARY or run make test to test a built binary")
	}
	home := t.TempDir()
	configDir := filepath.Join(home, ".config", "prompter")
	promptsDir := filepath.Join(configDir, "prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0o755))
	prompt, err := os.ReadFile("data/tampere_prompt.md")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(promptsDir, "describe_tampere.md"), prompt, 0o644))
	config := fmt.Sprintf("prompter:\n  transport:\n    type: stdio\n  storage:\n    provider: filesystem\n    filesystem:\n      prompts_directory: %q\n", promptsDir)
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "prompter.yaml"), []byte(config), 0o644))

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Dir = home
	cmd.Env = append(os.Environ(), "HOME="+home, "USERPROFILE="+home,
		"PROMPTER_STORAGE_PROVIDER=filesystem", "PROMPTER_LOG_FILE="+filepath.Join(home, "prompter.log"))
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		_ = stdin.Close()
		cancel()
		_ = cmd.Wait()
		if stderr.Len() > 0 {
			t.Logf("server stderr: %s", stderr.String())
		}
	})
	decoder := json.NewDecoder(stdout)
	send := func(file string) {
		t.Helper()
		data, err := os.ReadFile(filepath.Join("data", file))
		require.NoError(t, err)
		_, err = fmt.Fprintln(stdin, string(data))
		require.NoError(t, err)
	}
	receive := func(id int) json.RawMessage {
		t.Helper()
		var response struct {
			JSONRPC string          `json:"jsonrpc"`
			ID      int             `json:"id"`
			Result  json.RawMessage `json:"result"`
			Error   json.RawMessage `json:"error"`
		}
		require.NoError(t, decoder.Decode(&response), "missing or malformed response for request %d", id)
		require.Equal(t, "2.0", response.JSONRPC)
		require.Equal(t, id, response.ID)
		require.Empty(t, response.Error, "JSON-RPC error: %s", response.Error)
		require.NotEmpty(t, response.Result)
		require.NotEqual(t, "null", string(response.Result))
		return response.Result
	}

	send("client_initialize_call.json")
	var initialized struct {
		ServerInfo struct{ Name, Version string }
	}
	require.NoError(t, json.Unmarshal(receive(1), &initialized))
	require.Equal(t, "prompter", initialized.ServerInfo.Name)
	require.NotEmpty(t, initialized.ServerInfo.Version)
	send("client_initialized_call.json")
	send("client_prompts_list_call.json")
	var listed struct{ Prompts []struct{ Name string } }
	require.NoError(t, json.Unmarshal(receive(2), &listed))
	require.Len(t, listed.Prompts, 1)
	require.Equal(t, "describe_tampere", listed.Prompts[0].Name)
	send("client_prompts_get_tampere.json")
	var got struct {
		Messages []struct {
			Content struct{ Type, Text string }
		}
	}
	require.NoError(t, json.Unmarshal(receive(3), &got))
	require.NotEmpty(t, got.Messages)
	require.Equal(t, "text", got.Messages[0].Content.Type)
	require.Contains(t, got.Messages[0].Content.Text, "Tampere")
}
