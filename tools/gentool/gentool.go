package main

import (
	"flag"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const (
	//	Gorm Drivers mysql || postgres || sqlite || sqlserver
	DBMySQL     string = "mysql"
	DBPostgres  string = "postgres"
	DBSQLite    string = "sqlite"
	DBSQLServer string = "sqlserver"
)

func main() {
	dbDSN := flag.String("dsn", "", "consult[https://gorm.io/docs/connecting_to_the_database.html]")
	dbType := flag.String("db", "mysql", "You can input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]")
	tableList := flag.String("tables", "", "You can enter the required data table or leave it blank")
	outPath := flag.String("outPath", "./dao/query", "You can specify a directory for output")
	outFile := flag.String("outFile", "", "query code file name, default: gen.go")
	withUnitTest := flag.Bool("withUnitTest", true, "generate unit test for query code")
	modelsName := flag.String("modelsName", "", "generated model code's package name")
	fieldNullAble := flag.Bool("fieldNullable", false, "generate pointer when field is nullable")
	fieldWithIndexTag := flag.Bool("fieldWithIndexTag", true, "generate with gorm index tag")
	fieldWithTypeTag := flag.Bool("fieldWithTypeTag", false, "generate with gorm column type tag")
	flag.Parse()
	//dsn必须有
	if *dbDSN == "" {
		log.Fatalln("dsn must input")
		return
	}
	var tables []string
	if *tableList != "" {
		tables = strings.Split(*tableList, ",")
	}
	var db *gorm.DB
	var dbERR error
	switch *dbType {
	case DBMySQL:
		db, dbERR = gorm.Open(mysql.Open(*dbDSN), &gorm.Config{})
	case DBPostgres:
		db, dbERR = gorm.Open(postgres.Open(*dbDSN), &gorm.Config{})
	case DBSQLite:
		db, dbERR = gorm.Open(sqlite.Open(*dbDSN), &gorm.Config{})
	case DBSQLServer:
		db, dbERR = gorm.Open(sqlserver.Open(*dbDSN), &gorm.Config{})
	default:
		log.Fatalln("You can only enter Gorm Drivers mysql || postgres || sqlite || sqlserver")
		return
	}
	if dbERR != nil {
		log.Fatalln("Gorm.Open ERR:", dbERR)
		return
	}
	log.Println("NewGenerator Start")
	config := gen.Config{
		OutPath:           *outPath,
		OutFile:           *outFile,
		ModelPkgPath:      *modelsName,
		WithUnitTest:      *withUnitTest,
		FieldNullable:     *fieldNullAble,
		FieldWithIndexTag: *fieldWithIndexTag,
		FieldWithTypeTag:  *fieldWithTypeTag,
	}
	g := gen.NewGenerator(config)

	g.UseDB(db)

	if len(tables) == 0 {
		//Execute tasks for all tables in the database
		//tables = databases tables
		var err error
		tables, err = db.Migrator().GetTables()
		if err != nil {
			log.Fatalln("Gorm Migrator GetTables Err:", err)
			return
		}
	}
	//Execute some data table tasks
	for _, v := range tables {
		g.ApplyBasic(
			g.GenerateModel(v),
		)
	}

	// 执行并生成代码
	g.Execute()
	log.Println("NewGenerator End")
}
