package config

import "data-gather-analysis-service/lib"

type Config struct {
	Port struct {
		Display  int `yaml:"display"`
		Gather   int `yaml:"gather"`
		Analysis int `yaml:"analysis"`
	} `yaml:"port"`
	MQaddr string `yaml:"mqaddr"`
	lib.DatabaseConfig
}
