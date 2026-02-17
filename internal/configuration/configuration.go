package configuration

import (
	"fmt"

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

	return kfile.Configuration, nil
}
