package rpc

import (
	"fmt"

	"github.com/hkionline/prompter/internal/promptsdb"
)

func saveNewPrompt(msg Message, db promptsdb.Provider) Message {

	if msg.Params.Name == "" {
		// return an error for not supporting the prompt
		return errorMsg(msg.Id, 501, "missing a name for the new prompt")
	}

	err := db.Create(msg.Params.Arguments)

	if err != nil {
		// return an error if the prompt with a given name could not be found
		return errorMsg(msg.Id, 501, "prompt with the given name "+msg.Params.Name+" could not be created")
	}

	content := fmt.Sprintf("created new prompt with name '%s' and title '%s'", msg.Params.Arguments.Name, msg.Params.Arguments.Title)

	return Message{
		Version: JSONRPC_VERSION,
		Id:      msg.Id,
		Result: Result{
			Content: []Content{
				{
					Type: "text",
					Text: content,
				},
			},
			IsError: false,
		},
	}
}
