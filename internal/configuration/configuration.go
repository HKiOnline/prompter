package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

var (
	knf     = koanf.New(".") // Global active configuration instance of Koanf
	yparser = yaml.Parser()  // Koanf's yaml parser
)

// ConfigurationFile is used as the root for configuration files.
// This struct is not used within the service in any other purpose.
type ConfigurationFile struct {
	Configuration Configuration `koanf:"prompter" yaml:"prompter"`
}

// The configuration struct is the default data structure for all configurations.
// This is the struct you'll be mostly accessing from the service.
type Configuration struct {
	Transport TransportConfiguration          `yaml:"transport" koanf:"transport"`
	LogFile   string                          `yaml:"logFile" koanf:"logFile"`
	Storage   promptsdb.ProviderConfiguration `yaml:"storage" koanf:"storage"`
}

type TransportConfiguration struct {
	Type           string                      `yaml:"type" koanf:"type"`
	StreamableHTTP StreamableHTTPConfiguration `yaml:"streamable_http" koanf:"streamable_http"`
}

type StreamableHTTPConfiguration struct {
	Port int `yaml:"port" koanf:"port"`
}

// New default configuration for the service with a default configuration provider
func New(configFilePath string) (Configuration, error) {

	// Reset the global koanf instance to avoid state leaking between calls
	knf = koanf.New(".")

	// Load default configuration
	knf.Load(structs.Provider(ConfigurationFile{GetDefault()}, "koanf"), nil)

	// Load the specified config file
	if err := knf.Load(file.Provider(configFilePath), yparser); err != nil {
		return Configuration{}, fmt.Errorf("error loading config: %v", err)
	}

	if knf.Exists("prompter.http") {
		return Configuration{}, fmt.Errorf("legacy config key 'prompter.http' is not supported. Use 'prompter.transport.streamable_http'")
	}

	rawTransport := knf.Get("prompter.transport")
	if _, isScalarTransport := rawTransport.(string); isScalarTransport {
		return Configuration{}, fmt.Errorf("legacy config format for 'prompter.transport' is not supported. Use object format with 'prompter.transport.type'")
	}

	// Unmarshal the entire file, must be a yaml-file
	var kfile ConfigurationFile
	knf.Unmarshal("", &kfile)

	// Validate transport field
	if kfile.Configuration.Transport.Type != "stdio" && kfile.Configuration.Transport.Type != "streamable_http" {
		return Configuration{}, fmt.Errorf("invalid transport type: %s. Must be 'stdio' or 'streamable_http'", kfile.Configuration.Transport.Type)
	}

	// Validate storage provider
	provider := kfile.Configuration.Storage.Provider
	if provider == "" {
		provider = promptsdb.FILE_SYSTEM_PROVIDER
	}
	if provider != promptsdb.FILE_SYSTEM_PROVIDER && provider != promptsdb.SQLITE_PROVIDER {
		return Configuration{}, fmt.Errorf("invalid storage provider: %s. Must be 'filesystem' or 'sqlite'", provider)
	}

	// Apply environment variable overrides
	applyEnvOverrides(&kfile.Configuration)

	// Validate SQLite-specific config if sqlite provider is selected
	if kfile.Configuration.Storage.Provider == promptsdb.SQLITE_PROVIDER {
		if kfile.Configuration.Storage.SQLite.Path == "" {
			return Configuration{}, fmt.Errorf("sqlite storage provider requires storage.sqlite.path to be set")
		}
	}

	return kfile.Configuration, nil
}

// applyEnvOverrides overrides configuration values from environment variables.
// Supported env vars:
//   PROMPTER_STORAGE_PROVIDER   - storage provider (filesystem|sqlite)
//   PROMPTER_SQLITE_PATH        - SQLite database file path
//   PROMPTER_SQLITE_TIMEOUT     - SQLite connection timeout in seconds
//   PROMPTER_LOG_FILE           - log file path
func applyEnvOverrides(config *Configuration) {
	if v := os.Getenv("PROMPTER_STORAGE_PROVIDER"); v != "" {
		config.Storage.Provider = v
	}
	if v := os.Getenv("PROMPTER_SQLITE_PATH"); v != "" {
		config.Storage.SQLite.Path = v
	}
	if v := os.Getenv("PROMPTER_SQLITE_TIMEOUT"); v != "" {
		var timeout int
		if _, err := fmt.Sscanf(v, "%d", &timeout); err == nil {
			config.Storage.SQLite.Timeout = timeout
		}
	}
	if v := os.Getenv("PROMPTER_LOG_FILE"); v != "" {
		config.LogFile = v
	}
	if strings.ToLower(os.Getenv("PROMPTER_SQLITE_FOREIGN_KEYS")) == "true" {
		config.Storage.SQLite.ForeignKeys = true
	}
}
