package template

const createMethod = `
func new{{.StructName}}(db *gorm.DB) {{.NewStructName}} {
	_{{.NewStructName}} := {{.NewStructName}}{}

	_{{.NewStructName}}.{{.NewStructName}}Do.UseDB(db)
	_{{.NewStructName}}.{{.NewStructName}}Do.UseModel({{.StructInfo.Package}}.{{.StructInfo.Type}}{})

	{{if .HasMember}}tableName := _{{.NewStructName}}.{{.NewStructName}}Do.TableName(){{end}}
	{{range $p :=.Members}} _{{$.NewStructName}}.{{$p.Name}} = field.New{{$p.NewType}}(tableName, "{{$p.ColumnName}}")
	{{end}}
	
	return _{{.NewStructName}}
}

`

const defineMethodStruct = `type {{.NewStructName}}Do struct { gen.DO }`

const cloneMethod = `
func ({{.S}} {{.NewStructName}}) clone(db *gorm.DB) {{.NewStructName}} {
	{{.S}}.{{.NewStructName}}Do.ReplaceDB(db)
	return {{.S}}
}
`

const BaseStruct = createMethod + `
type {{.NewStructName}} struct {
	{{.NewStructName}}Do

	{{range $p :=.Members}}{{$p.Name}}  field.{{$p.NewType}}
	{{end}}
}

` + cloneMethod + defineMethodStruct

const BaseStructWithContext = createMethod + `
type {{.NewStructName}} struct {
	{{.NewStructName}}Do {{.NewStructName}}Do

	{{range $p :=.Members}}{{$p.Name}}  field.{{$p.NewType}}
	{{end}}
}

func ({{.S}} *{{.NewStructName}}) WithContext(ctx context.Context) *{{.NewStructName}}Do { return {{.S}}.{{.NewStructName}}Do.WithContext(ctx)}

` + cloneMethod + defineMethodStruct
