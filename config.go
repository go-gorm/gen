package gen

import (
	"fmt"
	"path/filepath"
	"strings"

	"gorm.io/gen/internal/check"
	"gorm.io/gen/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

type GenerateMode uint

const (
	// WithDefaultQuery create default query in generated code
	WithDefaultQuery GenerateMode = 1 << iota

	// WithoutContext generate code without context constrain
	WithoutContext
)

// Config generator's basic configuration
type Config struct {
	db *gorm.DB // db connection

	OutPath      string // query code path
	OutFile      string // query code file name, default: gen.go
	ModelPkgPath string // generated model code's package name
	WithUnitTest bool   // generate unit test for query code

	// generate model global configuration
	FieldNullable     bool // generate pointer when field is nullable
	FieldWithIndexTag bool // generate with gorm index tag
	FieldWithTypeTag  bool // generate with gorm column type tag

	Mode GenerateMode // generate mode

	queryPkgName string // generated query code's package name
	modelPkgPath string // model pkg path in target project
	dbNameOpts   []model.SchemaNameOpt

	// name strategy for syncing table from db
	tableNameNS func(tableName string) (targetTableName string)
	modelNameNS func(tableName string) (modelName string)
	fileNameNS  func(tableName string) (fielName string)

	dataTypeMap    map[string]func(detailType string) (dataType string)
	fieldJSONTagNS func(columnName string) (tagContent string)
	fieldNewTagNS  func(columnName string) (tagContent string)
}

// WithDbNameOpts set get database name function
func (cfg *Config) WithDbNameOpts(opts ...model.SchemaNameOpt) {
	if cfg.dbNameOpts == nil {
		cfg.dbNameOpts = opts
	} else {
		cfg.dbNameOpts = append(cfg.dbNameOpts, opts...)
	}
}

// WithTableNameStrategy specify table name naming strategy, only work when syncing table from db
func (cfg *Config) WithTableNameStrategy(ns func(tableName string) (targetTableName string)) {
	cfg.tableNameNS = ns
}

// WithModelNameStrategy specify model struct name naming strategy, only work when syncing table from db
func (cfg *Config) WithModelNameStrategy(ns func(tableName string) (modelName string)) {
	cfg.modelNameNS = ns
}

// WithFileNameStrategy specify file name naming strategy, only work when syncing table from db
func (cfg *Config) WithFileNameStrategy(ns func(tableName string) (fielName string)) {
	cfg.fileNameNS = ns
}

// WithDataTypeMap specify data type mapping relationship, only work when syncing table from db
func (cfg *Config) WithDataTypeMap(newMap map[string]func(detailType string) (dataType string)) {
	cfg.dataTypeMap = newMap
}

// WithJSONTagNameStrategy specify json tag naming strategy
func (cfg *Config) WithJSONTagNameStrategy(ns func(columnName string) (tagContent string)) {
	cfg.fieldJSONTagNS = ns
}

// WithNewTagNameStrategy specify new tag naming strategy
func (cfg *Config) WithNewTagNameStrategy(ns func(columnName string) (tagContent string)) {
	cfg.fieldNewTagNS = ns
}

// Revise format path and db
func (cfg *Config) Revise() (err error) {
	if strings.TrimSpace(cfg.ModelPkgPath) == "" {
		cfg.ModelPkgPath = check.DefaultModelPkg
	}

	cfg.OutPath, err = filepath.Abs(cfg.OutPath)
	if err != nil {
		return fmt.Errorf("outpath is invalid: %w", err)
	}
	if cfg.OutPath == "" {
		cfg.OutPath = "./query/"
	}
	if cfg.OutFile == "" {
		cfg.OutFile = cfg.OutPath + "/gen.go"
	}
	cfg.queryPkgName = filepath.Base(cfg.OutPath)

	if cfg.db == nil {
		cfg.db, _ = gorm.Open(tests.DummyDialector{})
	}

	return nil
}

func (cfg *Config) judgeMode(mode GenerateMode) bool { return cfg.Mode&mode != 0 }
