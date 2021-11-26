package models

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gen"
	"gorm.io/gorm"
	"time"
)

type TestStruct struct {
}

type Use1r struct {
	ID   int    `gorm:"primary_key"`
	Name string `gorm:"column:name"`
	Age  int    `gorm:"column:age;type:varchar(64)"`
}

type User struct {
	gorm.Model
	ID       int    `gorm:"primary_key"`
	Name     string `gorm:"column:name"`
	Age      int    `gorm:"column:age;type:varchar(64)"`
	Birthday time.Time
	GetMap   string
}

type People struct {
	gorm.Model
	Name    string
	age     int //non-exportable member will be ignored
	High    float64
	IsAdult bool
}
type Student struct {
	Pass1 *Password
	Pass2 Password
	Pass3 Password
	Pass4 PassPtr
	Pass5 *PassPtr
	//Pass6 PassArray
	Pass7 Pwd
	Pass8 *Pwd
	//Name  People
	Info Info
}

func _() {
	//s := new(Students)
	//s.Pass2
}

type NName string
type Info string

func (p *Info) Value() {
	fmt.Println("valuer")
}

type PassArray []Password

type PassPtr *Password

type Password string

func (p *Password) Scan(src interface{}) error {
	//*p = Password(fmt.Sprintf("@%v@", src))
	return nil
}

func (p Password) Value() (driver.Value, error) {
	//p = Password(strings.Trim(string(p), "@"))
	return nil, nil
}

type Pwd string

func (p Pwd) Scan(name interface{}) error {
	//*p = Password(fmt.Sprintf("@%v@", src))
	return nil
}

func (p Pwd) Value() error {
	//p = Password(strings.Trim(string(p), "@"))
	return nil
}
func (p Pwd) Test() (driver.Value, error) {
	//p = Password(strings.Trim(string(p), "@"))
	return nil, nil
}

type Test struct {
	TestInt     int
	TestInt8    int8
	TestInt16   int16
	TestInt32   int32
	TestInt64   int64
	TestUint    uint
	TestUint8   uint8
	TestUint16  *uint16
	TestUint32  uint32
	TestUint64  uint64
	TestString  string
	TestByte    byte
	TestFloat32 float32
	TestFloat64 float64
	TestBool    time.Time
}

type Sstring struct {
	TestString  string
	TestByte    byte
	TestFloat32 float32
	TestFloat64 float64
	TestBool    bool
}
type MapModel interface {
	////sql(select * from @@table)
	//GetMap1() gen.M

	// GetStringMap sql(select name from @@table where id in @ids and name=@names and key =@key)
	GetStringMap(ids []int, names string, key string) ([]map[string]interface{}, error)

	//select name from @@table where id=5
	GetInterface() (error, Tester)
}

type Model interface {

	// Where("name=@name and age=@age")
	WhereFindByName(name, col string, age int) (gen.T, error)

	// Where("name=@names and age=@ages")
	FindByNameAndAge(names []string, ages interface{}) ([]gen.T, error)

	//sql("select id,name,age from @@table where age>18")
	FindBySimpleName() (gen.T, error)

	//sql(select count(*) from @@table)
	GetCount() (int, error)

	//sql(select Birthday from @@table {{where}} id=@id{{end}})
	GotBirthDay(id int) (time.Time, error)

	//sql(select * from @@table)
	GetMap1() map[string]interface{}

	//sql(insert into @@table (name,age)values(@name,@age))
	InsertUser(name string, age int) error

	/*sql(select * from @@table
		{{if name=="admin"}}
			{{where}}
				id>0
				{{if age>18}}
					and age>18
				{{end}}
			{{end}}
		{{else if name=="root"}}
			{{where}}
				id>10
			{{end}}
	{{else if name=="hello"}}
			{{where}}
				id>16
			{{end}}
		{{else}}
			{{where}}
				id>50
			{{end}}
			{{if name=="user"}}
				and 1=1
			{{end}}
			and name = admin
		{{end}}
	)*/
	WhereInIF(name string, id, age int) (gen.T, error)

	//sql(update @@table {{set}}{{if name!=""}}name=@name{{end}}{{end}}where id=@id)
	HasSet(name string, id int) error

	/*
		select * from @@table
			{{if name=="aa{{where}}aa"}}
				id=@id
				{{if age>18 }}
					and age=@age
				{{else}}
					and name="{{where}}name"
				{{end}}
				{{if name=="adm\nin"}}
					and id>30
				{{else if name=="sss"}}
					and id>20
				{{end}}
			{{end}}
	*/
	RCE(name string, age, id int) gen.T

	//update @@table
	//	{{set}}
	//		update_time=now(),
	//		{{if name != ""}}
	//			name=@name
	//		{{end}}
	//	{{end}}
	//	{{where}}
	//		id=@id
	//	{{end}}
	UpdateName(name string, id int) error

	/*
			sql(select id,name,age from @@table where name="aaax

		\"xxx" and age >18
				{{if cond1}} and id=true  {{end}}
				{{if name != ""}}
				and @@column=@name{{end}}
			)
	*/
	FindByIDOrName(cond1 bool, id int, column, name string) (gen.T, error)

	//// Where("name=@name and age=@age")
	//FindByNameAndAge(name string, age int) (gen.T, error)
	////sql(select id,name,age from users where age>18)
	//FindBySimpleName() ([]gen.T, error)
	//
	////sql(select id,name,age from @@table where age>18
	////{{if cond1}}and id=@id {{end}}
	////{{if name == ""}}and @@col=@name{{end}})
	//FindByIDOrName(cond1 bool, id int, col, name string) (gen.T, error)
}

type UserMethod interface {
	// Where("name=@name and age=@age")
	FindUserByNameAndAge(name string, age int) (User, error)

	//sql(select id,name,age from users where age>18)
	FindUserBySimpleName() ([]gen.T, error)
}

type TestMethod interface {
	//// Where("name=@name and age=@age")
	//FindByNameAndAge(name string, age int) (gen.T, error)
	//
	////sql(select id,name,age from users where age>18)
	//FindBySimpleName() ([]gen.T, error)

	//select * from @@table
	//   {{where}}
	//		  id>0
	//        {{if cond}}id=@id {{end}}
	//        {{   if value != ""}}or @@key=@value{{end}}
	//    {{end}}
	FindByIDOrKey(cond bool, id int, key, value string) (gen.T, error)

	//select group_concat(name) from @@table
	FindNames() (string, error)

	//select * from @@table where id>@id
	FindOne(id int) (gen.T, error)
}
type Method interface {
	//// Where("name=@name and age=@age")
	FindByNameAndAge1(name string, age int) (gen.T, error)

	//sql(select * from user where id=@id)
	FindByID(id int) ([]gen.T, error)

	// SELECT * /* hint */ FROM `users`
	UseHint() []gen.T

	//select name from @@table
	FindName() (string, error)

	//select * from @@table
	//	{{where}}
	//		name=@name
	//	{{end}}
	WhereFind(name string) gen.T

	////select * from @@table
	////   {{where}}
	////        {{if cond}}id=@id {{end}}
	////        {{if value == nil}} or name=@value{{end}}
	////    {{end}}
	FindByIDOrKey2(cond bool, id int, key, value *Value) (gen.T, error)

	//select id,name,age from @@table
	//	{{where}}
	//		age>18,
	//		{{if ids != nil}}and id in @ids {{end}}
	//		{{if ids == nil}}and name=@nameString {{end}}
	// {{end}}
	FindByNotNil(ids []int, nameString string) (gen.T, error)

	//select * from users
	//{{where}}
	//		{{if name !=""}}
	//			{{if name =="admin"}}
	//				@@key=@val
	//			{{else}}
	//				name=@val
	//			{{end}}
	//		{{else}}
	//			@@key !=@val
	//		{{end}}
	//{{end}}
	FindCom(name, key, val string) gen.T

	//select id,name,age from @@table
	//	{{if cond}}
	//		{{where}}
	//			{{if name !=""}}
	//				name=@name
	//			{{end}}
	//		{{end}}
	//	{{else}}
	//		{{where}}
	//			@@key is NOT Null
	//			{{if name =="admin"}}
	//			{{end}}
	//		{{end}}
	//	{{end}}
	FindSomeWhere(cond bool, name, key, val string) gen.T
}

type Value string
type TestReturnMethod interface {
	//select * from @table where id>0
	ReturnOne() gen.T

	//select * from @table where id>10
	ReturnAll() []gen.T

	//select * from @table where id>10 /* comment */
	ReturnMap() []gen.M

	//select name from @table where id>10
	ReturnString() string

	//select name from @table where id>10
	ReturnStrings() []string

	//select count(*) from @table
	ReturnInt() int
}
type TestUpdate interface {

	//update @@table set name=@name where id=@id
	UpdateName(id int, name string) error

	//update @@table set name=@name where id=@id
	OnlyUpdate(id int, name string)

	/*
		update @@table set
		name=@name where id=@id and name like '%\n\n%'
	*/
	UpdateRaws(id int, name string) gen.RowsAffected

	// update @@table set name=@name where id=@id
	UpdateAffectedRows(id int, name string) (gen.RowsAffected, error)
}

type TestFor interface {

	//select * from users where age>18
	TestSimple() []gen.T

	// select * from @@table where
	//	{{for _,name:=range names}}
	//		name = @name and
	//{{end}}
	//1=1
	TestFor(names []string) (gen.T, error)

	//select * from @@table where
	//	{{for _,name:=range names}}
	//		@@name = @value and
	//{{end}}
	//1=1
	TestForKey(names []string, name, value string) (gen.T, error)

	// select * from @@table where
	//	{{for _,name:=range names}}
	//		{{if name !=""}}
	//			name = @name or
	//		{{end}}
	//{{end}}
	//1=2
	TestIfInFor(names []string, name string) (gen.T, error)

	// select * from @@table where
	//	{{if name !="" }}
	//		{{for _,forName:=range names}}
	//			name = @forName or
	//		{{end}}
	//{{end}}
	//1=2
	TestForInIf(names []string, name string) (gen.T, error)

	// select * from @@table
	//	{{where}}
	//		{{for _,forName:=range names}}
	//			or name = @forName
	//		{{end}}
	//{{end}}
	TestForInWhere(names []string, name, forName string) (gen.T, error)

	//select * from users
	//{{where}}
	//	{{for _,user :=range users}}
	//		name=@user.Name
	//	{{end}}
	//{{end}}
	TestForUserList(users []gen.T, name string) (gen.T, error)

	//select * from users
	//{{where}}
	//  {{if name !="xx"}}
	//		{{if name !="xx"}}
	//			name=@name
	//		{{end}}
	// {{end}}
	//{{end}}
	TestIfInIf(name string) gen.T

	// select * from @@table
	//	{{where}}
	//		{{for _,name := range names}}
	//			and name=@name
	//		{{end}}
	//		{{for _,id:=range ids}}
	//			and id=@id
	//		{{end}}
	//{{end}}
	TestMoreFor(names []string, ids []int) []gen.T

	// select * from @@table
	//	{{where}}
	//		{{for _,name := range names}}
	//			OR (name=@name
	//			{{for _,id:=range ids}}
	//				and id=@id
	//			{{end}}
	//			 and title !=@name)
	//		{{end}}
	//{{end}}
	TestMoreFor2(names []string, ids []int) []gen.T

	// update @@table
	//	{{set}}
	//		{{for _,user:=range users}}
	//			name=@user.Name,
	//		{{end}}
	//{{end}} where
	TestForInSet(users []gen.T) error

	// insert into @@table(name,age)values
	//		{{for index ,user:=range users}}
	//			{{if index >0}}
	//				,
	//			{{end}}
	//			(@user.Name,@user.Age)
	//		{{end}}
	TestInsertMoreInfo(users []gen.T) error

	//select * from @@table
	//{{where}}
	//	{{if name =="admin"}}
	//		(
	//		{{for index,user:=range users}}
	//			{{if index !=0}}
	//				and
	//			{{end}}
	//			name like @user.Name
	//		{{end}}
	//		)
	//	{{else if name !="guest"}}
	//		{{for index,guser:=range users}}
	//			{{if index ==0}}
	//				(
	//			{{else}}
	//				and
	//			{{end}}
	//			name = @guser.Name
	//		{{end}}
	//		)
	//	{{else}}
	//		name ="guest"
	//	{{end}}
	//{{end}}
	TestIfElseFor(name string, users []gen.T) error
}
type TestNumber interface {

	//select * from users where
	//{{if id>1}}
	// 	id=@id
	//{{else if  id == 2 *3}}
	// 	uid=@id
	//	{{else}}
	//   {{if len(name)>0}}
	//		name=@name
	// {{end}}
	//{{end}}
	TestNum(id int64, name string) gen.T
}

func (u User) Display() {
	value := u.Name
	fmt.Printf("userid:%d,Name:%s,Age:%f\n", u.ID, value, 2)
}

func test() {
	//var name string
	//if name != "" {
	//	 if name == "admin" {
	//	 		WhereClause0.WriteString(" u.Quote(key)=@val ")
	//	  } else {
	//		  WhereClause0.WriteString(" name=@val ")
	//	  }
	//   } else {
	// 	 WhereClause0.WriteString(" u.Quote(key) !=@val ")
	//    }
	// }

}
