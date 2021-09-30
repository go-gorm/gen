package template

const (
	BaseStruct = createMethod + `
	type {{.NewStructName}} struct {
		{{.NewStructName}}Do
		` + members + `
	}
	
	` + cloneMethod + relationship + defineMethodStruct

	BaseStructWithContext = createMethod + `
	type {{.NewStructName}} struct {
		{{.NewStructName}}Do {{.NewStructName}}Do
		` + members + `
	}
	
	func ({{.S}} *{{.NewStructName}}) WithContext(ctx context.Context) *{{.NewStructName}}Do { return {{.S}}.{{.NewStructName}}Do.WithContext(ctx)}
	
	` + cloneMethod + relationship + defineMethodStruct
)

const (
	createMethod = `
	func new{{.StructName}}(db *gorm.DB) {{.NewStructName}} {
		_{{.NewStructName}} := {{.NewStructName}}{}
	
		_{{.NewStructName}}.{{.NewStructName}}Do.UseDB(db)
		_{{.NewStructName}}.{{.NewStructName}}Do.UseModel({{.StructInfo.Package}}.{{.StructInfo.Type}}{})
	
		{{if .HasMember}}tableName := _{{.NewStructName}}.{{.NewStructName}}Do.TableName(){{end}}
		{{range .Members}} _{{$.NewStructName}}.{{.Name}} = field.New{{.GenType}}(tableName, "{{.ColumnName}}")
		{{end}}
		{{range .Relations.HasOne}}
			_{{$.NewStructName}}.{{.Name}} = {{$.NewStructName}}HasOne{{.Name}}{
			db: db.Session(&gorm.Session{}),

			{{.StructMemberInit}}
		}
		{{end}}
		{{- range .Relations.HasMany}}
			_{{$.NewStructName}}.{{.Name}} = {{$.NewStructName}}HasMany{{.Name}}{
			db: db.Session(&gorm.Session{}),

			{{.StructMemberInit}}
		}
		{{end}}
		{{- range .Relations.BelongsTo}}
			_{{$.NewStructName}}.{{.Name}} = {{$.NewStructName}}BelongsTo{{.Name}}{
			db: db.Session(&gorm.Session{}),
			
			{{.StructMemberInit}}
		}
		{{end}}
		{{- range .Relations.Many2Many}}
			_{{$.NewStructName}}.{{.Name}} = {{$.NewStructName}}Many2Many{{.Name}}{
			db: db.Session(&gorm.Session{}),
			
			{{.StructMemberInit}}
		}
		{{end}}
	
		return _{{.NewStructName}}
	}
	`
	members = `
	{{range .Members}}{{.Name}} field.{{.GenType}}
	{{end}}
	{{range .Relations.HasOne}}{{.Name}} {{$.NewStructName}}HasOne{{.Name}}
	{{end}}
	{{- range .Relations.HasMany}}{{.Name}} {{$.NewStructName}}HasMany{{.Name}}
	{{end}}
	{{- range .Relations.BelongsTo}}{{.Name}} {{$.NewStructName}}BelongsTo{{.Name}}
	{{end}}
	{{- range .Relations.Many2Many}}{{.Name}} {{$.NewStructName}}Many2Many{{.Name}}
	{{end}}
`
	cloneMethod = `
func ({{.S}} {{.NewStructName}}) clone(db *gorm.DB) {{.NewStructName}} {
	{{.S}}.{{.NewStructName}}Do.ReplaceDB(db)
	return {{.S}}
}
`
	relationship       = hasOneRelationship + hasManyRelationship + belongsToRelationship + many2ManyRelationship
	defineMethodStruct = `type {{.NewStructName}}Do struct { gen.DO }`
)

const (
	hasOneRelationship = `{{- $relationship := "HasOne"}}	
{{range .Relations.HasOne}}` + relationStruct + `{{end}}`
	hasManyRelationship = `{{- $relationship := "HasMany"}}	
{{range .Relations.HasMany}}` + relationStruct + `{{end}}`
	belongsToRelationship = `{{- $relationship := "BelongsTo"}}	
{{range .Relations.BelongsTo}}` + relationStruct + `{{end}}`
	many2ManyRelationship = `{{- $relationship := "Many2Many"}}	
{{range .Relations.Many2Many}}` + relationStruct + `{{end}}`
	relationStruct = `
type {{$.NewStructName}}{{$relationship}}{{.Name}} struct{
	db *gorm.DB
	
	field.RelationField
	
	{{.StructMember}}
}

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}) Where(conds ...field.Expr) *{{$.NewStructName}}{{$relationship}}{{.Name}} {
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

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}) WithContext(ctx context.Context) *{{$.NewStructName}}{{$relationship}}{{.Name}} {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}) Model(m *{{$.StructInfo.Package}}.{{$.StructInfo.Type}}) *{{$.NewStructName}}{{$relationship}}{{.Name}}Tx {
	return &{{$.NewStructName}}{{$relationship}}{{.Name}}Tx{a.db.Model(m).Association(a.Name())}
}

` + relationTx
	relationTx = `
type {{$.NewStructName}}{{$relationship}}{{.Name}}Tx struct{ tx *gorm.Association }

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}Tx) Find() (result {{if eq $relationship "HasMany" "Many2Many"}}[]{{end}}*{{.Type}}, err error) {
	return result, a.tx.Find(&result)
}

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}Tx) Append(values ...*{{.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}Tx) Replace(values ...*{{.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}Tx) Delete(values ...*{{.Type}}) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}Tx) Clear() error {
	return a.tx.Clear()
}

func (a {{$.NewStructName}}{{$relationship}}{{.Name}}Tx) Count() int64 {
	return a.tx.Count()
}
`
)
