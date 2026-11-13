package version

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAndCompose(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		semver  string
		core    string
		wantErr string
	}{
		{
			name:   "stable version",
			json:   `{"major":1,"minor":2,"patch":3,"preRelease":"","build":"","date":"2026-11-07T13:21:07Z"}`,
			semver: "1.2.3",
			core:   "1.2.3",
		},
		{
			name:   "zero components",
			json:   `{"major":0,"minor":0,"patch":0,"preRelease":"","build":"","date":"2026-11-07T13:21:07Z"}`,
			semver: "0.0.0",
			core:   "0.0.0",
		},
		{
			name:   "pre-release only",
			json:   `{"major":1,"minor":2,"patch":3,"preRelease":"rc.1","build":"","date":"2026-11-07T13:21:07Z"}`,
			semver: "1.2.3-rc.1",
			core:   "1.2.3-rc.1",
		},
		{
			name:   "build only",
			json:   `{"major":1,"minor":2,"patch":3,"preRelease":"","build":"abc.001","date":"2026-11-07T13:21:07Z"}`,
			semver: "1.2.3+abc.001",
			core:   "1.2.3",
		},
		{
			name:   "combined suffixes",
			json:   `{"major":1,"minor":2,"patch":3,"preRelease":"rc.1","build":"abc.001","date":"2026-11-07T13:21:07Z"}`,
			semver: "1.2.3-rc.1+abc.001",
			core:   "1.2.3-rc.1",
		},
		{
			name:   "local build stamp",
			json:   `{"major":0,"minor":8,"patch":0,"preRelease":"","build":"lb.20261107130122","date":"2026-11-07T13:01:22Z"}`,
			semver: "0.8.0+lb.20261107130122",
			core:   "0.8.0",
		},
		{
			name:   "ci build stamp",
			json:   `{"major":0,"minor":8,"patch":0,"preRelease":"","build":"ci.12345678.20261107130122","date":"2026-11-07T13:01:22Z"}`,
			semver: "0.8.0+ci.12345678.20261107130122",
			core:   "0.8.0",
		},
		{
			name:    "malformed JSON",
			json:    `{"major":1,`,
			wantErr: "invalid embedded version.json",
		},
		{
			name:    "missing field",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"","build":""}`,
			wantErr: "missing or null required field(s): date",
		},
		{
			name:    "null field",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":null,"build":"","date":"2026-11-07T13:21:07Z"}`,
			wantErr: "missing or null required field(s): preRelease",
		},
		{
			name:    "wrong type string for major",
			json:    `{"major":"1","minor":2,"patch":3,"preRelease":"","build":"","date":"2026-11-07T13:21:07Z"}`,
			wantErr: "invalid embedded version.json",
		},
		{
			name:    "fractional component",
			json:    `{"major":1.5,"minor":2,"patch":3,"preRelease":"","build":"","date":"2026-11-07T13:21:07Z"}`,
			wantErr: "invalid embedded version.json",
		},
		{
			name:    "negative component",
			json:    `{"major":-1,"minor":2,"patch":3,"preRelease":"","build":"","date":"2026-11-07T13:21:07Z"}`,
			wantErr: "nonnegative",
		},
		{
			name:    "pre-release numeric leading zero",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"rc.01","build":"","date":"2026-11-07T13:21:07Z"}`,
			wantErr: "leading zeroes",
		},
		{
			name:   "pre-release numeric zero allowed",
			json:   `{"major":1,"minor":2,"patch":3,"preRelease":"0","build":"","date":"2026-11-07T13:21:07Z"}`,
			semver: "1.2.3-0",
			core:   "1.2.3-0",
		},
		{
			name:    "invalid pre-release identifier chars",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"rc_1","build":"","date":"2026-11-07T13:21:07Z"}`,
			wantErr: "invalid preRelease",
		},
		{
			name:    "empty pre-release identifier",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"rc..1","build":"","date":"2026-11-07T13:21:07Z"}`,
			wantErr: "invalid preRelease",
		},
		{
			name:   "build numeric leading zero allowed",
			json:   `{"major":1,"minor":2,"patch":3,"preRelease":"","build":"001","date":"2026-11-07T13:21:07Z"}`,
			semver: "1.2.3+001",
			core:   "1.2.3",
		},
		{
			name:    "invalid build identifier chars",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"","build":"ab_c","date":"2026-11-07T13:21:07Z"}`,
			wantErr: "invalid build",
		},
		{
			name:    "date not UTC",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"","build":"","date":"2026-11-07T13:21:07+02:00"}`,
			wantErr: "invalid date",
		},
		{
			name:    "date with fractional seconds",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"","build":"","date":"2026-11-07T13:21:07.000Z"}`,
			wantErr: "invalid date",
		},
		{
			name:    "date only",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"","build":"","date":"2026-11-07"}`,
			wantErr: "invalid date",
		},
		{
			name:    "impossible date",
			json:    `{"major":1,"minor":2,"patch":3,"preRelease":"","build":"","date":"2026-13-40T99:99:99Z"}`,
			wantErr: "invalid date",
		},
	}

	BuildDate, BuildID = "", ""
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Load([]byte(tt.json))
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.semver, m.SemVer())
			assert.Equal(t, tt.core, m.Core())
		})
	}
}

func TestLoadAppliesStampOverrides(t *testing.T) {
	embedded := []byte(`{"major":0,"minor":8,"patch":0,"preRelease":"","build":"","date":"2026-11-07T13:21:07Z"}`)

	oldDate, oldID := BuildDate, BuildID
	defer func() { BuildDate, BuildID = oldDate, oldID }()

	BuildDate, BuildID = "", ""
	m, err := Load(embedded)
	require.NoError(t, err)
	assert.Equal(t, "0.8.0", m.SemVer())
	assert.Equal(t, "2026-11-07T13:21:07Z", m.Date)

	BuildDate, BuildID = "2026-11-07T14:00:00Z", "lb.20261107140000"
	m, err = Load(embedded)
	require.NoError(t, err)
	assert.Equal(t, "0.8.0+lb.20261107140000", m.SemVer())
	assert.Equal(t, "0.8.0", m.Core())
	assert.Equal(t, "2026-11-07T14:00:00Z", m.Date)
	assert.Equal(t, "prompter 0.8.0+lb.20261107140000 (2026-11-07T14:00:00Z)", m.Line())

	BuildDate, BuildID = "2026-11-07T14:00:00Z", "ci.12345678.20261107140000"
	m, err = Load(embedded)
	require.NoError(t, err)
	assert.Equal(t, "0.8.0+ci.12345678.20261107140000", m.SemVer())
	assert.Equal(t, "0.8.0", m.Core())

	BuildDate, BuildID = "not-a-date", ""
	_, err = Load(embedded)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid build date stamp")

	BuildDate, BuildID = "", "bad_build!"
	_, err = Load(embedded)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid build stamp")

	_, err = Load([]byte(`{invalid`))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid embedded version.json")
}

func TestEmbeddedVersionJSONIsValid(t *testing.T) {
	data, err := readRepoFile(t, "version.json")
	require.NoError(t, err)
	m, err := Load(data)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(m.Line(), "prompter "), "line format")
}
