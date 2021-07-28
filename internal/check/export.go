package check

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"gorm.io/gen/internal/parser"
)

// CheckStructs check the legitimacy of structures
func CheckStructs(db *gorm.DB, structs ...interface{}) (bases []*BaseStruct, err error) {
	if isDBUndefined(db) {
		return nil, fmt.Errorf("gen config db is undefined")
	}

	for _, st := range structs {
		structType := reflect.TypeOf(st)
		name := getTypeName(structType.String())
		base := &BaseStruct{
			S:             GetSimpleName(name),
			StructName:    name,
			NewStructName: strings.ToLower(name),
			StructInfo:    parser.Param{Type: name, Package: getPackageName(structType.String())},
			db:            db,
		}
		base.getMembers(st)
		base.getTableName(st)
		if e := base.check(); e != nil {
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
