package main

import (
	"fmt"

	"gorm.io/gen/internal/myquery"
)

// 运行: go run cmd/test/main.go
// 需要先运行: go run cmd/gormgen/main.go 生成 myquery 包
func main() {
	a := myquery.PowerSocket.ChargePoint.ChargingStation.City.Province
	// ChargePoint 配置了 PowerSockets 关系就会导致 ChargingStation 后的关系缺失
	b := myquery.ChargePoint.ChargingStation.City
	fmt.Println(a, b)
}
