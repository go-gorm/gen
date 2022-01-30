package postgres

import (
	"gorm.io/gen/internal/model"
	"gorm.io/gorm"
)

// t *TableInfo gorm.io/gen/internal/check.ITableInfo
// the schemaName is the database name, not the schema name of postgres

type TableInfo struct {
	Db *gorm.DB
}

func (t *TableInfo) GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error) {
	return result, t.Db.Raw(columnQuery, tableName).Scan(&result).Error
}

func (t *TableInfo) GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error) {
	return result, t.Db.Raw(indexQuery, tableName).Scan(&result).Error
}
