package check

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/inflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gorm.io/gen/internal/parser"
)

const (
	ModelPkg = "model"

	//query table structure
	columnQuery = "SELECT COLUMN_NAME ,COLUMN_COMMENT ,DATA_TYPE ,IS_NULLABLE ,COLUMN_KEY,COLUMN_TYPE,EXTRA" +
		" FROM information_schema.columns WHERE table_schema = ? AND table_name =?"
)

var dataType = map[string]string{
	"smallint":            "int32",
	"smallint unsigned":   "int32",
	"int":                 "int32",
	"int unsigned":        "int32",
	"bigint":              "int64",
	"bigint unsigned":     "int64",
	"varchar":             "string",
	"char":                "string",
	"date":                "time.Time",
	"datetime":            "time.Time",
	"bit(1)":              "[]uint8",
	"tinyint":             "int32",
	"tinyint unsigned":    "int32",
	"tinyint(1)":          "bool",
	"tinyint(1) unsigned": "bool",
	"json":                "string",
	"text":                "string",
	"timestamp":           "time.Time",
	"double":              "float64",
	"decimal":             "float64",
	"mediumtext":          "string",
	"longtext":            "string",
	"float":               "float32",
	"float unsigned":      "float32",
	"tinytext":            "string",
	"enum":                "string",
	"time":                "time.Time",
	"tinyblob":            "[]byte",
	"blob":                "[]byte",
	"mediumblob":          "[]byte",
	"longblob":            "[]byte",
	"integer":             "int32",
}

// GenBaseStructs generate db model by table name
func GenBaseStructs(db *gorm.DB, pkg string, tableName ...string) (bases []*BaseStruct, err error) {
	if isDBUndefined(db) {
		return nil, fmt.Errorf("gen config db is undefined")
	}
	if pkg == "" {
		pkg = ModelPkg
	}
	singular := singularModel(db.Config)
	dbName := getSchemaName(db)
	for _, tb := range tableName {
		columns, err := getTbColumns(db, dbName, tb)
		if err != nil {
			return nil, err
		}
		var base BaseStruct
		base.GenBaseStruct = true
		base.TableName = tb
		base.StructName = convertToModelName(singular, tb)
		base.StructInfo = parser.Param{Type: base.StructName, Package: pkg}
		for _, field := range columns {
			mt := dataType[field.DataType]
			base.Members = append(base.Members, &Member{
				Name:          nameToCamelCase(field.ColumnName),
				Type:          mt,
				ModelType:     mt,
				ColumnName:    field.ColumnName,
				ColumnComment: field.ColumnComment,
			})
		}

		base.NewStructName = strings.ToLower(base.StructName)
		base.S = string(base.NewStructName[0])
		_ = base.check()
		bases = append(bases, &base)
	}
	return
}

//Mysql
func getTbColumns(db *gorm.DB, schemaName string, tableName string) (result []*Column, err error) {
	err = db.Raw(columnQuery, schemaName, tableName).Scan(&result).Error
	return
}

// get mysql db' name
var dbNameReg = regexp.MustCompile(`/\w+\?`)

func getSchemaName(db *gorm.DB) string {
	if db == nil || db.Dialector == nil {
		return ""
	}
	myDia, ok := db.Dialector.(*mysql.Dialector)
	if !ok || myDia == nil || myDia.Config == nil {
		return ""
	}
	dbName := dbNameReg.FindString(myDia.DSN)
	if len(dbName) < 3 {
		return ""
	}
	return dbName[1 : len(dbName)-1]
}

// convert Table name or column name to camel case
func nameToCamelCase(name string) string {
	if name == "" {
		return name
	}
	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), " ", "")
}

func convertToModelName(singular bool, name string) string {
	cc := nameToCamelCase(name)
	if singular {
		return inflection.Singular(cc)
	}
	return cc
}

func singularModel(conf *gorm.Config) bool {
	if conf == nil || conf.NamingStrategy == nil {
		return false
	}
	if ns, ok := conf.NamingStrategy.(schema.NamingStrategy); ok && !ns.SingularTable {
		return true
	}
	return false
}
