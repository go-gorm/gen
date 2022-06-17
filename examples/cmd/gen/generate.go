package main

import (
	"gorm.io/gen"
	"gorm.io/gen/examples/conf"
	"gorm.io/gen/examples/dal"
)

func init() {
	dal.DB = dal.ConnectDB(conf.MySQLDSN).Debug()

	prepare(dal.DB) // prepare table for generate
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
	})

	g.UseDB(dal.DB)

	// generate all table from database
	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
