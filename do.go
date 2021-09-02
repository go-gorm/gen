package gen

import (
	"database/sql"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen/field"
)

var (
	createClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	queryClauses  = []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR"}
	updateClauses = []string{"UPDATE", "SET", "WHERE"}
	deleteClauses = []string{"DELETE", "FROM", "WHERE"}
)

func NewDO(db *gorm.DB) *DO { return &DO{db: db} }

// DO (data object): implement basic query methods
// the structure embedded with a *gorm.DB, and has a element item "alias" will be used when used as a sub query
type DO struct {
	db    *gorm.DB
	alias string // for subquery
}

type doOptions func(*gorm.DB) *gorm.DB

var (
	// Debug use DB in debug mode
	Debug doOptions = func(db *gorm.DB) *gorm.DB { return db.Debug() }
)

// UseDB specify a db connection(*gorm.DB)
func (d *DO) UseDB(db *gorm.DB, opts ...doOptions) {
	db = db.Session(new(gorm.Session))
	for _, opt := range opts {
		db = opt(db)
	}
	d.db = db
}

// UseModel specify a data model structure as a source for table name
func (d *DO) UseModel(model interface{}) {
	d.db = d.db.Model(model).Session(new(gorm.Session))
	_ = d.db.Statement.Parse(model)
}

// UseTable specify table name
func (d *DO) UseTable(tableName string) {
	d.db = d.db.Table(tableName).Session(new(gorm.Session))
}

// TableName return table name
func (d *DO) TableName() string {
	return d.db.Statement.Table
}

// UnderlyingDB return the underlying database connection
func (d *DO) UnderlyingDB() *gorm.DB {
	return d.db
}

// Quote return qutoed data
func (d *DO) Quote(raw string) string {
	return d.db.Statement.Quote(raw)
}

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

type stmtOpt func(*gorm.Statement) *gorm.Statement

var (
	// withFROM add FROM clause
	withFROM stmtOpt = func(stmt *gorm.Statement) *gorm.Statement {
		if stmt.Table == "" {
			_ = stmt.Parse(stmt.Model)
		}
		stmt.AddClause(clause.From{})
		return stmt
	}

	// // withSELECT add SELECT clause
	// withSELECT stmtOpt = func(stmt *gorm.Statement) *gorm.Statement {
	// 	if _, ok := stmt.Clauses["SELECT"]; !ok {
	// 		stmt.AddClause(clause.Select{Distinct: stmt.Distinct})
	// 	}
	// 	return stmt
	// }
)

// build FOR TEST. call statement.Build to combine all clauses in one statement
func (d *DO) build(opts ...stmtOpt) *gorm.Statement {
	stmt := d.db.Statement
	for _, opt := range opts {
		stmt = opt(stmt)
	}

	if _, ok := stmt.Clauses["SELECT"]; !ok && len(stmt.Selects) > 0 {
		stmt.AddClause(clause.Select{Distinct: stmt.Distinct, Expression: clause.Expr{SQL: strings.Join(stmt.Selects, ",")}})
	}

	findClauses := func() []string {
		for _, cs := range [][]string{createClauses, queryClauses, updateClauses, deleteClauses} {
			if _, ok := stmt.Clauses[cs[0]]; ok {
				return cs
			}
		}
		return queryClauses
	}

	stmt.Build(findClauses()...)
	return stmt
}

// underlyingDO return self
func (d *DO) underlyingDO() *DO { return d }

// Debug return a DO with db in debug mode
func (d *DO) Debug() Dao {
	return NewDO(d.db.Debug())
}

func (d *DO) Hints(hs ...Hint) Dao {
	return NewDO(d.db.Clauses(hintToExpression(hs)...))
}

// As alias cannot be heired, As must used on tail
func (d *DO) As(alias string) Dao {
	return &DO{db: d.db, alias: alias}
}

// Columns return columns for Subquery
func (*DO) Columns(cols ...field.Expr) columns {
	return cols
}

// ======================== chainable api ========================
func (d *DO) Not(conds ...Condition) Dao {
	return NewDO(d.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Not(condToExpression(conds)...)}}))
}

func (d *DO) Or(conds ...Condition) Dao {
	return NewDO(d.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Or(clause.And(condToExpression(conds)...))}}))
}

func (d *DO) Select(columns ...field.Expr) Dao {
	if len(columns) == 0 {
		return NewDO(d.db.Clauses(clause.Select{}))
	}
	return NewDO(d.db.Select(buildExpr(d.db.Statement, columns...)))
}

func (d *DO) Where(conds ...Condition) Dao {
	return NewDO(d.db.Clauses(clause.Where{Exprs: condToExpression(conds)}))
}

func (d *DO) Order(columns ...field.Expr) Dao {
	// lazy build Columns
	// if c, ok := d.db.Statement.Clauses[clause.OrderBy{}.Name()]; ok {
	// 	if order, ok := c.Expression.(clause.OrderBy); ok {
	// 		if expr, ok := order.Expression.(clause.CommaExpression); ok {
	// 			expr.Exprs = append(expr.Exprs, toExpression(columns)...)
	// 			return NewDO(d.db.Clauses(clause.OrderBy{Expression: expr}))
	// 		}
	// 	}
	// }
	// return NewDO(d.db.Clauses(clause.OrderBy{Expression: clause.CommaExpression{Exprs: toExpression(columns)}}))

	// eager build Columns
	orderArray := make([]string, len(columns))
	for i, c := range columns {
		orderArray[i] = c.BuildExpr(d.db.Statement)
	}
	return NewDO(d.db.Order(strings.Join(orderArray, ",")))
}

func (d *DO) Distinct(columns ...field.Expr) Dao {
	return NewDO(d.db.Distinct(toInterfaceSlice(toColumnFullName(d.db.Statement, columns...))...))
}

func (d *DO) Omit(columns ...field.Expr) Dao {
	return NewDO(d.db.Omit(toColNames(d.db.Statement, columns...)...))
}

func (d *DO) Group(column field.Expr) Dao {
	return NewDO(d.db.Group(column.Column().Name))
}

func (d *DO) Having(conds ...Condition) Dao {
	return NewDO(d.db.Clauses(clause.GroupBy{Having: condToExpression(conds)}))
}

func (d *DO) Limit(limit int) Dao {
	return NewDO(d.db.Limit(limit))
}

func (d *DO) Offset(offset int) Dao {
	return NewDO(d.db.Offset(offset))
}

func (d *DO) Scopes(funcs ...func(Dao) Dao) Dao {
	fcs := make([]func(*gorm.DB) *gorm.DB, len(funcs))
	for i, f := range funcs {
		fcs[i] = func(tx *gorm.DB) *gorm.DB { return f(NewDO(tx)).(*DO).db }
	}
	return NewDO(d.db.Scopes(fcs...))
}

func (d *DO) Unscoped() Dao {
	return NewDO(d.db.Unscoped())
}

func (d *DO) Join(table schema.Tabler, conds ...Condition) Dao {
	return d.join(table, clause.InnerJoin, conds...)
}

func (d *DO) LeftJoin(table schema.Tabler, conds ...Condition) Dao {
	return d.join(table, clause.LeftJoin, conds...)
}

func (d *DO) RightJoin(table schema.Tabler, conds ...Condition) Dao {
	return d.join(table, clause.RightJoin, conds...)
}

func (d *DO) join(table schema.Tabler, joinType clause.JoinType, conds ...Condition) Dao {
	from := getFromClause(d.db)
	from.Joins = append(from.Joins, clause.Join{
		Type:  joinType,
		Table: clause.Table{Name: table.TableName()},
		ON:    clause.Where{Exprs: condToExpression(conds)},
	})
	return NewDO(d.db.Clauses(from))
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
func (d *DO) Create(value interface{}) error {
	return d.db.Create(value).Error
}

func (d *DO) CreateInBatches(value interface{}, batchSize int) error {
	return d.db.CreateInBatches(value, batchSize).Error
}

func (d *DO) Save(value interface{}) error {
	return d.db.Save(value).Error
}

func (d *DO) First() (result interface{}, err error) {
	return d.singleQuery(d.db.First)
}

func (d *DO) Take() (result interface{}, err error) {
	return d.singleQuery(d.db.Take)
}

func (d *DO) Last() (result interface{}, err error) {
	return d.singleQuery(d.db.Last)
}

func (d *DO) singleQuery(query func(dest interface{}, conds ...interface{}) *gorm.DB) (result interface{}, err error) {
	result = d.newResultPointer()
	if err := query(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DO) Find() (results interface{}, err error) {
	return d.multiQuery(d.db.Find)
}

func (d *DO) multiQuery(query func(dest interface{}, conds ...interface{}) *gorm.DB) (results interface{}, err error) {
	resultsPtr := d.newResultSlicePointer()
	err = query(resultsPtr).Error
	return reflect.Indirect(reflect.ValueOf(resultsPtr)).Interface(), err
}

func (d *DO) FindInBatches(dest interface{}, batchSize int, fc func(tx Dao, batch int) error) error {
	return d.db.FindInBatches(dest, batchSize, func(tx *gorm.DB, batch int) error { return fc(NewDO(tx), batch) }).Error
}

// func (d *DO) FirstOrInit(dest interface{}, conds ...field.Expr) error {
// 	return d.db.Clauses(toExpression(conds)...).FirstOrInit(dest).Error
// }

// func (d *DO) FirstOrCreate(dest interface{}, conds ...field.Expr) error {
// 	return d.db.Clauses(toExpression(conds)...).FirstOrCreate(dest).Error
// }

func (d *DO) Update(column field.Expr, value interface{}) error {
	switch expr := column.RawExpr().(type) {
	case clause.Expression:
		return d.db.Update(column.Column().Name, expr).Error
	}

	switch value := value.(type) {
	case field.Expr:
		return d.db.Update(column.Column().Name, value.RawExpr()).Error
	case subQuery:
		return d.db.Update(column.Column().Name, value.UnderlyingDB()).Error
	default:
		return d.db.Update(column.Column().Name, value).Error
	}
}

func (d *DO) Updates(values interface{}) error {
	return d.db.Updates(values).Error
}

func (d *DO) UpdateColumn(column field.Expr, value interface{}) error {
	switch expr := column.RawExpr().(type) {
	case clause.Expression:
		return d.db.UpdateColumn(column.Column().Name, expr).Error
	}

	switch value := value.(type) {
	case field.Expr:
		return d.db.UpdateColumn(column.Column().Name, value.RawExpr()).Error
	case subQuery:
		return d.db.UpdateColumn(column.Column().Name, value.UnderlyingDB()).Error
	default:
		return d.db.UpdateColumn(column.Column().Name, value).Error
	}
}

func (d *DO) UpdateColumns(values interface{}) error {
	return d.db.UpdateColumns(values).Error
}

func (d *DO) Delete() error {
	return d.db.Delete(d.db.Statement.Model).Error
}

func (d *DO) Count() (count int64, err error) {
	return count, d.db.Count(&count).Error
}

func (d *DO) Row() *sql.Row {
	return d.db.Row()
}

func (d *DO) Rows() (*sql.Rows, error) {
	return d.db.Rows()
}

func (d *DO) Scan(dest interface{}) error {
	return d.db.Scan(dest).Error
}

func (d *DO) Pluck(column field.Expr, dest interface{}) error {
	return d.db.Pluck(column.Column().Name, dest).Error
}

func (d *DO) ScanRows(rows *sql.Rows, dest interface{}) error {
	return d.db.ScanRows(rows, dest)
}

func (d *DO) Transaction(fc func(Dao) error, opts ...*sql.TxOptions) error {
	return d.db.Transaction(func(tx *gorm.DB) error { return fc(NewDO(tx)) }, opts...)
}

func (d *DO) Begin(opts ...*sql.TxOptions) Dao {
	return NewDO(d.db.Begin(opts...))
}

func (d *DO) Commit() Dao {
	return NewDO(d.db.Commit())
}

func (d *DO) RollBack() Dao {
	return NewDO(d.db.Rollback())
}

func (d *DO) SavePoint(name string) Dao {
	return NewDO(d.db.SavePoint(name))
}

func (d *DO) RollBackTo(name string) Dao {
	return NewDO(d.db.RollbackTo(name))
}

func (d *DO) newResultPointer() interface{} {
	return reflect.New(d.getModel()).Interface()
}

func (d *DO) newResultSlicePointer() interface{} {
	return reflect.New(reflect.SliceOf(reflect.PtrTo(d.getModel()))).Interface()
}

func (d *DO) getModel() reflect.Type {
	return reflect.Indirect(reflect.ValueOf(d.db.Statement.Model)).Type()
}

func hintToExpression(hs []Hint) []clause.Expression {
	exprs := make([]clause.Expression, len(hs))
	for i, hint := range hs {
		exprs[i] = hint
	}
	return exprs
}

func condToExpression(conds []Condition) []clause.Expression {
	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		switch cond := cond.(type) {
		case subQuery:
			exprs = append(exprs, cond.underlyingDO().buildCondition()...)
		default:
			exprs = append(exprs, cond)
		}
	}
	return exprs
}

func toColumnFullName(stmt *gorm.Statement, columns ...field.Expr) []string {
	return buildColumn(stmt, columns, field.WithAll)
}

func toColNames(stmt *gorm.Statement, columns ...field.Expr) []string {
	return buildColumn(stmt, columns)
}

func buildColumn(stmt *gorm.Statement, cols []field.Expr, opts ...field.BuildOpt) []string {
	results := make([]string, len(cols))
	for i, c := range cols {
		results[i] = c.BuildColumn(stmt, opts...)
	}
	return results
}

func buildExpr(stmt *gorm.Statement, exprs ...field.Expr) []string {
	results := make([]string, len(exprs))
	for i, e := range exprs {
		results[i] = e.BuildExpr(stmt)
	}
	return results
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
// 	Table(u.Select(u.ID, u.Name).Where(u.Age.Gt(18))).Select()
// the above usage is equivalent to SQL statement:
// 	SELECT * FROM (SELECT `id`, `name` FROM `users_info` WHERE `age` > ?)"
func Table(subQueries ...subQuery) Dao {
	if len(subQueries) == 0 {
		return NewDO(nil)
	}

	tablePlaceholder := make([]string, len(subQueries))
	tableExprs := make([]interface{}, len(subQueries))
	for i, query := range subQueries {
		tablePlaceholder[i] = "(?)"

		do := query.underlyingDO()
		tableExprs[i] = do.db
		if do.alias != "" {
			tablePlaceholder[i] += " AS " + do.Quote(do.alias)
		}
	}

	return NewDO(subQueries[0].underlyingDO().db.
		Session(&gorm.Session{NewDB: true}).
		Table(strings.Join(tablePlaceholder, ", "), tableExprs...))
}

// ======================== sub query method ========================

type columns []field.Expr

// In accept query or value
func (cs columns) In(queryOrValue Condition) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}

	query, ok := queryOrValue.(subQuery)
	if !ok {
		return field.ContainsValue(cs, queryOrValue)
	}
	return field.ContainsSubQuery(cs, query.UnderlyingDB())
}

func (cs columns) NotIn(queryOrValue Condition) field.Expr {
	return field.Not(cs.In(queryOrValue))
}

func (cs columns) Eq(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.EqOp, cs[0], query.UnderlyingDB())
}

func (cs columns) Neq(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.NeqOp, cs[0], query.UnderlyingDB())
}

func (cs columns) Gt(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.GtOp, cs[0], query.UnderlyingDB())
}

func (cs columns) Gte(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.GteOp, cs[0], query.UnderlyingDB())
}

func (cs columns) Lt(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.LtOp, cs[0], query.UnderlyingDB())
}

func (cs columns) Lte(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.LteOp, cs[0], query.UnderlyingDB())
}
