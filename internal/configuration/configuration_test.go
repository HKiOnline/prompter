package configuration

import (
	"os"
	"strings"
	"testing"
)

func TestGetDefault(t *testing.T) {
	// Test that GetDefault returns a valid configuration with expected defaults
	config := GetDefault()

	// Verify transport is set to "stdio" by default
	if config.Transport.Type != "stdio" {
		t.Errorf("Expected transport type 'stdio', got '%s'", config.Transport.Type)
	}

	if config.Transport.StreamableHTTP.Port != 8080 {
		t.Errorf("Expected default streamable HTTP port 8080, got '%d'", config.Transport.StreamableHTTP.Port)
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
		Transport: TransportConfiguration{Type: "stdio"},
		LogFile:   "/tmp/prompter.log",
	}

	if config.Transport.Type != "stdio" {
		t.Errorf("Expected transport 'stdio', got '%s'", config.Transport.Type)
	}

	if config.LogFile != "/tmp/prompter.log" {
		t.Errorf("Expected log file '/tmp/prompter.log', got '%s'", config.LogFile)
	}
}

func TestDefaultConfigurationValues(t *testing.T) {
	// More specific test of default values
	defaultConfig := GetDefault()

	// Test that defaults are set correctly
	if defaultConfig.Transport.Type != "stdio" {
		t.Errorf("Expected default transport 'stdio', got '%s'", defaultConfig.Transport.Type)
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
		Transport: TransportConfiguration{Type: "stdio"},
		LogFile:   "/tmp/test.log",
	}

	if config.Transport.Type != "stdio" {
		t.Errorf("Transport field not set correctly")
	}

	if config.LogFile != "/tmp/test.log" {
		t.Errorf("LogFile field not set correctly")
	}

	// Test that the struct can be used for basic operations
	if config.Transport.Type == "" {
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
			Transport: TransportConfiguration{Type: "stdio"},
			LogFile:   "/tmp/prompter.log",
		},
	}

	if configFile.Configuration.Transport.Type != "stdio" {
		t.Errorf("Expected transport 'stdio', got '%s'", configFile.Configuration.Transport.Type)
	}
}

func TestConfigurationWithEmptyFields(t *testing.T) {
	// Test configuration with empty fields
	config := Configuration{}

	// These should be empty strings by default
	if config.Transport.Type != "" && config.LogFile != "" {
		t.Error("Expected empty fields for uninitialized configuration")
	}
}

func TestConfigurationWithCustomValues(t *testing.T) {
	// Test configuration with custom values
	config := Configuration{
		Transport: TransportConfiguration{Type: "custom"},
		LogFile:   "/custom/path.log",
	}

	if config.Transport.Type != "custom" {
		t.Errorf("Expected transport 'custom', got '%s'", config.Transport.Type)
	}

	if config.LogFile != "/custom/path.log" {
		t.Errorf("Expected log file '/custom/path.log', got '%s'", config.LogFile)
	}
}

func TestConfigurationStorageProvider(t *testing.T) {
	// Test storage provider configuration
	config := GetDefault()

	// Verify default storage provider is filesystem
	if config.Storage.Provider != "filesystem" {
		t.Errorf("Expected default storage provider 'filesystem', got '%s'", config.Storage.Provider)
	}

	// Verify filesystem directory is set
	if config.Storage.Filesystem.Directory == "" {
		t.Error("Expected filesystem directory to be set")
	}
}

func TestSetupWithValidConfigFile(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	configPath := tempDir + "/test.yaml"

	// Write a simple config file with only transport setting
	// (other settings will use defaults)
	configContent := `prompter:
  transport:
    type: "stdio"`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test Setup function
	config, err := New(configPath)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Verify the configuration was loaded correctly
	if config.Transport.Type != "stdio" {
		t.Errorf("Expected transport 'stdio', got '%s'", config.Transport.Type)
	}

	// Verify that default values are still applied for unset fields
	if config.LogFile == "" {
		t.Error("Expected log file to be set (from defaults)")
	}

	if config.Storage.Provider != "filesystem" {
		t.Errorf("Expected storage provider 'filesystem', got '%s'", config.Storage.Provider)
	}
}

func TestSetupRejectsLegacyScalarTransport(t *testing.T) {
	tempDir := t.TempDir()
	configPath := tempDir + "/test_legacy_transport.yaml"

	configContent := `prompter:
  transport: "stdio"`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := New(configPath)
	if err == nil {
		t.Fatalf("Expected error for legacy scalar transport config")
	}

	if !strings.Contains(err.Error(), "prompter.transport") {
		t.Fatalf("Expected legacy transport error, got: %v", err)
	}
}

func TestSetupRejectsLegacyTopLevelHTTP(t *testing.T) {
	tempDir := t.TempDir()
	configPath := tempDir + "/test_legacy_http.yaml"

	configContent := `prompter:
  transport:
    type: "streamable_http"
  http:
    port: 9090`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := New(configPath)
	if err == nil {
		t.Fatalf("Expected error for legacy top-level http config")
	}

	if !strings.Contains(err.Error(), "prompter.http") {
		t.Fatalf("Expected legacy http error, got: %v", err)
	}
}

func TestSetupWithSQLiteProvider(t *testing.T) {
	tempDir := t.TempDir()
	configPath := tempDir + "/test_sqlite.yaml"
	dbPath := tempDir + "/prompts.db"

	configContent := `prompter:
  transport:
    type: "stdio"
  storage:
    provider: "sqlite"
    sqlite:
      path: "` + dbPath + `"
      timeout: 60
      foreign_keys: true`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config, err := New(configPath)
	if err != nil {
		t.Fatalf("Setup with sqlite provider failed: %v", err)
	}

	if config.Storage.Provider != "sqlite" {
		t.Errorf("Expected storage provider 'sqlite', got '%s'", config.Storage.Provider)
	}
	if config.Storage.SQLite.Path != dbPath {
		t.Errorf("Expected sqlite path '%s', got '%s'", dbPath, config.Storage.SQLite.Path)
	}
	if config.Storage.SQLite.Timeout != 60 {
		t.Errorf("Expected timeout 60, got %d", config.Storage.SQLite.Timeout)
	}
	if !config.Storage.SQLite.ForeignKeys {
		t.Error("Expected foreign_keys to be true")
	}
}

func TestSetupRejectsInvalidProvider(t *testing.T) {
	tempDir := t.TempDir()
	configPath := tempDir + "/test_invalid_provider.yaml"

	configContent := `prompter:
  transport:
    type: "stdio"
  storage:
    provider: "invalid"`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := New(configPath)
	if err == nil {
		t.Fatal("Expected error for invalid storage provider")
	}
	if !strings.Contains(err.Error(), "invalid storage provider") {
		t.Fatalf("Expected invalid provider error, got: %v", err)
	}
}

func TestSetupSQLiteRequiresPath(t *testing.T) {
	tempDir := t.TempDir()
	configPath := tempDir + "/test_sqlite_no_path.yaml"

	// Explicitly set empty path to trigger validation
	configContent := `prompter:
  transport:
    type: "stdio"
  storage:
    provider: "sqlite"
    sqlite:
      path: ""`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := New(configPath)
	if err == nil {
		t.Fatal("Expected error for sqlite provider with empty path")
	}
	if !strings.Contains(err.Error(), "storage.sqlite.path") {
		t.Fatalf("Expected sqlite path error, got: %v", err)
	}
}

func TestEnvVarOverrides(t *testing.T) {
	tempDir := t.TempDir()
	configPath := tempDir + "/test_env.yaml"

	configContent := `prompter:
  transport:
    type: "stdio"`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set env vars
	os.Setenv("PROMPTER_STORAGE_PROVIDER", "sqlite")
	os.Setenv("PROMPTER_SQLITE_PATH", tempDir+"/env.db")
	os.Setenv("PROMPTER_SQLITE_TIMEOUT", "45")
	defer func() {
		os.Unsetenv("PROMPTER_STORAGE_PROVIDER")
		os.Unsetenv("PROMPTER_SQLITE_PATH")
		os.Unsetenv("PROMPTER_SQLITE_TIMEOUT")
	}()

	config, err := New(configPath)
	if err != nil {
		t.Fatalf("Setup with env overrides failed: %v", err)
	}

	if config.Storage.Provider != "sqlite" {
		t.Errorf("Expected provider 'sqlite' from env, got '%s'", config.Storage.Provider)
	}
	if config.Storage.SQLite.Path != tempDir+"/env.db" {
		t.Errorf("Expected sqlite path from env, got '%s'", config.Storage.SQLite.Path)
	}
	if config.Storage.SQLite.Timeout != 45 {
		t.Errorf("Expected timeout 45 from env, got %d", config.Storage.SQLite.Timeout)
	}
}

func TestDefaultSQLiteConfig(t *testing.T) {
	config := GetDefault()

	if config.Storage.SQLite.Path == "" {
		t.Error("Expected default SQLite path to be set")
	}
	if config.Storage.SQLite.Timeout != 30 {
		t.Errorf("Expected default SQLite timeout 30, got %d", config.Storage.SQLite.Timeout)
	}
	if !config.Storage.SQLite.ForeignKeys {
		t.Error("Expected default foreign_keys to be true")
	}
	if config.Storage.SQLite.MaxOpenConns != 10 {
		t.Errorf("Expected default max_open_conns 10, got %d", config.Storage.SQLite.MaxOpenConns)
	}
	if config.Storage.SQLite.MaxIdleConns != 5 {
		t.Errorf("Expected default max_idle_conns 5, got %d", config.Storage.SQLite.MaxIdleConns)
	}
}
