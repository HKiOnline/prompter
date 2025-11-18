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
	Transport string                          `yaml:"transport" koanf:"transport"`
	LogFile   string                          `yaml:"logFile" koanf:"logFile"`
	Storage   promptsdb.ProviderConfiguration `yaml:"storage" koanf:"storage"`
}

// Setup default configuration for the service with a default configuration provider
func Setup(configFilePath string) (Configuration, error) {

	// Load default configuration
	knf.Load(structs.Provider(ConfigurationFile{GetDefault()}, "koanf"), nil)

	// Load the specified config file
	if err := knf.Load(file.Provider(configFilePath), yparser); err != nil {
		return Configuration{}, fmt.Errorf("error loading config: %v", err)
	}

	// Unmarshal the entire file, must be a yaml-file
	var kfile ConfigurationFile
	knf.Unmarshal("", &kfile)

	return kfile.Configuration, nil
}
