package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func DefaultMethodTableName(structName string) *Method {
	return &Method{
		Receiver:   Param{IsPointer: true, Type: structName},
		MethodName: "TableName",
		Doc:        fmt.Sprint("TableName ", structName, "'s table name "),
		Result:     []Param{{Type: "string"}},
		Body:       fmt.Sprintf("{\n\treturn TableName%s\n} ", structName),
	}
}

// Method Apply to query struct and base struct custom method
type Method struct {
	Receiver   Param
	MethodName string
	Doc        string
	Params     []Param
	Result     []Param
	Body       string
}

// FuncSign function signature
func (m Method) FuncSign() string {
	return fmt.Sprintf("%s(%s) (%s)", m.MethodName, m.GetParamInTmpl(), m.GetResultParamInTmpl())
}

// GetBaseStructTmpl return method bind info string
func (m *Method) GetBaseStructTmpl() string {
	return m.Receiver.TmplString()
}

// GetParamInTmpl return param list
func (m *Method) GetParamInTmpl() string {
	return paramToString(m.Params)
}

// GetResultParamInTmpl return result list
func (m *Method) GetResultParamInTmpl() string {
	return paramToString(m.Result)
}

// paramToString param list to string used in tmpl
func paramToString(params []Param) string {
	res := make([]string, len(params))
	for i, param := range params {
		res[i] = param.TmplString()
	}
	return strings.Join(res, ",")
}

// DocComment return comment sql add "//" every line
func (m *Method) DocComment() string {
	return strings.Replace(strings.TrimSpace(m.Doc), "\n", "\n//", -1)
}

// DIYMethods user Custom methods bind to db base struct
type DIYMethods struct {
	BaseStructType string
	MethodName     string
	pkgPath        string
	currentFile    string
	pkgFiles       []string
	Methods        []*Method
}

func (m *DIYMethods) parserPath(path string) error {
	pathList := strings.Split(path, ".")
	if len(pathList) < 3 {
		return fmt.Errorf("parser diy method error")
	}

	m.pkgPath = strings.Join(pathList[:len(pathList)-2], ".")
	methodName := pathList[len(pathList)-1]
	m.MethodName = methodName[:len(methodName)-3]

	structName := pathList[len(pathList)-2]
	m.BaseStructType = strings.Trim(structName, "()*")
	return nil
}

// Visit ast visit function
func (m *DIYMethods) Visit(n ast.Node) (w ast.Visitor) {
	switch t := n.(type) {
	case *ast.FuncDecl:
		// check base struct and method name is expect
		structMeta := getParamList(t.Recv)
		if len(structMeta) != 1 {
			return
		}
		if structMeta[0].Type != m.BaseStructType {
			return
		}
		// if m.MethodName is null will generate all methods
		if m.MethodName != "" && m.MethodName != t.Name.Name {
			return
		}

		// use ast read bind start package is UNDEFINED ,set it null string
		structMeta[0].Package = ""
		m.Methods = append(m.Methods, &Method{
			Receiver:   structMeta[0],
			MethodName: t.Name.String(),
			Doc:        t.Doc.Text(),
			Body:       getBody(m.currentFile, int(t.Body.Pos()), int(t.Body.End())),
			Params:     getParamList(t.Type.Params),
			Result:     getParamList(t.Type.Results),
		})
	}

	return m
}

// read old file get method body
func getBody(fileName string, start, end int) string {
	f1, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "{}"
	}

	return string(f1[start-1 : end-1])
}

// LoadMethods ast read file get diy method
func (m *DIYMethods) LoadMethods() error {
	for _, filename := range m.pkgFiles {
		f, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("can't parse file %q: %s", filename, err)
		}
		m.currentFile = filename
		ast.Walk(m, f)
	}

	return nil
}
