package rpc

const (
	JSONRPC_VERSION   = "2.0"
	PROTOCOL_VERSION  = "2025-06-18"
	SERVER_INFO_NAME  = "prompter"
	SERVER_INFO_TITLE = "Prompter - MCP Server for prompts"
	SERVER_VERSION    = "0.0.0-alpha"
)

func initialization(msg Message) Message {
	/*
	   {
	     "jsonrpc": "2.0",
	     "id": 1,
	     "result": {
	       "protocolVersion": "2024-11-05",
	       "capabilities": {
	         "logging": {},
	         "prompts": {
	           "listChanged": true
	         },
	         "resources": {
	           "subscribe": true,
	           "listChanged": true
	         },
	         "tools": {
	           "listChanged": true
	         }
	       },
	       "serverInfo": {
	         "name": "ExampleServer",
	         "title": "Example Server Display Name",
	         "version": "1.0.0"
	       },
	       "instructions": "Optional instructions for the client"
	     }
	   }
	*/

	return Message{
		Version: JSONRPC_VERSION,
		Id:      msg.Id,
		Result: Result{
			ProtocolVersion: msg.Params.ProtocolVersion,
			Capabilities: Capabilities{
				Prompts: Prompts{
					ListChanged: true,
				},
				Tools: Tools{
					ListChanged: true,
				},
			},
			ServerInfo: Info{
				Name:    SERVER_INFO_NAME,
				Title:   SERVER_INFO_TITLE,
				Version: SERVER_VERSION,
			},
		},
	}

}
