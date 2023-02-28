package model

import (
	"bytes"
	"fmt"
	"gorm.io/gorm/schema"
	"reflect"
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
	Field       *schema.Field                                        // edit by hinego
}

// SetDataTypeMap set data type map
func (c *Column) SetDataTypeMap(m map[string]func(detailType string) (dataType string)) {
	c.dataTypeMap = m
}

// GetDataType get data type
func (c *Column) GetDataType() (fieldtype string) {
	if mapping, ok := c.dataTypeMap[c.DatabaseTypeName()]; ok {
		return mapping(c.columnType())
	}
	if c.UseScanType && c.ScanType() != nil {
		return c.ScanType().String()
	}
	return dataType.Get(c.DatabaseTypeName(), c.columnType())
}

// WithNS with name strategy
func (c *Column) WithNS(jsonTagNS, newTagNS func(columnName string) string) {
	c.jsonTagNS, c.newTagNS = jsonTagNS, newTagNS
	if c.jsonTagNS == nil {
		c.jsonTagNS = func(n string) string { return n }
	}
	if c.newTagNS == nil {
		c.newTagNS = func(string) string { return "" }
	}
}

// ToField convert to field // edit by hinego
func (c *Column) ToField(nullable, coverable, signable bool) *Field {
	var (
		FieldType  string
		DataType   = c.GetDataType()
		comment, _ = c.Comment()
		newTag     = c.newTagNS(c.Name()) + " "
		jsonTag    = c.jsonTagNS(c.Name())
	)
	if signable && strings.Contains(c.columnType(), "unsigned") && strings.HasPrefix(DataType, "int") {
		DataType = "u" + DataType
	}
	switch {
	case c.Name() == "deleted_at" && DataType == "time.Time":
		DataType = "gorm.DeletedAt"
	case coverable && c.needDefaultTag(c.defaultTagValue()):
		DataType = "*" + DataType
	case nullable:
		if n, ok := c.Nullable(); ok && n {
			DataType = "*" + DataType
		}
	}
	if c.Field != nil {
		DataType = c.Field.FieldType.String()
		FieldType = c.Field.Tag.Get("type")
		for k, v := range Parse(c.Field.Tag) {
			if k == "json" || k == "gorm" {
				continue
			}
			newTag += fmt.Sprintf(`%v:"%v" `, k, v)
		}
		if c.Field.Tag.Get("json") != "" {
			jsonTag = c.Field.Tag.Get("json")
		}
	}
	return &Field{
		Name:             c.Name(),
		CustomGenType:    FieldType,
		Type:             DataType,
		ColumnName:       c.Name(),
		MultilineComment: c.multilineComment(),
		GORMTag:          c.buildGormTag(),
		JSONTag:          jsonTag,
		NewTag:           newTag,
		ColumnComment:    comment,
	}
}

func (c *Column) multilineComment() bool {
	cm, ok := c.Comment()
	return ok && strings.Contains(cm, "\n")
}

func (c *Column) buildGormTag() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("column:%s", c.Name()))
	if c.Field != nil && c.Field.FieldType.String() == "decimal.Decimal" {
		buf.WriteString(fmt.Sprintf(";type:numeric"))
	} else {
		buf.WriteString(fmt.Sprintf(";type:%s", c.columnType()))
	}
	isPriKey, ok := c.PrimaryKey()
	isValidPriKey := ok && isPriKey
	if isValidPriKey {
		buf.WriteString(";primaryKey")
		if at, ok := c.AutoIncrement(); ok {
			buf.WriteString(fmt.Sprintf(";autoIncrement:%t", at))
		} else {
			buf.WriteString(fmt.Sprintf(";autoIncrement:false"))
		}
	} else if n, ok := c.Nullable(); ok && !n {
		buf.WriteString(";not null")
	}
	if c.Field != nil {
		if serializer, ok := c.Field.TagSettings["SERIALIZER"]; ok {
			buf.WriteString(";serializer:" + serializer)
		}
	}
	for _, idx := range c.Indexes {
		if idx == nil {
			continue
		}
		if pk, _ := idx.PrimaryKey(); pk { //ignore PrimaryKey
			continue
		}
		if uniq, _ := idx.Unique(); uniq {
			buf.WriteString(fmt.Sprintf(";uniqueIndex:%s,priority:%d", idx.Name(), idx.Priority))
		} else {
			buf.WriteString(fmt.Sprintf(";index:%s,priority:%d", idx.Name(), idx.Priority))
		}
	}

	if dtValue := c.defaultTagValue(); !isValidPriKey && c.needDefaultTag(dtValue) { // cannot set default tag for primary key
		buf.WriteString(fmt.Sprintf(`;default:%s`, dtValue))
	}
	return buf.String()
}

// needDefaultTag check if default tag needed
func (c *Column) needDefaultTag(defaultTagValue string) bool {
	if defaultTagValue == "" {
		return false
	}
	switch c.ScanType().Kind() {
	case reflect.Bool:
		return defaultTagValue != "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return defaultTagValue != "0"
	case reflect.String:
		return defaultTagValue != ""
	case reflect.Struct:
		return strings.Trim(defaultTagValue, "'0:- ") != ""
	}
	return c.Name() != "created_at" && c.Name() != "updated_at"
}

// defaultTagValue return gorm default tag's value
func (c *Column) defaultTagValue() string {
	value, ok := c.DefaultValue()
	if !ok {
		return ""
	}
	if value != "" && strings.TrimSpace(value) == "" {
		return "'" + value + "'"
	}
	return value
}

func (c *Column) columnType() (v string) {
	if cl, ok := c.ColumnType.ColumnType(); ok {
		return cl
	}
	return c.DatabaseTypeName()
}
