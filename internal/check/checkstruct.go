package check

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gorm.io/gen/field"
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
	Members       []*Member
	Source        sourceCode

	Relations field.Relations
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
		b.appendOrUpdateMember((&Member{
			Name:       f.Name,
			Type:       b.getMemberRealType(f.FieldType),
			ColumnName: f.DBName,
		}).Revise())
	}

	b.Relations = b.parseStructRelationShip(stmt.Schema.Relationships)

	return nil
}

func (b *BaseStruct) parseStructRelationShip(relationship schema.Relationships) field.Relations {
	return field.Relations{
		HasOne:    b.pullRelationShip(relationship.HasOne),
		BelongsTo: b.pullRelationShip(relationship.BelongsTo),
		HasMany:   b.pullRelationShip(relationship.HasMany),
		Many2Many: b.pullRelationShip(relationship.Many2Many),
	}
}

func (b *BaseStruct) pullRelationShip(relationships []*schema.Relationship) []*field.Relation {
	if len(relationships) == 0 {
		return nil
	}
	result := make([]*field.Relation, len(relationships))
	for i, relationship := range relationships {
		subRelationships := relationship.FieldSchema.Relationships
		relation := field.NewRelation(
			relationship.Name,
			strings.TrimLeft(relationship.Field.FieldType.String(), "[]*"),
			b.pullRelationShip(append(append(append(append(
				make([]*schema.Relationship, 0, 4),
				subRelationships.BelongsTo...),
				subRelationships.HasOne...),
				subRelationships.HasMany...),
				subRelationships.Many2Many...),
			)...,
		)
		result[i] = relation
	}
	return result
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
	return member.Kind().String()
}

func (b *BaseStruct) ReviseMemberName() {
	for _, m := range b.Members {
		m.ReviseKeyword()
	}
}

// check member if in BaseStruct update else append
func (b *BaseStruct) appendOrUpdateMember(member *Member) {
	if member.ColumnName == "" {
		return
	}
	for index, m := range b.Members {
		if m.Name == member.Name {
			b.Members[index] = member
			return
		}
	}
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

func GetNames(bases []*BaseStruct) (res []string) {
	for _, base := range bases {
		res = append(res, base.StructName)
	}
	return res
}

func isStructType(data reflect.Value) bool {
	return data.Kind() == reflect.Struct ||
		(data.Kind() == reflect.Ptr && data.Elem().Kind() == reflect.Struct)
}
