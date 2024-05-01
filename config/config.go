package config

import "data-gather-analysis-service/lib"

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	lib.DatabaseConfig
}
