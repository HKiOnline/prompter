package promptsdb

import (
	"fmt"
)

const (
	FILE_SYSTEM_PROVIDER = "filesystem"
	SQLITE_PROVIDER      = "sqlite"
)

// ProviderConstructor builds a Provider from the shared provider configuration.
type ProviderConstructor func(config ProviderConfiguration, logfile string) (Provider, error)

var providerConstructors = map[string]ProviderConstructor{}

// RegisterProvider registers a ProviderConstructor for the given provider name.
// Sub-packages (e.g. sqlite) call this from init() to avoid an import cycle,
// since they import promptsdb for the Provider interface.
func RegisterProvider(name string, ctor ProviderConstructor) {
	providerConstructors[name] = ctor
}

// New creates a Provider based on the given provider name and configuration.
// Falls back to the filesystem provider for unknown or empty provider names.
func New(dbProvider string, config ProviderConfiguration, logfile string) (Provider, error) {

	switch dbProvider {
	case SQLITE_PROVIDER:
		if ctor, ok := providerConstructors[SQLITE_PROVIDER]; ok {
			return ctor(config, logfile)
		}
		return nil, fmt.Errorf("promptsdb: no provider registered for %q (import the provider package so its init() runs, e.g. github.com/hkionline/prompter/internal/promptsdb/sqlite)", dbProvider)
	case FILE_SYSTEM_PROVIDER, "":
		return NewPromptsFsProvider(config.Filesystem.Directory, logfile)
	default:
		// Unknown provider: fall back to filesystem
		return NewPromptsFsProvider(config.Filesystem.Directory, logfile)
	}
}
