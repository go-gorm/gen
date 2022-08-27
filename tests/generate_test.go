package tests_test

import (
	"testing"

	"gorm.io/gen"
)

func TestGenerate(t *testing.T) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "dal/query",
		Mode:    gen.WithDefaultQuery,

		WithUnitTest: true,

		FieldNullable:     true,
		FieldCoverable:    true,
		FieldWithIndexTag: true,
	})

	g.UseDB(DB)

	g.WithJSONTagNameStrategy(func(c string) string { return "-" })

	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
