package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestParseConfig(t *testing.T) {
	testcfg := &Configuration{}
	c, err := yaml.Marshal(testcfg)
	assert.ErrorIs(t, err, nil)

	cfg, err := parseConfig(c)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}
	assert.Equal(t, cfg.LogFormat, "json")
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.Address, "0.0.0.0")
	assert.Equal(t, cfg.Port, "8080")
	assert.Equal(t, cfg.HealthUri, "/health")
	assert.Equal(t, cfg.StatsUri, "/stats")
	assert.Equal(t, cfg.ReadTimeout, "15s")
	assert.Equal(t, cfg.WriteTimeout, "15s")
	assert.Equal(t, cfg.CacheResults, false)
}

func TestParseConfigWithValues(t *testing.T) {
	testcfg := &Configuration{
		LogLevel:     "info",
		LogFormat:    "console",
		Address:      "127.0.0.1",
		Port:         "8443",
		HealthUri:    "/healthz",
		StatsUri:     "/stats",
		ReadTimeout:  "20s",
		WriteTimeout: "20s",
		StatsEnpoints: []string{
			"www.url1.com",
			"www.url2.com",
		},
		CacheResults: true,
	}
	c, err := yaml.Marshal(testcfg)
	assert.ErrorIs(t, err, nil)

	cfg, err := parseConfig(c)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}
	assert.Equal(t, cfg.LogFormat, "console")
	assert.Equal(t, cfg.LogLevel, "info")
	assert.Equal(t, cfg.Address, "127.0.0.1")
	assert.Equal(t, cfg.Port, "8443")
	assert.Equal(t, cfg.HealthUri, "/healthz")
	assert.Equal(t, cfg.StatsUri, "/stats")
	assert.Equal(t, cfg.ReadTimeout, "20s")
	assert.Equal(t, cfg.WriteTimeout, "20s")
	assert.NotEmpty(t, cfg.StatsEnpoints)
	assert.Equal(t, cfg.CacheResults, true)
}
