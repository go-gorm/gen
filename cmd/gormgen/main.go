package main

import (
	"gorm.io/gen"
	"gorm.io/gen/internal/mymodel"
)

// 运行: go run cmd/gormgen/main.go 生成 myquery 包
func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/myquery",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(nil)

	models := []interface{}{
		mymodel.PowerSocket{},
		mymodel.ChargePoint{},
		mymodel.ChargingStation{},
		mymodel.City{},
	}

	g.ApplyBasic(models...)
	g.ApplyInterface(func(Querier) {}, models...)
	g.Execute()
}

// Dynamic SQL
type Querier interface {
}
