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

// TODO remove pointer && query clone
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
	{{$d.StructName}} *{{$d.NewStructName}}
	{{end}}
}

func (q *Query) Transaction(fc func(db *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(Use(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *Query {
	return Use(q.db.Begin(opts...))
}

func (q *Query) Commit() error {
	return q.db.Commit().Error
}

func (q *Query) Rollback() error {
	return q.db.Rollback().Error
}

func (q *Query) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *Query) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}

`
