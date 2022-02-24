package field

import (
	"reflect"
	"strings"
	"sync"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"
)

var db, _ = gorm.Open(tests.DummyDialector{}, nil)

func GetStatement() *gorm.Statement {
	user, _ := schema.Parse(&User{}, &sync.Map{}, db.NamingStrategy)
	return &gorm.Statement{DB: db, Table: user.Table, Schema: user, Clauses: map[string]clause.Clause{}}
}

func CheckBuildExpr(t *testing.T, e Expr, result string, vars []interface{}) {
	stmt := GetStatement()

	e.expression().Build(stmt)

	sql := strings.TrimSpace(stmt.SQL.String())
	if sql != result {
		t.Errorf("SQL expects %v got %v", result, sql)
	}

	if !reflect.DeepEqual(stmt.Vars, vars) {
		t.Errorf("Vars expects %+v got %v", stmt.Vars, vars)
	}
}

func BuildToString(e Expr) (string, []interface{}) {
	stmt := GetStatement()

	e.expression().Build(stmt)

	return stmt.SQL.String(), stmt.Vars
}

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
