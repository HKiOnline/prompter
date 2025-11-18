package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
)

// Process the incoming JSON-RPC message and return a response
func Process(reqStr string, db promptsdb.Provider, p *plog.Plogger) string {

	var req Message
	var resp Message

	var noOp bool

	p.Write(plog.CLIENT, reqStr)

	err := json.Unmarshal([]byte(reqStr), &req)

	if err != nil {
		p.Write(plog.SERVER, "could not parse the message %s", err.Error())
		resp = errorMsg(-1, FAILURE_TO_PARSE_MESSAGE, "could not parse the message "+err.Error())

	} else {

		switch req.Method {
		case "initialize":
			p.Write(plog.CLIENT, "initialization requested")
			resp = initialization(req)

		case "notifications/initialized":
			p.Write(plog.CLIENT, "initialized")
			noOp = true

		case "ping":
			p.Write(plog.CLIENT, "ping")
			resp = ping(req)

		case "tools/list":
			p.Write(plog.CLIENT, "tools/list")
			resp = toolsList(req)

		case "tools/call":
			p.Write(plog.CLIENT, "tools/call")
			resp = toolsCall(req, db)

		case "prompts/list":
			p.Write(plog.CLIENT, "prompts/list")
			resp = promptsList(req, db)

		case "prompts/get":
			p.Write(plog.CLIENT, "prompts/list")
			resp = promptsGet(req, db)

		default:
			p.Write(plog.SERVER, req.Method+" is an unknown message")
			noOp = true
		}
	}

	if !noOp {
		respBytes, err := json.Marshal(resp)

		if err != nil {
			p.Write(plog.SERVER, "failed to marshal json response for the request", err.Error())
			return ""
		}

		respStr := fmt.Sprintln(string(respBytes))
		p.Write(plog.SERVER, respStr)

		return respStr
	} else {
		noOp = false
		return ""
	}
}
