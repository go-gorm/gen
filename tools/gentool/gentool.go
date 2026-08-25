// Command gentool generates GORM models and query code from a database schema.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// DBType database type
type DBType string

const (
	// dbMySQL Gorm Drivers mysql || postgres || sqlite || sqlserver
	dbMySQL      DBType = "mysql"
	dbPostgres   DBType = "postgres"
	dbSQLite     DBType = "sqlite"
	dbSQLServer  DBType = "sqlserver"
	dbClickHouse DBType = "clickhouse"
)
const (
	defaultQueryPath = "./dao/query"
)

// CmdParams is command line parameters
type CmdParams struct {
	DSN                 string   `yaml:"dsn"`          // consult[https://gorm.io/docs/connecting_to_the_database.html]"
	DB                  string   `yaml:"db"`           // input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
	Tables              []string `yaml:"tables"`       // enter the required data table or leave it blank
	OnlyModel           bool     `yaml:"onlyModel"`    // only generate model
	OutPath             string   `yaml:"outPath"`      // specify a directory for output
	OutFile             string   `yaml:"outFile"`      // query code file name, default: gen.go
	WithUnitTest        bool     `yaml:"withUnitTest"` // generate unit test for query code
	UnitTestTemplate    string   `yaml:"unitTestTemplate"`
	ModelPkgName        string   `yaml:"modelPkgName"`        // generated model code's package name
	FieldNullable       bool     `yaml:"fieldNullable"`       // generate with pointer when field is nullable
	FieldCoverable      bool     `yaml:"fieldCoverable"`      // generate with pointer when field has default value
	FieldWithIndexTag   bool     `yaml:"fieldWithIndexTag"`   // generate field with gorm index tag
	FieldWithTypeTag    bool     `yaml:"fieldWithTypeTag"`    // generate field with gorm column type tag
	FieldWithDefaultTag bool     `yaml:"fieldWithDefaultTag"` // generate field with gorm default tag
	FieldSignable       bool     `yaml:"fieldSignable"`       // detect integer field's unsigned type, adjust generated data type
	WithDefaultQuery    bool     `yaml:"withDefaultQuery"`    // create default query in generated code
	WithoutContext      bool     `yaml:"withoutContext"`      // generate code without context constrain
	WithQueryInterface  bool     `yaml:"withQueryInterface"`  // generate code with exported interface object
	WithGeneric         bool     `yaml:"withGeneric"`         // generate code with generic
	UseAny              bool     `yaml:"useAny"`              // emit "any" instead of "interface{}" in generated code (requires Go 1.18+)
}

func (c *CmdParams) revise() *CmdParams {
	if c == nil {
		return c
	}
	if c.DB == "" {
		c.DB = string(dbMySQL)
	}
	if c.OutPath == "" {
		c.OutPath = defaultQueryPath
	}
	if len(c.Tables) == 0 {
		return c
	}

	tableList := make([]string, 0, len(c.Tables))
	for _, tableName := range c.Tables {
		_tableName := strings.TrimSpace(tableName) // trim leading and trailing space in tableName
		if _tableName == "" {                      // skip empty tableName
			continue
		}
		tableList = append(tableList, _tableName)
	}
	c.Tables = tableList
	return c
}

// YamlConfig is yaml config struct
type YamlConfig struct {
	Version  string     `yaml:"version"`  //
	Database *CmdParams `yaml:"database"` //
}

func dialectorFor(t DBType, dsn string) (gorm.Dialector, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}

	switch t {
	case dbMySQL:
		return mysql.Open(dsn), nil
	case dbPostgres:
		return postgres.Open(dsn), nil
	case dbSQLite:
		return sqlite.Open(dsn), nil
	case dbSQLServer:
		return sqlserver.Open(dsn), nil
	case dbClickHouse:
		return clickhouse.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unknown db %q (supported: mysql, postgres, sqlite, sqlserver, clickhouse)", t)
	}
}

// connectDB chooses a driver and opens a GORM connection.
func connectDB(t DBType, dsn string) (*gorm.DB, error) {
	dialector, err := dialectorFor(t, dsn)
	if err != nil {
		return nil, err
	}
	return gorm.Open(dialector)
}

// genModels is gorm/gen generated models
func genModels(g *gen.Generator, db *gorm.DB, tables []string) (models []interface{}, err error) {
	if len(tables) == 0 {
		// Execute tasks for all tables in the database
		tables, err = db.Migrator().GetTables()
		if err != nil {
			return nil, fmt.Errorf("GORM migrator get all tables fail: %w", err)
		}
	}

	// Execute some data table tasks
	models = make([]interface{}, len(tables))
	for i, tableName := range tables {
		models[i] = g.GenerateModel(tableName)
	}
	return models, nil
}

func loadConfig(path string) (*CmdParams, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer func() { _ = file.Close() }()
	var yamlConfig YamlConfig
	if err = yaml.NewDecoder(file).Decode(&yamlConfig); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	if yamlConfig.Database == nil {
		return nil, fmt.Errorf("config is missing database section")
	}
	return yamlConfig.Database, nil
}

func parseArgs(fs *flag.FlagSet, args []string) (*CmdParams, error) {
	if fs == nil {
		return nil, fmt.Errorf("flag set cannot be nil")
	}
	var configPath, tableList string
	params := &CmdParams{}
	fs.StringVar(&configPath, "c", "", "path to gen.yml")
	fs.StringVar(&params.DSN, "dsn", "", "consult https://gorm.io/docs/connecting_to_the_database.html")
	fs.StringVar(&params.DB, "db", string(dbMySQL), "input mysql|postgres|sqlite|sqlserver|clickhouse")
	fs.StringVar(&tableList, "tables", "", "comma-separated tables, or empty for all")
	fs.BoolVar(&params.OnlyModel, "onlyModel", false, "only generate models (without query file)")
	fs.StringVar(&params.OutPath, "outPath", defaultQueryPath, "specify a directory for output")
	fs.StringVar(&params.OutFile, "outFile", "", "query code file name, default: gen.go")
	fs.BoolVar(&params.WithUnitTest, "withUnitTest", false, "generate unit test for query code")
	fs.StringVar(&params.UnitTestTemplate, "unitTestTemplate", "", "custom unit test template file path")
	fs.StringVar(&params.ModelPkgName, "modelPkgName", "", "generated model package name")
	fs.BoolVar(&params.FieldNullable, "fieldNullable", false, "generate pointers for nullable fields")
	fs.BoolVar(&params.FieldCoverable, "fieldCoverable", false, "generate pointers for fields with defaults")
	fs.BoolVar(&params.FieldWithIndexTag, "fieldWithIndexTag", false, "generate GORM index tags")
	fs.BoolVar(&params.FieldWithTypeTag, "fieldWithTypeTag", false, "generate GORM column type tags")
	fs.BoolVar(&params.FieldWithDefaultTag, "fieldWithDefaultTag", false, "generate GORM default tags")
	fs.BoolVar(&params.FieldSignable, "fieldSignable", false, "detect unsigned integer fields")
	fs.BoolVar(&params.WithDefaultQuery, "withDefaultQuery", false, "create the default query object")
	fs.BoolVar(&params.WithoutContext, "withoutContext", false, "generate APIs without context")
	fs.BoolVar(&params.WithQueryInterface, "withQueryInterface", false, "generate query interfaces")
	fs.BoolVar(&params.WithGeneric, "withGeneric", false, "generate generic APIs")
	fs.BoolVar(&params.UseAny, "useAny", false, `emit "any" instead of "interface{}"`)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if configPath != "" {
		return loadConfig(configPath)
	}
	if tableList != "" {
		params.Tables = strings.Split(tableList, ",")
	}
	return params, nil
}

func run(args []string) (err error) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			err = fmt.Errorf("generate code: %v", panicValue)
		}
	}()

	fs := flag.NewFlagSet("gentool", flag.ContinueOnError)
	config, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	config = config.revise()

	db, err := connectDB(DBType(config.DB), config.DSN)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}
	defer func() { _ = sqlDB.Close() }()

	var generateMode gen.GenerateMode
	if config.WithDefaultQuery {
		generateMode |= gen.WithDefaultQuery
	}
	if config.WithoutContext {
		generateMode |= gen.WithoutContext
	}
	if config.WithQueryInterface {
		generateMode |= gen.WithQueryInterface
	}
	if config.WithGeneric {
		generateMode |= gen.WithGeneric
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:             config.OutPath,
		OutFile:             config.OutFile,
		ModelPkgPath:        config.ModelPkgName,
		WithUnitTest:        config.WithUnitTest,
		UnitTestTemplate:    config.UnitTestTemplate,
		FieldNullable:       config.FieldNullable,
		FieldCoverable:      config.FieldCoverable,
		FieldWithIndexTag:   config.FieldWithIndexTag,
		FieldWithTypeTag:    config.FieldWithTypeTag,
		FieldWithDefaultTag: config.FieldWithDefaultTag,
		FieldSignable:       config.FieldSignable,
		UseAny:              config.UseAny,
		Mode:                generateMode,
	})

	g.UseDB(db)

	models, err := genModels(g, db, config.Tables)
	if err != nil {
		return fmt.Errorf("get table models: %w", err)
	}

	if !config.OnlyModel {
		g.ApplyBasic(models...)
	}

	g.Execute()
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
