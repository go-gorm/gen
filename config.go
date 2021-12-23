package gen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
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
	ProjectPath string // project root path, eg: /User/xxxProjectPath/dal/model
	OutPath      string // query code path
	OutFile      string // query code file name, default: gen.go
	ModelPkgPath string // generated model code's package name
	WithUnitTest bool   // generate unit test for query code

	// generate model global configuration
	FieldNullable     bool // generate pointer when field is nullable
	FieldWithIndexTag bool // generate with gorm index tag
	FieldWithTypeTag  bool // generate with gorm column type tag

	Mode GenerateMode // generate mode

	queryPkgName    string // generated query code's package name
	modelImportPath string
	dbNameOpts      []model.SchemaNameOpt

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

var moduleFullPath = func() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Path
}()

func (cfg *Config) Revise() (err error) {
	moduleName :=  getProjectModuleName(cfg.ProjectPath)
	if moduleName == "" {
		moduleName = "undefined"
	}
	if strings.TrimSpace(cfg.ModelPkgPath) == "" {
		cfg.OutPath, err = filepath.Abs(cfg.OutPath)
		if err != nil {
			fmt.Println("get abs of outPath error:", err.Error())
			return err
		}

		cfg.ModelPkgPath = check.DefaultModelPkg
		relPathOfOutPath, err := filepath.Rel(cfg.ProjectPath, cfg.OutPath)
		if err != nil {
			fmt.Println("get relPathOfOutPath error:", err.Error())
			return err
		}
		cfg.modelImportPath = filepath.Dir(filepath.Clean(moduleName+"/" + relPathOfOutPath)) + "/" + cfg.ModelPkgPath
	} else {
		cfg.ModelPkgPath, err = filepath.Abs(cfg.ModelPkgPath)
		if err != nil {
			return fmt.Errorf("get abs of ModelPkgPath failed: %v", err)
		}
		relPathOfModelPath, err := filepath.Rel(cfg.ProjectPath, cfg.ModelPkgPath)
		if err != nil {
			return fmt.Errorf("get  relPathOfModelPath failed: %v", err)
		}
		cfg.modelImportPath = filepath.Clean(moduleName + "/" + relPathOfModelPath + cfg.ModelPkgPath)
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

func getProjectModuleName(projectPath string) string {
	goModeFilePath := projectPath + "/go.mod"
	// check go.mod file
	_, err := os.Stat(goModeFilePath)
	if os.IsNotExist(err) {
		panic("go mod file not found" + err.Error())
	}
	firstLine := strings.TrimSpace(readLine(goModeFilePath, 1))
	if strings.HasPrefix(firstLine, "module") {
		return strings.TrimSpace(strings.ReplaceAll(firstLine, "module", ""))
	}
	return ""
}

func readLine(filePath string, lineNumber int) string {
	file, _ := os.Open(filePath)
	fileScanner := bufio.NewScanner(file)
	lineCount := 1
	for fileScanner.Scan() {
		if lineCount == lineNumber {
			return fileScanner.Text()
		}
		lineCount++
	}
	defer file.Close()
	return ""
}