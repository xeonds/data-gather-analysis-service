package model

type Data struct {
	ID   int     // ID，用于设备标识
	Data float64 // 数据
}

type Analysis struct {
	ID       int
	Max      float64
	Min      float64
	Avg      float64 // 平均值
	Variance float64 // 方差
}
