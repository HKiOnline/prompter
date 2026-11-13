package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hkionline/prompter/internal/configuration"
	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/hkionline/prompter/internal/server"
	"github.com/hkionline/prompter/internal/version"

	// blank import so the sqlite package's init() registers its constructor
	_ "github.com/hkionline/prompter/internal/promptsdb/sqlite"
)

//go:embed version.json
var embeddedVersionJSON []byte

func main() {
	os.Exit(run(os.Args[1:], embeddedVersionJSON, os.Stdout, os.Stderr))
}

// run handles CLI flags and startup, returning the process exit code.
// Flags and embedded metadata are handled before any configuration, logging,
// or storage initialization.
func run(args []string, embedded []byte, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("prompter", flag.ContinueOnError)
	fs.SetOutput(stderr)
	versionFlag := fs.Bool("version", false, "print version information and exit")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "prompter: unexpected argument %q\n", fs.Arg(0))
		return 2
	}

	meta, err := version.Load(embedded)
	if err != nil {
		fmt.Fprintf(stderr, "prompter: %s\n", err)
		return 1
	}

	if *versionFlag {
		fmt.Fprintln(stdout, meta.Line())
		return 0
	}

	return startServer(meta, stderr)
}

func startServer(meta version.Metadata, stderr io.Writer) int {

	// 1. Read configuration
	// 2. Setup service

	homeDir, err := os.UserHomeDir()

	if err != nil {
		fmt.Fprintf(stderr, "configuration failure: %s", err)
		return 1
	}

	config, err := configuration.New(homeDir + "/.config/prompter/prompter.yaml")

	if err != nil {
		fmt.Fprintf(stderr, "configuration failure: %s", err)
		return 1
	}

	log := plog.New(config.LogFile)

	log.Write(plog.SERVER, "configuration read")

	// Initialize database
	log.Write(plog.SERVER, "setting up prompts db")
	providerName := config.Storage.Provider
	if providerName == "" {
		providerName = promptsdb.FILE_SYSTEM_PROVIDER
	}
	db, err := promptsdb.New(providerName, config.Storage, config.LogFile)
	if err != nil {
		fmt.Fprintf(stderr, "failed to initialize prompts db connection: %s", err)
		return 1
	}

	// Create and start new prompter MCP-server. The MCP serverInfo.version is
	// the composed version without build metadata (MAJOR.MINOR.PATCH plus any
	// -preRelease).
	prompter := server.New(meta.Core(), &config, log, db)

	// Start the server with stdio transport
	err = prompter.Run(context.Background())
	if err != nil {
		fmt.Fprintf(stderr, "failed to run MCP server: %s", err)
		return 1
	}

	return 0
}
