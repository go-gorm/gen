package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gorm.io/gen"
	"gorm.io/gen/cmd/example/models"
)

type name struct {
	N string
}
type Stu struct {
	Name string
	Age  int
}

type MainMethod interface {
	// Where("name=@name and age=@age")
	FindByNameAndAge(name string, age int) (gen.T, error)

	//sql(select id,name,age from users where age>18)
	FindBySimpleName() ([]models.User, error)

	//sql(select id,name,age from @@table where age>18
	//{{if cond1}}and id=@id {{end}}
	//{{if name == ""}}and @@col=@name{{end}})
	FindByIDOrName(cond1 bool, id int, col, name string) (gen.T, error)
}

func main() {
	db, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g := gen.NewGenerator(gen.Config{
		OutPath: "./cmd/query/",
		//ModelPkgPath: "./cmd/model/v2",
		Mode: gen.WithoutContext | gen.WithDefaultQuery,
	})
	g.UseDB(db)

	//g.TableNames("login", "Users")
	//g.Tables(models.User{}, &models.People{})
	tblUser := g.GenerateModel("users")

	//g.ApplyBasic(tblUser, models.People{})
	//g.ApplyBasic(model.Person{})

	g.ApplyInterface(func(testFor models.TestFor, test models.Model, method models.TestMethod, test2 models.Test) {}, tblUser, g.GenerateModel("people"))
	//g.ApplyInterface(func(models.MapModel, models.Model) {}, g.GenerateModel("students"), Stu{})
	//g.ApplyByModel(models.User{}, func(method models.TestMethod, mainMethod MainMethod) {})

	//g.ApplyInterface(func(model models.TestNumber) {}, tblUser)

	//g.ApplyInterface(func(models.UserMethod, models.MapModel, models.TestReturnMethod) {}, models.User{})
	//g.ApplyInterface(func(models.Model) {}, models.User{}, g.GenerateModel("login"), g.GenerateModel("students", gen.FieldType("created_at", "models.Timestamp")))
	//g.ApplyInterface(func(model models.TestMethod) {}, models.User{}, g.GenerateModel("students"))
	//g.ApplyInterface(func(model models.TestMethod) {}, g.GenerateModel("login"))
	g.Execute()

}
