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
)

var (
	createClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	queryClauses  = []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR"}
	updateClauses = []string{"UPDATE", "SET", "WHERE"}
	deleteClauses = []string{"DELETE", "FROM", "WHERE"}
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

func (d *DO) ReplaceDB(db *gorm.DB) { d.db = db }

// UseModel specify a data model structure as a source for table name
func (d *DO) UseModel(model interface{}) {
	mt := reflect.TypeOf(model)
	if mt.Kind() == reflect.Ptr {
		mt = mt.Elem()
	}
	d.modelType = mt

	err := d.db.Statement.Parse(model)
	if err != nil {
		panic(fmt.Errorf("Cannot parse model: %+v", model))
	}
	d.schema = d.db.Statement.Schema
}

// UseTable specify table name
func (d *DO) UseTable(tableName string) { d.db = d.db.Table(tableName).Session(new(gorm.Session)) }

// TableName return table name
func (d DO) TableName() string {
	if d.schema == nil {
		return ""
	}
	return d.schema.Table
}

// Session replace db with new session
func (d *DO) Session(config *gorm.Session) Dao { return d.getInstance(d.db.Session(config)) }

// UnderlyingDB return the underlying database connection
func (d *DO) UnderlyingDB() *gorm.DB { return d.db }

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

// implements Condition
func (d *DO) BeCond() interface{} { return d.buildCondition() }
func (d *DO) CondError() error    { return nil }

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
	return &d
}

// Columns return columns for Subquery
func (*DO) Columns(cols ...field.Expr) columns { return cols }

// ======================== chainable api ========================
func (d *DO) Not(conds ...Condition) Dao {
	if len(conds) == 0 {
		return d
	}
	exprs, err := condToExpression(conds)
	if err != nil {
		return d.withError(err)
	}
	return d.getInstance(d.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Not(exprs...)}}))
}

func (d *DO) Or(conds ...Condition) Dao {
	if len(conds) == 0 {
		return d
	}
	exprs, err := condToExpression(conds)
	if err != nil {
		return d.withError(err)
	}
	return d.getInstance(d.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Or(clause.And(exprs...))}}))
}

func (d *DO) Select(columns ...field.Expr) Dao {
	if len(columns) == 0 {
		return d.getInstance(d.db.Clauses(clause.Select{}))
	}
	query, args := buildExpr(d.db.Statement, columns...)
	if len(args) == 0 {
		return d.getInstance(d.db.Select(query))
	}
	return d.getInstance(d.db.Select(strings.Join(query, ","), args...))
}

func (d *DO) Where(conds ...Condition) Dao {
	if len(conds) == 0 {
		return d
	}
	exprs, err := condToExpression(conds)
	if err != nil {
		return d.withError(err)
	}
	return d.getInstance(d.db.Clauses(clause.Where{Exprs: exprs}))
}

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
	return d.getInstance(d.db.Order(d.calcOrderValue(columns...)))
}

func (d *DO) calcOrderValue(columns ...field.Expr) string {
	// eager build Columns
	orderArray := make([]string, len(columns))
	for i, c := range columns {
		orderArray[i] = c.Build(d.db.Statement).String()
	}
	return strings.Join(orderArray, ",")
}

func (d *DO) Distinct(columns ...field.Expr) Dao {
	return d.getInstance(d.db.Distinct(toInterfaceSlice(toColExprFullName(d.db.Statement, columns...))...))
}

func (d *DO) Omit(columns ...field.Expr) Dao {
	if len(columns) == 0 {
		return d
	}
	return d.getInstance(d.db.Omit(getColumnName(columns...)...))
}

func (d *DO) Group(columns ...field.Expr) Dao {
	if len(columns) == 0 {
		return d
	}

	name := columns[0].BuildColumn(d.db.Statement, field.WithTable).String()
	for _, col := range columns[1:] {
		name += "," + col.BuildColumn(d.db.Statement, field.WithTable).String()
	}
	return d.getInstance(d.db.Group(name))
}

func (d *DO) Having(conds ...Condition) Dao {
	if len(conds) == 0 {
		return d
	}

	exprs, err := condToExpression(conds)
	if err != nil {
		return d.withError(err)
	}
	return d.getInstance(d.db.Clauses(clause.GroupBy{Having: exprs}))
}

func (d *DO) Limit(limit int) Dao {
	return d.getInstance(d.db.Limit(limit))
}

func (d *DO) Offset(offset int) Dao {
	return d.getInstance(d.db.Offset(offset))
}

func (d *DO) Scopes(funcs ...func(Dao) Dao) Dao {
	fcs := make([]func(*gorm.DB) *gorm.DB, len(funcs))
	for i, f := range funcs {
		sf := f
		fcs[i] = func(tx *gorm.DB) *gorm.DB { return sf(d.getInstance(tx)).(*DO).db }
	}
	return d.getInstance(d.db.Scopes(fcs...))
}

func (d *DO) Unscoped() Dao {
	return d.getInstance(d.db.Unscoped())
}

func (d *DO) Join(table schema.Tabler, conds ...field.Expr) Dao {
	return d.join(table, clause.InnerJoin, conds)
}

func (d *DO) LeftJoin(table schema.Tabler, conds ...field.Expr) Dao {
	return d.join(table, clause.LeftJoin, conds)
}

func (d *DO) RightJoin(table schema.Tabler, conds ...field.Expr) Dao {
	return d.join(table, clause.RightJoin, conds)
}

func (d *DO) join(table schema.Tabler, joinType clause.JoinType, conds []field.Expr) Dao {
	if len(conds) == 0 {
		return d.withError(ErrEmptyCondition)
	}

	from := getFromClause(d.db)
	from.Joins = append(from.Joins, clause.Join{
		Type:  joinType,
		Table: clause.Table{Name: table.TableName()},
		ON:    clause.Where{Exprs: toExpression(conds...)},
	})
	return d.getInstance(d.db.Clauses(from))
}

func (d *DO) Attrs(attrs ...field.AssignExpr) Dao {
	if len(attrs) == 0 {
		return d
	}
	return d.getInstance(d.db.Attrs(d.attrsValue(attrs)...))
}

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

func (d *DO) Joins(field field.RelationField) Dao {
	return d.getInstance(d.db.Joins(field.Path()))
}

// func (d *DO) Preload(column field.RelationPath, subQuery ...SubQuery) Dao {
// 	if len(subQuery) > 0 {
// 		return d.getInstance(d.db.Preload(string(column.Path()), subQuery[0].underlyingDB()))
// 	}
// 	return d.getInstance(d.db.Preload(string(column.Path())))
// }

func (d *DO) Preload(field field.RelationField) Dao {
	var args []interface{}
	if conds := field.GetConds(); len(conds) > 0 {
		args = append(args, toExpressionInterface(conds...)...)
	}
	if columns := field.GetOrderCol(); len(columns) > 0 {
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Order(d.calcOrderValue(columns...))
		})
	}
	if clauses := field.GetClauses(); len(clauses) > 0 {
		args = append(args, func(db *gorm.DB) *gorm.DB {
			return db.Clauses(clauses...)
		})
	}
	return d.getInstance(d.db.Preload(field.Path(), args...))
}

// UpdateFrom specify update sub query
func (d *DO) UpdateFrom(q subQuery) Dao {
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
func (d *DO) Create(value interface{}) error {
	return d.db.Create(value).Error
}

func (d *DO) CreateInBatches(value interface{}, batchSize int) error {
	return d.db.CreateInBatches(value, batchSize).Error
}

func (d *DO) Save(value interface{}) error {
	return d.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(value).Error
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

func (d *DO) FindInBatches(dest interface{}, batchSize int, fc func(tx Dao, batch int) error) error {
	return d.db.FindInBatches(dest, batchSize, func(tx *gorm.DB, batch int) error { return fc(d.getInstance(tx), batch) }).Error
}

func (d *DO) FirstOrInit() (result interface{}, err error) {
	return d.singleQuery(d.db.FirstOrInit)
}

func (d *DO) FirstOrCreate() (result interface{}, err error) {
	return d.singleQuery(d.db.FirstOrCreate)
}

func (d *DO) Update(column field.Expr, value interface{}) (info ResultInfo, err error) {
	tx := d.db.Model(d.newResultPointer())
	columnStr := column.BuildColumn(d.db.Statement, field.WithoutQuote).String()

	var result *gorm.DB
	switch value := value.(type) {
	case field.AssignExpr:
		result = tx.Update(columnStr, value.AssignExpr())
	case subQuery:
		result = tx.Update(columnStr, value.underlyingDB())
	default:
		result = tx.Update(columnStr, value)
	}
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

func (d *DO) UpdateSimple(columns ...field.AssignExpr) (info ResultInfo, err error) {
	if len(columns) == 0 {
		return
	}

	result := d.db.Model(d.newResultPointer()).Clauses(d.assignSet(columns)).Omit("*").Updates(map[string]interface{}{})
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

func (d *DO) Updates(value interface{}) (info ResultInfo, err error) {
	result := d.db.Model(d.newResultPointer()).Updates(value)
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

func (d *DO) UpdateColumn(column field.Expr, value interface{}) (info ResultInfo, err error) {
	tx := d.db.Model(d.newResultPointer())
	columnStr := column.BuildColumn(d.db.Statement, field.WithoutQuote).String()

	var result *gorm.DB
	switch value := value.(type) {
	case field.Expr:
		result = tx.UpdateColumn(columnStr, value.RawExpr())
	case subQuery:
		result = d.db.UpdateColumn(columnStr, value.underlyingDB())
	default:
		result = d.db.UpdateColumn(columnStr, value)
	}
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

func (d *DO) UpdateColumnSimple(columns ...field.AssignExpr) (info ResultInfo, err error) {
	if len(columns) == 0 {
		return
	}

	result := d.db.Model(d.newResultPointer()).Clauses(d.assignSet(columns)).Omit("*").UpdateColumns(map[string]interface{}{})
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

func (d *DO) UpdateColumns(value interface{}) (info ResultInfo, err error) {
	result := d.db.Model(d.newResultPointer()).UpdateColumns(value)
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

// assignSet fetch all set
func (d *DO) assignSet(exprs []field.AssignExpr) (set clause.Set) {
	for _, expr := range exprs {
		column := clause.Column{Name: string(expr.ColumnName())}
		if d.alias != "" {
			column.Table = d.alias
		}
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

func (d *DO) Delete() (info ResultInfo, err error) {
	result := d.db.Model(d.newResultPointer()).Delete(reflect.New(d.modelType).Interface())
	return ResultInfo{RowsAffected: result.RowsAffected, Error: result.Error}, result.Error
}

func (d *DO) Count() (count int64, err error) {
	return count, d.db.Session(&gorm.Session{}).Model(d.newResultPointer()).Count(&count).Error
}

func (d *DO) Row() *sql.Row {
	return d.db.Model(d.newResultPointer()).Row()
}

func (d *DO) Rows() (*sql.Rows, error) {
	return d.db.Model(d.newResultPointer()).Rows()
}

func (d *DO) Scan(dest interface{}) error {
	return d.db.Model(d.newResultPointer()).Scan(dest).Error
}

func (d *DO) Pluck(column field.Expr, dest interface{}) error {
	return d.db.Model(d.newResultPointer()).Pluck(column.ColumnName().String(), dest).Error
}

func (d *DO) ScanRows(rows *sql.Rows, dest interface{}) error {
	return d.db.Model(d.newResultPointer()).ScanRows(rows, dest)
}

func (d *DO) newResultPointer() interface{} {
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

func buildExpr(stmt *gorm.Statement, exprs ...field.Expr) (query []string, args []interface{}) {
	for _, e := range exprs {
		sql, vars := e.BuildWithArgs(stmt)
		query = append(query, sql.String())
		args = append(args, vars...)
	}
	return query, args
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
// 	Table(u.Select(u.ID, u.Name).Where(u.Age.Gt(18))).Select()
// the above usage is equivalent to SQL statement:
// 	SELECT * FROM (SELECT `id`, `name` FROM `users_info` WHERE `age` > ?)"
func Table(subQueries ...subQuery) Dao {
	if len(subQueries) == 0 {
		return &DO{}
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

	return &DO{
		db: subQueries[0].underlyingDO().db.Session(&gorm.Session{NewDB: true}).
			Table(strings.Join(tablePlaceholder, ", "), tableExprs...),
	}
}

// ======================== sub query method ========================

type columns []field.Expr

// Set assign value by subquery
func (cs columns) Set(query subQuery) field.AssignExpr {
	return field.AssignSubQuery(cs, query.underlyingDB())
}

// In accept query or value
func (cs columns) In(queryOrValue Condition) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}

	switch query := queryOrValue.(type) {
	case field.Value:
		return field.ContainsValue(cs, query)
	case subQuery:
		return field.ContainsSubQuery(cs, query.underlyingDB())
	default:
		return field.EmptyExpr()
	}
}

func (cs columns) NotIn(queryOrValue Condition) field.Expr {
	return field.Not(cs.In(queryOrValue))
}

func (cs columns) Eq(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.EqOp, cs[0], query.underlyingDB())
}

func (cs columns) Neq(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.NeqOp, cs[0], query.underlyingDB())
}

func (cs columns) Gt(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.GtOp, cs[0], query.underlyingDB())
}

func (cs columns) Gte(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.GteOp, cs[0], query.underlyingDB())
}

func (cs columns) Lt(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.LtOp, cs[0], query.underlyingDB())
}

func (cs columns) Lte(query subQuery) field.Expr {
	if len(cs) == 0 {
		return field.EmptyExpr()
	}
	return field.CompareSubQuery(field.LteOp, cs[0], query.underlyingDB())
}
