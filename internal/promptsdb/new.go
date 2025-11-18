package promptsdb

const (
	FILE_SYSTEM_PROVIDER = "fsProvider"
)

func New(dbProvider string, config ProviderConfiguration, logfile string) (Provider, error) {

	switch dbProvider {
	case FILE_SYSTEM_PROVIDER:
		return NewPromptsFsProvider(config.Filesystem.Directory, logfile)
	default:
		return NewPromptsFsProvider(config.Filesystem.Directory, logfile)
	}
}
