package promptsdb

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/hkionline/prompter/internal/plog"
	"gopkg.in/yaml.v3"
)

type FsProvider struct {
	mu sync.RWMutex
	// config
	cache map[string]Prompt // map of cached prompts identified by prompt id
	files map[string]string // map of prompt files identified by prompt id
	dir   string            // directory where prompt files are stored
}

type FsProviderConfiguration struct {
	Directory string `yaml:"prompts_directory"` // directory to save prompt files to
}

func NewPromptsFsProvider(promptsDir string, logfile string) (*FsProvider, error) {

	p := plog.New(logfile)

	p.Write(plog.SERVER, "setting up new prompts filesystem provider")

	// Load prompts to cache
	cache, files, err := loadCache(promptsDir, p)

	if err != nil {
		return &FsProvider{}, err
	}

	return &FsProvider{
		cache: cache,
		files: files,
		dir:   promptsDir,
	}, nil

}

func (f *FsProvider) Create(prompt Prompt) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	prompt.Id = prompt.Name

	// Write the prompt to a file
	fileName, err := savePrompt(prompt, f.dir)

	if err != nil {
		return err
	}

	// Add the prompt file to files map
	f.files[prompt.Id] = fileName

	// Add the prompt to cache
	f.cache[prompt.Id] = prompt

	return nil
}

func (f *FsProvider) Read(promptId string) (Prompt, error) {

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Read the cache and return the prompt

	if prompt, ok := f.cache[promptId]; ok {
		return prompt, nil
	} else {
		return Prompt{}, errors.New("could not read prompt: no prompt with the given id was found")
	}
}

func (f *FsProvider) Update(prompt Prompt) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	prompt.Id = prompt.Name

	// Write the prompt to a file
	fileName, err := savePrompt(prompt, f.dir)

	if err != nil {
		return err
	}

	// Add the prompt file to files map
	f.files[prompt.Id] = fileName

	// Add the prompt to cache
	f.cache[prompt.Id] = prompt

	return nil
}

func (f *FsProvider) Delete(promptId string) error {

	f.mu.Lock()
	defer f.mu.Unlock()

	// Remove the prompt from the cache
	delete(f.cache, promptId)

	// Remove the prompt file
	if promptFile, ok := f.files[promptId]; ok {
		return removePrompt(f.dir, promptFile)
	} else {
		return errors.New("could not delete prompt: no prompt file with the given id was found")
	}

}

func (f *FsProvider) List(query PromptQuery) ([]Prompt, error) {

	f.mu.RLock()
	defer f.mu.RUnlock()

	return slices.Collect(maps.Values(f.cache)), nil
}

func loadCache(fromDir string, p *plog.Plogger) (map[string]Prompt, map[string]string, error) {

	p.Write(plog.SERVER, "loading prompts from filesystem to populate the cache")

	cache := map[string]Prompt{}
	files := map[string]string{}

	dirEntries, err := os.ReadDir(fromDir)

	if err != nil {
		return cache, files, err
	}

	for _, entry := range dirEntries {

		if !entry.IsDir() {
			prompt, err := loadPrompt(filepath.Join(fromDir, entry.Name()), p)

			if err != nil {
				p.Write(plog.SERVER, err.Error())
				continue
			}

			cache[prompt.Id] = prompt
			files[prompt.Id] = entry.Name()
		}
	}

	return cache, files, nil
}

func loadPrompt(fromFile string, p *plog.Plogger) (Prompt, error) {

	p.Write(plog.SERVER, "loading prompt file "+fromFile)

	var prompt Prompt

	// Load configuration file
	file, err := os.ReadFile(fromFile)

	if err != nil {
		return prompt, err
	}

	err = yaml.Unmarshal(file, &prompt)

	if err != nil {
		return prompt, err
	}

	prompt.Id = prompt.Name

	re, err := regexp.Compile(`(?s)---.*?---\s*(.*)`)
	match := re.FindStringSubmatch(string(file))

	if err != nil {
		return prompt, err
	}

	if len(match) > 1 {
		// Extract and trim the unstructured text
		prompt.Content = strings.TrimSpace(match[1])
	} else {
		return prompt, fmt.Errorf("failed to extract prompt contents, matches (%d) =< 1", len(match))
	}

	return prompt, nil
}

func savePrompt(prompt Prompt, promptDir string) (string, error) {

	path := filepath.Join(promptDir, fmt.Sprintf("%s.%s", prompt.Id, "md"))

	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return path, err
	}

	file.WriteString("---\n")

	promptBytes, err := yaml.Marshal(prompt)

	if err != nil {
		return path, err
	}

	_, err = file.Write(promptBytes)

	if err != nil {
		return path, err
	}

	file.WriteString("---\n")
	file.WriteString(prompt.Content)

	return path, nil
}

func removePrompt(promptDir string, promptFile string) error {
	return os.Remove(filepath.Join(promptDir, promptFile))
}
