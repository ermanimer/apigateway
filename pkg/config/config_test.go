package config

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestReadConfigFromYaml(t *testing.T) {
	validPath, validConfig := createValidConfigFile(t)
	missingPath := createMissingConfigFile(t)
	emptyPath := createEmptyFile(t)

	tests := []struct {
		name        string
		path        string
		expected    Config
		expectedErr error
	}{
		{
			name:        "non-existing path",
			path:        "non-existing.yaml",
			expected:    Config{},
			expectedErr: os.ErrNotExist,
		},
		{
			name:        "empty config",
			path:        emptyPath,
			expected:    Config{},
			expectedErr: io.EOF,
		},
		{
			name:        "valid config",
			path:        validPath,
			expected:    validConfig,
			expectedErr: nil,
		},
		{
			name:        "missing config",
			path:        missingPath,
			expected:    Config{},
			expectedErr: ErrMissingConfig,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := ReadConfigFromYaml(test.path)
			require.ErrorIs(t, err, test.expectedErr)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestValidateUpstream(t *testing.T) {
	tests := []struct {
		name        string
		upstream    Upstream
		expectedErr error
	}{
		{
			name: "valid",
			upstream: Upstream{
				Pattern: "/api/",
				URL:     "http://localhost:8080",
			},
			expectedErr: nil,
		},
		{
			name: "missing pattern",
			upstream: Upstream{
				URL: "http://localhost:8080",
			},
			expectedErr: ErrMissingConfig,
		},
		{
			name: "missing url",
			upstream: Upstream{
				Pattern: "/api/",
			},
			expectedErr: ErrMissingConfig,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateUpstream(0, test.upstream)
			require.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestValidateUpstreams(t *testing.T) {
	tests := []struct {
		name        string
		upstreams   []Upstream
		expectedErr error
	}{
		{
			name: "valid",
			upstreams: []Upstream{
				{
					Pattern: "/api/",
					URL:     "http://localhost:8080",
				},
			},
			expectedErr: nil,
		},
		{
			name:        "missing upstreams",
			upstreams:   nil,
			expectedErr: ErrMissingConfig,
		},
		{
			name: "missing pattern and url",
			upstreams: []Upstream{
				{
					URL: "http://localhost:8080",
				},
				{
					Pattern: "/api/",
				},
			},
			expectedErr: ErrMissingConfig,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateUpstreams(test.upstreams)
			require.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestEnsureServerDefaults(t *testing.T) {
	expected := defaultServer
	actual := ensureServerDefaults(Server{})
	require.Equal(t, expected, actual)
}

func createConfigFile(t *testing.T, config Config) string {
	path := path.Join(t.TempDir(), "config.yaml")
	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()
	err = yaml.NewEncoder(file).Encode(config)
	require.NoError(t, err)
	return path
}

func createValidConfigFile(t *testing.T) (string, Config) {
	config := Config{
		Server: defaultServer,
		Upstreams: []Upstream{
			{
				Pattern: "/api/",
				URL:     "http://localhost:8080",
			},
		},
	}
	path := createConfigFile(t, config)
	return path, config
}

func createMissingConfigFile(t *testing.T) string {
	config := Config{
		Server: defaultServer,
		Upstreams: []Upstream{
			{
				URL: "http://localhost:8080",
			},
			{
				Pattern: "/api/",
			},
		},
	}
	path := createConfigFile(t, config)
	return path
}

func createEmptyFile(t *testing.T) string {
	path := path.Join(t.TempDir(), "empty")
	file, err := os.Create(path)
	require.NoError(t, err)
	defer file.Close()
	return path
}
