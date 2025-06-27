package mymodel

import (
	"gorm.io/gorm"
)

// 充电插座
type PowerSocket struct {
	gorm.Model

	ChargePointID uint
	ChargePoint   ChargePoint
}

// 充电桩
type ChargePoint struct {
	gorm.Model

	ChargingStationID uint
	ChargingStation   ChargingStation
	PowerSockets      []PowerSocket // 配置该关系后，会导致 query.ChargePoint.ChargingStation.City 关系缺失
}

// 充电站
type ChargingStation struct {
	gorm.Model

	CityID       uint
	City         City
	ChargePoints []ChargePoint
}

// 城市
type City struct {
	gorm.Model

	ProvinceID uint
	Province   Province
}

// 省份
type Province struct {
	gorm.Model
}
