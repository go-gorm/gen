package generate

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gorm.io/gen/internal/model"
)

/*
** The feature of mapping table from database server to Golang struct
** Provided by @qqxhb
 */

// convertIDNaming 将字段名中的ID相关部分改为Id
func convertIDNaming(name string) string {
	// 处理以ID结尾的情况（包括单独的ID和xxxID）
	if strings.HasSuffix(name, "ID") && len(name) >= 2 {
		// 如果整个字段名就是"ID"，转换为"Id"
		if name == "ID" {
			return "Id"
		}
		// 如果ID前面是大写字母或者是单词的开始部分，转换为Id
		if len(name) > 2 {
			if name[len(name)-3] >= 'A' && name[len(name)-3] <= 'Z' {
				return strings.TrimSuffix(name, "ID") + "Id"
			}
		}
		// 如果ID前面是小写字母，也转换为Id
		if len(name) > 2 {
			if name[len(name)-3] >= 'a' && name[len(name)-3] <= 'z' {
				return strings.TrimSuffix(name, "ID") + "Id"
			}
		}
	}
	return name
}

func getFields(db *gorm.DB, conf *model.Config, columns []*model.Column) (fields []*model.Field) {
	for _, col := range columns {
		col.SetDataTypeMap(conf.DataTypeMap)
		col.WithNS(conf.FieldJSONTagNS)

		m := col.ToField(conf.FieldNullable, conf.FieldCoverable, conf.FieldSignable, conf.FieldWithDefaultTag)

		if filterField(m, conf.FilterOpts) == nil {
			continue
		}
		if _, ok := col.ColumnType.ColumnType(); ok && !conf.FieldWithTypeTag { // remove type tag if FieldWithTypeTag == false
			m.GORMTag.Remove("type")
		}

		m = modifyField(m, conf.ModifyOpts)
		if ns, ok := db.NamingStrategy.(schema.NamingStrategy); ok {
			ns.SingularTable = true
			m.Name = ns.SchemaName(ns.TablePrefix + m.Name)
		} else if db.NamingStrategy != nil {
			m.Name = db.NamingStrategy.SchemaName(m.Name)
		}

		// 处理字段名中的ID，将其转换为Id
		m.Name = convertIDNaming(m.Name)

		fields = append(fields, m)
	}
	for _, create := range conf.CreateOpts {
		m := create.Operator()(nil)
		if m.Relation != nil {
			if m.Relation.Model() != nil {
				stmt := gorm.Statement{DB: db}
				_ = stmt.Parse(m.Relation.Model())
				if stmt.Schema != nil {
					m.Relation.AppendChildRelation(ParseStructRelationShip(&stmt.Schema.Relationships)...)
				}
			}
			m.Type = strings.ReplaceAll(m.Type, conf.ModelPkg+".", "") // remove modelPkg in field's Type, avoid import error
		}

		fields = append(fields, m)
	}
	return fields
}

func filterField(m *model.Field, opts []model.FieldOption) *model.Field {
	for _, opt := range opts {
		if opt.Operator()(m) == nil {
			return nil
		}
	}
	return m
}

func modifyField(m *model.Field, opts []model.FieldOption) *model.Field {
	for _, opt := range opts {
		m = opt.Operator()(m)
	}
	return m
}

// get mysql db' name
var modelNameReg = regexp.MustCompile(`^\w+$`)

func checkStructName(name string) error {
	if name == "" {
		return nil
	}
	if !modelNameReg.MatchString(name) {
		return fmt.Errorf("model name cannot contains invalid character")
	}
	if name[0] < 'A' || name[0] > 'Z' {
		return fmt.Errorf("model name must be initial capital")
	}
	return nil
}
