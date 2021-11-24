package main

import (
	"flag"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gen"
	"gorm.io/gorm"
	"log"
	"strings"
)

// DBType database type
type DBType string

const (
	//	Gorm Drivers mysql || postgres || sqlite || sqlserver
	DBMySQL     DBType = "mysql"
	DBPostgres  DBType = "postgres"
	DBSQLite    DBType = "sqlite"
	DBSQLServer DBType = "sqlserver"
)

func connectDB(t DBType, dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}

	switch t {
	case DBMySQL:
		return gorm.Open(mysql.Open(dsn))
	case DBPostgres:
		return gorm.Open(postgres.Open(dsn))
	case DBSQLite:
		return gorm.Open(sqlite.Open(dsn))
	case DBSQLServer:
		return gorm.Open(sqlserver.Open(dsn))
	default:
		return nil, fmt.Errorf("unknow db %q (support mysql || postgres || sqlite || sqlserver for now)", t)
	}
}

func getModels(g *gen.Generator, db *gorm.DB, tables []string) (models []interface{}, err error) {
	if len(tables) == 0 {
		//Execute tasks for all tables in the database
		tables, err = db.Migrator().GetTables()
		if err != nil {
			return nil, fmt.Errorf("GORM migrator get all tables fail: %w", err)
		}
	}

	//Execute some data table tasks
	models = make([]interface{}, len(tables))
	for i, tableName := range tables {
		models[i] = g.GenerateModel(tableName)
	}
	return models, nil
}

func main() {
	dsn := flag.String("dsn", "", "consult[https://gorm.io/docs/connecting_to_the_database.html]")
	dbType := flag.String("db", "mysql", "input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]")
	tableList := flag.String("tables", "", "enter the required data table or leave it blank")
	outPath := flag.String("outPath", "./dao/query", "specify a directory for output")
	outFile := flag.String("outFile", "", "query code file name, default: gen.go")
	withUnitTest := flag.Bool("withUnitTest", false, "generate unit test for query code")
	modelPkgName := flag.String("modelPkgName", "", "generated model code's package name")
	fieldNullable := flag.Bool("fieldNullable", false, "generate with pointer when field is nullable")
	fieldWithIndexTag := flag.Bool("fieldWithIndexTag", false, "generate field with gorm index tag")
	fieldWithTypeTag := flag.Bool("fieldWithTypeTag", false, "generate field with gorm column type tag")
	flag.Parse()

	db, err := connectDB(DBType(*dbType), *dsn)
	if err != nil {
		log.Fatalln("connect db server fail:", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           *outPath,
		OutFile:           *outFile,
		ModelPkgPath:      *modelPkgName,
		WithUnitTest:      *withUnitTest,
		FieldNullable:     *fieldNullable,
		FieldWithIndexTag: *fieldWithIndexTag,
		FieldWithTypeTag:  *fieldWithTypeTag,
	})

	g.UseDB(db)

	models, err := getModels(g, db, strings.Split(*tableList, ","))
	if err != nil {
		log.Fatalln("get tables info fail:", err)
	}

	g.ApplyBasic(models...)

	g.Execute()
}
