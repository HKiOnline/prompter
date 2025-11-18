package configuration

import (
	"testing"
)

func TestGetDefault(t *testing.T) {
	// Test that GetDefault returns a valid configuration with expected defaults
	config := GetDefault()

	// Verify transport is set to "stdio" by default
	if config.Transport != "stdio" {
		t.Errorf("Expected transport 'stdio', got '%s'", config.Transport)
	}

	// Verify log file path is set (we can't fully test the actual path without home dir)
	if config.LogFile == "" {
		t.Error("Expected log file to be set")
	}

	// Verify storage configuration is properly initialized
	if config.Storage.Provider == "" {
		t.Error("Expected storage provider to be set")
	}

	if config.Storage.Filesystem.Directory == "" {
		t.Error("Expected filesystem directory to be set")
	}
}

func TestConfigurationStruct(t *testing.T) {
	// Create a simple Configuration struct to test fields
	config := Configuration{
		Transport: "stdio",
		LogFile:   "/tmp/prompter.log",
	}

	if config.Transport != "stdio" {
		t.Errorf("Expected transport 'stdio', got '%s'", config.Transport)
	}

	if config.LogFile != "/tmp/prompter.log" {
		t.Errorf("Expected log file '/tmp/prompter.log', got '%s'", config.LogFile)
	}
}

func TestDefaultConfigurationValues(t *testing.T) {
	// More specific test of default values
	defaultConfig := GetDefault()

	// Test that defaults are set correctly
	if defaultConfig.Transport != "stdio" {
		t.Errorf("Expected default transport 'stdio', got '%s'", defaultConfig.Transport)
	}

	// Test that storage provider is filesystem
	if defaultConfig.Storage.Provider != "filesystem" {
		t.Errorf("Expected default storage provider 'filesystem', got '%s'", defaultConfig.Storage.Provider)
	}

	// Test that filesystem directory is set
	if defaultConfig.Storage.Filesystem.Directory == "" {
		t.Error("Expected filesystem directory to be set in defaults")
	}

	// Test that log file is set
	if defaultConfig.LogFile == "" {
		t.Error("Expected log file to be set in defaults")
	}
}

func TestConfigurationFields(t *testing.T) {
	// Test that the Configuration struct has all expected fields
	config := Configuration{
		Transport: "stdio",
		LogFile:   "/tmp/test.log",
	}

	if config.Transport != "stdio" {
		t.Errorf("Transport field not set correctly")
	}

	if config.LogFile != "/tmp/test.log" {
		t.Errorf("LogFile field not set correctly")
	}

	// Test that the struct can be used for basic operations
	if config.Transport == "" {
		t.Error("Transport field should not be empty")
	}

	if config.LogFile == "" {
		t.Error("LogFile field should not be empty")
	}
}

func TestConfigurationFileStruct(t *testing.T) {
	// Test the ConfigurationFile struct
	configFile := ConfigurationFile{
		Configuration: Configuration{
			Transport: "stdio",
			LogFile:   "/tmp/prompter.log",
		},
	}

	if configFile.Configuration.Transport != "stdio" {
		t.Errorf("Expected transport 'stdio', got '%s'", configFile.Configuration.Transport)
	}
}
