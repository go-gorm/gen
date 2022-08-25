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

type FieldParser interface {
	GetFieldGenType(f *schema.Field) string
}

type dummyFieldParser struct{}

func (dummyFieldParser) GetFieldGenType(*schema.Field) string { return "" }

// QueryStructMeta struct info in generated code
type QueryStructMeta struct {
	db *gorm.DB

	Generated       bool   // whether to generate db model
	FileName        string // generated file name
	S               string // the first letter(lower case)of simple Name (receiver)
	QueryStructName string // internal query struct name
	ModelStructName string // origin/model struct name
	TableName       string // table name in db server
	StructInfo      parser.Param
	Fields          []*model.Field
	Source          model.SourceCode
	ImportPkgPaths  []string
	ModelMethods    []*parser.Method // user custom method bind to db base struct

	interfaceMode bool
}

// parseStruct get all elements of struct with gorm's Parse, ignore unexported elements
func (b *QueryStructMeta) parseStruct(st interface{}) error {
	stmt := gorm.Statement{DB: b.db}

	err := stmt.Parse(st)
	if err != nil {
		return err
	}
	b.TableName = stmt.Table
	b.FileName = strings.ToLower(stmt.Table)

	var fp FieldParser = dummyFieldParser{}
	if fps, ok := st.(FieldParser); ok && fps != nil {
		fp = fps
	}
	for _, f := range stmt.Schema.Fields {
		b.appendOrUpdateField(&model.Field{
			Name:          f.Name,
			Type:          b.getFieldRealType(f.FieldType),
			ColumnName:    f.DBName,
			CustomGenType: fp.GetFieldGenType(f),
		})
	}
	for _, r := range ParseStructRelationShip(&stmt.Schema.Relationships) {
		r := r
		b.appendOrUpdateField(&model.Field{Relation: &r})
	}
	return nil
}

// getFieldRealType  get basic type of field
func (b *QueryStructMeta) getFieldRealType(f reflect.Type) string {
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
func (b *QueryStructMeta) ReviseFieldName() {
	for _, m := range b.Fields {
		m.EscapeKeyword()
	}
}

// check field if in BaseStruct update else append
func (b *QueryStructMeta) appendOrUpdateField(f *model.Field) {
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

func (b *QueryStructMeta) appendField(f *model.Field) { b.Fields = append(b.Fields, f) }

// HasField check if BaseStruct has fields
func (b *QueryStructMeta) HasField() bool { return len(b.Fields) > 0 }

// check if struct is exportable and if struct in main package and if field's type is regular
func (b *QueryStructMeta) check() (err error) {
	if b.StructInfo.InMainPkg() {
		return fmt.Errorf("can't generated data object for struct in main package, ignore:%s", b.ModelStructName)
	}
	if !isCapitalize(b.ModelStructName) {
		return fmt.Errorf("can't generated data object for non-exportable struct, ignore:%s", b.QueryStructName)
	}
	return nil
}

// Relations related field
func (b *QueryStructMeta) Relations() (result []field.Relation) {
	for _, f := range b.Fields {
		if f.IsRelation() {
			result = append(result, *f.Relation)
		}
	}
	return result
}

// StructComment struct comment
func (b *QueryStructMeta) StructComment() string {
	if b.TableName != "" {
		return fmt.Sprintf(`mapped from table <%s>`, b.TableName)
	}
	return `mapped from object`
}

// ReviseDIYMethod check diy method duplication name
func (b *QueryStructMeta) ReviseDIYMethod() error {
	var duplicateMethodName []string
	methods := make([]*parser.Method, 0, len(b.ModelMethods))
	methodMap := make(map[string]bool, len(b.ModelMethods))
	for _, method := range b.ModelMethods {
		if methodMap[method.MethodName] || method.MethodName == "TableName" {
			duplicateMethodName = append(duplicateMethodName, method.MethodName)
			continue
		}
		method.Receiver.Package = ""
		method.Receiver.Type = b.ModelStructName
		methods = append(methods, method)
		methodMap[method.MethodName] = true
	}
	b.ModelMethods = methods

	if len(duplicateMethodName) > 0 {
		return fmt.Errorf("can't generate struct with duplicated method, please check method name: %s", strings.Join(duplicateMethodName, ","))
	}
	return nil
}

func (b *QueryStructMeta) addMethodFromAddMethodOpt(methods ...interface{}) *QueryStructMeta {
	for _, method := range methods {
		modelMethods, err := parser.GetModelMethod(method, 4)
		if err != nil {
			panic("add diy method err:" + err.Error())
		}
		b.ModelMethods = append(b.ModelMethods, modelMethods.Methods...)
	}

	err := b.ReviseDIYMethod()
	if err != nil {
		b.db.Logger.Warn(context.Background(), err.Error())
	}
	return b
}

// IfaceMode object mode
func (b QueryStructMeta) IfaceMode(on bool) *QueryStructMeta {
	b.interfaceMode = on
	return &b
}

// ReturnObject return object in generated code
func (b *QueryStructMeta) ReturnObject() string {
	if b.interfaceMode {
		return fmt.Sprint("I", b.ModelStructName, "Do")
	}
	return fmt.Sprint("*", b.QueryStructName, "Do")
}

func isStructType(data reflect.Value) bool {
	return data.Kind() == reflect.Struct ||
		(data.Kind() == reflect.Ptr && data.Elem().Kind() == reflect.Struct)
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
		result[i] = *field.NewRelationWithType(field.RelationshipType(relationship.Type), relationship.Name, varType, childRelations...)
	}
	return result
}
