package main

import (
	"strings"

	"gorm.io/gen/helper"
)

var _ helper.Object = new(Demo)

type Demo struct {
	structName string
	tableName  string
	fileName   string
	fields     []helper.Field
}

func (d *Demo) TableName() string        { return d.tableName }
func (d *Demo) StructName() string       { return d.structName }
func (d *Demo) FileName() string         { return d.fileName }
func (d *Demo) ImportPkgPaths() []string { return nil }
func (d *Demo) Fields() []helper.Field   { return d.fields }

type DemoField struct {
	name    string
	typ     string
	gormTag string
	jsonTag string
	tag     string
	comment string
}

func (f *DemoField) Name() string       { return f.name }
func (f *DemoField) Type() string       { return f.typ }
func (f *DemoField) ColumnName() string { return strings.ToLower(f.name) }
func (f *DemoField) GORMTag() string    { return f.gormTag }
func (f *DemoField) JSONTag() string    { return f.jsonTag }
func (f *DemoField) Tag() string        { return f.tag }
func (f *DemoField) Comment() string    { return f.comment }
