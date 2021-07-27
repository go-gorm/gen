package check

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

// Member user input structures
type Member struct {
	Name          string
	Type          string
	NewType       string
	ColumnName    string
	ColumnComment string
	ModelType     string
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
