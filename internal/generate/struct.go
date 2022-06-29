package generate

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
)

// BaseStruct struct info in generated code
type BaseStruct struct {
	db *gorm.DB

	GenBaseStruct  bool   // whether to generate db model
	FileName       string // generated file name
	S              string // the first letter(lower case)of simple Name
	NewStructName  string // internal query struct name
	StructName     string // origin/model struct name
	TableName      string // table name in db server
	StructInfo     parser.Param
	Fields         []*model.Field
	Source         model.SourceCode
	ImportPkgPaths []string
	DIYMethods     []*parser.Method // user custom method bind to db base struct

	interfaceMode bool
}

// parseStruct get all elements of struct with gorm's Parse, ignore unexported elements
func (b *BaseStruct) parseStruct(st interface{}) error {
	stmt := gorm.Statement{DB: b.db}
	err := stmt.Parse(st)
	if err != nil {
		return err
	}
	b.TableName = stmt.Table
	b.FileName = strings.ToLower(stmt.Table)

	for _, f := range stmt.Schema.Fields {
		b.appendOrUpdateField(&model.Field{
			Name:       f.Name,
			Type:       b.getFieldRealType(f.FieldType),
			ColumnName: f.DBName,
		})
	}
	for _, r := range ParseStructRelationShip(&stmt.Schema.Relationships) {
		r := r
		b.appendOrUpdateField(&model.Field{Relation: &r})
	}
	return nil
}

// getFieldRealType  get basic type of field
func (b *BaseStruct) getFieldRealType(f reflect.Type) string {
	scanValuer := reflect.TypeOf((*field.ScanValuer)(nil)).Elem()
	if f.Implements(scanValuer) || reflect.New(f).Type().Implements(scanValuer) {
		return "field"
	}

	if f.Kind() == reflect.Ptr {
		f = f.Elem()
	}
	if f.String() == "time.Time" {
		return "time.Time"
	}
	if f.String() == "[]uint8" || f.String() == "json.RawMessage" {
		return "bytes"
	}
	return f.Kind().String()
}

// ReviseFieldName revise field name
func (b *BaseStruct) ReviseFieldName() {
	for _, m := range b.Fields {
		m.EscapeKeyword()
	}
}

// check field if in BaseStruct update else append
func (b *BaseStruct) appendOrUpdateField(f *model.Field) {
	if f.IsRelation() {
		b.appendField(f)
	}
	if f.ColumnName == "" {
		return
	}
	for i, m := range b.Fields {
		if m.Name == f.Name {
			b.Fields[i] = f
			return
		}
	}
	b.appendField(f)
}

func (b *BaseStruct) appendField(f *model.Field) { b.Fields = append(b.Fields, f) }

// HasField check if BaseStruct has fields
func (b *BaseStruct) HasField() bool { return len(b.Fields) > 0 }

// check if struct is exportable and if struct in main package and if field's type is regular
func (b *BaseStruct) check() (err error) {
	if b.StructInfo.InMainPkg() {
		return fmt.Errorf("can't generated data object for struct in main package, ignore:%s", b.StructName)
	}
	if !isCapitalize(b.StructName) {
		return fmt.Errorf("can't generated data object for non-exportable struct, ignore:%s", b.NewStructName)
	}
	return nil
}

// Relations related field
func (b *BaseStruct) Relations() (result []field.Relation) {
	for _, f := range b.Fields {
		if f.IsRelation() {
			result = append(result, *f.Relation)
		}
	}
	return result
}

// StructComment struct comment
func (b *BaseStruct) StructComment() string {
	if b.TableName != "" {
		return fmt.Sprintf(`mapped from table <%s>`, b.TableName)
	}
	return `mapped from object`
}

// ReviseDIYMethod check diy method duplication name
func (b *BaseStruct) ReviseDIYMethod() error {
	var duplicateMethodName []string

	methods := make([]*parser.Method, 0, len(b.DIYMethods))
	methodMap := make(map[string]bool, len(b.DIYMethods))
	for _, method := range b.DIYMethods {
		if methodMap[method.MethodName] || method.MethodName == "TableName" {
			duplicateMethodName = append(duplicateMethodName, method.MethodName)
			continue
		}
		method.BaseStruct.Package = ""
		method.BaseStruct.Type = b.StructName
		methods = append(methods, method)
		methodMap[method.MethodName] = true
	}
	b.DIYMethods = methods

	if len(duplicateMethodName) > 0 {
		return fmt.Errorf("can't generate struct with duplicated method, please check method name: %s", strings.Join(duplicateMethodName, ","))
	}
	return nil
}

// AddMethod generated model struct bind custom method, input a method of struct or a struct(bind all method of struct).
// eg: g.GenerateModel("users").AddMethod(user.IsEmpty,user.GetName) or g.GenerateModel("users").AddMethod(model.User)
func (b *BaseStruct) AddMethod(methods ...interface{}) *BaseStruct {
	for _, method := range methods {
		diyMethods, err := parser.GetDIYMethod(method)
		if err != nil {
			panic("add diy method err:" + err.Error())
		}
		b.DIYMethods = append(b.DIYMethods, diyMethods.Methods...)
	}

	err := b.ReviseDIYMethod()
	if err != nil {
		b.db.Logger.Warn(context.Background(), err.Error())
	}
	return b
}

// IfaceMode object mode
func (b BaseStruct) IfaceMode(on bool) *BaseStruct {
	b.interfaceMode = on
	return &b
}

// ReturnObject return object in generated code
func (b *BaseStruct) ReturnObject() string {
	if b.interfaceMode {
		return fmt.Sprint("I", b.StructName, "Do")
	}
	return fmt.Sprint("*", b.NewStructName, "Do")
}

// GetStructNames get struct names from base structs
func GetStructNames(bases []*BaseStruct) (names []string) {
	for _, base := range bases {
		names = append(names, base.StructName)
	}
	return names
}

func isStructType(data reflect.Value) bool {
	return data.Kind() == reflect.Struct ||
		(data.Kind() == reflect.Ptr && data.Elem().Kind() == reflect.Struct)
}

// ParseStructRelationShip parse struct's relationship
// No one should use it directly in project
func ParseStructRelationShip(relationship *schema.Relationships) []field.Relation {
	cache := make(map[string]bool)
	return append(append(append(append(
		make([]field.Relation, 0, 4),
		pullRelationShip(cache, relationship.HasOne)...),
		pullRelationShip(cache, relationship.HasMany)...),
		pullRelationShip(cache, relationship.BelongsTo)...),
		pullRelationShip(cache, relationship.Many2Many)...,
	)
}

func pullRelationShip(cache map[string]bool, relationships []*schema.Relationship) []field.Relation {
	if len(relationships) == 0 {
		return nil
	}
	result := make([]field.Relation, len(relationships))
	for i, relationship := range relationships {
		var childRelations []field.Relation
		varType := strings.TrimLeft(relationship.Field.FieldType.String(), "[]*")
		if !cache[varType] {
			cache[varType] = true
			childRelations = pullRelationShip(cache, append(append(append(append(
				make([]*schema.Relationship, 0, 4),
				relationship.FieldSchema.Relationships.BelongsTo...),
				relationship.FieldSchema.Relationships.HasOne...),
				relationship.FieldSchema.Relationships.HasMany...),
				relationship.FieldSchema.Relationships.Many2Many...),
			)
		}
		result[i] = *field.NewRelation(relationship.Name, varType, childRelations...)
	}
	return result
}
