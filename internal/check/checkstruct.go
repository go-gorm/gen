package check

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"gorm.io/gen/internal/parser"
	"gorm.io/gen/log"
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
	stmt := gorm.Statement{DB: b.db}
	_ = stmt.Parse(st)

	for _, field := range stmt.Schema.Fields {
		b.Members = append(b.Members, &Member{
			Name:       field.Name,
			Type:       DelPointerSym(field.FieldType.String()),
			ColumnName: field.DBName,
		})
	}
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
func (b *BaseStruct) check() (err error) {
	if b.StructInfo.InMainPkg() {
		err = fmt.Errorf("can't generated data object for struct in main package,ignored:%s", b.StructName)
		log.Println(err)
		return
	}
	if !isCapitalize(b.StructName) {
		err = fmt.Errorf("can't generated data object for non-exportable struct,ignore:%s", b.NewStructName)
		log.Println(err)
		return
	}
	for index, m := range b.Members {
		if !allowType(m.Type) {
			b.Members[index].Type = "field"
		}
		b.Members[index].NewType = getNewTypeName(m.Type)
	}
	return nil
}

func isStructType(data reflect.Value) bool {
	return data.Kind() == reflect.Struct ||
		(data.Kind() == reflect.Ptr && data.Elem().Kind() == reflect.Struct)
}
