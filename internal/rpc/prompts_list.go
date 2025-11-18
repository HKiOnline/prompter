package rpc

import "github.com/hkionline/prompter/internal/promptsdb"

func promptsList(msg Message, db promptsdb.Provider) Message {

	prompts, err := db.List(promptsdb.PromptQuery{All: true})

	if err != nil {
		// return an error if the prompt with a given name could not be found
		return errorMsg(msg.Id, 404, "prompts could not be listed: "+err.Error())
	}
	return Message{
		Version: JSONRPC_VERSION,
		Id:      msg.Id,
		Result: Result{
			Prompts: prompts,
		},
	}
}
