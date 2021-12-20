package model

import (
	"gorm.io/gorm"
)

type SchemaNameOpt func(*gorm.DB) string

var defaultSchemaNameOpt = SchemaNameOpt(func(db *gorm.DB) string {
	return db.Migrator().CurrentDatabase()
})

type FieldOpt interface{ Operator() func(*Field) *Field }

type ModifyFieldOpt func(*Field) *Field

func (o ModifyFieldOpt) Operator() func(*Field) *Field { return o }

type FilterFieldOpt ModifyFieldOpt

func (o FilterFieldOpt) Operator() func(*Field) *Field { return o }

type CreateFieldOpt ModifyFieldOpt

func (o CreateFieldOpt) Operator() func(*Field) *Field { return o }

func sortFieldOpt(opts []FieldOpt) (modifyOpts []FieldOpt, filterOpts []FieldOpt, createOpts []FieldOpt) {
	for _, opt := range opts {
		switch opt.(type) {
		case ModifyFieldOpt:
			modifyOpts = append(modifyOpts, opt)
		case FilterFieldOpt:
			filterOpts = append(filterOpts, opt)
		case CreateFieldOpt:
			createOpts = append(createOpts, opt)
		}
	}
	return
}
