package generate

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gorm.io/gen/internal/model"
)

// ITableInfo table info interface
type ITableInfo interface {
	GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error)

	GetTbIndex(schemaName string, tableName string) (indexes []gorm.Index, err error)
}

func getITableInfo(db *gorm.DB) ITableInfo {
	return &defaultTableInfo{db: db}
}

func getTblColumns(db *gorm.DB, schemaName string, tableName string, indexTag bool) (result []*model.Column, err error) {
	if db == nil {
		return nil, errors.New("gorm db is nil")
	}

	mt := getITableInfo(db)
	result, err = mt.GetTbColumns(schemaName, tableName)
	if err != nil {
		return nil, err
	}
	if !indexTag || len(result) == 0 {
		return result, nil
	}

	index, err := mt.GetTbIndex(schemaName, tableName)
	if err != nil { //ignore find index err
		db.Logger.Warn(context.Background(), "GetTbIndex for %s,err=%s", tableName, err.Error())
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

type defaultTableInfo struct {
	db *gorm.DB
}

// GetTbColumns  struct
func (t *defaultTableInfo) GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error) {
	types, err := t.db.Migrator().ColumnTypes(tableName)
	if err != nil {
		return nil, err
	}
	for _, column := range types {
		result = append(result, &model.Column{ColumnType: column, TableName: tableName, UseScanType: t.db.Dialector.Name() != "mysql" && t.db.Dialector.Name() != "sqlite"})
	}
	return result, nil
}

// GetTbIndex  index
func (t *defaultTableInfo) GetTbIndex(schemaName string, tableName string) (indexes []gorm.Index, err error) {
	return t.db.Migrator().GetIndexes(tableName)
}
