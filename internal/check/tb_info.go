package check

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gorm.io/gen/internal/model"
)

const (
	//query table structure
	columnQuery = "SELECT COLUMN_NAME,COLUMN_COMMENT,DATA_TYPE,IS_NULLABLE,COLUMN_KEY,COLUMN_TYPE,COLUMN_DEFAULT,EXTRA " +
		"FROM information_schema.COLUMNS " +
		"WHERE TABLE_SCHEMA = ? AND TABLE_NAME =? " +
		"ORDER BY ORDINAL_POSITION"

	//query table index
	indexQuery = "SELECT TABLE_NAME,COLUMN_NAME,INDEX_NAME,SEQ_IN_INDEX,NON_UNIQUE " +
		"FROM information_schema.STATISTICS " +
		"WHERE TABLE_SCHEMA = ? AND TABLE_NAME =?"
)

type ITableInfo interface {
	GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error)

	GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error)
}

func getITableInfo(db *gorm.DB) ITableInfo {
	return &mysqlTableInfo{db: db}
}

func getTbColumns(db *gorm.DB, schemaName string, tableName string, indexTag bool) (result []*model.Column, err error) {
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
		c.Indexes = im[c.ColumnName]
	}
	return result, nil
}

type mysqlTableInfo struct {
	db *gorm.DB
}

//GetTbColumns Mysql struct
func (t *mysqlTableInfo) GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error) {
	return result, t.db.Raw(columnQuery, schemaName, tableName).Scan(&result).Error
}

//GetTbIndex Mysql index
func (t *mysqlTableInfo) GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error) {
	return result, t.db.Raw(indexQuery, schemaName, tableName).Scan(&result).Error
}
