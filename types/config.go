package types

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Algo    string
	Servers []string
}

func NewConfig(path string) *Config {

	var cfg Config
	data, err := yaml.Marshal(path)
	if err != nil {
		panic("could not parse the config")
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		panic("could not marshal the config")
	}
	return &cfg
}
