package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hkionline/prompter/internal/configuration"
	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/hkionline/prompter/internal/server"
)

func main() {

	// 1. Read configuration
	// 2. Setup service

	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration failure: %s", err)
		os.Exit(-1)
	}

	config, err := configuration.New(homeDir + "/.config/prompter/prompter.yaml")

	log := plog.New(config.LogFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration failure: %s", err)
		os.Exit(-1)
	}

	log.Write(plog.SERVER, "configuration read")

	// Initialize database
	log.Write(plog.SERVER, "setting up prompts db")
	db, err := promptsdb.New(promptsdb.FILE_SYSTEM_PROVIDER, config.Storage, config.LogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize prompts db connection: %s", err)
		os.Exit(-1)
	}

	// Create and start new prompter MCP-server
	prompter := server.New("0.5.0", &config, log, db)

	// Start the server with stdio transport
	err = prompter.Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run MCP server: %s", err)
		os.Exit(-1)
	}
}
