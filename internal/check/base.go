package check

import "bytes"

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
)

type source int

const (
	Struct source = iota
	TableName
)

var keywords = []string{
	"UnderlyingDB", "UseDB", "UseModel", "UseTable", "Quote", "Debug", "TableName",
	"As", "Not", "Or", "Build", "Columns", "Hints",
	"Distinct", "Omit",
	"Select", "Where", "Order", "Group", "Having", "Limit", "Offset",
	"Join", "LeftJoin", "RightJoin",
	"Save", "Create", "CreateInBatches",
	"Update", "Updates", "UpdateColumn", "UpdateColumns",
	"Find", "FindInBatches", "First", "Take", "Last", "Pluck", "Count",
	"Scan", "ScanRows", "Row", "Rows",
	"Delete", "Unscoped",
	"Transaction", "Begin", "Commit", "SavePoint", "RollBack", "RollBackTo", "Scopes",
}

// Member user input structures
type Member struct {
	Name          string
	Type          string
	NewType       string
	ColumnName    string
	ColumnComment string
	ModelType     string
	JSONTag       string
	GORMTag       string
	NewTag        string
}

// AllowType check Member Type
func (m *Member) AllowType() bool {
	switch m.Type {
	case "string", "bytes":
		return true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	case "float64", "float32":
		return true
	case "bool":
		return true
	case "time.Time":
		return true
	default:
		return false
	}
}

// fix special type and get newType
func (m *Member) Revise() *Member {
	if contains(m.Name, keywords) {
		m.Name += "_"
	}
	if !m.AllowType() {
		m.Type = "field"
	}
	m.NewType = getNewTypeName(m.Type)

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
	Extra         string `gorm:"column:EXTRA"`
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
