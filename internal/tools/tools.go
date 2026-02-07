package tools

import (
	"context"
	"fmt"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	CREATE_PROMPT = "create_prompt" // tool call name for creating and storing a new prompt
)

// ToolHandler handles MCP tool requests
type ToolHandler struct {
	db     promptsdb.Provider
	logger *plog.Plogger
}

// NewToolHandler creates a new ToolHandler instance
func NewToolHandler(db promptsdb.Provider, logger *plog.Plogger) *ToolHandler {
	return &ToolHandler{
		db:     db,
		logger: logger,
	}
}

// CreatePromptTool creates the tool definition for creating prompts
func (h *ToolHandler) CreatePromptTool() *mcp.Tool {
	// Create the tool definition
	return &mcp.Tool{
		Name:        CREATE_PROMPT,
		Title:       "Create prompt",
		Description: "Create and save a new prompt",
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"name": {
					Type:        "string",
					Description: "Computer readable name for the prompt. White space should be replaced with underscores and special characters omitted.",
				},
				"title": {
					Type:        "string",
					Description: "Human readable display name for the prompt. This should be kept short and to the point.",
				},
				"description": {
					Type:        "string",
					Description: "Human readable explanation what the prompt is for expanding the title.",
				},
				"content": {
					Type:        "string",
					Description: "Full content of the prompt to be created and stored.",
				},
			},
		},
	}
}

// HandleCall handles the tools/call request
func (h *ToolHandler) HandleCall(ctx context.Context, ss *mcp.ServerSession, req *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResult, error) {
	h.logger.Write(plog.CLIENT, "tools/call")

	if req.Name == CREATE_PROMPT {
		return h.handleSaveNewPrompt(req)
	}

	h.logger.Write(plog.SERVER, "Unsupported tool: %s", req.Name)
	return nil, fmt.Errorf("unsupported tool: %s", req.Name)
}

// handleSaveNewPrompt handles the saveNewPrompt tool call
func (h *ToolHandler) handleSaveNewPrompt(req *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResult, error) {
	if req.Arguments == nil {
		h.logger.Write(plog.SERVER, "Missing arguments for saveNewPrompt")
		return nil, fmt.Errorf("missing arguments for saveNewPrompt")
	}

	// Extract arguments from the map
	args := req.Arguments
	name, _ := args["name"].(string)
	title, _ := args["title"].(string)
	description, _ := args["description"].(string)
	content, _ := args["content"].(string)

	if name == "" {
		h.logger.Write(plog.SERVER, "Missing prompt name")
		return nil, fmt.Errorf("missing prompt name")
	}

	// Create new prompt
	prompt := promptsdb.Prompt{
		Name:        name,
		Title:       title,
		Description: description,
		Content:     content,
	}

	err := h.db.Create(prompt)
	if err != nil {
		h.logger.Write(plog.SERVER, "Failed to create prompt: %s", err.Error())
		return nil, fmt.Errorf("failed to create prompt: %w", err)
	}

	responseText := fmt.Sprintf("created new prompt with name '%s' and title '%s'", name, title)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: responseText,
			},
		},
		IsError: false,
	}, nil
}
