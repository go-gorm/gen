// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"

	"gorm.io/plugin/dbresolver"

	"gorm.io/gen/tests/.gen/dal_4/model"

	"time"
)

func newUser(db *gorm.DB, opts ...gen.DOOption) user {
	_user := user{}

	_user.userDo.UseDB(db, opts...)
	_user.userDo.UseModel(&model.User{})

	tableName := _user.userDo.TableName()
	_user.ALL = field.NewAsterisk(tableName)
	_user.ID = field.NewInt64(tableName, "id")
	_user.CreatedAt = field.NewTime(tableName, "created_at")
	_user.Name = field.NewString(tableName, "name")
	_user.Address = field.NewString(tableName, "address")
	_user.RegisterTime = field.NewTime(tableName, "register_time")
	_user.Alive = field.NewBool(tableName, "alive")
	_user.CompanyID = field.NewInt64(tableName, "company_id")
	_user.PrivateURL = field.NewString(tableName, "private_url")

	_user.fillFieldMap()

	return _user
}

type user struct {
	userDo userDo

	ALL          field.Asterisk
	ID           field.Int64
	CreatedAt    field.Time
	Name         field.String // oneline
	Address      field.String
	RegisterTime field.Time
	/*
		multiline
		line1
		line2
	*/
	Alive      field.Bool
	CompanyID  field.Int64
	PrivateURL field.String

	fieldMap map[string]field.Expr
}

func (u user) Table(newTableName string) *user {
	u.userDo.UseTable(newTableName)
	return u.updateTableName(newTableName)
}

func (u user) As(alias string) *user {
	u.userDo.DO = *(u.userDo.As(alias).(*gen.DO))
	return u.updateTableName(alias)
}

func (u *user) updateTableName(table string) *user {
	u.ALL = field.NewAsterisk(table)
	u.ID = field.NewInt64(table, "id")
	u.CreatedAt = field.NewTime(table, "created_at")
	u.Name = field.NewString(table, "name")
	u.Address = field.NewString(table, "address")
	u.RegisterTime = field.NewTime(table, "register_time")
	u.Alive = field.NewBool(table, "alive")
	u.CompanyID = field.NewInt64(table, "company_id")
	u.PrivateURL = field.NewString(table, "private_url")

	u.fillFieldMap()

	return u
}

func (u *user) WithContext(ctx context.Context) IUserDo { return u.userDo.WithContext(ctx) }

func (u user) TableName() string { return u.userDo.TableName() }

func (u user) Alias() string { return u.userDo.Alias() }

func (u user) Columns(cols ...field.Expr) gen.Columns { return u.userDo.Columns(cols...) }

func (u *user) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := u.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (u *user) fillFieldMap() {
	u.fieldMap = make(map[string]field.Expr, 8)
	u.fieldMap["id"] = u.ID
	u.fieldMap["created_at"] = u.CreatedAt
	u.fieldMap["name"] = u.Name
	u.fieldMap["address"] = u.Address
	u.fieldMap["register_time"] = u.RegisterTime
	u.fieldMap["alive"] = u.Alive
	u.fieldMap["company_id"] = u.CompanyID
	u.fieldMap["private_url"] = u.PrivateURL
}

func (u user) clone(db *gorm.DB) user {
	u.userDo.ReplaceConnPool(db.Statement.ConnPool)
	return u
}

func (u user) replaceDB(db *gorm.DB) user {
	u.userDo.ReplaceDB(db)
	return u
}

type userDo struct{ gen.DO }

type IUserDo interface {
	gen.SubQuery
	Debug() IUserDo
	WithContext(ctx context.Context) IUserDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IUserDo
	WriteDB() IUserDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IUserDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IUserDo
	Not(conds ...gen.Condition) IUserDo
	Or(conds ...gen.Condition) IUserDo
	Select(conds ...field.Expr) IUserDo
	Where(conds ...gen.Condition) IUserDo
	Order(conds ...field.Expr) IUserDo
	Distinct(cols ...field.Expr) IUserDo
	Omit(cols ...field.Expr) IUserDo
	Join(table schema.Tabler, on ...field.Expr) IUserDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IUserDo
	RightJoin(table schema.Tabler, on ...field.Expr) IUserDo
	Group(cols ...field.Expr) IUserDo
	Having(conds ...gen.Condition) IUserDo
	Limit(limit int) IUserDo
	Offset(offset int) IUserDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IUserDo
	Unscoped() IUserDo
	Create(values ...*model.User) error
	CreateInBatches(values []*model.User, batchSize int) error
	Save(values ...*model.User) error
	First() (*model.User, error)
	Take() (*model.User, error)
	Last() (*model.User, error)
	Find() ([]*model.User, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.User, err error)
	FindInBatches(result *[]*model.User, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.User) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IUserDo
	Assign(attrs ...field.AssignExpr) IUserDo
	Joins(fields ...field.RelationField) IUserDo
	Preload(fields ...field.RelationField) IUserDo
	FirstOrInit() (*model.User, error)
	FirstOrCreate() (*model.User, error)
	FindByPage(offset int, limit int) (result []*model.User, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Rows() (*sql.Rows, error)
	Row() *sql.Row
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IUserDo
	UnderlyingDB() *gorm.DB
	schema.Tabler

	FindByUsers(user model.User) (result []model.User)
	FindByComplexIf(user *model.User) (result []model.User)
	FindByIfTime(start time.Time) (result []model.User)
	TestFor(names []string) (result model.User, err error)
	TestForKey(names []string, name string, value string) (result model.User, err error)
	TestForOr(names []string) (result model.User, err error)
	TestIfInFor(names []string, name string) (result model.User, err error)
	TestForInIf(names []string, name string) (result model.User, err error)
	TestForInWhere(names []string, name string, forName string) (result model.User, err error)
	TestForUserList(users []*model.User, name string) (result model.User, err error)
	TestForMap(param map[string]string, name string) (result model.User, err error)
	TestIfInIf(name string) (result model.User)
	TestMoreFor(names []string, ids []int) (result []model.User)
	TestMoreFor2(names []string, ids []int) (result []model.User)
	TestForInSet(users []model.User) (err error)
	TestInsertMoreInfo(users []model.User) (err error)
	TestIfElseFor(name string, users []model.User) (err error)
	TestForLike(names []string) (result []model.User)
	AddUser(name string, age int) (result sql.Result, err error)
	AddUser1(name string, age int) (rowsAffected int64, err error)
	AddUser2(name string, age int) (rowsAffected int64)
	AddUser3(name string, age int) (result sql.Result)
	AddUser4(name string, age int) (row *sql.Row)
	AddUser5(name string, age int) (rows *sql.Rows)
	AddUser6(name string, age int) (rows *sql.Rows, err error)
	FindByID(id int) (result model.User)
	LikeSearch(name string) (result *model.User)
	InSearch(names []string) (result []*model.User)
	ColumnSearch(name string, names []string) (result []*model.User)
}

// FindByUsers
//
// select * from @@table
// {{where}}
// {{if user.Name !=""}}
// name=@user.Name
// {{end}}
// {{end}}
func (u userDo) FindByUsers(user model.User) (result []model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	if user.Name != "" {
		params = append(params, user.Name)
		whereSQL0.WriteString("name=? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// FindByComplexIf
//
// select * from @@table
// {{where}}
// {{if user != nil && user.Name !=""}}
// name=@user.Name
// {{end}}
// {{end}}
func (u userDo) FindByComplexIf(user *model.User) (result []model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	if user != nil && user.Name != "" {
		params = append(params, user.Name)
		whereSQL0.WriteString("name=? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// FindByIfTime
//
// select * from @@table
// {{if !start.IsZero()}}
// created_at > start
// {{end}}
func (u userDo) FindByIfTime(start time.Time) (result []model.User) {
	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	if !start.IsZero() {
		generateSQL.WriteString("created_at > start ")
	}

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String()).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// TestFor
//
// select * from @@table where
// {{for _,name:=range names}}
// name = @name and
// {{end}}
// 1=1
func (u userDo) TestFor(names []string) (result model.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users where ")
	for _, name := range names {
		params = append(params, name)
		generateSQL.WriteString("name = ? and ")
	}
	generateSQL.WriteString("1=1 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestForKey
//
// select * from @@table where
// {{for _,name:=range names}}
// or @@name = @value
// {{end}}
// and 1=1
func (u userDo) TestForKey(names []string, name string, value string) (result model.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users where ")
	for _, name := range names {
		params = append(params, value)
		generateSQL.WriteString("or " + u.Quote(name) + " = ? ")
	}
	generateSQL.WriteString("and 1=1 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestForOr
//
// select * from @@table
// {{where}}
// (
// {{for _,name:=range names}}
// name = @name or
// {{end}}
// {{end}}
// )
func (u userDo) TestForOr(names []string) (result model.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	whereSQL0.WriteString("( ")
	for _, name := range names {
		params = append(params, name)
		whereSQL0.WriteString("name = ? or ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)
	generateSQL.WriteString(") ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestIfInFor
//
// select * from @@table where
// {{for _,name:=range names}}
// {{if name !=""}}
// name = @name or
// {{end}}
// {{end}}
// 1=2
func (u userDo) TestIfInFor(names []string, name string) (result model.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users where ")
	for _, name := range names {
		if name != "" {
			params = append(params, name)
			generateSQL.WriteString("name = ? or ")
		}
	}
	generateSQL.WriteString("1=2 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestForInIf
//
// select * from @@table where
// {{if name !="" }}
// {{for _,forName:=range names}}
// name = @forName or
// {{end}}
// {{end}}
// 1=2
func (u userDo) TestForInIf(names []string, name string) (result model.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users where ")
	if name != "" {
		for _, forName := range names {
			params = append(params, forName)
			generateSQL.WriteString("name = ? or ")
		}
	}
	generateSQL.WriteString("1=2 ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestForInWhere
//
// select * from @@table
// {{where}}
// {{for _,forName:=range names}}
// or name = @forName
// {{end}}
// {{end}}
func (u userDo) TestForInWhere(names []string, name string, forName string) (result model.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	for _, forName := range names {
		params = append(params, forName)
		whereSQL0.WriteString("or name = ? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestForUserList
//
// select * from users
// {{where}}
// {{for _,user :=range users}}
// name=@user.Name
// {{end}}
// {{end}}
func (u userDo) TestForUserList(users []*model.User, name string) (result model.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	for _, user := range users {
		params = append(params, user.Name)
		whereSQL0.WriteString("name=? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestForMap
//
// select * from users
// {{where}}
// {{for key,value :=range param}}
// @@key=@value
// {{end}}
// {{end}}
func (u userDo) TestForMap(param map[string]string, name string) (result model.User, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	for key, value := range param {
		params = append(params, value)
		whereSQL0.WriteString(u.Quote(key) + "=? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestIfInIf
//
// select * from users
// {{where}}
// {{if name !="xx"}}
// {{if name !="xx"}}
// name=@name
// {{end}}
// {{end}}
// {{end}}
func (u userDo) TestIfInIf(name string) (result model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	if name != "xx" {
		if name != "xx" {
			params = append(params, name)
			whereSQL0.WriteString("name=? ")
		}
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// TestMoreFor
//
// select * from @@table
// {{where}}
// {{for _,name := range names}}
// and name=@name
// {{end}}
// {{for _,id:=range ids}}
// and id=@id
// {{end}}
// {{end}}
func (u userDo) TestMoreFor(names []string, ids []int) (result []model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	for _, name := range names {
		params = append(params, name)
		whereSQL0.WriteString("and name=? ")
	}
	for _, id := range ids {
		params = append(params, id)
		whereSQL0.WriteString("and id=? ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// TestMoreFor2
//
// select * from @@table
// {{where}}
// {{for _,name := range names}}
// OR (name=@name
// {{for _,id:=range ids}}
// and id=@id
// {{end}}
// and title !=@name)
// {{end}}
// {{end}}
func (u userDo) TestMoreFor2(names []string, ids []int) (result []model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	for _, name := range names {
		params = append(params, name)
		whereSQL0.WriteString("OR (name=? ")
		for _, id := range ids {
			params = append(params, id)
			whereSQL0.WriteString("and id=? ")
		}
		params = append(params, name)
		whereSQL0.WriteString("and title !=?) ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// TestForInSet
//
// update @@table
// {{set}}
// {{for _,user:=range users}}
// name=@user.Name,
// {{end}}
// {{end}} where
func (u userDo) TestForInSet(users []model.User) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("update users ")
	var setSQL0 strings.Builder
	for _, user := range users {
		params = append(params, user.Name)
		setSQL0.WriteString("name=?, ")
	}
	helper.JoinSetBuilder(&generateSQL, setSQL0)
	generateSQL.WriteString("where ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestInsertMoreInfo
//
// insert into @@table(name,age)values
// {{for index ,user:=range users}}
// {{if index >0}}
// ,
// {{end}}
// (@user.Name,@user.Age)
// {{end}}
func (u userDo) TestInsertMoreInfo(users []model.User) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("insert into users(name,age)values ")
	for index, user := range users {
		if index > 0 {
			generateSQL.WriteString(", ")
		}
		params = append(params, user.Name)
		params = append(params, user.Age)
		generateSQL.WriteString("(?,?) ")
	}

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestIfElseFor
//
// select * from @@table
// {{where}}
// {{if name =="admin"}}
// (
// {{for index,user:=range users}}
// {{if index !=0}}
// and
// {{end}}
// name like @user.Name
// {{end}}
// )
// {{else if name !="guest"}}
// {{for index,guser:=range users}}
// {{if index ==0}}
// (
// {{else}}
// and
// {{end}}
// name = @guser.Name
// {{end}}
// )
// {{else}}
// name ="guest"
// {{end}}
// {{end}}
func (u userDo) TestIfElseFor(name string, users []model.User) (err error) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	if name == "admin" {
		whereSQL0.WriteString("( ")
		for index, user := range users {
			if index != 0 {
				whereSQL0.WriteString("and ")
			}
			params = append(params, user.Name)
			whereSQL0.WriteString("name like ? ")
		}
		whereSQL0.WriteString(") ")
	} else if name != "guest" {
		for index, guser := range users {
			if index == 0 {
				whereSQL0.WriteString("( ")
			} else {
				whereSQL0.WriteString("and ")
			}
			params = append(params, guser.Name)
			whereSQL0.WriteString("name = ? ")
		}
		whereSQL0.WriteString(") ")
	} else {
		whereSQL0.WriteString("name =\"guest\" ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	err = executeSQL.Error

	return
}

// TestForLike
//
// select * from @@table
// {{where}}
// {{for _,name:=range names}}
// name like concat("%",@name,"%") or
// {{end}}
// {{end}}
func (u userDo) TestForLike(names []string) (result []model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	generateSQL.WriteString("select * from users ")
	var whereSQL0 strings.Builder
	for _, name := range names {
		params = append(params, name)
		whereSQL0.WriteString("name like concat(\"%\",?,\"%\") or ")
	}
	helper.JoinWhereBuilder(&generateSQL, whereSQL0)

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// AddUser
//
// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
func (u userDo) AddUser(name string, age int) (result sql.Result, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	params = append(params, age)
	generateSQL.WriteString("INSERT INTO users (name,age) VALUES (?,?) ON DUPLICATE KEY UPDATE age=VALUES(age) ")

	stmt := u.UnderlyingDB().Statement
	result, err = stmt.ConnPool.ExecContext(stmt.Context, generateSQL.String(), params...) // ignore_security_alert

	return
}

// AddUser1
//
// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
func (u userDo) AddUser1(name string, age int) (rowsAffected int64, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	params = append(params, age)
	generateSQL.WriteString("INSERT INTO users (name,age) VALUES (?,?) ON DUPLICATE KEY UPDATE age=VALUES(age) ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	rowsAffected = executeSQL.RowsAffected
	err = executeSQL.Error

	return
}

// AddUser2
//
// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
func (u userDo) AddUser2(name string, age int) (rowsAffected int64) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	params = append(params, age)
	generateSQL.WriteString("INSERT INTO users (name,age) VALUES (?,?) ON DUPLICATE KEY UPDATE age=VALUES(age) ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Exec(generateSQL.String(), params...) // ignore_security_alert
	rowsAffected = executeSQL.RowsAffected

	return
}

// AddUser3
//
// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
func (u userDo) AddUser3(name string, age int) (result sql.Result) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	params = append(params, age)
	generateSQL.WriteString("INSERT INTO users (name,age) VALUES (?,?) ON DUPLICATE KEY UPDATE age=VALUES(age) ")

	stmt := u.UnderlyingDB().Statement
	result, _ = stmt.ConnPool.ExecContext(stmt.Context, generateSQL.String(), params...) // ignore_security_alert

	return
}

// AddUser4
//
// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
func (u userDo) AddUser4(name string, age int) (row *sql.Row) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	params = append(params, age)
	generateSQL.WriteString("INSERT INTO users (name,age) VALUES (?,?) ON DUPLICATE KEY UPDATE age=VALUES(age) ")

	row = u.UnderlyingDB().Raw(generateSQL.String(), params...).Row() // ignore_security_alert

	return
}

// AddUser5
//
// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
func (u userDo) AddUser5(name string, age int) (rows *sql.Rows) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	params = append(params, age)
	generateSQL.WriteString("INSERT INTO users (name,age) VALUES (?,?) ON DUPLICATE KEY UPDATE age=VALUES(age) ")

	rows, _ = u.UnderlyingDB().Raw(generateSQL.String(), params...).Rows() // ignore_security_alert

	return
}

// AddUser6
//
// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
func (u userDo) AddUser6(name string, age int) (rows *sql.Rows, err error) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	params = append(params, age)
	generateSQL.WriteString("INSERT INTO users (name,age) VALUES (?,?) ON DUPLICATE KEY UPDATE age=VALUES(age) ")

	rows, err = u.UnderlyingDB().Raw(generateSQL.String(), params...).Rows() // ignore_security_alert

	return
}

// FindByID
//
// select * from users where id=@id
func (u userDo) FindByID(id int) (result model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, id)
	generateSQL.WriteString("select * from users where id=? ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// LikeSearch
//
// SELECT * FROM @@table where name LIKE concat('%',@name,'%')
func (u userDo) LikeSearch(name string) (result *model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, name)
	generateSQL.WriteString("SELECT * FROM users where name LIKE concat('%',?,'%') ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Take(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// InSearch
//
// select * from @@table where name in @names
func (u userDo) InSearch(names []string) (result []*model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, names)
	generateSQL.WriteString("select * from users where name in ? ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

// ColumnSearch
//
// select * from @@table where @@name in @names
func (u userDo) ColumnSearch(name string, names []string) (result []*model.User) {
	var params []interface{}

	var generateSQL strings.Builder
	params = append(params, names)
	generateSQL.WriteString("select * from users where " + u.Quote(name) + " in ? ")

	var executeSQL *gorm.DB
	executeSQL = u.UnderlyingDB().Raw(generateSQL.String(), params...).Find(&result) // ignore_security_alert
	_ = executeSQL

	return
}

func (u userDo) Debug() IUserDo {
	return u.withDO(u.DO.Debug())
}

func (u userDo) WithContext(ctx context.Context) IUserDo {
	return u.withDO(u.DO.WithContext(ctx))
}

func (u userDo) ReadDB() IUserDo {
	return u.Clauses(dbresolver.Read)
}

func (u userDo) WriteDB() IUserDo {
	return u.Clauses(dbresolver.Write)
}

func (u userDo) Session(config *gorm.Session) IUserDo {
	return u.withDO(u.DO.Session(config))
}

func (u userDo) Clauses(conds ...clause.Expression) IUserDo {
	return u.withDO(u.DO.Clauses(conds...))
}

func (u userDo) Returning(value interface{}, columns ...string) IUserDo {
	return u.withDO(u.DO.Returning(value, columns...))
}

func (u userDo) Not(conds ...gen.Condition) IUserDo {
	return u.withDO(u.DO.Not(conds...))
}

func (u userDo) Or(conds ...gen.Condition) IUserDo {
	return u.withDO(u.DO.Or(conds...))
}

func (u userDo) Select(conds ...field.Expr) IUserDo {
	return u.withDO(u.DO.Select(conds...))
}

func (u userDo) Where(conds ...gen.Condition) IUserDo {
	return u.withDO(u.DO.Where(conds...))
}

func (u userDo) Order(conds ...field.Expr) IUserDo {
	return u.withDO(u.DO.Order(conds...))
}

func (u userDo) Distinct(cols ...field.Expr) IUserDo {
	return u.withDO(u.DO.Distinct(cols...))
}

func (u userDo) Omit(cols ...field.Expr) IUserDo {
	return u.withDO(u.DO.Omit(cols...))
}

func (u userDo) Join(table schema.Tabler, on ...field.Expr) IUserDo {
	return u.withDO(u.DO.Join(table, on...))
}

func (u userDo) LeftJoin(table schema.Tabler, on ...field.Expr) IUserDo {
	return u.withDO(u.DO.LeftJoin(table, on...))
}

func (u userDo) RightJoin(table schema.Tabler, on ...field.Expr) IUserDo {
	return u.withDO(u.DO.RightJoin(table, on...))
}

func (u userDo) Group(cols ...field.Expr) IUserDo {
	return u.withDO(u.DO.Group(cols...))
}

func (u userDo) Having(conds ...gen.Condition) IUserDo {
	return u.withDO(u.DO.Having(conds...))
}

func (u userDo) Limit(limit int) IUserDo {
	return u.withDO(u.DO.Limit(limit))
}

func (u userDo) Offset(offset int) IUserDo {
	return u.withDO(u.DO.Offset(offset))
}

func (u userDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IUserDo {
	return u.withDO(u.DO.Scopes(funcs...))
}

func (u userDo) Unscoped() IUserDo {
	return u.withDO(u.DO.Unscoped())
}

func (u userDo) Create(values ...*model.User) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Create(values)
}

func (u userDo) CreateInBatches(values []*model.User, batchSize int) error {
	return u.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (u userDo) Save(values ...*model.User) error {
	if len(values) == 0 {
		return nil
	}
	return u.DO.Save(values)
}

func (u userDo) First() (*model.User, error) {
	if result, err := u.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Take() (*model.User, error) {
	if result, err := u.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Last() (*model.User, error) {
	if result, err := u.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) Find() ([]*model.User, error) {
	result, err := u.DO.Find()
	return result.([]*model.User), err
}

func (u userDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.User, err error) {
	buf := make([]*model.User, 0, batchSize)
	err = u.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (u userDo) FindInBatches(result *[]*model.User, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return u.DO.FindInBatches(result, batchSize, fc)
}

func (u userDo) Attrs(attrs ...field.AssignExpr) IUserDo {
	return u.withDO(u.DO.Attrs(attrs...))
}

func (u userDo) Assign(attrs ...field.AssignExpr) IUserDo {
	return u.withDO(u.DO.Assign(attrs...))
}

func (u userDo) Joins(fields ...field.RelationField) IUserDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Joins(_f))
	}
	return &u
}

func (u userDo) Preload(fields ...field.RelationField) IUserDo {
	for _, _f := range fields {
		u = *u.withDO(u.DO.Preload(_f))
	}
	return &u
}

func (u userDo) FirstOrInit() (*model.User, error) {
	if result, err := u.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) FirstOrCreate() (*model.User, error) {
	if result, err := u.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.User), nil
	}
}

func (u userDo) FindByPage(offset int, limit int) (result []*model.User, count int64, err error) {
	result, err = u.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = u.Offset(-1).Limit(-1).Count()
	return
}

func (u userDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = u.Count()
	if err != nil {
		return
	}

	err = u.Offset(offset).Limit(limit).Scan(result)
	return
}

func (u userDo) Scan(result interface{}) (err error) {
	return u.DO.Scan(result)
}

func (u userDo) Delete(models ...*model.User) (result gen.ResultInfo, err error) {
	return u.DO.Delete(models)
}

func (u *userDo) withDO(do gen.Dao) *userDo {
	u.DO = *do.(*gen.DO)
	return u
}
