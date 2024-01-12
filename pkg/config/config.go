package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Address         string        `yaml:"address"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	MaxHeaderBytes  int           `yaml:"maxHeaderBytes"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
}

type Upstream struct {
	Pattern    string `yaml:"pattern"`
	TrimPrefix bool   `yaml:"trim_pefix"`
	URL        string `yaml:"url"`
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

var ErrMissingConfig = fmt.Errorf("missing config")

func ReadConfigFromYaml(path string) (Config, error) {
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

	err = validateUpstreams(c.Upstreams)
	if err != nil {
		return Config{}, err
	}

	c.Server = ensureServerDefaults(c.Server)

	return c, nil
}

func validateUpstream(index int, u Upstream) error {
	if len(u.Pattern) == 0 {
		return fmt.Errorf("%w: services[%d].pattern", ErrMissingConfig, index)
	}
	if len(u.URL) == 0 {
		return fmt.Errorf("%w: services[%d].url", ErrMissingConfig, index)
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
