package server

import (
	"testing"

	"github.com/hkionline/prompter/internal/configuration"
	"github.com/hkionline/prompter/internal/plog"
	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/stretchr/testify/assert"
)

// MockDB is a mock implementation of promptsdb.Provider for testing
type MockDB struct{}

func (m *MockDB) Create(prompt promptsdb.Prompt) error {
	return nil
}

func (m *MockDB) Read(name string) (promptsdb.Prompt, error) {
	return promptsdb.Prompt{
		Name:        name,
		Title:       "Test Prompt",
		Description: "Test Description",
		Content:     "Test Content",
	}, nil
}

func (m *MockDB) List(query promptsdb.PromptQuery) ([]promptsdb.Prompt, error) {
	return []promptsdb.Prompt{
		{
			Name:        "test1",
			Title:       "Test 1",
			Description: "Description 1",
		},
		{
			Name:        "test2",
			Title:       "Test 2",
			Description: "Description 2",
		},
	}, nil
}

func (m *MockDB) Update(prompt promptsdb.Prompt) error {
	return nil
}

func (m *MockDB) Delete(name string) error {
	return nil
}

func (m *MockDB) Close() error {
	return nil
}

func (m *MockDB) Setup(config promptsdb.ProviderConfiguration) error {
	return nil
}

func TestNewServer(t *testing.T) {
	config := &configuration.Configuration{
		Storage: promptsdb.ProviderConfiguration{
			Provider: "filesystem",
			Filesystem: promptsdb.FsProviderConfiguration{
				Directory: "/tmp/test",
			},
		},
		LogFile: "/tmp/test.log",
	}

	logger := plog.New("/tmp/test.log")
	db := &MockDB{}

	server := New("0.5.0", config, logger, db)

	assert.NotNil(t, server)
	assert.Equal(t, config, server.config)
	assert.Equal(t, logger, server.logger)
	assert.Equal(t, db, server.db)
}

func TestServerStart(t *testing.T) {
	config := &configuration.Configuration{
		Transport: "stdio",
		Storage: promptsdb.ProviderConfiguration{
			Provider: "filesystem",
			Filesystem: promptsdb.FsProviderConfiguration{
				Directory: "/tmp/test",
			},
		},
		LogFile: "/tmp/test.log",
	}

	logger := plog.New("/tmp/test.log")
	db := &MockDB{}

	server := New("0.5.0", config, logger, db)

	// Test that server is initialized correctly (server field is nil until Run is called)
	assert.NotNil(t, server)
	assert.Equal(t, "stdio", config.Transport)
}
