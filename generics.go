package gen

import (
	"context"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type IGenericsDo[T any, E any] interface {
	SubQuery
	Debug() T
	WithContext(ctx context.Context) T
	WithResult(fc func(tx Dao)) ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() T
	WriteDB() T
	As(alias string) T
	Session(config *gorm.Session) T
	Columns(cols ...field.Expr) Columns
	Clauses(conds ...clause.Expression) T
	Not(conds ...Condition) T
	Or(conds ...Condition) T
	Select(conds ...field.Expr) T
	Where(conds ...Condition) T
	Order(conds ...field.Expr) T
	Distinct(cols ...field.Expr) T
	Omit(cols ...field.Expr) T
	Join(table schema.Tabler, on ...field.Expr) T
	LeftJoin(table schema.Tabler, on ...field.Expr) T
	RightJoin(table schema.Tabler, on ...field.Expr) T
	Group(cols ...field.Expr) T
	Having(conds ...Condition) T
	Limit(limit int) T
	Offset(offset int) T
	Count() (count int64, err error)
	Scopes(funcs ...func(Dao) Dao) T
	Unscoped() T
	Create(values ...E) error
	CreateInBatches(values []E, batchSize int) error
	Save(values ...E) error
	First() (E, error)
	Take() (E, error)
	Last() (E, error)
	Find() ([]E, error)
	FindInBatch(batchSize int, fc func(tx Dao, batch int) error) (results []E, err error)
	FindInBatches(result *[]E, batchSize int, fc func(tx Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...E) (info ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info ResultInfo, err error)
	Updates(value interface{}) (info ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info ResultInfo, err error)
	UpdateColumns(value interface{}) (info ResultInfo, err error)
	UpdateFrom(q SubQuery) Dao
	Attrs(attrs ...field.AssignExpr) T
	Assign(attrs ...field.AssignExpr) T
	Joins(fields ...field.RelationField) T
	Preload(fields ...field.RelationField) T
	FirstOrInit() (E, error)
	FirstOrCreate() (E, error)
	FindByPage(offset int, limit int) (result []E, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) T
	UnderlyingDB() *gorm.DB
	schema.Tabler
	UseDB(db *gorm.DB, opts ...DOOption)
	UseModel(model interface{})
	UseTable(tableName string)
	Alias() string
	ReplaceConnPool(pool gorm.ConnPool)
	SetDo(T IGenericsDo[T, E]) T
}

func NewGenericsDo[T IGenericsDo[T, E], E any](realDo T) IGenericsDo[T, E] {
	return &genericsDo[T, E]{realDo: realDo}
}

type genericsDo[T IGenericsDo[T, E], E any] struct {
	DO
	realDo T
}

func (b genericsDo[T, E]) Debug() T {
	return b.withDO(b.DO.Debug())
}
func (b genericsDo[T, E]) As(alias string) T {
	return b.withDO(b.DO.As(alias))
}

func (b genericsDo[T, E]) WithContext(ctx context.Context) T {
	return b.withDO(b.DO.WithContext(ctx))
}

func (b genericsDo[T, E]) ReadDB() T {
	return b.Clauses(dbresolver.Read)
}

func (b genericsDo[T, E]) WriteDB() T {
	return b.Clauses(dbresolver.Write)
}

func (b genericsDo[T, E]) Session(config *gorm.Session) T {
	return b.withDO(b.DO.Session(config))
}

func (b genericsDo[T, E]) Clauses(conds ...clause.Expression) T {
	return b.withDO(b.DO.Clauses(conds...))
}

func (b genericsDo[T, E]) Returning(value interface{}, columns ...string) T {
	return b.withDO(b.DO.Returning(value, columns...))
}

func (b genericsDo[T, E]) Not(conds ...Condition) T {
	return b.withDO(b.DO.Not(conds...))
}

func (b genericsDo[T, E]) Or(conds ...Condition) T {
	return b.withDO(b.DO.Or(conds...))
}

func (b genericsDo[T, E]) Select(conds ...field.Expr) T {
	return b.withDO(b.DO.Select(conds...))
}

func (b genericsDo[T, E]) Where(conds ...Condition) T {
	return b.withDO(b.DO.Where(conds...))
}

func (b genericsDo[T, E]) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) T {
	return b.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (b genericsDo[T, E]) Order(conds ...field.Expr) T {
	return b.withDO(b.DO.Order(conds...))
}

func (b genericsDo[T, E]) Distinct(cols ...field.Expr) T {
	return b.withDO(b.DO.Distinct(cols...))
}

func (b genericsDo[T, E]) Omit(cols ...field.Expr) T {
	return b.withDO(b.DO.Omit(cols...))
}

func (b genericsDo[T, E]) Join(table schema.Tabler, on ...field.Expr) T {
	return b.withDO(b.DO.Join(table, on...))
}

func (b genericsDo[T, E]) LeftJoin(table schema.Tabler, on ...field.Expr) T {
	return b.withDO(b.DO.LeftJoin(table, on...))
}

func (b genericsDo[T, E]) RightJoin(table schema.Tabler, on ...field.Expr) T {
	return b.withDO(b.DO.RightJoin(table, on...))
}

func (b genericsDo[T, E]) Group(cols ...field.Expr) T {
	return b.withDO(b.DO.Group(cols...))
}

func (b genericsDo[T, E]) Having(conds ...Condition) T {
	return b.withDO(b.DO.Having(conds...))
}

func (b genericsDo[T, E]) Limit(limit int) T {
	return b.withDO(b.DO.Limit(limit))
}

func (b genericsDo[T, E]) Offset(offset int) T {
	return b.withDO(b.DO.Offset(offset))
}

func (b genericsDo[T, E]) Scopes(funcs ...func(Dao) Dao) T {
	return b.withDO(b.DO.Scopes(funcs...))
}

func (b genericsDo[T, E]) Unscoped() T {
	return b.withDO(b.DO.Unscoped())
}

func (b genericsDo[T, E]) Create(values ...E) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Create(values)
}

func (b genericsDo[T, E]) CreateInBatches(values []E, batchSize int) error {
	return b.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (b genericsDo[T, E]) Save(values ...E) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Save(values)
}

func (b genericsDo[T, E]) First() (E, error) {
	var e E
	if result, err := b.DO.First(); err != nil {
		return e, err
	} else {
		return result.(E), nil
	}
}

func (b genericsDo[T, E]) Take() (E, error) {
	var e E
	if result, err := b.DO.Take(); err != nil {
		return e, err
	} else {
		return result.(E), nil
	}
}

func (b genericsDo[T, E]) Last() (E, error) {
	var e E
	if result, err := b.DO.Last(); err != nil {
		return e, err
	} else {
		return result.(E), nil
	}
}

func (b genericsDo[T, E]) Find() ([]E, error) {
	result, err := b.DO.Find()
	return result.([]E), err
}

func (b genericsDo[T, E]) FindInBatch(batchSize int, fc func(tx Dao, batch int) error) (results []E, err error) {
	buf := make([]E, 0, batchSize)
	err = b.DO.FindInBatches(&buf, batchSize, func(tx Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (b genericsDo[T, E]) FindInBatches(result *[]E, batchSize int, fc func(tx Dao, batch int) error) error {
	return b.DO.FindInBatches(result, batchSize, fc)
}

func (b genericsDo[T, E]) Attrs(attrs ...field.AssignExpr) T {
	return b.withDO(b.DO.Attrs(attrs...))
}

func (b genericsDo[T, E]) Assign(attrs ...field.AssignExpr) T {
	return b.withDO(b.DO.Assign(attrs...))
}

func (b genericsDo[T, E]) Joins(fields ...field.RelationField) T {
	for _, _f := range fields {
		b.realDo = b.withDO(b.DO.Joins(_f))
	}
	return b.realDo
}

func (b genericsDo[T, E]) Preload(fields ...field.RelationField) T {
	for _, _f := range fields {
		b.realDo = b.withDO(b.DO.Preload(_f))
	}
	return b.realDo
}

func (b genericsDo[T, E]) FirstOrInit() (E, error) {
	var e E
	if result, err := b.DO.FirstOrInit(); err != nil {
		return e, err
	} else {
		return result.(E), nil
	}
}

func (b genericsDo[T, E]) FirstOrCreate() (E, error) {
	var e E
	if result, err := b.DO.FirstOrCreate(); err != nil {
		return e, err
	} else {
		return result.(E), nil
	}
}

func (b genericsDo[T, E]) FindByPage(offset int, limit int) (result []E, count int64, err error) {
	result, err = b.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = b.Offset(-1).Limit(-1).Count()
	return
}

func (b genericsDo[T, E]) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = b.Count()
	if err != nil {
		return
	}
	err = b.Offset(offset).Limit(limit).Scan(result)
	return
}

func (b genericsDo[T, E]) Scan(result interface{}) (err error) {
	return b.DO.Scan(result)
}

func (b genericsDo[T, E]) Delete(models ...E) (result ResultInfo, err error) {
	return b.DO.Delete(models)
}

func (b genericsDo[T, E]) SetDo(do IGenericsDo[T, E]) T {
	return b.realDo.SetDo(do)
}

func (b genericsDo[T, E]) withDO(do Dao) T {
	b.DO = *do.(*DO)
	return b.SetDo(&b)
}
