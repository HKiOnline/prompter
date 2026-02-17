package server

import (
	"context"
	"fmt"

	"github.com/hkionline/prompter/internal/configuration"
	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/prompts"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/hkionline/prompter/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Prompter implements the MCP server using the official Go SDK
type Prompter struct {
	version string
	config  *configuration.Configuration
	logger  *plog.Plogger
	db      promptsdb.Provider
	server  *mcp.Server
	prompts *prompts.PromptHandler
	tools   *tools.ToolHandler
}

// New creates a new SDKServer instance
func New(version string, config *configuration.Configuration, logger *plog.Plogger, db promptsdb.Provider) *Prompter {
	return &Prompter{
		version: version,
		config:  config,
		logger:  logger,
		db:      db,
	}
}

// Run starts the MCP server with the configured transport
func (s *Prompter) Run(ctx context.Context) error {
	s.logger.Write(plog.SERVER, "initializing prompter MCP-server")

	// Create MCP server instance
	server := mcp.NewServer("prompter", s.version, &mcp.ServerOptions{})

	s.server = server

	s.logger.Write(plog.SERVER, "attaching capability handlers to the server")

	// Initialize handlers
	s.prompts = prompts.NewPromptHandler(s.db, s.logger)
	s.tools = tools.NewToolHandler(s.db, s.logger)

	// Add all prompts from the database to the server
	promptsList, err := s.db.List(promptsdb.PromptQuery{All: true})
	if err != nil {
		s.logger.Write(plog.SERVER, "failed to load prompts for registration: %s", err.Error())
	} else {
		for _, prompt := range promptsList {
			s.server.AddPrompts(
				&mcp.ServerPrompt{
					Prompt: &mcp.Prompt{
						Name:        prompt.Name,
						Title:       prompt.Title,
						Description: prompt.Description,
					},
					Handler: s.prompts.HandleGet,
				},
			)
		}
	}

	// Add tools to the server
	tool := s.tools.CreatePromptTool()

	s.server.AddTools(
		&mcp.ServerTool{
			Tool:    tool,
			Handler: s.tools.HandleCall,
		},
	)

	// Create transport based on configuration
	s.logger.Write(plog.SERVER, "creating transport: %s", s.config.Transport.Type)
	trans, err := newTransport(s.config.Transport.Type)
	if err != nil {
		return fmt.Errorf("failed to create transport: %w", err)
	}

	s.logger.Write(plog.SERVER, "starting the MCP server with %s transport", s.config.Transport.Type)
	return trans.start(ctx, server, s.config)
}

// GetServer returns the underlying MCP server instance
func (s *Prompter) GetServer() *mcp.Server {
	return s.server
}
