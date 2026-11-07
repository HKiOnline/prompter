package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testEmbedded = `{"major":0,"minor":8,"patch":0,"preRelease":"","build":"","date":"2026-11-07T13:21:07Z"}`

func TestRunVersionFlag(t *testing.T) {
	for _, arg := range []string{"--version", "-version"} {
		t.Run(arg, func(t *testing.T) {
			var out, errBuf bytes.Buffer
			code := run([]string{arg}, []byte(testEmbedded), &out, &errBuf)
			assert.Equal(t, 0, code)
			assert.Equal(t, "prompter 0.8.0 (2026-11-07T13:21:07Z)\n", out.String())
			assert.Empty(t, errBuf.String())
		})
	}
}

func TestRunUnexpectedArgument(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run([]string{"bogus"}, []byte(testEmbedded), &out, &errBuf)
	assert.NotEqual(t, 0, code)
	assert.Contains(t, errBuf.String(), `unexpected argument "bogus"`)
	assert.Empty(t, out.String())
}

func TestRunHelp(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run([]string{"-h"}, []byte(testEmbedded), &out, &errBuf)
	assert.Equal(t, 0, code)
	assert.Contains(t, errBuf.String(), "version")
	assert.Empty(t, out.String())
}

func TestRunInvalidMetadata(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run([]string{"--version"}, []byte(`{"major":-1,"minor":0,"patch":0,"preRelease":"","build":"","date":"2026-11-07T13:21:07Z"}`), &out, &errBuf)
	assert.NotEqual(t, 0, code)
	assert.Contains(t, errBuf.String(), "invalid embedded version.json")
	// no version output or protocol data on stdout
	assert.Empty(t, out.String())
}

func TestRunVersionFalseProceedsToStartup(t *testing.T) {
	// --version=false must proceed to ordinary startup; with an unusable
	// configuration this fails during configuration, proving the version
	// flag did not short-circuit startup.
	home := t.TempDir()
	cfgDir := filepath.Join(home, ".config", "prompter")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "prompter.yaml"), []byte("prompter:\n  transport:\n    type: \"bogus\"\n"), 0o644))
	t.Setenv("HOME", home)

	var out, errBuf bytes.Buffer
	code := run([]string{"--version=false"}, []byte(testEmbedded), &out, &errBuf)
	assert.NotEqual(t, 0, code)
	assert.Contains(t, errBuf.String(), "configuration failure")
	assert.Empty(t, out.String())
}
