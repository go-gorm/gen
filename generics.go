package gen

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// IGenericsDo generic query interface
type IGenericsDo[T any, E any] interface {
	SubQuery
	Debug() T
	WithContext(ctx context.Context) T
	WithResult(fc func(tx Dao)) ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() T
	WriteDB() T
	As(alias string) Dao
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
	Rows() (*sql.Rows, error)
	Row() *sql.Row
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) T
	UnderlyingDB() *gorm.DB
	schema.Tabler
	ToSQL(queryFn func(T)) string
}

// IWithDO is a generic interface for types that can be associated with a DAO (Data Access Object)
// It enables method chaining by returning the concrete implementation type
type IWithDO[T any] interface {
	// WithDO associates a DAO instance with the implementing type
	//   - do: The DAO (Data Access Object) to be used for data operations
	// Returns a new instance of the implementing type with the DAO set
	// This enables fluent-style method chaining while maintaining type safety
	WithDO(do Dao) T
}

// GenericsDo base impl of IGenericsDo
type GenericsDo[T IGenericsDo[T, E], E any] struct {
	DO
	IWithDO[T]
}

// Debug ...
func (b GenericsDo[T, E]) Debug() T {
	return b.withDO(b.DO.Debug())
}

// WithContext ...
func (b GenericsDo[T, E]) WithContext(ctx context.Context) T {
	return b.withDO(b.DO.WithContext(ctx))
}

// ReadDB ...
func (b GenericsDo[T, E]) ReadDB() T {
	return b.Clauses(dbresolver.Read)
}

// WriteDB ...
func (b GenericsDo[T, E]) WriteDB() T {
	return b.Clauses(dbresolver.Write)
}

// Session ...
func (b GenericsDo[T, E]) Session(config *gorm.Session) T {
	return b.withDO(b.DO.Session(config))
}

// Clauses ...
func (b GenericsDo[T, E]) Clauses(conds ...clause.Expression) T {
	return b.withDO(b.DO.Clauses(conds...))
}

// Returning ...
func (b GenericsDo[T, E]) Returning(value interface{}, columns ...string) T {
	return b.withDO(b.DO.Returning(value, columns...))
}

// Not ...
func (b GenericsDo[T, E]) Not(conds ...Condition) T {
	return b.withDO(b.DO.Not(conds...))
}

// Or ...
func (b GenericsDo[T, E]) Or(conds ...Condition) T {
	return b.withDO(b.DO.Or(conds...))
}

// Select ...
func (b GenericsDo[T, E]) Select(conds ...field.Expr) T {
	return b.withDO(b.DO.Select(conds...))
}

// Where ...
func (b GenericsDo[T, E]) Where(conds ...Condition) T {
	return b.withDO(b.DO.Where(conds...))
}

// Order ...
func (b GenericsDo[T, E]) Order(conds ...field.Expr) T {
	return b.withDO(b.DO.Order(conds...))
}

// Distinct ...
func (b GenericsDo[T, E]) Distinct(cols ...field.Expr) T {
	return b.withDO(b.DO.Distinct(cols...))
}

// Omit ...
func (b GenericsDo[T, E]) Omit(cols ...field.Expr) T {
	return b.withDO(b.DO.Omit(cols...))
}

// Join ...
func (b GenericsDo[T, E]) Join(table schema.Tabler, on ...field.Expr) T {
	return b.withDO(b.DO.Join(table, on...))
}

// LeftJoin ...
func (b GenericsDo[T, E]) LeftJoin(table schema.Tabler, on ...field.Expr) T {
	return b.withDO(b.DO.LeftJoin(table, on...))
}

// RightJoin ...
func (b GenericsDo[T, E]) RightJoin(table schema.Tabler, on ...field.Expr) T {
	return b.withDO(b.DO.RightJoin(table, on...))
}

// Group ...
func (b GenericsDo[T, E]) Group(cols ...field.Expr) T {
	return b.withDO(b.DO.Group(cols...))
}

// Having ...
func (b GenericsDo[T, E]) Having(conds ...Condition) T {
	return b.withDO(b.DO.Having(conds...))
}

// Limit ...
func (b GenericsDo[T, E]) Limit(limit int) T {
	return b.withDO(b.DO.Limit(limit))
}

// Offset ...
func (b GenericsDo[T, E]) Offset(offset int) T {
	return b.withDO(b.DO.Offset(offset))
}

// Scopes ...
func (b GenericsDo[T, E]) Scopes(funcs ...func(Dao) Dao) T {
	return b.withDO(b.DO.Scopes(funcs...))
}

// Unscoped ...
func (b GenericsDo[T, E]) Unscoped() T {
	return b.withDO(b.DO.Unscoped())
}

// Create ...
func (b GenericsDo[T, E]) Create(values ...E) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Create(values)
}

// CreateInBatches ...
func (b GenericsDo[T, E]) CreateInBatches(values []E, batchSize int) error {
	return b.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (b GenericsDo[T, E]) Save(values ...E) error {
	if len(values) == 0 {
		return nil
	}
	return b.DO.Save(values)
}

// First ...
func (b GenericsDo[T, E]) First() (E, error) {
	var e E
	result, err := b.DO.First()
	if err != nil {
		return e, err
	}
	v, ok := result.(E)
	if !ok {
		return e, fmt.Errorf("gen: First type mismatch, expected %T, got %T", new(E), result)
	}
	return v, nil
}

// Take ...
func (b GenericsDo[T, E]) Take() (E, error) {
	var e E
	result, err := b.DO.Take()
	if err != nil {
		return e, err
	}
	v, ok := result.(E)
	if !ok {
		return e, fmt.Errorf("gen: Take type mismatch, expected %T, got %T", new(E), result)
	}
	return v, nil
}

// Last ...
func (b GenericsDo[T, E]) Last() (E, error) {
	var e E
	result, err := b.DO.Last()
	if err != nil {
		return e, err
	}
	v, ok := result.(E)
	if !ok {
		return e, fmt.Errorf("gen: Last type mismatch, expected %T, got %T", new(E), result)
	}
	return v, nil
}

// Find ...
func (b GenericsDo[T, E]) Find() ([]E, error) {
	result, err := b.DO.Find()
	if err != nil {
		return nil, err
	}
	v, ok := result.([]E)
	if !ok {
		return nil, fmt.Errorf("gen: Find type mismatch, expected []E, got %T", result)
	}
	return v, nil
}

// FindInBatch ...
func (b GenericsDo[T, E]) FindInBatch(batchSize int, fc func(tx Dao, batch int) error) (results []E, err error) {
	if batchSize <= 0 {
		return nil, fmt.Errorf("gen: invalid batchSize %d, must be > 0", batchSize)
	}
	buf := make([]E, 0, batchSize)
	err = b.DO.FindInBatches(&buf, batchSize, func(tx Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

// FindInBatches ...
func (b GenericsDo[T, E]) FindInBatches(result *[]E, batchSize int, fc func(tx Dao, batch int) error) error {
	return b.DO.FindInBatches(result, batchSize, fc)
}

// Attrs ...
func (b GenericsDo[T, E]) Attrs(attrs ...field.AssignExpr) T {
	return b.withDO(b.DO.Attrs(attrs...))
}

// Assign ...
func (b GenericsDo[T, E]) Assign(attrs ...field.AssignExpr) T {
	return b.withDO(b.DO.Assign(attrs...))
}

// Joins ...
func (b GenericsDo[T, E]) Joins(fields ...field.RelationField) T {
	var do Dao = &b.DO
	for _, _f := range fields {
		do = do.Joins(_f)
	}
	return b.withDO(do)
}

// Preload ...
func (b GenericsDo[T, E]) Preload(fields ...field.RelationField) T {
	var do Dao = &b.DO
	for _, _f := range fields {
		do = do.Preload(_f)
	}
	return b.withDO(do)
}

// FirstOrInit ...
func (b GenericsDo[T, E]) FirstOrInit() (E, error) {
	var e E
	result, err := b.DO.FirstOrInit()
	if err != nil {
		return e, err
	}
	v, ok := result.(E)
	if !ok {
		return e, fmt.Errorf("gen: FirstOrInit type mismatch, expected %T, got %T", new(E), result)
	}
	return v, nil
}

// FirstOrCreate ...
func (b GenericsDo[T, E]) FirstOrCreate() (E, error) {
	var e E
	result, err := b.DO.FirstOrCreate()
	if err != nil {
		return e, err
	}
	v, ok := result.(E)
	if !ok {
		return e, fmt.Errorf("gen: FirstOrCreate type mismatch, expected %T, got %T", new(E), result)
	}
	return v, nil
}

// FindByPage ...
func (b GenericsDo[T, E]) FindByPage(offset int, limit int) (result []E, count int64, err error) {
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

// ScanByPage ...
func (b GenericsDo[T, E]) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = b.Count()
	if err != nil {
		return
	}
	err = b.Offset(offset).Limit(limit).Scan(result)
	return
}

// Scan ...
func (b GenericsDo[T, E]) Scan(result interface{}) (err error) {
	return b.DO.Scan(result)
}

// Delete ...
func (b GenericsDo[T, E]) Delete(models ...E) (result ResultInfo, err error) {
	return b.DO.Delete(models)
}

// ToSQL ...
func (b GenericsDo[T, E]) ToSQL(queryFn func(T)) string {
	b.db = b.db.Session(&gorm.Session{DryRun: true, SkipDefaultTransaction: true, NewDB: true}).Table(b.TableName())
	queryFn(b.withDO(&b.DO))
	db := b.db
	stmt := db.Statement
	return db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
}

func (b *GenericsDo[T, E]) withDO(do Dao) T {
	if b.IWithDO == nil {
		panic("gen: IWithDO is nil")
	}
	return b.WithDO(do)
}

// WithDOFunc is a function type that implements the IWithDO interface through
// functional options. It provides a mechanism to inject DAO dependencies while
// maintaining type safety with generics.
//
// The generic type T should be constrained to pointer receivers of concrete
// types to ensure proper value semantics and immutability.
type WithDOFunc[T any] func(do Dao) T

// WithDO executes the receiver function with the provided DAO instance,
// returning a new instance of type T. This method enables fluent chaining
// while preserving immutability by design.
//
// Parameters:
//
//	do - Data Access Object instance (non-nil) to be injected
//
// Returns:
//
//	New instance of type T initialized with the DAO dependency
//
// Example:
//
//	constructor := func(do Dao) *UserService {
//	    return &UserService{dao: do}
//	}
//	service := WithDOFunc[*UserService](constructor).WithDO(db)
func (b WithDOFunc[T]) WithDO(do Dao) T {
	return b(do)
}
