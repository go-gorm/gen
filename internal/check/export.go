package check

import (
	"context"
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
)

// ConvertStructs convert to base structures
func ConvertStructs(db *gorm.DB, structs ...interface{}) (bases []*BaseStruct, err error) {
	for _, st := range structs {
		if base, ok := st.(*BaseStruct); ok {
			bases = append(bases, base)
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

		base := &BaseStruct{
			S:             getPureName(name),
			StructName:    name,
			NewStructName: uncaptialize(newStructName),
			StructInfo:    parser.Param{PkgPath: structType.PkgPath(), Type: name, Package: getPackageName(structType.String())},
			Source:        model.Struct,
			db:            db,
		}
		if err := base.parseStruct(st); err != nil {
			return nil, fmt.Errorf("transform struct [%s.%s] error:%s", base.StructInfo.Package, name, err)
		}
		if err := base.check(); err != nil {
			db.Logger.Warn(context.Background(), err.Error())
			continue
		}

		bases = append(bases, base)
	}
	return
}

// BuildDIYMethod check the legitimacy of interfaces
func BuildDIYMethod(f *parser.InterfaceSet, s *BaseStruct, data []*InterfaceMethod) (checkResults []*InterfaceMethod, err error) {
	for _, interfaceInfo := range f.Interfaces {
		if interfaceInfo.MatchStruct(s.StructName) {
			for _, method := range interfaceInfo.Methods {
				t := &InterfaceMethod{
					S:             s.S,
					TargetStruct:  s.NewStructName,
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
