package field_test

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/field"
)

var db, _ = gorm.Open(tests.DummyDialector{}, nil)

type User struct {
	gorm.Model
	Name string
	Age  uint
	// Birthday  *time.Time
	// Account   Account
	// Pets      []*Pet
	// Toys      []Toy `gorm:"polymorphic:Owner"`
	// CompanyID *int
	// Company   Company
	// ManagerID *uint
	// Manager   *User
	// Team      []User     `gorm:"foreignkey:ManagerID"`
	// Languages []Language `gorm:"many2many:UserSpeak;"`
	// Friends   []*User    `gorm:"many2many:user_friends;"`
	// Active    bool
}

var _ field.ScanValuer = new(password)

type password string

func (p *password) Scan(src interface{}) error {
	*p = password(fmt.Sprintf("this is password {%q}", src))
	return nil
}
func (p *password) Value() (driver.Value, error) {
	return strings.TrimPrefix(strings.TrimSuffix(string(*p), "}"), "this is password {"), nil
}

func checkBuildExpr(t *testing.T, e field.Expr, result string, vars []interface{}) {
	user, _ := schema.Parse(&User{}, &sync.Map{}, db.NamingStrategy)
	stmt := &gorm.Statement{DB: db, Table: user.Table, Schema: user, Clauses: map[string]clause.Clause{}}

	e.Build(stmt)

	sql := strings.TrimSpace(stmt.SQL.String())
	if sql != result {
		t.Errorf("SQL expects %v got %v", result, sql)
	}

	if !reflect.DeepEqual(stmt.Vars, vars) {
		t.Errorf("Vars expects %+v got %v", stmt.Vars, vars)
	}
}

func TestExpr_Build(t *testing.T) {
	timeData, _ := time.Parse("2006-01-02 15:04:05", "2021-06-29 15:11:49")
	p := password("i am password")

	testcases := []struct {
		Expr         field.Expr
		ExpectedVars []interface{}
		Result       string
	}{
		// ======================== generic ========================
		{
			Expr:         field.NewField("user", "password").Eq(&p),
			ExpectedVars: []interface{}{&p},
			Result:       "`user`.`password` = ?",
		},
		{
			Expr:   field.NewField("", "id").EqCol(field.NewField("", "new_id")),
			Result: "`id` = `new_id`",
		},
		{
			Expr:   field.NewField("", "id").NeqCol(field.NewField("", "new_id")),
			Result: "`id` <> `new_id`",
		},
		{
			Expr:   field.NewField("", "id").LtCol(field.NewField("", "new_id")),
			Result: "`id` < `new_id`",
		},
		{
			Expr:   field.NewField("", "id").LteCol(field.NewField("", "new_id")),
			Result: "`id` <= `new_id`",
		},
		{
			Expr:   field.NewField("", "id").GtCol(field.NewField("", "new_id")),
			Result: "`id` > `new_id`",
		},
		{
			Expr:   field.NewField("", "id").GteCol(field.NewField("", "new_id")),
			Result: "`id` >= `new_id`",
		},
		{
			Expr:   field.NewField("", "id").EqCol(field.NewField("", "new_id").Avg()),
			Result: "`id` = AVG(`new_id`)",
		},
		{
			Expr:   field.NewField("", "id").EqCol(field.NewField("", "new_id").WithTable("tableB")),
			Result: "`id` = `tableB`.`new_id`",
		},
		{
			Expr:   field.NewField("", "id").EqCol(field.NewField("", "new_id").WithTable("tableB")),
			Result: "`id` = `tableB`.`new_id`",
		},
		{
			Expr:   field.NewField("", "id").NeqCol(field.NewField("", "new_id").WithTable("tableB")),
			Result: "`id` <> `tableB`.`new_id`",
		},
		{
			Expr:   field.NewField("", "id").IsNull(),
			Result: "`id` IS NULL",
		},
		{
			Expr:   field.NewField("", "id").IsNotNull(),
			Result: "`id` IS NOT NULL",
		},
		// ======================== integer ========================
		{
			Expr:   field.NewUint("", "id"),
			Result: "`id`",
		},
		{
			Expr:         field.NewUint("user", "id").Sum().Gt(100),
			ExpectedVars: []interface{}{float64(100)},
			Result:       "SUM(`user`.`id`) > ?",
		},
		{
			Expr:   field.NewUint("", "i`d"),
			Result: "`i`d`", // should be `i``d`
		},
		{
			Expr:   field.NewUint("", "id").Avg(),
			Result: "AVG(`id`)",
		},
		{
			Expr:   field.NewUint("", "id").Desc(),
			Result: "`id` DESC",
		},
		{
			Expr:   field.NewUint("", "id").As("number"),
			Result: "`id` AS `number`",
		},
		{
			Expr:   field.NewUint("", "id").Avg().As("number"),
			Result: "AVG(`id`) AS `number`",
		},
		{
			Expr:         field.NewUint("", "id").Eq(10),
			ExpectedVars: []interface{}{uint(10)},
			Result:       "`id` = ?",
		},
		{
			Expr:         field.NewUint("", "id").Neq(10),
			ExpectedVars: []interface{}{uint(10)},
			Result:       "`id` <> ?",
		},
		{
			Expr:         field.NewUint("", "id").Gt(1),
			ExpectedVars: []interface{}{uint(1)},
			Result:       "`id` > ?",
		},
		{
			Expr:         field.NewUint("", "id").Gte(1),
			ExpectedVars: []interface{}{uint(1)},
			Result:       "`id` >= ?",
		},
		{
			Expr:         field.NewUint("", "id").Lt(1),
			ExpectedVars: []interface{}{uint(1)},
			Result:       "`id` < ?",
		},
		{
			Expr:         field.NewUint("", "id").Lte(1),
			ExpectedVars: []interface{}{uint(1)},
			Result:       "`id` <= ?",
		},
		{
			Expr:         field.NewUint("", "id").Mod(7),
			ExpectedVars: []interface{}{uint(7)},
			Result:       "`id`%?",
		},
		{
			Expr:         field.And(field.NewUint("", "id").Gt(1), field.NewUint("", "id").Lt(10)),
			ExpectedVars: []interface{}{uint(1), uint(10)},
			Result:       "(`id` > ? AND `id` < ?)",
		},
		{
			Expr:         field.Or(field.NewUint("", "id").Lt(4), field.NewUint("", "id").Gt(6)),
			ExpectedVars: []interface{}{uint(4), uint(6)},
			Result:       "(`id` < ? OR `id` > ?)",
		},
		{
			Expr:         field.NewUint("", "id").In(1, 2, 3),
			ExpectedVars: []interface{}{uint(1), uint(2), uint(3)},
			Result:       "`id` IN (?,?,?)",
		},
		{
			Expr:         field.NewUint("", "id").NotIn(1, 2, 3),
			ExpectedVars: []interface{}{uint(1), uint(2), uint(3)},
			Result:       "NOT `id` IN (?,?,?)",
		},
		{
			Expr:         field.NewUint("", "id").Between(1, 10),
			ExpectedVars: []interface{}{uint(1), uint(10)},
			Result:       "`id` BETWEEN ? AND ?",
		},
		{
			Expr:   field.NewUint("", "id").Count(),
			Result: "COUNT(`id`)",
		},
		{
			Expr:   field.NewUint("", "id").Count().As("UserID"),
			Result: "COUNT(`id`) AS `UserID`",
		},
		{
			Expr:   field.NewUint("", "id").Distinct(),
			Result: "DISTINCT `id`",
		},
		{
			Expr:   field.NewUint("", "id").Distinct().Count(),
			Result: "COUNT(DISTINCT `id`)",
		},
		{
			Expr:   field.NewUint("", "id").Distinct().Count().As("UserID"),
			Result: "COUNT(DISTINCT `id`) AS `UserID`",
		},
		{
			Expr:         field.NewInt("", "age").RightShift(3),
			ExpectedVars: []interface{}{3},
			Result:       "`age`>>?",
		},
		{
			Expr:         field.NewInt("", "age").LeftShift(3),
			ExpectedVars: []interface{}{3},
			Result:       "`age`<<?",
		},
		// ======================== float ========================
		{
			Expr:         field.NewFloat64("", "score").Add(3.0),
			ExpectedVars: []interface{}{float64(3.0)},
			Result:       "`score`+?",
		},
		{
			Expr:         field.NewFloat64("", "score").Sub(3.0),
			ExpectedVars: []interface{}{float64(3.0)},
			Result:       "`score`-?",
		},
		{
			Expr:         field.NewFloat64("", "score").Mul(3.0),
			ExpectedVars: []interface{}{float64(3.0)},
			Result:       "`score`*?",
		},
		{
			Expr:         field.NewFloat64("", "score").Div(3.0),
			ExpectedVars: []interface{}{float64(3.0)},
			Result:       "`score`/?",
		},
		{
			Expr:         field.NewFloat64("", "score").FloorDiv(3.0),
			ExpectedVars: []interface{}{float64(3.0)},
			Result:       "`score` DIV ?",
		},
		// ======================== string ========================
		{
			Expr:         field.NewString("", "name").Eq("tom"),
			ExpectedVars: []interface{}{"tom"},
			Result:       "`name` = ?",
		},
		{
			Expr:         field.NewString("", "name").Neq("tom"),
			ExpectedVars: []interface{}{"tom"},
			Result:       "`name` <> ?",
		},
		{
			Expr:         field.NewString("", "name").Like("%%tom%%"),
			ExpectedVars: []interface{}{"%%tom%%"},
			Result:       "`name` LIKE ?",
		},
		{
			Expr:         field.NewString("", "name").NotLike("%%tom%%"),
			ExpectedVars: []interface{}{"%%tom%%"},
			Result:       "NOT `name` LIKE ?",
		},
		{
			Expr:         field.NewString("", "name").Regexp(".*"),
			ExpectedVars: []interface{}{".*"},
			Result:       "`name` REGEXP ?",
		},
		{
			Expr:         field.NewString("", "name").NotRegxp(".*"),
			ExpectedVars: []interface{}{".*"},
			Result:       "NOT `name` REGEXP ?",
		},
		// ======================== time ========================
		{
			Expr:         field.NewTime("", "creatAt").Eq(timeData),
			ExpectedVars: []interface{}{timeData},
			Result:       "`creatAt` = ?",
		},
		{
			Expr:         field.NewTime("", "creatAt").Gt(timeData),
			ExpectedVars: []interface{}{timeData},
			Result:       "`creatAt` > ?",
		},
		{
			Expr:         field.NewTime("", "creatAt").Gte(timeData),
			ExpectedVars: []interface{}{timeData},
			Result:       "`creatAt` >= ?",
		},
		{
			Expr:         field.NewTime("", "creatAt").Lt(timeData),
			ExpectedVars: []interface{}{timeData},
			Result:       "`creatAt` < ?",
		},
		{
			Expr:         field.NewTime("", "creatAt").Lte(timeData),
			ExpectedVars: []interface{}{timeData},
			Result:       "`creatAt` <= ?",
		},
		{
			Expr:         field.NewTime("", "creatAt").Between(timeData, timeData.Add(24*time.Hour)),
			ExpectedVars: []interface{}{timeData, timeData.Add(24 * time.Hour)},
			Result:       "`creatAt` BETWEEN ? AND ?",
		},
		{
			Expr:         field.NewTime("", "creatAt").Add(24 * time.Hour),
			ExpectedVars: []interface{}{time.Duration(24 * time.Hour).Microseconds()},
			Result:       "DATE_ADD(`creatAt`, INTERVAL ? MICROSECOND)",
		},
		{
			Expr:         field.NewTime("", "creatAt").Sub(24 * time.Hour),
			ExpectedVars: []interface{}{time.Duration(24 * time.Hour).Microseconds()},
			Result:       "DATE_SUB(`creatAt`, INTERVAL ? MICROSECOND)",
		},
		// ======================== bool ========================
		{
			Expr:   field.NewBool("", "male").Not(),
			Result: "NOT `male`",
		},
		{
			Expr:         field.NewBool("", "male").Is(true),
			ExpectedVars: []interface{}{true},
			Result:       "`male` IS ?",
		},
		{
			Expr:         field.NewBool("", "male").And(true),
			ExpectedVars: []interface{}{true},
			Result:       "`male` AND ?",
		},
		{
			Expr:         field.NewBool("", "male").Or(true),
			ExpectedVars: []interface{}{true},
			Result:       "`male` OR ?",
		},
	}

	for _, testcase := range testcases {
		checkBuildExpr(t, testcase.Expr, testcase.Result, testcase.ExpectedVars)
	}
}

func BenchmarkExpr_Count(b *testing.B) {
	id := field.NewUint("", "id")
	for i := 0; i < b.N; i++ {
		n := id.Count()
		_ = n
	}
}
