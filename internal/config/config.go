package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Config defines the configuration for the process manager.
type Config struct {
	LogLevel string `yaml:"log_level"`
}

// LoadConfig loads the configuration from the given file path.
func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read configuration file")
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal configuration file")
	}

	return config, nil
}

// SaveConfig saves the configuration to the given file path.
func SaveConfig(config *Config, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "failed to marshal configuration data")
	}

	err = ioutil.WriteFile(filePath, data, 0o644)
	if err != nil {
		return errors.Wrap(err, "failed to write configuration file")
	}

	return nil
}
