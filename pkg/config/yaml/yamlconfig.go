package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

// Config implements the Config interface.
// It gets key-value configs from a YAML file and
// set them as environmental variables.
type Config struct {
}

// NewConfig read the YAML file configs, and return an
// instance of the yaml-Config object.
func NewConfig(configFileName string) (*Config, error) {
	config, err := readYAML(configFileName)
	if err != nil {
		return nil, xerrors.Errorf("NewConfig: %w", err)
	}
	setConfigToENV(config)

	return &Config{}, nil
}

// GetConfig get the config value given it's key.
func (c *Config) GetConfig(key string) string {
	return os.Getenv(key)
}

func setConfigToENV(config map[string]interface{}) {
	for key, val := range config {
		// For each of the stated configuraton, set it as an environment variable.
		os.Setenv(key, fmt.Sprintf("%v", val))
	}
}

// readYAML read data from a YAML file into a map.
func readYAML(yamlFileName string) (map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		return nil, xerrors.Errorf("readYaml: %w", err)
	}

	yamlReceiver := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile, &yamlReceiver)
	if err != nil {
		return nil, xerrors.Errorf("readYaml: %w", err)
	}

	configMap := make(map[string]interface{})
	err = mapstructure.Decode(yamlReceiver, &configMap)
	if err != nil {
		return nil, xerrors.Errorf("readYaml: %w", err)
	}

	return configMap, nil
}
