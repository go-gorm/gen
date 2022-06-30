package generate

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gorm.io/gen/internal/model"
)

// ITableInfo table info interface
type ITableInfo interface {
	GetTableColumns(schemaName string, tableName string) (result []*model.Column, err error)

	GetTableIndex(schemaName string, tableName string) (indexes []gorm.Index, err error)
}

func getTableInfo(db *gorm.DB) ITableInfo {
	return &tableInfo{db}
}

func getTableColumns(db *gorm.DB, schemaName string, tableName string, indexTag bool) (result []*model.Column, err error) {
	if db == nil {
		return nil, errors.New("gorm db is nil")
	}

	mt := getTableInfo(db)
	result, err = mt.GetTableColumns(schemaName, tableName)
	if err != nil {
		return nil, err
	}
	if !indexTag || len(result) == 0 {
		return result, nil
	}

	index, err := mt.GetTableIndex(schemaName, tableName)
	if err != nil { //ignore find index err
		db.Logger.Warn(context.Background(), "GetTableIndex for %s,err=%s", tableName, err.Error())
		return result, nil
	}
	if len(index) == 0 {
		return result, nil
	}

	im := model.GroupByColumn(index)
	for _, c := range result {
		c.Indexes = im[c.Name()]
	}
	return result, nil
}

type tableInfo struct{ *gorm.DB }

// GetTableColumns  struct
func (t *tableInfo) GetTableColumns(schemaName string, tableName string) (result []*model.Column, err error) {
	types, err := t.Migrator().ColumnTypes(tableName)
	if err != nil {
		return nil, err
	}
	for _, column := range types {
		result = append(result, &model.Column{ColumnType: column, TableName: tableName, UseScanType: t.Dialector.Name() != "mysql" && t.Dialector.Name() != "sqlite"})
	}
	return result, nil
}

// GetTableIndex  index
func (t *tableInfo) GetTableIndex(schemaName string, tableName string) (indexes []gorm.Index, err error) {
	return t.Migrator().GetIndexes(tableName)
}
