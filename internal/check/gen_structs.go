package check

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/internal/parser"
)

/*
** The feature of mapping table from database server to Golang struct
** Provided by @qqxhb
 */

const (
	ModelPkg = "model"

	//query table structure
	columnQuery = "SELECT COLUMN_NAME ,COLUMN_COMMENT ,DATA_TYPE ,IS_NULLABLE ,COLUMN_KEY,COLUMN_TYPE,COLUMN_DEFAULT,EXTRA" +
		" FROM information_schema.columns WHERE table_schema = ? AND table_name =?"
)

var (
	defaultDataType             = "string"
	dataType        dataTypeMap = map[string]func(detailType string) string{
		"int":        func(string) string { return "int32" },
		"integer":    func(string) string { return "int32" },
		"smallint":   func(string) string { return "int32" },
		"mediumint":  func(string) string { return "int32" },
		"bigint":     func(string) string { return "int64" },
		"float":      func(string) string { return "float32" },
		"double":     func(string) string { return "float64" },
		"decimal":    func(string) string { return "float64" },
		"char":       func(string) string { return "string" },
		"varchar":    func(string) string { return "string" },
		"tinytext":   func(string) string { return "string" },
		"mediumtext": func(string) string { return "string" },
		"longtext":   func(string) string { return "string" },
		"tinyblob":   func(string) string { return "[]byte" },
		"blob":       func(string) string { return "[]byte" },
		"mediumblob": func(string) string { return "[]byte" },
		"longblob":   func(string) string { return "[]byte" },
		"text":       func(string) string { return "string" },
		"json":       func(string) string { return "string" },
		"enum":       func(string) string { return "string" },
		"time":       func(string) string { return "time.Time" },
		"date":       func(string) string { return "time.Time" },
		"datetime":   func(string) string { return "time.Time" },
		"timestamp":  func(string) string { return "time.Time" },
		"year":       func(string) string { return "int32" },
		"bit":        func(string) string { return "[]uint8" },
		"boolean":    func(string) string { return "bool" },
		"tinyint": func(detailType string) string {
			if strings.HasPrefix(detailType, "tinyint(1)") {
				return "bool"
			}
			return "int32"
		},
	}
)

type dataTypeMap map[string]func(string) string

// TODO diy type map global or single
func (m dataTypeMap) Get(dataType, detailType string) string {
	if convert, ok := m[dataType]; ok {
		return convert(detailType)
	}
	return defaultDataType
}

type SchemaNameOpt func(*gorm.DB) string
type MemberOpt func(*Member) *Member

// GenBaseStructs generate db model by table name
func GenBaseStructs(db *gorm.DB, pkg, tableName, modelName string, schemaNameOpts []SchemaNameOpt, memberOpts []MemberOpt) (bases *BaseStruct, err error) {
	if _, ok := db.Config.Dialector.(tests.DummyDialector); ok {
		return nil, fmt.Errorf("UseDB() is necessary to generate model struct [%s] from database table [%s]", modelName, tableName)
	}

	if err = checkModelName(modelName); err != nil {
		return nil, fmt.Errorf("model name %q is invalid: %w", modelName, err)
	}
	if pkg == "" {
		pkg = ModelPkg
	}
	pkg = filepath.Base(pkg)
	dbName := getSchemaName(db, schemaNameOpts...)
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

	// TODO create new member
	for _, field := range columns {
		m := modifyMember(toMember(field), memberOpts)
		if m == nil {
			continue
		}
		m.Name = db.NamingStrategy.SchemaName(m.Name)

		base.Members = append(base.Members, m.Revise())
	}

	return &base, nil
}

func toMember(field *Column) *Member {
	memberType := dataType.Get(field.DataType, field.ColumnType)
	if field.ColumnName == "deleted_at" && memberType == "time.Time" {
		memberType = "gorm.DeletedAt"
	}
	return &Member{
		Name:             field.ColumnName,
		Type:             memberType,
		ModelType:        memberType,
		ColumnName:       field.ColumnName,
		ColumnComment:    field.ColumnComment,
		MultilineComment: containMultiline(field.ColumnComment),
		GORMTag:          BuildGormTag(field.ColumnName, field.ColumnDefault),
		JSONTag:          field.ColumnName,
	}
}

func BuildGormTag(columnName, columnDefault string) string {
	if columnDefault == "" {
		return fmt.Sprintf("column:%s", columnName)
	}
	return fmt.Sprintf("column:%s;default:%s", columnName, columnDefault)
}

func modifyMember(m *Member, opts []MemberOpt) *Member {
	for _, opt := range opts {
		m = opt(m)
		if m == nil {
			break
		}
	}
	return m
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
		if name := opt(db); name != "" {
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

// get mysql db' name
var modelNameReg = regexp.MustCompile(`^\w+$`)

func checkModelName(name string) error {
	if name == "" {
		return nil
	}
	if !modelNameReg.MatchString(name) {
		return fmt.Errorf("model name cannot contains invalid character")
	}
	if name[0] < 'A' || name[0] > 'Z' {
		return fmt.Errorf("model name must be initial capital")
	}
	return nil
}
