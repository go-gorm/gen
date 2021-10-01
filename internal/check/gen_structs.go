package check

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/field"
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
		" FROM information_schema.columns WHERE table_schema = ? AND table_name =? ORDER BY ORDINAL_POSITION"
)

type SchemaNameOpt func(*gorm.DB) string

// GenBaseStructs generate db model by table name
func GenBaseStructs(db *gorm.DB, pkg, tableName, modelName string, schemaNameOpts []SchemaNameOpt, memberOpts []MemberOpt, nullable bool) (bases *BaseStruct, err error) {
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

		Relations: field.Relations{},
	}

	modifyOpts, filterOpts, createOpts := sortOpt(memberOpts)
	for _, field := range columns {
		m := field.toMember(nullable)

		if filterMember(m, filterOpts) == nil {
			continue
		}

		m = modifyMember(m, modifyOpts)
		m.Name = db.NamingStrategy.SchemaName(m.Name)

		base.Members = append(base.Members, m)
	}

	for _, create := range createOpts {
		m := create.self()(nil)

		if m.Relation != nil {
			base.Relations.Accept(m.Relation)
		} else { // Relation Field do not need SchemaName convert
			m.Name = db.NamingStrategy.SchemaName(m.Name)
		}

		base.Members = append(base.Members, m)
	}

	return &base, nil
}

func filterMember(m *Member, opts []MemberOpt) *Member {
	for _, opt := range opts {
		if opt.self()(m) == nil {
			return nil
		}
	}
	return m
}

func modifyMember(m *Member, opts []MemberOpt) *Member {
	for _, opt := range opts {
		m = opt.self()(m)
	}
	return m
}

//Mysql
func getTbColumns(db *gorm.DB, schemaName string, tableName string) (result []*Column, err error) {
	return result, db.Raw(columnQuery, schemaName, tableName).Scan(&result).Error
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
