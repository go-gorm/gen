package gen

import (
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen/field"
)

var _ Condition = Dao(&DO{})
var _ Condition = field.Expr(field.NewField("", "field"))

type (
	// Condition query condition
	// field.Expr and subquery(Dao) are expect value
	Condition clause.Expression

	// Hint hints.Hints and hints.IndexHint are expect value
	Hint clause.Expression
)

type subQuery interface {
	UnderlyingDB() *gorm.DB
	underlyingDO() *DO
	build(opts ...stmtOpt) *gorm.Statement
}

// Dao CRUD methods
type Dao interface {
	clause.Expression
	subQuery

	As(alias string) Dao

	Not(conds ...Condition) Dao
	Or(conds ...Condition) Dao

	Select(columns ...field.Expr) Dao
	Where(conds ...Condition) Dao
	Order(columns ...field.Expr) Dao
	Distinct(columns ...field.Expr) Dao
	Omit(columns ...field.Expr) Dao
	Join(table schema.Tabler, conds ...Condition) Dao
	LeftJoin(table schema.Tabler, conds ...Condition) Dao
	RightJoin(table schema.Tabler, conds ...Condition) Dao
	Group(columns field.Expr) Dao
	Having(conds ...Condition) Dao
	Limit(limit int) Dao
	Offset(offset int) Dao
	Scopes(funcs ...func(Dao) Dao) Dao
	// Preload(conds ...field.Expr) Dao
	// Attrs(attrs ...interface{}) Dao
	// Assign(attrs ...interface{}) Dao
	Unscoped() Dao

	Create(value interface{}) error
	CreateInBatches(value interface{}, batchSize int) error
	Save(value interface{}) error
	First(dest interface{}, conds ...field.Expr) error
	Last(dest interface{}, conds ...field.Expr) error
	Take(dest interface{}, conds ...field.Expr) error
	Find(dest interface{}, conds ...field.Expr) error
	FindInBatches(dest interface{}, batchSize int, fc func(tx Dao, batch int) error) error
	FirstOrInit(dest interface{}, conds ...field.Expr) error
	FirstOrCreate(dest interface{}, conds ...field.Expr) error
	Update(column field.Expr, value interface{}) error
	Updates(values interface{}) error
	UpdateColumn(column field.Expr, value interface{}) error
	UpdateColumns(values interface{}) error
	Delete(value interface{}, conds ...field.Expr) error
	Count(count *int64) error
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) error
	Pluck(column field.Expr, dest interface{}) error
	ScanRows(rows *sql.Rows, dest interface{}) error

	Transaction(fc func(tx Dao) error, opts ...*sql.TxOptions) error
	Begin(opts ...*sql.TxOptions) Dao
	Commit() Dao
	RollBack() Dao
	SavePoint(name string) Dao
	RollBackTo(name string) Dao
}
