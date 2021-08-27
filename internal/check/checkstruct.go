package check

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/parser"
)

// BaseStruct struct info in generated code
type BaseStruct struct {
	GenBaseStruct bool   // whether to generate db model
	S             string // the first letter(lower case)of simple Name
	NewStructName string // new struct name
	StructName    string // origin struct name
	TableName     string
	StructInfo    parser.Param
	Members       []*Member
	Source        source
	db            *gorm.DB
}

// getMembers get all elements of struct with gorm's Parse, ignore unexport elements
func (b *BaseStruct) getMembers(st interface{}) {
	// TODO find struct member's basic type

	stmt := gorm.Statement{DB: b.db}
	_ = stmt.Parse(st)

	for _, f := range stmt.Schema.Fields {
		b.appendOrUpdateMember(&Member{
			Name:       f.Name,
			Type:       b.getMemberRealType(f.FieldType),
			ColumnName: f.DBName,
		})
	}
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
	return member.Kind().String()
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

// getTableName get table name with gorm's Parse
func (b *BaseStruct) getTableName(st interface{}) {
	stmt := gorm.Statement{DB: b.db}
	_ = stmt.Parse(st)
	b.TableName = stmt.Table
}

// HasMember check if BaseStruct has members
func (b *BaseStruct) HasMember() bool {
	return len(b.Members) > 0
}

// check if struct is exportable and if struct in main package and if member's type is regular
func (b *BaseStruct) checkOrFix() (err error) {
	if b.StructInfo.InMainPkg() {
		return fmt.Errorf("can't generated data object for struct in main package, ignored:%s", b.StructName)
	}
	if !isCapitalize(b.StructName) {
		return fmt.Errorf("can't generated data object for non-exportable struct, ignore:%s", b.NewStructName)
	}
	for _, m := range b.Members {
		if m.IsGormDeleteAt() {
			m.Type = "time.Time"
		}
		if !m.AllowType() {
			m.Type = "field"
		}
		m.NewType = getNewTypeName(m.Type)
	}
	return nil
}

func isStructType(data reflect.Value) bool {
	return data.Kind() == reflect.Struct ||
		(data.Kind() == reflect.Ptr && data.Elem().Kind() == reflect.Struct)
}
