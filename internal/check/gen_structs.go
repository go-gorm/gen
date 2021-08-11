package check

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

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

type SchemaNameOpt func(db *gorm.DB) string

// GenBaseStructs generate db model by table name
func GenBaseStructs(db *gorm.DB, pkg, tableName, modelName string, schemaNameOpt ...SchemaNameOpt) (bases *BaseStruct, err error) {
	if isDBUndefined(db) {
		return nil, fmt.Errorf("gen config db is undefined")
	}
	if !isModelNameValid(modelName) {
		return nil, fmt.Errorf("model name %q is invalid", modelName)
	}
	if pkg == "" {
		pkg = ModelPkg
	}
	dbName := getSchemaName(db, schemaNameOpt...)
	columns, err := getTbColumns(db, dbName, tableName)
	if err != nil {
		return nil, err
	}
	base := BaseStruct{
		Source:        TableName,
		GenBaseStruct: true,
		TableName:     tableName,
		StructName:    modelName,
		NewStructName: uncaptialize(modelName),
		S:             strings.ToLower(modelName[0:1]),
		StructInfo:    parser.Param{Type: modelName, Package: pkg},
	}
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

	_ = base.check()
	return &base, nil
}

//Mysql
func getTbColumns(db *gorm.DB, schemaName string, tableName string) (result []*Column, err error) {
	err = db.Raw(columnQuery, schemaName, tableName).Scan(&result).Error
	return
}

// get mysql db' name
var dbNameReg = regexp.MustCompile(`/\w+\??`)

func getSchemaName(db *gorm.DB, opts ...SchemaNameOpt) string {
	for _, opt := range opts {
		name := opt(db)
		if name != "" {
			return name
		}
	}
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
	end := len(dbName)
	if strings.HasSuffix(dbName, "?") {
		end--
	}
	return dbName[1:end]
}

// convert Table name or column name to camel case
func nameToCamelCase(name string) string {
	if name == "" {
		return name
	}
	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(name, "_", " ")), " ", "")
}

// get mysql db' name
var modelNameReg = regexp.MustCompile(`^\w+$`)

func isModelNameValid(name string) bool {
	if name == "" {
		return true
	}
	if !modelNameReg.MatchString(name) {
		return false
	}
	// must be captialize
	return name[0] >= 'A' && name[0] <= 'Z'
}
