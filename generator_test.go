package gen

import (
	"database/sql"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"

	"gorm.io/gen/field"
)

func TestConfig(t *testing.T) {
	_ = &Config{
		db: nil,

		OutPath: "path",
		OutFile: "",

		ModelPkgName: "models",

		queryPkgName: "query",
	}
}

// test data

var db, _ = gorm.Open(tests.DummyDialector{}, nil)

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
	DO

	ID         field.Uint
	Name       field.String
	Age        field.Int
	Score      field.Float64
	Address    field.String
	Famous     field.Bool
	RegisterAt field.Time
}

func (u user) Debug() *user {
	u.DO = *u.DO.Debug().(*DO)
	return &u
}

func (u user) Hints(hs ...Hint) *user {
	u.DO = *u.DO.Hints(hs...).(*DO)
	return &u
}

func (u user) Not(conds ...Condition) *user {
	u.DO = *u.DO.Not(conds...).(*DO)
	return &u
}

func (u user) Or(conds ...Condition) *user {
	u.DO = *u.DO.Or(conds...).(*DO)
	return &u
}

func (u user) Select(conds ...field.Expr) *user {
	u.DO = *u.DO.Select(conds...).(*DO)
	return &u
}

func (u user) Where(conds ...Condition) *user {
	u.DO = *u.DO.Where(conds...).(*DO)
	return &u
}

func (u user) Order(conds ...field.Expr) *user {
	u.DO = *u.DO.Order(conds...).(*DO)
	return &u
}

func (u user) Distinct(cols ...field.Expr) *user {
	u.DO = *u.DO.Distinct(cols...).(*DO)
	return &u
}

func (u user) Omit(cols ...field.Expr) *user {
	u.DO = *u.DO.Omit(cols...).(*DO)
	return &u
}

func (u user) Join(table schema.Tabler, on ...Condition) *user {
	u.DO = *u.DO.Join(table, on...).(*DO)
	return &u
}

func (u user) LeftJoin(table schema.Tabler, on ...Condition) *user {
	u.DO = *u.DO.LeftJoin(table, on...).(*DO)
	return &u
}

func (u user) RightJoin(table schema.Tabler, on ...Condition) *user {
	u.DO = *u.DO.RightJoin(table, on...).(*DO)
	return &u
}

func (u user) Group(col field.Expr) *user {
	u.DO = *u.DO.Group(col).(*DO)
	return &u
}

func (u user) Having(conds ...Condition) *user {
	u.DO = *u.DO.Having(conds...).(*DO)
	return &u
}

func (u user) Limit(limit int) *user {
	u.DO = *u.DO.Limit(limit).(*DO)
	return &u
}

func (u user) Offset(offset int) *user {
	u.DO = *u.DO.Offset(offset).(*DO)
	return &u
}

func (u user) Scopes(funcs ...func(Dao) Dao) *user {
	u.DO = *u.DO.Scopes(funcs...).(*DO)
	return &u
}

func (u user) Unscoped() *user {
	u.DO = *u.DO.Unscoped().(*DO)
	return &u
}

func (u user) Create(value *User) error {
	return u.DO.Create(value)
}

func (u user) CreateInBatches(values []*User, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

func (u user) Save(value *User) error {
	return u.DO.Save(value)
}

func (u user) First(conds ...field.Expr) (*User, error) {
	result := new(User)
	if err := u.DO.First(result, conds...); err != nil {
		return nil, err
	}
	return result, nil
}

func (u user) Last(conds ...field.Expr) (*User, error) {
	result := new(User)
	if err := u.DO.Last(result, conds...); err != nil {
		return nil, err
	}
	return result, nil
}

func (u user) Take(conds ...field.Expr) (*User, error) {
	result := new(User)
	if err := u.DO.Take(result, conds...); err != nil {
		return nil, err
	}
	return result, nil
}

func (u user) Find(conds ...field.Expr) (result []*User, err error) {
	return result, u.DO.Find(&result, conds...)
}

func (u user) FindInBatches(result []*User, batchSize int, fc func(tx Dao, batch int) error) error {
	return u.DO.FindInBatches(&result, batchSize, fc)
}

func (u user) FindByPage(offset int, limit int) (result []*User, count int64, err error) {
	count, err = u.DO.Count()
	if err != nil {
		return
	}
	err = u.DO.Offset(offset).Limit(limit).Find(&result)
	return
}

func (u user) Delete(conds ...field.Expr) error {
	return u.DO.Delete(conds...)
}

func (u user) Begin(opts ...*sql.TxOptions) *user {
	u.DO = *u.DO.Begin(opts...).(*DO)
	return &u
}

func (u user) Commit() *user {
	u.DO = *u.DO.Commit().(*DO)
	return &u
}

func (u user) RollBack() *user {
	u.DO = *u.DO.Commit().(*DO)
	return &u
}

func (u user) SavePoint(name string) *user {
	u.DO = *u.DO.SavePoint(name).(*DO)
	return &u
}

func (u user) RollBackTo(name string) *user {
	u.DO = *u.DO.RollBackTo(name).(*DO)
	return &u
}

var u = func() *user {
	u := user{
		ID:         field.NewUint("", "id"),
		Name:       field.NewString("", "name"),
		Age:        field.NewInt("", "age"),
		Score:      field.NewFloat64("", "score"),
		Address:    field.NewString("", "address"),
		Famous:     field.NewBool("", "famous"),
		RegisterAt: field.NewTime("", "register_at"),
	}
	u.UseDB(db.Session(&gorm.Session{DryRun: true}))
	u.UseModel(User{})
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
