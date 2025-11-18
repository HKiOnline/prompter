package rpc

import "github.com/hkionline/prompter/internal/promptsdb"

type Message struct {
	Version string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Method  string `json:"method,omitempty"`
	Params  Params `json:"params,omitzero"`
	Result  Result `json:"result,omitzero"`
	Error   Error  `json:"error,omitzero"`
}

type Params struct {
	Name            string           `json:"name,omitempty"`
	ProtocolVersion string           `json:"protocolVersion,omitempty"`
	Capabilities    Capabilities     `json:"capabilities,omitzero"`
	Arguments       promptsdb.Prompt `json:"arguments,omitzero"`
	ClientInfo      Info             `json:"clientInfo,omitzero"`
}

type Result struct {
	ProtocolVersion string             `json:"protocolVersion,omitempty"`
	Description     string             `json:"description,omitempty"`
	Capabilities    Capabilities       `json:"capabilities,omitzero"`
	Tools           []Tool             `json:"tools,omitzero"`
	Prompts         []promptsdb.Prompt `json:"prompts,omitzero"`
	ServerInfo      Info               `json:"serverInfo,omitzero"`
	Instructions    string             `json:"instructions,omitempty"`
	Content         []Content          `json:"content,omitzero"`
	Messages        []ResultMessage    `json:"messages,omitzero"`
	IsError         bool               `json:"isError,omitempty"`
}

type Content struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type ResultMessage struct {
	Role    string  `json:"role,omitempty"`
	Content Content `json:"content,omitzero"`
}

type Capabilities struct {
	Roots       Roots       `json:"roots,omitzero"`
	Sampling    Sampling    `json:"sampling,omitzero"`
	Elicitation Elicitation `json:"elicitation,omitzero"`
	Logging     Logging     `json:"logging,omitzero"`
	Prompts     Prompts     `json:"prompts,omitzero"`
	Resources   Resources   `json:"resources,omitzero"`
	Tools       Tools       `json:"tools,omitzero"`
}

type Roots struct {
	ListChanged bool `json:"listChanged"`
}

type Sampling struct{}

type Elicitation struct{}

type Logging struct{}

type Prompts struct {
	ListChanged bool `json:"listChanged"`
}

type Prompt struct {
	Name        string     `json:"name,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Arguments   []Argument `json:"arguments,omitzero"`
}

type Argument struct{}

type Resources struct {
	Subscribe   bool `json:"subscribe"`
	ListChanged bool `json:"listChanged"`
}

type Tools struct {
	ListChanged bool `json:"listChanged"`
}

type Tool struct {
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties NewPromptProperties `json:"properties,omitzero"`
}

type NewPromptProperties struct {
	Name        Property `json:"name"`
	Title       Property `json:"title"`
	Description Property `json:"description"`
	Content     Property `json:"content"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Info struct {
	Name    string `json:"name,omitempty"`
	Title   string `json:"title,omitempty"`
	Version string `json:"version,omitempty"`
}

type Error struct {
	Code    int    `json:"error"`
	Message string `json:"title"`
}

// type Data struct {}
