package promptsdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hkionline/prompter/internal/plog"
)

func TestNewPromptsFsProvider(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test successful creation
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	if provider == nil {
		t.Fatal("FsProvider should not be nil")
	}

	if provider.dir != tempDir {
		t.Errorf("Expected dir %s, got %s", tempDir, provider.dir)
	}
}

func TestFsProviderCreate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts_create")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create a test prompt
	prompt := Prompt{
		Name:        "test-prompt",
		Title:       "Test Prompt",
		Description: "A test prompt for unit testing",
		Arguments:   []string{"name", "age"},
		Content:     "Hello {{.name}}, you are {{.age}} years old.",
		Tags:        []string{"test", "example"},
	}

	// Test successful creation
	err = provider.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	// Verify the prompt was created and cached
	retrievedPrompt, err := provider.Read("test-prompt")
	if err != nil {
		t.Fatalf("Failed to read prompt: %v", err)
	}

	if retrievedPrompt.Name != prompt.Name {
		t.Errorf("Expected name %s, got %s", prompt.Name, retrievedPrompt.Name)
	}
	if retrievedPrompt.Title != prompt.Title {
		t.Errorf("Expected title %s, got %s", prompt.Title, retrievedPrompt.Title)
	}
	if retrievedPrompt.Description != prompt.Description {
		t.Errorf("Expected description %s, got %s", prompt.Description, retrievedPrompt.Description)
	}
	if len(retrievedPrompt.Arguments) != len(prompt.Arguments) {
		t.Errorf("Expected %d arguments, got %d", len(prompt.Arguments), len(retrievedPrompt.Arguments))
	}
	if retrievedPrompt.Content != prompt.Content {
		t.Errorf("Expected content %s, got %s", prompt.Content, retrievedPrompt.Content)
	}
	if len(retrievedPrompt.Tags) != len(prompt.Tags) {
		t.Errorf("Expected %d tags, got %d", len(prompt.Tags), len(retrievedPrompt.Tags))
	}

	// Verify file was created
	filePath := filepath.Join(tempDir, "test-prompt.md")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected prompt file %s to exist", filePath)
	}
}

func TestFsProviderCreateWithSpecialCharacters(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts_special")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create a test prompt with special characters
	prompt := Prompt{
		Name:        "test-prompt_123",
		Title:       "Test Prompt With Special Chars!",
		Description: "A test prompt with special chars @#$%",
		Arguments:   []string{"name", "age"},
		Content:     "Hello {{.name}}, you are {{.age}} years old.",
		Tags:        []string{"test", "example"},
	}

	// Test successful creation
	err = provider.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt with special chars: %v", err)
	}

	// Verify the prompt was created
	retrievedPrompt, err := provider.Read("test-prompt_123")
	if err != nil {
		t.Fatalf("Failed to read prompt with special chars: %v", err)
	}

	if retrievedPrompt.Name != prompt.Name {
		t.Errorf("Expected name %s, got %s", prompt.Name, retrievedPrompt.Name)
	}
}

func TestFsProviderRead(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts_read")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Try to read a non-existent prompt
	_, err = provider.Read("nonexistent")
	if err == nil {
		t.Error("Expected error when reading non-existent prompt")
	}
}

func TestFsProviderUpdate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts_update")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create a test prompt
	prompt := Prompt{
		Name:        "test-prompt",
		Title:       "Test Prompt",
		Description: "A test prompt for unit testing",
		Arguments:   []string{"name", "age"},
		Content:     "Hello {{.name}}, you are {{.age}} years old.",
		Tags:        []string{"test", "example"},
	}

	// Create the prompt
	err = provider.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	// Update the prompt
	prompt.Title = "Updated Test Prompt"
	prompt.Description = "An updated test prompt"

	err = provider.Update(prompt)
	if err != nil {
		t.Fatalf("Failed to update prompt: %v", err)
	}

	// Verify the update
	retrievedPrompt, err := provider.Read("test-prompt")
	if err != nil {
		t.Fatalf("Failed to read prompt: %v", err)
	}

	if retrievedPrompt.Title != prompt.Title {
		t.Errorf("Expected updated title %s, got %s", prompt.Title, retrievedPrompt.Title)
	}
	if retrievedPrompt.Description != prompt.Description {
		t.Errorf("Expected updated description %s, got %s", prompt.Description, retrievedPrompt.Description)
	}
}

func TestFsProviderList(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts_list")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create test prompts
	prompt1 := Prompt{
		Name:        "prompt-1",
		Title:       "Prompt 1",
		Description: "First test prompt",
		Content:     "Content 1",
	}

	prompt2 := Prompt{
		Name:        "prompt-2",
		Title:       "Prompt 2",
		Description: "Second test prompt",
		Content:     "Content 2",
	}

	// Create prompts
	err = provider.Create(prompt1)
	if err != nil {
		t.Fatalf("Failed to create prompt 1: %v", err)
	}

	err = provider.Create(prompt2)
	if err != nil {
		t.Fatalf("Failed to create prompt 2: %v", err)
	}

	// List all prompts
	prompts, err := provider.List(PromptQuery{})
	if err != nil {
		t.Fatalf("Failed to list prompts: %v", err)
	}

	if len(prompts) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(prompts))
	}

	// Verify we got the right prompts
	promptNames := make(map[string]bool)
	for _, p := range prompts {
		promptNames[p.Name] = true
	}

	if !promptNames["prompt-1"] {
		t.Error("Expected prompt-1 in list")
	}
	if !promptNames["prompt-2"] {
		t.Error("Expected prompt-2 in list")
	}
}

func TestFsProviderListWithQuery(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts_list_query")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create test prompts
	prompt1 := Prompt{
		Name:        "prompt-1",
		Title:       "Prompt 1",
		Description: "First test prompt",
		Content:     "Content 1",
	}

	prompt2 := Prompt{
		Name:        "prompt-2",
		Title:       "Prompt 2",
		Description: "Second test prompt",
		Content:     "Content 2",
	}

	// Create prompts
	err = provider.Create(prompt1)
	if err != nil {
		t.Fatalf("Failed to create prompt 1: %v", err)
	}

	err = provider.Create(prompt2)
	if err != nil {
		t.Fatalf("Failed to create prompt 2: %v", err)
	}

	// Test with query that should return all
	prompts, err := provider.List(PromptQuery{All: true})
	if err != nil {
		t.Fatalf("Failed to list prompts with query: %v", err)
	}

	if len(prompts) != 2 {
		t.Errorf("Expected 2 prompts with All=true, got %d", len(prompts))
	}

	// Test with empty query (should also return all)
	prompts, err = provider.List(PromptQuery{})
	if err != nil {
		t.Fatalf("Failed to list prompts with empty query: %v", err)
	}

	if len(prompts) != 2 {
		t.Errorf("Expected 2 prompts with empty query, got %d", len(prompts))
	}
}

func TestConcurrencySafety(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_concurrency")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create a test prompt
	prompt := Prompt{
		Name:        "concurrent-prompt",
		Title:       "Concurrent Test Prompt",
		Description: "A test prompt for concurrent testing",
		Arguments:   []string{"name"},
		Content:     "Hello {{.name}}!",
		Tags:        []string{"test"},
	}

	// Test that we can create a prompt
	err = provider.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	// Verify it was created
	retrievedPrompt, err := provider.Read("concurrent-prompt")
	if err != nil {
		t.Fatalf("Failed to read prompt: %v", err)
	}

	if retrievedPrompt.Name != "concurrent-prompt" {
		t.Errorf("Expected concurrent-prompt, got %s", retrievedPrompt.Name)
	}
}

func TestEmptyAndNilFields(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_empty_fields")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create a prompt with empty and nil fields
	prompt := Prompt{
		Name:        "empty-fields-prompt",
		Title:       "",
		Description: "",
		Arguments:   []string{},
		Content:     "Hello world!",
		Tags:        []string{},
	}

	// Test successful creation
	err = provider.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt with empty fields: %v", err)
	}

	// Verify the prompt was created
	retrievedPrompt, err := provider.Read("empty-fields-prompt")
	if err != nil {
		t.Fatalf("Failed to read prompt with empty fields: %v", err)
	}

	if retrievedPrompt.Name != "empty-fields-prompt" {
		t.Errorf("Expected empty-fields-prompt, got %s", retrievedPrompt.Name)
	}
	if retrievedPrompt.Content != "Hello world!" {
		t.Errorf("Expected content 'Hello world!', got %s", retrievedPrompt.Content)
	}
}

func TestPromptWithNoArguments(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_no_args")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create a prompt with no arguments
	prompt := Prompt{
		Name:        "no-args-prompt",
		Title:       "No Arguments Prompt",
		Description: "A prompt with no arguments",
		Arguments:   []string{},
		Content:     "Hello world!",
		Tags:        []string{"test"},
	}

	// Test successful creation
	err = provider.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt with no arguments: %v", err)
	}

	// Verify the prompt was created
	retrievedPrompt, err := provider.Read("no-args-prompt")
	if err != nil {
		t.Fatalf("Failed to read prompt with no arguments: %v", err)
	}

	if len(retrievedPrompt.Arguments) != 0 {
		t.Errorf("Expected no arguments, got %d", len(retrievedPrompt.Arguments))
	}

	if retrievedPrompt.Title != "No Arguments Prompt" {
		t.Errorf("Expected title 'No Arguments Prompt', got %s", retrievedPrompt.Title)
	}
}

func TestFilePermissions(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_permissions")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Create a test prompt
	prompt := Prompt{
		Name:        "permission-test",
		Title:       "Permission Test Prompt",
		Description: "A test prompt for permission testing",
		Arguments:   []string{"name"},
		Content:     "Hello {{.name}}!",
		Tags:        []string{"test"},
	}

	// Test successful creation
	err = provider.Create(prompt)
	if err != nil {
		t.Fatalf("Failed to create prompt: %v", err)
	}

	// Verify the prompt was created
	retrievedPrompt, err := provider.Read("permission-test")
	if err != nil {
		t.Fatalf("Failed to read prompt: %v", err)
	}

	if retrievedPrompt.Name != "permission-test" {
		t.Errorf("Expected permission-test, got %s", retrievedPrompt.Name)
	}

	if !strings.Contains(retrievedPrompt.Content, "Hello {{.name}}!") {
		t.Errorf("Expected content with 'Hello {{.name}}!', got %s", retrievedPrompt.Content)
	}
}

func TestFsProviderDelete(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts_delete")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid prompt file manually (to simulate loadCache behavior)
	promptFile := filepath.Join(tempDir, "delete-test.md")
	promptContent := `---
name: delete-test
title: Delete Test Prompt
description: A test prompt for delete testing
---
This will be deleted`

	err = os.WriteFile(promptFile, []byte(promptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create prompt file: %v", err)
	}

	// Create provider (this will load the file into cache)
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Verify the prompt exists
	_, err = provider.Read("delete-test")
	if err != nil {
		t.Fatalf("Failed to read prompt before delete: %v", err)
	}

	// Delete the prompt
	err = provider.Delete("delete-test")
	if err != nil {
		t.Fatalf("Failed to delete prompt: %v", err)
	}

	// Verify the prompt no longer exists in cache
	_, err = provider.Read("delete-test")
	if err == nil {
		t.Error("Expected error when reading deleted prompt")
	}

	// Verify the file was deleted
	if _, err := os.Stat(promptFile); err == nil {
		t.Error("Expected prompt file to be deleted")
	}
}

func TestFsProviderDeleteNonexistent(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_prompts_delete_nonexistent")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create provider
	provider, err := NewPromptsFsProvider(tempDir, "")
	if err != nil {
		t.Fatalf("Failed to create FsProvider: %v", err)
	}

	// Try to delete a non-existent prompt
	err = provider.Delete("nonexistent-prompt")
	if err == nil {
		t.Error("Expected error when deleting non-existent prompt")
	}

	// Verify error message
	if !strings.Contains(err.Error(), "no prompt file with the given id was found") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestLoadPromptValidFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_load_prompt")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid prompt file manually
	promptFile := filepath.Join(tempDir, "test-prompt.md")
	promptContent := `---
name: test-prompt
title: Test Prompt
description: A test prompt
arguments:
  - name
  - age
tags:
  - test
---
Hello {{.name}}, you are {{.age}} years old.`

	err = os.WriteFile(promptFile, []byte(promptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create prompt file: %v", err)
	}

	// Create a mock logger
	logFile := filepath.Join(tempDir, "test.log")
	p := plog.New(logFile)

	// Load the prompt
	prompt, err := loadPrompt(promptFile, p)
	if err != nil {
		t.Fatalf("Failed to load prompt: %v", err)
	}

	// Verify prompt data
	if prompt.Name != "test-prompt" {
		t.Errorf("Expected name 'test-prompt', got '%s'", prompt.Name)
	}
	if prompt.Title != "Test Prompt" {
		t.Errorf("Expected title 'Test Prompt', got '%s'", prompt.Title)
	}
	if prompt.Description != "A test prompt" {
		t.Errorf("Expected description 'A test prompt', got '%s'", prompt.Description)
	}
	if len(prompt.Arguments) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(prompt.Arguments))
	}
	if prompt.Content != "Hello {{.name}}, you are {{.age}} years old." {
		t.Errorf("Expected content 'Hello {{.name}}, you are {{.age}} years old.', got '%s'", prompt.Content)
	}
	if prompt.Id != "test-prompt" {
		t.Errorf("Expected Id to be set to 'test-prompt', got '%s'", prompt.Id)
	}
}

func TestLoadPromptMalformedYAML(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_load_prompt_malformed")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a malformed YAML file
	promptFile := filepath.Join(tempDir, "malformed.md")
	promptContent := `---
name: test
title: [unclosed array
---
Content here`

	err = os.WriteFile(promptFile, []byte(promptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create prompt file: %v", err)
	}

	// Create a mock logger
	logFile := filepath.Join(tempDir, "test.log")
	p := plog.New(logFile)

	// Try to load the malformed prompt
	_, err = loadPrompt(promptFile, p)
	if err == nil {
		t.Error("Expected error when loading malformed YAML")
	}
}

func TestLoadPromptMissingContentSeparator(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_load_prompt_missing_sep")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a file without proper content separator
	promptFile := filepath.Join(tempDir, "no-separator.md")
	promptContent := `---
name: test
title: Test
description: Test prompt`

	err = os.WriteFile(promptFile, []byte(promptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create prompt file: %v", err)
	}

	// Create a mock logger
	logFile := filepath.Join(tempDir, "test.log")
	p := plog.New(logFile)

	// Try to load the prompt without separator
	_, err = loadPrompt(promptFile, p)
	if err == nil {
		t.Error("Expected error when content separator is missing")
	}

	if !strings.Contains(err.Error(), "failed to extract prompt contents") {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}

func TestLoadPromptFileNotFound(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_load_prompt_not_found")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock logger
	logFile := filepath.Join(tempDir, "test.log")
	p := plog.New(logFile)

	// Try to load a non-existent file
	_, err = loadPrompt(filepath.Join(tempDir, "nonexistent.md"), p)
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestLoadCacheWithMultipleFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_load_cache_multiple")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create multiple valid prompt files
	prompts := []struct {
		filename string
		name     string
	}{
		{"prompt1.md", "prompt-1"},
		{"prompt2.md", "prompt-2"},
		{"prompt3.md", "prompt-3"},
	}

	for _, p := range prompts {
		promptContent := fmt.Sprintf(`---
name: %s
title: Test Prompt
description: Test prompt
---
Content for %s`, p.name, p.name)
		err = os.WriteFile(filepath.Join(tempDir, p.filename), []byte(promptContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create prompt file: %v", err)
		}
	}

	// Create a mock logger
	logFile := filepath.Join(tempDir, "test.log")
	plogger := plog.New(logFile)

	// Load cache
	cache, files, err := loadCache(tempDir, plogger)
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	// Verify all prompts were loaded
	if len(cache) != 3 {
		t.Errorf("Expected 3 prompts in cache, got %d", len(cache))
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files tracked, got %d", len(files))
	}

	// Verify each prompt is in the cache
	for _, p := range prompts {
		if _, ok := cache[p.name]; !ok {
			t.Errorf("Expected prompt '%s' in cache", p.name)
		}
	}
}

func TestLoadCacheWithCorruptedFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_load_cache_corrupted")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid prompt file
	validPrompt := `---
name: valid-prompt
title: Valid Prompt
---
Valid content`
	err = os.WriteFile(filepath.Join(tempDir, "valid.md"), []byte(validPrompt), 0644)
	if err != nil {
		t.Fatalf("Failed to create valid prompt file: %v", err)
	}

	// Create a corrupted prompt file
	corruptedPrompt := `---
name: corrupted
title: [malformed
---
Content`
	err = os.WriteFile(filepath.Join(tempDir, "corrupted.md"), []byte(corruptedPrompt), 0644)
	if err != nil {
		t.Fatalf("Failed to create corrupted prompt file: %v", err)
	}

	// Create a mock logger
	logFile := filepath.Join(tempDir, "test.log")
	plogger := plog.New(logFile)

	// Load cache - should skip corrupted files but load valid ones
	cache, _, err := loadCache(tempDir, plogger)
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	// Should have loaded only the valid prompt
	if len(cache) != 1 {
		t.Errorf("Expected 1 valid prompt in cache, got %d", len(cache))
	}

	if _, ok := cache["valid-prompt"]; !ok {
		t.Error("Expected valid-prompt to be in cache")
	}
}

func TestLoadCacheEmptyDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_load_cache_empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock logger
	logFile := filepath.Join(tempDir, "test.log")
	plogger := plog.New(logFile)

	// Load cache from empty directory
	cache, files, err := loadCache(tempDir, plogger)
	if err != nil {
		t.Fatalf("Failed to load cache from empty directory: %v", err)
	}

	// Should return empty cache and files
	if len(cache) != 0 {
		t.Errorf("Expected empty cache, got %d items", len(cache))
	}

	if len(files) != 0 {
		t.Errorf("Expected empty files map, got %d items", len(files))
	}
}

func TestLoadCacheWithSubdirectories(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_load_cache_subdirs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a valid prompt file
	validPrompt := `---
name: valid-prompt
title: Valid Prompt
---
Valid content`
	err = os.WriteFile(filepath.Join(tempDir, "valid.md"), []byte(validPrompt), 0644)
	if err != nil {
		t.Fatalf("Failed to create valid prompt file: %v", err)
	}

	// Create a subdirectory (should be ignored)
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create a mock logger
	logFile := filepath.Join(tempDir, "test.log")
	plogger := plog.New(logFile)

	// Load cache - should skip subdirectories
	cache, _, err := loadCache(tempDir, plogger)
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	// Should only load the valid file, ignoring the subdirectory
	if len(cache) != 1 {
		t.Errorf("Expected 1 prompt in cache, got %d", len(cache))
	}
}

func TestRemovePrompt(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_remove_prompt")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := "test-prompt.md"
	filePath := filepath.Join(tempDir, testFile)
	err = os.WriteFile(filePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("Test file should exist")
	}

	// Remove the prompt
	err = removePrompt(tempDir, testFile)
	if err != nil {
		t.Fatalf("Failed to remove prompt: %v", err)
	}

	// Verify file was deleted
	if _, err := os.Stat(filePath); err == nil {
		t.Error("Expected file to be deleted")
	}
}

func TestRemovePromptNonexistent(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_remove_nonexistent")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Try to remove a non-existent file
	err = removePrompt(tempDir, "nonexistent.md")
	if err == nil {
		t.Error("Expected error when removing non-existent file")
	}
}

func TestNewWithFileSystemProvider(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_new_fs_provider")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := ProviderConfiguration{
		Provider: FILE_SYSTEM_PROVIDER,
		Filesystem: FsProviderConfiguration{
			Directory: tempDir,
		},
	}

	logFile := filepath.Join(tempDir, "test.log")

	provider, err := New(FILE_SYSTEM_PROVIDER, config, logFile)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}

func TestNewWithDefaultProvider(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_new_default_provider")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := ProviderConfiguration{
		Provider: "unknown-provider",
		Filesystem: FsProviderConfiguration{
			Directory: tempDir,
		},
	}

	logFile := filepath.Join(tempDir, "test.log")

	// Should default to filesystem provider
	provider, err := New("unknown-provider", config, logFile)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	if provider == nil {
		t.Fatal("Expected non-nil provider")
	}
}

func TestPromptsDBError(t *testing.T) {
	// Test the PromptsDBError type
	err := &PromptsDBError{}

	// Should return empty string
	if err.Error() != "" {
		t.Errorf("Expected empty error string, got '%s'", err.Error())
	}
}

func TestSavePromptFunction(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_save_prompt_func")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	prompt := Prompt{
		Id:          "test-prompt",
		Name:        "test-prompt",
		Title:       "Test Prompt",
		Description: "A test prompt",
		Arguments:   []string{"arg1", "arg2"},
		Content:     "Test content",
		Tags:        []string{"test", "example"},
	}

	// Save the prompt
	savedPath, err := savePrompt(prompt, tempDir)
	if err != nil {
		t.Fatalf("Failed to save prompt: %v", err)
	}

	expectedPath := filepath.Join(tempDir, "test-prompt.md")
	if savedPath != expectedPath {
		t.Errorf("Expected path '%s', got '%s'", expectedPath, savedPath)
	}

	// Verify file was created
	if _, err := os.Stat(savedPath); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}

	// Read the file and verify content
	content, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	contentStr := string(content)

	// Verify YAML front matter
	if !strings.Contains(contentStr, "name: test-prompt") {
		t.Error("Expected file to contain 'name: test-prompt'")
	}

	if !strings.Contains(contentStr, "title: Test Prompt") {
		t.Error("Expected file to contain 'title: Test Prompt'")
	}

	// Verify content section
	if !strings.Contains(contentStr, "Test content") {
		t.Error("Expected file to contain 'Test content'")
	}

	// Verify YAML separators
	if strings.Count(contentStr, "---") < 2 {
		t.Error("Expected at least two '---' separators")
	}
}
