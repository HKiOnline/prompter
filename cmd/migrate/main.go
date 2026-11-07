//go:build sqlite

// Command migrate migrates prompts from the filesystem provider to a SQLite database.
//
// Usage:
//
//	prompter-migrate -src /path/to/prompts -dst /path/to/prompts.db [flags]
//
// Flags:
//
//	-src        Source directory containing prompt .md files (required)
//	-dst        Destination SQLite database file path (required)
//	-dry-run    Validate without writing to the database
//	-incremental Only migrate prompts not already present in the target database
//	-verbose    Print detailed progress information
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hkionline/prompter/internal/promptsdb"
	"github.com/hkionline/prompter/internal/promptsdb/sqlite"
	"gopkg.in/yaml.v3"
)

type migrationReport struct {
	Total     int
	Migrated  int
	Skipped   int
	Failed    int
	Conflicts int
	Warnings  []string
	Errors    []string
}

func main() {
	srcDir := flag.String("src", "", "Source directory containing prompt .md files (required)")
	dstPath := flag.String("dst", "", "Destination SQLite database file path (required)")
	dryRun := flag.Bool("dry-run", false, "Validate without writing to the database")
	incremental := flag.Bool("incremental", false, "Only migrate prompts not already in the target database")
	verbose := flag.Bool("verbose", false, "Print detailed progress")
	flag.Parse()

	if *srcDir == "" || *dstPath == "" {
		fmt.Fprintf(os.Stderr, "Usage: prompter-migrate -src <prompts-dir> -dst <db-file> [--dry-run] [--incremental] [--verbose]\n")
		os.Exit(1)
	}

	report, err := runMigration(*srcDir, *dstPath, *dryRun, *incremental, *verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Migration failed: %s\n", err)
		os.Exit(1)
	}

	printReport(report, *dryRun)

	if report.Failed > 0 {
		os.Exit(1)
	}
}

func runMigration(srcDir, dstPath string, dryRun, incremental, verbose bool) (*migrationReport, error) {
	report := &migrationReport{}

	// Read all prompt files from source directory
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read source directory %s: %w", srcDir, err)
	}

	var prompts []promptsdb.Prompt
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(srcDir, entry.Name())
		prompt, err := loadPromptFile(filePath)
		if err != nil {
			report.Failed++
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", entry.Name(), err.Error()))
			if verbose {
				fmt.Fprintf(os.Stderr, "SKIP (corrupted): %s - %s\n", entry.Name(), err.Error())
			}
			continue
		}

		prompts = append(prompts, prompt)
		report.Total++
	}

	if verbose {
		fmt.Printf("Found %d prompt file(s) in %s\n", report.Total, srcDir)
	}

	if dryRun {
		if verbose {
			for _, p := range prompts {
				fmt.Printf("  [dry-run] Would migrate: %s (%s)\n", p.Name, p.Title)
			}
		}
		report.Migrated = len(prompts)
		fmt.Printf("Dry run complete: %d prompt(s) would be migrated to %s\n", len(prompts), dstPath)
		return report, nil
	}

	// Create SQLite provider for the destination
	provider, err := sqlite.NewSQLiteProvider(sqlite.SQLiteProviderConfiguration{
		Path:        dstPath,
		Timeout:     30,
		ForeignKeys: true,
	}, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create SQLite provider: %w", err)
	}

	sqliteProvider, ok := provider.(*sqlite.SQLiteProvider)
	if !ok {
		return nil, fmt.Errorf("unexpected provider type: %T", provider)
	}
	defer sqliteProvider.Close()

	// Migrate each prompt
	for i, prompt := range prompts {
		if verbose {
			fmt.Printf("[%d/%d] Migrating: %s\n", i+1, len(prompts), prompt.Name)
		}

		// Check for conflicts in incremental mode
		if incremental {
			_, err := sqliteProvider.Read(prompt.Name)
			if err == nil {
				report.Conflicts++
				report.Warnings = append(report.Warnings, fmt.Sprintf("conflict skipped: %s already exists", prompt.Name))
				if verbose {
					fmt.Printf("  SKIP (conflict): %s already exists\n", prompt.Name)
				}
				continue
			}
		}

		// Try to create; if it already exists, try update
		err := sqliteProvider.Create(prompt)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				report.Conflicts++
				report.Warnings = append(report.Warnings, fmt.Sprintf("conflict skipped: %s already exists", prompt.Name))
				if verbose {
					fmt.Printf("  SKIP (conflict): %s\n", prompt.Name)
				}
				continue
			}
			report.Failed++
			report.Errors = append(report.Errors, fmt.Sprintf("%s: %s", prompt.Name, err.Error()))
			if verbose {
				fmt.Fprintf(os.Stderr, "  FAIL: %s - %s\n", prompt.Name, err.Error())
			}
			continue
		}

		report.Migrated++
	}

	report.Skipped = report.Total - report.Migrated - report.Conflicts
	if report.Skipped < 0 {
		report.Skipped = 0
	}

	return report, nil
}

func loadPromptFile(filePath string) (promptsdb.Prompt, error) {
	var prompt promptsdb.Prompt

	data, err := os.ReadFile(filePath)
	if err != nil {
		return prompt, err
	}

	re := regexp.MustCompile(`(?s)^---\s*\n(.*?)\n---\s*\n?(.*)$`)
	match := re.FindStringSubmatch(string(data))
	if len(match) < 3 {
		return prompt, fmt.Errorf("failed to parse prompt file %s: missing frontmatter", filePath)
	}

	if err := yaml.Unmarshal([]byte(match[1]), &prompt); err != nil {
		return prompt, fmt.Errorf("failed to parse YAML frontmatter: %w", err)
	}

	prompt.Id = prompt.Name
	prompt.Content = strings.TrimSpace(match[2])

	if prompt.Name == "" {
		return prompt, fmt.Errorf("prompt name is empty")
	}
	if prompt.Content == "" {
		return prompt, fmt.Errorf("prompt content is empty")
	}

	return prompt, nil
}

func printReport(r *migrationReport, dryRun bool) {
	mode := "Migration"
	if dryRun {
		mode = "Dry run"
	}
	fmt.Printf("\n=== %s Report ===\n", mode)
	fmt.Printf("Total files:    %d\n", r.Total)
	fmt.Printf("Migrated:       %d\n", r.Migrated)
	fmt.Printf("Skipped:        %d\n", r.Skipped)
	fmt.Printf("Conflicts:      %d\n", r.Conflicts)
	fmt.Printf("Failed:         %d\n", r.Failed)

	if len(r.Warnings) > 0 {
		fmt.Printf("\nWarnings:\n")
		for _, w := range r.Warnings {
			fmt.Printf("  - %s\n", w)
		}
	}

	if len(r.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, e := range r.Errors {
			fmt.Printf("  - %s\n", e)
		}
	}
}
