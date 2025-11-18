package promptsdb

type Prompt struct {
	Id          string   `json:"-" yaml:"id"`                              // Unique computer readable datastorage engine identifier
	Name        string   `json:"name,omitempty" yaml:"name"`               // Unique programmatic or logical name used to invoke the prompt
	Title       string   `json:"title,omitempty" yaml:"title"`             // Human readable title of the prompt
	Description string   `json:"description,omitempty" yaml:"description"` // Human readable longer explanation what the prompt is
	Arguments   []string `json:"arguments,omitzero" yaml:"arguments"`      // Arguments used in invoking the prompt
	Content     string   `json:"content" yaml:"-"`                         // The contents of the actual prompt
	Tags        []string `json:"-" yaml:"tags"`                            // Tags for the prompt, can be used for example in completion suggestions
}

type PromptQuery struct {
	All            bool
	NameStartsWith string
	NameContains   string
	IndexFrom      int
	IndexTo        int
}

type PromptsDBError struct{}

func (p *PromptsDBError) Error() string {

	return ""
}
