package check

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"gorm.io/gen/internal/check/mysql"
	"gorm.io/gen/internal/check/postgres"
	"gorm.io/gen/internal/model"
)

type ITableInfo interface {
	GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error)

	GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error)
}

func getITableInfo(db *gorm.DB) (ITableInfo, error) {
	// keep the original logic
	// only set postgres to special tableInfo impl
	switch db.Dialector.Name() {
	case "mysql":
		fallthrough
	case "sqlite":
		fallthrough
	case "sqlserver":
		return &mysql.TableInfo{Db: db}, nil
	case "postgres":
		return &postgres.TableInfo{Db: db}, nil

	default:
		return nil, fmt.Errorf("unsupported database: %s", db.Dialector.Name())
	}
}

func getTblColumns(db *gorm.DB, schemaName string, tableName string, indexTag bool) (result []*model.Column, err error) {
	if db == nil {
		return nil, errors.New("gorm db is nil")
	}

	mt, err := getITableInfo(db)
	if err != nil {
		return nil, err
	}

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
