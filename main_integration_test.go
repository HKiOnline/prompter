package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildBinary compiles the prompter executable into a temporary directory.
// The package directory is the repository root, so "go build ." embeds the
// real version.json.
func buildBinary(t *testing.T, ldflags string) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "prompter")
	args := []string{"build", "-o", bin}
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	args = append(args, ".")
	cmd := exec.Command("go", args...)
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "go build failed: %s", out)
	return bin
}

// writeUnusableConfig creates a HOME with a configuration file that fails
// validation, so any code path that reads configuration must fail.
func writeUnusableConfig(t *testing.T, home string) {
	t.Helper()
	cfgDir := filepath.Join(home, ".config", "prompter")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(cfgDir, "prompter.yaml"),
		[]byte("prompter:\n  transport:\n    type: \"bogus\"\n"), 0o644))
}

func writeValidConfig(t *testing.T, home string) string {
	t.Helper()
	cfgDir := filepath.Join(home, ".config", "prompter")
	promptsDir := filepath.Join(cfgDir, "prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(cfgDir, "prompter.yaml"),
		[]byte(fmt.Sprintf("prompter:\n  transport:\n    type: \"stdio\"\n    streamable_http:\n      port: 8080\n  storage:\n    provider: \"filesystem\"\n    filesystem:\n      prompts_directory: %q\n", promptsDir)), 0o644))
	return promptsDir
}

func dirListing(t *testing.T, root string) map[string]bool {
	t.Helper()
	listing := map[string]bool{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		listing[rel] = true
		return nil
	})
	require.NoError(t, err)
	return listing
}

// TestVersionFlagSubprocess runs the built executable outside the source
// tree, with an altered runtime version.json and an unusable user
// configuration, and verifies both flag forms exit promptly with the
// embedded version on stdout, empty stderr, and no created files.
func TestVersionFlagSubprocess(t *testing.T) {
	bin := buildBinary(t, "")

	workDir := t.TempDir()
	// altered runtime version.json next to where the binary runs: the
	// embedded metadata must win
	require.NoError(t, os.WriteFile(filepath.Join(workDir, "version.json"), []byte(`{"major":9,"minor":9,"patch":9}`), 0o644))

	home := t.TempDir()
	writeUnusableConfig(t, home)
	before := dirListing(t, home)

	for _, arg := range []string{"--version", "-version"} {
		t.Run(arg, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, bin, arg)
			cmd.Dir = workDir
			cmd.Env = append(os.Environ(), "HOME="+home)
			// stdin is the null device: the flag must not wait for stdin
			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			require.NoError(t, err, "exit status / timeout; stderr: %s", stderr.String())
			assert.Equal(t, "prompter 0.8.0 (2026-11-07T13:21:07Z)\n", stdout.String())
			assert.Empty(t, stderr.String())
		})
	}

	// no logs or storage may have been created by the version flag
	assert.Equal(t, before, dirListing(t, home), "version flag must not create files")
}

// TestVersionFlagStampedSubprocess verifies a build-stamped binary reports
// the +build suffix and the stamped build date.
func TestVersionFlagStampedSubprocess(t *testing.T) {
	bin := buildBinary(t, "-X github.com/hkionline/prompter/internal/version.BuildDate=2026-11-07T14:00:00Z -X github.com/hkionline/prompter/internal/version.BuildID=ci.12345678.20261107140000")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "--version")
	cmd.Dir = t.TempDir()
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	assert.Equal(t, "prompter 0.8.0+ci.12345678.20261107140000 (2026-11-07T14:00:00Z)\n", string(out))
}

const initializeRequest = `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}` + "\n"

// initializeServer starts the binary with a valid stdio configuration, sends
// an MCP initialize request, and returns the raw response line.
func initializeServer(t *testing.T, bin string) string {
	t.Helper()

	home := t.TempDir()
	writeValidConfig(t, home)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Dir = t.TempDir()
	cmd.Env = append(os.Environ(), "HOME="+home)
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)
	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	require.NoError(t, cmd.Start())

	lines := make(chan string, 10)
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()

	_, err = stdin.Write([]byte(initializeRequest))
	require.NoError(t, err)

	select {
	case line := <-lines:
		require.NoError(t, stdin.Close())
		_ = cmd.Wait()
		return line
	case <-time.After(15 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatalf("timed out waiting for initialize response; stderr: %s", stderr.String())
		return ""
	}
}

// TestMCPInitializeReportsCoreVersion verifies serverInfo.version is the
// MAJOR.MINOR.PATCH core version (no +build suffix) and that stdout contains
// only protocol output.
func TestMCPInitializeReportsCoreVersion(t *testing.T) {
	bin := buildBinary(t, "-X github.com/hkionline/prompter/internal/version.BuildDate=2026-11-07T14:00:00Z -X github.com/hkionline/prompter/internal/version.BuildID=lb.20261107140000")

	line := initializeServer(t, bin)

	// no version banner: the line must be pure JSON-RPC
	var response struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			ProtocolVersion string `json:"protocolVersion"`
			ServerInfo      struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"serverInfo"`
		} `json:"result"`
		Error any `json:"error"`
	}
	require.NoError(t, json.Unmarshal([]byte(line), &response), "stdout line must be JSON: %q", line)
	assert.Equal(t, "2.0", response.JSONRPC)
	assert.Nil(t, response.Error)
	assert.Equal(t, "prompter", response.Result.ServerInfo.Name)
	assert.Equal(t, "0.8.0", response.Result.ServerInfo.Version, "serverInfo.version must be the core version without +build")
}

// TestCLIAndMCPVersionsAgree verifies the MCP serverInfo.version equals the
// CLI --version output minus any +build suffix.
func TestCLIAndMCPVersionsAgree(t *testing.T) {
	bin := buildBinary(t, "-X github.com/hkionline/prompter/internal/version.BuildDate=2026-11-07T14:00:00Z -X github.com/hkionline/prompter/internal/version.BuildID=ci.12345678.20261107140000")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, bin, "--version").Output()
	require.NoError(t, err)
	cliLine := strings.TrimSpace(string(out))
	require.True(t, strings.HasPrefix(cliLine, "prompter "), cliLine)
	cliSemVer := strings.Fields(cliLine)[1]
	require.Contains(t, cliSemVer, "+ci.12345678.", "stamped binary CLI must show the build suffix")
	wantCore := strings.SplitN(cliSemVer, "+", 2)[0]

	line := initializeServer(t, bin)
	var response struct {
		Result struct {
			ServerInfo struct {
				Version string `json:"version"`
			} `json:"serverInfo"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal([]byte(line), &response))
	assert.Equal(t, wantCore, response.Result.ServerInfo.Version)
}
