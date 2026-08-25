package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/build"
	goparser "go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// InterfacePath interface path
type InterfacePath struct {
	Name     string   // unqualified interface type name
	FullName string   // reflected package-qualified type name
	Files    []string // Go source files that may contain the interface declaration
	Package  string   // import path containing the interface
}

type goListPackage struct {
	Name    string   `json:"Name"`
	Import  string   `json:"ImportPath"`
	Dir     string   `json:"Dir"`
	GoFiles []string `json:"GoFiles"`
}

// GetInterfacePath get interface's directory path and all files it contains
func GetInterfacePath(v interface{}) (paths []*InterfacePath, err error) {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Func {
		err = fmt.Errorf("model param is not function:%s", value.String())
		return
	}

	callerDirs, err := interfaceCallerDirs()
	if err != nil {
		return nil, err
	}

	for i := 0; i < value.Type().NumIn(); i++ {
		var path InterfacePath
		arg := value.Type().In(i)
		path.FullName = arg.String()

		// keep the last model
		for _, n := range strings.Split(arg.String(), ".") {
			path.Name = n
		}

		p, loadErr := loadInterfacePackage(arg, path.Name, callerDirs)
		if loadErr != nil {
			return nil, loadErr
		}

		path.Package = p.Import
		for _, file := range p.GoFiles {
			goFile := file
			if !filepath.IsAbs(goFile) {
				goFile = filepath.Join(p.Dir, goFile)
			}
			if fileExists(goFile) {
				path.Files = append(path.Files, goFile)
			}
		}

		if len(path.Files) == 0 {
			err = fmt.Errorf("interface file not found:%s", value.String())
			return
		}

		paths = append(paths, &path)
	}

	return
}

func loadInterfacePackage(arg reflect.Type, interfaceName string, callerDirs []string) (*goListPackage, error) {
	pattern := arg.PkgPath()
	if isMainPackage(arg) {
		pattern = "."
	}
	if pattern == "" {
		return nil, fmt.Errorf("interface package not found:%s", arg.String())
	}

	var firstErr error
	for _, callerDir := range callerDirs {
		pkg, err := loadPackageInDir(pattern, callerDir)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if !isMainPackage(arg) || (pkg.Name == "main" && packageHasInterface(pkg.GoFiles, pkg.Dir, interfaceName)) {
			return pkg, nil
		}
	}

	if firstErr != nil && !isMainPackage(arg) {
		return nil, firstErr
	}
	return nil, fmt.Errorf("load interface package %s fail: interface %s not found", pattern, arg.String())
}

func loadPackageInDir(pattern, callerDir string) (*goListPackage, error) {
	cmd := exec.Command("go", "list", "-json", pattern)
	cmd.Dir = callerDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("load interface package %s fail:%w\n%s", pattern, err, bytes.TrimSpace(output))
	}

	var pkg goListPackage
	if err := json.Unmarshal(output, &pkg); err != nil {
		return nil, fmt.Errorf("parse go list output for %s fail:%w", pattern, err)
	}
	return &pkg, nil
}

func packageHasInterface(files []string, callerDir, interfaceName string) bool {
	for _, file := range files {
		goFile := file
		if !filepath.IsAbs(goFile) {
			goFile = filepath.Join(callerDir, goFile)
		}
		if fileHasInterface(goFile, interfaceName) {
			return true
		}
	}
	return false
}

func fileHasInterface(file, interfaceName string) bool {
	f, err := goparser.ParseFile(token.NewFileSet(), file, nil, 0)
	if err != nil {
		return false
	}
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != interfaceName {
				continue
			}
			if _, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				return true
			}
		}
	}
	return false
}

func isMainPackage(arg reflect.Type) bool {
	return strings.Split(arg.String(), ".")[0] == "main"
}

func interfaceCallerDir() (string, error) {
	callerDirs, err := interfaceCallerDirs()
	if err != nil {
		return "", err
	}
	return callerDirs[0], nil
}

func interfaceCallerDirs() ([]string, error) {
	_, parserFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("interface caller not found")
	}
	parserPackageDir := filepath.ToSlash(filepath.Dir(parserFile))

	var dirs []string
	seen := map[string]bool{}
	for skip := 1; ; skip++ {
		pc, file, _, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		if file == "" {
			continue
		}
		fn := ""
		if runtimeFn := runtime.FuncForPC(pc); runtimeFn != nil {
			fn = runtimeFn.Name()
		}
		if isGenInternalFrame(file, parserPackageDir, fn) {
			continue
		}
		dir := filepath.Dir(file)
		if !seen[dir] {
			dirs = append(dirs, dir)
			seen[dir] = true
		}
	}
	if len(dirs) > 0 {
		return dirs, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("interface caller not found:%w", err)
	}
	return []string{wd}, nil
}

func isGenInternalFrame(file, parserPackageDir, fn string) bool {
	if strings.HasPrefix(fn, "gorm.io/gen/internal/parser.") {
		return true
	}
	if strings.HasPrefix(fn, "gorm.io/gen.(*Generator).") {
		return true
	}
	return fn == "" && filepath.ToSlash(filepath.Dir(file)) == parserPackageDir
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetModelMethod get diy methods
func GetModelMethod(v interface{}) (method *DIYMethods, err error) {
	method = new(DIYMethods)

	// get diy method info by input value, must input a function or a struct
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Func:
		fullPath := runtime.FuncForPC(value.Pointer()).Name()
		err = method.parserPath(fullPath)
		if err != nil {
			return nil, err
		}
	case reflect.Struct:
		method.pkgPath = value.Type().PkgPath()
		method.BaseStructType = value.Type().Name()
	default:
		return nil, fmt.Errorf("method param must be a function or struct")
	}

	var p *build.Package

	// if struct in main file
	ctx := build.Default
	if method.pkgPath == "main" {
		var skip int
		var file string
		for {
			_, file, _, _ = runtime.Caller(skip)
			if !(strings.Contains(file, "gorm/gen/generator.go") || strings.Contains(file, "gorm/gen/internal")) || file == "" {
				break
			}
			skip++
		}
		p, err = ctx.ImportDir(filepath.Dir(file), build.ImportComment)
	} else {
		p, err = ctx.Import(method.pkgPath, "", build.ImportComment)
	}
	if err != nil {
		return nil, fmt.Errorf("diy method dir not found:%s.%s %w", method.pkgPath, method.MethodName, err)
	}

	for _, file := range p.GoFiles {
		goFile := p.Dir + "/" + file
		if fileExists(goFile) {
			method.pkgFiles = append(method.pkgFiles, goFile)
		}
	}
	if len(method.pkgFiles) == 0 {
		return nil, fmt.Errorf("diy method file not found:%s.%s", method.pkgPath, method.MethodName)
	}

	// read files got methods
	return method, method.LoadMethods()
}
