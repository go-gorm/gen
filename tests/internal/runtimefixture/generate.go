package runtimefixture

import (
	"gorm.io/gen"
	fixture "gorm.io/gen/tests/fixture/model"
)

// Generate writes the checked-in runtime query fixture to outPath.
func Generate(outPath string) {
	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.ApplyBasic(fixture.Company{}, fixture.Order{})
	g.ApplyInterface(func(fixture.UserMethods) {}, fixture.User{})
	g.Execute()
}
