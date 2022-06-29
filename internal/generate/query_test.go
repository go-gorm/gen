package generate

import (
	"reflect"
	"testing"

	"gorm.io/gen/internal/parser"
)

func checkCustomMethod(t *testing.T, expMethods []*parser.Method, methods ...interface{}) {
	base := &QueryStructMeta{
		GenBaseStruct: true,
		FileName:      "users",
		S:             "u",
		NewStructName: "user",
		StructName:    "User",
		TableName:     "users",
		StructInfo: parser.Param{
			PkgPath: "",
			Package: "model",
			Name:    "",
			Type:    "User",
		},
	}
	base.AddMethod(methods...)
	genMethods := base.DIYMethods
	if len(expMethods) != len(genMethods) {
		t.Errorf("custom method length unexpectate exp:%v got:%v", len(expMethods), len(genMethods))
	}

	for _, expMethod := range expMethods {
		pass := false
		for _, genMethod := range genMethods {
			if genMethod.MethodName == expMethod.MethodName {
				switch {
				case !reflect.DeepEqual(genMethod.Receiver, expMethod.Receiver):
					t.Errorf("gen custom method BaseStruct unexpectate \nexp:%v \ngot:%v\n", expMethod.Receiver, genMethod.Receiver)
				case !reflect.DeepEqual(genMethod.Doc, expMethod.Doc):
					t.Errorf("gen custom method Doc unexpectate \nexp:%v \ngot:%v\n", expMethod.Doc, genMethod.Doc)
				case !reflect.DeepEqual(genMethod.Params, expMethod.Params):
					t.Errorf("gen custom method Params unexpectate \nexp:%v \ngot:%v\n", expMethod.Params, genMethod.Params)
				case !reflect.DeepEqual(genMethod.Result, expMethod.Result):
					t.Errorf("gen custom method Result unexpectate \nexp:%v \ngot:%v\n", expMethod.Result, genMethod.Result)
				case !reflect.DeepEqual(genMethod.Body, expMethod.Body):
					t.Errorf("gen custom method Body unexpectate \nexp:%v \ngot:%v\n", expMethod.Body, genMethod.Body)
				default:
					pass = true
				}
				break
			}
		}
		if !pass {
			t.Errorf("gen custom method unexpectate \nexp:%v \n but not found in gen methods", expMethod)
		}

	}
}

func TestBaseStruct_AddMethod(t *testing.T) {
	u := OnlyForTestUser{}
	testcases := []struct {
		MethodOrStruct interface{}
		CustomMethods  []*parser.Method
	}{
		{
			MethodOrStruct: u.IsEmpty, // param is a function
			CustomMethods: []*parser.Method{
				{
					Receiver: parser.Param{
						Name:      "u",
						Type:      "User",
						IsPointer: true,
					},
					MethodName: "IsEmpty",
					Doc:        "IsEmpty is a custom method\n",
					Result: []parser.Param{
						{
							Type: "bool",
						},
					},
					Body: "{\n\tif u == nil {\n\t\treturn true\n\t}\n\n\treturn u.ID == 0\n}",
				},
			},
		},
		{
			MethodOrStruct: u, // param is a struct
			CustomMethods: []*parser.Method{
				{
					Receiver: parser.Param{
						Name:      "u",
						Type:      "User",
						IsPointer: true,
					},
					MethodName: "IsEmpty",
					Doc:        "IsEmpty is a custom method\n",
					Result: []parser.Param{
						{
							Type: "bool",
						},
					},
					Body: "{\n\tif u == nil {\n\t\treturn true\n\t}\n\n\treturn u.ID == 0\n}",
				},
				{
					Receiver: parser.Param{
						Name:      "u",
						Type:      "User",
						IsPointer: true,
					},
					MethodName: "SetName",
					Doc:        "SetName set user name\n",
					Params: []parser.Param{
						{
							Name: "name",
							Type: "string",
						},
					},
					Body: "{\n\tu.Name = name\n}",
				},
				{
					Receiver: parser.Param{
						Name:      "u",
						Type:      "User",
						IsPointer: true,
					},
					MethodName: "GetName",
					Doc:        "GetName get to lower name\n",
					Result: []parser.Param{
						{
							Type: "string",
						},
					},
					Body: "{\n\treturn strings.ToLower(u.Name)\n}",
				},
			},
		},
	}
	for _, testcase := range testcases {
		checkCustomMethod(t, testcase.CustomMethods, testcase.MethodOrStruct)
	}

}
