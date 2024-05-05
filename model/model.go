package model

type Data struct {
	DataID int     `gorm:"primaryKey,autoIncrement"` // 数据ID
	Device int     // ID，用于设备标识
	Data   float64 // 数据
}

type Analysis struct {
	AnalysisID int `gorm:"primaryKey,autoIncrement"` // 分析ID
	Device     int
	Max        float64
	Min        float64
	Avg        float64 // 平均值
	Variance   float64 // 方差
}
