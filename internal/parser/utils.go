package parser

import "go/ast"

func getParamList(field *ast.FieldList) []Param {
	if field == nil {
		return nil
	}
	var pars []Param
	if len(field.List) < 1 {
		return nil
	}
	for _, field := range field.List {
		if field.Names == nil {
			par := Param{}
			par.astGetParamType(field)
			pars = append(pars, par)
			continue
		}
		for _, name := range field.Names {
			par := Param{
				Name: name.Name,
			}
			par.astGetParamType(field)
			pars = append(pars, par)
			continue
		}
	}
	return pars
}
