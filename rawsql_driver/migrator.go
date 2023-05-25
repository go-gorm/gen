package rawsql_driver

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
)

type Migrator struct {
	migrator.Migrator
	Dialector
}

func (m Migrator) ColumnTypes(value interface{}) ([]gorm.ColumnType, error) {
	columnTypes := make([]gorm.ColumnType, 0)
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var (
			_, tableName = m.CurrentSchema(stmt, stmt.Table)
		)
		table, ok := m.tables[tableName]
		if ok && table != nil {
			columnTypes = table.ColumnTypes
		}
		return nil
	})
	return columnTypes, err
}

func (m Migrator) GetIndexes(value interface{}) ([]gorm.Index, error) {
	indexes := make([]gorm.Index, 0)
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var (
			_, tableName = m.CurrentSchema(stmt, stmt.Table)
		)
		table, ok := m.tables[tableName]
		if ok && table != nil {
			indexes = table.Indexes
		}
		return nil
	})
	return indexes, err
}

func (m Migrator) GetTables() (tableList []string, err error) {
	tableList = make([]string, 0, len(m.tables))
	for tb, _ := range m.tables {
		tableList = append(tableList, tb)
	}
	return tableList, nil
}
func (m Migrator) CurrentSchema(stmt *gorm.Statement, table string) (string, string) {
	if tables := strings.Split(table, `.`); len(tables) == 2 {
		return tables[0], tables[1]
	}
	m.DB = m.DB.Table(table)
	return "", table
}
