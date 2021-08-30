package check

import (
	"context"
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"gorm.io/gen/internal/parser"
)

// CheckStructs check the legitimacy of structures
func CheckStructs(db *gorm.DB, structs ...interface{}) (bases []*BaseStruct, err error) {
	if isDBUnset(db) {
		return nil, fmt.Errorf("gen config db is undefined")
	}

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
		base := &BaseStruct{
			S:             GetSimpleName(name),
			StructName:    name,
			NewStructName: uncaptialize(name),
			StructInfo:    parser.Param{Type: name, Package: getPackageName(structType.String())},
			Source:        Struct,
			db:            db,
		}
		if err := base.transStruct(st); err != nil {
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

// CheckInterface check the legitimacy of interfaces
func CheckInterface(f *parser.InterfaceSet, s *BaseStruct) (checkResults []*InterfaceMethod, err error) {
	for _, interfaceInfo := range f.Interfaces {
		for _, method := range interfaceInfo.Methods {
			t := &InterfaceMethod{
				S:             s.S,
				MethodStruct:  s.NewStructName,
				OriginStruct:  s.StructInfo,
				MethodName:    method.MethodName,
				Params:        method.Params,
				Doc:           method.Doc,
				ExecuteResult: "_",
				Table:         s.TableName,
				InterfaceName: interfaceInfo.Name,
				Package:       getPackageName(interfaceInfo.Package),
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
			checkResults = append(checkResults, t)
		}
	}
	return
}
