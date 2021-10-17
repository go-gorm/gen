package gen

import (
	"reflect"
	"regexp"
	"strings"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/check"
	"gorm.io/gen/internal/models"
	"gorm.io/gorm/schema"
)

var ns = schema.NamingStrategy{}

var (
	// FieldNew add new field
	FieldNew = func(fieldName, fieldType, fieldTag string) models.CreateMemberOpt {
		return func(*models.Member) *models.Member {
			return &models.Member{
				Name:         fieldName,
				Type:         fieldType,
				OverwriteTag: fieldTag,
			}
		}
	}
	// FieldIgnore ignore some columns by name
	FieldIgnore = func(columnNames ...string) models.FilterMemberOpt {
		return func(m *models.Member) *models.Member {
			for _, name := range columnNames {
				if m.ColumnName == name {
					return nil
				}
			}
			return m
		}
	}
	// FieldIgnoreReg ignore some columns by reg rule
	FieldIgnoreReg = func(columnNameRegs ...string) models.FilterMemberOpt {
		regs := make([]regexp.Regexp, len(columnNameRegs))
		for i, reg := range columnNameRegs {
			regs[i] = *regexp.MustCompile(reg)
		}
		return func(m *models.Member) *models.Member {
			for _, reg := range regs {
				if reg.MatchString(m.ColumnName) {
					return nil
				}
			}
			return m
		}
	}
	// FieldRename specify field name in generated struct
	FieldRename = func(columnName string, newName string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			if m.ColumnName == columnName {
				m.Name = newName
			}
			return m
		}
	}
	// FieldType specify field type in generated struct
	FieldType = func(columnName string, newType string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			if m.ColumnName == columnName {
				m.Type = newType
			}
			return m
		}
	}
	// FieldIgnoreType ignore some columns by reg rule
	FieldTypeReg = func(columnNameReg string, newType string) models.ModifyMemberOpt {
		reg := regexp.MustCompile(columnNameReg)
		return func(m *models.Member) *models.Member {
			if reg.MatchString(m.ColumnName) {
				m.Type = newType
			}
			return m
		}
	}
	// FieldTag specify json tag and gorm tag
	FieldTag = func(columnName string, gormTag, jsonTag string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			if m.ColumnName == columnName {
				m.GORMTag, m.JSONTag = gormTag, jsonTag
			}
			return m
		}
	}
	// FieldJSONTag specify json tag
	FieldJSONTag = func(columnName string, jsonTag string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			if m.ColumnName == columnName {
				m.JSONTag = jsonTag
			}
			return m
		}
	}
	// FieldGORMTag specify gorm tag
	FieldGORMTag = func(columnName string, gormTag string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			if m.ColumnName == columnName {
				m.GORMTag = gormTag
			}
			return m
		}
	}
	// FieldNewTag add new tag
	FieldNewTag = func(columnName string, newTag string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			if m.ColumnName == columnName {
				m.NewTag += " " + newTag
			}
			return m
		}
	}
	// FieldTrimPrefix trim column name's prefix
	FieldTrimPrefix = func(prefix string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			m.Name = strings.TrimPrefix(m.Name, prefix)
			return m
		}
	}
	// FieldTrimSuffix trim column name's suffix
	FieldTrimSuffix = func(suffix string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			m.Name = strings.TrimSuffix(m.Name, suffix)
			return m
		}
	}
	// FieldAddPrefix add prefix to struct's memeber name
	FieldAddPrefix = func(prefix string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			m.Name = prefix + m.Name
			return m
		}
	}
	// FieldAddSuffix add suffix to struct's memeber name
	FieldAddSuffix = func(suffix string) models.ModifyMemberOpt {
		return func(m *models.Member) *models.Member {
			m.Name += suffix
			return m
		}
	}
	FieldRelate = func(relationship field.RelationshipType, fieldName string, table *check.BaseStruct, config *field.RelateConfig) models.CreateMemberOpt {
		if config == nil {
			config = &field.RelateConfig{}
		}
		if config.JSONTag == "" {
			config.JSONTag = ns.ColumnName("", fieldName)
		}
		return func(*models.Member) *models.Member {
			return &models.Member{
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
	FieldRelateModel = func(relationship field.RelationshipType, fieldName string, model interface{}, config *field.RelateConfig) models.CreateMemberOpt {
		st := reflect.TypeOf(model)
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

		return func(*models.Member) *models.Member {
			return &models.Member{
				Name:         fieldName,
				Type:         config.RelateFieldPrefix(relationship) + fieldType,
				JSONTag:      config.JSONTag,
				GORMTag:      config.GORMTag,
				NewTag:       config.NewTag,
				OverwriteTag: config.OverwriteTag,

				Relation: field.NewRelationWithModel(relationship, fieldName, fieldType, model),
			}
		}
	}
)
