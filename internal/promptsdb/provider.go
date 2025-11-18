package promptsdb

type Provider interface {
	Create(prompt Prompt) error
	Read(promptId string) (Prompt, error)
	Update(prompt Prompt) error
	Delete(promptId string) error
	List(query PromptQuery) ([]Prompt, error)
}

type ProviderConfiguration struct {
	Provider   string                  `yaml:"provider"`
	Filesystem FsProviderConfiguration `yaml:"filesystem"`
}
