package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_LOGLEVEL     = "debug"
	DEFAULT_STATSURI     = "/stats"
	DEFAULT_READTIMEOUT  = "15s"
	DEFAULT_WRITETIMEOUT = "15s"
	DEFAULT_LOGFORMAT    = "json"
	DEFAULT_HEALTHURI    = "/health"
	DEFAULT_PORT         = "8080"
	DEFAULT_ADDRESS      = "0.0.0.0"
)

type Configuration struct {
	LogLevel      string   `yaml:"logLevel,omitempty"`
	StatsEnpoints []string `yaml:"statsEndpoints,omitempty"`
	ReadTimeout   string   `yaml:"readTimeout,omitempty"`
	WriteTimeout  string   `yaml:"writeTimeout,omitempty"`
	LogFormat     string   `yaml:"logFormat,omitempty"`
	HealthUri     string   `yaml:"healthUri,omitempty"`
	StatsUri      string   `yaml:"statsUri,omitempty"`
	Port          string   `yaml:"port,omitempty"`
	Address       string   `yaml:"address,omitempty"`
	CacheResults  bool     `yaml:"cacheResults,omitempty"`
}

func ParseConfig(path string) (*Configuration, error) {
	contents, err := getFileContents(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open config file")
	}
	return parseConfig(contents)
}

func parseConfig(contents []byte) (*Configuration, error) {
	var cfg Configuration
	if err := yaml.Unmarshal(contents, &cfg); err != nil {
		return nil, errors.Wrap(err, "unable to decode config file contents")
	}
	SetDefaults(&cfg)
	return &cfg, nil
}

func getFileContents(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func SetDefaults(cfg *Configuration) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = DEFAULT_LOGLEVEL
	}
	if cfg.ReadTimeout == "" {
		cfg.ReadTimeout = DEFAULT_READTIMEOUT
	}
	if cfg.WriteTimeout == "" {
		cfg.WriteTimeout = DEFAULT_WRITETIMEOUT
	}
	if cfg.LogFormat == "" {
		cfg.LogFormat = DEFAULT_LOGFORMAT
	}
	if cfg.HealthUri == "" {
		cfg.HealthUri = DEFAULT_HEALTHURI
	}
	if cfg.StatsUri == "" {
		cfg.StatsUri = DEFAULT_STATSURI
	}
	if cfg.Address == "" {
		cfg.Address = DEFAULT_ADDRESS
	}
	if cfg.Port == "" {
		cfg.Port = DEFAULT_PORT
	}
}
