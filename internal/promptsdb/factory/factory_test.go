//go:build sqlite

package factory

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hkionline/prompter/internal/promptsdb"
)

func TestFactoryFilesystemProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "factory_fs_")
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

	provider, err := New(promptsdb.FILE_SYSTEM_PROVIDER, config, "")
	if err != nil {
		t.Fatalf("Factory returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}

func TestFactorySQLiteProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "factory_sqlite_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	config := promptsdb.ProviderConfiguration{
		Provider: promptsdb.SQLITE_PROVIDER,
		SQLite: promptsdb.SQLiteConfiguration{
			Path:    dbPath,
			Timeout: 30,
		},
	}

	provider, err := New(promptsdb.SQLITE_PROVIDER, config, "")
	if err != nil {
		t.Fatalf("Factory returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}

	// Verify CRUD works through the interface
	prompt := promptsdb.Prompt{
		Name:    "factory-test",
		Title:   "Factory Test",
		Content: "Test content",
		Tags:    []string{"test"},
	}
	if err := provider.Create(prompt); err != nil {
		t.Fatalf("Create via interface failed: %v", err)
	}

	got, err := provider.Read("factory-test")
	if err != nil {
		t.Fatalf("Read via interface failed: %v", err)
	}
	if got.Name != "factory-test" {
		t.Errorf("Expected factory-test, got %s", got.Name)
	}
}

func TestFactoryDefaultFallback(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "factory_default_")
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
	provider, err := New("unknown", config, "")
	if err != nil {
		t.Fatalf("Factory returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}

func TestFactoryEmptyProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "factory_empty_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := promptsdb.ProviderConfiguration{
		Filesystem: promptsdb.FsProviderConfiguration{
			Directory: tempDir,
		},
	}

	provider, err := New("", config, "")
	if err != nil {
		t.Fatalf("Factory returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}
