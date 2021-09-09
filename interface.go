package gen

import (
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"gorm.io/gen/field"
)

type (
	// Condition query condition
	// field.Expr and subquery are expect value
	Condition interface{ ConditionMark() }
)

type subQuery interface {
	underlyingDB() *gorm.DB
	underlyingDO() *DO

	ConditionMark()
}

// Dao CRUD methods
type Dao interface {
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
	First() (result interface{}, err error)
	Take() (result interface{}, err error)
	Last() (result interface{}, err error)
	Find() (results interface{}, err error)
	FindInBatches(dest interface{}, batchSize int, fc func(tx Dao, batch int) error) error
	// FirstOrInit(dest interface{}) error
	// FirstOrCreate(dest interface{}) error
	Model(model interface{}) Dao
	Update(column field.Expr, value interface{}) error
	UpdateSimple(column field.Expr) error
	Updates(values interface{}) error
	UpdateColumn(column field.Expr, value interface{}) error
	UpdateColumns(values interface{}) error
	Delete() error
	Count() (int64, error)
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
