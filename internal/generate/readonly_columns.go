package generate

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"gorm.io/gen/internal/model"
)

func markReadOnlyColumns(db *gorm.DB, schemaName, tableName string, columns []*model.Column) error {
	if db == nil || len(columns) == 0 {
		return nil
	}

	readOnlyColumns := map[string]struct{}{}
	var err error
	switch db.Dialector.Name() {
	case "mysql":
		readOnlyColumns, err = getMySQLGeneratedColumns(db, tableName)
	case "postgres":
		readOnlyColumns, err = getPostgresGeneratedColumns(db, schemaName, tableName)
	case "sqlite":
		readOnlyColumns, err = getSQLiteGeneratedColumns(db, tableName)
	case "clickhouse":
		readOnlyColumns, err = getClickHouseGeneratedColumns(db, schemaName, tableName)
	default:
		return nil
	}
	if err != nil {
		return err
	}
	if len(readOnlyColumns) == 0 {
		return nil
	}

	for _, col := range columns {
		if col == nil {
			continue
		}
		_, ok := readOnlyColumns[col.Name()]
		col.ReadOnly = ok
	}
	return nil
}

func getMySQLGeneratedColumns(db *gorm.DB, tableName string) (map[string]struct{}, error) {
	if db == nil {
		return nil, nil
	}

	currentDB := db.Migrator().CurrentDatabase()
	readOnlyColumns := map[string]struct{}{}
	rows := make([]struct {
		ColumnName string `gorm:"column:COLUMN_NAME"`
	}, 0)

	err := db.Raw(
		`SELECT COLUMN_NAME
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = ?
	AND TABLE_NAME = ?
	AND (
		UPPER(COALESCE(GENERATION_TYPE, '')) IN ('VIRTUAL', 'STORED')
		OR UPPER(COALESCE(EXTRA, '')) LIKE '%GENERATED%'
	)`,
		currentDB, tableName,
	).Scan(&rows).Error
	if err != nil {
		rows = rows[:0]
		err = db.Raw(
			`SELECT COLUMN_NAME
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = ?
	AND TABLE_NAME = ?
	AND UPPER(COALESCE(EXTRA, '')) LIKE '%GENERATED%'`,
			currentDB, tableName,
		).Scan(&rows).Error
		if err != nil {
			return nil, err
		}
	}

	for _, row := range rows {
		if row.ColumnName == "" {
			continue
		}
		readOnlyColumns[row.ColumnName] = struct{}{}
	}
	return readOnlyColumns, nil
}

func getPostgresGeneratedColumns(db *gorm.DB, schemaName, tableName string) (map[string]struct{}, error) {
	if db == nil {
		return nil, nil
	}
	if strings.TrimSpace(schemaName) == "" {
		schemaName = "public"
	}

	rows := make([]struct {
		ColumnName string `gorm:"column:column_name"`
	}, 0)
	err := db.Raw(
		`SELECT column_name
FROM information_schema.columns
WHERE table_schema = ?
	AND table_name = ?
	AND UPPER(COALESCE(is_generated, '')) IN ('ALWAYS', 'YES')`,
		schemaName, tableName,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	readOnlyColumns := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		if row.ColumnName == "" {
			continue
		}
		readOnlyColumns[row.ColumnName] = struct{}{}
	}
	return readOnlyColumns, nil
}

func getSQLiteGeneratedColumns(db *gorm.DB, tableName string) (map[string]struct{}, error) {
	if db == nil {
		return nil, nil
	}

	escapedTableName := strings.ReplaceAll(tableName, `'`, `''`)
	rows := make([]struct {
		Name   string `gorm:"column:name"`
		Hidden int64  `gorm:"column:hidden"`
	}, 0)
	// SQLite exposes generated columns through PRAGMA table_xinfo:
	// hidden=2 for VIRTUAL generated columns, hidden=3 for STORED generated columns.
	err := db.Raw(fmt.Sprintf("PRAGMA table_xinfo('%s')", escapedTableName)).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	readOnlyColumns := make(map[string]struct{})
	for _, row := range rows {
		if row.Name == "" {
			continue
		}
		if row.Hidden == 2 || row.Hidden == 3 {
			readOnlyColumns[row.Name] = struct{}{}
		}
	}
	return readOnlyColumns, nil
}

func getClickHouseGeneratedColumns(db *gorm.DB, schemaName, tableName string) (map[string]struct{}, error) {
	if db == nil {
		return nil, nil
	}

	database := strings.TrimSpace(schemaName)
	if database == "" {
		database = db.Migrator().CurrentDatabase()
	}
	rows := make([]struct {
		Name        string `gorm:"column:name"`
		DefaultKind string `gorm:"column:default_kind"`
	}, 0)
	// ClickHouse generated/read-only columns are represented by default_kind:
	// MATERIALIZED, ALIAS and EPHEMERAL are not writable by regular INSERT/UPDATE.
	err := db.Raw(
		`SELECT name, default_kind
FROM system.columns
WHERE database = ?
	AND table = ?
	AND UPPER(COALESCE(default_kind, '')) IN ('MATERIALIZED', 'ALIAS', 'EPHEMERAL')`,
		database, tableName,
	).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	readOnlyColumns := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		if row.Name == "" {
			continue
		}
		readOnlyColumns[row.Name] = struct{}{}
	}
	return readOnlyColumns, nil
}
