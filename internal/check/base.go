package check

import (
	"bytes"
	"fmt"
	"strings"
)

type Status int

const (
	UNKNOWN Status = iota
	SQL
	DATA
	VARIABLE
	WHERE
	IF
	SET
	ELSE
	ELSEIF
	END
	BOOL
	INT
	STRING
	TIME
	OTHER
	EXPRESSION
	LOGICAL
	NIL
)

type sourceCode int

const (
	Struct sourceCode = iota
	TableName
)

var keywords = []string{
	"UnderlyingDB", "UseDB", "UseModel", "UseTable", "Quote", "Debug", "TableName", "WithContext",
	"As", "Not", "Or", "Build", "Columns", "Hints",
	"Distinct", "Omit",
	"Select", "Where", "Order", "Group", "Having", "Limit", "Offset",
	"Join", "LeftJoin", "RightJoin",
	"Save", "Create", "CreateInBatches",
	"Update", "Updates", "UpdateColumn", "UpdateColumns",
	"Find", "FindInBatches", "First", "Take", "Last", "Pluck", "Count",
	"Scan", "ScanRows", "Row", "Rows",
	"Delete", "Unscoped",
	"Scopes",
}

var (
	defaultDataType             = "string"
	dataType        dataTypeMap = map[string]func(detailType string) string{
		"int":        func(string) string { return "int32" },
		"integer":    func(string) string { return "int32" },
		"smallint":   func(string) string { return "int32" },
		"mediumint":  func(string) string { return "int32" },
		"bigint":     func(string) string { return "int64" },
		"float":      func(string) string { return "float32" },
		"double":     func(string) string { return "float64" },
		"decimal":    func(string) string { return "float64" },
		"char":       func(string) string { return "string" },
		"varchar":    func(string) string { return "string" },
		"tinytext":   func(string) string { return "string" },
		"mediumtext": func(string) string { return "string" },
		"longtext":   func(string) string { return "string" },
		"binary":     func(string) string { return "[]byte" },
		"varbinary":  func(string) string { return "[]byte" },
		"tinyblob":   func(string) string { return "[]byte" },
		"blob":       func(string) string { return "[]byte" },
		"mediumblob": func(string) string { return "[]byte" },
		"longblob":   func(string) string { return "[]byte" },
		"text":       func(string) string { return "string" },
		"json":       func(string) string { return "string" },
		"enum":       func(string) string { return "string" },
		"time":       func(string) string { return "time.Time" },
		"date":       func(string) string { return "time.Time" },
		"datetime":   func(string) string { return "time.Time" },
		"timestamp":  func(string) string { return "time.Time" },
		"year":       func(string) string { return "int32" },
		"bit":        func(string) string { return "[]uint8" },
		"boolean":    func(string) string { return "bool" },
		"tinyint": func(detailType string) string {
			if strings.HasPrefix(detailType, "tinyint(1)") {
				return "bool"
			}
			return "int32"
		},
	}
)

type dataTypeMap map[string]func(string) string

// TODO diy type map global or single
func (m dataTypeMap) Get(dataType, detailType string) string {
	if convert, ok := m[dataType]; ok {
		return convert(detailType)
	}
	return defaultDataType
}

// Member user input structures
type Member struct {
	Name             string
	Type             string
	ColumnName       string
	ColumnComment    string
	MultilineComment bool
	JSONTag          string
	GORMTag          string
	NewTag           string
	OverwriteTag     string
}

func (m *Member) GenType() string {
	switch m.Type {
	case "string", "bytes":
		return strings.Title(m.Type)
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return strings.Title(m.Type)
	case "float64", "float32":
		return strings.Title(m.Type)
	case "bool":
		return strings.Title(m.Type)
	case "time.Time":
		return "Time"
	default:
		return "Field"
	}
}

func (m *Member) EscapeKeyword() *Member {
	if contains(m.Name, keywords) {
		m.Name += "_"
	}
	return m
}

// Column table column's info
type Column struct {
	TableName     string `gorm:"column:TABLE_NAME"`
	ColumnName    string `gorm:"column:COLUMN_NAME"`
	ColumnComment string `gorm:"column:COLUMN_COMMENT"`
	DataType      string `gorm:"column:DATA_TYPE"`
	ColumnKey     string `gorm:"column:COLUMN_KEY"`
	ColumnType    string `gorm:"column:COLUMN_TYPE"`
	ColumnDefault string `gorm:"column:COLUMN_DEFAULT"`
	Extra         string `gorm:"column:EXTRA"`
	IsNullable    string `gorm:"column:IS_NULLABLE"`
}

func (c *Column) IsPrimaryKey() bool {
	if c == nil {
		return false
	}
	if c.ColumnKey == "PRI" {
		return true
	}
	return false
}

func (c *Column) AutoIncrement() bool {
	if c == nil {
		return false
	}
	if c.Extra == "auto_increment" {
		return true
	}
	return false
}

func (c *Column) toMember(nullable bool) *Member {
	memberType := dataType.Get(c.DataType, c.ColumnType)
	if c.ColumnName == "deleted_at" && memberType == "time.Time" {
		memberType = "gorm.DeletedAt"
	}
	if nullable && c.IsNullable == "YES" {
		memberType = "*" + memberType
	}
	return &Member{
		Name:             c.ColumnName,
		Type:             memberType,
		ColumnName:       c.ColumnName,
		ColumnComment:    c.ColumnComment,
		MultilineComment: c.multilineComment(),
		GORMTag:          c.buildGormTag(),
		JSONTag:          c.ColumnName,
	}
}

func (c *Column) multilineComment() bool { return strings.Contains(c.ColumnComment, "\n") }

func (c *Column) buildGormTag() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("column:%s;type:%s", c.ColumnName, c.ColumnType))
	if c.IsPrimaryKey() {
		buf.WriteString(";primaryKey")
		if !c.AutoIncrement() {
			// integer PrioritizedPrimaryField enables AutoIncrement by default,
			// if not, we need to turn off autoIncrement for the fields
			buf.WriteString(";autoIncrement:false")
		}
	}
	if c.ColumnDefault != "" {
		buf.WriteString(fmt.Sprintf(";default:%s", c.ColumnDefault))
	}
	return buf.String()
}

type sql struct{ bytes.Buffer }

func (s *sql) WriteSql(b byte) {
	switch b {
	case '\n', '\t', ' ':
		if s.Len() == 0 || s.Bytes()[s.Len()-1] != ' ' {
			_ = s.WriteByte(' ')
		}
	default:
		_ = s.WriteByte(b)
	}
}

func (s *sql) Dump() string {
	defer s.Reset()
	return s.String()
}
