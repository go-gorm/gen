package gen

import (
	"gorm.io/gen/field"
	"gorm.io/gorm/schema"
	"log"
	"sync"
)

type Schema struct {
	Schema map[string]*schema.Schema
	Model  map[string]any
	*Generator
}

func (r *Schema) GetSchema(name string) *schema.Schema {
	if ret, ok := r.Schema[name]; ok {
		return ret
	} else {
		return nil
	}
}
func (r *Schema) GetModel(name string) any {
	if ret, ok := r.Model[name]; ok {
		return ret
	} else {
		return nil
	}
}
func (r *Schema) LinkModel(dst ...any) (err error) {
	var (
		parse *schema.Schema
	)
	if err = r.db.AutoMigrate(dst...); err != nil {
		return
	}
	for _, v := range dst {
		if parse, err = schema.Parse(v, &sync.Map{}, schema.NamingStrategy{}); err != nil {
			return
		}
		r.Schema[parse.Table] = parse
		r.Schema[parse.ModelType.String()] = parse
		r.Model[parse.Table] = v
	}
	return nil
}
func (r *Schema) GetModelOpt(table string) (opt []ModelOpt) {
	var (
		schema1 = r.GetSchema(table)
		schema2 *schema.Schema
	)
	if schema1 == nil {
		return
	}
	for _, item := range schema1.Relationships.Relations {
		log.Println(field.RelationshipType(item.Type))
		if schema2 = r.GetSchema(item.FieldSchema.Table); schema2 != nil {
			opt = append(opt, FieldRelate(field.RelationshipType(item.Type), item.Name, r.GenerateModel(schema2.Table), &field.RelateConfig{
				OverwriteTag: string(item.Field.Tag),
			}))
		}
	}
	for _, item := range schema1.Fields {
		if item.Tag.Get("gorm") == "-" {
			opt = append(opt, FieldNew(item.Name, item.FieldType.String(), string(item.Tag)))
		}
	}
	return opt
}
