package rpc

import "github.com/hkionline/prompter/internal/promptsdb"

func toolsCall(msg Message, db promptsdb.Provider) Message {

	switch msg.Params.Name {
	case CREATE_PROMPT:
		return saveNewPrompt(msg, db)
	default:
		errMsg := errorMsg(msg.Id, UNSUPPORTED_TOOL, UNSUPPORTED_TOOL_MSG)
		errMsg.Result.IsError = true
		return errMsg
	}

}
