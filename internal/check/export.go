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
		base := &BaseStruct{
			S:             GetSimpleName(structType.Name()),
			StructName:    structType.Name(),
			NewStructName: strings.ToLower(structType.Name()),
			StructInfo:    parser.Param{Type: structType.Name(), Package: getPackageName(structType.String())},
			db:            db,
		}
		base.getMembers(st)
		base.getTableName(st)
		if e := base.checkStructAndMembers(); e != nil {
			continue
		}

		bases = append(bases, base)
	}
	return
}

// CheckInterface check the legitimacy of interfaces
func CheckInterface(f *parser.InterfaceSet, s *BaseStruct) (checkResult []*InterfaceMethod, err error) {
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
			err = t.checkParams(method.Params)
			if err != nil {
				return
			}
			err = t.checkResult(method.Result)
			if err != nil {
				return
			}
			err = t.checkSQL()
			if err != nil {
				return
			}
			checkResult = append(checkResult, t)
		}
	}
	return
}
