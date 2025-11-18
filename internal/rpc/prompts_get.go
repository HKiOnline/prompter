package rpc

import (
	"github.com/hkionline/prompter/internal/promptsdb"
)

func promptsGet(msg Message, db promptsdb.Provider) Message {

	if msg.Params.Name == "" {
		// return an error for not supporting the prompt
		return errorMsg(msg.Id, 501, "missing prompt name")
	}

	prompt, err := db.Read(msg.Params.Name)

	if err != nil {
		// return an error if the prompt with a given name could not be found
		return errorMsg(msg.Id, 404, "prompt with the given name "+msg.Params.Name+" could not be found")
	}

	return Message{
		Version: JSONRPC_VERSION,
		Id:      msg.Id,
		Result: Result{
			Description: prompt.Description,
			Messages: []ResultMessage{
				{
					Role: "user",
					Content: Content{
						Type: "text",
						Text: prompt.Content,
					},
				},
			},
		},
	}

}
