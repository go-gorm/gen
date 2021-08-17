package parser

import (
	"fmt"
	"go/build"
	"os"
	"reflect"
	"runtime"
	"strings"
)

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

		if strings.Split(arg.String(), ".")[0] == "main" {
			_, file, _, ok := runtime.Caller(3)
			if ok {
				path.Files = append(path.Files, file)
			}
			paths = append(paths, &path)
			continue
		}

		ctx := build.Default
		var p *build.Package
		p, err = ctx.Import(arg.PkgPath(), "", build.ImportComment)
		if err != nil {
			return
		}

		for _, file := range p.GoFiles {
			goFile := fmt.Sprintf("%s/%s", p.Dir, file)
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
