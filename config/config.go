package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/previousnext/m8s/utils/k8s/pod/sidecar"
)

// Load the config file.
func Load(file string) (Config, error) {
	var cfg Config

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to load config")
	}

	err = yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		return cfg, errors.Wrap(err, "failed to unmarshal config")
	}

	return cfg, nil
}

// LoadWithDefaults will load the config file and apply defaults if not set.
func LoadWithDefaults(file string) (Config, error) {
	cfg, err := Load(file)
	if err != nil {
		return cfg, err
	}

	if cfg.Port == 0 {
		cfg.Port = sidecar.DefaultPort
	}

	return cfg, nil
}
