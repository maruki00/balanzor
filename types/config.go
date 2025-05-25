package types

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Algo    string   `yaml:"algo"`
	Servers []string `yaml:"servers"`
}

func NewConfig(path string) (*Config, error) {
	yamlData, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("could not open the config file")
	}
	var cfg Config
	err = yaml.Unmarshal(yamlData, &cfg)
	if err != nil {
		return nil, errors.New("could not parse the config")
	}
	return &cfg, nil
}
