package rpc

const (
	CREATE_PROMPT = "create_prompt" // tool call name for creating and storing a new prompt
)

func toolsList(msg Message) Message {

	return Message{
		Version: JSONRPC_VERSION,
		Id:      msg.Id,
		Result: Result{
			Tools: []Tool{
				Tool{
					Name:        CREATE_PROMPT,
					Title:       "Create prompt",
					Description: "Create and save a new prompt",
					InputSchema: InputSchema{
						Type: "object",
						Properties: NewPromptProperties{
							Name: Property{
								Type:        "string",
								Description: "Computer readable name for the prompt. White space should be replaced with underscores and special characters omitted.",
							},
							Title: Property{
								Type:        "string",
								Description: "Human readable display name for the prompt. This should be kept short and to the point.",
							},
							Description: Property{
								Type:        "string",
								Description: "Human readable explanation what the prompt is for expanding the title.",
							},
							Content: Property{
								Type:        "string",
								Description: "Full content of the prompt to be created and stored.",
							},
						},
					},
				},
			},
		},
	}

}
