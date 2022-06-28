package helper

import (
	"errors"
	"fmt"
)

// Object an object interface
type Object interface {
	// TableName return table name
	TableName() string
	// StructName return struct name
	StructName() string
	// FileName return field name
	FileName() string
	// ImportPkgPaths return need import package path
	ImportPkgPaths() []string

	// Fields return field array
	Fields() []Field
}

// Field a field interface
type Field interface {
	// Name return field name
	Name() string
	// Type return field type
	Type() string

	// ColumnName return column name
	ColumnName() string
	// GORMTag return gorm tag
	GORMTag() string
	// JSONTag return json tag
	JSONTag() string
	// Tag return field tag
	Tag() string

	// Comment return comment
	Comment() string
}

// CheckObject check ojbect
func CheckObject(obj Object) error {
	if obj.StructName() == "" {
		return errors.New("Object's StructName() cannot be empty")
	}

	for _, field := range obj.Fields() {
		switch "" {
		case field.Name():
			return fmt.Errorf("Object %s's Field.Name() cannot be empty", obj.StructName())
		case field.Type():
			return fmt.Errorf("Object %s's Field.Type() cannot be empty", obj.StructName())
		}
	}
	return nil
}
