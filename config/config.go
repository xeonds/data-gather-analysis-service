package config

import "data-gather-analysis-service/lib"

type Config struct {
	Port struct {
		Display int `yaml:"display"` // 数据展示网页面板端口
	} `yaml:"port"`
	MQaddr        string `yaml:"mqaddr"`
	DetectorCount int    `yaml:"detector_count"`
	lib.DatabaseConfig
}
