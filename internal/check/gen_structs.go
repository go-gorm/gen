package check

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/helper"
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

// GenBaseStruct generate db model by table name
func GenBaseStruct(db *gorm.DB, conf model.Conf) (base *BaseStruct, err error) {
	modelPkg := conf.ModelPkg
	tablePrefix := conf.TablePrefix
	tableName := conf.TableName
	structName := conf.ModelName

	if _, ok := db.Config.Dialector.(tests.DummyDialector); ok {
		return nil, fmt.Errorf("UseDB() is necessary to generate model struct [%s] from database table [%s]", structName, tableName)
	}

	if conf.ModelNameNS != nil {
		structName = conf.ModelNameNS(tableName)
	}
	if err = checkStructName(structName); err != nil {
		return nil, fmt.Errorf("model name %q is invalid: %w", structName, err)
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

	base = &BaseStruct{
		Source:         model.Table,
		GenBaseStruct:  true,
		FileName:       fileName,
		TableName:      tableName,
		StructName:     structName,
		NewStructName:  uncaptialize(structName),
		S:              strings.ToLower(structName[0:1]),
		StructInfo:     parser.Param{Type: structName, Package: modelPkg},
		ImportPkgPaths: conf.ImportPkgPaths,
	}

	modifyOpts, filterOpts, createOpts := conf.SortOpt()
	for _, col := range columns {
		col.SetDataTypeMap(conf.DataTypeMap)
		col.WithNS(conf.FieldJSONTagNS, conf.FieldNewTagNS)

		m := col.ToField(conf.FieldNullable, conf.FieldCoverable, conf.FieldSignable)

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
		} else if db.NamingStrategy != nil {
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

func checkStructName(name string) error {
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

func GenBaseStructFromObject(obj helper.Object, conf model.Conf) (*BaseStruct, error) {
	err := helper.CheckObject(obj)
	if err != nil {
		return nil, err
	}

	pkgName := filepath.Base(conf.ModelPkg)
	if pkgName == "" {
		pkgName = DefaultModelPkg
	}

	tableName := obj.TableName()
	if conf.TableNameNS != nil {
		tableName = conf.TableNameNS(tableName)
	}

	structName := obj.StructName()
	if conf.ModelNameNS != nil {
		structName = conf.ModelNameNS(structName)
	}

	fileName := obj.FileName()
	if fileName == "" {
		fileName = tableName
	}
	if fileName == "" {
		fileName = structName
	}
	if conf.FileNameNS != nil {
		fileName = conf.FileNameNS(fileName)
	} else {
		fileName = schema.NamingStrategy{SingularTable: true}.TableName(fileName)
	}

	base := &BaseStruct{
		Source:         model.Object,
		GenBaseStruct:  true,
		FileName:       fileName,
		TableName:      tableName,
		StructName:     structName,
		NewStructName:  uncaptialize(structName),
		S:              strings.ToLower(structName[0:1]),
		StructInfo:     parser.Param{Type: structName, Package: pkgName},
		ImportPkgPaths: append(conf.ImportPkgPaths, obj.ImportPkgPaths()...),
	}

	for _, field := range obj.Fields() {
		base.Fields = append(base.Fields, &model.Field{
			Name:             field.Name(),
			Type:             field.Type(),
			ColumnName:       field.ColumnName(),
			GORMTag:          field.GORMTag(),
			JSONTag:          field.JSONTag(),
			NewTag:           field.Tag(),
			ColumnComment:    field.Comment(),
			MultilineComment: strings.Contains(field.Comment(), "\n"),
		})
	}
	return base, nil
}
