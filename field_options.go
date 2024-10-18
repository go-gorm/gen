package gen

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/generate"
	"gorm.io/gen/internal/model"
	"gorm.io/gorm/schema"
)

// ModelOpt field option
type ModelOpt = model.Option

// Field exported model.Field
type Field = *model.Field

var ns = schema.NamingStrategy{}

var (
	FieldModify = func(opt func(Field) Field) model.ModifyFieldOpt {
		return func(f *model.Field) *model.Field {
			return opt(f)
		}
	}

	// WithDataTypesNullType configures the types of fields to use their datatypes nullable counterparts.
	/**
	 *
	 * @param {boolean} all - If true, all basic types of fields will be replaced with their `datatypes.Null[T]` types.
	 *                        If false, only fields that are allowed to be null will be replaced with `datatypes.Null[T]` types.
	 *
	 * Examples:
	 *
	 * When `all` is true:
	 * - `int64` will be replaced with `datatypes.NullInt64`
	 * - `string` will be replaced with `datatypes.NullString`
	 *
	 * When `all` is false:
	 * - Only fields that can be null (e.g., `*string` or `*int`) will be replaced with `datatypes.Null[T]` types.
	 *
	 * Note:
	 * Ensure that proper error handling is implemented when converting
	 * fields to their `datatypes.Null[T]` types to avoid runtime issues.
	 */
	WithDataTypesNullType = func(all bool) model.ModifyFieldOpt {
		return func(f *model.Field) *model.Field {
			ft := f.Type
			nullable := false
			if strings.HasPrefix(ft, "*") {
				nullable = true
				ft = strings.TrimLeft(ft, "*")
			}
			if !all && !nullable {
				return f
			}
			switch ft {
			case "time.Time", "string", "int", "int8", "int16",
				"int32", "int64", "uint", "uint8", "uint16", "uint32",
				"uint64", "float64", "float32", "byte", "bool":
				ft = fmt.Sprintf("datatypes.Null[%s]", ft)
			default:
				return f
			}
			f.CustomGenType = f.GenType()
			f.Type = ft
			return f
		}
	}

	// FieldNew add new field (any type your want)
	FieldNew = func(fieldName, fieldType string, fieldTag field.Tag) model.CreateFieldOpt {
		return func(*model.Field) *model.Field {
			return &model.Field{
				Name: fieldName,
				Type: fieldType,
				Tag:  fieldTag,
			}
		}
	}
	// FieldIgnore ignore some columns by name
	FieldIgnore = func(columnNames ...string) model.FilterFieldOpt {
		return func(m *model.Field) *model.Field {
			for _, name := range columnNames {
				if m.ColumnName == name {
					return nil
				}
			}
			return m
		}
	}
	// FieldIgnoreReg ignore some columns by RegExp
	FieldIgnoreReg = func(columnNameRegs ...string) model.FilterFieldOpt {
		regs := make([]regexp.Regexp, len(columnNameRegs))
		for i, reg := range columnNameRegs {
			regs[i] = *regexp.MustCompile(reg)
		}
		return func(m *model.Field) *model.Field {
			for _, reg := range regs {
				if reg.MatchString(m.ColumnName) {
					return nil
				}
			}
			return m
		}
	}
	// FieldRename specify field name in generated struct
	FieldRename = func(columnName string, newName string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if m.ColumnName == columnName {
				m.Name = newName
			}
			return m
		}
	}
	// FieldComment specify field comment in generated struct
	FieldComment = func(columnName string, comment string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if m.ColumnName == columnName {
				m.ColumnComment = comment
				m.MultilineComment = strings.Contains(comment, "\n")
			}
			return m
		}
	}
	// FieldType specify field type in generated struct
	FieldType = func(columnName string, newType string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if m.ColumnName == columnName {
				m.Type = newType
			}
			return m
		}
	}
	// FieldTypeReg specify field type in generated struct by RegExp
	FieldTypeReg = func(columnNameReg string, newType string) model.ModifyFieldOpt {
		reg := regexp.MustCompile(columnNameReg)
		return func(m *model.Field) *model.Field {
			if reg.MatchString(m.ColumnName) {
				m.Type = newType
			}
			return m
		}
	}
	// FieldGenType specify field gen type in generated dao
	FieldGenType = func(columnName string, newType string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if m.ColumnName == columnName {
				m.CustomGenType = newType
			}
			return m
		}
	}
	// FieldGenTypeReg specify field gen type in generated dao  by RegExp
	FieldGenTypeReg = func(columnNameReg string, newType string) model.ModifyFieldOpt {
		reg := regexp.MustCompile(columnNameReg)
		return func(m *model.Field) *model.Field {
			if reg.MatchString(m.ColumnName) {
				m.CustomGenType = newType
			}
			return m
		}
	}
	// FieldTag specify GORM tag and JSON tag
	FieldTag = func(columnName string, tagFunc func(tag field.Tag) field.Tag) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if m.ColumnName == columnName {
				m.Tag = tagFunc(m.Tag)
			}
			return m
		}
	}
	// FieldJSONTag specify JSON tag
	FieldJSONTag = func(columnName string, jsonTag string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if m.ColumnName == columnName {
				m.Tag.Set(field.TagKeyJson, jsonTag)
			}
			return m
		}
	}
	// FieldJSONTagWithNS specify JSON tag with name strategy
	FieldJSONTagWithNS = func(schemaName func(columnName string) (tagContent string)) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if schemaName != nil {
				m.Tag.Set(field.TagKeyJson, schemaName(m.ColumnName))

			}
			return m
		}
	}
	// FieldGORMTag specify GORM tag
	FieldGORMTag = func(columnName string, gormTag func(tag field.GormTag) field.GormTag) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if m.ColumnName == columnName {
				m.GORMTag = gormTag(m.GORMTag)
			}
			return m
		}
	}
	// FieldGORMTagReg specify GORM tag by RegExp
	FieldGORMTagReg = func(columnNameReg string, gormTag func(tag field.GormTag) field.GormTag) model.ModifyFieldOpt {
		reg := regexp.MustCompile(columnNameReg)
		return func(m *model.Field) *model.Field {
			if reg.MatchString(m.ColumnName) {
				m.GORMTag = gormTag(m.GORMTag)
			}
			return m
		}
	}
	// FieldNewTag add new tag
	FieldNewTag = func(columnName string, newTag field.Tag) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if m.ColumnName == columnName {
				for k, v := range newTag {
					m.Tag.Set(k, v)
				}
			}
			return m
		}
	}
	// FieldNewTagWithNS add new tag with name strategy
	FieldNewTagWithNS = func(tagName string, schemaName func(columnName string) string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			if schemaName == nil {
				schemaName = func(name string) string { return name }
			}
			m.Tag.Set(tagName, schemaName(m.ColumnName))
			return m
		}
	}
	// FieldTrimPrefix trim column name's prefix
	FieldTrimPrefix = func(prefix string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			m.Name = strings.TrimPrefix(m.Name, prefix)
			return m
		}
	}
	// FieldTrimSuffix trim column name's suffix
	FieldTrimSuffix = func(suffix string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			m.Name = strings.TrimSuffix(m.Name, suffix)
			return m
		}
	}
	// FieldAddPrefix add prefix to struct's memeber name
	FieldAddPrefix = func(prefix string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			m.Name = prefix + m.Name
			return m
		}
	}
	// FieldAddSuffix add suffix to struct's memeber name
	FieldAddSuffix = func(suffix string) model.ModifyFieldOpt {
		return func(m *model.Field) *model.Field {
			m.Name += suffix
			return m
		}
	}
	// FieldRelate relate to table in database
	FieldRelate = func(relationship field.RelationshipType, fieldName string, table *generate.QueryStructMeta, config *field.RelateConfig) model.CreateFieldOpt {
		if config == nil {
			config = &field.RelateConfig{}
		}

		return func(*model.Field) *model.Field {
			return &model.Field{
				Name:    fieldName,
				Type:    config.RelateFieldPrefix(relationship) + table.StructInfo.Type,
				Tag:     config.GetTag(fieldName),
				GORMTag: config.GORMTag,
				Relation: field.NewRelationWithType(
					relationship, fieldName, table.StructInfo.Package+"."+table.StructInfo.Type,
					table.Relations()...),
			}
		}
	}
	// FieldRelateModel relate to exist table model
	FieldRelateModel = func(relationship field.RelationshipType, fieldName string, relModel interface{}, config *field.RelateConfig) model.CreateFieldOpt {
		st := reflect.TypeOf(relModel)
		if st.Kind() == reflect.Ptr {
			st = st.Elem()
		}
		fieldType := st.String()

		if config == nil {
			config = &field.RelateConfig{}
		}

		return func(*model.Field) *model.Field {
			return &model.Field{
				Name:     fieldName,
				Type:     config.RelateFieldPrefix(relationship) + fieldType,
				GORMTag:  config.GORMTag,
				Tag:      config.GetTag(fieldName),
				Relation: field.NewRelationWithModel(relationship, fieldName, fieldType, relModel),
			}
		}
	}

	// WithMethod add custom method for table model
	WithMethod = func(methods ...interface{}) model.AddMethodOpt {
		return func() []interface{} { return methods }
	}
)

var (
	DefaultMethodTableWithNamer = (&defaultModel{}).TableName
)

type defaultModel struct {
}

func (*defaultModel) TableName(namer schema.Namer) string {
	if namer == nil {
		return "@@table"
	}
	return namer.TableName("@@table")
}
