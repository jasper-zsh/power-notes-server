package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	DB DBConfig `yaml:"database"`
}

type DBConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

func LoadConfig() (*Config, error) {
	b, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
