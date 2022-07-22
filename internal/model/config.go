package model

import (
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// Config model configuration
type Config struct {
	ModelPkg    string
	TablePrefix string
	TableName   string
	ModelName   string

	ImportPkgPaths []string
	ModelOpts      []Option

	NameStrategy
	FieldConfig
	MethodConfig
}

// NameStrategy name strategy
type NameStrategy struct {
	SchemaNameOpts []SchemaNameOpt

	TableNameNS func(tableName string) string
	ModelNameNS func(tableName string) string
	FileNameNS  func(tableName string) string
}

// FieldConfig field configuration
type FieldConfig struct {
	DataTypeMap map[string]func(detailType string) (dataType string)

	FieldNullable     bool // generate pointer when field is nullable
	FieldCoverable    bool // generate pointer when field has default value
	FieldSignable     bool // detect integer field's unsigned type, adjust generated data type
	FieldWithIndexTag bool // generate with gorm index tag
	FieldWithTypeTag  bool // generate with gorm column type tag

	FieldJSONTagNS func(columnName string) string
	FieldNewTagNS  func(columnName string) string

	ModifyOpts []FieldOption
	FilterOpts []FieldOption
	CreateOpts []FieldOption
}

// MethodConfig method configuration
type MethodConfig struct {
	MethodOpts []MethodOption
}

// Preprocess revise invalid field
func (cfg *Config) Preprocess() *Config {
	if cfg.ModelPkg == "" {
		cfg.ModelPkg = DefaultModelPkg
	}
	cfg.ModelPkg = filepath.Base(cfg.ModelPkg)

	cfg.ModifyOpts, cfg.FilterOpts, cfg.CreateOpts, cfg.MethodOpts = sortOptions(cfg.ModelOpts)

	return cfg
}

// GetNames get names
func (cfg *Config) GetNames() (tableName, structName, fileName string) {
	tableName, structName = cfg.TableName, cfg.ModelName

	if cfg.ModelNameNS != nil {
		structName = cfg.ModelNameNS(tableName)
	}

	if cfg.TableNameNS != nil {
		tableName = cfg.TableNameNS(tableName)
	}
	if !strings.HasPrefix(tableName, cfg.TablePrefix) {
		tableName = cfg.TablePrefix + tableName
	}

	fileName = strings.ToLower(tableName)
	if cfg.FileNameNS != nil {
		fileName = cfg.FileNameNS(cfg.TableName)
	}

	return
}

// GetModelMethods get diy method from option
func (cfg *Config) GetModelMethods() (methods []interface{}) {
	if cfg == nil {
		return
	}

	for _, opt := range cfg.MethodOpts {
		methods = append(methods, opt.Methods()...)
	}
	return
}

// GetSchemaName get schema name
func (cfg *Config) GetSchemaName(db *gorm.DB) string {
	if cfg == nil {
		return ""
	}

	for _, opt := range cfg.SchemaNameOpts {
		if name := opt(db); name != "" {
			return name
		}
	}
	return ""
}
