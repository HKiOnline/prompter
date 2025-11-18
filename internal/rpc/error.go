package rpc

const (
	GENERAL_ERROR_CODE = iota + 100
	FAILURE_TO_PARSE_MESSAGE
	UNSUPPORTED_VERSION
	UNSUPPORTED_METHOD
	UNSUPPORTED_PROMPT
	UNSUPPORTED_TOOL
)

const (
	GENERAL_ERROR_MSG            = "internal server failure"
	FAILURE_TO_PARSE_MESSAGE_MSG = "failed to parse the json-rpc message"
	UNSUPPORTED_VERSION_MSG      = "unsupported mcp version"
	UNSUPPORTED_METHOD_MSG       = "unsupported json-rpc method"
	UNSUPPORTED_PROMPT_MSG       = "unsupported prompt capability"
	UNSUPPORTED_TOOL_MSG         = "unsupported tools capability"
)

func errorMsg(id int, code int, errMsg string) Message {
	return Message{
		Version: JSONRPC_VERSION,
		Id:      id,
		Error: Error{
			Code:    code,
			Message: errMsg,
		},
	}
}
