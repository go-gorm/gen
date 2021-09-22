package template

const createMethod = `
func new{{.StructName}}(db *gorm.DB) {{.NewStructName}} {
	_{{.NewStructName}} := {{.NewStructName}}{}

	_{{.NewStructName}}.{{.NewStructName}}Do.UseDB(db)
	_{{.NewStructName}}.{{.NewStructName}}Do.UseModel({{.StructInfo.Package}}.{{.StructInfo.Type}}{})

	{{if .HasMember}}tableName := _{{.NewStructName}}.{{.NewStructName}}Do.TableName(){{end}}
	{{range $p :=.Members}} _{{$.NewStructName}}.{{$p.Name}} = field.New{{$p.NewType}}(tableName, "{{$p.ColumnName}}")
	{{end}}

	{{if .HasMember}}tableName := _{{.NewStructName}}.{{.NewStructName}}Do.TableName(){{end}}
	{{range $p :=.Members}} _{{$.NewStructName}}.{{$p.Name}} = field.New{{$p.NewType}}(tableName, "{{$p.ColumnName}}")
	{{end}}
	{{if .HasRelation}}
	{{range $p :=.Relations}}_{{$.NewStructName}}.{{$p.Name}} = {{$p.RelationShip}}{{$p.Name}}{
		Relation: field.NewRelation("{{$p.Name}}"),
		db:       db.Session(&gorm.Session{}),
	}
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
	{{.NewStructName}}Do` +
	Members + `
}

` + cloneMethod + defineMethodStruct

const BaseStructWithContext = createMethod + `
type {{.NewStructName}} struct {
	{{.NewStructName}}Do {{.NewStructName}}Do` +
	Members + `
}

func ({{.S}} *{{.NewStructName}}) WithContext(ctx context.Context) *{{.NewStructName}}Do { return {{.S}}.{{.NewStructName}}Do.WithContext(ctx)}

` + cloneMethod + defineMethodStruct

const Members = `
	{{range $p :=.Members}}{{$p.Name}}  field.{{$p.NewType}}
	{{end}}
	{{range $p :=.Relations}}{{$p.Name}} {{$p.RelationShip}}{{$p.Name}}
	{{end}}
`

const SingleRelation = Relation + SingleRelateFind
const ManyRelation = Relation + ManyRelateFind

const RelateStruct = `
type {{$p.RelationShip}}{{$p.Name}} struct {
	field.Relation

	db *gorm.DB
}

func (a {{$p.RelationShip}}{{$p.Name}}) Model(m *{{.StructInfo.Package}}.{{.StructInfo.Type}}) *{{$p.RelationShip}}{{$p.Name}}Tx {
	return &{{$p.RelationShip}}{{$p.Name}}Tx{a.db.Model(m).Association(string(a.Path()))}
}
`

const Relation = RelateStruct + `
type {{$p.RelationShip}}{{$p.Name}}Tx struct{ tx *gorm.Association }

func (a {{$p.RelationShip}}{{$p.Name}}Tx) Append(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}})) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a {{$p.RelationShip}}{{$p.Name}}Tx) Replace(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}})) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}


func (a {{$p.RelationShip}}{{$p.Name}}Tx) Delete(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}})) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a {{$p.RelationShip}}{{$p.Name}}Tx) Clear() error {
	return a.tx.Clear()
}

func (a {{$p.RelationShip}}{{$p.Name}}Tx) Count() int64 {
	return a.tx.Count()
}
`
const SingleRelateFind = `
func (a {{$p.RelationShip}}{{$p.Name}}Tx) Find() (result *{{.StructInfo.Package}}.{{.StructInfo.Type}}), err error) {
	return result, a.tx.Find(&result)
}
`

const ManyRelateFind = `
func (a {{$p.RelationShip}}{{$p.Name}}Tx) Find() (result []*{{.StructInfo.Package}}.{{.StructInfo.Type}}), err error) {
	return result, a.tx.Find(&result)
}
`
