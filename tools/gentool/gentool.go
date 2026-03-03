package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
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

const defaultConfigFileName = "gen.yml"

const defaultConfigTemplate = `version: "0.1"
database:
  # consult[https://gorm.io/docs/connecting_to_the_database.html]"
  dsn : "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  # input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
  db  : "mysql"
  # enter the required data table or leave it blank.You can input :
  # tables  :
  #   - orders
  #   - users
  #   - goods
  tables  :
  # only generate models (without query file)
  onlyModel : false
  # specify a directory for output
  outPath :  "./dao/query"
  # query code file name, default: gen.go
  outFile :  ""
  # generate unit test for query code
  withUnitTest  : false
  unitTestTemplate : ""
  # generated model code's package name
  modelPkgName  : ""
  # generate with pointer when field is nullable
  fieldNullable : false
  # generate with pointer when field has default value
  fieldCoverable : false
  # generate field with gorm index tag
  fieldWithIndexTag : false
  # generate field with gorm column type tag
  fieldWithTypeTag  : false
  # generate field with gorm default tag
  fieldWithDefaultTag : false
  # detect integer field's unsigned type, adjust generated data type
  fieldSignable  : false
`

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

// connectDB choose db type for connection to database
func connectDB(t DBType, dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn cannot be empty")
	}

	switch t {
	case dbMySQL:
		return gorm.Open(mysql.Open(dsn))
	case dbPostgres:
		return gorm.Open(postgres.Open(dsn))
	case dbSQLite:
		return gorm.Open(sqlite.Open(dsn))
	case dbSQLServer:
		return gorm.Open(sqlserver.Open(dsn))
	case dbClickHouse:
		return gorm.Open(clickhouse.Open(dsn))
	default:
		return nil, fmt.Errorf("unknow db %q (support mysql || postgres || sqlite || sqlserver for now)", t)
	}
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

// parseCmdFromYaml parse cmd param from yaml
func parseCmdFromYaml(path string) (*CmdParams, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config fail: %w", err)
	}
	defer file.Close() // nolint
	var yamlConfig YamlConfig
	if err = yaml.NewDecoder(file).Decode(&yamlConfig); err != nil {
		return nil, fmt.Errorf("decode config fail: %w", err)
	}
	return yamlConfig.Database, nil
}

func parseGenArgs(args []string) (*CmdParams, error) {
	fs := flag.NewFlagSet("gentool gen", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	configPath := fs.String("c", "", "is path for gen.yml")
	dsn := fs.String("dsn", "", "consult[https://gorm.io/docs/connecting_to_the_database.html]")
	db := fs.String("db", string(dbMySQL), "input mysql|postgres|sqlite|sqlserver|clickhouse. consult[https://gorm.io/docs/connecting_to_the_database.html]")
	tableList := fs.String("tables", "", "enter the required data table or leave it blank")
	onlyModel := fs.Bool("onlyModel", false, "only generate models (without query file)")
	outPath := fs.String("outPath", defaultQueryPath, "specify a directory for output")
	outFile := fs.String("outFile", "", "query code file name, default: gen.go")
	withUnitTest := fs.Bool("withUnitTest", false, "generate unit test for query code")
	unitTestTemplate := fs.String("unitTestTemplate", "", "custom unit test template file path for query code")
	modelPkgName := fs.String("modelPkgName", "", "generated model code's package name")
	fieldNullable := fs.Bool("fieldNullable", false, "generate with pointer when field is nullable")
	fieldCoverable := fs.Bool("fieldCoverable", false, "generate with pointer when field has default value")
	fieldWithIndexTag := fs.Bool("fieldWithIndexTag", false, "generate field with gorm index tag")
	fieldWithTypeTag := fs.Bool("fieldWithTypeTag", false, "generate field with gorm column type tag")
	fieldWithDefaultTag := fs.Bool("fieldWithDefaultTag", false, "generate field with gorm default tag")
	fieldSignable := fs.Bool("fieldSignable", false, "detect integer field's unsigned type, adjust generated data type")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, flag.ErrHelp
		}
		return nil, err
	}

	if *configPath != "" {
		return parseCmdFromYaml(*configPath)
	}

	cmdParse := &CmdParams{
		DSN:                 strings.TrimSpace(*dsn),
		DB:                  strings.TrimSpace(*db),
		OnlyModel:           *onlyModel,
		OutPath:             strings.TrimSpace(*outPath),
		OutFile:             strings.TrimSpace(*outFile),
		WithUnitTest:        *withUnitTest,
		UnitTestTemplate:    strings.TrimSpace(*unitTestTemplate),
		ModelPkgName:        strings.TrimSpace(*modelPkgName),
		FieldNullable:       *fieldNullable,
		FieldCoverable:      *fieldCoverable,
		FieldWithIndexTag:   *fieldWithIndexTag,
		FieldWithTypeTag:    *fieldWithTypeTag,
		FieldWithDefaultTag: *fieldWithDefaultTag,
		FieldSignable:       *fieldSignable,
	}
	if strings.TrimSpace(*tableList) != "" {
		cmdParse.Tables = strings.Split(*tableList, ",")
	}
	return cmdParse, nil
}

func runGen(args []string) error {
	config, err := parseGenArgs(args)
	if err != nil {
		return err
	}
	config = config.revise()
	if config == nil {
		return fmt.Errorf("parse config fail")
	}

	db, err := connectDB(DBType(config.DB), config.DSN)
	if err != nil {
		return fmt.Errorf("connect db server fail: %w", err)
	}

	cfg := gen.Config{
		OutPath:           config.OutPath,
		OutFile:           config.OutFile,
		ModelPkgPath:      config.ModelPkgName,
		WithUnitTest:      config.WithUnitTest,
		FieldNullable:     config.FieldNullable,
		FieldCoverable:    config.FieldCoverable,
		FieldWithIndexTag: config.FieldWithIndexTag,
		FieldWithTypeTag:  config.FieldWithTypeTag,
		FieldSignable:     config.FieldSignable,
	}
	setConfigFieldIfExists(&cfg, "UnitTestTemplate", config.UnitTestTemplate)
	setConfigFieldIfExists(&cfg, "FieldWithDefaultTag", config.FieldWithDefaultTag)

	g := gen.NewGenerator(cfg)

	g.UseDB(db)

	models, err := genModels(g, db, config.Tables)
	if err != nil {
		return fmt.Errorf("get tables info fail: %w", err)
	}

	if !config.OnlyModel {
		g.ApplyBasic(models...)
	}

	g.Execute()
	return nil
}

func setConfigFieldIfExists(cfg *gen.Config, field string, value any) {
	v := reflect.ValueOf(cfg).Elem()
	f := v.FieldByName(field)
	if !f.IsValid() || !f.CanSet() {
		return
	}
	val := reflect.ValueOf(value)
	if !val.IsValid() {
		return
	}
	if val.Type().AssignableTo(f.Type()) {
		f.Set(val)
		return
	}
	if val.Type().ConvertibleTo(f.Type()) {
		f.Set(val.Convert(f.Type()))
	}
}

func runInit(args []string) error {
	fs := flag.NewFlagSet("gentool init", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	out := fs.String("o", defaultConfigFileName, "output path for config file")
	force := fs.Bool("f", false, "overwrite if file exists")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return flag.ErrHelp
		}
		return err
	}

	target := strings.TrimSpace(*out)
	if target == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	if !*force {
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("%s already exists (use -f to overwrite)", target)
		}
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil && filepath.Dir(target) != "." {
		return fmt.Errorf("create directory fail: %w", err)
	}

	if err := os.WriteFile(target, []byte(defaultConfigTemplate), 0o644); err != nil {
		return fmt.Errorf("write config fail: %w", err)
	}
	fmt.Fprintf(os.Stdout, "created %s\n", target)
	return nil
}

func runDoctor(args []string) error {
	config, err := parseGenArgs(args)
	if err != nil {
		return err
	}
	config = config.revise()
	if config == nil {
		return fmt.Errorf("parse config fail")
	}

	db, err := connectDB(DBType(config.DB), config.DSN)
	if err != nil {
		return fmt.Errorf("connect db server fail: %w", err)
	}

	allTables, err := db.Migrator().GetTables()
	if err != nil {
		return fmt.Errorf("get tables fail: %w", err)
	}

	tableSet := make(map[string]struct{}, len(allTables))
	for _, t := range allTables {
		tableSet[t] = struct{}{}
	}

	missing := make([]string, 0)
	for _, t := range config.Tables {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, ok := tableSet[t]; !ok {
			missing = append(missing, t)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("tables not found: %s", strings.Join(missing, ","))
	}

	outDir := strings.TrimSpace(config.OutPath)
	if outDir == "" {
		outDir = defaultQueryPath
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("ensure outPath fail: %w", err)
	}

	fmt.Fprintf(os.Stdout, "db=%s tables=%d outPath=%s\n", config.DB, len(allTables), outDir)
	if len(config.Tables) > 0 {
		fmt.Fprintf(os.Stdout, "selectedTables=%d\n", len(config.Tables))
	}
	return nil
}

func buildVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	for _, dep := range bi.Deps {
		if dep.Path == "gorm.io/gen/tools/gentool" && dep.Version != "" {
			return dep.Version
		}
	}
	if bi.Main.Path == "gorm.io/gen/tools/gentool" && bi.Main.Version != "" {
		return bi.Main.Version
	}
	if bi.Main.Version != "" {
		return bi.Main.Version
	}
	return "dev"
}

func printMainUsage(w io.Writer) {
	fmt.Fprintln(w, "gentool is a binary helper for gorm.io/gen")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  gentool [flags]              (same as: gentool gen [flags])")
	fmt.Fprintln(w, "  gentool gen [flags]          generate code")
	fmt.Fprintln(w, "  gentool init [-o path] [-f]  create a gen.yml template")
	fmt.Fprintln(w, "  gentool doctor [flags]       validate config and database connectivity")
	fmt.Fprintln(w, "  gentool version              print version")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Docs: https://gorm.io/gen/index.html")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		if err := runGen(args); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return
			}
			log.Fatalln(err)
		}
		return
	}

	switch args[0] {
	case "gen":
		if err := runGen(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return
			}
			log.Fatalln(err)
		}
	case "init":
		if err := runInit(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return
			}
			log.Fatalln(err)
		}
	case "doctor":
		if err := runDoctor(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return
			}
			log.Fatalln(err)
		}
	case "version":
		fmt.Fprintf(os.Stdout, "%s\n", buildVersion())
	case "help", "-h", "--help":
		printMainUsage(os.Stdout)
	default:
		printMainUsage(os.Stderr)
		log.Fatalln("unknown command:", args[0])
	}
}
