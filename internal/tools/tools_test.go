package tools

import (
	"context"
	"sync"
	"testing"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

var testServer *mcp.Server
var testSession *mcp.ServerSession
var setupMutex sync.Mutex

func setupTestServer() {
	setupMutex.Lock()
	defer setupMutex.Unlock()

	if testServer == nil {
		testServer = mcp.NewServer("test", "1.0.0", &mcp.ServerOptions{})
	}
	if testSession == nil {
		var err error
		testSession, err = testServer.Connect(context.Background(), mcp.NewStdioTransport())
		if err != nil {
			panic("Failed to create test server session:" + err.Error())
		}
	}
}

// MockDB is a mock implementation of promptsdb.Provider for testing
type MockDB struct{}

func (m *MockDB) Create(prompt promptsdb.Prompt) error {
	if prompt.Name == "error" {
		return assert.AnError
	}
	return nil
}

func (m *MockDB) Read(name string) (promptsdb.Prompt, error) {
	return promptsdb.Prompt{}, nil
}

func (m *MockDB) List(query promptsdb.PromptQuery) ([]promptsdb.Prompt, error) {
	return []promptsdb.Prompt{}, nil
}

func (m *MockDB) Update(prompt promptsdb.Prompt) error {
	return nil
}

func (m *MockDB) Delete(name string) error {
	return nil
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) Setup(config promptsdb.ProviderConfiguration) error {
	return nil
}

func TestNewToolHandler(t *testing.T) {
	db := &MockDB{}
	logger := plog.New("/tmp/test.log")

	handler := NewToolHandler(db, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, db, handler.db)
	assert.Equal(t, logger, handler.logger)
}

func TestHandleList(t *testing.T) {
	setupTestServer()

	db := &MockDB{}
	logger := plog.New("/tmp/test.log")
	handler := NewToolHandler(db, logger)

	tool := handler.CreatePromptTool()

	assert.NotNil(t, tool)
	assert.Equal(t, CREATE_PROMPT, tool.Name)
	assert.Equal(t, "Create prompt", tool.Title)
	assert.Equal(t, "Create and save a new prompt", tool.Description)
	assert.Equal(t, "object", tool.InputSchema.Type)
}

func TestHandleCallSaveNewPrompt(t *testing.T) {
	setupTestServer()

	db := &MockDB{}
	logger := plog.New("/tmp/test.log")
	handler := NewToolHandler(db, logger)

	req := &mcp.CallToolParamsFor[map[string]any]{
		Name: CREATE_PROMPT,
		Arguments: map[string]any{
			"name":        "test_prompt",
			"title":       "Test Prompt",
			"description": "Test Description",
			"content":     "Test Content",
		},
	}
	resp, err := handler.HandleCall(context.Background(), testSession, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Content, 1)
	if textContent, ok := resp.Content[0].(*mcp.TextContent); ok {
		assert.Contains(t, textContent.Text, "test_prompt")
		assert.Contains(t, textContent.Text, "Test Prompt")
	} else {
		t.Error("Expected TextContent type")
	}
	assert.False(t, resp.IsError)
}

func TestHandleCallSaveNewPromptError(t *testing.T) {
	setupTestServer()

	db := &MockDB{}
	logger := plog.New("/tmp/test.log")
	handler := NewToolHandler(db, logger)

	req := &mcp.CallToolParamsFor[map[string]any]{
		Name: CREATE_PROMPT,
		Arguments: map[string]any{
			"name": "error",
		},
	}
	resp, err := handler.HandleCall(context.Background(), testSession, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestHandleCallSaveNewPromptMissingName(t *testing.T) {
	setupTestServer()

	db := &MockDB{}
	logger := plog.New("/tmp/test.log")
	handler := NewToolHandler(db, logger)

	req := &mcp.CallToolParamsFor[map[string]any]{
		Name: CREATE_PROMPT,
		Arguments: map[string]any{
			"name": "",
		},
	}
	resp, err := handler.HandleCall(context.Background(), testSession, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestHandleCallSaveNewPromptMissingArguments(t *testing.T) {
	setupTestServer()

	db := &MockDB{}
	logger := plog.New("/tmp/test.log")
	handler := NewToolHandler(db, logger)

	req := &mcp.CallToolParamsFor[map[string]any]{
		Name: CREATE_PROMPT,
	}
	resp, err := handler.HandleCall(context.Background(), testSession, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestHandleCallUnsupportedTool(t *testing.T) {
	setupTestServer()

	db := &MockDB{}
	logger := plog.New("/tmp/test.log")
	handler := NewToolHandler(db, logger)

	req := &mcp.CallToolParamsFor[map[string]any]{
		Name: "unsupported_tool",
	}
	resp, err := handler.HandleCall(context.Background(), testSession, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
