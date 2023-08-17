package main

import (
	"github.com/dieagenturverwaltung/gorm-gen"
	"github.com/dieagenturverwaltung/gorm-gen/examples/dal/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
		Mode:    gen.WithDefaultQuery,
	})

	// generate from struct in project
	g.ApplyBasic(model.Customer{})

	g.Execute()
}
