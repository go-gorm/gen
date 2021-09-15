package template

const DefaultQueryTmpl = `
var (
	Q *Query
	{{range $name,$d :=.Data -}}
	{{$d.StructName}} *{{$d.NewStructName}}
	{{end -}}
)

func SetDefault(db *gorm.DB) {
	Q = Use(db)
	{{range $name,$d :=.Data -}}
	{{$d.StructName}} = Q.{{$d.StructName}}
	{{end -}}
}

`

const QueryTmpl = `
func Use(db *gorm.DB) *Query {
	return &Query{
		db: db,
		{{range $name,$d :=.Data -}}
		{{$d.StructName}}: new{{$d.StructName}}(db),
		{{end -}}
	}
}

type Query struct{
	db *gorm.DB

	{{range $name,$d :=.Data -}}
	{{$d.StructName}} {{$d.NewStructName}}
	{{end}}
}

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		{{range $name,$d :=.Data -}}
		{{$d.StructName}}: q.{{$d.StructName}}.clone(db),
		{{end}}
	}
}

type queryCtx struct{ 
	{{range $name,$d :=.Data -}}
	{{$d.StructName}} {{$d.NewStructName}}Do
	{{end}}
}

func (q *Query) WithContext(ctx context.Context) *queryCtx  {
	return &queryCtx{
		{{range $name,$d :=.Data -}}
		{{$d.StructName}}: q.{{$d.StructName}}.{{$d.NewStructName}}Do,
		{{end}}
	}
}

func (q *Query) Transaction(fc func(db *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *queryTx {
	return &queryTx{q.clone(q.db.Begin(opts...))}
}

type queryTx struct{ *Query }

func (q *queryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *queryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *queryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *queryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}

`
