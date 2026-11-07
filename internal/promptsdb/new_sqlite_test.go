//go:build sqlite

package promptsdb_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hkionline/prompter/internal/promptsdb"

	// blank import so the sqlite package's init() registers its constructor
	_ "github.com/hkionline/prompter/internal/promptsdb/sqlite"
)

func TestNewSQLiteProvider(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "new_sqlite_")
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

	provider, err := promptsdb.New(promptsdb.SQLITE_PROVIDER, config, "")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}

	// Verify CRUD works through the interface
	prompt := promptsdb.Prompt{
		Name:    "new-test",
		Title:   "New Test",
		Content: "Test content",
		Tags:    []string{"test"},
	}
	if err := provider.Create(prompt); err != nil {
		t.Fatalf("Create via interface failed: %v", err)
	}

	got, err := provider.Read("new-test")
	if err != nil {
		t.Fatalf("Read via interface failed: %v", err)
	}
	if got.Name != "new-test" {
		t.Errorf("Expected new-test, got %s", got.Name)
	}
}
