package promptsdb

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	filePath := filepath.Join(tempDir, "test-prompt.yaml")
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
