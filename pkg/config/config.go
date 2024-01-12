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

	var config Config
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	err = validateUpstreams(config.Upstreams)
	if err != nil {
		return Config{}, err
	}

	config.Server = ensureServerDefaults(config.Server)

	return config, nil
}

func validateUpstream(index int, config Upstream) error {
	if len(config.Pattern) == 0 {
		return fmt.Errorf("%w: services[%d].pattern", ErrMissingConfig, index)
	}
	if len(config.URL) == 0 {
		return fmt.Errorf("%w: services[%d].url", ErrMissingConfig, index)
	}
	return nil
}

func validateUpstreams(configs []Upstream) error {
	if len(configs) == 0 {
		return fmt.Errorf("%w: services", ErrMissingConfig)
	}
	for index, config := range configs {
		err := validateUpstream(index, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func ensureServerDefaults(config Server) Server {
	if len(config.Address) == 0 {
		config.Address = defaultServer.Address
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = defaultServer.ReadTimeout
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = defaultServer.WriteTimeout
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = defaultServer.IdleTimeout
	}
	if config.MaxHeaderBytes == 0 {
		config.MaxHeaderBytes = defaultServer.MaxHeaderBytes
	}
	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = defaultServer.ShutdownTimeout
	}
	return config
}
