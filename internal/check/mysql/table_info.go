package mysql

import (
	"gorm.io/gen/internal/model"
	"gorm.io/gorm"
)

type TableInfo struct {
	Db *gorm.DB
}

//GetTbColumns Mysql struct
func (t *TableInfo) GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error) {
	return result, t.Db.Raw(columnQuery, schemaName, tableName).Scan(&result).Error
}

//GetTbIndex Mysql index
func (t *TableInfo) GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error) {
	return result, t.Db.Raw(indexQuery, schemaName, tableName).Scan(&result).Error
}
