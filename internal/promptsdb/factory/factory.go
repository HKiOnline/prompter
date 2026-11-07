package factory

import (
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/hkionline/prompter/internal/promptsdb/sqlite"
)

// New creates a Provider based on the given provider name and configuration.
// Falls back to the filesystem provider for unknown or empty provider names.
func New(dbProvider string, config promptsdb.ProviderConfiguration, logfile string) (promptsdb.Provider, error) {

	switch dbProvider {
	case promptsdb.SQLITE_PROVIDER:
		return sqlite.NewSQLiteProvider(sqlite.SQLiteProviderConfiguration{
			Path:         config.SQLite.Path,
			Timeout:      config.SQLite.Timeout,
			ForeignKeys:  config.SQLite.ForeignKeys,
			MaxOpenConns: config.SQLite.MaxOpenConns,
			MaxIdleConns: config.SQLite.MaxIdleConns,
		}, logfile)
	case promptsdb.FILE_SYSTEM_PROVIDER, "":
		return promptsdb.NewPromptsFsProvider(config.Filesystem.Directory, logfile)
	default:
		// Unknown provider: fall back to filesystem
		return promptsdb.NewPromptsFsProvider(config.Filesystem.Directory, logfile)
	}
}
