package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Address         string        `yaml:"address"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes  int           `yaml:"max_header_bytes"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type Upstream struct {
	Pattern     string `yaml:"pattern"`
	StripPrefix bool   `yaml:"strip_prefix"`
	URL         string `yaml:"url"`
}

type Config struct {
	Server    Server     `yaml:"server"`
	Upstreams []Upstream `yaml:"upstreams"`
}

var defaultServer = Server{
	Address:         ":8080",
	ReadTimeout:     5 * time.Second,
	WriteTimeout:    10 * time.Second,
	IdleTimeout:     120 * time.Second,
	MaxHeaderBytes:  1048576,
	ShutdownTimeout: 10 * time.Second,
}

var (
	ErrMissingConfig = fmt.Errorf("missing config")
	ErrInvalidConfig = fmt.Errorf("invalid config")
)

func ReadFromYaml(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var c Config
	err = yaml.NewDecoder(file).Decode(&c)
	if err != nil {
		return Config{}, err
	}

	c.Server = ensureServerDefaults(c.Server)

	err = validate(c.Server, c.Upstreams)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

func ensureServerDefaults(s Server) Server {
	if len(s.Address) == 0 {
		s.Address = defaultServer.Address
	}
	if s.ReadTimeout == 0 {
		s.ReadTimeout = defaultServer.ReadTimeout
	}
	if s.WriteTimeout == 0 {
		s.WriteTimeout = defaultServer.WriteTimeout
	}
	if s.IdleTimeout == 0 {
		s.IdleTimeout = defaultServer.IdleTimeout
	}
	if s.MaxHeaderBytes == 0 {
		s.MaxHeaderBytes = defaultServer.MaxHeaderBytes
	}
	if s.ShutdownTimeout == 0 {
		s.ShutdownTimeout = defaultServer.ShutdownTimeout
	}
	return s
}

func validateAddress(address string) error {
	_, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return fmt.Errorf("%w: server.address: %v", ErrInvalidConfig, err)
	}
	return nil
}

func validatePattern(pattern string) error {
	if !strings.HasPrefix(pattern, "/") {
		return fmt.Errorf("pattern must start with /")
	}
	if !strings.HasSuffix(pattern, "/") {
		return fmt.Errorf("pattern must end with /")
	}
	if len(pattern) < 3 {
		return fmt.Errorf("pattern must be at least 3 characters long")
	}
	return nil
}

func validateURL(rawURL string) error {
	_, err := url.Parse(rawURL)
	return err
}

func validateUpstream(index int, u Upstream) error {
	if len(u.Pattern) == 0 {
		return fmt.Errorf("%w: services[%d].pattern", ErrMissingConfig, index)
	}
	err := validatePattern(u.Pattern)
	if err != nil {
		return fmt.Errorf("%w: services[%d].pattern: %v", ErrInvalidConfig, index, err)
	}
	if len(u.URL) == 0 {
		return fmt.Errorf("%w: services[%d].url", ErrMissingConfig, index)
	}
	err = validateURL(u.URL)
	if err != nil {
		return fmt.Errorf("%w: services[%d].url: %v", ErrInvalidConfig, index, err)
	}
	return nil
}

func validateUpstreams(uu []Upstream) error {
	if len(uu) == 0 {
		return fmt.Errorf("%w: services", ErrMissingConfig)
	}
	for index, u := range uu {
		err := validateUpstream(index, u)
		if err != nil {
			return err
		}
	}
	return nil
}

func validate(s Server, uu []Upstream) error {
	err := validateAddress(s.Address)
	if err != nil {
		return err
	}
	return validateUpstreams(uu)
}
