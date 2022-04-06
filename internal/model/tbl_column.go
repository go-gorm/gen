package model

import (
	"bytes"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Column table column's info
type Column struct {
	gorm.ColumnType
	TableName   string                                               `gorm:"column:TABLE_NAME"`
	Indexes     []*Index                                             `gorm:"-"`
	UseScanType bool                                                 `gorm:"-"`
	dataTypeMap map[string]func(detailType string) (dataType string) `gorm:"-"`
	jsonTagNS   func(columnName string) string                       `gorm:"-"`
	newTagNS    func(columnName string) string                       `gorm:"-"`
}

func (c *Column) SetDataTypeMap(m map[string]func(detailType string) (dataType string)) {
	c.dataTypeMap = m
}

func (c *Column) GetDataType() (fieldtype string) {
	if mapping, ok := c.dataTypeMap[c.DatabaseTypeName()]; ok {
		return mapping(c.columnType())
	}
	if c.UseScanType && c.ScanType() != nil {
		return c.ScanType().String()
	}
	return dataType.Get(c.DatabaseTypeName(), c.columnType())
}

func (c *Column) WithNS(jsonTagNS, newTagNS func(columnName string) string) {
	c.jsonTagNS, c.newTagNS = jsonTagNS, newTagNS
	if c.jsonTagNS == nil {
		c.jsonTagNS = func(n string) string { return n }
	}
	if c.newTagNS == nil {
		c.newTagNS = func(string) string { return "" }
	}
}

func (c *Column) ToField(nullable, coverable bool) *Field {
	fieldType := c.GetDataType()
	switch {
	case c.Name() == "deleted_at" && fieldType == "time.Time":
		fieldType = "gorm.DeletedAt"
	case nullable:
		if n, ok := c.Nullable(); ok && n {
			fieldType = "*" + fieldType
		}
	case coverable && c.withDefaultValue():
		fieldType = "*" + fieldType
	}
	f := &Field{
		Name:             c.Name(),
		Type:             fieldType,
		ColumnName:       c.Name(),
		MultilineComment: c.multilineComment(),
		GORMTag:          c.buildGormTag(),
		JSONTag:          c.jsonTagNS(c.Name()),
		NewTag:           c.newTagNS(c.Name()),
	}
	if c, ok := c.Comment(); ok {
		f.ColumnComment = c
	}
	return f
}

func (c *Column) multilineComment() bool {
	cm, ok := c.Comment()
	return ok && strings.Contains(cm, "\n")
}

func (c *Column) buildGormTag() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("column:%s;type:%s", c.Name(), c.columnType()))
	if p, ok := c.PrimaryKey(); ok && p {
		buf.WriteString(";primaryKey")
		if at, ok := c.AutoIncrement(); ok {
			buf.WriteString(fmt.Sprintf(";autoIncrement:%t", at))
		}
	} else if n, ok := c.Nullable(); ok && !n {
		buf.WriteString(";not null")
	}

	for _, idx := range c.Indexes {
		if idx == nil || idx.IsPrimaryKey() {
			continue
		}
		if idx.IsUnique() {
			buf.WriteString(fmt.Sprintf(";uniqueIndex:%s,priority:%d", idx.IndexName, idx.SeqInIndex))
		} else {
			buf.WriteString(fmt.Sprintf(";index:%s,priority:%d", idx.IndexName, idx.SeqInIndex))
		}
	}
	if c.withDefaultValue() {
		buf.WriteString(fmt.Sprintf(";default:%s", c.defaultValue()))
	}
	return buf.String()
}

// withDefaultValue check if col has default value and not created_at or updated_at
func (c *Column) withDefaultValue() (normal bool) {
	return c.defaultValue() != "" && c.Name() != "created_at" && c.Name() != "updated_at"
}

func (c *Column) defaultValue() string {
	df, ok := c.DefaultValue()
	if !ok {
		return ""
	}
	if strings.TrimSpace(df) == "" {
		return "'" + df + "'"
	}
	return df
}

func (c *Column) columnType() (v string) {
	if cl, ok := c.ColumnType.ColumnType(); ok {
		return cl
	}
	return c.DatabaseTypeName()
}
