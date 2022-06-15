package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// InterfaceSet ...
type InterfaceSet struct {
	Interfaces []InterfaceInfo
}

// InterfaceInfo ...
type InterfaceInfo struct {
	Name        string
	Doc         string
	Methods     []*Method
	Package     string
	ApplyStruct []string
}

func (i *InterfaceInfo) IsMatchStruct(name string) bool {
	for _, s := range i.ApplyStruct {
		if s == name {
			return true
		}
	}
	return false
}

// ParseFile get interface's info from source file
func (i *InterfaceSet) ParseFile(paths []*InterfacePath, structNames []string) error {
	for _, path := range paths {
		for _, file := range path.Files {
			absFilePath, err := filepath.Abs(file)
			if err != nil {
				return fmt.Errorf("file not found：%s", file)
			}

			err = i.getInterfaceFromFile(absFilePath, path.Name, path.FullName, structNames)
			if err != nil {
				return fmt.Errorf("can't get interface from %s:%s", path.FullName, err)
			}
		}
	}
	return nil
}

// Visit ast visit function
func (i *InterfaceSet) Visit(n ast.Node) (w ast.Visitor) {
	switch n := n.(type) {
	case *ast.TypeSpec:
		if data, ok := n.Type.(*ast.InterfaceType); ok {
			r := InterfaceInfo{
				Methods: []*Method{},
			}
			methods := data.Methods.List
			r.Name = n.Name.Name
			r.Doc = n.Doc.Text()

			for _, m := range methods {
				for _, name := range m.Names {
					r.Methods = append(r.Methods, &Method{
						MethodName: name.Name,
						Doc:        m.Doc.Text(),
						Params:     getParamList(m.Type.(*ast.FuncType).Params),
						Result:     getParamList(m.Type.(*ast.FuncType).Results),
					})
				}

			}
			i.Interfaces = append(i.Interfaces, r)
		}
	}
	return i
}

// getInterfaceFromFile get interfaces
// get all interfaces from file and compare with specified name
func (i *InterfaceSet) getInterfaceFromFile(filename string, name, Package string, structNames []string) error {
	fileset := token.NewFileSet()
	f, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("can't parse file %q: %s", filename, err)
	}

	astResult := new(InterfaceSet)
	ast.Walk(astResult, f)

	for _, info := range astResult.Interfaces {
		if name == info.Name {
			info.Package = Package
			info.ApplyStruct = structNames
			i.Interfaces = append(i.Interfaces, info)
		}
	}

	return nil
}

// Param parameters in method
type Param struct { // (user model.User)
	PkgPath   string // package's path: internal/model
	Package   string // package's name: model
	Name      string // param's name: user
	Type      string // param's type: User
	IsArray   bool   // is array or not
	IsPointer bool   // is pointer or not
}

func (p *Param) Eq(q Param) bool {
	return p.Package == q.Package && p.Type == q.Type
}

func (p *Param) IsError() bool {
	return p.Type == "error"
}

func (p *Param) IsGenM() bool {
	return p.Package == "gen" && p.Type == "M"
}

func (p *Param) IsGenRowsAffected() bool {
	return p.Package == "gen" && p.Type == "RowsAffected"
}

func (p *Param) IsMap() bool {
	return strings.HasPrefix(p.Type, "map[")
}

func (p *Param) IsGenT() bool {
	return p.Package == "gen" && p.Type == "T"
}

func (p *Param) IsInterface() bool {
	return p.Type == "interface{}"
}

func (p *Param) IsNull() bool {
	return p.Package == "" && p.Type == "" && p.Name == ""
}

func (p *Param) InMainPkg() bool {
	return p.Package == "main"
}

func (p *Param) IsTime() bool {
	return p.Package == "time" && p.Type == "Time"
}

func (p *Param) SetName(name string) {
	p.Name = name
}

func (p *Param) TypeName() string {
	if p.IsArray {
		return "[]" + p.Type
	}
	return p.Type
}

// TmplString param to string in tmpl
func (p *Param) TmplString() string {
	var res strings.Builder
	if p.Name != "" {
		res.WriteString(p.Name)
		res.WriteString(" ")
	}

	if p.IsArray {
		res.WriteString("[]")
	}
	if p.IsPointer {
		res.WriteString("*")
	}
	if p.Package != "" {
		res.WriteString(p.Package)
		res.WriteString(".")
	}
	res.WriteString(p.Type)
	return res.String()
}

func (p *Param) AllowType() bool {
	switch p.Type {
	case "string", "bytes":
		return true
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	case "float64", "float32":
		return true
	case "bool":
		return true
	case "time.Time":
		return true
	default:
		return false
	}
}

func (p *Param) astGetParamType(param *ast.Field) {
	switch v := param.Type.(type) {
	case *ast.Ident:
		p.Type = v.Name
		if v.Obj != nil {
			p.Package = "UNDEFINED" // set a placeholder
		}
	case *ast.SelectorExpr:
		p.astGetEltType(v)
	case *ast.ArrayType:
		p.astGetEltType(v.Elt)
		p.IsArray = true
	case *ast.Ellipsis:
		p.astGetEltType(v.Elt)
		p.IsArray = true
	case *ast.MapType:
		p.astGetMapType(v)
	case *ast.InterfaceType:
		p.Type = "interface{}"
	case *ast.StarExpr:
		p.IsPointer = true
		p.astGetEltType(v.X)
	default:
		log.Fatalf("unknow param type: %+v", v)
	}
}

func (p *Param) astGetEltType(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		p.Type = v.Name
		if v.Obj != nil {
			p.Package = "UNDEFINED"
		}
	case *ast.SelectorExpr:
		p.Type = v.Sel.Name
		p.astGetPackageName(v.X)
	case *ast.MapType:
		p.astGetMapType(v)
	case *ast.StarExpr:
		p.IsPointer = true
		p.astGetEltType(v.X)
	case *ast.InterfaceType:
		p.Type = "interface{}"
	default:
		log.Fatalf("unknow param type: %+v", v)
	}
	return p.Type
}

func (p *Param) astGetPackageName(expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		p.Package = v.Name
	}

}

func (p *Param) astGetMapType(expr *ast.MapType) string {
	p.Type = fmt.Sprintf("map[%s]%s", astGetType(expr.Key), astGetType(expr.Value))
	return ""
}

func astGetType(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.InterfaceType:
		return "interface{}"
	}
	return ""

}

// Method Apply to query struct and base struct custom method
type Method struct {
	BaseStruct Param
	MethodName string
	Doc        string
	Params     []Param
	Result     []Param
	Body       string
}

//GetParamInTmpl return method bind info string
func (m *Method) GetBaseStructTmpl() string {
	return m.BaseStruct.TmplString()
}

//GetParamInTmpl return param list
func (m *Method) GetParamInTmpl() string {
	return paramToString(m.Params)
}

// GetResultParamInTmpl return result list
func (m *Method) GetResultParamInTmpl() string {
	return paramToString(m.Result)
}

// paramToString param list to string used in tmpl
func paramToString(params []Param) string {
	var res []string
	for _, param := range params {
		res = append(res, param.TmplString())
	}
	return strings.Join(res, ",")
}

// DocComment return comment sql add "//" every line
func (m *Method) DocComment() string {
	return strings.Replace(strings.TrimSpace(m.Doc), "\n", "\n//", -1)
}

// CustomMethods user Custom methods bind to db base struct
type CustomMethods struct {
	BaseStructType string
	MethodName     string
	pkgPath        string
	currentFile    string
	pkgFiles       []string
	Methods        []*Method
}

func (m *CustomMethods) parserPath(path string) error {
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
func (m *CustomMethods) Visit(n ast.Node) (w ast.Visitor) {
	switch t := n.(type) {
	case *ast.FuncDecl:
		// check base struct and method name is expect
		baseStruct := getParamList(t.Recv)
		if len(baseStruct) != 1 {
			return
		}
		if baseStruct[0].Type != m.BaseStructType {
			return
		}
		// if m.MethodName is null will generate all methods
		if m.MethodName != "" && m.MethodName != t.Name.Name {
			return
		}

		// use ast read bind start package is UNDEFINED ,set it null string
		baseStruct[0].Package = ""
		m.Methods = append(m.Methods, &Method{
			BaseStruct: baseStruct[0],
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

	return string(f1[start-1 : end])
}

// LoadMethods ast read file get diy method
func (m *CustomMethods) LoadMethods() error {
	for _, filename := range m.pkgFiles {
		fileset := token.NewFileSet()
		f, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("can't parse file %q: %s", filename, err)
		}
		m.currentFile = filename
		ast.Walk(m, f)
	}

	return nil
}
