package parser

import "go/ast"

func getParamList(fields *ast.FieldList) []Param {
	if fields == nil {
		return nil
	}
	var pars []Param
	if len(fields.List) < 1 {
		return nil
	}
	for _, field := range fields.List {
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

func fixParamPackagePath(imports map[string]string, params []Param) {
	for i := range params {
		if importPath, exist := imports[params[i].Package]; exist {
			params[i].PkgPath = importPath
		}
	}
}
