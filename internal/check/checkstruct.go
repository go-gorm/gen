package check

import (
	"fmt"
	"strings"

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
	db            *gorm.DB
}

// getMembers get all elements of struct with gorm's Parse, ignore unexport elements
func (b *BaseStruct) getMembers(st interface{}) {
	stmt := gorm.Statement{DB: b.db}
	_ = stmt.Parse(st)

	for _, field := range stmt.Schema.Fields {
		b.Members = append(b.Members, &Member{
			Name:       field.Name,
			Type:       DelPointer(field.FieldType.String()),
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

// checkStructAndMembers check if struct is exportable and if member's type is regular
func (b *BaseStruct) checkStructAndMembers() (err error) {
	if !isCapitalize(b.StructName) {
		err = fmt.Errorf("ignoring non exportable struct name:%s", b.NewStructName)
		log.Println(err)
		return
	}
	for index, m := range b.Members {
		if !allowType(m.Type) {
			b.Members[index].Type = "field"
		}
		b.Members[index].NewType = getNewType(m.Type)
	}
	return nil
}

func getNewType(t string) string {
	var newType string
	for _, s := range strings.Split(t, ".") {
		newType = s
	}
	return strings.Title(newType)
}
