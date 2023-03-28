package gen

import (
	"encoding/json"

	"gorm.io/gen/internal/model"
	"gorm.io/gen/internal/parser"
)

type Method struct {
	Receiver   Param
	MethodName string
	Doc        string
	Params     []Param
	Result     []Param
	Body       string
}

type Param struct { // (user model.User)
	PkgPath   string // package's path: internal/model
	Package   string // package's name: model
	Name      string // param's name: user
	Type      string // param's type: User
	IsArray   bool   // is array or not
	IsPointer bool   // is pointer or not
}

var (
	WithDIYMethod = func(method *Method) model.AddMethodOpt {
		return func() []interface{} {
			if method == nil {
				return []interface{}{}
			}
			var dst parser.Method
			{
				src, err := json.Marshal(method)
				if err != nil {
					return []interface{}{}
				}
				err = json.Unmarshal(src, &dst)
				if err != nil {
					return []interface{}{}
				}
			}
			return []interface{}{&dst}
		}
	}
)
