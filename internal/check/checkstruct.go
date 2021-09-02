package check

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"gorm.io/gen/field"
	"gorm.io/gen/internal/parser"
)

var keywords = []string{
	"UnderlyingDB", "UseDB", "UseModel", "UseTable", "Quote", "Debug", "TableName",
	"As", "Not", "Or", "Build", "Columns", "Hints",
	"Distinct", "Omit",
	"Select", "Where", "Order", "Group", "Having", "Limit", "Offset",
	"Join", "LeftJoin", "RightJoin",
	"Save", "Create", "CreateInBatches",
	"Update", "Updates", "UpdateColumn", "UpdateColumns",
	"Find", "FindInBatches", "First", "Take", "Last", "Pluck", "Count",
	"Scan", "ScanRows", "Row", "Rows",
	"Delete", "Unscoped",
	"Transaction", "Begin", "Commit", "SavePoint", "RollBack", "RollBackTo", "Scopes",
}

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

// transStruct get all elements of struct with gorm's Parse, ignore unexported elements
func (b *BaseStruct) transStruct(st interface{}) error {
	stmt := gorm.Statement{DB: b.db}
	err := stmt.Parse(st)
	if err != nil {
		return err
	}
	b.TableName = stmt.Table

	for _, f := range stmt.Schema.Fields {
		b.appendOrUpdateMember(&Member{
			Name:       f.Name,
			Type:       b.getMemberRealType(f.FieldType),
			ColumnName: f.DBName,
		})
	}

	b.fixMember()
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

// HasMember check if BaseStruct has members
func (b *BaseStruct) HasMember() bool {
	return len(b.Members) > 0
}

// check if struct is exportable and if struct in main package and if member's type is regular
func (b *BaseStruct) check() (err error) {
	if b.StructInfo.InMainPkg() {
		return fmt.Errorf("can't generated data object for struct in main package, ignored:%s", b.StructName)
	}
	if !isCapitalize(b.StructName) {
		return fmt.Errorf("can't generated data object for non-exportable struct, ignore:%s", b.NewStructName)
	}

	return nil
}

// fixMember fix special type and get newType
func (b *BaseStruct) fixMember() {
	for _, m := range b.Members {
		if contains(m.Name, keywords) {
			m.Name += "_"
		}
		if !m.AllowType() {
			m.Type = "field"
		}

		m.NewType = getNewTypeName(m.Type)
	}
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
