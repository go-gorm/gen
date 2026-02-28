package parser

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"golang.org/x/tools/go/packages"
)

// InterfacePath interface path
type InterfacePath struct {
	Name     string
	FullName string
	Files    []string
	Package  string
}

// GetInterfacePath get interface's directory path and all files it contains
func GetInterfacePath(v interface{}) (paths []*InterfacePath, err error) {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Func {
		err = fmt.Errorf("model param is not function:%s", value.String())
		return
	}

	for i := 0; i < value.Type().NumIn(); i++ {
		var path InterfacePath
		arg := value.Type().In(i)
		path.FullName = arg.String()

		// keep the last model
		for _, n := range strings.Split(arg.String(), ".") {
			path.Name = n
		}

		cfg := &packages.Config{Mode: packages.NeedFiles}
		var pkgs []*packages.Package
		if strings.Split(arg.String(), ".")[0] == "main" {
			var skip int
			var file string
			for {
				_, file, _, _ = runtime.Caller(skip)
				if !(strings.Contains(file, "gorm/gen/generator.go") || strings.Contains(file, "gorm/gen/internal")) || file == "" {
					break
				}
				skip++
			}
			cfg.Dir = filepath.Dir(file)
			pkgs, err = packages.Load(cfg, ".")
		} else {
			pkgs, err = packages.Load(cfg, arg.PkgPath())
		}
		if err != nil {
			return nil, err
		}
		if len(pkgs) == 0 {
			return nil, fmt.Errorf("interface package not found:%s", arg.PkgPath())
		}
		if len(pkgs[0].Errors) > 0 {
			return nil, fmt.Errorf("load package %s fail: %w", arg.PkgPath(), pkgs[0].Errors[0])
		}
		path.Package = pkgs[0].PkgPath
		for _, file := range pkgs[0].GoFiles {
			if fileExists(file) {
				path.Files = append(path.Files, file)
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
