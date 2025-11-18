package main

import (
	"fmt"
	"os"

	"github.com/hkionline/prompter/internal/configuration"
	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/transport/stdio"
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

	go stdio.Serve(config)
	select {}
}
