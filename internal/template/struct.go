package template

const (
	BaseStruct = createMethod + `
	type {{.NewStructName}} struct {
		{{.NewStructName}}Do
		` + members + `
	}
	
	` + getFieldMethod + cloneMethod + relationship + defineMethodStruct

	BaseStructWithContext = createMethod + `
	type {{.NewStructName}} struct {
		{{.NewStructName}}Do {{.NewStructName}}Do
		` + members + `
	}
	
	func ({{.S}} *{{.NewStructName}}) WithContext(ctx context.Context) *{{.NewStructName}}Do { return {{.S}}.{{.NewStructName}}Do.WithContext(ctx)}

	func ({{.S}} {{.NewStructName}}) TableName() string { return {{.S}}.{{.NewStructName}}Do.TableName()} 

	
	
	` + getFieldMethod + cloneMethod + relationship + defineMethodStruct
)

const (
	createMethod = `
	func new{{.StructName}}(db *gorm.DB) {{.NewStructName}} {
		_{{.NewStructName}} := {{.NewStructName}}{}
	
		_{{.NewStructName}}.{{.NewStructName}}Do.UseDB(db)
		_{{.NewStructName}}.{{.NewStructName}}Do.UseModel(&{{.StructInfo.Package}}.{{.StructInfo.Type}}{})
	
		{{if .HasMember}}tableName := _{{.NewStructName}}.{{.NewStructName}}Do.TableName(){{end}}
		_{{$.NewStructName}}.ALL = field.NewField(tableName, "*")
		{{range .Members -}}
		{{if not .IsRelation -}}
			_{{$.NewStructName}}.{{.Name}} = field.New{{.GenType}}(tableName, "{{.ColumnName}}")
		{{- else -}}
			_{{$.NewStructName}}.{{.Relation.Name}} = {{$.NewStructName}}{{.Relation.RelationshipName}}{{.Relation.Name}}{
				db: db.Session(&gorm.Session{}),

				{{.Relation.StructMemberInit}}
			}
		{{end}}
		{{end}}

		{{range .Members -}}
		{{if not .IsRelation -}}
			_{{$.NewStructName}}.fieldMap["{{.ColumnName}}"] = _{{$.NewStructName}}.{{.Name}}
		{{end -}}
		{{end}}
		return _{{.NewStructName}}
	}
	`
	members = `

	ALL field.Field
	{{range .Members -}}
	{{if not .IsRelation -}}
		{{.Name}} field.{{.GenType}}
	{{- else -}}
		{{.Relation.Name}} {{$.NewStructName}}{{.Relation.RelationshipName}}{{.Relation.Name}}
	{{end}}
	{{end}}

	fieldMap  map[string]field.Expr
`
	cloneMethod = `
func ({{.S}} {{.NewStructName}}) clone(db *gorm.DB) {{.NewStructName}} {
	{{.S}}.{{.NewStructName}}Do.ReplaceDB(db)
	return {{.S}}
}
`
	getFieldMethod = `
func ({{.S}} *{{.NewStructName}}) GetFieldByName(fieldName string) (field.Expr, bool) {
	field, ok := {{.S}}.fieldMap[fieldName]
	return field, ok
}
`
	relationship = `{{range .Members}}{{if .IsRelation}}` +
		`{{- $relation := .Relation }}{{- $relationship := $relation.RelationshipName}}` +
		relationStruct + relationTx +
		`{{end}}{{end}}`
	defineMethodStruct = `type {{.NewStructName}}Do struct { gen.DO }`
)

const (
	relationStruct = `
type {{$.NewStructName}}{{$relationship}}{{$relation.Name}} struct{
	db *gorm.DB
	
	field.RelationField
	
	{{$relation.StructMember}}
}

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}) Where(conds ...field.Expr) *{{$.NewStructName}}{{$relationship}}{{$relation.Name}} {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}) WithContext(ctx context.Context) *{{$.NewStructName}}{{$relationship}}{{$relation.Name}} {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}) Model(m *{{$.StructInfo.Package}}.{{$.StructInfo.Type}}) *{{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx {
	return &{{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx{a.db.Model(m).Association(a.Name())}
}

`
	relationTx = `
type {{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx struct{ tx *gorm.Association }

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx) Find() (result {{if eq $relationship "HasMany" "Many2Many"}}[]{{end}}*{{$relation.Type}}, err error) {
	return result, a.tx.Find(&result)
}

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx) Append(values ...*{{$relation.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx) Replace(values ...*{{$relation.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx) Delete(values ...*{{$relation.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx) Clear() error {
	return a.tx.Clear()
}

func (a {{$.NewStructName}}{{$relationship}}{{$relation.Name}}Tx) Count() int64 {
	return a.tx.Count()
}
`
)
