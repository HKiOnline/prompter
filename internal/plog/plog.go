package plog

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
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
	attrs    []slog.Attr
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

// Enabled implements slog.Handler.
func (h *Plogger) Enabled(_ context.Context, _ slog.Level) bool {
	// Always enabled for simplicity; level filtering can be added here if needed.
	return true
}

// WithAttrs implements slog.Handler.
func (h *Plogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Plogger{
		plogFile: h.plogFile,
		attrs:    append(h.attrs, attrs...),
	}
}

// WithGroup implements slog.Handler.
func (h *Plogger) WithGroup(name string) slog.Handler {
	// Not implemented for simplicity; group support can be added if needed.
	return h
}

// Handle implements slog.Handler.
func (h *Plogger) Handle(_ context.Context, r slog.Record) error {
	/*
		file, err := os.OpenFile(h.plogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening log file:", err)
			return err
		}
		defer file.Close()

		sb := strings.Builder{}
		sb.WriteString(r.Time.Format(time.RFC3339) + " ")
		sb.WriteString("[" + r.PCFile() + "]: ") // or use a custom sender field if you have one

		// Append the log message
		sb.WriteString(r.Message)

		// Append all attributes
		r.Attrs(func(a slog.Attr) bool {
			sb.WriteString(" :: " + a.Key + "=" + a.Value.String())
			return true
		})

		// Append newlines from all attrs in h.attrs (context)
		for _, attr := range h.attrs {
			sb.WriteString(" :: " + attr.Key + "=" + attr.Value.String())
		}

		_, err = file.WriteString(sb.String() + "\n")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error writing to log file:", err)
			return err
		}
	*/

	return nil
}
