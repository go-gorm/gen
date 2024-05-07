package helper

import (
	"errors"
	"fmt"

	"gorm.io/gen/field"
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
	// Methods return method array
	Methods() []Method
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
	Tag() field.Tag

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

	for _, method := range obj.Methods() {
		if method.Name() == "" {
			return fmt.Errorf("object %s's Method.Name() cannot be empty", obj.StructName())
		}
	}

	return nil
}

// Method an object method interface
type Method interface {
	// Name return func name
	Name() string
	// Receiver return func receiver
	Receiver() Param
	// Comment return func comment
	Comment() string
	// Params return func input args
	Params() []Param
	// Result return func return args
	Result() []Param
	// Body return func body
	Body() string
}

// Param an method param interface
type Param interface {
	// PackagePath return package path
	PackagePath() string
	// PackageName return package name
	PackageName() string
	// TypeName return param type name
	TypeName() string
	// IsPointer return if param type is pointer
	IsPointer() bool
	// IsArray return if param type is array
	IsArray() bool
	// Name return param name
	Name() string
}
