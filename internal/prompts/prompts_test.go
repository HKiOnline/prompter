package prompts

import (
	"context"
	"testing"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

// MockDB is a mock implementation of promptsdb.Provider for testing
type MockDB struct {
	prompts map[string]promptsdb.Prompt
}

func NewMockDB(prompts []promptsdb.Prompt) *MockDB {
	db := &MockDB{
		prompts: make(map[string]promptsdb.Prompt),
	}
	for _, p := range prompts {
		db.prompts[p.Name] = p
	}
	return db
}

func (m *MockDB) Create(prompt promptsdb.Prompt) error {
	if prompt.Name == "error" {
		return assert.AnError
	}
	m.prompts[prompt.Name] = prompt
	return nil
}

func (m *MockDB) Read(name string) (promptsdb.Prompt, error) {
	if name == "error" {
		return promptsdb.Prompt{}, assert.AnError
	}
	if p, ok := m.prompts[name]; ok {
		return p, nil
	}
	return promptsdb.Prompt{}, assert.AnError
}

func (m *MockDB) List(query promptsdb.PromptQuery) ([]promptsdb.Prompt, error) {
	results := make([]promptsdb.Prompt, 0, len(m.prompts))
	for _, p := range m.prompts {
		results = append(results, p)
	}
	return results, nil
}

func (m *MockDB) Update(prompt promptsdb.Prompt) error {
	if prompt.Name == "error" {
		return assert.AnError
	}
	m.prompts[prompt.Name] = prompt
	return nil
}

func (m *MockDB) Delete(name string) error {
	if name == "error" {
		return assert.AnError
	}
	delete(m.prompts, name)
	return nil
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) Setup(config promptsdb.ProviderConfiguration) error {
	return nil
}

func TestNewPromptHandler(t *testing.T) {
	db := NewMockDB([]promptsdb.Prompt{})
	logger := plog.New("/tmp/test.log")

	handler := NewPromptHandler(db, logger)

	assert.NotNil(t, handler)
	assert.Equal(t, db, handler.db)
	assert.Equal(t, logger, handler.logger)
}

func TestHandleList(t *testing.T) {
	testPrompts := []promptsdb.Prompt{
		{
			Name:        "test1",
			Title:       "Test 1",
			Description: "Description 1",
		},
		{
			Name:        "test2",
			Title:       "Test 2",
			Description: "Description 2",
		},
	}
	db := NewMockDB(testPrompts)
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.ListPromptsParams{}
	resp, err := handler.HandleList(context.Background(), nil, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Prompts, 2)
	assert.Equal(t, "test1", resp.Prompts[0].Name)
	assert.Equal(t, "test2", resp.Prompts[1].Name)
}

func TestHandleGet(t *testing.T) {
	testPrompt := promptsdb.Prompt{
		Name:        "test",
		Title:       "Test Prompt",
		Description: "Test Description",
		Content:     "Test Content",
	}
	db := NewMockDB([]promptsdb.Prompt{testPrompt})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name: "test",
	}
	resp, err := handler.HandleGet(context.Background(), nil, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 1)
	assert.Equal(t, mcp.Role("user"), resp.Messages[0].Role)
	if textContent, ok := resp.Messages[0].Content.(*mcp.TextContent); ok {
		assert.Equal(t, "Test Content", textContent.Text)
	} else {
		t.Error("Expected TextContent type")
	}
}

func TestHandleGetError(t *testing.T) {
	db := NewMockDB([]promptsdb.Prompt{})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name: "error",
	}
	resp, err := handler.HandleGet(context.Background(), nil, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestHandleGetEmptyName(t *testing.T) {
	db := NewMockDB([]promptsdb.Prompt{})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name: "",
	}
	resp, err := handler.HandleGet(context.Background(), nil, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
