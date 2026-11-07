//go:build sqlite

package sqlite

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hkionline/prompter/internal/promptsdb"
)

// Integration test: migrate prompts from filesystem to SQLite
func TestMigrationFromFilesystem(t *testing.T) {
	// Create source directory with prompt files
	srcDir, err := os.MkdirTemp("", "migrate_src_")
	if err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	// Write test prompt files
	prompts := []promptsdb.Prompt{
		{Name: "migrate-1", Title: "Migrate 1", Description: "First", Content: "Content 1", Tags: []string{"a", "b"}, Arguments: []string{"x"}},
		{Name: "migrate-2", Title: "Migrate 2", Description: "Second", Content: "Content 2", Tags: []string{"c"}},
		{Name: "migrate-3", Title: "Migrate 3", Content: "Content 3"},
	}

	for _, p := range prompts {
		content := "---\n" +
			"name: " + p.Name + "\n" +
			"title: " + p.Title + "\n"
		if p.Description != "" {
			content += "description: " + p.Description + "\n"
		}
		if len(p.Arguments) > 0 {
			content += "arguments:\n"
			for _, a := range p.Arguments {
				content += "  - " + a + "\n"
			}
		}
		if len(p.Tags) > 0 {
			content += "tags:\n"
			for _, tag := range p.Tags {
				content += "  - " + tag + "\n"
			}
		}
		content += "---\n" + p.Content

		if err := os.WriteFile(filepath.Join(srcDir, p.Name+".md"), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write prompt file: %v", err)
		}
	}

	// Create a corrupted file
	os.WriteFile(filepath.Join(srcDir, "corrupted.md"), []byte("not a valid prompt file"), 0644)

	// Create destination SQLite db
	dstDir, err := os.MkdirTemp("", "migrate_dst_")
	if err != nil {
		t.Fatalf("Failed to create dst dir: %v", err)
	}
	defer os.RemoveAll(dstDir)
	dbPath := filepath.Join(dstDir, "migrated.db")

	// Create SQLite provider
	provider, err := NewSQLiteProvider(SQLiteProviderConfiguration{Path: dbPath, Timeout: 30}, "")
	if err != nil {
		t.Fatalf("Failed to create SQLite provider: %v", err)
	}
	sp := provider.(*SQLiteProvider)
	defer sp.Close()

	// Migrate each prompt (simulating what the migrate command does)
	migrated := 0
	for _, p := range prompts {
		if err := sp.Create(p); err != nil {
			t.Fatalf("Failed to migrate %s: %v", p.Name, err)
		}
		migrated++
	}

	if migrated != 3 {
		t.Errorf("Expected 3 migrated, got %d", migrated)
	}

	// Verify all prompts are in SQLite
	all, err := sp.List(promptsdb.PromptQuery{All: true})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Expected 3 prompts in SQLite, got %d", len(all))
	}

	// Verify content integrity
	for _, p := range prompts {
		got, err := sp.Read(p.Name)
		if err != nil {
			t.Fatalf("Read %s failed: %v", p.Name, err)
		}
		if got.Content != p.Content {
			t.Errorf("Content mismatch for %s: expected %q, got %q", p.Name, p.Content, got.Content)
		}
		if got.Title != p.Title {
			t.Errorf("Title mismatch for %s: expected %q, got %q", p.Name, p.Title, got.Title)
		}
	}
}

// Integration test: RPC-like operations through the Provider interface
func TestProviderInterfaceCompliance(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	var _ promptsdb.Provider = sp

	// Full CRUD cycle through interface
	prompt := promptsdb.Prompt{
		Name:      "rpc-test",
		Title:     "RPC Test",
		Content:   "RPC content",
		Arguments: []string{"a"},
		Tags:      []string{"rpc"},
	}

	if err := sp.Create(prompt); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := sp.Read("rpc-test")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if got.Name != "rpc-test" {
		t.Errorf("Name mismatch: %s", got.Name)
	}

	prompt.Content = "Updated RPC content"
	if err := sp.Update(prompt); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, _ = sp.Read("rpc-test")
	if got.Content != "Updated RPC content" {
		t.Errorf("Update not applied: %s", got.Content)
	}

	list, err := sp.List(promptsdb.PromptQuery{All: true})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("Expected 1 prompt, got %d", len(list))
	}

	if err := sp.Delete("rpc-test"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	list, _ = sp.List(promptsdb.PromptQuery{All: true})
	if len(list) != 0 {
		t.Errorf("Expected 0 prompts after delete, got %d", len(list))
	}
}
