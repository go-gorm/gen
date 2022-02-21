package check

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
)

/*
** The feature of mapping table from database server to Golang struct
** Provided by @qqxhb
 */

const (
	DefaultModelPkg = "model"
)

// GenBaseStructs generate db model by table name
func GenBaseStructs(db *gorm.DB, conf model.Conf) (bases *BaseStruct, err error) {
	modelPkg := conf.ModelPkg
	tablePrefix := conf.TablePrefix
	tableName := conf.TableName
	modelName := conf.ModelName

	if _, ok := db.Config.Dialector.(tests.DummyDialector); ok {
		return nil, fmt.Errorf("UseDB() is necessary to generate model struct [%s] from database table [%s]", modelName, tableName)
	}

	if conf.ModelNameNS != nil {
		modelName = conf.ModelNameNS(tableName)
	}
	if err = checkModelName(modelName); err != nil {
		return nil, fmt.Errorf("model name %q is invalid: %w", modelName, err)
	}

	if conf.TableNameNS != nil {
		tableName = conf.TableNameNS(tableName)
	}
	if !strings.HasPrefix(tableName, tablePrefix) {
		tableName = tablePrefix + tableName
	}

	fileName := strings.ToLower(tableName)
	if conf.FileNameNS != nil {
		fileName = conf.FileNameNS(conf.TableName)
	}

	columns, err := getTblColumns(db, conf.GetSchemaName(db), tableName, conf.FieldWithIndexTag)
	if err != nil {
		return nil, err
	}

	if modelPkg == "" {
		modelPkg = DefaultModelPkg
	}
	modelPkg = filepath.Base(modelPkg)

	base := &BaseStruct{
		Source:         model.Table,
		GenBaseStruct:  true,
		FileName:       fileName,
		TableName:      tableName,
		StructName:     modelName,
		NewStructName:  uncaptialize(modelName),
		S:              strings.ToLower(modelName[0:1]),
		StructInfo:     parser.Param{Type: modelName, Package: modelPkg},
		ImportPkgPaths: conf.ImportPkgPaths,
	}

	modifyOpts, filterOpts, createOpts := conf.SortOpt()
	for _, col := range columns {
		col.SetDataTypeMap(conf.DataTypeMap)
		col.WithNS(conf.FieldJSONTagNS, conf.FieldNewTagNS)

		m := col.ToField(conf.FieldNullable, conf.FieldCoverable)

		if filterField(m, filterOpts) == nil {
			continue
		}
		if t, ok := col.ColumnType.ColumnType(); ok && !conf.FieldWithTypeTag { // remove type tag if FieldWithTypeTag == false
			m.GORMTag = strings.ReplaceAll(m.GORMTag, ";type:"+t, "")
		}

		m = modifyField(m, modifyOpts)
		if ns, ok := db.NamingStrategy.(schema.NamingStrategy); ok {
			ns.SingularTable = true
			m.Name = ns.SchemaName(ns.TablePrefix + m.Name)
		} else {
			m.Name = db.NamingStrategy.SchemaName(m.Name)
		}

		base.Fields = append(base.Fields, m)
	}

	for _, create := range createOpts {
		m := create.Operator()(nil)

		if m.Relation != nil {
			if m.Relation.Model() != nil {
				stmt := gorm.Statement{DB: db}
				_ = stmt.Parse(m.Relation.Model())
				if stmt.Schema != nil {
					m.Relation.AppendChildRelation(ParseStructRelationShip(&stmt.Schema.Relationships)...)
				}
			}
			m.Type = strings.ReplaceAll(m.Type, modelPkg+".", "") // remove modelPkg in field's Type, avoid import error
		}

		base.Fields = append(base.Fields, m)
	}

	return base, nil
}

func filterField(m *model.Field, opts []model.FieldOpt) *model.Field {
	for _, opt := range opts {
		if opt.Operator()(m) == nil {
			return nil
		}
	}
	return m
}

func modifyField(m *model.Field, opts []model.FieldOpt) *model.Field {
	for _, opt := range opts {
		m = opt.Operator()(m)
	}
	return m
}

// get mysql db' name
var modelNameReg = regexp.MustCompile(`^\w+$`)

func checkModelName(name string) error {
	if name == "" {
		return nil
	}
	if !modelNameReg.MatchString(name) {
		return fmt.Errorf("model name cannot contains invalid character")
	}
	if name[0] < 'A' || name[0] > 'Z' {
		return fmt.Errorf("model name must be initial capital")
	}
	return nil
}
