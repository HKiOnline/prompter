package rpc

func ping(msg Message) Message {

	return Message{
		Version: JSONRPC_VERSION,
		Id:      msg.Id,
		Result:  Result{},
	}

}
