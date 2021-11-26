package main

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen/cmd/example/models"
	"gorm.io/gen/cmd/model"
	"gorm.io/gen/cmd/query"
	"gorm.io/gorm"
	"strings"
)

func main() {
	d, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	db := query.Use(d.Debug())
	u := db.User
	var user *models.User
	var users []*models.User
	var ctx = context.Background()
	var err error
	//res, err := u.WithContext(ctx).FindName()
	//fmt.Println(err, res)
	//fmt.Println("-------")
	//res, err := u.UpdateAffectedRows(12, "lzq")
	//fmt.Println(res)
	//fmt.Println("err:", err)
	//
	//res = u.UpdateRaws(13, "kang")
	//fmt.Println(res)
	//
	//users = u.ReturnAll()
	//fmt.Println(users)
	//
	//user = u.WithContext(ctx).ReturnOne()
	fmt.Println(user)
	//user, err = u.WithContext(ctx).FindOne(1000)
	//
	//fmt.Println(user)
	//fmt.Println("------")

	//user = u.WithContext(ctx).ReturnOne()
	//fmt.Println("ReturnOne", user)
	//
	//users = u.WithContext(ctx).ReturnAll()
	//fmt.Println("ReturnAll", users)
	//
	//str := u.WithContext(ctx).ReturnString()
	//fmt.Println("ReturnString", str)
	//
	//strs := u.WithContext(ctx).ReturnStrings()
	//fmt.Println("ReturnStrings", strs)
	//
	//count := u.WithContext(ctx).ReturnInt()
	//fmt.Println("ReturnInt", count)
	//
	//stu, _ := db.Student.WithContext(ctx).Find()
	//fmt.Println(stu)
	////res2 := u.GetMap()
	//fmt.Println(res2)
	//fmt.Println("------------------")
	//res, _ := u.GetStringMap([]int{2, 3, 4})
	//fmt.Println(reflect.TypeOf(res))
	//
	////s, _ := u.Where(u.Commit_.Eq("aaa")).Find()
	//s, _ = u.Where(u.Name.Eq("aa")).First()
	//fmt.Println(s)
	//
	////_ = u.UpdateName("ttt", 5)
	//
	//field := field2.NewString("xx", "yy-	`")
	//r, _ := u.Where(field.Eq("zz")).Find()
	//fmt.Println(r)

	fmt.Println("=++++++++++++++=")
	//stu := models.Students{Num: 1, Name: "zhangSan"}
	//d.AutoMigrate(&models.Students{})
	//result := d.Debug().Create(&stu)
	//fmt.Println(result)

	//user, err = u.FindByNameAndAge("test", 18) //nolint
	//if err != nil {
	//	err = fmt.Errorf("FindByNameAndAge field :%s", err)
	//	return
	//}
	//users, err = u.FindByID(5) //nolint
	//if err != nil {
	//	err = fmt.Errorf("FindByID field :%s", err)
	//	return
	//}
	//user, err = u.FindByIDOrKey(true, 2, "name", "Tom")
	//if err != nil {
	//	err = fmt.Errorf("FindByIDOrKey field :%s", err)
	//	return
	//}

	_, _ = user, users

	//users, err := db.User.Select(db.User.ID, db.User.Name, db.User.Age).Where(db.User.Age.Gt(18), db.User.ID.Gt(4)).Find()
	//if err != nil {
	//	fmt.Println("find users fail:", err)
	//	return
	//}
	//for _, user := range users {
	//	user.Display()
	//}
	////user, _ := db.User.FindByIDOrName(true, 2, "name", "")
	////fmt.Println(user)
	////
	////u := db.User.RCE("admin", 5, 5)
	////fmt.Println(u)
	//s := "      aa     bbb\n\n\n\n\n        ccc d  e"
	//fmt.Println(strings.Join(strings.Fields(s), " "))
	//
	//fmt.Println("------")
	////b := db.User.GetMap()
	////fmt.Println(b)
	////
	////w, _ := db.User.WhereInIF("admin", 0, 29)
	////fmt.Println(w)
	//
	//k, _ := db.User.FindByIDOrKey(true, 1, "name", "admin")
	//fmt.Println("kk", k)
	//fmt.Println(getStructName("xxxx"))
	////fmt.Println("================================")
	////name := "admin"
	////user1, _ := db.User.FindByNameAndAge(name, 26)
	////user1.Display()
	////
	////user2, err := db.User.InsertUser("gen", 1)
	////if err != nil {
	////	fmt.Println(err)
	////}
	////user2.Display()
	////c, err := db.User.GetCount()
	////fmt.Println(err)
	////fmt.Println(c)
	////
	////fmt.Println("db.User.GetMap")
	////m, err := db.User.GetMap()
	////fmt.Println(err)
	////for _, val := range m {
	////	val.Display()
	////}
	////fmt.Println(m)
	////
	////fmt.Println("db.People.GetMap")
	////p, err := db.People.GetMap()
	////fmt.Println(err)
	////for _, val := range p {
	////	fmt.Printf("people:%+v\n", val)
	//}
	//fmt.Println(p)
	names := []string{
		"admin",
		"test",
		"xxx",
	}
	ids := []int{1, 2, 3, 4, 5}
	us := []model.User{
		{
			Name: "gen",
			Age:  18,
		},
		{
			Name: "gorm",
			Age:  22,
		},
	}
	_ = names
	_ = ids
	_ = us
	res, err := u.WithContext(ctx).TestFor(names)
	//res := u.WithContext(ctx).TestMoreFor(name, ids)
	//res := u.WithContext(ctx).TestMoreFor2(name, ids)
	//err = u.WithContext(ctx).TestInsertMoreInfo(us)
	//res, err := u.WithContext(ctx).TestIfElseFor(name, "xxxxyi")
	fmt.Println(res)
	fmt.Println(err)

}
func getStructName(t string) string {
	list := strings.Split(t, ".")
	return list[len(list)-1]
}
