package template

// DefaultQuery default query
const DefaultQuery = `
var (
	Q =new(Query)
	{{range $name,$d :=.Data -}}
	{{$d.ModelStructName}} *{{$d.QueryStructName}}
	{{end -}}
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db,opts...)
	{{range $name,$d :=.Data -}}
	{{$d.ModelStructName}} = &Q.{{$d.ModelStructName}}
	{{end -}}
}

`

// QueryMethod query method template
const QueryMethod = `
func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db: db,
		{{range $name,$d :=.Data -}}
		{{$d.ModelStructName}}: new{{$d.ModelStructName}}(db,opts...),
		{{end -}}
	}
}

type Query struct{
	db *gorm.DB

	{{range $name,$d :=.Data -}}
	{{$d.ModelStructName}} {{$d.QueryStructName}}
	{{end}}
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db: db,
		{{range $name,$d :=.Data -}}
		{{$d.ModelStructName}}: q.{{$d.ModelStructName}}.clone(db),
		{{end}}
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db: db,
		{{range $name,$d :=.Data -}}
		{{$d.ModelStructName}}: q.{{$d.ModelStructName}}.replaceDB(db),
		{{end}}
	}
}

type queryCtx struct{ 
	{{range $name,$d :=.Data -}}
	{{$d.ModelStructName}} {{$d.ReturnObject}}
	{{end}}
}

func (q *Query) WithContext(ctx context.Context) *queryCtx  {
	return &queryCtx{
		{{range $name,$d :=.Data -}}
		{{$d.ModelStructName}}: q.{{$d.ModelStructName}}.WithContext(ctx),
		{{end}}
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}

`

// QueryMethodTest query method test template
const QueryMethodTest = `

const _gen_test_db_name = "gen_test.db"

var _gen_test_db *gorm.DB
var _gen_test_once sync.Once

func init() {
	InitializeDB()
	_gen_test_db.AutoMigrate(&_another{})
}

func InitializeDB() {
	_gen_test_once.Do(func() {
		var err error
		_gen_test_db, err = gorm.Open(sqlite.Open(_gen_test_db_name), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("open sqlite %q fail: %w", _gen_test_db_name, err))
		}
	})
}

func assert(t *testing.T, methodName string, res, exp interface{}) {
	if !reflect.DeepEqual(res, exp) {
		t.Errorf("%v() gotResult = %v, want %v", methodName, res, exp)
	}
}

type _another struct {
	ID uint64 ` + "`" + `gorm:"primaryKey"` + "`" + `
}

func (*_another) TableName() string { return "another_for_unit_test" }

func Test_Available(t *testing.T) {
	if !Use(_gen_test_db).Available() {
		t.Errorf("query.Available() == false")
	}
}

func Test_WithContext(t *testing.T) {
	query := Use(_gen_test_db)
	if !query.Available() {
		t.Errorf("query Use(_gen_test_db) fail: query.Available() == false")
	}

	type Content string
	var key, value Content = "gen_tag", "unit_test"
	qCtx := query.WithContext(context.WithValue(context.Background(), key, value))

	for _, ctx := range []context.Context{
		{{range $name,$d :=.Data -}}
		qCtx.{{$d.ModelStructName}}.UnderlyingDB().Statement.Context,
		{{end}}
	} {
		if v := ctx.Value(key); v != value {
			t.Errorf("get value from context fail, expect %q, got %q", value, v)
		}
	}
}

func Test_Transaction(t *testing.T) {
	query := Use(_gen_test_db)
	if !query.Available() {
		t.Errorf("query Use(_gen_test_db) fail: query.Available() == false")
	}

	err := query.Transaction(func(tx *Query) error { return nil })
	if err != nil {
		t.Errorf("query.Transaction execute fail: %s", err)
	}

	tx := query.Begin()

	err = tx.SavePoint("point")
	if err != nil {
		t.Errorf("query tx SavePoint fail: %s", err)
	}
	err = tx.RollbackTo("point")
	if err != nil {
		t.Errorf("query tx RollbackTo fail: %s", err)
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("query tx Commit fail: %s", err)
	}

	err = query.Begin().Rollback()
	if err != nil {
		t.Errorf("query tx Rollback fail: %s", err)
	}
}
`
