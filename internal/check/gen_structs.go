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
func GenBaseStructs(db *gorm.DB, conf model.DBConf) (bases *BaseStruct, err error) {
	modelName, tableName := conf.ModelName, conf.TableName

	if _, ok := db.Config.Dialector.(tests.DummyDialector); ok {
		return nil, fmt.Errorf("UseDB() is necessary to generate model struct [%s] from database table [%s]", modelName, tableName)
	}

	if err = checkModelName(modelName); err != nil {
		return nil, fmt.Errorf("model name %q is invalid: %w", modelName, err)
	}

	modelPkg := conf.ModelPkg
	if modelPkg == "" {
		modelPkg = DefaultModelPkg
	}
	modelPkg = filepath.Base(modelPkg)

	columns, err := getTbColumns(db, conf.GetSchemaName(db), tableName, conf.FieldWithIndexTag)
	if err != nil {
		return nil, err
	}

	base := BaseStruct{
		Source:        model.TableName,
		GenBaseStruct: true,
		TableName:     tableName,
		StructName:    modelName,
		NewStructName: uncaptialize(modelName),
		S:             strings.ToLower(modelName[0:1]),
		StructInfo:    parser.Param{Type: modelName, Package: modelPkg},
	}

	modifyOpts, filterOpts, createOpts := conf.SortOpt()
	for _, field := range columns {
		field.SetDataTypeMap(conf.DataTypeMap)
		field.WithNS(conf.FieldJSONTagNS, conf.FieldNewTagNS)
		m := field.ToMember(conf.FieldNullable)

		if filterMember(m, filterOpts) == nil {
			continue
		}

		if !conf.FieldWithTypeTag { // remove type tag if FieldWithTypeTag == false
			m.GORMTag = strings.ReplaceAll(m.GORMTag, ";type:"+field.ColumnType, "")
		}

		m = modifyMember(m, modifyOpts)
		if ns, ok := db.NamingStrategy.(schema.NamingStrategy); ok {
			ns.SingularTable = true
			m.Name = ns.SchemaName(m.Name)
		} else {
			m.Name = db.NamingStrategy.SchemaName(m.Name)
		}

		base.Members = append(base.Members, m)
	}

	for _, create := range createOpts {
		m := create.Self()(nil)

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

		base.Members = append(base.Members, m)
	}

	return &base, nil
}

func filterMember(m *model.Member, opts []model.MemberOpt) *model.Member {
	for _, opt := range opts {
		if opt.Self()(m) == nil {
			return nil
		}
	}
	return m
}

func modifyMember(m *model.Member, opts []model.MemberOpt) *model.Member {
	for _, opt := range opts {
		m = opt.Self()(m)
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
