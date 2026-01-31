package server

import (
	"context"

	"github.com/hkionline/prompter/internal/configuration"
	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/prompts"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/hkionline/prompter/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SDKServer implements the MCP server using the official Go SDK
type SDKServer struct {
	config  *configuration.Configuration
	logger  *plog.Plogger
	db      promptsdb.Provider
	server  *mcp.Server
	prompts *prompts.PromptHandler
	tools   *tools.ToolHandler
}

// NewServer creates a new SDKServer instance
func NewServer(config *configuration.Configuration, logger *plog.Plogger, db promptsdb.Provider) *SDKServer {
	return &SDKServer{
		config: config,
		logger: logger,
		db:     db,
	}
}

// Start initializes and starts the MCP server
func (s *SDKServer) Start() error {
	s.logger.Write(plog.SERVER, "Initializing MCP SDK server")

	// Create MCP server instance
	server := mcp.NewServer("prompter", "0.0.0-alpha", &mcp.ServerOptions{})

	s.server = server

	// Initialize handlers
	s.prompts = prompts.NewPromptHandler(s.db, s.logger)
	s.tools = tools.NewToolHandler(s.db, s.logger)

	// Register handlers
	s.registerHandlers()

	s.logger.Write(plog.SERVER, "MCP SDK server initialized successfully")
	return nil
}

// registerHandlers registers all MCP handlers with the server
func (s *SDKServer) registerHandlers() {
	// Initialize handlers
	s.prompts = prompts.NewPromptHandler(s.db, s.logger)
	s.tools = tools.NewToolHandler(s.db, s.logger)

	// Add all prompts from the database to the server
	promptsList, err := s.db.List(promptsdb.PromptQuery{All: true})
	if err != nil {
		s.logger.Write(plog.SERVER, "Failed to load prompts for registration: %s", err.Error())
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
	s.server.AddTools(
		&mcp.ServerTool{
			Tool: &mcp.Tool{
				Name:        "create_prompt",
				Title:       "Create prompt",
				Description: "Create and save a new prompt",
			},
			Handler: s.tools.HandleCall,
		},
	)
}

// Run starts the MCP server with stdio transport
func (s *SDKServer) Run(ctx context.Context) error {
	s.logger.Write(plog.SERVER, "Starting MCP server with stdio transport")

	transport := mcp.NewStdioTransport()
	return s.server.Run(ctx, transport)
}

// GetServer returns the underlying MCP server instance
func (s *SDKServer) GetServer() *mcp.Server {
	return s.server
}
