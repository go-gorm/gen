# GORM/GEN

[![GoVersion](https://img.shields.io/github/go-mod/go-version/go-gorm/gen)](https://github.com/go-gorm/gen/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/go-gorm/gen)](https://github.com/go-gorm/gen/releases)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/gorm.io/gen?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-gorm/gen)](https://goreportcard.com/report/github.com/go-gorm/gen)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![OpenIssue](https://img.shields.io/github/issues/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aopen+is%3Aissue)
[![ClosedIssue](https://img.shields.io/github/issues-closed/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aissue+is%3Aclosed)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/go-gorm/gen)](https://www.tickgit.com/browse?repo=github.com/go-gorm/gen)

基于 [GORM](https://github.com/go-gorm/gorm), 更安全更友好的ORM工具。

## Overview

- 自动生成CRUD和DIY方法
- 自动根据表结构生成model
- 完全兼容GORM
- 更安全、更友好
- 多种生成代码模式

## Contents


## 安装

安装GEN前，需要安装好GO并配置你的工作环境。

1.安装完Go(version 1.14+)之后，通过下面的命令安装gen。

```bash
go get -u gorm.io/gen
```

2.导入到你的工程:

```go
import "gorm.io/gen"
```

## 快速开始

**注⚠️**: 这里所有的教程都是在 `WithContext` 模式下写的. 如果你用的是`WithoutContext` 模式,则可以删除所有的 `WithContext(ctx)` ，这样代码看起来会更简洁.

```bash
# assume the following code in generate.go file
$ cat generate.go
```

```go
package main

import "gorm.io/gen"

// generate code
func main() {
    // specify the output directory (default: "./query")
    // ### if you want to query without context constrain, set mode gen.WithoutContext ###
    g := gen.NewGenerator(gen.Config{
        OutPath: "../dal/query",
        /* Mode: gen.WithoutContext|gen.WithDefaultQuery*/
        //if you want the nullable field generation property to be pointer type, set FieldNullable true
        /* FieldNullable: true,*/
        //if you want to generate index tags from database, set FieldWithIndexTag true
        /* FieldWithIndexTag: true,*/
        //if you want to generate type tags from database, set FieldWithTypeTag true
        /* FieldWithTypeTag: true,*/
        //if you need unit tests for query code, set WithUnitTest true
        /* WithUnitTest: true, */
    })
  
    // reuse the database connection in Project or create a connection here
    // if you want to use GenerateModel/GenerateModelAs, UseDB is necessray or it will panic
    // db, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
    g.UseDB(db)
  
    // apply basic crud api on structs or table models which is specified by table name with function
    // GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
    g.ApplyBasic(model.User{}, g.GenerateModel("company"), g.GenerateModelAs("people", "Person", gen.FieldIgnore("address")))
    
    // apply diy interfaces on structs or table models
    g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))

    // execute the action of code generation
    g.Execute()
}
```

生成Model:

- `gen.WithoutContext` 非 `WithContext` 模式生成
- `gen.WithDefaultQuery` 生成默认全局查询变量

### 项目路径

最佳实践项目模板:

```bash
demo
├── cmd
│   └── generate
│       └── generate.go # execute it will generate codes
├── dal
│   ├── dal.go # create connections with database server here
│   ├── model
│   │   ├── method.go # DIY method interfaces
│   │   └── model.go  # store struct which corresponding to the database table
│   └── query  # generated code's directory
|       ├── user.gen.go # generated code for user
│       └── gen.go # generated code
├── biz
│   └── query.go # call function in dal/gorm_generated.go and query databases
├── config
│   └── config.go # DSN for database server
├── generate.sh # a shell to execute cmd/generate
├── go.mod
├── go.sum
└── main.go
```

## API 示例

### 生成

#### 生成Model

```go
// generate a model struct map to table `people` in database
g.GenerateModel("people")

// generate a struct and specify struct's name
g.GenerateModelAs("people", "People")

// add option to ignore field
g.GenerateModel("people", gen.FieldIgnore("address"), gen.FieldType("id", "int64"))

// generate all tables, ex: g.ApplyBasic(g.GenerateAllTable()...)
g.GenerateAllTable()
```

字段生成 **Options**

```go
FieldNew           // create new field
FieldIgnore        // ignore field
FieldIgnoreReg     // ignore field (match with regexp)
FieldRename        // rename field in struct
FieldType          // specify field type
FieldTypeReg       // specify field type (match with regexp)
FieldTag           // specify gorm and json tag
FieldJSONTag       // specify json tag
FieldGORMTag       // specify gorm tag
FieldNewTag        // append new tag
FieldNewTagWithNS  // specify new tag with name strategy
FieldTrimPrefix    // trim column prefix
FieldTrimSuffix    // trim column suffix
FieldAddPrefix     // add prefix to struct member's name
FieldAddSuffix     // add suffix to struct member's name
FieldRelate        // specify relationship with other tables
FieldRelateModel   // specify relationship with exist models
```

#### 类型映射

自定义数据库字段类型和go类型的映射关系.

```go
dataMap := map[string]func(detailType string) (dataType string){
  "int": func(detailType string) (dataType string) { return "int64" },
  // bool mapping
  "tinyint": func(detailType string) (dataType string) {
    if strings.HasPrefix(detailType, "tinyint(1)") {
      return "bool"
    }
    return "int8"
  },
}

g.WithDataTypeMap(dataMap)
```

### 字段表达式

#### 创建字段

实际上你需要手动创建字段，因为都会在生成代码自动创建。

| Field Type | Detail Type           | Create Function               | Supported Query Method                                       |
| ---------- | --------------------- | ------------------------------ | ------------------------------------------------------------ |
| generic    | field                 | NewField                       | IsNull/IsNotNull/Count/Eq/Neq/Gt/Gte/Lt/Lte/Like             |
| int        | int/int8/.../int64    | NewInt/NewInt8/.../NewInt64    | Eq/Neq/Gt/Gte/Lt/Lte/In/NotIn/Between/NotBetween/Like/NotLike/Add/Sub/Mul/Div/Mod/FloorDiv/RightShift/LeftShift/BitXor/BitAnd/BitOr/BitFlip |
| uint       | uint/uint8/.../uint64 | NewUint/NewUint8/.../NewUint64 | same with int                                                |
| float      | float32/float64       | NewFloat32/NewFloat64          | Eq/Neq/Gt/Gte/Lt/Lte/In/NotIn/Between/NotBetween/Like/NotLike/Add/Sub/Mul/Div/FloorDiv |
| string     | string/[]byte         | NewString/NewBytes             | Eq/Neq/Gt/Gte/Lt/Lte/Between/NotBetween/In(val/NotIn(val/Like/NotLike/Regexp/NotRegxp/FindInSet/FindInSetWith |
| bool       | bool                  | NewBool                        | Not/Is/And/Or/Xor/BitXor/BitAnd/BitOr                        |
| time       | time.Time             | NewTime                        | Eq/Neq/Gt/Gte/Lt/Lte/Between/NotBetween/In/NotIn/Add/Sub     |

创建字段示例:

```go
import "gorm.io/gen/field"

// create a new generic field map to `generic_a`
a := field.NewField("table_name", "generic_a")

// create a field map to `id`
i := field.NewInt("user", "id")

// create a field map to `address`
s := field.NewString("user", "address")

// create a field map to `create_time`
t := field.NewTime("user", "create_time")
```

### CRUD 接口

生成基础model `user` 和 `DB`.

```go
// generated code
// generated code
// generated code
package query

import "gorm.io/gen"

// struct map to table `users` 
type user struct {
    gen.DO
    ID       field.Uint
    Name     field.String
    Age      field.Int
    Address  field.Field
    Birthday field.Time
}

// struct collection
type DB struct {
    db       *gorm.DB
    User     *user
}
```

#### 创建

##### 创建记录

```go
// u refer to query.user
user := model.User{Name: "Modi", Age: 18, Birthday: time.Now()}

u := query.Use(db).User
err := u.WithContext(ctx).Create(&user) // pass pointer of data to Create

err // returns error
```

##### 选择字段创建

自定义哪些字段需要插入。

```go
u := query.Use(db).User
u.WithContext(ctx).Select(u.Name, u.Age).Create(&user)
// INSERT INTO `users` (`name`,`age`) VALUES ("modi", 18)
```

自定义创建时需要忽略的字段。

```go
u := query.Use(db).User
u.WithContext(ctx).Omit(u.Name, u.Age).Create(&user)
// INSERT INTO `users` (`Address`, `Birthday`) VALUES ("2021-08-17 20:54:12.000", 18)
```

##### 批量创建

 `Create` 方法支持批量创建，参数只要是对应model的slice就可以. GORM会通过一条语句高效创建并返回所有的主键赋值给slice的Model.

```go
var users = []model.User{{Name: "modi"}, {Name: "zhangqiang"}, {Name: "songyuan"}}
query.Use(db).User.WithContext(ctx).Create(&users)

for _, user := range users {
    user.ID // 1,2,3
}
```

`CreateInBatches`可以指定批量创建的大小, e.g:

```go
var users = []User{{Name: "modi_1"}, ...., {Name: "modi_10000"}}

// batch size 100
query.Use(db).User.WithContext(ctx).CreateInBatches(users, 100)
```

也可以通过全局配置方式，在初始化gorm时设置 `CreateBatchSize` in `gorm.Config` / `gorm.Session`

```go
db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
    CreateBatchSize: 1000,
})
// OR
db = db.Session(&gorm.Session{CreateBatchSize: 1000})

u := query.NewUser(db)

var users = []User{{Name: "modi_1"}, ...., {Name: "modi_5000"}}

u.WithContext(ctx).Create(&users)
// INSERT INTO users xxx (5 batches)
```

#### 查询

##### 单个数据查询

自动生成 `First`, `Take`, `Last` 三个查询单条数据的方法。 执行的sql后面会自动添加 `LIMIT 1` ，如果没有查到数据会返回错误： `ErrRecordNotFound` 。

```go
u := query.Use(db).User

// Get the first record ordered by primary key
user, err := u.WithContext(ctx).First()
// SELECT * FROM users ORDER BY id LIMIT 1;

// Get one record, no specified order
user, err := u.WithContext(ctx).Take()
// SELECT * FROM users LIMIT 1;

// Get last record, ordered by primary key desc
user, err := u.WithContext(ctx).Last()
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

// check error ErrRecordNotFound
errors.Is(err, gorm.ErrRecordNotFound)
```

##### 根据主键查询数据

```go
u := query.Use(db).User

user, err := u.WithContext(ctx).Where(u.ID.Eq(10)).First()
// SELECT * FROM users WHERE id = 10;

users, err := u.WithContext(ctx).Where(u.ID.In(1,2,3)).Find()
// SELECT * FROM users WHERE id IN (1,2,3);
```

如果是string类型的主键，比如UUID等:

```go
user, err := u.WithContext(ctx).Where(u.ID.Eq("1b74413f-f3b8-409f-ac47-e8c062e3472a")).First()
// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
```

##### 查询所有数据

```go
u := query.Use(db).User

// Get all records
users, err := u.WithContext(ctx).Find()
// SELECT * FROM users;
```

##### 条件

###### 基础查询

```go
u := query.Use(db).User

// Get first matched record
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).First()
// SELECT * FROM users WHERE name = 'modi' ORDER BY id LIMIT 1;

// Get all matched records
users, err := u.WithContext(ctx).Where(u.Name.Neq("modi")).Find()
// SELECT * FROM users WHERE name <> 'modi';

// IN
users, err := u.WithContext(ctx).Where(u.Name.In("modi", "zhangqiang")).Find()
// SELECT * FROM users WHERE name IN ('modi','zhangqiang');

// LIKE
users, err := u.WithContext(ctx).Where(u.Name.Like("%modi%")).Find()
// SELECT * FROM users WHERE name LIKE '%modi%';

// AND
users, err := u.WithContext(ctx).Where(u.Name.Eq("modi"), u.Age.Gte(17)).Find()
// SELECT * FROM users WHERE name = 'modi' AND age >= 17;

// Time
users, err := u.WithContext(ctx).Where(u.Birthday.Gt(birthTime).Find()
// SELECT * FROM users WHERE birthday > '2000-01-01 00:00:00';

// BETWEEN
users, err := u.WithContext(ctx).Where(u.Birthday.Between(lastWeek, today)).Find()
// SELECT * FROM users WHERE birthday BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';
```

###### Not 

```go
u := query.Use(db).User

user, err := u.WithContext(ctx).Not(u.Name.Eq("modi")).First()
// SELECT * FROM users WHERE NOT name = "modi" ORDER BY id LIMIT 1;

// Not In
users, err := u.WithContext(ctx).Not(u.Name.In("modi", "zhangqiang")).Find()
// SELECT * FROM users WHERE name NOT IN ("modi", "zhangqiang");

// Not In slice of primary keys
user, err := u.WithContext(ctx).Not(u.ID.In(1,2,3)).First()
// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;
```

###### Or

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(u.Role.Eq("admin")).Or(u.Role.Eq("super_admin")).Find()
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
```

###### Group
组合where或or 构建复杂查询

```go
p := query.Use(db).Pizza

pizzas, err := p.WithContext(ctx).Where(
    p.WithContext(ctx).Where(p.Pizza.Eq("pepperoni")).
        Where(p.Where(p.Size.Eq("small")).Or(p.Size.Eq("medium"))),
).Or(
    p.WithContext(ctx).Where(p.Pizza.Eq("hawaiian")).Where(p.Size.Eq("xlarge")),
).Find()

// SELECT * FROM `pizzas` WHERE (pizza = "pepperoni" AND (size = "small" OR size = "medium")) OR (pizza = "hawaiian" AND size = "xlarge")
```

###### 指定字段查询

通过`Select` 可以选择你要查询的字段，否则就是查询所有字段。

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Select(u.Name, u.Age).Find()
// SELECT name, age FROM users;

u.WithContext(ctx).Select(u.Age.Avg()).Rows()
// SELECT Avg(age) FROM users;
```

###### 元组查询
例如多字段IN

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(u.Columns(u.ID, u.Name).In(field.Values([][]inferface{}{{1, "modi"}, {2, "zhangqiang"}}))).Find()
// SELECT * FROM `users` WHERE (`id`, `name`) IN ((1,'humodi'),(2,'tom'));
```

###### JSON 查询

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(gen.Cond(datatypes.JSONQuery("attributes").HasKey("role"))...).Find()
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`,'$.role') IS NOT NULL;
```

###### Order

指定查询的排序方式

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Order(u.Age.Desc(), u.Name).Find()
// SELECT * FROM users ORDER BY age DESC, name;

// Multiple orders
users, err := u.WithContext(ctx).Order(u.Age.Desc()).Order(u.Name).Find()
// SELECT * FROM users ORDER BY age DESC, name;
```

###### Limit & Offset

分页查询，`Limit`限制最大条数，`Offset`指定数据的其实位置。

```go
u := query.Use(db).User

urers, err := u.WithContext(ctx).Limit(3).Find()
// SELECT * FROM users LIMIT 3;

// Cancel limit condition with -1
users, err := u.WithContext(ctx).Limit(10).Limit(-1).Find()
// SELECT * FROM users;

users, err := u.WithContext(ctx).Offset(3).Find()
// SELECT * FROM users OFFSET 3;

users, err := u.WithContext(ctx).Limit(10).Offset(5).Find()
// SELECT * FROM users OFFSET 5 LIMIT 10;

// Cancel offset condition with -1
users, err := u.WithContext(ctx).Offset(10).Offset(-1).Find()
// SELECT * FROM users;
```

###### Group By & Having

```go
u := query.Use(db).User

type Result struct {
    Date  time.Time
    Total int
}

var result Result

err := u.WithContext(ctx).Select(u.Name, u.Age.Sum().As("total")).Where(u.Name.Like("%modi%")).Group(u.Name).Scan(&result)
// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "%modi%" GROUP BY `name`

err := u.WithContext(ctx).Select(u.Name, u.Age.Sum().As("total")).Group(u.Name).Having(u.Name.Eq("group")).Scan(&result)
// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"

rows, err := u.WithContext(ctx).Select(u.Birthday.As("date"), u.Age.Sum().As("total")).Group(u.Birthday).Rows()
for rows.Next() {
  ...
}

o := query.Use(db).Order

rows, err := o.WithContext(ctx).Select(o.CreateAt.Date().As("date"), o.Amount.Sum().As("total")).Group(o.CreateAt.Date()).Having(u.Amount.Sum().Gt(100)).Rows()
for rows.Next() {
  ...
}

var results []Result

o.WithContext(ctx).Select(o.CreateAt.Date().As("date"), o.WithContext(ctx).Amount.Sum().As("total")).Group(o.CreateAt.Date()).Having(u.Amount.Sum().Gt(100)).Scan(&results)
```

###### Distinct

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Distinct(u.Name, u.Age).Order(u.Name, u.Age.Desc()).Find()
```

`Distinct` works with `Pluck` and `Count` too

###### Joins

联表查询，`Join`是指`inner join`,还有`LeftJoin`和`RightJoin`

```go
u := query.Use(db).User
e := query.Use(db).Email
c := query.Use(db).CreditCard

type Result struct {
    Name  string
    Email string
}

var result Result

err := u.WithContext(ctx).Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Scan(&result)
// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

rows, err := u.WithContext(ctx).Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Rows()
for rows.Next() {
  ...
}

var results []Result

err := u.WithContext(ctx).Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Scan(&results)

// multiple joins with parameter
users := u.WithContext(ctx).Join(e, e.UserID.EqCol(u.id), e.Email.Eq("modi@example.org")).Join(c, c.UserID.EqCol(u.ID)).Where(c.Number.Eq("411111111111")).Find()
```

##### 子查询


```go
o := query.Use(db).Order
u := query.Use(db).User

orders, err := o.WithContext(ctx).Where(u.Columns(o.Amount).Gt(o.WithContext(ctx).Select(o.Amount.Avg())).Find()
// SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

subQuery := u.WithContext(ctx).Select(u.Age.Avg()).Where(u.Name.Like("name%"))
users, err := u.WithContext(ctx).Select(u.Age.Avg().As("avgage")).Group(u.Name).Having(u.Columns(u.Age.Avg()).Gt(subQuery).Find()
// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")
```

###### From 子查询

通过`Table`方法构建出的子查询，可以直接放到From语句中:

```go
u := query.Use(db).User
p := query.Use(db).Pet

users, err := gen.Table(u.WithContext(ctx).Select(u.Name, u.Age).As("u")).Where(u.Age.Eq(18)).Find()
// SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE `age` = 18

subQuery1 := u.WithContext(ctx).Select(u.Name)
subQuery2 := p.WithContext(ctx).Select(p.Name)
users, err := gen.Table(subQuery1.As("u"), subQuery2.As("p")).Find()
db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&User{})
// SELECT * FROM (SELECT `name` FROM `users`) as u, (SELECT `name` FROM `pets`) as p
```

###### 更新子查询

通过子查询更新表字段

```go
u := query.Use(db).User
c := query.Use(db).Company

u.WithContext(ctx).Update(u.CompanyName, c.Select(c.Name).Where(c.ID.EqCol(u.CompanyID)))
// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);

u.WithContext(ctx).Where(u.Name.Eq("modi")).Update(u.CompanyName, c.Select(c.Name).Where(c.ID.EqCol(u.CompanyID)))
```

###### 多字段更新子查询

针对mysql提供同时更新多个字段的子查询:

```go
u := query.Use(db).User
c := query.Use(db).Company

ua := u.As("u")
ca := u.As("c")

ua.WithContext(ctx).UpdateFrom(ca.WithContext(ctx).Select(c.ID, c.Address, c.Phone).Where(c.ID.Gt(100))).
Where(ua.CompanyID.EqCol(ca.ID)).
UpdateSimple(
  ua.Address.SetCol(ca.Address),
  ua.Phone.SetCol(ca.Phone),
)
// UPDATE `users` AS `u`,(
//   SELECT `company`.`id`,`company`.`address`,`company`.`phone` 
//   FROM `company` WHERE `company`.`id` > 100 AND `company`.`deleted_at` IS NULL
// ) AS `c` 
// SET `u`.`address`=`c`.`address`,`c`.`phone`=`c`.`phone`,`updated_at`='2021-11-11 11:11:11.111'
// WHERE `u`.`company_id` = `c`.`id`
```

##### 事务

多个操作需要在一个事务中完成的情况.

```go
q := query.Use(db)

q.Transaction(func(tx *query.Query) error {
  if _, err := tx.User.WithContext(ctx).Where(tx.User.ID.Eq(100)).Delete(); err != nil {
    return err
  }
  if _, err := tx.Article.WithContext(ctx).Create(&model.User{Name:"modi"}); err != nil {
    return err
  }
  return nil
})
```

###### 嵌套事务

GEN 支持潜逃事务，在一个大事务中嵌套子事务。

```go
q := query.Use(db)

q.Transaction(func(tx *query.Query) error {
  tx.User.WithContext(ctx).Create(&user1)

  tx.Transaction(func(tx2 *query.Query) error {
    tx2.User.WithContext(ctx).Create(&user2)
    return errors.New("rollback user2") // Rollback user2
  })

  tx.Transaction(func(tx2 *query.Query) error {
    tx2.User.WithContext(ctx).Create(&user3)
    return nil
  })

  return nil
})

// Commit user1, user3
```

###### 手动事务

```go
q := query.Use(db)

// begin a transaction
tx := q.Begin()

// do some database operations in the transaction (use 'tx' from this point, not 'db')
tx.User.WithContext(ctx).Create(...)

// ...

// rollback the transaction in case of error
tx.Rollback()

// Or commit the transaction
tx.Commit()
```

For example:

```go
q := query.Use(db)

func doSomething(ctx context.Context, users ...*model.User) (err error) {
    tx := q.Begin()
    defer func() {
        if recover() != nil || err != nil {
            _ = tx.Rollback()
        }
    }()

    err = tx.User.WithContext(ctx).Create(users...)
    if err != nil {
        return
    }
    return tx.Commit()
}
```

###### 保存点

`SavePoint`, `RollbackTo` 可以保存或者回滚事务点:

```go
tx := q.Begin()
txCtx = tx.WithContext(ctx)

txCtx.User.Create(&user1)

tx.SavePoint("sp1")
txCtx.Create(&user2)
tx.RollbackTo("sp1") // Rollback user2

tx.Commit() // Commit user1
```

##### Advanced Query

###### 迭代

GEN支持通过Row迭代取值

```go
u := query.Use(db).User
do := u.WithContext(ctx)
rows, err := do.Where(u.Name.Eq("modi")).Rows()
defer rows.Close()

for rows.Next() {
    var user User
    // ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
    do.ScanRows(rows, &user)

    // do something
}
```

###### 批量查询


```go
u := query.Use(db).User

// batch size 100
err := u.WithContext(ctx).Where(u.ID.Gt(9)).FindInBatches(&results, 100, func(tx gen.Dao, batch int) error {
    for _, result := range results {
      // batch processing found records
    }
  
    // build a new `u` to use it's api
    // queryUsery := query.NewUser(tx.UnderlyingDB())

    tx.Save(&results)

    batch // Batch 1, 2, 3

    // returns error will stop future batches
    return nil
})
```

###### Pluck

从数据库中查询单个列并扫描成一个切片或者基础类型

```go
u := query.Use(db).User

var ages []int64
u.WithContext(ctx).Pluck(u.Age, &ages)

var names []string
u.WithContext(ctx).Pluck(u.Name, &names)

// Distinct Pluck
u.WithContext(ctx).Distinct().Pluck(u.Name, &names)
// SELECT DISTINCT `name` FROM `users`

// Requesting more than one column, use `Scan` or `Find` like this:
db.WithContext(ctx).Select(u.Name, u.Age).Scan(&users)
users, err := db.Select(u.Name, u.Age).Find()
```

###### Scopes

可以声明一些常用的或者公用的条件方法，然后通过`Scopes` 查询

```go
o := query.Use(db).Order

func AmountGreaterThan1000(tx gen.Dao) gen.Dao {
    return tx.Where(o.Amount.Gt(1000))
}

func PaidWithCreditCard(tx gen.Dao) gen.Dao {
    return tx.Where(o.PayModeSign.Eq("C"))
}

func PaidWithCod(tx gen.Dao) gen.Dao {
    return tx.Where(o.PayModeSign.Eq("C"))
}

func OrderStatus(status []string) func (tx gen.Dao) gen.Dao {
    return func (tx gen.Dao) gen.Dao {
      return tx.Where(o.Status.In(status...))
    }
}

orders, err := o.WithContext(ctx).Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find()
// Find all credit card orders and amount greater than 1000

orders, err := o.WithContext(ctx).Scopes(AmountGreaterThan1000, PaidWithCod).Find()
// Find all COD orders and amount greater than 1000

orders, err := o.WithContext(ctx).Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find()
// Find all paid, shipped orders that amount greater than 1000
```

###### Count


```go
u := query.Use(db).User

count, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).Or(u.Name.Eq("zhangqiang")).Count()
// SELECT count(1) FROM users WHERE name = 'modi' OR name = 'zhangqiang'

count, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).Count()
// SELECT count(1) FROM users WHERE name = 'modi'; (count)

// Count with Distinct
u.WithContext(ctx).Distinct(u.Name).Count()
// SELECT COUNT(DISTINCT(`name`)) FROM `users`
```

###### FirstOrInit

获取匹配的第一条数据或用给定条件初始化一个实例

```go
u := query.Use(db).User

// User not found, initialize it with give conditions
user, err := u.WithContext(ctx).Where(u.Name.Eq("non_existing")).FirstOrInit()
// user -> User{Name: "non_existing"}

// Found user with `name` = `modi`
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).FirstOrInit()
// user -> User{ID: 1, Name: "modi", Age: 17}
```

如果希望初始化的实例包含一些非查询条件的属性，则可以通过`Attrs`指定

```go
u := query.Use(db).User

// User not found, initialize it with give conditions and Attrs
user, err := u.WithContext(ctx).Where(u.Name.Eq("non_existing")).Attrs(u.Age.Value(20)).FirstOrInit()
// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
// user -> User{Name: "non_existing", Age: 20}

// User not found, initialize it with give conditions and Attrs
user, err := u.WithContext(ctx).Where(u.Name.Eq("non_existing")).Attrs(u.Age.Value(20)).FirstOrInit()
// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
// user -> User{Name: "non_existing", Age: 20}

// Found user with `name` = `modi`, attributes will be ignored
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).Attrs(u.Age.Value(20)).FirstOrInit()
// SELECT * FROM USERS WHERE name = modi' ORDER BY id LIMIT 1;
// user -> User{ID: 1, Name: "modi", Age: 17}
```

`Assign` 则是无论有没有找到记录，都用指定的属性进行覆盖已有的属性

```go
// User not found, initialize it with give conditions and Assign attributes
user, err := u.WithContext(ctx).Where(u.Name.Eq("non_existing")).Assign(u.Age.Value(20)).FirstOrInit()
// user -> User{Name: "non_existing", Age: 20}

// Found user with `name` = `modi`, update it with Assign attributes
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).Assign(u.Age.Value(20)).FirstOrInit()
// SELECT * FROM USERS WHERE name = modi' ORDER BY id LIMIT 1;
// user -> User{ID: 111, Name: "modi", Age: 20}
```

###### FirstOrCreate

获取第一条匹配的记录或在给定条件下创建一条新记录

```go
u := query.Use(db).User

// User not found, create a new record with give conditions
user, err := u.WithContext(ctx).Where(u.Name.Eq("non_existing")).FirstOrCreate()
// INSERT INTO "users" (name) VALUES ("non_existing");
// user -> User{ID: 112, Name: "non_existing"}

// Found user with `name` = `modi`
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).FirstOrCreate()
// user -> User{ID: 111, Name: "modi", "Age": 18}
```

如果希望创建的实例包含一些非查询条件的属性，则可以通过`Attrs`指定

```go
u := query.Use(db).User

// User not found, create it with give conditions and Attrs
user, err := u.WithContext(ctx).Where(u.Name.Eq("non_existing")).Attrs(u.Age.Value(20)).FirstOrCreate()
// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
// user -> User{ID: 112, Name: "non_existing", Age: 20}

// Found user with `name` = `modi`, attributes will be ignored
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).Attrs(u.Age.Value(20)).FirstOrCreate()
// SELECT * FROM users WHERE name = 'modi' ORDER BY id LIMIT 1;
// user -> User{ID: 111, Name: "modi", Age: 18}
```

`Assign` 则是无论有没有找到记录，都用指定的属性进行覆盖并且入库

```go
u := query.Use(db).User

// User not found, initialize it with give conditions and Assign attributes
user, err := u.WithContext(ctx).Where(u.Name.Eq("non_existing")).Assign(u.Age.Value(20)).FirstOrCreate()
// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
// INSERT INTO "users" (name, age) VALUES ("non_existing", 20);
// user -> User{ID: 112, Name: "non_existing", Age: 20}

// Found user with `name` = `modi`, update it with Assign attributes
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).Assign(u.Age.Value(20)).FirstOrCreate(&user)
// SELECT * FROM users WHERE name = 'modi' ORDER BY id LIMIT 1;
// UPDATE users SET age=20 WHERE id = 111;
// user -> User{ID: 111, Name: "modi", Age: 20}
```

#### 关联

GEN将像GORM一样自动保存关联（(BelongsTo/HasOne/HasMany/Many2Many) 。

##### Relation

There are 4 kind of relationship.

```go
const (
    HasOne    RelationshipType = RelationshipType(schema.HasOne)    // HasOneRel has one relationship
    HasMany   RelationshipType = RelationshipType(schema.HasMany)   // HasManyRel has many relationships
    BelongsTo RelationshipType = RelationshipType(schema.BelongsTo) // BelongsToRel belongs to relationship
    Many2Many RelationshipType = RelationshipType(schema.Many2Many) // Many2ManyRel many to many relationship
)
```

###### 关联已经存在的Model

```go
package model

// exist model
type Customer struct {
    gorm.Model
    CreditCards []CreditCard `gorm:"foreignKey:CustomerRefer"`
}

type CreditCard struct {
    gorm.Model
    Number        string
    CustomerRefer uint
}
```

GEN 会检查解析这些关联关系:

```go
// specify model
g.ApplyBasic(model.Customer{}, model.CreditCard{})

// assoications will be detected and converted to code 
package query

type customer struct {
    ...
    CreditCards customerHasManyCreditCards
}

type creditCard struct{
    ...
}
```

###### 和数据库表关联

必须使用 `gen.FieldRelate`声明

```go
card := g.GenerateModel("credit_cards")
customer := g.GenerateModel("customers", gen.FieldRelate(field.HasMany, "CreditCards", b, 
    &field.RelateConfig{
        // RelateSlice: true,
        GORMTag: "foreignKey:CustomerRefer",
    }),
)

g.ApplyBasic(card, custormer)
```

GEN 会生成申明的关联属性:

```go
// customers
type Customer struct {
    ID          int64          `gorm:"column:id;type:bigint(20) unsigned;primaryKey" json:"id"`
    CreatedAt   time.Time      `gorm:"column:created_at;type:datetime(3)" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"column:updated_at;type:datetime(3)" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3)" json:"deleted_at"`
    CreditCards []CreditCard   `gorm:"foreignKey:CustomerRefer" json:"credit_cards"`
}


// credit_cards
type CreditCard struct {
    ID            int64          `gorm:"column:id;type:bigint(20) unsigned;primaryKey" json:"id"`
    CreatedAt     time.Time      `gorm:"column:created_at;type:datetime(3)" json:"created_at"`
    UpdatedAt     time.Time      `gorm:"column:updated_at;type:datetime(3)" json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3)" json:"deleted_at"`
    CustomerRefer int64          `gorm:"column:customer_refer;type:bigint(20) unsigned" json:"customer_refer"`
}
```

如果是已经存在的关联model, 则可以用`gen.FieldRelateModel` 声明.

```go
customer := g.GenerateModel("customers", gen.FieldRelateModel(field.HasMany, "CreditCards", model.CreditCard{}, 
    &field.RelateConfig{
        // RelateSlice: true,
        GORMTag: "foreignKey:CustomerRefer",
    }),
)

g.ApplyBasic(custormer)
```

###### 关联配置

```go
type RelateConfig struct {
    // specify field's type
    RelatePointer      bool // ex: CreditCard  *CreditCard
    RelateSlice        bool // ex: CreditCards []CreditCard
    RelateSlicePointer bool // ex: CreditCards []*CreditCard

    JSONTag      string // related field's JSON tag
    GORMTag      string // related field's GORM tag
    NewTag       string // related field's new tag
    OverwriteTag string // related field's tag
}
```

##### 操作

###### 跳过自动创建关联

```go
user := model.User{
  Name:            "modi",
  BillingAddress:  Address{Address1: "Billing Address - Address 1"},
  ShippingAddress: Address{Address1: "Shipping Address - Address 1"},
  Emails:          []Email{
    {Email: "modi@example.com"},
    {Email: "modi-2@example.com"},
  },
  Languages:       []Language{
    {Name: "ZH"},
    {Name: "EN"},
  },
}

u := query.Use(db).User

u.WithContext(ctx).Select(u.Name).Create(&user)
// INSERT INTO "users" (name) VALUES ("jinzhu", 1, 2);

u.WithContext(ctx).Omit(u.BillingAddress.Field()).Create(&user)
// Skip create BillingAddress when creating a user

u.WithContext(ctx).Omit(u.BillingAddress.Field("Address1")).Create(&user)
// Skip create BillingAddress.Address1 when creating a user

u.WithContext(ctx).Omit(field.AssociationFields).Create(&user)
// Skip all associations when creating a user
```

Method `Field` will join a serious field name with ''.", for example: `u.BillingAddress.Field("Address1", "Street")` equals to `BillingAddress.Address1.Street`

###### 查询关联

```go
u := query.Use(db).User

languages, err = u.Languages.Model(&user).Find()
```

查询指定条件的关联

```go
q := query.Use(db)
u := q.User

languages, err = u.Languages.Where(q.Language.Name.In([]string{"ZH","EN"})).Model(&user).Find()
```

###### 添加关联


```go
u := query.Use(db).User

u.Languages.Model(&user).Append(&languageZH, &languageEN)

u.Languages.Model(&user).Append(&Language{Name: "DE"})

u.CreditCards.Model(&user).Append(&CreditCard{Number: "411111111111"})
```

###### 替换关联

```go
u.Languages.Model(&user).Replace(&languageZH, &languageEN)
```

###### 删除关联

删除存在的关联，不会删除数据

```go
u := query.Use(db).User

u.Languages.Model(&user).Delete(&languageZH, &languageEN)

u.Languages.Model(&user).Delete([]*Language{&languageZH, &languageEN}...)
```

###### 清楚关联

清楚所有的关联，不会删除数据

```go
u.Languages.Model(&user).Clear()
```

###### 统计关联

```go
u.Languages.Model(&user).Count()
```

###### 删除指定关联

删除制定条件数据并删除关联数据:

```go
u := query.Use(db).User

// delete user's account when deleting user
u.Select(u.Account).Delete(&user)

// delete user's Orders, CreditCards relations when deleting user
db.Select(u.Orders.Field(), u.CreditCards.Field()).Delete(&user)

// delete user's has one/many/many2many relations when deleting user
db.Select(field.AssociationsFields).Delete(&user)
```

##### 预加载

###### Preload

GEN 支持通过 `Preload`加载关联数据:

```go
type User struct {
  gorm.Model
  Username string
  Orders   []Order
}

type Order struct {
  gorm.Model
  UserID uint
  Price  float64
}

q := query.Use(db)
u := q.User
o := q.Order

// Preload Orders when find users
users, err := u.WithContext(ctx).Preload(u.Orders).Find()
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4);

users, err := u.WithContext(ctx).Preload(u.Orders).Preload(u.Profile).Preload(u.Role).Find()
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4); // has many
// SELECT * FROM profiles WHERE user_id IN (1,2,3,4); // has one
// SELECT * FROM roles WHERE id IN (4,5,6); // belongs to
```

###### Preload All

`clause.Associations` 通过`Preload` 预加载所有的关联数据:

```go
type User struct {
  gorm.Model
  Name       string
  CompanyID  uint
  Company    Company
  Role       Role
  Orders     []Order
}

users, err := u.WithContext(ctx).Preload(field.Associations).Find()
```

`clause.Associations` 不会加载嵌套关联, 潜逃关联家在可以用 [Nested Preloading](#nested_preloading) e.g:

```go
users, err := u.WithContext(ctx).Preload(u.Orders.OrderItems.Product).Find()
```

###### 根据条件预加载

```go
q := query.Use(db)
u := q.User
o := q.Order

// Preload Orders with conditions
users, err := u.WithContext(ctx).Preload(u.Orders.On(o.State.NotIn("cancelled")).Find()
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2,3,4) AND state NOT IN ('cancelled');

users, err := u.WithContext(ctx).Where(u.State.Eq("active")).Preload(u.Orders.On(o.State.NotIn("cancelled")).Find()
// SELECT * FROM users WHERE state = 'active';
// SELECT * FROM orders WHERE user_id IN (1,2) AND state NOT IN ('cancelled');

users, err := u.WithContext(ctx).Preload(u.Orders.Order(o.ID.Desc(), o.CreateTime).Find()
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2) Order By id DESC, create_time;

users, err := u.WithContext(ctx).Preload(u.Orders.On(o.State.Eq("on")).Order(o.ID.Desc()).Find()
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2) AND state = "on" Order By id DESC;

users, err := u.WithContext(ctx).Preload(u.Orders.Clauses(hints.UseIndex("idx_order_id"))).Find()
// SELECT * FROM users;
// SELECT * FROM orders WHERE user_id IN (1,2) USE INDEX (`idx_order_id`);
```

###### 潜逃预加载


```go
db.Preload(u.Orders.OrderItems.Product).Preload(u.CreditCard).Find(&users)

// Customize Preload conditions for `Orders`
// And GEN won't preload unmatched order's OrderItems then
db.Preload(u.Orders.On(o.State.Eq("paid"))).Preload(u.Orders.OrderItems).Find(&users)
```

#### 更新

##### 更新单字段

`Update`方法更新单个字段。需要注意的是必须指定更新条件否则会报错`ErrMissingWhereClause`:

```go
u := query.Use(db).User

// Update with conditions
u.WithContext(ctx).Where(u.Activate.Is(true)).Update(u.Name, "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

// Update with conditions
u.WithContext(ctx).Where(u.Activate.Is(true)).Update(u.Age, u.Age.Add(1))
// or
u.WithContext(ctx).Where(u.Activate.Is(true)).UpdateSimple(u.Age.Add(1))
// UPDATE users SET age=age+1, updated_at='2013-11-17 21:34:10' WHERE active=true;

u.WithContext(ctx).Where(u.Activate.Is(true)).UpdateSimple(u.Age.Zero())
// UPDATE users SET age=0, updated_at='2013-11-17 21:34:10' WHERE active=true;
```

##### 更新多字段

`Updates` 支持 `struct` 和 `map[string]interface{}`类型，更新多个字段，但是会忽略其中的零值属性

```go
u := query.Use(db).User

// Update attributes with `map`
u.WithContext(ctx).Where(u.ID.Eq(111)).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

// Update attributes with `struct`
u.WithContext(ctx).Where(u.ID.Eq(111)).Updates(model.User{Name: "hello", Age: 18, Active: false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

// Update with expression
u.WithContext(ctx).Where(u.ID.Eq(111)).UpdateSimple(u.Age.Add(1), u.Number.Add(1))
// UPDATE users SET age=age+1,number=number+1, updated_at='2013-11-17 21:34:10' WHERE id=111;

u.WithContext(ctx).Where(u.Activate.Is(true)).UpdateSimple(u.Age.Value(17), u.Number.Zero(), u.Birthday.Null())
// UPDATE users SET age=17, number=0, birthday=NULL, updated_at='2013-11-17 21:34:10' WHERE active=true;
```

> **NOTE** When update with struct, GEN will only update non-zero fields, you might want to use `map` to update attributes or use `Select` to specify fields to update

##### 选择更新的字段

通过 `Select`, `Omit`选择需要更新的字段或者需要忽略更新的字段

```go
u := query.Use(db).User

// Select with Map
// User's ID is `111`:
u.WithContext(ctx).Select(u.Name).Where(u.ID.Eq(111)).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello' WHERE id=111;

u.WithContext(ctx).Omit(u.Name).Where(u.ID.Eq(111)).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;

result, err := u.WithContext(ctx).Where(u.ID.Eq(111)).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})

result.RowsAffected // affect rows number
err                 // error
```

#### 删除

##### 删除记录

```go
e := query.Use(db).Email

// Email's ID is `10`
e.WithContext(ctx).Where(e.ID.Eq(10)).Delete()
// DELETE from emails where id = 10;

// Delete with additional conditions
e.WithContext(ctx).Where(e.ID.Eq(10), e.Name.Eq("modi")).Delete()
// DELETE from emails where id = 10 AND name = "modi";

result, err := e.WithContext(ctx).Where(e.ID.Eq(10), e.Name.Eq("modi")).Delete()

result.RowsAffected // affect rows number
err                 // error
```

##### 根据主键删除


```go
u.WithContext(ctx).Where(u.ID.In(1,2,3)).Delete()
// DELETE FROM users WHERE id IN (1,2,3);
```

##### 批量删除

没有指定主键会删除松油匹配的数据

```go
e := query.Use(db).Email

e.WithContext(ctx).Where(e.Name.Like("%modi%")).Delete()
// DELETE from emails where email LIKE "%modi%";
```

##### 软删除

如果你的model中有`gorm.DeletedAt` 字段，则会自动执行软删除。也就是不会删除数据，只是把该字段的指设置为当前时间。

```go
// Batch Delete
u.WithContext(ctx).Where(u.Age.Eq(20)).Delete()
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// Soft deleted records will be ignored when querying
users, err := u.WithContext(ctx).Where(u.Age.Eq(20)).Find()
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
```

If you don’t want to include `gorm.Model`, you can enable the soft delete feature like:

```go
type User struct {
    ID      int
    Deleted gorm.DeletedAt
    Name    string
}
```

##### 查询包含软删除的记录

可以通过 `Unscoped`实现

```go
users, err := db.WithContext(ctx).Unscoped().Where(u.Age.Eq(20)).Find()
// SELECT * FROM users WHERE age = 20;
```

##### 永久删除

通过 `Unscoped`可以直接删除数据，而不是标记删除

```go
o.WithContext(ctx).Unscoped().Where(o.ID.Eq(10)).Delete()
// DELETE FROM orders WHERE id=10;
```

### DIY 方法

#### 接口定义

自定义方法，需要通过接口定义。在方上通过注释的方式描述具体的查询逻辑，复杂的可以直接用`sql()`，简单的直接用`where()`。
如果想要写一些原始注释，可以先写注释然后换行，在写`sql`或者`where`。

```go
type Method interface {
    // where("name=@name and age=@age")
    SimpleFindByNameAndAge(name string, age int) (gen.T, error)

    // FindUserToMap query by id and return id->instance
    // 
    // sql(select * from users where id=@id)
    FindUserToMap(id int) (gen.M, error)
    
    // InsertValue insert into users (name,age) values (@name,@age)
    InsertValue(age int, name string) error
}
```

返回值可以是 `gen.T`/`gen.M`/`gen.RowsAffected`。 `gen.T` 表示单个model, `[]gen.T`表示的是slice model，`gen.M`表示的是`map[string]interface{}`,当然也可以是其他类型Gen不会转换。除了返回一个值外，还可以返回一个err。

##### 语法

###### 占位符

- `gen.T` represents specified `struct` or `table`
- `gen.M` represents `map[string]interface`
- `gen.RowsAffected` represents SQL executed `rowsAffected` (type:int64)
- `@@table`  represents table's name (if method's parameter doesn't contains variable `table`, GEN will generate `table` from model struct)
- `@@<columnName>` represents column's name or table's name
- `@<name>` represents normal query variable

###### 模板

逻辑操作必须包裹在`{{}}`中，如`{{if}}`,结束语句必须是 `{{end}}`, 所有的语句都可以嵌套。

- `if`/`else if`/`else` the condition accept a bool parameter or operation expression which conforms to Golang syntax.
- `where` The `where` clause will be inserted only if the child elements return something. The key word  `and` or `or`  in front of clause will be removed. And `and` will be added automatically when there is no junction keyword between query condition clause.
- `Set` The  `set` clause will be inserted only if the child elements return something. The `,` in front of columns array will be removed.And `,` will be added automatically when there is no junction keyword between query coulmns.
- `...` Coming soon

###### `If` clause

```sql
{{if cond1}}
    // do something here
{{else if cond2}}
    // do something here
{{else}}
    // do something here
{{end}}
```

Use case in raw SQL:

```go
// select * from users where {{if name !=""}} name=@name{{end}}
methond(name string) (gen.T,error) 
```

Use case in raw SQL template:

```sql
select * from @@table where
{{if age>60}}
    status="older"
{{else if age>30}}
    status="middle-ager"
{{else if age>18}}
    status="younger"
{{else}}
    {{if sex=="male"}}
        status="boys"
    {{else}}
        status="girls"
    {{end}}
{{end}}
```

###### `Where` clause

```sql
{{where}}
    // do something here
{{end}}
```

Use case in raw SQL

```go
// select * from {{where}}id=@id{{end}}
methond(id int) error
```

Use case in raw SQL template

```sql
select * from @@table 
{{where}}
    {{if cond}}id=@id {{end}}
    {{if name != ""}}@@key=@value{{end}}
{{end}}
```

###### `Set` clause

```sql
{{set}}
    // sepecify update expression here
{{end}}
```

Use case in raw SQL

```go
// update users {{set}}name=@name{{end}}
methond() error
```

Use case in raw SQL template

```sql
update @@table 
{{set}}
    {{if name!=""}} name=@name {{end}}
    {{if age>0}} age=@age {{end}}
{{end}}
where id=@id
```

##### Method interface example

```go
type Method interface {
    // Where("name=@name and age=@age")
    SimpleFindByNameAndAge(name string, age int) (gen.T, error)
    
    // select * from users where id=@id
    FindUserToMap(id int) (gen.M, error)
    
    // sql(insert into @@table (name,age) values (@name,@age) )
    InsertValue(age int, name string) error
    
    // select name from @@table where id=@id
    FindNameByID(id int) string
    
    // select * from @@table
    //  {{where}}
    //      id>0
    //      {{if cond}}id=@id {{end}}
    //      {{if key!="" && value != ""}} or @@key=@value{{end}}
    //  {{end}}
    FindByIDOrCustom(cond bool, id int, key, value string) ([]gen.T, error)
    
    // update @@table
    //  {{set}}
    //      update_time=now()
    //      {{if name != ""}}
    //          name=@name
    //      {{end}}
    //  {{end}}
    //  {{where}}
    //      id=@id
    //  {{end}}
    UpdateName(name string, id int) (gen.RowsAffected,error)
}
```

#### 智能选择字段

GEN 查询的时候会自动选择你的model定义的字段

```go
type User struct {
  ID     uint
  Name   string
  Age    int
  Gender string
  // hundreds of fields
}

type APIUser struct {
  ID   uint
  Name string
}

type Method interface{
    // select * from user
    FindSome() ([]APIUser, error)
}

apiusers, err := u.WithContext(ctx).Limit(10).FindSome()
// SELECT `id`, `name` FROM `users` LIMIT 10
```

### 高级教程

#### Hints

Hints可以用来优化查询计划，比如指定索引后者强制索引等。

```go
import "gorm.io/hints"

u := query.Use(db).User

users, err := u.WithContext(ctx).Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find()
// SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`
```


```go
import "gorm.io/hints"

u := query.Use(db).User

users, err := u.WithContext(ctx).Clauses(hints.UseIndex("idx_user_name")).Find()
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

users, err := u.WithContext(ctx).Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find()
// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)"
```

## 命令工具

安装gen命令行工具:

```bash
go install gorm.io/gen/tools/gentool@latest
```

使用参数:

```bash
$ gentool -h
Usage of gentool:
  -db string
      input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html] (default "mysql")
  -dsn string
      consult[https://gorm.io/docs/connecting_to_the_database.html]
  -fieldNullable
      generate with pointer when field is nullable
  -fieldWithIndexTag
      generate field with gorm index tag
  -fieldWithTypeTag
      generate field with gorm column type tag
  -modelPkgName string
      generated model code's package name
  -outFile string
      query code file name, default: gen.go
  -outPath string
      specify a directory for output (default "./dao/query")
  -tables string
      enter the required data table or leave it blank
  -withUnitTest
      generate unit test for query code
```

示例:

``` bash
gentool -dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" -tables "orders,doctor"
```

## Maintainers

[@riverchu](https://github.com/riverchu) [@idersec](https://github.com/idersec) [@qqxhb](https://github.com/qqxhb) [@dino-ma](https://github.com/dino-ma)

[@jinzhu](https://github.com/jinzhu)

## Contributing

You can help to deliver a better GORM/GEN

## License

Released under the [MIT License](https://github.com/go-gorm/gen/blob/master/License)
