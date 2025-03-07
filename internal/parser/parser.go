package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// InterfaceSet ...
type InterfaceSet struct {
	Interfaces []InterfaceInfo
	imports    map[string]string // package name -> quoted "package path"
}

// InterfaceInfo ...
type InterfaceInfo struct {
	Name        string
	Doc         string
	Methods     []*Method
	Package     string
	ApplyStruct []string
}

// MatchStruct ...
func (i *InterfaceInfo) MatchStruct(name string) bool {
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
				return fmt.Errorf("file not found: %s", file)
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
	case *ast.ImportSpec:
		importName, _ := strconv.Unquote(n.Path.Value)
		importName = path.Base(importName)
		if n.Name != nil {
			name := n.Name.Name
			// ignore dummy imports
			// TODO: full support for dot imports requires type checking the whole package
			if name == "_" || name == "." {
				return i
			}
			importName = name
		}
		i.imports[importName] = n.Path.Value
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
					method := &Method{
						MethodName: name.Name,
						Doc:        m.Doc.Text(),
						Params:     getParamList(m.Type.(*ast.FuncType).Params),
						Result:     getParamList(m.Type.(*ast.FuncType).Results),
					}
					fixParamPackagePath(i.imports, method.Params)
					r.Methods = append(r.Methods, method)
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

	astResult := &InterfaceSet{imports: make(map[string]string)}
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

// Eq if param equal to another
func (p *Param) Eq(q Param) bool {
	return p.Package == q.Package && p.Type == q.Type
}

// IsError ...
func (p *Param) IsError() bool {
	return p.Type == "error"
}

// IsGenM ...
func (p *Param) IsGenM() bool {
	return p.Package == "gen" && p.Type == "M"
}

// IsGenRowsAffected ...
func (p *Param) IsGenRowsAffected() bool {
	return p.Package == "gen" && p.Type == "RowsAffected"
}

// IsMap ...
func (p *Param) IsMap() bool {
	return strings.HasPrefix(p.Type, "map[")
}

// IsGenT ...
func (p *Param) IsGenT() bool {
	return p.Package == "gen" && p.Type == "T"
}

// IsInterface ...
func (p *Param) IsInterface() bool {
	return p.Type == "interface{}"
}

// IsNull ...
func (p *Param) IsNull() bool {
	return p.Package == "" && p.Type == "" && p.Name == ""
}

// InMainPkg ...
func (p *Param) InMainPkg() bool {
	return p.Package == "main"
}

// IsTime ...
func (p *Param) IsTime() bool {
	return p.Package == "time" && p.Type == "Time"
}

// IsSQLResult ...
func (p *Param) IsSQLResult() bool {
	return (p.Package == "sql" && p.Type == "Result") || (p.Package == "gen" && p.Type == "SQLResult")
}

// IsSQLRow ...
func (p *Param) IsSQLRow() bool {
	return (p.Package == "sql" && p.Type == "Row") || (p.Package == "gen" && p.Type == "SQLRow")
}

// IsSQLRows ...
func (p *Param) IsSQLRows() bool {
	return (p.Package == "sql" && p.Type == "Rows") || (p.Package == "gen" && p.Type == "SQLRows")
}

// SetName ...
func (p *Param) SetName(name string) {
	p.Name = name
}

// TypeName ...
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

// IsBaseType judge whether the param type is basic type
func (p *Param) IsBaseType() bool {
	switch p.Type {
	case "string", "byte":
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
	case *ast.IndexExpr:
		p.astGetEltType(v.X)
	case *ast.IndexListExpr:
		p.astGetEltType(v.X)
	default:
		log.Printf("Unsupported param type: %+v", v)
	}
}

func (p *Param) astGetEltType(expr ast.Expr) {
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
	case *ast.ArrayType:
		p.astGetEltType(v.Elt)
		p.Type = "[]" + p.Type
	case *ast.IndexExpr:
		p.astGetEltType(v.X)
	default:
		log.Printf("Unsupported param type: %+v", v)
	}
}

func (p *Param) astGetPackageName(expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		p.Package = v.Name
	}
}

func (p *Param) astGetMapType(expr *ast.MapType) {
	p.Type = fmt.Sprintf("map[%s]%s", astGetType(expr.Key), astGetType(expr.Value))
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
