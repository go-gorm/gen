package generate

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gen/internal/parser"
)

//GetTestParamInTmpl return param list
func (m *InterfaceMethod) GetTestParamInTmpl() string {
	return testParamToString(m.Params)
}

// GetTestResultParamInTmpl return result list
func (m *InterfaceMethod) GetTestResultParamInTmpl() string {
	var res []string
	for i := range m.Result {
		tmplString := fmt.Sprintf("res%d", i+1)
		res = append(res, tmplString)
	}
	return strings.Join(res, ",")
}

// testParamToString param list to string used in tmpl
func testParamToString(params []parser.Param) string {
	var res []string
	for i, param := range params {
		// TODO manage array and pointer
		typ := param.Type
		if param.Package != "" {
			typ = param.Package + "." + typ
		}
		if param.IsArray {
			typ = "[]" + typ
		}
		if param.IsPointer {
			typ = "*" + typ
		}
		res = append(res, fmt.Sprintf("tt.Input.Args[%d].(%s)", i, typ))
	}
	return strings.Join(res, ",")
}

// GetAssertInTmpl assert in diy test
func (m *InterfaceMethod) GetAssertInTmpl() string {
	var res []string
	for i := range m.Result {
		tmplString := fmt.Sprintf("assert(t, %v, res%d, tt.Expectation.Ret[%d])", strconv.Quote(m.MethodName), i+1, i)
		res = append(res, tmplString)
	}
	return strings.Join(res, "\n")
}
