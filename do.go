package gen

import (
	"database/sql"
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
	Debug = func(db *gorm.DB) *gorm.DB { return db.Debug() }
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
	Emit(methodDiy)
	return d.db
}

// Quote return qutoed data
func (d *DO) Quote(raw string) string {
	return d.db.Statement.Quote(raw)
}

// Build implement the interface of claues.Expression
// only call WHERE clause's Build
func (d *DO) Build(builder clause.Builder) {
	for _, e := range d.buildWhere() {
		e.Build(builder)
	}
}

func (d *DO) buildWhere() []clause.Expression {
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
	// 		stmt.AddClause(clause.Select{})
	// 	}
	// 	return stmt
	// }
)

// buildStmt call statement.Build to combine all clauses in one statement
func (d *DO) buildStmt(opts ...stmtOpt) *gorm.Statement {
	stmt := d.db.Statement
	for _, opt := range opts {
		stmt = opt(stmt)
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

// func (s *DO) subQueryExpr() clause.Expr {
// 	stmt := s.buildStmt(withFROM, withSELECT)
// 	return clause.Expr{SQL: "(" + stmt.SQL.String() + ")", Vars: stmt.Vars}
// }

// Debug return a DO with db in debug mode
func (d *DO) Debug() Dao {
	return NewDO(d.db.Debug())
}

// As alias cannot be heired, As must used on tail
func (d *DO) As(alias string) Dao {
	return &DO{db: d.db, alias: alias}
}

// ======================== chainable api ========================
func (d *DO) Not(conds ...Condition) Dao {
	return NewDO(d.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Not(condToExpression(conds...)...)}}))
}

func (d *DO) Or(conds ...Condition) Dao {
	return NewDO(d.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Or(clause.And(condToExpression(conds...)...))}}))
}

func (d *DO) Select(columns ...field.Expr) Dao {
	Emit(methodSelect)
	if len(columns) == 0 {
		return NewDO(d.db.Clauses(clause.Select{}))
	}
	// return NewDO(d.db.Clauses(clause.Select{Expression: clause.CommaExpression{Exprs: toExpression(columns...)}}))
	return NewDO(d.db.Select(buildExpr(d.db.Statement, columns...)))
}

func (d *DO) Where(conds ...Condition) Dao {
	Emit(methodWhere)
	var exprs = make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		switch cond := cond.(type) {
		case *DO:
			exprs = append(exprs, cond.buildWhere()...)
		default:
			exprs = append(exprs, cond)
		}
	}
	return NewDO(d.db.Clauses(clause.Where{Exprs: exprs}))
}

func (d *DO) Order(columns ...field.Expr) Dao {
	Emit(methodOrder)
	return NewDO(d.db.Clauses(clause.OrderBy{Expression: clause.CommaExpression{Exprs: toExpression(columns...)}}))
}

func (d *DO) Distinct(columns ...field.Expr) Dao {
	Emit(methodDistinct)
	return NewDO(d.db.Distinct(toInterfaceSlice(toColumnFullName(d.db.Statement, columns...))...))
}

func (d *DO) Omit(columns ...field.Expr) Dao {
	Emit(methodOmit)
	return NewDO(d.db.Omit(toColNames(d.db.Statement, columns...)...))
}

func (d *DO) Group(column field.Expr) Dao {
	Emit(methodGroup)
	return NewDO(d.db.Group(column.Column().Name))
}

func (d *DO) Having(conds ...Condition) Dao {
	Emit(methodHaving)
	return NewDO(d.db.Clauses(clause.GroupBy{Having: condToExpression(conds...)}))
}

func (d *DO) Limit(limit int) Dao {
	Emit(methodLimit)
	return NewDO(d.db.Limit(limit))
}

func (d *DO) Offset(offset int) Dao {
	Emit(methodOffset)
	return NewDO(d.db.Offset(offset))
}

func (d *DO) Scopes(funcs ...func(Dao) Dao) Dao {
	Emit(methodScopes)
	var result Dao = d
	for _, f := range funcs {
		result = f(result)
	}
	return result
}

func (d *DO) Unscoped() Dao {
	Emit(methodUnscoped)
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
	Emit(methodJoin)
	var exprs = make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		switch cond := cond.(type) {
		case *DO:
			exprs = append(exprs, cond.buildWhere()...)
		default:
			exprs = append(exprs, cond)
		}
	}

	join := clause.Join{Type: joinType, Table: clause.Table{Name: table.TableName()}, ON: clause.Where{
		Exprs: exprs,
	}}
	from := getFromClause(d.db)
	from.Joins = append(from.Joins, join)
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
	Emit(methodCreate)
	return d.db.Create(value).Error
}

func (d *DO) CreateInBatches(value interface{}, batchSize int) error {
	Emit(methodCreateInBatches)
	return d.db.CreateInBatches(value, batchSize).Error
}

func (d *DO) Save(value interface{}) error {
	Emit(methodSave)
	return d.db.Save(value).Error
}

func (d *DO) First(dest interface{}, conds ...field.Expr) error {
	Emit(methodFirst)
	return d.db.Clauses(toExpression(conds...)...).First(dest).Error
}

func (d *DO) Take(dest interface{}, conds ...field.Expr) error {
	Emit(methodTake)
	return d.db.Clauses(toExpression(conds...)...).Take(dest).Error
}

func (d *DO) Last(dest interface{}, conds ...field.Expr) error {
	Emit(methodLast)
	return d.db.Clauses(toExpression(conds...)...).Last(dest).Error
}

func (d *DO) Find(dest interface{}, conds ...field.Expr) error {
	Emit(methodFind)
	return d.db.Clauses(toExpression(conds...)...).Find(dest).Error
}

func (d *DO) FindInBatches(dest interface{}, batchSize int, fc func(tx Dao, batch int) error) error {
	Emit(methodFindInBatches)
	return d.db.FindInBatches(dest, batchSize, func(tx *gorm.DB, batch int) error { return fc(NewDO(tx), batch) }).Error
}

func (d *DO) FirstOrInit(dest interface{}, conds ...field.Expr) error {
	Emit(methodFirstOrInit)
	return d.db.Clauses(toExpression(conds...)...).FirstOrInit(dest).Error
}

func (d *DO) FirstOrCreate(dest interface{}, conds ...field.Expr) error {
	Emit(methodFirstOrCreate)
	return d.db.Clauses(toExpression(conds...)...).FirstOrCreate(dest).Error
}

func (d *DO) Update(column field.Expr, value interface{}) error {
	Emit(methodUpdate)
	switch expr := column.RawExpr().(type) {
	case clause.Expression:
		return d.db.Update(column.Column().Name, expr).Error
	}

	switch value := value.(type) {
	case field.Expr:
		return d.db.Update(column.Column().Name, value.RawExpr()).Error
	case *DO:
		return d.db.Update(column.Column().Name, value.db).Error
	default:
		return d.db.Update(column.Column().Name, value).Error
	}
}

func (d *DO) Updates(values interface{}) error {
	Emit(methodUpdates)
	return d.db.Updates(values).Error
}

func (d *DO) UpdateColumn(column field.Expr, value interface{}) error {
	Emit(methodUpdateColumn)
	switch expr := column.RawExpr().(type) {
	case clause.Expression:
		return d.db.UpdateColumn(column.Column().Name, expr).Error
	}

	switch value := value.(type) {
	case field.Expr:
		return d.db.UpdateColumn(column.Column().Name, value.RawExpr()).Error
	case *DO:
		return d.db.UpdateColumn(column.Column().Name, value.db).Error
	default:
		return d.db.UpdateColumn(column.Column().Name, value).Error
	}
}

func (d *DO) UpdateColumns(values interface{}) error {
	Emit(methodUpdateColumns)
	return d.db.UpdateColumns(values).Error
}

func (d *DO) Delete(value interface{}, conds ...field.Expr) error {
	Emit(methodDelete)
	return d.db.Clauses(toExpression(conds...)...).Delete(value).Error
}

func (d *DO) Count(count *int64) error {
	Emit(methodCount)
	return d.db.Count(count).Error
}

func (d *DO) Row() *sql.Row {
	Emit(methodRow)
	return d.db.Row()
}

func (d *DO) Rows() (*sql.Rows, error) {
	Emit(methodRows)
	return d.db.Rows()
}

func (d *DO) Scan(dest interface{}) error {
	Emit(methodScan)
	return d.db.Scan(dest).Error
}

func (d *DO) Pluck(column field.Expr, dest interface{}) error {
	Emit(methodPluck)
	return d.db.Pluck(column.Column().Name, dest).Error
}

func (d *DO) ScanRows(rows *sql.Rows, dest interface{}) error {
	Emit(methodScanRows)
	return d.db.ScanRows(rows, dest)
}

func (d *DO) Transaction(fc func(Dao) error, opts ...*sql.TxOptions) error {
	Emit(methodTransaction)
	return d.db.Transaction(func(tx *gorm.DB) error { return fc(NewDO(tx)) }, opts...)
}

func (d *DO) Begin(opts ...*sql.TxOptions) Dao {
	Emit(methodBegin)
	return NewDO(d.db.Begin(opts...))
}

func (d *DO) Commit() Dao {
	Emit(methodCommit)
	return NewDO(d.db.Commit())
}

func (d *DO) RollBack() Dao {
	Emit(methodRollback)
	return NewDO(d.db.Rollback())
}

func (d *DO) SavePoint(name string) Dao {
	Emit(methodSavePoint)
	return NewDO(d.db.SavePoint(name))
}

func (d *DO) RollBackTo(name string) Dao {
	Emit(methodRollbackTo)
	return NewDO(d.db.RollbackTo(name))
}

func toExpression(conds ...field.Expr) []clause.Expression {
	exprs := make([]clause.Expression, len(conds))
	for i, cond := range conds {
		exprs[i] = cond
	}
	return exprs
}

func condToExpression(conds ...Condition) []clause.Expression {
	exprs := make([]clause.Expression, len(conds))
	for i, cond := range conds {
		exprs[i] = cond
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
func Table(subQueries ...Dao) Dao {
	if len(subQueries) == 0 {
		return NewDO(nil)
	}

	tablePlaceholder := make([]string, len(subQueries))
	tableExprs := make([]interface{}, len(subQueries))
	for i, query := range subQueries {
		tablePlaceholder[i] = "(?)"

		do := query.(*DO)
		tableExprs[i] = do.db
		if do.alias != "" {
			tablePlaceholder[i] += " AS " + do.Quote(do.alias)
		}
	}

	db := subQueries[0].(*DO).db
	return NewDO(db.Session(&gorm.Session{NewDB: true}).Table(strings.Join(tablePlaceholder, ", "), tableExprs...))
}

// ======================== sub query method ========================

// In return a subquery expression
// the params conds must contains 2 parameters. Input should be in the order of field expression and sub query
// the function will painc if the last item is not of sub query type
func In(conds ...Condition) field.Expr {
	switch len(conds) {
	case 0, 1:
		return field.ContainsSubQuery(nil, nil)
	default:
		columns := condToExpr(conds[:len(conds)-1]...)
		query := conds[len(conds)-1].(subQuery)
		return field.ContainsSubQuery(columns, query.UnderlyingDB())
	}
}

func Eq(column field.Expr, query subQuery) field.Expr {
	return field.CompareSubQuery(field.EqOp, column, query.UnderlyingDB())
}

func Gt(column field.Expr, query subQuery) field.Expr {
	return field.CompareSubQuery(field.GtOp, column, query.UnderlyingDB())
}

func Gte(column field.Expr, query subQuery) field.Expr {
	return field.CompareSubQuery(field.GteOp, column, query.UnderlyingDB())
}

func Lt(column field.Expr, query subQuery) field.Expr {
	return field.CompareSubQuery(field.LtOp, column, query.UnderlyingDB())
}

func Lte(column field.Expr, query subQuery) field.Expr {
	return field.CompareSubQuery(field.LteOp, column, query.UnderlyingDB())
}

func condToExpr(conds ...Condition) []field.Expr {
	exprs := make([]field.Expr, len(conds))
	for i, cond := range conds {
		exprs[i] = cond.(field.Expr)
	}
	return exprs
}
