package rpc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
)

func TestProcessValidInitialize(t *testing.T) {
	// Create a mock logger
	tempDir, err := os.MkdirTemp("", "test_logs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logfile := filepath.Join(tempDir, "test.log")
	p := plog.New(logfile)

	// Create a valid initialize request
	reqStr := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "initialize"
	}`

	// Create a mock db provider
	mockDB := &MockProvider{}

	// Process the request
	respStr := Process(reqStr, mockDB, p)

	// Verify response is not empty and contains expected fields
	if respStr == "" {
		t.Error("Expected non-empty response for initialize request")
	}

	// Parse the response to validate structure
	var resp Message
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Version != "2.0" {
		t.Errorf("Expected version '2.0', got '%s'", resp.Version)
	}

	if resp.Id != 1 {
		t.Errorf("Expected id 1, got %d", resp.Id)
	}

	if resp.Result.ProtocolVersion != "" {
		t.Errorf("Expected empty protocol version, got '%s'", resp.Result.ProtocolVersion)
	}

	if resp.Result.ServerInfo.Name != "prompter" {
		t.Errorf("Expected server name 'prompter', got '%s'", resp.Result.ServerInfo.Name)
	}
}

func TestProcessValidPing(t *testing.T) {
	// Create a mock logger
	tempDir, err := os.MkdirTemp("", "test_logs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logfile := filepath.Join(tempDir, "test.log")
	p := plog.New(logfile)

	// Create a valid ping request
	reqStr := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "ping"
	}`

	// Create a mock db provider
	mockDB := &MockProvider{}

	// Process the request
	respStr := Process(reqStr, mockDB, p)

	// Verify response is not empty and contains expected fields
	if respStr == "" {
		t.Error("Expected non-empty response for ping request")
	}

	// Parse the response to validate structure
	var resp Message
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Version != "2.0" {
		t.Errorf("Expected version '2.0', got '%s'", resp.Version)
	}

	if resp.Id != 1 {
		t.Errorf("Expected id 1, got %d", resp.Id)
	}
}

func TestProcessValidPromptsList(t *testing.T) {
	// Create a mock logger
	tempDir, err := os.MkdirTemp("", "test_logs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logfile := filepath.Join(tempDir, "test.log")
	p := plog.New(logfile)

	// Create a valid prompts/list request
	reqStr := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "prompts/list"
	}`

	// Create a mock db provider that returns test prompts
	mockDB := &MockProvider{
		prompts: []promptsdb.Prompt{
			{
				Name:        "test-prompt",
				Title:       "Test Prompt",
				Description: "A test prompt for unit testing",
				Content:     "Hello {{.name}}!",
			},
		},
	}

	// Process the request
	respStr := Process(reqStr, mockDB, p)

	// Verify response is not empty and contains expected fields
	if respStr == "" {
		t.Error("Expected non-empty response for prompts/list request")
	}

	// Parse the response to validate structure
	var resp Message
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Version != "2.0" {
		t.Errorf("Expected version '2.0', got '%s'", resp.Version)
	}

	if resp.Id != 1 {
		t.Errorf("Expected id 1, got %d", resp.Id)
	}

	if len(resp.Result.Prompts) != 1 {
		t.Errorf("Expected 1 prompt in result, got %d", len(resp.Result.Prompts))
	}

	if resp.Result.Prompts[0].Name != "test-prompt" {
		t.Errorf("Expected prompt name 'test-prompt', got '%s'", resp.Result.Prompts[0].Name)
	}
}

func TestProcessValidPromptsGet(t *testing.T) {
	// Create a mock logger
	tempDir, err := os.MkdirTemp("", "test_logs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logfile := filepath.Join(tempDir, "test.log")
	p := plog.New(logfile)

	// Create a valid prompts/get request
	reqStr := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "prompts/get",
		"params": {
			"name": "test-prompt"
		}
	}`

	// Create a mock db provider that returns test prompts
	mockDB := &MockProvider{
		prompts: []promptsdb.Prompt{
			{
				Name:        "test-prompt",
				Title:       "Test Prompt",
				Description: "A test prompt for unit testing",
				Content:     "Hello {{.name}}!",
			},
		},
	}

	// Process the request
	respStr := Process(reqStr, mockDB, p)

	// Verify response is not empty and contains expected fields
	if respStr == "" {
		t.Error("Expected non-empty response for prompts/get request")
	}

	// Parse the response to validate structure
	var resp Message
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Version != "2.0" {
		t.Errorf("Expected version '2.0', got '%s'", resp.Version)
	}

	if resp.Id != 1 {
		t.Errorf("Expected id 1, got %d", resp.Id)
	}

	if len(resp.Result.Messages) != 1 {
		t.Errorf("Expected 1 message in result, got %d", len(resp.Result.Messages))
	}

	if resp.Result.Messages[0].Content.Text != "Hello {{.name}}!" {
		t.Errorf("Expected content 'Hello {{.name}}!', got '%s'", resp.Result.Messages[0].Content.Text)
	}
}

func TestProcessValidToolsList(t *testing.T) {
	// Create a mock logger
	tempDir, err := os.MkdirTemp("", "test_logs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logfile := filepath.Join(tempDir, "test.log")
	p := plog.New(logfile)

	// Create a valid tools/list request
	reqStr := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "tools/list"
	}`

	// Process the request
	respStr := Process(reqStr, &MockProvider{}, p)

	// Verify response is not empty and contains expected fields
	if respStr == "" {
		t.Error("Expected non-empty response for tools/list request")
	}

	// Parse the response to validate structure
	var resp Message
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Version != "2.0" {
		t.Errorf("Expected version '2.0', got '%s'", resp.Version)
	}

	if resp.Id != 1 {
		t.Errorf("Expected id 1, got %d", resp.Id)
	}

	if len(resp.Result.Tools) != 1 {
		t.Errorf("Expected 1 tool in result, got %d", len(resp.Result.Tools))
	}

	if resp.Result.Tools[0].Name != "create_prompt" {
		t.Errorf("Expected tool name 'create_prompt', got '%s'", resp.Result.Tools[0].Name)
	}
}

func TestProcessInvalidJSON(t *testing.T) {
	// Create a mock logger
	tempDir, err := os.MkdirTemp("", "test_logs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logfile := filepath.Join(tempDir, "test.log")
	p := plog.New(logfile)

	// Create invalid JSON
	reqStr := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "initialize"
		// Missing closing brace
	}`

	// Process the request
	respStr := Process(reqStr, &MockProvider{}, p)

	// Should return error response
	if respStr == "" {
		t.Error("Expected error response for invalid JSON")
	}

	// Parse the response to validate structure
	var resp Message
	err = json.Unmarshal([]byte(respStr), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp.Error.Code != 101 {
		t.Errorf("Expected error code 101, got %d", resp.Error.Code)
	}
}

func TestProcessUnknownMethod(t *testing.T) {
	// Create a mock logger
	tempDir, err := os.MkdirTemp("", "test_logs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logfile := filepath.Join(tempDir, "test.log")
	p := plog.New(logfile)

	// Create request with unknown method
	reqStr := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "unknown/method"
	}`

	// Process the request
	respStr := Process(reqStr, &MockProvider{}, p)

	// Should return empty response (noOp=true)
	if respStr != "" {
		t.Error("Expected empty response for unknown method")
	}
}

func TestProcessNotificationsInitialized(t *testing.T) {
	// Create a mock logger
	tempDir, err := os.MkdirTemp("", "test_logs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	logfile := filepath.Join(tempDir, "test.log")
	p := plog.New(logfile)

	// Create notifications/initialized request
	reqStr := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "notifications/initialized"
	}`

	// Process the request
	respStr := Process(reqStr, &MockProvider{}, p)

	// Should return empty response (noOp=true)
	if respStr != "" {
		t.Error("Expected empty response for notifications/initialized")
	}
}

// Mock provider implementation for testing
type MockProvider struct {
	prompts []promptsdb.Prompt
}

func (m *MockProvider) Create(prompt promptsdb.Prompt) error {
	return nil
}

func (m *MockProvider) Read(promptId string) (promptsdb.Prompt, error) {
	for _, p := range m.prompts {
		if p.Name == promptId {
			return p, nil
		}
	}
	return promptsdb.Prompt{}, &promptsdb.PromptsDBError{}
}

func (m *MockProvider) Update(prompt promptsdb.Prompt) error {
	return nil
}

func (m *MockProvider) Delete(promptId string) error {
	return nil
}

func (m *MockProvider) List(query promptsdb.PromptQuery) ([]promptsdb.Prompt, error) {
	return m.prompts, nil
}
