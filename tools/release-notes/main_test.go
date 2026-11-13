package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestPrepare(t *testing.T) {
	metadata := `{"major":0,"minor":8,"patch":0,"preRelease":"%s","build":"","date":"2026-11-07T13:21:07Z"}`
	for _, tc := range []struct {
		name, pre, tag, changelog, wantError string
	}{
		{"stable", "", "v0.8.0", "## [0.8.0] - 2026-11-13\n### Added\n- Current feature\n\n## [0.7.0] - 2026-11-07\n- Old feature", ""},
		{"preview", "", "", "## [0.8.0] - 2026-11-13\n- Current feature", ""},
		{"prerelease", "rc.1", "v0.8.0-rc.1", "## [0.8.0-rc.1] - 2026-11-13\n- Current feature", ""},
		{"wrong tag", "", "v0.7.0", "## [0.8.0] - 2026-11-13\n- Current feature", "tag"},
		{"missing prefix", "", "0.8.0", "## [0.8.0] - 2026-11-13\n- Current feature", "tag"},
		{"prerelease mismatch", "", "v0.8.0-rc.1", "## [0.8.0] - 2026-11-13\n- Current feature", "tag"},
		{"stale latest", "", "v0.8.0", "## [0.9.0] - 2026-11-13\n- Future\n## [0.8.0] - 2026-11-07\n- Old", "latest changelog"},
		{"unreleased", "", "v0.8.0", "## [Unreleased]\n## [0.8.0] - 2026-11-13\n- Feature", "Unreleased"},
		{"missing date", "", "v0.8.0", "## [0.8.0]\n- Feature", "heading"},
		{"invalid date", "", "v0.8.0", "## [0.8.0] - 2026-02-30\n- Feature", "date"},
		{"empty section", "", "v0.8.0", "## [0.8.0] - 2026-11-13\n\n## [0.7.0] - 2026-11-07\n- Old", "empty"},
		{"no heading", "", "v0.8.0", "# Changelog", "no release"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v, notes, err := prepare([]byte(fmt.Sprintf(metadata, tc.pre)), tc.changelog, tc.tag)
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("expected error containing %q, got %v", tc.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(v, "0.8.0") || !strings.Contains(notes, "Current feature") || strings.Contains(notes, "Old feature") {
				t.Fatalf("wrong release section: %s", notes)
			}
			if tc.pre != "" && strings.Contains(notes, "brew install") {
				t.Fatal("prerelease must not tell users to install the stable cask")
			}
		})
	}
}

func TestPrepareInvalidMetadata(t *testing.T) {
	if _, _, err := prepare([]byte(`{"major":0}`), "", "v0.8.0"); err == nil {
		t.Fatal("invalid version metadata accepted")
	}
}
