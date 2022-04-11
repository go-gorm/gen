package main

import (
	"gorm.io/gen"
	"gorm.io/gen/helper"
)

var detail, data helper.Object

func init() {
	detail = &Demo{
		structName: "Detail",
		fileName:   "diy_data_detail",
		fields: []helper.Field{
			&DemoField{
				name:    "Username",
				typ:     "string",
				jsonTag: "username",
				comment: "用户名",
			},
			&DemoField{
				name:    "Age",
				typ:     "uint",
				jsonTag: "age",
				comment: "用户年龄",
			},
			&DemoField{
				name:    "Phone",
				typ:     "string",
				jsonTag: "phone",
				comment: "手机号",
			},
		},
	}

	data = &Demo{
		structName: "Data",
		tableName:  "data",
		fileName:   "diy_data",
		fields: []helper.Field{
			&DemoField{
				name:    "ID",
				typ:     "uint",
				gormTag: "column:id;type:bigint unsigned;primaryKey;autoIncrement:true",
				jsonTag: "id",
				tag:     `kms:"enc:aes"`,
				comment: "主键",
			},
			&DemoField{
				name:    "UserInfo",
				typ:     "[]Detail",
				jsonTag: "user_info",
				comment: "用户信息",
			},
			&DemoField{
				name:    "Remark",
				typ:     "json.RawMessage",
				gormTag: "column:detail",
				jsonTag: "remark",
				tag:     `kms:"enc:aes"`,
				comment: "备注\n详细信息",
			},
		},
	}
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "/tmp/gentest/query",
		ModelPkgPath: "/tmp/gentest/demo",
	})

	g.GenerateModelFrom(detail)

	g.ApplyBasic(g.GenerateModelFrom(data))

	g.Execute()
}
