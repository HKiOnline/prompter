package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hkionline/prompter/internal/promptsdb"
)

func GetDefault() Configuration {

	homeDir, err := os.UserHomeDir()

	defaultPrompterDir := "/.config/prompter"
	defaultPrompterLogFile := "/prompter.log"
	defaultPromptsDir := "/prompts"

	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration defaults failure: %s", err)
		os.Exit(-1)
	}

	promptsDir := filepath.Join(homeDir, defaultPrompterDir, defaultPromptsDir)
	logFile := filepath.Join(homeDir, defaultPrompterDir, defaultPrompterLogFile)

	return Configuration{
		Transport: "stdio",
		LogFile:   logFile,
		HTTP: HTTPConfiguration{
			Port: 8080,
		},
		Storage: promptsdb.ProviderConfiguration{
			Provider: "filesystem",
			Filesystem: promptsdb.FsProviderConfiguration{
				Directory: promptsDir,
			},
		},
	}
}
