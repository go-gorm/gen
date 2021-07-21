package gen

import (
	"database/sql"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"gorm.io/gen/field"
)

var (
	createClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	queryClauses  = []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR"}
	updateClauses = []string{"UPDATE", "SET", "WHERE"}
	deleteClauses = []string{"DELETE", "FROM", "WHERE"}
)

func NewStage(db *gorm.DB) *Stage { return &Stage{db: db} }

// Stage implement basic query methods
// the structure embedded with a *gorm.DB, and has a element item "alias" will be used when used as a sub query
type Stage struct {
	db    *gorm.DB
	alias string // for subquery
}

// UseDB specify a db connection(*gorm.DB)
func (s *Stage) UseDB(db *gorm.DB) {
	s.db = db.Session(new(gorm.Session))
}

// UseModel specify a data model structure as a source for table name
func (s *Stage) UseModel(model interface{}) {
	s.db = s.db.Model(model).Session(new(gorm.Session))
	_ = s.db.Statement.Parse(model)
}

// UseTable specify table name
func (s *Stage) UseTable(tableName string) {
	s.db = s.db.Table(tableName).Session(new(gorm.Session))
}

// Table return table name
func (s *Stage) Table() string {
	return s.db.Statement.Table
}

// UnderlyingDB return the underlying database connection
func (s *Stage) UnderlyingDB() *gorm.DB {
	Emit(diyM)
	return s.db
}

// Build implement the interface of claues.Expression
// only call WHERE clause's Build
func (s *Stage) Build(builder clause.Builder) {
	for _, e := range s.buildWhere() {
		e.Build(builder)
	}
}

func (s *Stage) buildWhere() []clause.Expression {
	return s.db.Statement.BuildCondition(s.db)
}

type stmtOpt func(*gorm.Statement) *gorm.Statement

var (
	// withFROM 增加FROM子句
	withFROM stmtOpt = func(stmt *gorm.Statement) *gorm.Statement {
		if stmt.Table == "" {
			_ = stmt.Parse(stmt.Model)
		}
		stmt.AddClause(clause.From{})
		return stmt
	}

	// // withSELECT 增加SELECT子句
	// withSELECT stmtOpt = func(stmt *gorm.Statement) *gorm.Statement {
	// 	if _, ok := stmt.Clauses["SELECT"]; !ok {
	// 		stmt.AddClause(clause.Select{})
	// 	}
	// 	return stmt
	// }
)

// buildStmt call statement.Build to combine all clauses in one statement
func (s *Stage) buildStmt(opts ...stmtOpt) *gorm.Statement {
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

// func (s *Stage) subQueryExpr() clause.Expr {
// 	stmt := s.buildStmt(withFROM, withSELECT)
// 	return clause.Expr{SQL: "(" + stmt.SQL.String() + ")", Vars: stmt.Vars}
// }

// As 指定的值不可继承，因此需要在结尾使用
func (s *Stage) As(alias string) Dao {
	return &Stage{db: s.db, alias: alias}
}

// ======================== 逻辑操作 ========================
func (s *Stage) Not(conds ...Condition) Dao {
	return NewStage(s.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Not(condToExpression(conds...)...)}}))
}

func (s *Stage) Or(conds ...Condition) Dao {
	return NewStage(s.db.Clauses(clause.Where{Exprs: []clause.Expression{clause.Or(clause.And(condToExpression(conds...)...))}}))
}

// ======================== chainable api ========================
func (s *Stage) Select(columns ...field.Expr) Dao {
	Emit(selectM)
	if len(columns) == 0 {
		return NewStage(s.db.Clauses(clause.Select{}))
	}
	return NewStage(s.db.Clauses(clause.Select{Expression: CommaExpression{Exprs: toExpression(columns...)}}))
}

func (s *Stage) Where(conds ...Condition) Dao {
	Emit(whereM)
	var exprs = make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		switch cond := cond.(type) {
		case *Stage:
			exprs = append(exprs, cond.buildWhere()...)
		default:
			exprs = append(exprs, cond)
		}
	}
	return NewStage(s.db.Clauses(clause.Where{Exprs: exprs}))
}

func (s *Stage) Order(columns ...field.Expr) Dao {
	Emit(orderM)
	return NewStage(s.db.Clauses(clause.OrderBy{Expression: CommaExpression{Exprs: toExpression(columns...)}}))
}

func (s *Stage) Distinct(columns ...field.Expr) Dao {
	Emit(distinctM)
	return NewStage(s.db.Distinct(toInterfaceSlice(toColNames(s.db.Statement, columns...))))
}

func (s *Stage) Omit(columns ...field.Expr) Dao {
	Emit(omitM)
	return NewStage(s.db.Omit(toColNames(s.db.Statement, columns...)...))
}

func (s *Stage) Group(column field.Expr) Dao {
	Emit(groupM)
	return NewStage(s.db.Group(column.Column().Name))
}

func (s *Stage) Having(conds ...Condition) Dao {
	Emit(havingM)
	return NewStage(s.db.Clauses(clause.GroupBy{Having: condToExpression(conds...)}))
}

func (s *Stage) Limit(limit int) Dao {
	Emit(limitM)
	return NewStage(s.db.Limit(limit))
}

func (s *Stage) Offset(offset int) Dao {
	Emit(offsetM)
	return NewStage(s.db.Offset(offset))
}

func (s *Stage) Scopes(funcs ...func(Dao) Dao) Dao {
	Emit(scopesM)
	var result Dao = s
	for _, f := range funcs {
		result = f(result)
	}
	return result
}

func (s *Stage) Unscoped() Dao {
	Emit(unscopedM)
	return NewStage(s.db.Unscoped())
}

// ======================== finisher api ========================
func (s *Stage) Create(value interface{}) error {
	Emit(createM)
	return s.db.Create(value).Error
}

func (s *Stage) CreateInBatches(value interface{}, batchSize int) error {
	Emit(createInBatchesM)
	return s.db.CreateInBatches(value, batchSize).Error
}

func (s *Stage) Save(value interface{}) error {
	Emit(saveM)
	return s.db.Save(value).Error
}

func (s *Stage) First(dest interface{}, conds ...field.Expr) error {
	Emit(firstM)
	return s.db.Clauses(toExpression(conds...)...).First(dest).Error
}

func (s *Stage) Take(dest interface{}, conds ...field.Expr) error {
	Emit(takeM)
	return s.db.Clauses(toExpression(conds...)...).Take(dest).Error
}

func (s *Stage) Last(dest interface{}, conds ...field.Expr) error {
	Emit(lastM)
	return s.db.Clauses(toExpression(conds...)...).Last(dest).Error
}

func (s *Stage) Find(dest interface{}, conds ...field.Expr) error {
	Emit(findM)
	return s.db.Clauses(toExpression(conds...)...).Find(dest).Error
}

func (s *Stage) FindInBatches(dest interface{}, batchSize int, fc func(tx Dao, batch int) error) error {
	Emit(findInBatchesM)
	return s.db.FindInBatches(dest, batchSize, func(tx *gorm.DB, batch int) error { return fc(NewStage(tx), batch) }).Error
}

func (s *Stage) FirstOrInit(dest interface{}, conds ...field.Expr) error {
	Emit(firstOrInitM)
	return s.db.Clauses(toExpression(conds...)...).FirstOrInit(dest).Error
}

func (s *Stage) FirstOrCreate(dest interface{}, conds ...field.Expr) error {
	Emit(firstOrCreateM)
	return s.db.Clauses(toExpression(conds...)...).FirstOrCreate(dest).Error
}

func (s *Stage) Update(column field.Expr, value interface{}) error {
	Emit(updateM)
	switch expr := column.RawExpr().(type) {
	case clause.Expression:
		return s.db.Update(column.Column().Name, expr).Error
	}

	switch value := value.(type) {
	case field.Expr:
		return s.db.Update(column.Column().Name, value.RawExpr()).Error
	case *Stage:
		return s.db.Update(column.Column().Name, value.db).Error
	default:
		return s.db.Update(column.Column().Name, value).Error
	}
}

func (s *Stage) Updates(values interface{}) error {
	Emit(updatesM)
	if vs, ok := values.(map[string]interface{}); ok {
		return s.db.Updates(toQuotedSubQueryMap(s.db.Statement, vs)).Error
	}
	return s.db.Updates(values).Error
}

func (s *Stage) UpdateColumn(column field.Expr, value interface{}) error {
	Emit(updateColumnM)
	switch expr := column.RawExpr().(type) {
	case clause.Expression:
		return s.db.UpdateColumn(column.Column().Name, expr).Error
	}

	switch value := value.(type) {
	case field.Expr:
		return s.db.UpdateColumn(column.Column().Name, value.RawExpr()).Error
	case *Stage:
		return s.db.UpdateColumn(column.Column().Name, value.db).Error
	default:
		return s.db.UpdateColumn(column.Column().Name, value).Error
	}
}

func (s *Stage) UpdateColumns(values interface{}) error {
	Emit(updateColumnsM)
	if vs, ok := values.(map[string]interface{}); ok {
		return s.db.UpdateColumns(toQuotedSubQueryMap(s.db.Statement, vs)).Error
	}
	return s.db.UpdateColumns(values).Error
}

func (s *Stage) Delete(value interface{}, conds ...field.Expr) error {
	Emit(deleteM)
	return s.db.Clauses(toExpression(conds...)...).Delete(value).Error
}

func (s *Stage) Count(count *int64) error {
	Emit(countM)
	return s.db.Count(count).Error
}

func (s *Stage) Row() *sql.Row {
	Emit(rowM)
	return s.db.Row()
}

func (s *Stage) Rows() (*sql.Rows, error) {
	Emit(rowsM)
	return s.db.Rows()
}

func (s *Stage) Scan(dest interface{}) error {
	Emit(scanM)
	return s.db.Scan(dest).Error
}

func (s *Stage) Pluck(column field.Expr, dest interface{}) error {
	Emit(pluckM)
	return s.db.Pluck(column.Column().Name, dest).Error
}

func (s *Stage) ScanRows(rows *sql.Rows, dest interface{}) error {
	Emit(scanRowsM)
	return s.db.ScanRows(rows, dest)
}

func (s *Stage) Transaction(fc func(Dao) error, opts ...*sql.TxOptions) error {
	Emit(transactionM)
	return s.db.Transaction(func(tx *gorm.DB) error { return fc(NewStage(tx)) }, opts...)
}

func (s *Stage) Begin(opts ...*sql.TxOptions) Dao {
	Emit(beginM)
	return NewStage(s.db.Begin(opts...))
}

func (s *Stage) Commit() Dao {
	Emit(commitM)
	return NewStage(s.db.Commit())
}

func (s *Stage) RollBack() Dao {
	Emit(rollbackM)
	return NewStage(s.db.Rollback())
}

func (s *Stage) SavePoint(name string) Dao {
	Emit(savePointM)
	return NewStage(s.db.SavePoint(name))
}

func (s *Stage) RollBackTo(name string) Dao {
	Emit(rollbackToM)
	return NewStage(s.db.RollbackTo(name))
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

func toQuotedSubQueryMap(stmt *gorm.Statement, values map[string]interface{}) map[string]interface{} {
	escapedValues := make(map[string]interface{}, len(values))
	for k, v := range values {
		if query, ok := v.(*Stage); ok {
			v = query.db
		}
		escapedValues[stmt.Quote(k)] = v
	}
	return escapedValues
}

// ======================== 临时数据结构 ========================
// 逗号分割的表达式
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
		return NewStage(nil)
	}

	tablePlaceholder := make([]string, len(subQueries))
	tableExprs := make([]interface{}, len(subQueries))
	for i, query := range subQueries {
		tablePlaceholder[i] = "(?)"

		stage := query.(*Stage)
		tableExprs[i] = stage.db
		if stage.alias != "" {
			tablePlaceholder[i] += " AS " + stage.db.Statement.Quote(stage.alias)
		}
	}

	db := subQueries[0].(*Stage).db
	return NewStage(db.Session(&gorm.Session{NewDB: true}).Table(strings.Join(tablePlaceholder, ", "), tableExprs...))
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
