package helper

import (
	"errors"
	"fmt"
)

type Object interface {
	PkgName() string
	TableName() string
	StructName() string
	ImportPkgPaths() []string

	Fields() []Field
}

type Field interface {
	Name() string
	Type() string
	
	GORMTag() string
	JSONTag() string
	Tag() string

	Comment() string

	PkgPath() string
	PkgAlias() string
}

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
