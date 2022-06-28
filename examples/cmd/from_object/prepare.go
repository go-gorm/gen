package main

import (
	"strings"

	"gorm.io/gen/helper"
)

var _ helper.Object = new(Demo)

// Demo demo structure
type Demo struct {
	structName string
	tableName  string
	fileName   string
	fields     []helper.Field
}

// TableName return table name
func (d *Demo) TableName() string { return d.tableName }

// StructName return struct name
func (d *Demo) StructName() string { return d.structName }

// FileName return file name
func (d *Demo) FileName() string { return d.fileName }

// ImportPkgPaths return import package paths
func (d *Demo) ImportPkgPaths() []string { return nil }

// Fields return fields
func (d *Demo) Fields() []helper.Field { return d.fields }

// DemoField demo field
type DemoField struct {
	name    string
	typ     string
	gormTag string
	jsonTag string
	tag     string
	comment string
}

// Name return name
func (f *DemoField) Name() string { return f.name }

// Type return field type
func (f *DemoField) Type() string { return f.typ }

// ColumnName return column name
func (f *DemoField) ColumnName() string { return strings.ToLower(f.name) }

// GORMTag return gorm tag
func (f *DemoField) GORMTag() string { return f.gormTag }

// JSONTag return json tag
func (f *DemoField) JSONTag() string { return f.jsonTag }

// Tag return new tag
func (f *DemoField) Tag() string { return f.tag }

// Comment return comment
func (f *DemoField) Comment() string { return f.comment }
