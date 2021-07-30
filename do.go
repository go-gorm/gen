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
func (s *DO) UseDB(db *gorm.DB, opts ...doOptions) {
	db = db.Session(new(gorm.Session))
	for _, opt := range opts {
		db = opt(db)
	}
	s.db = db
}

// UseModel specify a data model structure as a source for table name
func (s *DO) UseModel(model interface{}) {
	s.db = s.db.Model(model).Session(new(gorm.Session))
	_ = s.db.Statement.Parse(model)
}

// UseTable specify table name
func (s *DO) UseTable(tableName string) {
	s.db = s.db.Table(tableName).Session(new(gorm.Session))
}

// Table return table name
func (s *DO) TableName() string {
	return s.db.Statement.Table
}

// UnderlyingDB return the underlying database connection
func (s *DO) UnderlyingDB() *gorm.DB {
	Emit(methodDiy)
	return s.db
}

// Build implement the interface of claues.Expression
// only call WHERE clause's Build
func (s *DO) Build(builder clause.Builder) {
	for _, e := range s.buildWhere() {
		e.Build(builder)
	}
}

func (s *DO) buildWhere() []clause.Expression {
	return s.db.Statement.BuildCondition(s.db)
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
func (s *DO) buildStmt(opts ...stmtOpt) *gorm.Statement {
	stmt := s.db.Statement
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

// As alias cannot be heired, As must used on tail
func (s *DO) As(alias string) Dao {
	return &DO{db: s.db, alias: alias}
}

// ======================== chainable api ========================
func (s *DO) Not(conds ...Condition) Dao {
	return NewDO(s.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Not(condToExpression(conds...)...)}}))
}

func (s *DO) Or(conds ...Condition) Dao {
	return NewDO(s.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Or(clause.And(condToExpression(conds...)...))}}))
}

func (s *DO) Select(columns ...field.Expr) Dao {
	Emit(methodSelect)
	if len(columns) == 0 {
		return NewDO(s.db.Clauses(clause.Select{}))
	}
	return NewDO(s.db.Clauses(clause.Select{Expression: CommaExpression{Exprs: toExpression(columns...)}}))
}

func (s *DO) Where(conds ...Condition) Dao {
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
	return NewDO(s.db.Clauses(clause.Where{Exprs: exprs}))
}

func (s *DO) Order(columns ...field.Expr) Dao {
	Emit(methodOrder)
	return NewDO(s.db.Clauses(clause.OrderBy{Expression: CommaExpression{Exprs: toExpression(columns...)}}))
}

func (s *DO) Distinct(columns ...field.Expr) Dao {
	Emit(methodDistinct)
	return NewDO(s.db.Distinct(toInterfaceSlice(toColNames(s.db.Statement, columns...))))
}

func (s *DO) Omit(columns ...field.Expr) Dao {
	Emit(methodOmit)
	return NewDO(s.db.Omit(toColNames(s.db.Statement, columns...)...))
}

func (s *DO) Group(column field.Expr) Dao {
	Emit(methodGroup)
	return NewDO(s.db.Group(column.Column().Name))
}

func (s *DO) Having(conds ...Condition) Dao {
	Emit(methodHaving)
	return NewDO(s.db.Clauses(clause.GroupBy{Having: condToExpression(conds...)}))
}

func (s *DO) Limit(limit int) Dao {
	Emit(methodLimit)
	return NewDO(s.db.Limit(limit))
}

func (s *DO) Offset(offset int) Dao {
	Emit(methodOffset)
	return NewDO(s.db.Offset(offset))
}

func (s *DO) Scopes(funcs ...func(Dao) Dao) Dao {
	Emit(methodScopes)
	var result Dao = s
	for _, f := range funcs {
		result = f(result)
	}
	return result
}

func (s *DO) Unscoped() Dao {
	Emit(methodUnscoped)
	return NewDO(s.db.Unscoped())
}

func (s *DO) Join(table schema.Tabler, conds ...Condition) Dao {
	return s.join(table, clause.InnerJoin, conds...)
}

func (s *DO) LeftJoin(table schema.Tabler, conds ...Condition) Dao {
	return s.join(table, clause.LeftJoin, conds...)
}

func (s *DO) RightJoin(table schema.Tabler, conds ...Condition) Dao {
	return s.join(table, clause.RightJoin, conds...)
}

func (s *DO) join(table schema.Tabler, joinType clause.JoinType, conds ...Condition) Dao {
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
	from := getFromClause(s.db)
	if from == nil {
		from = &clause.From{}
	}
	from.Joins = append(from.Joins, join)
	return NewDO(s.db.Clauses(from))
}

func getFromClause(db *gorm.DB) *clause.From {
	if db == nil {
		return nil
	}
	c, ok := db.Statement.Clauses[clause.From{}.Name()]
	if !ok || c.Expression == nil {
		return nil
	}
	if from, ok := c.Expression.(clause.From); ok {
		return &from
	}
	return nil
}

// ======================== finisher api ========================
func (s *DO) Create(value interface{}) error {
	Emit(methodCreate)
	return s.db.Create(value).Error
}

func (s *DO) CreateInBatches(value interface{}, batchSize int) error {
	Emit(methodCreateInBatches)
	return s.db.CreateInBatches(value, batchSize).Error
}

func (s *DO) Save(value interface{}) error {
	Emit(methodSave)
	return s.db.Save(value).Error
}

func (s *DO) First(dest interface{}, conds ...field.Expr) error {
	Emit(methodFirst)
	return s.db.Clauses(toExpression(conds...)...).First(dest).Error
}

func (s *DO) Take(dest interface{}, conds ...field.Expr) error {
	Emit(methodTake)
	return s.db.Clauses(toExpression(conds...)...).Take(dest).Error
}

func (s *DO) Last(dest interface{}, conds ...field.Expr) error {
	Emit(methodLast)
	return s.db.Clauses(toExpression(conds...)...).Last(dest).Error
}

func (s *DO) Find(dest interface{}, conds ...field.Expr) error {
	Emit(methodFind)
	return s.db.Clauses(toExpression(conds...)...).Find(dest).Error
}

func (s *DO) FindInBatches(dest interface{}, batchSize int, fc func(tx Dao, batch int) error) error {
	Emit(methodFindInBatches)
	return s.db.FindInBatches(dest, batchSize, func(tx *gorm.DB, batch int) error { return fc(NewDO(tx), batch) }).Error
}

func (s *DO) FirstOrInit(dest interface{}, conds ...field.Expr) error {
	Emit(methodFirstOrInit)
	return s.db.Clauses(toExpression(conds...)...).FirstOrInit(dest).Error
}

func (s *DO) FirstOrCreate(dest interface{}, conds ...field.Expr) error {
	Emit(methodFirstOrCreate)
	return s.db.Clauses(toExpression(conds...)...).FirstOrCreate(dest).Error
}

func (s *DO) Update(column field.Expr, value interface{}) error {
	Emit(methodUpdate)
	switch expr := column.RawExpr().(type) {
	case clause.Expression:
		return s.db.Update(column.Column().Name, expr).Error
	}

	switch value := value.(type) {
	case field.Expr:
		return s.db.Update(column.Column().Name, value.RawExpr()).Error
	case *DO:
		return s.db.Update(column.Column().Name, value.db).Error
	default:
		return s.db.Update(column.Column().Name, value).Error
	}
}

func (s *DO) Updates(values interface{}) error {
	Emit(methodUpdates)
	return s.db.Updates(values).Error
}

func (s *DO) UpdateColumn(column field.Expr, value interface{}) error {
	Emit(methodUpdateColumn)
	switch expr := column.RawExpr().(type) {
	case clause.Expression:
		return s.db.UpdateColumn(column.Column().Name, expr).Error
	}

	switch value := value.(type) {
	case field.Expr:
		return s.db.UpdateColumn(column.Column().Name, value.RawExpr()).Error
	case *DO:
		return s.db.UpdateColumn(column.Column().Name, value.db).Error
	default:
		return s.db.UpdateColumn(column.Column().Name, value).Error
	}
}

func (s *DO) UpdateColumns(values interface{}) error {
	Emit(methodUpdateColumns)
	return s.db.UpdateColumns(values).Error
}

func (s *DO) Delete(value interface{}, conds ...field.Expr) error {
	Emit(methodDelete)
	return s.db.Clauses(toExpression(conds...)...).Delete(value).Error
}

func (s *DO) Count(count *int64) error {
	Emit(methodCount)
	return s.db.Count(count).Error
}

func (s *DO) Row() *sql.Row {
	Emit(methodRow)
	return s.db.Row()
}

func (s *DO) Rows() (*sql.Rows, error) {
	Emit(methodRows)
	return s.db.Rows()
}

func (s *DO) Scan(dest interface{}) error {
	Emit(methodScan)
	return s.db.Scan(dest).Error
}

func (s *DO) Pluck(column field.Expr, dest interface{}) error {
	Emit(methodPluck)
	return s.db.Pluck(column.Column().Name, dest).Error
}

func (s *DO) ScanRows(rows *sql.Rows, dest interface{}) error {
	Emit(methodScanRows)
	return s.db.ScanRows(rows, dest)
}

func (s *DO) Transaction(fc func(Dao) error, opts ...*sql.TxOptions) error {
	Emit(methodTransaction)
	return s.db.Transaction(func(tx *gorm.DB) error { return fc(NewDO(tx)) }, opts...)
}

func (s *DO) Begin(opts ...*sql.TxOptions) Dao {
	Emit(methodBegin)
	return NewDO(s.db.Begin(opts...))
}

func (s *DO) Commit() Dao {
	Emit(methodCommit)
	return NewDO(s.db.Commit())
}

func (s *DO) RollBack() Dao {
	Emit(methodRollback)
	return NewDO(s.db.Rollback())
}

func (s *DO) SavePoint(name string) Dao {
	Emit(methodSavePoint)
	return NewDO(s.db.SavePoint(name))
}

func (s *DO) RollBackTo(name string) Dao {
	Emit(methodRollbackTo)
	return NewDO(s.db.RollbackTo(name))
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

func toColNames(stmt *gorm.Statement, columns ...field.Expr) []string {
	names := make([]string, len(columns))
	for i, col := range columns {
		names[i] = col.BuildColumn(stmt)
	}
	return names
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
	default:
		return nil
	}
}

// ======================== temporary ========================
// CommaExpression comma expression
type CommaExpression struct {
	Exprs []clause.Expression
}

func (comma CommaExpression) Build(builder clause.Builder) {
	for idx, expr := range comma.Exprs {
		if idx > 0 {
			_, _ = builder.WriteString(", ")
		}
		expr.Build(builder)
	}
}

// ======================== New Table ========================

// Table return a new table produced by subquery,
// the return value has to be used as root node
//
// 	Table(u.Select(u.ID, u.Name).Where(u.Age.Gt(18))).Select()
// the above usage is equaivalent to SQL statement:
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
			tablePlaceholder[i] += " AS " + do.db.Statement.Quote(do.alias)
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
