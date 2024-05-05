package config

import "data-gather-analysis-service/lib"

type Config struct {
	Port struct {
		Display int // 数据展示网页面板端口
	}
	MQaddr        string
	DetectorCount int
	lib.DatabaseConfig
}
