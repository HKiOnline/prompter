// Command release-notes validates release metadata and prepares the curated
// GitHub release body. It never creates a tag or publishes anything.
package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hkionline/prompter/internal/version"
)

func main() {
	tag := flag.String("tag", "", "release tag to validate (required in release CI)")
	output := flag.String("output", "", "write Markdown to this file instead of stdout")
	versionOnly := flag.Bool("version", false, "print the validated committed version instead of release notes")
	flag.Parse()
	if flag.NArg() != 0 {
		fail(fmt.Errorf("unexpected positional arguments"))
	}
	metadata, err := os.ReadFile("version.json")
	if err != nil {
		fail(err)
	}
	changelog, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		fail(err)
	}
	v, notes, err := prepare(metadata, string(changelog), *tag)
	if err != nil {
		fail(err)
	}
	if *versionOnly {
		fmt.Println(v)
	} else if *output == "" {
		fmt.Print(notes)
	} else if err := os.WriteFile(*output, []byte(notes), 0o644); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "release preparation:", err)
	os.Exit(1)
}

var heading = regexp.MustCompile(`^## \[([^\]]+)\] - (\d{4}-\d{2}-\d{2})$`)

func prepare(data []byte, changelog, tag string) (string, string, error) {
	m, err := version.Load(data)
	if err != nil {
		return "", "", err
	}
	if tag != "" && tag != "v"+m.SemVer() {
		return "", "", fmt.Errorf("tag %q does not match version.json: expected v%s", tag, m.SemVer())
	}
	lines := strings.Split(strings.ReplaceAll(changelog, "\r\n", "\n"), "\n")
	start, end := -1, len(lines)
	for i, line := range lines {
		if !strings.HasPrefix(line, "## ") {
			continue
		}
		if strings.Contains(strings.ToLower(line), "unreleased") {
			return "", "", fmt.Errorf("CHANGELOG.md must not contain an Unreleased section")
		}
		if start >= 0 {
			if end == len(lines) {
				end = i
			}
			continue
		}
		match := heading.FindStringSubmatch(line)
		if match == nil {
			return "", "", fmt.Errorf("invalid latest changelog heading: %s", line)
		}
		if match[1] != m.Core() {
			return "", "", fmt.Errorf("latest changelog version %q does not match version.json %q", match[1], m.Core())
		}
		if _, err := time.Parse("2006-01-02", match[2]); err != nil {
			return "", "", fmt.Errorf("invalid changelog release date: %w", err)
		}
		start = i + 1
	}
	if start < 0 {
		return "", "", fmt.Errorf("no release heading in CHANGELOG.md")
	}
	changes := strings.TrimSpace(strings.Join(lines[start:end], "\n"))
	if changes == "" {
		return "", "", fmt.Errorf("latest changelog section is empty")
	}
	installation := "## Install on macOS (recommended)\n\n" +
		"```sh\nbrew install --cask hkionline/tap/prompter\nprompter --version\n```\n\n" +
		"For an existing installation, run `brew update` then `brew upgrade --cask hkionline/tap/prompter`. " +
		"Homebrew selects the native Apple Silicon or Intel binary. The unsigned cask removes quarantine from the installed executable.\n\n"
	if m.PreRelease != "" {
		installation = "## Install this prerelease\n\nDownload the archive for your platform below. " +
			"This prerelease does not update the stable Homebrew cask.\n\n"
	}
	notes := installation + "## Other downloads and setup\n\n" +
		"Release assets contain macOS/Linux `.tar.gz` archives and Windows `.zip` archives. " +
		"Choose `arm64` for ARM64, `amd64` for 64-bit Intel/AMD, or `386` for 32-bit x86 (Linux/Windows only). " +
		"Extract `prompter` (or `prompter.exe`) and place it on your PATH. " +
		"Use the attached SHA-256 checksums to verify downloads. GitHub's source archives are not prebuilt binaries.\n\n" +
		"SQLite support is included; filesystem storage remains the default. Configure your MCP client with the executable's absolute path " +
		"(`command -v prompter` on macOS/Linux). See the [client configuration guide](https://github.com/hkionline/prompter/blob/v" + m.SemVer() + "/docs/client-configuration.md).\n\n" +
		"## Changes\n\n" + changes + "\n"
	return m.SemVer(), notes, nil
}
