package model

import (
	"strings"

	"gorm.io/gorm"
)

// SchemaNameOpt schema name option
type SchemaNameOpt func(*gorm.DB) string

// TableSchemaNameOpt schema name option
type TableSchemaNameOpt func(db *gorm.DB, table string) string

var defaultSchemaNameOpt = SchemaNameOpt(func(db *gorm.DB) string {
	return db.Migrator().CurrentDatabase()
})

func DefaultTableSchemaNameOpt(opt SchemaNameOpt, table string) (TableSchemaNameOpt, string) {
	schema := ""
	if tables := strings.Split(table, `.`); len(tables) == 2 {
		table = tables[1]
		schema = tables[0]
	}
	return func(db *gorm.DB, tb string) string {
		if schema != "" && tb == table {
			return schema
		}
		if opt != nil {
			return opt(db)
		}
		return ""
	}, table
}

// FieldOpt field option
type FieldOpt interface{ Operator() func(*Field) *Field }

// ModifyFieldOpt modify field option
type ModifyFieldOpt func(*Field) *Field

// Operator implement for FieldOpt
func (o ModifyFieldOpt) Operator() func(*Field) *Field { return o }

// FilterFieldOpt filter field option
type FilterFieldOpt ModifyFieldOpt

// Operator implement for FieldOpt
func (o FilterFieldOpt) Operator() func(*Field) *Field { return o }

// CreateFieldOpt create field option
type CreateFieldOpt ModifyFieldOpt

// Operator implement for FieldOpt
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
