package helper

import "errors"

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
	// if obj.PkgName() == "" {
	// 	return errors.New("Object's PkgName() cannot be empty")
	// }
	if obj.StructName() == "" {
		return errors.New("Object's StructName() cannot be empty")
	}

	return nil
}
