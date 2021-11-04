package check

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gorm.io/gen/internal/model"
)

const (
	//query db table list
	allTableQuery = "SELECT TABLE_NAME FROM information_schema.tables where TABLE_SCHEMA=?"
	tablesQuery   = "SELECT TABLE_NAME FROM information_schema.tables where TABLE_SCHEMA=? and TABLE_NAME in (?)"
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
	GetALLTables(schemaName string) (result []*model.Table, err error)

	GetTables(schemaName string, tableNames []string) (result []*model.Table, err error)

	GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error)

	GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error)
}

func getITableInfo(db *gorm.DB) ITableInfo {
	return &mysqlTableInfo{db: db}
}

func GetALLTables(db *gorm.DB, schemaName string) (result []*model.Table, err error) {
	if db == nil {
		return nil, errors.New("gorm db is nil")
	}

	mt := getITableInfo(db)
	result, err = mt.GetALLTables(schemaName)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetTables(db *gorm.DB, schemaName string, tableNames []string) (result []*model.Table, err error) {
	if db == nil {
		return nil, errors.New("gorm db is nil")
	}

	mt := getITableInfo(db)
	result, err = mt.GetTables(schemaName, tableNames)
	if err != nil {
		return nil, err
	}

	return result, nil
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

// GetALLTables Mysql Table List
func (t *mysqlTableInfo) GetALLTables(schemaName string) (result []*model.Table, err error) {
	return result, t.db.Raw(allTableQuery, schemaName).Scan(&result).Error
}

// GetTables Mysql Table List
func (t *mysqlTableInfo) GetTables(schemaName string, tableNames []string) (result []*model.Table, err error) {
	return result, t.db.Raw(tablesQuery, schemaName, tableNames).Scan(&result).Error
}

//GetTbColumns Mysql struct
func (t *mysqlTableInfo) GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error) {
	return result, t.db.Raw(columnQuery, schemaName, tableName).Scan(&result).Error
}

//GetTbIndex Mysql index
func (t *mysqlTableInfo) GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error) {
	return result, t.db.Raw(indexQuery, schemaName, tableName).Scan(&result).Error
}
