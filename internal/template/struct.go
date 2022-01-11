package template

const (
	BaseStruct = createMethod + `
	type {{.NewStructName}} struct {
		{{.NewStructName}}Do
		` + fields + `
	}
	` + tableMethod + asMethond + updateFieldMethod + getFieldMethod + fillFieldMapMethod + cloneMethod + relationship + defineMethodStruct

	BaseStructWithContext = createMethod + `
	type {{.NewStructName}} struct {
		{{.NewStructName}}Do {{.NewStructName}}Do
		` + fields + `
	}
	` + tableMethod + asMethond + updateFieldMethod + `
	
	func ({{.S}} *{{.NewStructName}}) WithContext(ctx context.Context) I{{.StructName}}Do { return {{.S}}.{{.NewStructName}}Do.WithContext(ctx)}

	func ({{.S}} {{.NewStructName}}) TableName() string { return {{.S}}.{{.NewStructName}}Do.TableName()} 

	` + getFieldMethod + fillFieldMapMethod + cloneMethod + relationship + defineMethodStruct + defineDoInterface
)

const (
	createMethod = `
	func new{{.StructName}}(db *gorm.DB) {{.NewStructName}} {
		_{{.NewStructName}} := {{.NewStructName}}{}
	
		_{{.NewStructName}}.{{.NewStructName}}Do.UseDB(db)
		_{{.NewStructName}}.{{.NewStructName}}Do.UseModel(&{{.StructInfo.Package}}.{{.StructInfo.Type}}{})
	
		tableName := _{{.NewStructName}}.{{.NewStructName}}Do.TableName()
		_{{$.NewStructName}}.ALL = field.NewField(tableName, "*")
		{{range .Fields -}}
		{{if not .IsRelation -}}
			{{- if .ColumnName -}}_{{$.NewStructName}}.{{.Name}} = field.New{{.GenType}}(tableName, "{{.ColumnName}}"){{- end -}}
		{{- else -}}
			_{{$.NewStructName}}.{{.Relation.Name}} = {{$.NewStructName}}{{.Relation.RelationshipName}}{{.Relation.Name}}{
				db: db.Session(&gorm.Session{}),

				{{.Relation.StructFieldInit}}
			}
		{{end}}
		{{end}}

		_{{$.NewStructName}}.fillFieldMap()
		
		return _{{.NewStructName}}
	}
	`
	fields = `
	ALL field.Field
	{{range .Fields -}}
	{{if not .IsRelation -}}
		{{- if .ColumnName -}}{{.Name}} field.{{.GenType}}{{- end -}}
	{{- else -}}
		{{.Relation.Name}} {{$.NewStructName}}{{.Relation.RelationshipName}}{{.Relation.Name}}
	{{end}}
	{{end}}

	fieldMap  map[string]field.Expr
`
	tableMethod = `
func ({{.S}} {{.NewStructName}}) Table(newTableName string) *{{.NewStructName}} { 
	{{.S}}.{{.NewStructName}}Do.UseTable(newTableName)
	return {{.S}}.updateTableName(newTableName)
}
`

	asMethond = `	
func ({{.S}} {{.NewStructName}}) As(alias string) *{{.NewStructName}} { 
	{{.S}}.{{.NewStructName}}Do.DO = *({{.S}}.{{.NewStructName}}Do.As(alias).(*gen.DO))
	return {{.S}}.updateTableName(alias)
}
`
	updateFieldMethod = `
func ({{.S}} *{{.NewStructName}}) updateTableName(table string) *{{.NewStructName}} { 
	{{.S}}.ALL = field.NewField(table, "*")
	{{range .Fields -}}
	{{if not .IsRelation -}}
		{{- if .ColumnName -}}{{$.S}}.{{.Name}} = field.New{{.GenType}}(table, "{{.ColumnName}}"){{- end -}}
	{{end}}
	{{end}}
	
	{{.S}}.fillFieldMap()

	return {{.S}}
}
`

	cloneMethod = `
func ({{.S}} {{.NewStructName}}) clone(db *gorm.DB) {{.NewStructName}} {
	{{.S}}.{{.NewStructName}}Do.ReplaceDB(db)
	return {{.S}}
}
`
	getFieldMethod = `
func ({{.S}} *{{.NewStructName}}) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := {{.S}}.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe,ok := _f.(field.OrderExpr)
	return _oe,ok
}
`
	relationship = `{{range .Fields}}{{if .IsRelation}}` +
		`{{- $relation := .Relation }}{{- $relationship := $relation.RelationshipName}}` +
		relationStruct + relationTx +
		`{{end}}{{end}}`
	defineMethodStruct = `type {{.NewStructName}}Do struct { gen.DO }`

	fillFieldMapMethod = `
func ({{.S}} *{{.NewStructName}}) fillFieldMap() {
	{{.S}}.fieldMap =  make(map[string]field.Expr, {{len .Fields}})
	{{range .Fields -}}
	{{if not .IsRelation -}}
		{{- if .ColumnName -}}{{$.S}}.fieldMap["{{.ColumnName}}"] = {{$.S}}.{{.Name}}{{- end -}}
	{{end}}
	{{end -}}
}
`

	defineDoInterface = `

type I{{.StructName}}Do interface {
	Debug() I{{.StructName}}Do
	WithContext(ctx context.Context) I{{.StructName}}Do
	Clauses(conds ...clause.Expression) I{{.StructName}}Do
	Not(conds ...gen.Condition) I{{.StructName}}Do
	Or(conds ...gen.Condition) I{{.StructName}}Do
	Select(conds ...field.Expr) I{{.StructName}}Do
	Where(conds ...gen.Condition) I{{.StructName}}Do
	Order(conds ...field.Expr) I{{.StructName}}Do
	Distinct(cols ...field.Expr) I{{.StructName}}Do
	Omit(cols ...field.Expr) I{{.StructName}}Do
	Join(table schema.Tabler, on ...field.Expr) I{{.StructName}}Do
	LeftJoin(table schema.Tabler, on ...field.Expr) I{{.StructName}}Do
	RightJoin(table schema.Tabler, on ...field.Expr) I{{.StructName}}Do
	Group(cols ...field.Expr) I{{.StructName}}Do
	Having(conds ...gen.Condition) I{{.StructName}}Do
	Limit(limit int) I{{.StructName}}Do
	Offset(offset int) I{{.StructName}}Do
	Scopes(funcs ...func(gen.Dao) gen.Dao) I{{.StructName}}Do
	Unscoped() I{{.StructName}}Do
	Create(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) error
	CreateInBatches(values []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int) error
	Save(values ...*{{.StructInfo.Package}}.{{.StructInfo.Type}}) error
	First() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	Take() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	Last() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	Find() ([]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, err error)
	FindInBatches(result *[]*{{.StructInfo.Package}}.{{.StructInfo.Type}}, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Attrs(attrs ...field.AssignExpr) I{{.StructName}}Do
	Assign(attrs ...field.AssignExpr) I{{.StructName}}Do
	Joins(field field.RelationField) I{{.StructName}}Do
	Preload(field field.RelationField) I{{.StructName}}Do
	FirstOrInit() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	FirstOrCreate() (*{{.StructInfo.Package}}.{{.StructInfo.Type}}, error)
	FindByPage(offset int, limit int) (result []*{{.StructInfo.Package}}.{{.StructInfo.Type}}, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
}
`
)

const (
	relationStruct = `
type {{$.NewStructName}}{{$relationship}}{{$relation.Name}} struct{
	db *gorm.DB
	
	field.RelationField
	
	{{$relation.StructField}}
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
