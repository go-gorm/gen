package diy_method

import (
	"gorm.io/gen"
	"time"
)

type InsertMethod interface {

	// AddUser
	//
	// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
	AddUser(name string, age int) (gen.SQLResult, error)

	// AddUser1
	//
	// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
	AddUser1(name string, age int) (gen.RowsAffected, error)

	// AddUser2
	//
	// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
	AddUser2(name string, age int) gen.RowsAffected

	// AddUser3
	//
	// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
	AddUser3(name string, age int) gen.SQLResult

	// AddUser4
	//
	// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
	AddUser4(name string, age int) gen.SQLRow

	// AddUser5
	//
	// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
	AddUser5(name string, age int) gen.SQLRows

	// AddUser6
	//
	// INSERT INTO users (name,age) VALUES (@name,@age) ON DUPLICATE KEY UPDATE age=VALUES(age)
	AddUser6(name string, age int) (gen.SQLRows, error)
}

type SelectMethod interface {

	// FindByID
	//
	// select * from users where id=@id
	FindByID(id int) gen.T

	// LikeSearch
	//
	// SELECT * FROM @@table where name LIKE concat('%',@name,'%')
	LikeSearch(name string) *gen.T

	// InSearch
	//
	// select * from @@table where name in @names
	InSearch(names []string) []*gen.T

	// ColumnSearch
	//
	// select * from @@table where @@name in @names
	ColumnSearch(name string, names []string) []*gen.T
}

type TrimTest interface {

	// TestTrim
	//
	// select * from @@table where
	// {{trim}}
	// {{for _,name :=range list}}
	// name = @name or
	// {{end}}
	// {{end}}
	TestTrim(list []string) []gen.T

	// TestTrimInWhere
	//
	// select * from @@table
	// {{where}}
	// {{trim}}
	// {{for _,name :=range list}}
	// name = @name or
	// {{end}}
	// {{end}}
	// {{end}}
	TestTrimInWhere(list []string) []gen.T

	// TestInsert
	//
	// insert into users (name,age) values
	// {{trim}}
	// {{for key,value :=range data}}
	// (@key,@value),
	// {{end}}
	// {{end}}
	TestInsert(data gen.M) error
}

type TestIF interface {
	// FindByUsers
	//
	//select * from @@table
	//{{where}}
	//{{if user.Name !=""}}
	//name=@user.Name
	//{{end}}
	//{{end}}
	FindByUsers(user gen.T) []gen.T

	// FindByComplexIf
	//
	//select * from @@table
	//{{where}}
	//{{if user != nil && user.Name !=""}}
	//name=@user.Name
	//{{end}}
	//{{end}}
	FindByComplexIf(user *gen.T) []gen.T

	// FindByIfTime
	//
	// select * from @@table
	//{{if !start.IsZero()}}
	//created_at > start
	//{{end}}
	FindByIfTime(start time.Time) []gen.T
}

type TestFor interface {
	// TestFor
	//
	//select * from @@table where
	//{{for _,name:=range names}}
	//name = @name and
	//{{end}}
	//1=1
	TestFor(names []string) (gen.T, error)

	// TestForKey
	//
	// select * from @@table where
	//{{for _,name:=range names}}
	//or @@name = @value
	//{{end}}
	//and 1=1
	TestForKey(names []string, name, value string) (gen.T, error)

	// TestForOr
	//
	//select * from @@table
	//{{where}}
	//(
	//{{for _,name:=range names}}
	//name = @name or
	//{{end}}
	//{{end}}
	//)
	TestForOr(names []string) (gen.T, error)

	// TestIfInFor
	//
	//select * from @@table where
	//{{for _,name:=range names}}
	//{{if name !=""}}
	//name = @name or
	//{{end}}
	//{{end}}
	//1=2
	TestIfInFor(names []string, name string) (gen.T, error)

	// TestForInIf
	//
	//select * from @@table where
	//{{if name !="" }}
	//{{for _,forName:=range names}}
	//name = @forName or
	//{{end}}
	//{{end}}
	//1=2
	TestForInIf(names []string, name string) (gen.T, error)

	// TestForInWhere
	//
	//select * from @@table
	//{{where}}
	//{{for _,forName:=range names}}
	//or name = @forName
	//{{end}}
	//{{end}}
	TestForInWhere(names []string, name, forName string) (gen.T, error)

	// TestForUserList
	//
	//select * from users
	//{{where}}
	//{{for _,user :=range users}}
	//name=@user.Name
	//{{end}}
	//{{end}}
	TestForUserList(users []*gen.T, name string) (gen.T, error)

	// TestForMap
	//
	//select * from users
	//{{where}}
	//{{for key,value :=range param}}
	//@@key=@value
	//{{end}}
	//{{end}}
	TestForMap(param map[string]string, name string) (gen.T, error)

	// TestIfInIf
	//
	//select * from users
	//{{where}}
	//{{if name !="xx"}}
	//{{if name !="xx"}}
	//name=@name
	//{{end}}
	//{{end}}
	//{{end}}
	TestIfInIf(name string) gen.T

	// TestMoreFor
	//
	//select * from @@table
	//{{where}}
	//{{for _,name := range names}}
	//and name=@name
	//{{end}}
	//{{for _,id:=range ids}}
	//and id=@id
	//{{end}}
	//{{end}}
	TestMoreFor(names []string, ids []int) []gen.T

	// TestMoreFor2
	//
	//select * from @@table
	//{{where}}
	//{{for _,name := range names}}
	//OR (name=@name
	//{{for _,id:=range ids}}
	//and id=@id
	//{{end}}
	// and title !=@name)
	//{{end}}
	// {{end}}
	TestMoreFor2(names []string, ids []int) []gen.T

	// TestForInSet
	//
	// update @@table
	//{{set}}
	//{{for _,user:=range users}}
	//name=@user.Name,
	//{{end}}
	// {{end}} where
	TestForInSet(users []gen.T) error

	// TestInsertMoreInfo
	//
	// insert into @@table(name,age)values
	//{{for index ,user:=range users}}
	//{{if index >0}}
	//,
	//{{end}}
	//(@user.Name,@user.Age)
	//{{end}}
	TestInsertMoreInfo(users []gen.T) error

	// TestIfElseFor
	//
	// select * from @@table
	// {{where}}
	//{{if name =="admin"}}
	//(
	//{{for index,user:=range users}}
	//{{if index !=0}}
	//and
	//{{end}}
	//name like @user.Name
	//{{end}}
	//)
	//{{else if name !="guest"}}
	//{{for index,guser:=range users}}
	//{{if index ==0}}
	//(
	//{{else}}
	//and
	//{{end}}
	//name = @guser.Name
	//{{end}}
	//)
	//{{else}}
	//name ="guest"
	//{{end}}
	// {{end}}
	TestIfElseFor(name string, users []gen.T) error

	// TestForLike
	//
	// select * from @@table
	// {{where}}
	//{{for _,name:=range names}}
	//name like concat("%",@name,"%") or
	//{{end}}
	// {{end}}
	TestForLike(names []string) []gen.T
}
