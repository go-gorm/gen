package main

import (
	"gorm.io/gen"
	"gorm.io/gen/helper"
)

func init() {
	// prepare() // prepare table for generate
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "/tmp/gentest/query",
		ModelPkgPath: "/tmp/gentest/demo",
	})

	base := &Demo{
		pkgName:    "demo",
		structName: "DemoStruct1",
		tableName:  "demoTable1",
		fields: []helper.Field{
			&DemoField{
				name:    "ID",
				typ:     "uint",
				gormTag: "column:id;type:bigint unsigned;primaryKey;autoIncrement:true",
				jsonTag: "id",
				tag:     ` kms:"enc:aes"`,
				comment: "主键",
			},
		},
	}

	demo := &Demo{
		pkgName:    "demo",
		structName: "DemoStruct",
		tableName:  "demoTable",
		fields: []helper.Field{
			&DemoField{
				name:    "ID",
				typ:     "uint",
				gormTag: "column:id;type:bigint unsigned;primaryKey;autoIncrement:true",
				jsonTag: "id",
				tag:     ` kms:"enc:aes"`,
				comment: "主键",
			},
			&DemoField{
				name:    "Username",
				typ:     "[]DemoStruct1",
				gormTag: "column:username",
				jsonTag: "username",
				comment: "用户名",
			},
			&DemoField{
				name:    "Detail",
				typ:     "json.RawMessage",
				gormTag: "column:detail",
				jsonTag: "detail",
				tag:     ` kms:"enc:aes"`,
				comment: "用户敏感数据\n需加密",
			},
		},
	}

	g.GenerateModelFrom(base)

	g.ApplyBasic(g.GenerateModelFrom(demo))

	g.Execute()
}
