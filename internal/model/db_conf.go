package model

import (
	"gorm.io/gorm"
)

type DBConf struct {
	ModelPkg  string
	TableName string
	ModelName string

	SchemaNameOpts []SchemaNameOpt
	MemberOpts     []MemberOpt

	GenerateModelConfig
}

type GenerateModelConfig struct {
	FieldNullable     bool // generate pointer when field is nullable
	FieldWithIndexTag bool // generate with gorm index tag
	FieldWithTypeTag  bool // generate with gorm column type tagl
}

func (cf *DBConf) SortOpt() (modifyOpts []MemberOpt, filterOpts []MemberOpt, createOpts []MemberOpt) {
	if cf == nil {
		return
	}
	return sortOpt(cf.MemberOpts)
}

func (cf *DBConf) GetSchemaName(db *gorm.DB) string {
	if cf == nil {
		return defaultMysqlSchemaNameOpt(db)
	}
	for _, opt := range cf.SchemaNameOpts {
		if name := opt(db); name != "" {
			return name
		}
	}
	return defaultMysqlSchemaNameOpt(db)
}
