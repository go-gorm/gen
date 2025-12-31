package main

import (
	"examples/conf"
	"examples/dal"

	"gorm.io/gen"
)

func init() {
	dal.DB = dal.ConnectDB(conf.SQLiteDBName).Debug()

	prepare(dal.DB) // prepare table for generate
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../../dal/query",
		ModelPkgPath: "../../dal/model",
	})

	g.UseDB(dal.DB)

	// auto registry to models
	g.WithAutoRegistry()

	// generate all table from database
	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
