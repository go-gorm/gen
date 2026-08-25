package model

import (
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// Config model configuration
type Config struct {
	ModelPkg    string // package name for generated models
	TablePrefix string // prefix applied to normalized table names that do not already contain it
	TableName   string // source database table name
	ModelName   string // requested generated model type name

	ImportPkgPaths []string // additional quoted imports required by generated models
	ModelOpts      []Option // unsorted field and method options from the public configuration

	NameStrategy
	FieldConfig
	MethodConfig
}

// NameStrategy name strategy
type NameStrategy struct {
	SchemaNameOpts []SchemaNameOpt // ordered schema resolvers; the first non-empty result wins

	TableNameNS func(tableName string) string // maps the source table to its generated table name
	ModelNameNS func(tableName string) string // maps the source table to a Go model name
	FileNameNS  func(tableName string) string // maps the source table to a generated file name
}

// FieldConfig field configuration
type FieldConfig struct {
	DataTypeMap map[string]func(columnType gorm.ColumnType) (dataType string) // overrides database-to-Go type mappings

	FieldNullable       bool // generate pointer when field is nullable
	FieldCoverable      bool // generate pointer when field has default value
	FieldSignable       bool // detect integer field's unsigned type, adjust generated data type
	FieldWithIndexTag   bool // generate with gorm index tag
	FieldWithTypeTag    bool // generate with gorm column type tag
	FieldWithDefaultTag bool // includes database defaults in generated GORM tags

	FieldJSONTagNS func(columnName string) string // maps database columns to JSON tag values

	ModifyOpts []FieldOption // transforms existing fields after introspection
	FilterOpts []FieldOption // removes existing fields after introspection
	CreateOpts []FieldOption // adds synthetic fields after introspection
}

// MethodConfig method configuration
type MethodConfig struct {
	// MethodOpts contains custom model methods to copy into generated models.
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
	if tableName != "" && !strings.HasPrefix(tableName, cfg.TablePrefix) {
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
