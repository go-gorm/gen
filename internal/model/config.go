package model

import (
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// Config model configuration
type Config struct {
	ModelPkg  string // package name for generated models
	TableName string // source database table name
	ModelName string // requested generated model type name

	// TablePrefix prefixes normalized table names that do not already start
	// with it. The prefix mirrors gorm's schema.NamingStrategy.TablePrefix so
	// generated TableName() values match the runtime table mapping.
	//
	// The "already starts with it" check follows gorm's naming convention but
	// cannot tell an intentionally prefixed table from a coincidental one:
	// with prefix "user" a source table named "users" is kept as-is instead
	// of becoming "userusers", and with prefix "t_" the source tables
	// "users" and "t_users" both normalize to "t_users". The generator
	// therefore rejects such collisions when models are registered instead of
	// silently emitting code that queries the wrong physical table; use
	// WithTableNameStrategy (TableNameNS) to disambiguate.
	TablePrefix string

	ImportPkgPaths []string // additional quoted imports required by generated models
	ModelOpts      []Option // unsorted field and method options from the public configuration

	NameStrategy
	FieldConfig
	MethodConfig
	TemplateConfig
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

// TemplateConfig custom template configuration
type TemplateConfig struct {
	TemplateOpts []TemplateOption
}

// Preprocess revise invalid field
func (cfg *Config) Preprocess() *Config {
	if cfg.ModelPkg == "" {
		cfg.ModelPkg = DefaultModelPkg
	}
	cfg.ModelPkg = filepath.Base(cfg.ModelPkg)

	cfg.ModifyOpts, cfg.FilterOpts, cfg.CreateOpts, cfg.MethodOpts, cfg.TemplateOpts = sortOptions(cfg.ModelOpts)

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
	// Skip only when the (already normalized) table name starts with the
	// prefix, matching gorm's naming convention for pre-prefixed tables. An
	// empty prefix never applies, and coincidental prefixes (prefix "user",
	// table "users") are indistinguishable from intentional ones at this
	// layer; the generator detects the resulting cross-table collisions at
	// registration time.
	if tableName != "" && cfg.TablePrefix != "" && !strings.HasPrefix(tableName, cfg.TablePrefix) {
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

// GetCustomTemplates get diy custom template from option
func (cfg *Config) GetCustomTemplates() (templates []string) {
	if cfg == nil {
		return
	}

	for _, opt := range cfg.TemplateOpts {
		templates = append(templates, opt.Templates()...)
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
