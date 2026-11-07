package promptsdb

type Provider interface {
	Create(prompt Prompt) error
	Read(promptId string) (Prompt, error)
	Update(prompt Prompt) error
	Delete(promptId string) error
	List(query PromptQuery) ([]Prompt, error)
}

type ProviderConfiguration struct {
	Provider   string                  `yaml:"provider" koanf:"provider"`
	Filesystem FsProviderConfiguration `yaml:"filesystem" koanf:"filesystem"`
	SQLite     SQLiteConfiguration     `yaml:"sqlite" koanf:"sqlite"`
}

// SQLiteConfiguration holds SQLite-specific storage configuration.
type SQLiteConfiguration struct {
	Path         string `yaml:"path" koanf:"path"`
	Timeout      int    `yaml:"timeout" koanf:"timeout"`
	ForeignKeys  bool   `yaml:"foreign_keys" koanf:"foreign_keys"`
	MaxOpenConns int    `yaml:"max_open_conns" koanf:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns" koanf:"max_idle_conns"`
}
