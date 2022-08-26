package gen

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen/field"
	"gorm.io/gen/helper"
)

// ResultInfo query/execute info
type ResultInfo struct {
	RowsAffected int64
	Error        error
}

var _ Dao = new(DO)

// DO (data object): implement basic query methods
// the structure embedded with a *gorm.DB, and has a element item "alias" will be used when used as a sub query
type DO struct {
	db        *gorm.DB
	alias     string // for subquery
	modelType reflect.Type
	schema    *schema.Schema

	backfillData interface{}
}

func (d DO) getInstance(db *gorm.DB) *DO {
	d.db = db
	return &d
}

type doOptions func(*gorm.DB) *gorm.DB

var (
	// Debug use DB in debug mode
	Debug doOptions = func(db *gorm.DB) *gorm.DB { return db.Debug() }
)

// UseDB specify a db connection(*gorm.DB)
func (d *DO) UseDB(db *gorm.DB, opts ...doOptions) {
	db = db.Session(&gorm.Session{Context: context.Background()})
	for _, opt := range opts {
		db = opt(db)
	}
	d.db = db
}

// ReplaceDB replace db connection
func (d *DO) ReplaceDB(db *gorm.DB) { d.db = db }

// UseModel specify a data model structure as a source for table name
func (d *DO) UseModel(model interface{}) {
	d.modelType = d.indirect(model)

	err := d.db.Statement.Parse(model)
	if err != nil {
		panic(fmt.Errorf("Cannot parse model: %+v\n%w", model, err))
	}
	d.schema = d.db.Statement.Schema
}

func (d *DO) indirect(value interface{}) reflect.Type {
	mt := reflect.TypeOf(value)
	if mt.Kind() == reflect.Ptr {
		mt = mt.Elem()
	}
	return mt
}

// UseTable specify table name
func (d *DO) UseTable(tableName string) {
	d.db = d.db.Table(tableName).Session(new(gorm.Session))
	d.schema.Table = tableName
}

// TableName return table name
func (d DO) TableName() string {
	if d.schema == nil {
		return ""
	}
	return d.schema.Table
}

// Returning backfill data
func (d DO) Returning(value interface{}, columns ...string) Dao {
	d.backfillData = value

	var targetCulumns []clause.Column
	for _, column := range columns {
		targetCulumns = append(targetCulumns, clause.Column{Name: column})
	}
	d.db = d.db.Clauses(clause.Returning{Columns: targetCulumns})
	return &d
}

// Session replace db with new session
func (d *DO) Session(config *gorm.Session) Dao { return d.getInstance(d.db.Session(config)) }

// UnderlyingDB return the underlying database connection
func (d *DO) UnderlyingDB() *gorm.DB { return d.underlyingDB() }

// Quote return qutoed data
func (d *DO) Quote(raw string) string { return d.db.Statement.Quote(raw) }

// Build implement the interface of claues.Expression
// only call WHERE clause's Build
func (d *DO) Build(builder clause.Builder) {
	for _, e := range d.buildCondition() {
		e.Build(builder)
	}
}

func (d *DO) buildCondition() []clause.Expression {
	return d.db.Statement.BuildCondition(d.db)
}

// underlyingDO return self
func (d *DO) underlyingDO() *DO { return d }

// underlyingDB return self.db
func (d *DO) underlyingDB() *gorm.DB { return d.db }

func (d *DO) withError(err error) *DO {
	if err == nil {
		return d
	}

	newDB := d.db.Session(new(gorm.Session))
	_ = newDB.AddError(err)
	return d.getInstance(newDB)
}

// BeCond implements Condition
func (d *DO) BeCond() interface{} { return d.buildCondition() }

// CondError implements Condition
func (d *DO) CondError() error { return nil }

// Debug return a DO with db in debug mode
func (d *DO) Debug() Dao { return d.getInstance(d.db.Debug()) }

// WithContext return a DO with db with context
func (d *DO) WithContext(ctx context.Context) Dao { return d.getInstance(d.db.WithContext(ctx)) }

// Clauses specify Clauses
func (d *DO) Clauses(conds ...clause.Expression) Dao {
	if err := checkConds(conds); err != nil {
		newDB := d.db.Session(new(gorm.Session))
		_ = newDB.AddError(err)
		return d.getInstance(newDB)
	}
	return d.getInstance(d.db.Clauses(conds...))
}

// As alias cannot be heired, As must used on tail
func (d DO) As(alias string) Dao {
	d.alias = alias
	d.db = d.db.Table(fmt.Sprintf("%s AS %s", d.Quote(d.TableName()), d.Quote(alias)))
	return &d
}

// Alias return alias name
func (d *DO) Alias() string { return d.alias }

// Columns return columns for Subquery
func (*DO) Columns(cols ...field.Expr) Columns { return cols }

// ======================== chainable api ========================

// Not ...
func (d *DO) Not(conds ...Condition) Dao {
	exprs, err := condToExpression(conds)
	if err != nil {
		return d.withError(err)
	}
	if len(exprs) == 0 {
		return d
	}
	return d.getInstance(d.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Not(exprs...)}}))
}

// Or ...
func (d *DO) Or(conds ...Condition) Dao {
	exprs, err := condToExpression(conds)
	if err != nil {
		return d.withError(err)
	}
	if len(exprs) == 0 {
		return d
	}
	return d.getInstance(d.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Or(clause.And(exprs...))}}))
}

// Select ...
func (d *DO) Select(columns ...field.Expr) Dao {
	if len(columns) == 0 {
		return d.getInstance(d.db.Clauses(clause.Select{}))
	}
	query, args := buildExpr4Select(d.db.Statement, columns...)
	return d.getInstance(d.db.Select(query, args...))
}

// Where ...
func (d *DO) Where(conds ...Condition) Dao {
	exprs, err := condToExpression(conds)
	if err != nil {
		return d.withError(err)
	}
	if len(exprs) == 0 {
		return d
	}
	return d.getInstance(d.db.Clauses(clause.Where{Exprs: exprs}))
}

// Order ...
func (d *DO) Order(columns ...field.Expr) Dao {
	// lazy build Columns
	// if c, ok := d.db.Statement.Clauses[clause.OrderBy{}.Name()]; ok {
	// 	if order, ok := c.Expression.(clause.OrderBy); ok {
	// 		if expr, ok := order.Expression.(clause.CommaExpression); ok {
	// 			expr.Exprs = append(expr.Exprs, toExpression(columns)...)
	// 			return d.newInstance(d.db.Clauses(clause.OrderBy{Expression: expr}))
	// 		}
	// 	}
	// }
	// return d.newInstance(d.db.Clauses(clause.OrderBy{Expression: clause.CommaExpression{Exprs: toExpression(columns)}}))
	if len(columns) == 0 {
		return d
	}
	return d.getInstance(d.db.Order(d.toOrderValue(columns...)))
}

func (d *DO) toOrderValue(columns ...field.Expr) string {
	// eager build Columns
	orderArray := make([]string, len(columns))
	for i, c := range columns {
		orderArray[i] = c.Build(d.db.Statement).String()
	}
	return strings.Join(orderArray, ",")
}

// Distinct ...
func (d *DO) Distinct(columns ...field.Expr) Dao {
	return d.getInstance(d.db.Distinct(toInterfaceSlice(toColExprFullName(d.db.Statement, columns...))...))
}

// Omit ...
func (d *DO) Omit(columns ...field.Expr) Dao {
	if len(columns) == 0 {
		return d
	}
	return d.getInstance(d.db.Omit(getColumnName(columns...)...))
}

// Group ...
func (d *DO) Group(columns ...field.Expr) Dao {
	if len(columns) == 0 {
		return d
	}
	name := string(columns[0].Build(d.db.Statement))
	for _, col := range columns[1:] {
		name += "," + string(col.Build(d.db.Statement))
	}
	return d.getInstance(d.db.Group(name))
}

// Having ...
func (d *DO) Having(conds ...Condition) Dao {
	exprs, err := condToExpression(conds)
	if err != nil {
		return d.withError(err)
	}
	if len(exprs) == 0 {
		return d
	}
	return d.getInstance(d.db.Clauses(clause.GroupBy{Having: exprs}))
}

// Limit ...
func (d *DO) Limit(limit int) Dao {
	return d.getInstance(d.db.Limit(limit))
}

// Offset ...
func (d *DO) Offset(offset int) Dao {
	return d.getInstance(d.db.Offset(offset))
}

// Scopes ...
func (d *DO) Scopes(funcs ...func(Dao) Dao) Dao {
	fcs := make([]func(*gorm.DB) *gorm.DB, len(funcs))
	for i, f := range funcs {
		sf := f
		fcs[i] = func(tx *gorm.DB) *gorm.DB { return sf(d.getInstance(tx)).(*DO).db }
	}
	return d.getInstance(d.db.Scopes(fcs...))
}

// Unscoped ...
func (d *DO) Unscoped() Dao {
	return d.getInstance(d.db.Unscoped())
}

// Join ...
func (d *DO) Join(table schema.Tabler, conds ...field.Expr) Dao {
	return d.join(table, clause.InnerJoin, conds)
}

// LeftJoin ...
func (d *DO) LeftJoin(table schema.Tabler, conds ...field.Expr) Dao {
	return d.join(table, clause.LeftJoin, conds)
}

// RightJoin ...
func (d *DO) RightJoin(table schema.Tabler, conds ...field.Expr) Dao {
	return d.join(table, clause.RightJoin, conds)
}

func (d *DO) join(table schema.Tabler, joinType clause.JoinType, conds []field.Expr) Dao {
	if len(conds) == 0 {
		return d.withError(ErrEmptyCondition)
	}

	join := clause.Join{
		Type:  joinType,
		Table: clause.Table{Name: table.TableName()},
		ON:    clause.Where{Exprs: toExpression(conds...)},
	}
	if do, ok := table.(Dao); ok {
		join.Expression = helper.NewJoinTblExpr(join, Table(do).underlyingDB().Statement.TableExpr)
	}
	if al, ok := table.(interface{ Alias() string }); ok {
		join.Table.Alias = al.Alias()
	}

	from := getFromClause(d.db)
	from.Joins = append(from.Joins, join)
	return d.getInstance(d.db.Clauses(from))
}

// Attrs ...
func (d *DO) Attrs(attrs ...field.AssignExpr) Dao {
	if len(attrs) == 0 {
		return d
	}
	return d.getInstance(d.db.Attrs(d.attrsValue(attrs)...))
}

// Assign ...
func (d *DO) Assign(attrs ...field.AssignExpr) Dao {
	if len(attrs) == 0 {
		return d
	}
	return d.getInstance(d.db.Assign(d.attrsValue(attrs)...))
}

func (d *DO) attrsValue(attrs []field.AssignExpr) []interface{} {
	values := make([]interface{}, 0, len(attrs))
	for _, attr := range attrs {
		if expr, ok := attr.AssignExpr().(clause.Eq); ok {
			values = append(values, expr)
		}
	}
	return values
}

// Joins ...
func (d *DO) Joins(field field.RelationField) Dao {
	var args []interface{}

	if conds := field.GetConds(); len(conds) > 0 {
		var exprs []clause.Expression
		for _, oe := range toExpression(conds...) {
			switch e := oe.(type) {
			case clause.Eq:
				if c, ok := e.Column.(clause.Column); ok {
					c.Table = field.Name()
					e.Column = c
				}
				exprs = append(exprs, e)
			case clause.Neq:
				if c, ok := e.Column.(clause.Column); ok {
					c.Table = field.Name()
					e.Column = c
				}
				exprs = append(exprs, e)
			case clause.Gt:
				if c, ok := e.Column.(clause.Column); ok {
					c.Table = field.Name()
					e.Column = c
				}
				exprs = append(exprs, e)
			case clause.Gte:
				if c, ok := e.Column.(clause.Column); ok {
					c.Table = field.Name()
					e.Column = c
				}
				exprs = append(exprs, e)
			case clause.Lt:
				if c, ok := e.Column.(clause.Column); ok {
					c.Table = field.Name()
					e.Column = c
				}
				exprs = append(exprs, e)
			case clause.Lte:
				if c, ok := e.Column.(clause.Column); ok {
					c.Table = field.Name()
					e.Column = c
				}
				exprs = append(exprs, e)
			case clause.Like:
				if c, ok := e.Column.(clause.Column); ok {
					c.Table = field.Name()
					e.Column = c
				}
				exprs = append(exprs, e)
			}
		}

		args = append(args, d.db.Clauses(clause.Where{
			Exprs: exprs,
		}))
	}
	if columns := field.GetSelects(); len(columns) > 0 {
		colNames := make([]string, len(columns))
		for i, c := range columns {
			colNames[i] = string(c.ColumnName())
		}
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Select(colNames)
		})
	}
	if columns := field.GetOrderCol(); len(columns) > 0 {
		var os []string
		for _, oe := range columns {
			switch e := oe.RawExpr().(type) {
			case clause.Expr:
				vs := []interface{}{}
				for _, v := range e.Vars {
					if c, ok := v.(clause.Column); ok {
						vs = append(vs, clause.Column{
							Table: field.Name(),
							Name:  c.Name,
							Alias: c.Alias,
							Raw:   c.Raw,
						})
					}
				}
				e.Vars = vs
				newStmt := &gorm.Statement{DB: d.db.Statement.DB, Table: d.db.Statement.Table, Schema: d.db.Statement.Schema}
				e.Build(newStmt)
				os = append(os, newStmt.SQL.String())
			}
		}
		args = append(args, d.db.Order(strings.Join(os, ",")))
	}
	if clauses := field.GetClauses(); len(clauses) > 0 {
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clauses...)
		})
	}
	if funcs := field.GetScopes(); len(funcs) > 0 {
		for _, f := range funcs {
			args = append(args, (func(*gorm.DB) *gorm.DB)(f))
		}
	}
	if offset, limit := field.GetPage(); offset|limit != 0 {
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Offset(offset).Limit(limit)
		})
	}

	return d.getInstance(d.db.Joins(field.Path(), args...))
}

// func (d *DO) Preload(column field.RelationPath, subQuery ...SubQuery) Dao {
// 	if len(subQuery) > 0 {
// 		return d.getInstance(d.db.Preload(string(column.Path()), subQuery[0].underlyingDB()))
// 	}
// 	return d.getInstance(d.db.Preload(string(column.Path())))
// }

// Preload ...
func (d *DO) Preload(field field.RelationField) Dao {
	var args []interface{}
	if conds := field.GetConds(); len(conds) > 0 {
		args = append(args, toExpressionInterface(conds...)...)
	}
	if columns := field.GetSelects(); len(columns) > 0 {
		colNames := make([]string, len(columns))
		for i, c := range columns {
			colNames[i] = string(c.ColumnName())
		}
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Select(colNames)
		})
	}
	if columns := field.GetOrderCol(); len(columns) > 0 {
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Order(d.toOrderValue(columns...))
		})
	}
	if clauses := field.GetClauses(); len(clauses) > 0 {
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clauses...)
		})
	}
	if funcs := field.GetScopes(); len(funcs) > 0 {
		for _, f := range funcs {
			args = append(args, (func(*gorm.DB) *gorm.DB)(f))
		}
	}
	if offset, limit := field.GetPage(); offset|limit != 0 {
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Offset(offset).Limit(limit)
		})
	}
	return d.getInstance(d.db.Preload(field.Path(), args...))
}

// UpdateFrom specify update sub query
func (d *DO) UpdateFrom(q SubQuery) Dao {
	var tableName strings.Builder
	d.db.Statement.QuoteTo(&tableName, d.TableName())
	if d.alias != "" {
		tableName.WriteString(" AS ")
		d.db.Statement.QuoteTo(&tableName, d.alias)
	}

	tableName.WriteByte(',')
	if _, ok := q.underlyingDB().Statement.Clauses["SELECT"]; ok || len(q.underlyingDB().Statement.Selects) > 0 {
		tableName.WriteString("(" + q.underlyingDB().ToSQL(func(tx *gorm.DB) *gorm.DB { return tx.Table(q.underlyingDO().TableName()).Find(nil) }) + ")")
	} else {
		d.db.Statement.QuoteTo(&tableName, q.underlyingDO().TableName())
	}
	if alias := q.underlyingDO().alias; alias != "" {
		tableName.WriteString(" AS ")
		d.db.Statement.QuoteTo(&tableName, alias)
	}

	return d.getInstance(d.db.Clauses(clause.Update{Table: clause.Table{Name: tableName.String(), Raw: true}}))
}

func getFromClause(db *gorm.DB) *clause.From {
	if db == nil || db.Statement == nil {
		return &clause.From{}
	}
	c, ok := db.Statement.Clauses[clause.From{}.Name()]
	if !ok || c.Expression == nil {
		return &clause.From{}
	}
	from, ok := c.Expression.(clause.From)
	if !ok {
		return &clause.From{}
	}
	return &from
}

// ======================== finisher api ========================

// Create ...
func (d *DO) Create(value interface{}) error {
	return d.db.Create(value).Error
}

// CreateInBatches ...
func (d *DO) CreateInBatches(value interface{}, batchSize int) error {
	return d.db.CreateInBatches(value, batchSize).Error
}

// Save ...
func (d *DO) Save(value interface{}) error {
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(value).Error
}

// First ...
func (d *DO) First() (result interface{}, err error) {
	return d.singleQuery(d.db.First)
}

// Take ...
func (d *DO) Take() (result interface{}, err error) {
	return d.singleQuery(d.db.Take)
}

// Last ...
func (d *DO) Last() (result interface{}, err error) {
	return d.singleQuery(d.db.Last)
}

func (d *DO) singleQuery(query func(dest interface{}, conds ...interface{}) *gorm.DB) (result interface{}, err error) {
	if d.modelType == nil {
		return d.singleScan()
	}

	result = d.newResultPointer()
	if err := query(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DO) singleScan() (result interface{}, err error) {
	result = map[string]interface{}{}
	err = d.db.Scan(result).Error
	return
}

// Find ...
func (d *DO) Find() (results interface{}, err error) {
	return d.multiQuery(d.db.Find)
}

func (d *DO) multiQuery(query func(dest interface{}, conds ...interface{}) *gorm.DB) (results interface{}, err error) {
	if d.modelType == nil {
		return d.findToMap()
	}

	resultsPtr := d.newResultSlicePointer()
	err = query(resultsPtr).Error
	return reflect.Indirect(reflect.ValueOf(resultsPtr)).Interface(), err
}

func (d *DO) findToMap() (interface{}, error) {
	var results []map[string]interface{}
	err := d.db.Find(&results).Error
	return results, err
}

// FindInBatches ...
func (d *DO) FindInBatches(dest interface{}, batchSize int, fc func(tx Dao, batch int) error) error {
	return d.db.FindInBatches(dest, batchSize, func(tx *gorm.DB, batch int) error { return fc(d.getInstance(tx), batch) }).Error
}

// FirstOrInit ...
func (d *DO) FirstOrInit() (result interface{}, err error) {
	return d.singleQuery(d.db.FirstOrInit)
}

// FirstOrCreate ...
func (d *DO) FirstOrCreate() (result interface{}, err error) {
	return d.singleQuery(d.db.FirstOrCreate)
}

// Update ...
func (d *DO) Update(column field.Expr, value interface{}) (info ResultInfo, err error) {
	tx := d.db.Model(d.newResultPointer())
	columnStr := column.BuildColumn(d.db.Statement, field.WithoutQuote).String()

	var result *gorm.DB
	switch value := value.(type) {
	case field.AssignExpr:
		result = tx.Update(columnStr, value.AssignExpr())
	case SubQuery:
		result = tx.Update(columnStr, value.underlyingDB())
	default:
		result = tx.Update(columnStr, value)
	}
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

// UpdateSimple ...
func (d *DO) UpdateSimple(columns ...field.AssignExpr) (info ResultInfo, err error) {
	if len(columns) == 0 {
		return
	}

	result := d.db.Model(d.newResultPointer()).Clauses(d.assignSet(columns)).Omit("*").Updates(map[string]interface{}{})
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

// Updates ...
func (d *DO) Updates(value interface{}) (info ResultInfo, err error) {
	var rawTyp, valTyp reflect.Type

	rawTyp = reflect.TypeOf(value)
	if rawTyp.Kind() == reflect.Ptr {
		valTyp = rawTyp.Elem()
	} else {
		valTyp = rawTyp
	}

	tx := d.db
	if d.backfillData != nil {
		tx = tx.Model(d.backfillData)
	}
	switch {
	case valTyp != d.modelType: // different type with model
		if d.backfillData == nil {
			tx = tx.Model(d.newResultPointer())
		}
	case rawTyp.Kind() == reflect.Ptr: // ignore ptr value
	default: // for fixing "reflect.Value.Addr of unaddressable value" panic
		ptr := reflect.New(d.modelType)
		ptr.Elem().Set(reflect.ValueOf(value))
		value = ptr.Interface()
	}
	result := tx.Updates(value)
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

// UpdateColumn ...
func (d *DO) UpdateColumn(column field.Expr, value interface{}) (info ResultInfo, err error) {
	tx := d.db.Model(d.newResultPointer())
	columnStr := column.BuildColumn(d.db.Statement, field.WithoutQuote).String()

	var result *gorm.DB
	switch value := value.(type) {
	case field.Expr:
		result = tx.UpdateColumn(columnStr, value.RawExpr())
	case SubQuery:
		result = d.db.UpdateColumn(columnStr, value.underlyingDB())
	default:
		result = d.db.UpdateColumn(columnStr, value)
	}
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

// UpdateColumnSimple ...
func (d *DO) UpdateColumnSimple(columns ...field.AssignExpr) (info ResultInfo, err error) {
	if len(columns) == 0 {
		return
	}

	result := d.db.Model(d.newResultPointer()).Clauses(d.assignSet(columns)).Omit("*").UpdateColumns(map[string]interface{}{})
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

// UpdateColumns ...
func (d *DO) UpdateColumns(value interface{}) (info ResultInfo, err error) {
	result := d.db.Model(d.newResultPointer()).UpdateColumns(value)
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

// assignSet fetch all set
func (d *DO) assignSet(exprs []field.AssignExpr) (set clause.Set) {
	for _, expr := range exprs {
		column := clause.Column{Table: d.alias, Name: string(expr.ColumnName())}
		switch e := expr.AssignExpr().(type) {
		case clause.Expr:
			set = append(set, clause.Assignment{Column: column, Value: e})
		case clause.Eq:
			set = append(set, clause.Assignment{Column: column, Value: e.Value})
		case clause.Set:
			set = append(set, e...)
		}
	}

	stmt := d.db.Session(&gorm.Session{}).Statement
	stmt.Dest = map[string]interface{}{}
	return append(set, callbacks.ConvertToAssignments(stmt)...)
}

// Delete ...
func (d *DO) Delete(models ...interface{}) (info ResultInfo, err error) {
	var result *gorm.DB
	if len(models) == 0 || reflect.ValueOf(models[0]).Len() == 0 {
		result = d.db.Model(d.newResultPointer()).Delete(reflect.New(d.modelType).Interface())
	} else {
		targets := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(d.modelType)), 0, len(models))
		value := reflect.ValueOf(models[0])
		for i := 0; i < value.Len(); i++ {
			targets = reflect.Append(targets, value.Index(i))
		}
		result = d.db.Delete(targets.Interface())
	}
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

// Count ...
func (d *DO) Count() (count int64, err error) {
	return count, d.db.Session(&gorm.Session{}).Model(d.newResultPointer()).Count(&count).Error
}

// Row ...
func (d *DO) Row() *sql.Row {
	return d.db.Model(d.newResultPointer()).Row()
}

// Rows ...
func (d *DO) Rows() (*sql.Rows, error) {
	return d.db.Model(d.newResultPointer()).Rows()
}

// Scan ...
func (d *DO) Scan(dest interface{}) error {
	return d.db.Model(d.newResultPointer()).Scan(dest).Error
}

// Pluck ...
func (d *DO) Pluck(column field.Expr, dest interface{}) error {
	return d.db.Model(d.newResultPointer()).Pluck(column.ColumnName().String(), dest).Error
}

// ScanRows ...
func (d *DO) ScanRows(rows *sql.Rows, dest interface{}) error {
	return d.db.Model(d.newResultPointer()).ScanRows(rows, dest)
}

// WithResult ...
func (d DO) WithResult(fc func(tx Dao)) ResultInfo {
	d.db = d.db.Set("", "")
	fc(&d)
	return ResultInfo{RowsAffected: d.db.RowsAffected, Error: d.db.Error}
}

func (d *DO) newResultPointer() interface{} {
	if d.backfillData != nil {
		return d.backfillData
	}
	if d.modelType == nil {
		return nil
	}
	return reflect.New(d.modelType).Interface()
}

func (d *DO) newResultSlicePointer() interface{} {
	return reflect.New(reflect.SliceOf(reflect.PtrTo(d.modelType))).Interface()
}

func toColExprFullName(stmt *gorm.Statement, columns ...field.Expr) []string {
	return buildColExpr(stmt, columns, field.WithAll)
}

func getColumnName(columns ...field.Expr) (result []string) {
	for _, c := range columns {
		result = append(result, c.ColumnName().String())
	}
	return result
}

func buildColExpr(stmt *gorm.Statement, cols []field.Expr, opts ...field.BuildOpt) []string {
	results := make([]string, len(cols))
	for i, c := range cols {
		switch c.RawExpr().(type) {
		case clause.Column:
			results[i] = c.BuildColumn(stmt, opts...).String()
		case clause.Expression:
			sql, args := c.BuildWithArgs(stmt)
			results[i] = stmt.Dialector.Explain(sql.String(), args...)
		}
	}
	return results
}

func buildExpr4Select(stmt *gorm.Statement, exprs ...field.Expr) (query string, args []interface{}) {
	if len(exprs) == 0 {
		return "", nil
	}

	var queryItems []string
	for _, e := range exprs {
		sql, vars := e.BuildWithArgs(stmt)
		queryItems = append(queryItems, sql.String())
		args = append(args, vars...)
	}
	if len(args) == 0 {
		return queryItems[0], toInterfaceSlice(queryItems[1:])
	}
	return strings.Join(queryItems, ","), args
}

func toExpression(exprs ...field.Expr) []clause.Expression {
	result := make([]clause.Expression, len(exprs))
	for i, e := range exprs {
		result[i] = singleExpr(e)
	}
	return result
}

func toExpressionInterface(exprs ...field.Expr) []interface{} {
	result := make([]interface{}, len(exprs))
	for i, e := range exprs {
		result[i] = singleExpr(e)
	}
	return result
}

func singleExpr(e field.Expr) clause.Expression {
	switch v := e.RawExpr().(type) {
	case clause.Expression:
		return v
	case clause.Column:
		return clause.NamedExpr{SQL: "?", Vars: []interface{}{v}}
	default:
		return clause.Expr{}
	}
}

func toInterfaceSlice(value interface{}) []interface{} {
	switch v := value.(type) {
	case string:
		return []interface{}{v}
	case []string:
		res := make([]interface{}, len(v))
		for i, item := range v {
			res[i] = item
		}
		return res
	case []clause.Column:
		res := make([]interface{}, len(v))
		for i, item := range v {
			res[i] = item
		}
		return res
	default:
		return nil
	}
}

// ======================== New Table ========================

// Table return a new table produced by subquery,
// the return value has to be used as root node
//
//	Table(u.Select(u.ID, u.Name).Where(u.Age.Gt(18))).Select()
//
// the above usage is equivalent to SQL statement:
//
//	SELECT * FROM (SELECT `id`, `name` FROM `users_info` WHERE `age` > ?)"
func Table(subQueries ...SubQuery) Dao {
	if len(subQueries) == 0 {
		return &DO{}
	}

	tablePlaceholder := make([]string, len(subQueries))
	tableExprs := make([]interface{}, len(subQueries))
	for i, query := range subQueries {
		tablePlaceholder[i] = "(?)"

		do := query.underlyingDO()
		// ignore alias, or will misuse with sub query alias
		tableExprs[i] = do.db.Table(do.TableName())
		if do.alias != "" {
			tablePlaceholder[i] += " AS " + do.Quote(do.alias)
		}
	}

	return &DO{
		db: subQueries[0].underlyingDO().db.Session(&gorm.Session{NewDB: true}).
			Table(strings.Join(tablePlaceholder, ", "), tableExprs...),
	}
}

// ======================== sub query method ========================

// Columns columns array
type Columns []field.Expr

// Set assign value by subquery
func (cs Columns) Set(query SubQuery) field.AssignExpr {
	return field.AssignSubQuery(cs, query.underlyingDB())
}

// In accept query or value
func (cs Columns) In(queryOrValue Condition) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}

	switch query := queryOrValue.(type) {
	case field.Value:
		return field.ContainsValue(cs, query)
	case SubQuery:
		return field.ContainsSubQuery(cs, query.underlyingDB())
	default:
		return field.EmptyExpr()
	}
}

// NotIn ...
func (cs Columns) NotIn(queryOrValue Condition) field.Expr {
	return field.Not(cs.In(queryOrValue))
}

// Eq ...
func (cs Columns) Eq(query SubQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.EqOp, cs[0], query.underlyingDB())
}

// Neq ...
func (cs Columns) Neq(query SubQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.NeqOp, cs[0], query.underlyingDB())
}

// Gt ...
func (cs Columns) Gt(query SubQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.GtOp, cs[0], query.underlyingDB())
}

// Gte ...
func (cs Columns) Gte(query SubQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.GteOp, cs[0], query.underlyingDB())
}

// Lt ...
func (cs Columns) Lt(query SubQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.LtOp, cs[0], query.underlyingDB())
}

// Lte ...
func (cs Columns) Lte(query SubQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.LteOp, cs[0], query.underlyingDB())
}
