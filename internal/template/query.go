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
		{{$d.StructName}}: New{{$d.StructName}}(db),
		{{end -}}
	}
}

type Query struct{
	db *gorm.DB

	{{range $name,$d :=.Data -}}
	{{$d.StructName}} *{{$d.NewStructName}}
	{{end}}
}

func (q *Query) Transaction(fc func(db *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.withTx(tx)) }, opts...)
}

func (q Query) Begin(opts ...*sql.TxOptions) *Query {
	q.db = q.db.Begin(opts...)
	return &q
}

func (q Query) Commit() *Query {
	q.db = q.db.Commit()
	return &q
}

func (q Query) Rollback() *Query {
	q.db = q.db.Rollback()
	return &q
}

func (q Query) SavePoint(name string) *Query {
	q.db = q.db.SavePoint(name)
	return &q
}

func (q Query) RollbackTo(name string) *Query {
	q.db = q.db.RollbackTo(name)
	return &q
}

func (q Query) withTx(tx *gorm.DB) *Query {
	q.db = tx
	return &q
}

`
