package main

import (
	"fmt"

	"gorm.io/gen/internal/myquery"
)

func main() {
	a := myquery.PowerSocket.ChargePoint.ChargingStation.City.Province
	// ChargePoint 配置了 PowerSockets 关系就会导致 ChargingStation 后的关系缺失
	b := myquery.ChargePoint.ChargingStation.City
	fmt.Println(a, b)
}
