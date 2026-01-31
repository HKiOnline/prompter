package prompts

import (
	"context"
	"fmt"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// PromptHandler handles MCP prompt requests
type PromptHandler struct {
	db     promptsdb.Provider
	logger *plog.Plogger
}

// NewPromptHandler creates a new PromptHandler instance
func NewPromptHandler(db promptsdb.Provider, logger *plog.Plogger) *PromptHandler {
	return &PromptHandler{
		db:     db,
		logger: logger,
	}
}

// HandleList handles the prompts/list request
func (h *PromptHandler) HandleList(ctx context.Context, ss *mcp.ServerSession, req *mcp.ListPromptsParams) (*mcp.ListPromptsResult, error) {
	h.logger.Write(plog.CLIENT, "prompts/list")

	prompts, err := h.db.List(promptsdb.PromptQuery{All: true})
	if err != nil {
		h.logger.Write(plog.SERVER, "Error listing prompts: %s", err.Error())
		return nil, fmt.Errorf("failed to list prompts: %w", err)
	}

	// Convert to MCP prompt format
	mcpPrompts := make([]*mcp.Prompt, len(prompts))
	for i, prompt := range prompts {
		mcpPrompts[i] = &mcp.Prompt{
			Name:        prompt.Name,
			Title:       prompt.Title,
			Description: prompt.Description,
		}
	}

	return &mcp.ListPromptsResult{
		Prompts: mcpPrompts,
	}, nil
}

// HandleGet handles the prompts/get request
func (h *PromptHandler) HandleGet(ctx context.Context, ss *mcp.ServerSession, req *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	h.logger.Write(plog.CLIENT, "prompts/get")

	if req.Name == "" {
		h.logger.Write(plog.SERVER, "Missing prompt name")
		return nil, fmt.Errorf("missing prompt name")
	}

	prompt, err := h.db.Read(req.Name)
	if err != nil {
		h.logger.Write(plog.SERVER, "Prompt not found: %s", err.Error())
		return nil, fmt.Errorf("prompt with name %s not found: %w", req.Name, err)
	}

	return &mcp.GetPromptResult{
		Description: prompt.Description,
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: prompt.Content,
				},
			},
		},
	}, nil
}