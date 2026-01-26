package plog

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	CLIENT = "client"
	SERVER = "prompter"
)

type Plogger struct {
	plogFile string
}

func New(plogFilePath string) *Plogger {
	return &Plogger{plogFile: plogFilePath}
}

func (p Plogger) Write(sender string, messages ...string) {

	file, err := os.OpenFile(p.plogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	errs := bufio.NewWriter(os.Stderr)

	if err != nil {
		errs.WriteString(fmt.Sprintln("Error opening log file:", err))
		return
	}

	defer file.Close()

	sb := strings.Builder{}
	sb.WriteString(time.Now().Format(time.RFC3339) + " ")
	sb.WriteString("[" + sender + "]: ")

	for i, message := range messages {

		if i == 0 {
			sb.WriteString(" " + message)
		} else {
			sb.WriteString(" :: " + message)
		}
	}

	_, err = file.WriteString(sb.String() + "\n")

	if err != nil {
		errs.WriteString(fmt.Sprintln("Error writing to log file:", err))
	}
}
