package gen

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/check"
	"gorm.io/gen/internal/model"
	"gorm.io/gorm/schema"
)

var ns = schema.NamingStrategy{}

var (
	// FieldNew add new field (any type your want)
	FieldNew = func(fieldName, fieldType, fieldTag string) model.CreateMemberOpt {
		return func(*model.Member) *model.Member {
			return &model.Member{
				Name:         fieldName,
				Type:         fieldType,
				OverwriteTag: fieldTag,
			}
		}
	}
	// FieldIgnore ignore some columns by name
	FieldIgnore = func(columnNames ...string) model.FilterMemberOpt {
		return func(m *model.Member) *model.Member {
			for _, name := range columnNames {
				if m.ColumnName == name {
					return nil
				}
			}
			return m
		}
	}
	// FieldIgnoreReg ignore some columns by RegExp
	FieldIgnoreReg = func(columnNameRegs ...string) model.FilterMemberOpt {
		regs := make([]regexp.Regexp, len(columnNameRegs))
		for i, reg := range columnNameRegs {
			regs[i] = *regexp.MustCompile(reg)
		}
		return func(m *model.Member) *model.Member {
			for _, reg := range regs {
				if reg.MatchString(m.ColumnName) {
					return nil
				}
			}
			return m
		}
	}
	// FieldRename specify field name in generated struct
	FieldRename = func(columnName string, newName string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			if m.ColumnName == columnName {
				m.Name = newName
			}
			return m
		}
	}
	// FieldType specify field type in generated struct
	FieldType = func(columnName string, newType string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			if m.ColumnName == columnName {
				m.Type = newType
			}
			return m
		}
	}
	// FieldIgnoreType ignore some columns by RegExp
	FieldTypeReg = func(columnNameReg string, newType string) model.ModifyMemberOpt {
		reg := regexp.MustCompile(columnNameReg)
		return func(m *model.Member) *model.Member {
			if reg.MatchString(m.ColumnName) {
				m.Type = newType
			}
			return m
		}
	}
	// FieldTag specify GORM tag and JSON tag
	FieldTag = func(columnName string, gormTag, jsonTag string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			if m.ColumnName == columnName {
				m.GORMTag, m.JSONTag = gormTag, jsonTag
			}
			return m
		}
	}
	// FieldJSONTag specify JSON tag
	FieldJSONTag = func(columnName string, jsonTag string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			if m.ColumnName == columnName {
				m.JSONTag = jsonTag
			}
			return m
		}
	}
	// FieldGORMTag specify GORM tag
	FieldGORMTag = func(columnName string, gormTag string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			if m.ColumnName == columnName {
				m.GORMTag = gormTag
			}
			return m
		}
	}
	// FieldNewTag add new tag
	FieldNewTag = func(columnName string, newTag string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			if m.ColumnName == columnName {
				m.NewTag += " " + newTag
			}
			return m
		}
	}
	// FieldNewTagWithNS add new tag with name strategy
	FieldNewTagWithNS = func(tagName string, schemaName func(columnName string) string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			if schemaName == nil {
				schemaName = func(name string) string { return name }
			}
			m.NewTag = fmt.Sprintf(`%s %s:"%s"`, m.NewTag, tagName, schemaName(m.ColumnName))
			return m
		}
	}
	// FieldTrimPrefix trim column name's prefix
	FieldTrimPrefix = func(prefix string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			m.Name = strings.TrimPrefix(m.Name, prefix)
			return m
		}
	}
	// FieldTrimSuffix trim column name's suffix
	FieldTrimSuffix = func(suffix string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			m.Name = strings.TrimSuffix(m.Name, suffix)
			return m
		}
	}
	// FieldAddPrefix add prefix to struct's memeber name
	FieldAddPrefix = func(prefix string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			m.Name = prefix + m.Name
			return m
		}
	}
	// FieldAddSuffix add suffix to struct's memeber name
	FieldAddSuffix = func(suffix string) model.ModifyMemberOpt {
		return func(m *model.Member) *model.Member {
			m.Name += suffix
			return m
		}
	}
	// FieldRelate relate to table in database
	FieldRelate = func(relationship field.RelationshipType, fieldName string, table *check.BaseStruct, config *field.RelateConfig) model.CreateMemberOpt {
		if config == nil {
			config = &field.RelateConfig{}
		}
		if config.JSONTag == "" {
			config.JSONTag = ns.ColumnName("", fieldName)
		}
		return func(*model.Member) *model.Member {
			return &model.Member{
				Name:         fieldName,
				Type:         config.RelateFieldPrefix(relationship) + table.StructInfo.Type,
				JSONTag:      config.JSONTag,
				GORMTag:      config.GORMTag,
				NewTag:       config.NewTag,
				OverwriteTag: config.OverwriteTag,

				Relation: field.NewRelationWithType(
					relationship, fieldName, table.StructInfo.Package+"."+table.StructInfo.Type,
					table.Relations()...),
			}
		}
	}
	// FieldRelateModel relate to exist table model
	FieldRelateModel = func(relationship field.RelationshipType, fieldName string, relModel interface{}, config *field.RelateConfig) model.CreateMemberOpt {
		st := reflect.TypeOf(relModel)
		if st.Kind() == reflect.Ptr {
			st = st.Elem()
		}
		fieldType := st.String()

		if config == nil {
			config = &field.RelateConfig{}
		}
		if config.JSONTag == "" {
			config.JSONTag = ns.ColumnName("", fieldName)
		}

		return func(*model.Member) *model.Member {
			return &model.Member{
				Name:         fieldName,
				Type:         config.RelateFieldPrefix(relationship) + fieldType,
				JSONTag:      config.JSONTag,
				GORMTag:      config.GORMTag,
				NewTag:       config.NewTag,
				OverwriteTag: config.OverwriteTag,

				Relation: field.NewRelationWithModel(relationship, fieldName, fieldType, relModel),
			}
		}
	}
)
