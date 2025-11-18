package stdio

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	"github.com/hkionline/prompter/internal/configuration"
	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/hkionline/prompter/internal/rpc"
)

// Server for stdio transport layer
func Serve(configuration configuration.Configuration) {
	// fmt.Println("Starting stdio server")

	p := plog.New(configuration.LogFile)

	p.Write(plog.SERVER, "starting stdio server")

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	p.Write(plog.SERVER, "setting up prompts db")
	db, err := promptsdb.New(promptsdb.FILE_SYSTEM_PROVIDER, configuration.Storage, configuration.LogFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize prompts db connection: %s", err)
		os.Exit(-1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	p.Write(plog.SERVER, "stdio server set up, starting to listen messages")

	go func() {
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				p.Write(plog.SERVER, "Error reading from stdin", err.Error())
				return
			}
			p.Write(plog.CLIENT, message)

			resp := rpc.Process(message, db, p)

			if resp != "" {
				p.Write(plog.SERVER, resp)

				_, err = writer.WriteString(resp)
				if err != nil {
					p.Write(plog.SERVER, "Error writing to stdout", err.Error())
					return
				}
			}
			err = writer.Flush()
			if err != nil {
				p.Write(plog.SERVER, "Error flushing buffer", err.Error())
				return
			}
		}
	}()

	go func() {
		s := <-signalChan
		p.Write(plog.SERVER, "Received "+s.String()+" signal, shutting down stdio server.")
		writer.Flush()
		os.Exit(0)
	}()

	select {}
}
