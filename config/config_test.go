package config

import (
	"testing"
)

const (
	cfgStr = `{
  "local": {
    "url": "mongodb://localhost:27017/nrsv",
    "collection": "verses"
  }
}
`
	filename = "../config.json"
	env      = "local"
)

var expectedConfig Mongo = Mongo{"mongodb://localhost:27017/nrsv", "verses"}

type stringReader string

func (r stringReader) Read() ([]byte, error) {
	return []byte(r), nil
}

func TestParseConfig(t *testing.T) {
	sr := stringReader(cfgStr)
	cfg, err := readConfig(sr, env)
	if err != nil {
		t.Error(err)
	}
	if cfg != expectedConfig {
		t.Errorf("incorrect config value: %v\n", cfg)
	}
	t.Log(cfg)
}

func TestParseConfigFromFile(t *testing.T) {
	fr := fileReader(filename)
	cfg, err := readConfig(fr, env)
	if err != nil {
		t.Error(err)
	}
	if cfg != expectedConfig {
		t.Errorf("incorrect config value: %v\n", cfg)
	}
	t.Log(cfg)
}

func TestPublicParseConfigFromFile(t *testing.T) {
	cfg, _ := ReadEnvironmentFromFile(filename, env)
	if cfg != expectedConfig {
		t.Errorf("incorrect config value: %v\n", cfg)
	}
	t.Log(cfg)
}
