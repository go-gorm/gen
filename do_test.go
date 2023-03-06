package gen

import (
	"reflect"
	"strings"
	"testing"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"

	"gorm.io/gen/field"
)

var (
	createClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	queryClauses  = []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT", "FOR"}
	updateClauses = []string{"UPDATE", "SET", "WHERE"}
	deleteClauses = []string{"DELETE", "FROM", "WHERE"}
)

type stmtOpt func(*gorm.Statement) *gorm.Statement

var (
	// withFROM add FROM clause
	withFROM stmtOpt = func(stmt *gorm.Statement) *gorm.Statement {
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

func checkBuildExpr(t *testing.T, e SubQuery, opts []stmtOpt, result string, vars []interface{}) {
	stmt := build(e.underlyingDB().Statement, opts...)

	sql := strings.TrimSpace(stmt.SQL.String())
	if sql != result {
		t.Errorf("SQL expects %v got %v", result, sql)
	}

	if !reflect.DeepEqual(stmt.Vars, vars) {
		t.Errorf("Vars expects %+v got %v", vars, stmt.Vars)
	}
}

func build(stmt *gorm.Statement, opts ...stmtOpt) *gorm.Statement {
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

func TestDO_methods(t *testing.T) {
	testcases := []struct {
		Expr         SubQuery
		Opts         []stmtOpt
		ExpectedVars []interface{}
		Result       string
	}{
		{
			Expr:   u.Select(),
			Result: "SELECT *",
		},
		{
			Expr:   u.Select(u.ID, u.Name),
			Result: "SELECT `id`,`name`",
		},
		{
			Expr:   u.Distinct(u.Name),
			Result: "SELECT DISTINCT `name`",
		},
		{
			Expr:   teacher.Distinct(teacher.ID, teacher.Name),
			Result: "SELECT DISTINCT `teacher`.`id`,`teacher`.`name`",
		},
		{
			Expr:   teacher.Select(teacher.ID, teacher.Name).Distinct(),
			Result: "SELECT DISTINCT `teacher`.`id`,`teacher`.`name`",
		},
		{
			Expr:   teacher.Distinct().Select(teacher.ID, teacher.Name),
			Result: "SELECT DISTINCT `teacher`.`id`,`teacher`.`name`",
		},
		{
			Expr:   teacher.Select(teacher.Name.As("n")).Distinct(),
			Result: "SELECT DISTINCT `teacher`.`name` AS `n`",
		},
		{
			Expr:   teacher.Select(field.ALL),
			Result: "SELECT *",
		},
		{
			Expr:   teacher.Select(field.ALL.Count()),
			Result: "SELECT COUNT(*)",
		},
		{
			Expr:   teacher.Select(field.ALL.Distinct().Count()),
			Result: "SELECT COUNT(DISTINCT *)",
		},
		{
			Expr:   teacher.Select(teacher.ALL),
			Result: "SELECT `teacher`.*",
		},
		{
			Expr:   teacher.Select(teacher.ALL.Count()),
			Result: "SELECT COUNT(`teacher`.*)",
		},
		{
			Expr:   teacher.Select(teacher.ALL.Distinct().Count()),
			Result: "SELECT COUNT(DISTINCT `teacher`.*)",
		},
		{
			Expr:   teacher.Select(teacher.ID.As("i"), teacher.Name.As("n")).Distinct(),
			Result: "SELECT DISTINCT `teacher`.`id` AS `i`,`teacher`.`name` AS `n`",
		},
		{
			Expr:         u.Where(u.ID.Eq(10)),
			ExpectedVars: []interface{}{uint(10)},
			Result:       "WHERE `id` = ?",
		},
		{
			Expr:         u.Where(u.Name.Eq("tom"), u.Age.Gt(18)),
			ExpectedVars: []interface{}{"tom", 18},
			Result:       "WHERE `name` = ? AND `age` > ?",
		},
		{
			Expr:   u.Order(u.ID),
			Result: "ORDER BY `id`",
		},
		{
			Expr:   u.Order(u.ID.Desc()),
			Result: "ORDER BY `id` DESC",
		},
		{
			Expr:   u.Order(u.ID.Desc(), u.Age),
			Result: "ORDER BY `id` DESC,`age`",
		},
		{
			Expr:   u.Order(u.ID.Desc()).Order(u.Age),
			Result: "ORDER BY `id` DESC,`age`",
		},
		{
			Expr:   u.Clauses(hints.New("hint")).Select(),
			Result: "SELECT /*+ hint */ *",
		},
		{
			Expr:   u.Clauses(hints.Comment("select", "hint")).Select(),
			Result: "SELECT /* hint */ *",
		},
		{
			Expr:   u.Clauses(hints.CommentBefore("select", "hint")).Select(),
			Result: "/* hint */ SELECT *",
		},
		{
			Expr:   u.Clauses(hints.CommentAfter("select", "hint")).Select(),
			Result: "SELECT * /* hint */",
		},
		{
			Expr:         u.Clauses(hints.CommentAfter("where", "hint")).Select().Where(u.ID.Gt(0)),
			ExpectedVars: []interface{}{uint(0)},
			Result:       "SELECT * WHERE `id` > ? /* hint */",
		},
		{
			Expr:   u.Clauses(hints.UseIndex("user_name")).Select(),
			Opts:   []stmtOpt{withFROM},
			Result: "SELECT * FROM `users_info` USE INDEX (`user_name`)",
		},
		{
			Expr:   u.Clauses(hints.ForceIndex("user_name", "user_id").ForJoin()).Select(),
			Opts:   []stmtOpt{withFROM},
			Result: "SELECT * FROM `users_info` FORCE INDEX FOR JOIN (`user_name`,`user_id`)",
		},
		{
			Expr: u.Clauses(
				hints.ForceIndex("user_name", "user_id").ForJoin(),
				hints.IgnoreIndex("user_name").ForGroupBy(),
			).Select(),
			Opts:   []stmtOpt{withFROM},
			Result: "SELECT * FROM `users_info` FORCE INDEX FOR JOIN (`user_name`,`user_id`) IGNORE INDEX FOR GROUP BY (`user_name`)",
		},
		// ======================== where conditions ========================
		{
			Expr:         u.Where(u.Where(u.ID.Neq(0)), u.Where(u.Age.Gt(18))),
			ExpectedVars: []interface{}{uint(0), 18},
			Result:       "WHERE `id` <> ? AND `age` > ?",
		},
		{
			Expr:         u.Where(u.Age.Lte(18)).Or(u.Where(u.Name.Eq("tom"))),
			ExpectedVars: []interface{}{18, "tom"},
			Result:       "WHERE `age` <= ? OR `name` = ?",
		},
		{
			Expr:         u.Where(u.Age.Lte(18)).Or(u.Name.Eq("tom"), u.Famous.Is(true)),
			ExpectedVars: []interface{}{18, "tom", true},
			Result:       "WHERE `age` <= ? OR (`name` = ? AND `famous` = ?)",
		},
		{
			Expr:         u.Where(u.Columns(u.ID, u.Age).In(field.Values([][]int{{1, 18}, {2, 19}}))),
			ExpectedVars: []interface{}{1, 18, 2, 19},
			Result:       "WHERE (`id`, `age`) IN ((?,?),(?,?))",
		},
		{
			Expr:         u.Where(u.Columns(u.ID, u.Age).NotIn(field.Values([][]int{{1, 18}, {2, 19}}))),
			ExpectedVars: []interface{}{1, 18, 2, 19},
			Result:       "WHERE NOT (`id`, `age`) IN ((?,?),(?,?))",
		},
		{
			Expr:         u.Where(u.Columns(u.ID, u.Name).In(field.Values([][]interface{}{{1, "modi"}, {2, "tom"}}))),
			ExpectedVars: []interface{}{1, "modi", 2, "tom"},
			Result:       "WHERE (`id`, `name`) IN ((?,?),(?,?))",
		},
		{
			Expr:         u.Where(u.Where(u.Name.Eq("tom"), u.Famous.Is(true))).Or(u.Age.Lte(18)),
			ExpectedVars: []interface{}{"tom", true, 18},
			Result:       "WHERE (`name` = ? AND `famous` = ?) OR `age` <= ?",
		},
		{
			Expr:         u.Where(u.Name.Eq("tom")).Where(u.Where(u.Famous.Is(true)).Or(u.Age.Lte(18))),
			ExpectedVars: []interface{}{"tom", true, 18},
			Result:       "WHERE `name` = ? AND (`famous` = ? OR `age` <= ?)",
		},
		{
			Expr:         u.Where(u.Name.Eq("tom"), field.Or(u.Famous.Is(true), u.Age.Lte(18))),
			ExpectedVars: []interface{}{"tom", true, 18},
			Result:       "WHERE `name` = ? AND (`famous` = ? OR `age` <= ?)",
		},
		{
			Expr:         u.Where(Cond(datatypes.JSONQuery("attributes").HasKey("role", "name"))...),
			ExpectedVars: []interface{}{"$.role.name"},
			Result:       "WHERE JSON_EXTRACT(`attributes`,?) IS NOT NULL",
		},
		{
			Expr: u.Where(
				u.Where(u.ID.Neq(0)).Where(u.Score.Gt(89.9)),
				u.Where(u.Age.Gt(18)).Where(u.Address.Eq("New York")),
			),
			ExpectedVars: []interface{}{uint(0), 89.9, 18, "New York"},
			Result:       "WHERE (`id` <> ? AND `score` > ?) AND (`age` > ? AND `address` = ?)",
		},
		{
			Expr: u.Where(
				u.Where(u.Age.Gt(18)).Where(u.Where(u.Famous.Is(true)).Or(u.Score.Gte(100.0))),
			).Or(
				u.Where(u.Age.Lte(18)).Where(u.Name.Eq("tom")),
			),
			ExpectedVars: []interface{}{18, true, 100.0, 18, "tom"},
			Result:       "WHERE (`age` > ? AND (`famous` = ? OR `score` >= ?)) OR (`age` <= ? AND `name` = ?)",
		},
		{
			Expr:         u.Select(u.ID, u.Name).Where(u.Age.Gt(18), u.Score.Gte(100)),
			ExpectedVars: []interface{}{18, 100.0},
			Result:       "SELECT `id`,`name` WHERE `age` > ? AND `score` >= ?",
		},
		{
			Expr:         u.Select().Where(Cond(datatypes.JSONQuery("attributes").HasKey("role"))...),
			Opts:         []stmtOpt{withFROM},
			ExpectedVars: []interface{}{"$.role"},
			Result:       "SELECT * FROM `users_info` WHERE JSON_EXTRACT(`attributes`,?) IS NOT NULL",
		},
		// ======================== subquery ========================
		{
			Expr:   u.Select().Where(u.Columns(u.ID).Eq(u.Select(u.ID.Max()))),
			Result: "SELECT * WHERE `id` = (SELECT MAX(`id`) FROM `users_info`)",
		},
		{
			Expr: u.Select().Where(
				u.Columns(u.ID).Eq(
					u.Select(u.ID).Where(
						u.Columns(u.Name).Eq(
							student.Select(student.Name).Where(student.ID.Eq(1)),
						),
					),
				)),
			ExpectedVars: []interface{}{int64(1)},
			Result:       "SELECT * WHERE `id` = (SELECT `id` FROM `users_info` WHERE `name` = (SELECT `student`.`name` FROM `student` WHERE `student`.`id` = ?))",
		},
		{
			Expr:   u.Select().Where(u.Columns(u.ID).Neq(u.Select(u.ID.Max()))),
			Result: "SELECT * WHERE `id` <> (SELECT MAX(`id`) FROM `users_info`)",
		},
		{
			Expr:         u.Select(u.ID).Where(u.Columns(u.Score.Mul(2)).Lte(u.Select(u.Score.Avg()))),
			ExpectedVars: []interface{}{2.0},
			Result:       "SELECT `id` WHERE `score`*? <= (SELECT AVG(`score`) FROM `users_info`)",
		},
		{
			Expr:   u.Select(u.ID).Where(u.Columns(u.Score).Gt(u.Select(u.Score.Avg()))),
			Result: "SELECT `id` WHERE `score` > (SELECT AVG(`score`) FROM `users_info`)",
		},
		{
			Expr:         u.Select(u.ID, u.Name).Where(u.Columns(u.Score).Lte(u.Select(u.Score.Avg()).Where(u.Age.Gte(18)))),
			ExpectedVars: []interface{}{18},
			Result:       "SELECT `id`,`name` WHERE `score` <= (SELECT AVG(`score`) FROM `users_info` WHERE `age` >= ?)",
		},
		{
			Expr:         u.Select(u.ID).Where(u.Columns(u.Score).In(u.Select(u.Score).Where(u.Age.Gte(18)))),
			ExpectedVars: []interface{}{18},
			Result:       "SELECT `id` WHERE `score` IN (SELECT `score` FROM `users_info` WHERE `age` >= ?)",
		},
		{
			Expr:         u.Select(u.ID).Where(u.Columns(u.ID, u.Age).In(u.Select(u.ID, u.Age).Where(u.Score.Eq(100)))),
			ExpectedVars: []interface{}{100.0},
			Result:       "SELECT `id` WHERE (`id`,`age`) IN (SELECT `id`,`age` FROM `users_info` WHERE `score` = ?)",
		},
		{
			Expr:         u.Select(u.Age.Avg().As("avgage")).Group(u.Name).Having(u.Columns(u.Age.Avg()).Gt(u.Select(u.Age.Avg()).Where(u.Name.Like("name%")))),
			Opts:         []stmtOpt{withFROM},
			ExpectedVars: []interface{}{"name%"},
			Result:       "SELECT AVG(`age`) AS `avgage` FROM `users_info` GROUP BY `name` HAVING AVG(`age`) > (SELECT AVG(`age`) FROM `users_info` WHERE `name` LIKE ?)",
		},
		{
			Expr:   student.Select().LeftJoin(teacher, teacher.ID.EqCol(student.Instructor)).Group(student.ID),
			Result: "SELECT * FROM `student` LEFT JOIN `teacher` ON `teacher`.`id` = `student`.`instructor` GROUP BY `student`.`id`",
		},
		// ======================== from subquery ========================
		{
			Expr:         Table(u.Select(u.ID, u.Name).Where(u.Age.Gt(18))).Select(),
			Opts:         []stmtOpt{withFROM},
			ExpectedVars: []interface{}{18},
			Result:       "SELECT * FROM (SELECT `id`,`name` FROM `users_info` WHERE `age` > ?)",
		},
		{
			Expr:         Table(u.Select(u.ID).Where(u.Age.Gt(18)), u.Select(u.ID).Where(u.Score.Gte(100))).Select(),
			Opts:         []stmtOpt{withFROM},
			ExpectedVars: []interface{}{18, 100.0},
			Result:       "SELECT * FROM (SELECT `id` FROM `users_info` WHERE `age` > ?), (SELECT `id` FROM `users_info` WHERE `score` >= ?)",
		},
		{
			Expr:         Table(u.Select().Where(u.Age.Gt(18)), u.Where(u.Score.Gte(100))).Select(),
			Opts:         []stmtOpt{withFROM},
			ExpectedVars: []interface{}{18, 100.0},
			Result:       "SELECT * FROM (SELECT * FROM `users_info` WHERE `age` > ?), (SELECT * FROM `users_info` WHERE `score` >= ?)",
		},
		{
			Expr:         Table(u.Select().Where(u.Age.Gt(18)).As("a"), u.Where(u.Score.Gte(100)).As("b")).Select(),
			Opts:         []stmtOpt{withFROM},
			ExpectedVars: []interface{}{18, 100.0},
			Result:       "SELECT * FROM (SELECT * FROM `users_info` WHERE `age` > ?) AS `a`, (SELECT * FROM `users_info` WHERE `score` >= ?) AS `b`",
		},

		// ======================== join subquery ========================
		{
			Expr:   student.Join(teacher, student.Instructor.EqCol(teacher.ID)).Select(),
			Result: "SELECT * FROM `student` INNER JOIN `teacher` ON `student`.`instructor` = `teacher`.`id`",
		},
		{
			Expr:         student.LeftJoin(teacher, student.Instructor.EqCol(teacher.ID)).Where(teacher.ID.Gt(0)).Select(student.Name, teacher.Name),
			Result:       "SELECT `student`.`name`,`teacher`.`name` FROM `student` LEFT JOIN `teacher` ON `student`.`instructor` = `teacher`.`id` WHERE `teacher`.`id` > ?",
			ExpectedVars: []interface{}{int64(0)},
		},
		{
			Expr:         student.RightJoin(teacher, student.Instructor.EqCol(teacher.ID), student.ID.Eq(666)).Select(),
			Result:       "SELECT * FROM `student` RIGHT JOIN `teacher` ON `student`.`instructor` = `teacher`.`id` AND `student`.`id` = ?",
			ExpectedVars: []interface{}{int64(666)},
		},
		{
			Expr:         student.Join(teacher, student.Instructor.EqCol(teacher.ID)).LeftJoin(teacher, student.ID.Gt(100)).Select(student.ID, student.Name, teacher.Name.As("teacher_name")),
			Result:       "SELECT `student`.`id`,`student`.`name`,`teacher`.`name` AS `teacher_name` FROM `student` INNER JOIN `teacher` ON `student`.`instructor` = `teacher`.`id` LEFT JOIN `teacher` ON `student`.`id` > ?",
			ExpectedVars: []interface{}{int64(100)},
		},
	}

	// _ = u.Update(u.Age, u.Age.Add(1))
	// _ = u.Update(u.Age, gorm.Expr("age+1"))
	// _ = u.UpdateSimple(u.Age.Add(1))

	// _ = u.Find(u.ID.In(1, 2, 3))

	for _, testcase := range testcases {
		checkBuildExpr(t, testcase.Expr, testcase.Opts, testcase.Result, testcase.ExpectedVars)
	}
}
