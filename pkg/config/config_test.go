package config

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestReadFromYaml(t *testing.T) {
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
			name:        "empty config file",
			path:        emptyPath,
			expected:    Config{},
			expectedErr: io.EOF,
		},
		{
			name:        "valid config file",
			path:        validPath,
			expected:    validConfig,
			expectedErr: nil,
		},
		{
			name:        "missing config file",
			path:        missingPath,
			expected:    Config{},
			expectedErr: ErrMissingConfig,
		},
		{
			name:        "invalid config file",
			path:        createInvalidConfigFile(t),
			expected:    Config{},
			expectedErr: ErrInvalidConfig,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := ReadFromYaml(test.path)
			require.ErrorIs(t, err, test.expectedErr)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestEnsureServerDefaults(t *testing.T) {
	expected := defaultServer
	actual := ensureServerDefaults(Server{})
	require.Equal(t, expected, actual)
}

func TestValidateAddress(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		expectedErr error
	}{
		{
			name:        "valid",
			address:     ":8080",
			expectedErr: nil,
		},
		{
			name:        "invalid",
			address:     "invalid",
			expectedErr: ErrInvalidConfig,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateAddress(test.address)
			require.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestValidatePattern(t *testing.T) {
	tests := []struct {
		name      string
		pattern   string
		expectErr bool
	}{
		{
			name:      "valid",
			pattern:   "/service1/",
			expectErr: false,
		},
		{
			name:      "invalid prefix",
			pattern:   "invalid/",
			expectErr: true,
		},
		{
			name:      "invalid suffix",
			pattern:   "/invalid",
			expectErr: true,
		},
		{
			name:      "invalid length",
			pattern:   "//",
			expectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validatePattern(test.pattern)
			require.Equal(t, test.expectErr, err != nil)
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expectErr bool
	}{
		{
			name:      "valid",
			url:       "http://localhost:8081",
			expectErr: false,
		},
		{
			name:      "invalid",
			url:       "://localhost:8081",
			expectErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateURL(test.url)
			require.Equal(t, test.expectErr, err != nil)
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
				Pattern: "/service1/",
				URL:     "http://localhost:8081",
			},
			expectedErr: nil,
		},
		{
			name: "missing pattern",
			upstream: Upstream{
				URL: "http://localhost:8081",
			},
			expectedErr: ErrMissingConfig,
		},
		{
			name: "invalid pattern",
			upstream: Upstream{
				Pattern: "invalid",
				URL:     "http://localhost:8081",
			},
			expectedErr: ErrInvalidConfig,
		},
		{
			name: "missing url",
			upstream: Upstream{
				Pattern: "/service1/",
			},
			expectedErr: ErrMissingConfig,
		},
		{
			name: "invalid url",
			upstream: Upstream{
				Pattern: "/service1/",
				URL:     "://localhost:8081",
			},
			expectedErr: ErrInvalidConfig,
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
					Pattern: "/service1/",
					URL:     "http://localhost:8081",
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
			name: "invalid upstreams",
			upstreams: []Upstream{
				{
					Pattern: "invalid",
					URL:     "://localhost:8081",
				},
			},
			expectedErr: ErrInvalidConfig,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateUpstreams(test.upstreams)
			require.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		server      Server
		upstreams   []Upstream
		expectedErr error
	}{
		{
			name: "valid",
			server: Server{
				Address: ":8080",
			},
			upstreams: []Upstream{
				{
					Pattern: "/service1/",
					URL:     "http://localhost:8081",
				},
			},
			expectedErr: nil,
		},
		{
			name: "invalid server",
			server: Server{
				Address: "invalid",
			},
			upstreams: []Upstream{
				{
					Pattern: "/service1/",
					URL:     "http://localhost:8081",
				},
			},
			expectedErr: ErrInvalidConfig,
		},
		{
			name: "missing upstreams",
			server: Server{
				Address: ":8080",
			},
			upstreams:   nil,
			expectedErr: ErrMissingConfig,
		},
		{
			name: "invalid upstreams",
			server: Server{
				Address: ":8080",
			},
			upstreams: []Upstream{
				{
					Pattern: "invalid",
					URL:     "://localhost:8081",
				},
			},
			expectedErr: ErrInvalidConfig,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validate(test.server, test.upstreams)
			require.ErrorIs(t, err, test.expectedErr)
		})
	}
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
				Pattern: "/service1/",
				URL:     "http://localhost:8081",
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
				URL: "http://localhost:8081",
			},
			{
				Pattern: "/service1/",
			},
		},
	}
	path := createConfigFile(t, config)
	return path
}

func createInvalidConfigFile(t *testing.T) string {
	config := Config{
		Server: Server{
			Address: "invalid",
		},
		Upstreams: []Upstream{
			{
				Pattern: "invalid",
				URL:     ":://localhost:8081",
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
