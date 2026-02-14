package prompts

import (
	"context"
	"testing"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestDateFunctionVariations(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		args     map[string]string
		expected string
	}{
		{
			name:     "Simple date",
			content:  "Date: {{date}}",
			args:     map[string]string{},
			expected: `Date: \d{4}-\d{2}-\d{2}`,
		},
		{
			name:     "Date with text",
			content:  "Today is {{date}} and it's a beautiful day!",
			args:     map[string]string{},
			expected: `Today is \d{4}-\d{2}-\d{2} and it's a beautiful day!`,
		},
		{
			name:     "Multiple dates",
			content:  "Start: {{date}}, End: {{date}}",
			args:     map[string]string{},
			expected: `Start: \d{4}-\d{2}-\d{2}, End: \d{4}-\d{2}-\d{2}`,
		},
		{
			name:     "Date with arguments",
			content:  "Hello {{.name}}, today is {{date}}",
			args:     map[string]string{"name": "Alice"},
			expected: `Hello Alice, today is \d{4}-\d{2}-\d{2}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testPrompt := promptsdb.Prompt{
				Name:        "test-" + tc.name,
				Title:       "Test " + tc.name,
				Description: "Test description",
				Content:     tc.content,
			}

			db := NewMockDB([]promptsdb.Prompt{testPrompt})
			logger := plog.New("/tmp/test.log")
			handler := NewPromptHandler(db, logger)

			req := &mcp.GetPromptParams{
				Name:      testPrompt.Name,
				Arguments: tc.args,
			}

			resp, err := handler.HandleGet(context.Background(), nil, req)

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Len(t, resp.Messages, 1)

			if textContent, ok := resp.Messages[0].Content.(*mcp.TextContent); ok {
				assert.Regexp(t, tc.expected, textContent.Text)
			} else {
				t.Error("Expected TextContent type")
			}
		})
	}
}
