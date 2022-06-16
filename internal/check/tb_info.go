package check

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"

	"gorm.io/gen/internal/model"
)

const (
	//query table index
	mysqlIndexQuery = "SELECT TABLE_NAME,COLUMN_NAME,INDEX_NAME,SEQ_IN_INDEX,NON_UNIQUE " +
		"FROM information_schema.STATISTICS " +
		"WHERE TABLE_SCHEMA = ? AND TABLE_NAME =?"
	postgresIndexQuery = "select tablename,indexname,tablespace,indexdef from pg_indexes where tablename=?"
	mysql              = "mysql"
	postgres           = "postgres"
)

type ITableInfo interface {
	GetTbColumns(schemaName string, tableName string) (result []*model.Column, err error)

	GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error)
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

//GetTbColumns  struct
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

//GetTbIndex  index
func (t *defaultTableInfo) GetTbIndex(schemaName string, tableName string) (result []*model.Index, err error) {
	dn := t.db.Dialector.Name()

	switch dn {
	case mysql:
		return t.getMysqlIndex(schemaName, tableName)
	case postgres:
		return t.getPostgresIndex(tableName)
	default:
		return nil, fmt.Errorf("%s dose not support index", dn)
	}

	//if dn := t.db.Dialector.Name(); dn != "mysql" {
	//	return nil, fmt.Errorf("%s dose not support index", dn)
	//}
	//return result, t.db.Raw(mysqlIndexQuery, schemaName, tableName).Scan(&result).Error
}

func (t *defaultTableInfo) getMysqlIndex(schemaName string, tableName string) (result []*model.Index, err error) {
	return result, t.db.Raw(mysqlIndexQuery, schemaName, tableName).Scan(&result).Error
}
func (t *defaultTableInfo) getPostgresIndex(tableName string) (result []*model.Index, err error) {
	indexs := make([]struct {
		TableName string `gorm:"column:tablename"`
		IndexName string `gorm:"column:indexname"`
		IndexDef  string `gorm:"column:indexdef"`
	}, 0, 5)
	err = t.db.Raw(postgresIndexQuery, tableName).Scan(&indexs).Error
	if err != nil {
		return nil, err
	}
	reg := regexp.MustCompile(`\([^\(\)]*\)`) // reg get the index columns
	if reg == nil {
		return nil, fmt.Errorf("regexp err")
	}
	for _, index := range indexs {
		if index.IndexName == index.TableName+"_pkey" {
			continue // pk should break
		}
		columns := reg.FindAllString(index.IndexDef, -1)
		if len(columns) < 1 {
			return nil, fmt.Errorf("columns is zero")
		}
		columns[0] = strings.TrimFunc(columns[0], func(r rune) bool {
			if string(r) == `(` || string(r) == `)` {
				return true
			}
			return false
		})
		for i, s := range strings.Split(columns[0] /*get the first */, ",") {
			ind := &model.Index{}
			ind.IndexName = index.IndexName
			ind.TableName = index.TableName
			ind.NonUnique = 1
			if strings.Contains(index.IndexDef, "CREATE UNIQUE") {
				ind.NonUnique = 0
			}
			ind.ColumnName = strings.TrimSpace(s)
			ind.SeqInIndex = int32(i + 1)
			result = append(result, ind)
		}

	}
	return result, err
}
