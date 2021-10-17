package models

import (
	"bytes"
	"strings"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/utils"
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

type SourceCode int

const (
	Struct SourceCode = iota
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

	Relation *field.Relation
}

func (m *Member) IsRelation() bool { return m.Relation != nil }

func (m *Member) GenType() string {
	if m.IsRelation() {
		return m.Type
	}

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
	if utils.ListContain(m.Name, keywords) {
		m.Name += "_"
	}
	return m
}

type Sql struct{ bytes.Buffer }

func (s *Sql) WriteSql(b byte) {
	switch b {
	case '\n', '\t', ' ':
		if s.Len() == 0 || s.Bytes()[s.Len()-1] != ' ' {
			_ = s.WriteByte(' ')
		}
	default:
		_ = s.WriteByte(b)
	}
}

func (s *Sql) Dump() string {
	defer s.Reset()
	return s.String()
}
