//go:build sqlite

package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/hkionline/prompter/internal/promptsdb"
)

func setupTestProvider(t *testing.T) (*SQLiteProvider, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "sqlite_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "test.db")
	provider, err := NewSQLiteProvider(SQLiteProviderConfiguration{
		Path:        dbPath,
		Timeout:     30,
		ForeignKeys: true,
	}, "")
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create SQLiteProvider: %v", err)
	}

	sp, ok := provider.(*SQLiteProvider)
	if !ok {
		os.RemoveAll(tempDir)
		t.Fatalf("Expected *SQLiteProvider, got %T", provider)
	}

	cleanup := func() {
		sp.Close()
		os.RemoveAll(tempDir)
	}

	return sp, cleanup
}

func testPrompt(name string) promptsdb.Prompt {
	return promptsdb.Prompt{
		Name:        name,
		Title:       "Test " + name,
		Description: "A test prompt",
		Arguments:   []string{"arg1", "arg2"},
		Content:     "Hello {{.arg1}}, you are {{.arg2}}.",
		Tags:        []string{"test", "unit"},
	}
}

// 6.1 Unit tests for SQLiteProvider methods

func TestSQLiteProviderCreate(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := testPrompt("test-create")
	err := sp.Create(prompt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify it can be read back
	read, err := sp.Read("test-create")
	if err != nil {
		t.Fatalf("Read after Create failed: %v", err)
	}
	if read.Name != prompt.Name {
		t.Errorf("Expected name %s, got %s", prompt.Name, read.Name)
	}
	if read.Content != prompt.Content {
		t.Errorf("Expected content %q, got %q", prompt.Content, read.Content)
	}
	if len(read.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(read.Tags))
	}
	if len(read.Arguments) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(read.Arguments))
	}
}

func TestSQLiteProviderCreateDuplicate(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := testPrompt("dup-test")
	if err := sp.Create(prompt); err != nil {
		t.Fatalf("First create failed: %v", err)
	}

	err := sp.Create(prompt)
	if err == nil {
		t.Fatal("Expected error on duplicate create, got nil")
	}
}

func TestSQLiteProviderRead(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := testPrompt("read-test")
	if err := sp.Create(prompt); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := sp.Read("read-test")
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if got.Title != prompt.Title {
		t.Errorf("Expected title %q, got %q", prompt.Title, got.Title)
	}
}

func TestSQLiteProviderReadNotFound(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	_, err := sp.Read("nonexistent")
	if err == nil {
		t.Fatal("Expected error reading nonexistent prompt, got nil")
	}
}

func TestSQLiteProviderUpdate(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := testPrompt("update-test")
	if err := sp.Create(prompt); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	prompt.Title = "Updated Title"
	prompt.Content = "Updated content"
	if err := sp.Update(prompt); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, err := sp.Read("update-test")
	if err != nil {
		t.Fatalf("Read after update failed: %v", err)
	}
	if got.Title != "Updated Title" {
		t.Errorf("Expected updated title, got %q", got.Title)
	}
	if got.Content != "Updated content" {
		t.Errorf("Expected updated content, got %q", got.Content)
	}
}

func TestSQLiteProviderUpdateNotFound(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := testPrompt("nonexistent-update")
	err := sp.Update(prompt)
	if err == nil {
		t.Fatal("Expected error updating nonexistent prompt, got nil")
	}
}

func TestSQLiteProviderDelete(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := testPrompt("delete-test")
	if err := sp.Create(prompt); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := sp.Delete("delete-test"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := sp.Read("delete-test")
	if err == nil {
		t.Fatal("Expected error reading deleted prompt, got nil")
	}
}

func TestSQLiteProviderDeleteNotFound(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	err := sp.Delete("nonexistent-delete")
	if err == nil {
		t.Fatal("Expected error deleting nonexistent prompt, got nil")
	}
}

func TestSQLiteProviderList(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	names := []string{"alpha", "beta", "gamma"}
	for _, n := range names {
		if err := sp.Create(testPrompt(n)); err != nil {
			t.Fatalf("Create %s failed: %v", n, err)
		}
	}

	prompts, err := sp.List(promptsdb.PromptQuery{All: true})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(prompts) != 3 {
		t.Errorf("Expected 3 prompts, got %d", len(prompts))
	}

	// Should be sorted by name
	if prompts[0].Name != "alpha" || prompts[1].Name != "beta" || prompts[2].Name != "gamma" {
		t.Errorf("Expected sorted order alpha,beta,gamma, got %s,%s,%s", prompts[0].Name, prompts[1].Name, prompts[2].Name)
	}
}

func TestSQLiteProviderListEmpty(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompts, err := sp.List(promptsdb.PromptQuery{All: true})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(prompts) != 0 {
		t.Errorf("Expected 0 prompts, got %d", len(prompts))
	}
}

func TestSQLiteProviderListNameStartsWith(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	sp.Create(testPrompt("foo-one"))
	sp.Create(testPrompt("foo-two"))
	sp.Create(testPrompt("bar-one"))

	prompts, err := sp.List(promptsdb.PromptQuery{NameStartsWith: "foo"})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(prompts) != 2 {
		t.Errorf("Expected 2 prompts starting with 'foo', got %d", len(prompts))
	}
}

func TestSQLiteProviderListNameContains(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	sp.Create(testPrompt("hello-world"))
	sp.Create(testPrompt("goodbye-world"))
	sp.Create(testPrompt("hello-again"))

	prompts, err := sp.List(promptsdb.PromptQuery{NameContains: "world"})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(prompts) != 2 {
		t.Errorf("Expected 2 prompts containing 'world', got %d", len(prompts))
	}
}

// 6.2 Test database schema and constraints

func TestSchemaCreated(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	// Verify prompts table exists
	var name string
	err := sp.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='prompts'").Scan(&name)
	if err != nil {
		t.Fatalf("prompts table not found: %v", err)
	}

	// Verify indexes exist
	rows, err := sp.db.Query("SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='prompts'")
	if err != nil {
		t.Fatalf("Failed to query indexes: %v", err)
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var idx string
		rows.Scan(&idx)
		indexes = append(indexes, idx)
	}
	if len(indexes) < 2 {
		t.Errorf("Expected at least 2 indexes, got %d", len(indexes))
	}
}

func TestSchemaVersion(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	var version int
	err := sp.db.QueryRow("SELECT MAX(version) FROM schema_version").Scan(&version)
	if err != nil {
		t.Fatalf("Failed to read schema version: %v", err)
	}
	if version != 1 {
		t.Errorf("Expected schema version 1, got %d", version)
	}
}

func TestUniqueNameConstraint(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := testPrompt("unique-test")
	if err := sp.Create(prompt); err != nil {
		t.Fatalf("First create failed: %v", err)
	}

	// Try to insert directly with same name but different id
	_, err := sp.db.Exec(
		"INSERT INTO prompts (id, name, content, created_at, updated_at) VALUES (?, ?, ?, datetime('now'), datetime('now'))",
		"different-id", "unique-test", "content",
	)
	if err == nil {
		t.Fatal("Expected unique constraint violation, got nil")
	}
}

// 6.3 Test error handling scenarios

func TestCreateValidationEmptyName(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := promptsdb.Prompt{Name: "", Content: "some content"}
	err := sp.Create(prompt)
	if err == nil {
		t.Fatal("Expected validation error for empty name, got nil")
	}
}

func TestCreateValidationEmptyContent(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	prompt := promptsdb.Prompt{Name: "test", Content: ""}
	err := sp.Create(prompt)
	if err == nil {
		t.Fatal("Expected validation error for empty content, got nil")
	}
}

func TestNewProviderEmptyPath(t *testing.T) {
	_, err := NewSQLiteProvider(SQLiteProviderConfiguration{Path: ""}, "")
	if err == nil {
		t.Fatal("Expected error for empty path, got nil")
	}
}

// 6.4 Test concurrent access scenarios

func TestConcurrentReads(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	if err := sp.Create(testPrompt("concurrent")); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	var wg sync.WaitGroup
	errors := make(chan error, 20)

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := sp.Read("concurrent")
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent read error: %v", err)
	}
}

func TestConcurrentWrites(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			name := fmt.Sprintf("concurrent-%d", idx)
			if err := sp.Create(testPrompt(name)); err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent write error: %v", err)
	}

	prompts, err := sp.List(promptsdb.PromptQuery{All: true})
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(prompts) != 10 {
		t.Errorf("Expected 10 prompts after concurrent writes, got %d", len(prompts))
	}
}

// 6.5 Test transaction rollback behavior

func TestTransactionRollback(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	tx, err := sp.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx failed: %v", err)
	}

	// Insert a valid prompt in the transaction
	if err := CreateInTx(tx, testPrompt("tx-prompt")); err != nil {
		t.Fatalf("CreateInTx failed: %v", err)
	}

	// Rollback
	if err := tx.Rollback(); err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Verify the prompt was not persisted
	_, err = sp.Read("tx-prompt")
	if err == nil {
		t.Fatal("Expected prompt to be rolled back, but it was found")
	}
}

func TestTransactionCommit(t *testing.T) {
	sp, cleanup := setupTestProvider(t)
	defer cleanup()

	tx, err := sp.BeginTx()
	if err != nil {
		t.Fatalf("BeginTx failed: %v", err)
	}

	if err := CreateInTx(tx, testPrompt("tx-commit")); err != nil {
		t.Fatalf("CreateInTx failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	// Verify the prompt was persisted
	got, err := sp.Read("tx-commit")
	if err != nil {
		t.Fatalf("Expected prompt to be persisted after commit: %v", err)
	}
	if got.Name != "tx-commit" {
		t.Errorf("Expected name tx-commit, got %s", got.Name)
	}
}

// Persistence across restarts
func TestPersistenceAcrossReopen(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sqlite_persist_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "persist.db")

	// Create and populate
	p1, err := NewSQLiteProvider(SQLiteProviderConfiguration{Path: dbPath, Timeout: 30}, "")
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}
	sp1 := p1.(*SQLiteProvider)
	if err := sp1.Create(testPrompt("persist-test")); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	sp1.Close()

	// Reopen and verify
	p2, err := NewSQLiteProvider(SQLiteProviderConfiguration{Path: dbPath, Timeout: 30}, "")
	if err != nil {
		t.Fatalf("Failed to reopen provider: %v", err)
	}
	sp2 := p2.(*SQLiteProvider)
	defer sp2.Close()

	got, err := sp2.Read("persist-test")
	if err != nil {
		t.Fatalf("Failed to read after reopen: %v", err)
	}
	if got.Name != "persist-test" {
		t.Errorf("Expected persist-test, got %s", got.Name)
	}
}
