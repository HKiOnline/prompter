//go:build !sqlite

package sqlite

import (
	"errors"

	"github.com/hkionline/prompter/internal/promptsdb"
)

// SQLiteProviderConfiguration holds SQLite-specific configuration options.
// This is a stub used when the "sqlite" build tag is not set.
type SQLiteProviderConfiguration struct {
	Path         string `yaml:"path"`
	Timeout      int    `yaml:"timeout"`
	ForeignKeys  bool   `yaml:"foreign_keys"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// NewSQLiteProvider returns an error when built without the "sqlite" build tag.
func NewSQLiteProvider(config SQLiteProviderConfiguration, logfile string) (promptsdb.Provider, error) {
	return nil, errors.New("sqlite provider: this binary was built without SQLite support; rebuild with -tags sqlite")
}
