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

	NameStrategy
	FieldConfig
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

	FieldOpts []FieldOpt
}

// Revise revise invalid field
func (cfg *Config) Revise() *Config {
	if cfg.ModelPkg == "" {
		cfg.ModelPkg = DefaultModelPkg
	}
	cfg.ModelPkg = filepath.Base(cfg.ModelPkg)
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

// SortOpt sort option
func (cfg *Config) SortOpt() (modifyOpts []FieldOpt, filterOpts []FieldOpt, createOpts []FieldOpt) {
	if cfg == nil {
		return
	}
	return sortFieldOpt(cfg.FieldOpts)
}

// GetSchemaName get schema name
func (cfg *Config) GetSchemaName(db *gorm.DB) string {
	if cfg == nil {
		return defaultSchemaNameOpt(db)
	}
	for _, opt := range cfg.SchemaNameOpts {
		if name := opt(db); name != "" {
			return name
		}
	}
	return defaultSchemaNameOpt(db)
}
