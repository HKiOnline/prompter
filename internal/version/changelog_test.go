package version

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// readRepoFile reads a file from the repository root relative to this
// package's directory (internal/version).
func readRepoFile(t *testing.T, name string) ([]byte, error) {
	t.Helper()
	return os.ReadFile("../../" + name)
}

var (
	changelogHeadingRe = regexp.MustCompile(`^## \[([^\]]+)\](?: - (\S+))?\s*$`)
	unreleasedRe       = regexp.MustCompile(`(?i)^## \[?Unreleased\]?\s*$`)
)

// TestEmbeddedVersionMatchesChangelog verifies that the committed
// version.json version matches the latest versioned CHANGELOG heading and
// that no Unreleased section exists. The date is build time, so only the
// version components are compared.
func TestEmbeddedVersionMatchesChangelog(t *testing.T) {
	data, err := readRepoFile(t, "version.json")
	require.NoError(t, err)

	BuildDate, BuildID = "", ""
	m, err := Load(data)
	require.NoError(t, err)

	changelog, err := readRepoFile(t, "CHANGELOG.md")
	require.NoError(t, err)

	// CHANGELOG.md is ordered newest first: the first versioned heading is
	// the latest release.
	latest := ""
	for _, line := range regexp.MustCompile(`\n`).Split(string(changelog), -1) {
		if unreleasedRe.MatchString(line) {
			t.Fatalf("CHANGELOG.md must not contain an Unreleased heading, found: %q", line)
		}
		if latest == "" {
			if match := changelogHeadingRe.FindStringSubmatch(line); match != nil {
				latest = match[1]
			}
		}
	}
	require.NotEmpty(t, latest, "no versioned heading found in CHANGELOG.md")
	assert.Equal(t, latest, m.Core(), "embedded version.json version must match the latest CHANGELOG heading")
}
