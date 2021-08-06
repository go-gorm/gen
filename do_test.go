package gen

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/field"
)

var db, _ = gorm.Open(tests.DummyDialector{}, nil)

func init() {
	db = db.Debug()

	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		UpdateClauses: []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"},
		DeleteClauses: []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"},
	})
}

// UserRaw user data struct
type UserRaw struct {
	ID         uint `gorm:"primary_key"`
	Name       string
	Age        int
	Score      float64
	Address    string
	Famous     bool
	RegisterAt time.Time
}

func (UserRaw) TableName() string {
	return "users_info"
}

// StudentRaw student data struct
type StudentRaw struct {
	ID         int64 `gorm:"primary_key"`
	Name       string
	Age        int
	Instructor int64 //导师
}

func (StudentRaw) TableName() string {
	return "student"
}

// Teacher teacher data struct
type TeacherRaw struct {
	ID   int64 `gorm:"primary_key"`
	Name string
}

func (TeacherRaw) TableName() string {
	return "teacher"
}

type User struct {
	DO

	ID         field.Uint
	Name       field.String
	Age        field.Int
	Score      field.Float64
	Address    field.String
	Famous     field.Bool
	RegisterAt field.Time
}

var u = func() *User {
	u := User{
		ID:         field.NewUint("", "id"),
		Name:       field.NewString("", "name"),
		Age:        field.NewInt("", "age"),
		Score:      field.NewFloat64("", "score"),
		Address:    field.NewString("", "address"),
		Famous:     field.NewBool("", "famous"),
		RegisterAt: field.NewTime("", "register_at"),
	}
	u.UseDB(db.Session(&gorm.Session{DryRun: true}))
	u.UseModel(UserRaw{})
	return &u
}()

type Student struct {
	DO

	ID         field.Int64
	Name       field.String
	Age        field.Int
	Instructor field.Int64
}

var student = func() *Student {
	s := Student{
		ID:         field.NewInt64("student", "id"),
		Name:       field.NewString("student", "name"),
		Age:        field.NewInt("student", "age"),
		Instructor: field.NewInt64("student", "instructor"),
	}
	s.UseDB(db.Session(&gorm.Session{DryRun: true}))
	s.UseModel(StudentRaw{})
	return &s
}()

type Teacher struct {
	DO

	ID   field.Int64
	Name field.String
}

var teacher = func() *Teacher {
	t := Teacher{
		ID:   field.NewInt64("teacher", "id"),
		Name: field.NewString("teacher", "name"),
	}
	t.UseDB(db.Session(&gorm.Session{DryRun: true}))
	t.UseModel(TeacherRaw{})
	return &t
}()

func checkBuildExpr(t *testing.T, e Dao, opts []stmtOpt, result string, vars []interface{}) {
	stmt := e.(*DO).buildStmt(opts...)

	sql := strings.TrimSpace(stmt.SQL.String())
	if sql != result {
		t.Errorf("Sql expects %v got %v", result, sql)
	}

	if !reflect.DeepEqual(stmt.Vars, vars) {
		t.Errorf("Vars expects %+v got %v", vars, stmt.Vars)
	}
}

func TestDO_methods(t *testing.T) {
	testcases := []struct {
		Expr         Dao
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
			Result: "SELECT `id`, `name`",
		},
		{
			Expr:   u.Distinct(u.Name),
			Result: "SELECT `name`",
		},
		{
			Expr:   teacher.Distinct(teacher.ID, teacher.Name),
			Result: "SELECT `teacher`.`id`, `teacher`.`name`",
		},
		{
			Expr:   teacher.Select(teacher.ID, teacher.Name).Distinct(),
			Result: "SELECT `teacher`.`id`, `teacher`.`name`",
		},
		{
			Expr:   teacher.Select(teacher.Name.As("n")).Distinct(),
			Result: "SELECT `teacher`.`name` AS `n`",
		},
		{
			Expr:   teacher.Select(teacher.ID.As("i"), teacher.Name.As("n")).Distinct(),
			Result: "SELECT `teacher`.`id` AS `i`, `teacher`.`name` AS `n`",
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
			Expr:         u.Order(u.ID),
			ExpectedVars: nil,
			Result:       "ORDER BY `id`",
		},
		{
			Expr:         u.Order(u.ID.Desc()),
			ExpectedVars: nil,
			Result:       "ORDER BY `id` DESC",
		},
		{
			Expr:         u.Order(u.ID.Desc(), u.Age),
			ExpectedVars: nil,
			Result:       "ORDER BY `id` DESC, `age`",
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
			Result:       "WHERE `age` <= ? OR (`name` = ? AND `famous` IS ?)",
		},
		{
			Expr:         u.Where(u.Where(u.Name.Eq("tom"), u.Famous.Is(true))).Or(u.Age.Lte(18)),
			ExpectedVars: []interface{}{"tom", true, 18},
			Result:       "WHERE (`name` = ? AND `famous` IS ?) OR `age` <= ?",
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
			Result:       "WHERE (`age` > ? AND (`famous` IS ? OR `score` >= ?)) OR (`age` <= ? AND `name` = ?)",
		},
		{
			Expr:         u.Select(u.ID, u.Name).Where(u.Age.Gt(18), u.Score.Gte(100)),
			ExpectedVars: []interface{}{18, 100.0},
			Result:       "SELECT `id`, `name` WHERE `age` > ? AND `score` >= ?",
		},
		// ======================== subquery ========================
		{
			Expr:         u.Select().Where(Eq(u.ID, u.Select(u.ID.Max()))),
			ExpectedVars: nil,
			Result:       "SELECT * WHERE `id` = (SELECT MAX(`id`) FROM `users_info`)",
		},
		{
			Expr:         u.Select(u.ID).Where(Gt(u.Score, u.Select(u.Score.Avg()))),
			ExpectedVars: nil,
			Result:       "SELECT `id` WHERE `score` > (SELECT AVG(`score`) FROM `users_info`)",
		},
		{
			Expr:         u.Select(u.ID, u.Name).Where(Lte(u.Score, u.Select(u.Score.Avg()).Where(u.Age.Gte(18)))),
			ExpectedVars: []interface{}{18},
			Result:       "SELECT `id`, `name` WHERE `score` <= (SELECT AVG(`score`) FROM `users_info` WHERE `age` >= ?)",
		},
		{
			Expr:         u.Select(u.ID).Where(In(u.Score, u.Select(u.Score).Where(u.Age.Gte(18)))),
			ExpectedVars: []interface{}{18},
			Result:       "SELECT `id` WHERE `score` IN (SELECT `score` FROM `users_info` WHERE `age` >= ?)",
		},
		{
			Expr:         u.Select(u.ID).Where(In(u.ID, u.Age, u.Select(u.ID, u.Age).Where(u.Score.Eq(100)))),
			ExpectedVars: []interface{}{100.0},
			Result:       "SELECT `id` WHERE (`id`, `age`) IN (SELECT `id`, `age` FROM `users_info` WHERE `score` = ?)",
		},
		{
			Expr:         u.Select(u.Age.Avg().As("avgage")).Group(u.Name).Having(Gt(u.Age.Avg(), u.Select(u.Age.Avg()).Where(u.Name.Like("name%")))),
			Opts:         []stmtOpt{withFROM},
			ExpectedVars: []interface{}{"name%"},
			Result:       "SELECT AVG(`age`) AS `avgage` FROM `users_info` GROUP BY `name` HAVING AVG(`age`) > (SELECT AVG(`age`) FROM `users_info` WHERE `name` LIKE ?)",
		},
		// ======================== from subquery ========================
		{
			Expr:         Table(u.Select(u.ID, u.Name).Where(u.Age.Gt(18))).Select(),
			Opts:         []stmtOpt{withFROM},
			ExpectedVars: []interface{}{18},
			Result:       "SELECT * FROM (SELECT `id`, `name` FROM `users_info` WHERE `age` > ?)",
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
			Result:       "SELECT `student`.`name`, `teacher`.`name` FROM `student` LEFT JOIN `teacher` ON `student`.`instructor` = `teacher`.`id` WHERE `teacher`.`id` > ?",
			ExpectedVars: []interface{}{int64(0)},
		},
		{
			Expr:         student.RightJoin(teacher, student.Instructor.EqCol(teacher.ID), student.ID.Eq(666)).Select(),
			Result:       "SELECT * FROM `student` RIGHT JOIN `teacher` ON `student`.`instructor` = `teacher`.`id` AND `student`.`id` = ?",
			ExpectedVars: []interface{}{int64(666)},
		},
		{
			Expr:         student.Join(teacher, student.Instructor.EqCol(teacher.ID)).LeftJoin(teacher, student.ID.Gt(100)).Select(student.ID, student.Name, teacher.Name.As("teacher_name")),
			Result:       "SELECT `student`.`id`, `student`.`name`, `teacher`.`name` AS `teacher_name` FROM `student` INNER JOIN `teacher` ON `student`.`instructor` = `teacher`.`id` LEFT JOIN `teacher` ON `student`.`id` > ?",
			ExpectedVars: []interface{}{int64(100)},
		},
	}

	// _ = u.Update(u.Age, u.Age.Add(1))
	// _ = u.Update(u.Age, gorm.Expr("age+1"))

	// _ = u.Find(u.ID.In(1, 2, 3))

	for _, testcase := range testcases {
		checkBuildExpr(t, testcase.Expr, testcase.Opts, testcase.Result, testcase.ExpectedVars)
	}
}
