package check

import (
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

	GenBaseStruct bool   // whether to generate db model
	S             string // the first letter(lower case)of simple Name
	NewStructName string // new struct name
	StructName    string // origin struct name
	TableName     string
	StructInfo    parser.Param
	Members       []*model.Member
	Source        model.SourceCode
}

// parseStruct get all elements of struct with gorm's Parse, ignore unexported elements
func (b *BaseStruct) parseStruct(st interface{}) error {
	stmt := gorm.Statement{DB: b.db}
	err := stmt.Parse(st)
	if err != nil {
		return err
	}
	b.TableName = stmt.Table

	for _, f := range stmt.Schema.Fields {
		b.appendOrUpdateMember(&model.Member{
			Name:       f.Name,
			Type:       b.getMemberRealType(f.FieldType),
			ColumnName: f.DBName,
		})
	}
	for _, r := range ParseStructRelationShip(&stmt.Schema.Relationships) {
		r := r
		b.appendOrUpdateMember(&model.Member{Relation: &r})
	}
	return nil
}

// getMemberRealType  get basic type of member
func (b *BaseStruct) getMemberRealType(member reflect.Type) string {
	scanValuer := reflect.TypeOf((*field.ScanValuer)(nil)).Elem()
	if member.Implements(scanValuer) || reflect.New(member).Type().Implements(scanValuer) {
		return "field"
	}

	if member.Kind() == reflect.Ptr {
		member = member.Elem()
	}
	if member.String() == "time.Time" {
		return "time.Time"
	}
	if member.String() == "[]uint8" || member.String() == "json.RawMessage" {
		return "bytes"
	}
	return member.Kind().String()
}

func (b *BaseStruct) ReviseMemberName() {
	for _, m := range b.Members {
		m.EscapeKeyword()
	}
}

// check member if in BaseStruct update else append
func (b *BaseStruct) appendOrUpdateMember(member *model.Member) {
	if member.IsRelation() {
		b.appendMember(member)
	}
	if member.ColumnName == "" {
		return
	}
	for index, m := range b.Members {
		if m.Name == member.Name {
			b.Members[index] = member
			return
		}
	}
	b.appendMember(member)
}

func (b *BaseStruct) appendMember(member *model.Member) {
	b.Members = append(b.Members, member)
}

// HasMember check if BaseStruct has members
func (b *BaseStruct) HasMember() bool { return len(b.Members) > 0 }

// check if struct is exportable and if struct in main package and if member's type is regular
func (b *BaseStruct) check() (err error) {
	if b.StructInfo.InMainPkg() {
		return fmt.Errorf("can't generated data object for struct in main package, ignore:%s", b.StructName)
	}
	if !isCapitalize(b.StructName) {
		return fmt.Errorf("can't generated data object for non-exportable struct, ignore:%s", b.NewStructName)
	}

	return nil
}

func (b *BaseStruct) Relations() []field.Relation {
	result := make([]field.Relation, 0, 4)
	for _, m := range b.Members {
		if m.IsRelation() {
			result = append(result, *m.Relation)
		}
	}
	return result
}

func GetStructNames(bases []*BaseStruct) (res []string) {
	for _, base := range bases {
		res = append(res, base.StructName)
	}
	return res
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
