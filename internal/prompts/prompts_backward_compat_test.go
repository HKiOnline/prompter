package prompts

import (
	"context"
	"testing"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestBackwardCompatibility(t *testing.T) {
	// Test that existing prompts without templates still work
	testPrompt := promptsdb.Prompt{
		Name:        "legacy-prompt",
		Title:       "Legacy Prompt",
		Description: "Prompt without templates",
		Content:     "This is a simple prompt without any template syntax.",
	}

	db := NewMockDB([]promptsdb.Prompt{testPrompt})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name: "legacy-prompt",
	}

	resp, err := handler.HandleGet(context.Background(), nil, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 1)
	assert.Equal(t, mcp.Role("user"), resp.Messages[0].Role)

	if textContent, ok := resp.Messages[0].Content.(*mcp.TextContent); ok {
		// Should return the original content unchanged
		assert.Equal(t, "This is a simple prompt without any template syntax.", textContent.Text)
	} else {
		t.Error("Expected TextContent type")
	}
}

func TestBackwardCompatibilityWithArguments(t *testing.T) {
	// Test that existing prompts with arguments but no template functions still work
	testPrompt := promptsdb.Prompt{
		Name:        "legacy-with-args",
		Title:       "Legacy Prompt with Args",
		Description: "Prompt with arguments but no template functions",
		Content:     "Hello {{.name}}, welcome to our system!",
	}

	db := NewMockDB([]promptsdb.Prompt{testPrompt})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name: "legacy-with-args",
		Arguments: map[string]string{
			"name": "Bob",
		},
	}

	resp, err := handler.HandleGet(context.Background(), nil, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 1)
	assert.Equal(t, mcp.Role("user"), resp.Messages[0].Role)

	if textContent, ok := resp.Messages[0].Content.(*mcp.TextContent); ok {
		// Should process the template with arguments
		assert.Equal(t, "Hello Bob, welcome to our system!", textContent.Text)
	} else {
		t.Error("Expected TextContent type")
	}
}

func TestInvalidTemplateBackwardCompatibility(t *testing.T) {
	// Test that prompts with invalid template syntax return original content
	testPrompt := promptsdb.Prompt{
		Name:        "invalid-template",
		Title:       "Invalid Template",
		Description: "Prompt with invalid template syntax",
		Content:     "This has {{invalid syntax",
	}

	db := NewMockDB([]promptsdb.Prompt{testPrompt})
	logger := plog.New("/tmp/test.log")
	handler := NewPromptHandler(db, logger)

	req := &mcp.GetPromptParams{
		Name: "invalid-template",
	}

	resp, err := handler.HandleGet(context.Background(), nil, req)

	// Should not return error, but should return original content
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 1)

	if textContent, ok := resp.Messages[0].Content.(*mcp.TextContent); ok {
		// Should return original content when template parsing fails
		assert.Equal(t, "This has {{invalid syntax", textContent.Text)
	} else {
		t.Error("Expected TextContent type")
	}
}
