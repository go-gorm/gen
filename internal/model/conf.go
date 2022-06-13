package model

import (
	"gorm.io/gorm"
)

// FieldConf field configuration
type FieldConf struct {
	DataTypeMap map[string]func(detailType string) (dataType string)

	FieldNullable     bool // generate pointer when field is nullable
	FieldCoverable    bool // generate pointer when field has default value
	FieldSignable     bool // detect integer field's unsigned type, adjust generated data type
	FieldWithIndexTag bool // generate with gorm index tag
	FieldWithTypeTag  bool // generate with gorm column type tag

	FieldJSONTagNS func(columnName string) string
	FieldNewTagNS  func(columnName string) string

	FieldOpts []FieldOpt
}

// Conf model configuration
type Conf struct {
	ModelPkg    string
	TablePrefix string
	TableName   string
	ModelName   string

	ImportPkgPaths []string

	SchemaNameOpts []SchemaNameOpt
	TableNameNS    func(tableName string) string
	ModelNameNS    func(tableName string) string
	FileNameNS     func(tableName string) string

	FieldConf
}

func (cf *Conf) SortOpt() (modifyOpts []FieldOpt, filterOpts []FieldOpt, createOpts []FieldOpt) {
	if cf == nil {
		return
	}
	return sortFieldOpt(cf.FieldOpts)
}

func (cf *Conf) GetSchemaName(db *gorm.DB) string {
	if cf == nil {
		return defaultSchemaNameOpt(db)
	}
	for _, opt := range cf.SchemaNameOpts {
		if name := opt(db); name != "" {
			return name
		}
	}
	return defaultSchemaNameOpt(db)
}
