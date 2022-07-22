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

func getFields(db *gorm.DB, conf *model.Config, columns []*model.Column) (fields []*model.Field) {
	for _, col := range columns {
		col.SetDataTypeMap(conf.DataTypeMap)
		col.WithNS(conf.FieldJSONTagNS, conf.FieldNewTagNS)

		m := col.ToField(conf.FieldNullable, conf.FieldCoverable, conf.FieldSignable)

		if filterField(m, conf.FilterOpts) == nil {
			continue
		}
		if t, ok := col.ColumnType.ColumnType(); ok && !conf.FieldWithTypeTag { // remove type tag if FieldWithTypeTag == false
			m.GORMTag = strings.ReplaceAll(m.GORMTag, ";type:"+t, "")
		}

		m = modifyField(m, conf.ModifyOpts)
		if ns, ok := db.NamingStrategy.(schema.NamingStrategy); ok {
			ns.SingularTable = true
			m.Name = ns.SchemaName(ns.TablePrefix + m.Name)
		} else if db.NamingStrategy != nil {
			m.Name = db.NamingStrategy.SchemaName(m.Name)
		}

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
