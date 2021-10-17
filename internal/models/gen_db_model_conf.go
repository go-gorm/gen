package models

import (
	"gorm.io/gorm"
)

type DbModelConf struct {
	ModelPkg       string
	TableName      string
	ModelName      string
	SchemaNameOpts []SchemaNameOpt
	MemberOpts     []MemberOpt
	Nullable       bool
	IndexTag       bool
}

func (cf *DbModelConf) SortOpt() (modifyOpts []MemberOpt, filterOpts []MemberOpt, createOpts []MemberOpt) {
	if cf == nil {
		return
	}
	return sortOpt(cf.MemberOpts)
}

func (cf *DbModelConf) GetSchemaName(db *gorm.DB) string {
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
