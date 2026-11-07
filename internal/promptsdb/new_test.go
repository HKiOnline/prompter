package promptsdb_test

import (
	"os"
	"testing"

	"github.com/hkionline/prompter/internal/promptsdb"
)

func TestNewFilesystemProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "new_fs_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := promptsdb.ProviderConfiguration{
		Provider: promptsdb.FILE_SYSTEM_PROVIDER,
		Filesystem: promptsdb.FsProviderConfiguration{
			Directory: tempDir,
		},
	}

	provider, err := promptsdb.New(promptsdb.FILE_SYSTEM_PROVIDER, config, "")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}

func TestNewDefaultFallback(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "new_default_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := promptsdb.ProviderConfiguration{
		Provider: "unknown",
		Filesystem: promptsdb.FsProviderConfiguration{
			Directory: tempDir,
		},
	}

	// Unknown provider should fall back to filesystem
	provider, err := promptsdb.New("unknown", config, "")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}

func TestNewEmptyProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "new_empty_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := promptsdb.ProviderConfiguration{
		Filesystem: promptsdb.FsProviderConfiguration{
			Directory: tempDir,
		},
	}

	provider, err := promptsdb.New("", config, "")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}
