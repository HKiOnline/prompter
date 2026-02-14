package prompts

import (
	"context"
	"testing"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetWithTemplate(t *testing.T) {
	// Create a test prompt with template syntax
	testPrompt := promptsdb.Prompt{
		Name:        "test-template",
		Title:       "Test Template",
		Description: "Test Description",
		Content:     "Today is {{date}}. Hello {{.name}}!",
	}

	db := NewMockDB([]promptsdb.Prompt{testPrompt})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name: "test-template",
		Arguments: map[string]string{
			"name": "World",
		},
	}

	resp, err := handler.HandleGet(context.Background(), nil, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 1)
	assert.Equal(t, mcp.Role("user"), resp.Messages[0].Role)

	if textContent, ok := resp.Messages[0].Content.(*mcp.TextContent); ok {
		// Check that the template was processed
		content := textContent.Text
		assert.Contains(t, content, "Today is ")
		assert.Contains(t, content, "Hello World!")

		// Check that date is in correct format (YYYY-MM-DD)
		assert.Regexp(t, `Today is \d{4}-\d{2}-\d{2}`, content)

		// Debug output
		t.Logf("Template result: %s", content)
	} else {
		t.Error("Expected TextContent type")
	}
}

func TestHandleGetWithTemplateNoArgs(t *testing.T) {
	// Create a test prompt with template syntax but no arguments
	testPrompt := promptsdb.Prompt{
		Name:        "test-template-no-args",
		Title:       "Test Template No Args",
		Description: "Test Description",
		Content:     "Current date: {{date}}",
	}

	db := NewMockDB([]promptsdb.Prompt{testPrompt})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name:      "test-template-no-args",
		Arguments: map[string]string{},
	}

	resp, err := handler.HandleGet(context.Background(), nil, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 1)
	assert.Equal(t, mcp.Role("user"), resp.Messages[0].Role)

	if textContent, ok := resp.Messages[0].Content.(*mcp.TextContent); ok {
		// Check that the template was processed
		content := textContent.Text
		assert.Contains(t, content, "Current date: ")

		// Check that date is in correct format (YYYY-MM-DD)
		assert.Regexp(t, `Current date: \d{4}-\d{2}-\d{2}`, content)
	} else {
		t.Error("Expected TextContent type")
	}
}

func TestHandleGetWithInvalidTemplate(t *testing.T) {
	// Create a test prompt with invalid template syntax
	testPrompt := promptsdb.Prompt{
		Name:        "test-invalid-template",
		Title:       "Test Invalid Template",
		Description: "Test Description",
		Content:     "Hello {{.InvalidSyntax",
	}

	db := NewMockDB([]promptsdb.Prompt{testPrompt})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name: "test-invalid-template",
	}

	resp, err := handler.HandleGet(context.Background(), nil, req)

	// Should not return error, but should return original content
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 1)

	if textContent, ok := resp.Messages[0].Content.(*mcp.TextContent); ok {
		// Should return original content when template parsing fails
		assert.Equal(t, "Hello {{.InvalidSyntax", textContent.Text)
	} else {
		t.Error("Expected TextContent type")
	}
}
