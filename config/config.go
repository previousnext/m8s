package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"github.com/pkg/errors"
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