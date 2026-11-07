// Package version provides parsing, validation, and formatting of the
// embedded version metadata (root version.json).
package version

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// BuildDate and BuildID are stamped at build time via -ldflags -X (by the
// Makefile or GoReleaser). Empty values mean the embedded version.json
// values are used. BuildID uses dot-separated SemVer build identifiers:
// lb.{YYYYMMDDHHMMSS} for local builds, ci.{pipeline-id}.{YYYYMMDDHHMMSS}
// in CI when a pipeline ID is available.
var BuildDate string
var BuildID string

// Metadata is the schema of the root version.json file. All six fields are
// required; absent or null fields are rejected during decoding.
type Metadata struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	PreRelease string `json:"preRelease"`
	Build      string `json:"build"`
	Date       string `json:"date"`
}

// UnmarshalJSON decodes Metadata, distinguishing absent/null fields from
// legitimate zero/empty values: all six fields are required.
func (m *Metadata) UnmarshalJSON(data []byte) error {
	var raw struct {
		Major      *int    `json:"major"`
		Minor      *int    `json:"minor"`
		Patch      *int    `json:"patch"`
		PreRelease *string `json:"preRelease"`
		Build      *string `json:"build"`
		Date       *string `json:"date"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	missing := []string{}
	if raw.Major == nil {
		missing = append(missing, "major")
	}
	if raw.Minor == nil {
		missing = append(missing, "minor")
	}
	if raw.Patch == nil {
		missing = append(missing, "patch")
	}
	if raw.PreRelease == nil {
		missing = append(missing, "preRelease")
	}
	if raw.Build == nil {
		missing = append(missing, "build")
	}
	if raw.Date == nil {
		missing = append(missing, "date")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing or null required field(s): %s", strings.Join(missing, ", "))
	}
	m.Major, m.Minor, m.Patch = *raw.Major, *raw.Minor, *raw.Patch
	m.PreRelease, m.Build, m.Date = *raw.PreRelease, *raw.Build, *raw.Date
	return nil
}

var (
	identifierRe = regexp.MustCompile(`^[0-9A-Za-z-]+$`)
	numericRe    = regexp.MustCompile(`^[0-9]+$`)
	// numeric pre-release identifiers must not have leading zeroes
	preNumericRe = regexp.MustCompile(`^(0|[1-9][0-9]*)$`)
	// date must be exactly YYYY-MM-DDTHH:mm:ssZ (ISO 8601, UTC)
	dateFormatRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)
)

// Validate checks the metadata against the SemVer 2.0.0 rules and the
// required date format.
func (m Metadata) Validate() error {
	if m.Major < 0 || m.Minor < 0 || m.Patch < 0 {
		return fmt.Errorf("version components must be nonnegative integers, got %d.%d.%d", m.Major, m.Minor, m.Patch)
	}
	if err := validatePreRelease(m.PreRelease); err != nil {
		return fmt.Errorf("invalid preRelease: %w", err)
	}
	if err := validateBuild(m.Build); err != nil {
		return fmt.Errorf("invalid build: %w", err)
	}
	if err := validateDate(m.Date); err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}
	return nil
}

func validatePreRelease(s string) error {
	if s == "" {
		return nil
	}
	for _, id := range strings.Split(s, ".") {
		if !identifierRe.MatchString(id) {
			return fmt.Errorf("identifier %q must be nonempty ASCII alphanumeric/hyphen", id)
		}
		if numericRe.MatchString(id) && !preNumericRe.MatchString(id) {
			return fmt.Errorf("numeric identifier %q must not have leading zeroes", id)
		}
	}
	return nil
}

func validateBuild(s string) error {
	if s == "" {
		return nil
	}
	for _, id := range strings.Split(s, ".") {
		if !identifierRe.MatchString(id) {
			return fmt.Errorf("identifier %q must be nonempty ASCII alphanumeric/hyphen", id)
		}
		// numeric build identifiers may have leading zeroes
	}
	return nil
}

func validateDate(s string) error {
	if !dateFormatRe.MatchString(s) {
		return fmt.Errorf("date %q must have format YYYY-MM-DDTHH:mm:ssZ (UTC)", s)
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("date %q is not a valid calendar date/time: %w", s, err)
	}
	if t.Location() != time.UTC {
		return fmt.Errorf("date %q must be UTC", s)
	}
	return nil
}

// SemVer composes major.minor.patch[-preRelease][+build] per SemVer 2.0.0,
// without a leading v.
func (m Metadata) SemVer() string {
	v := fmt.Sprintf("%d.%d.%d", m.Major, m.Minor, m.Patch)
	if m.PreRelease != "" {
		v += "-" + m.PreRelease
	}
	if m.Build != "" {
		v += "+" + m.Build
	}
	return v
}

// Core composes major.minor.patch[-preRelease] without build metadata.
// This is the version reported as MCP serverInfo.version.
func (m Metadata) Core() string {
	v := fmt.Sprintf("%d.%d.%d", m.Major, m.Minor, m.Patch)
	if m.PreRelease != "" {
		v += "-" + m.PreRelease
	}
	return v
}

// Line returns the CLI version output line: prompter <semver> (<date>).
func (m Metadata) Line() string {
	return fmt.Sprintf("prompter %s (%s)", m.SemVer(), m.Date)
}

// Load parses and validates the embedded metadata, then applies the
// build-time stamp overrides (BuildDate/BuildID) when they are set.
func Load(embedded []byte) (Metadata, error) {
	var m Metadata
	if err := json.Unmarshal(embedded, &m); err != nil {
		return Metadata{}, fmt.Errorf("invalid embedded version.json: %w", err)
	}
	if err := m.Validate(); err != nil {
		return Metadata{}, fmt.Errorf("invalid embedded version.json: %w", err)
	}
	if BuildDate != "" {
		if err := validateDate(BuildDate); err != nil {
			return Metadata{}, fmt.Errorf("invalid build date stamp: %w", err)
		}
		m.Date = BuildDate
	}
	if BuildID != "" {
		if err := validateBuild(BuildID); err != nil {
			return Metadata{}, fmt.Errorf("invalid build stamp: %w", err)
		}
		m.Build = BuildID
	}
	return m, nil
}
