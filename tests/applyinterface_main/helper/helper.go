package helper

import "gorm.io/gen"

type MainMethod interface {
	FindWrongPackage() error
}

type MainUser struct {
	ID int
}

//go:noinline
func Apply(g *gen.Generator, fc interface{}, models ...interface{}) {
	g.ApplyInterface(fc, models...)
}
