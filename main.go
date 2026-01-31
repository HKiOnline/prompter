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

	config, err := configuration.Setup(homeDir + "/.config/prompter/prompter.yaml")

	p := plog.New(config.LogFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration failure: %s", err)
		os.Exit(-1)
	}

	p.Write(plog.SERVER, "configuration read")

	// Initialize database
	p.Write(plog.SERVER, "setting up prompts db")
	db, err := promptsdb.New(promptsdb.FILE_SYSTEM_PROVIDER, config.Storage, config.LogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize prompts db connection: %s", err)
		os.Exit(-1)
	}

	// Create and start new SDK server
	sdkServer := server.NewServer(&config, p, db)
	err = sdkServer.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start MCP server: %s", err)
		os.Exit(-1)
	}

	// Start the server with stdio transport
	err = sdkServer.Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run MCP server: %s", err)
		os.Exit(-1)
	}
}
