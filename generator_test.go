package gen

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/field"
)

func TestConfig(t *testing.T) {
	_ = &Config{
		db: nil,

		OutPath: "path",
		OutFile: "",

		ModelPkgPath: "models",

		queryPkgName: "query",
	}
}

// test data
type mysqlDialectors struct{ tests.DummyDialector }

func (mysqlDialectors) Name() string {
	return "mysql"
}

var db, _ = gorm.Open(mysqlDialectors{}, nil)

func init() {
	db = db.Debug()

	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		UpdateClauses: []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"},
		DeleteClauses: []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"},
	})
}

// User user data struct
type User struct {
	ID         uint `gorm:"primary_key"`
	Name       string
	Age        int
	Score      float64
	Address    string
	Famous     bool
	RegisterAt time.Time
}

func (User) TableName() string {
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

type user struct {
	userDo

	ALL        field.Asterisk
	ID         field.Uint
	Name       field.String
	Age        field.Int
	Score      field.Float64
	Address    field.String
	Famous     field.Bool
	RegisterAt field.Time
}

type userDo struct{ DO }

func (u userDo) Debug() *userDo {
	u.DO = *u.DO.Debug().(*DO)
	return &u
}

func (u userDo) WithContext(ctx context.Context) *userDo {
	u.DO = *u.DO.WithContext(ctx).(*DO)
	return &u
}

func (u userDo) Clauses(conds ...clause.Expression) *userDo {
	u.DO = *u.DO.Clauses(conds...).(*DO)
	return &u
}

func (u userDo) Not(conds ...Condition) *userDo {
	u.DO = *u.DO.Not(conds...).(*DO)
	return &u
}

func (u userDo) Or(conds ...Condition) *userDo {
	u.DO = *u.DO.Or(conds...).(*DO)
	return &u
}

func (u userDo) Select(conds ...field.Expr) *userDo {
	u.DO = *u.DO.Select(conds...).(*DO)
	return &u
}

func (u userDo) Where(conds ...Condition) *userDo {
	u.DO = *u.DO.Where(conds...).(*DO)
	return &u
}

func (u userDo) Order(conds ...field.Expr) *userDo {
	u.DO = *u.DO.Order(conds...).(*DO)
	return &u
}

func (u userDo) Distinct(cols ...field.Expr) *userDo {
	u.DO = *u.DO.Distinct(cols...).(*DO)
	return &u
}

func (u userDo) Omit(cols ...field.Expr) *userDo {
	u.DO = *u.DO.Omit(cols...).(*DO)
	return &u
}

func (u userDo) Join(table schema.Tabler, on ...field.Expr) *userDo {
	u.DO = *u.DO.Join(table, on...).(*DO)
	return &u
}

func (u userDo) LeftJoin(table schema.Tabler, on ...field.Expr) *userDo {
	u.DO = *u.DO.LeftJoin(table, on...).(*DO)
	return &u
}

func (u userDo) RightJoin(table schema.Tabler, on ...field.Expr) *userDo {
	u.DO = *u.DO.RightJoin(table, on...).(*DO)
	return &u
}

func (u userDo) Group(col field.Expr) *userDo {
	u.DO = *u.DO.Group(col).(*DO)
	return &u
}

func (u userDo) Having(conds ...Condition) *userDo {
	u.DO = *u.DO.Having(conds...).(*DO)
	return &u
}

func (u userDo) Limit(limit int) *userDo {
	u.DO = *u.DO.Limit(limit).(*DO)
	return &u
}

func (u userDo) Offset(offset int) *userDo {
	u.DO = *u.DO.Offset(offset).(*DO)
	return &u
}

func (u userDo) Scopes(funcs ...func(Dao) Dao) *userDo {
	u.DO = *u.DO.Scopes(funcs...).(*DO)
	return &u
}

func (u userDo) Unscoped() *userDo {
	u.DO = *u.DO.Unscoped().(*DO)
	return &u
}

func (u userDo) Create(values ...*user) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userDo) CreateInBatches(values []*user, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

func (u userDo) Save(values ...*user) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userDo) First() (*user, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*user), nil
	}
}

func (u userDo) Take() (*user, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*user), nil
	}
}

func (u userDo) Last() (*user, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*user), nil
	}
}

func (u userDo) Find() ([]*user, error) {
	result, err := u.DO.Find()
	return result.([]*user), err
}

func (u userDo) FindInBatches(result []*user, batchSize int, fc func(tx Dao, batch int) error) error {
	return u.DO.FindInBatches(&result, batchSize, fc)
}

func (u userDo) FindByPage(offset int, limit int) (result []*user, count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	result, err = u.Offset(offset).Limit(limit).Find()
	return
}

var u = func() *user {
	u := user{
		ALL:        field.NewAsterisk(""),
		ID:         field.NewUint("", "id"),
		Name:       field.NewString("", "name"),
		Age:        field.NewInt("", "age"),
		Score:      field.NewFloat64("", "score"),
		Address:    field.NewString("", "address"),
		Famous:     field.NewBool("", "famous"),
		RegisterAt: field.NewTime("", "register_at"),
	}
	u.UseDB(db.Session(&gorm.Session{Context: context.Background(), DryRun: true}))
	u.UseModel(User{})
	return &u
}()

type Student struct {
	DO

	ALL        field.Asterisk
	ID         field.Int64
	Name       field.String
	Age        field.Int
	Instructor field.Int64
}

var student = func() *Student {
	s := Student{
		ALL:        field.NewAsterisk("student"),
		ID:         field.NewInt64("student", "id"),
		Name:       field.NewString("student", "name"),
		Age:        field.NewInt("student", "age"),
		Instructor: field.NewInt64("student", "instructor"),
	}
	s.UseDB(db.Session(&gorm.Session{Context: context.Background(), DryRun: true}))
	s.UseModel(StudentRaw{})
	return &s
}()

type Teacher struct {
	DO

	ALL  field.Asterisk
	ID   field.Int64
	Name field.String
}

var teacher = func() Teacher {
	t := Teacher{
		ALL:  field.NewAsterisk("teacher"),
		ID:   field.NewInt64("teacher", "id"),
		Name: field.NewString("teacher", "name"),
	}
	t.UseDB(db.Session(&gorm.Session{Context: context.Background(), DryRun: true}))
	t.UseModel(TeacherRaw{})
	return t
}()
