package model

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// Column table column's info
type Column struct {
	gorm.ColumnType
	TableName   string                                                        `gorm:"column:TABLE_NAME"`
	Indexes     []*Index                                                      `gorm:"-"`
	UseScanType bool                                                          `gorm:"-"`
	dataTypeMap map[string]func(columnType gorm.ColumnType) (dataType string) `gorm:"-"`
	jsonTagNS   func(columnName string) string                                `gorm:"-"`
}

// SetDataTypeMap set data type map
func (c *Column) SetDataTypeMap(m map[string]func(columnType gorm.ColumnType) (dataType string)) {
	c.dataTypeMap = m
}

// GetDataType get data type
func (c *Column) GetDataType() (fieldtype string) {
	if mapping, ok := c.dataTypeMap[c.DatabaseTypeName()]; ok {
		return mapping(c.ColumnType)
	}
	if c.UseScanType && c.ScanType() != nil {
		return c.ScanType().String()
	}
	return dataType.Get(c.DatabaseTypeName(), c.columnType())
}

// WithNS with name strategy
func (c *Column) WithNS(jsonTagNS func(columnName string) string) {
	c.jsonTagNS = jsonTagNS
	if c.jsonTagNS == nil {
		c.jsonTagNS = func(n string) string { return n }
	}
}

// ToField convert to field
func (c *Column) ToField(nullable, coverable, signable bool) *Field {
	fieldType := c.GetDataType()
	if signable && strings.Contains(c.columnType(), "unsigned") && strings.HasPrefix(fieldType, "int") {
		fieldType = "u" + fieldType
	}
	switch {
	case c.Name() == "deleted_at" && fieldType == "time.Time":
		fieldType = "gorm.DeletedAt"
	case coverable && c.needDefaultTag(c.defaultTagValue()):
		fieldType = "*" + fieldType
	case nullable && !strings.HasPrefix(fieldType, "*"):
		if n, ok := c.Nullable(); ok && n {
			fieldType = "*" + fieldType
		}
	}

	var comment string
	if c, ok := c.Comment(); ok {
		comment = c
	}

	return &Field{
		Name:             c.Name(),
		Type:             fieldType,
		ColumnName:       c.Name(),
		MultilineComment: c.multilineComment(),
		GORMTag:          c.buildGormTag(),
		Tag:              map[string]string{field.TagKeyJson: c.jsonTagNS(c.Name())},
		ColumnComment:    comment,
	}
}

func (c *Column) multilineComment() bool {
	cm, ok := c.Comment()
	return ok && strings.Contains(cm, "\n")
}

func (c *Column) buildGormTag() field.GormTag {
	tag := field.GormTag{
		field.TagKeyGormColumn: []string{c.Name()},
		field.TagKeyGormType:   []string{c.columnType()},
	}
	isPriKey, ok := c.PrimaryKey()
	isValidPriKey := ok && isPriKey
	if isValidPriKey {
		tag.Set(field.TagKeyGormPrimaryKey, "")
		if at, ok := c.AutoIncrement(); ok {
			tag.Set(field.TagKeyGormAutoIncrement, fmt.Sprintf("%t", at))
		}
	} else if n, ok := c.Nullable(); ok && !n {
		tag.Set(field.TagKeyGormNotNull, "")
	}

	for _, idx := range c.Indexes {
		if idx == nil {
			continue
		}
		if pk, _ := idx.PrimaryKey(); pk { //ignore PrimaryKey
			continue
		}
		if uniq, _ := idx.Unique(); uniq {
			tag.Append(field.TagKeyGormUniqueIndex, fmt.Sprintf("%s,priority:%d", idx.Name(), idx.Priority))
		} else {
			tag.Append(field.TagKeyGormIndex, fmt.Sprintf("%s,priority:%d", idx.Name(), idx.Priority))
		}
	}

	if dtValue := c.defaultTagValue(); c.needDefaultTag(dtValue) { // cannot set default tag for primary key
		tag.Set(field.TagKeyGormDefault, dtValue)
	}
	if comment, ok := c.Comment(); ok && comment != "" {
		if c.multilineComment() {
			comment = strings.ReplaceAll(comment, "\n", "\\n")
		}
		tag.Set(field.TagKeyGormComment, comment)
	}
	return tag
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
		return strings.Trim(defaultTagValue, "'0:-") != ""
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
