package generate

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/field"
	"gorm.io/gen/helper"
	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
)

// GetQueryStructMeta generate db model by table name
func GetQueryStructMeta(db *gorm.DB, conf *model.Config) (*QueryStructMeta, error) {
	if _, ok := db.Config.Dialector.(tests.DummyDialector); ok {
		return nil, fmt.Errorf("UseDB() is necessary to generate model struct [%s] from database table [%s]", conf.ModelName, conf.TableName)
	}

	conf = conf.Preprocess()
	tableName, structName, fileName := conf.GetNames()
	if tableName == "" {
		return nil, nil
	}
	if err := checkStructName(structName); err != nil {
		return nil, fmt.Errorf("model name %q is invalid: %w", structName, err)
	}

	columns, err := getTableColumns(db, conf.GetSchemaName(db), tableName, conf.FieldWithIndexTag)
	if err != nil {
		return nil, err
	}

	return (&QueryStructMeta{
		db:              db,
		Source:          model.Table,
		Generated:       true,
		FileName:        fileName,
		TableName:       tableName,
		ModelStructName: structName,
		QueryStructName: uncaptialize(structName),
		S:               strings.ToLower(structName[0:1]),
		StructInfo:      parser.Param{Type: structName, Package: conf.ModelPkg},
		ImportPkgPaths:  conf.ImportPkgPaths,
		Fields:          getFields(db, conf, columns),
	}).addMethodFromAddMethodOpt(conf.GetModelMethods()...), nil
}

// GetQueryStructMetaFromObject generate base struct from object
func GetQueryStructMetaFromObject(obj helper.Object, conf *model.Config) (*QueryStructMeta, error) {
	err := helper.CheckObject(obj)
	if err != nil {
		return nil, err
	}

	conf = conf.Preprocess()

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

	fields := make([]*model.Field, 0, 16)
	for _, field := range obj.Fields() {
		fields = append(fields, &model.Field{
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

	return &QueryStructMeta{
		Source:          model.Object,
		Generated:       true,
		FileName:        fileName,
		TableName:       tableName,
		ModelStructName: structName,
		QueryStructName: uncaptialize(structName),
		S:               strings.ToLower(structName[0:1]),
		StructInfo:      parser.Param{Type: structName, Package: conf.ModelPkg},
		ImportPkgPaths:  append(conf.ImportPkgPaths, obj.ImportPkgPaths()...),
		Fields:          fields,
	}, nil
}

// ConvertStructs convert to base structures
func ConvertStructs(db *gorm.DB, structs ...interface{}) (metas []*QueryStructMeta, err error) {
	for _, st := range structs {
		if isNil(st) {
			continue
		}
		if base, ok := st.(*QueryStructMeta); ok {
			metas = append(metas, base)
			continue
		}
		if !isStructType(reflect.ValueOf(st)) {
			return nil, fmt.Errorf("%s is not a struct", reflect.TypeOf(st).String())
		}

		structType := reflect.TypeOf(st)
		name := getStructName(structType.String())
		newStructName := name
		if st, ok := st.(interface{ GenInternalDoName() string }); ok {
			newStructName = st.GenInternalDoName()
		}

		meta := &QueryStructMeta{
			S:               getPureName(name),
			ModelStructName: name,
			QueryStructName: uncaptialize(newStructName),
			StructInfo:      parser.Param{PkgPath: structType.PkgPath(), Type: name, Package: getPackageName(structType.String())},
			Source:          model.Struct,
			db:              db,
		}
		if err := meta.parseStruct(st); err != nil {
			return nil, fmt.Errorf("transform struct [%s.%s] error:%s", meta.StructInfo.Package, name, err)
		}
		if err := meta.check(); err != nil {
			db.Logger.Warn(context.Background(), err.Error())
			continue
		}

		metas = append(metas, meta)
	}
	return
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	// if v is not ptr, return false(i is not nil)
	// if v is ptr, return v.IsNil()
	v := reflect.ValueOf(i)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

// BuildDIYMethod check the legitimacy of interfaces
func BuildDIYMethod(f *parser.InterfaceSet, s *QueryStructMeta, data []*InterfaceMethod) (checkResults []*InterfaceMethod, err error) {
	for _, interfaceInfo := range f.Interfaces {
		if interfaceInfo.MatchStruct(s.ModelStructName) {
			for _, method := range interfaceInfo.Methods {
				t := &InterfaceMethod{
					S:             s.S,
					TargetStruct:  s.QueryStructName,
					OriginStruct:  s.StructInfo,
					MethodName:    method.MethodName,
					Params:        method.Params,
					Doc:           method.Doc,
					Table:         s.TableName,
					InterfaceName: interfaceInfo.Name,
					Package:       getPackageName(interfaceInfo.Package),
				}
				if err = t.checkMethod(data, s); err != nil {
					return nil, err
				}
				if err = t.checkParams(method.Params); err != nil {
					return
				}
				if err = t.checkResult(method.Result); err != nil {
					return
				}
				if err = t.checkSQL(); err != nil {
					return
				}
				_, err = t.Section.BuildSQL()
				if err != nil {
					err = fmt.Errorf("sql [%s] build err:%w", t.SQLString, err)
					return
				}
				checkResults = append(checkResults, t)
			}
		}
	}
	return
}

// ParseStructRelationShip parse struct's relationship
// No one should use it directly in project
func ParseStructRelationShip(relationship *schema.Relationships) []field.Relation {
	cache := make(map[string]bool)
	return append(append(append(append(
		make([]field.Relation, 0, 4),
		pullRelationShip(cache, relationship.HasOne)...),
		pullRelationShip(cache, relationship.HasMany)...),
		pullRelationShip(cache, relationship.BelongsTo)...),
		pullRelationShip(cache, relationship.Many2Many)...,
	)
}

// GetStructNames get struct names from base structs
func GetStructNames(bases []*QueryStructMeta) (names []string) {
	for _, base := range bases {
		names = append(names, base.ModelStructName)
	}
	return names
}
