# <span id="gormgen">GORM/GEN</span>

[![GoVersion](https://img.shields.io/github/go-mod/go-version/go-gorm/gen)](https://github.com/go-gorm/gen/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/go-gorm/gen)](https://github.com/go-gorm/gen/releases)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/gorm.io/gen?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-gorm/gen)](https://goreportcard.com/report/github.com/go-gorm/gen)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![OpenIssue](https://img.shields.io/github/issues/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aopen+is%3Aissue)
[![ClosedIssue](https://img.shields.io/github/issues-closed/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aissue+is%3Aclosed)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/go-gorm/gen)](https://www.tickgit.com/browse?repo=github.com/go-gorm/gen)

基于 [GORM](https://github.com/go-gorm/gorm), 更安全更友好的ORM工具。

## <span id="multilingual-readme">更多语言版本的 README</span>

[English Version](./README.md) | [中文版本](./README.zh-CN.md)

## <span id="overview">概览</span>

- 自动生成 CRUD 和 DIY 方法
- 自动根据表结构生成模型（model）代码
- 事务、嵌套事务、保存点、回滚事务点
- 完全兼容 GORM
- 更安全、更友好
- 多种生成代码模式

## <span id="contents">目录</span>

- [GORM/GEN](#gormgen)
  - [更多语言版本的 README](#multilingual-readme)
  - [概览](#overview)
  - [目录](#contents)
  - [安装](#installation)
  - [快速开始](#quick-start)
    - [项目路径](#project-directory)
  - [API 示例](#api-examples)
    - [生成](#generate)
      - [生成模式](#generate-mode)
      - [模型生成](#generate-model)
      - [类型映射](#data-mapping)
    - [字段表达式](#field-expression)
      - [创建字段](#create-field)
    - [CRUD 接口](#crud-api)
      - [创建](#create)
        - [创建记录](#create-record)
        - [选择字段创建](#create-record-with-selected-fields)
        - [批量创建](#batch-insert)
      - [查询](#query)
        - [单条数据](#retrieving-a-single-object)
        - [根据主键查询数据](#retrieving-objects-with-primary-key)
        - [查询所有数据](#retrieving-all-objects)
        - [条件查询](#conditions)
          - [字符串和基础查询套件](#string-conditions)
          - [内联条件查询](#inline-condition)
          - [取反（not）查询](#not-conditions)
          - [或（or）查询](#or-conditions)
          - [组合查询](#group-conditions)
          - [指定字段查询](#selecting-specific-fields)
          - [元组查询](#tuple-query)
          - [JSON 查询](#json-query)
          - [Order 排序](#order)
          - [分页查询(Limit & Offset)](#limit--offset)
          - [分组查询（Group By & Having）](#group-by--having)
          - [去重（Distinct）](#distinct)
          - [联表查询（Joins）](#joins)
        - [子查询](#subquery)
          - [From 子查询](#from-subquery)
          - [从子查询更新](#update-from-subquery)
          - [从子查询更新多个字段](#update-multiple-columns-from-subquery)
        - [事务](#transaction)
          - [嵌套事务](#nested-transactions)
          - [手动事务](#transactions-by-manual)
          - [保存点/回滚](#savepointrollbackto)
        - [高级查询](#advanced-query)
          - [迭代](#iteration)
          - [批量查询](#findinbatches)
          - [Pluck 方法](#pluck)
          - [Scopes 查询](#scopes)
          - [Count 计数](#count)
          - [首条匹配或指定查询实例初始化条件（FirstOrInit）](#firstorinit)
          - [首条匹配或指定创建实例初始化条件（FirstOrCreate）](#firstorcreate)
      - [关联关系（Association）](#association)
        - [关联](#relation)
          - [关联已存在的模型](#relate-to-exist-model)
          - [和数据库表关联](#relate-to-table-in-database)
          - [关联配置](#relate-config)
        - [操作](#operation)
          - [跳过自动创建关联](#skip-auto-createupdate)
          - [查询关联](#find-associations)
          - [添加关联](#append-associations)
          - [替换关联](#replace-associations)
          - [删除关联](#delete-associations)
          - [清除关联](#clear-associations)
          - [统计关联](#count-associations)
          - [删除指定关联](#delete-with-select)
        - [预加载](#preloading)
          - [预加载（Preload）](#preload)
          - [预加载全部数据（Preload All）](#preload-all)
          - [预加载指定列](#preload-with-select)
          - [根据条件预加载](#preload-with-conditions)
          - [嵌套预加载](#nested-preloading)
      - [更新](#update)
        - [更新单个字段](#update-single-column)
        - [更新多个字段](#updates-multiple-columns)
        - [更新指定字段](#update-selected-fields)
      - [删除](#delete)
        - [删除记录](#delete-record)
        - [根据主键删除](#delete-with-primary-key)
        - [批量删除](#batch-delete)
        - [软删除](#soft-delete)
        - [查询包含软删除的记录](#find-soft-deleted-records)
        - [永久删除](#delete-permanently)
    - [自定义方法](#diy-method)
      - [接口定义](#method-interface)
        - [模板语法](#syntax-of-template)
          - [占位符](#placeholder)
          - [模板](#template)
          - [`If` 子句](#if-clause)
          - [`Where` 子句](#where-clause)
          - [`Set` 子句](#set-clause)
          - [`For` 子句](#for-clause)
        - [方法接口示例](#method-interface-example)
      - [单元测试](#unit-test)
      - [智能选择字段](#smart-select-fields)
    - [高级教程](#advanced-topics)
      - [查询优化提示（Hints）](#hints)
  - [二进制命令行工具安装](#binary)
  - [维护者](#maintainers)
  - [如何参与贡献](#contributing)
  - [开源许可协议](#license)

## <span id="installation">安装</span>

安装 GEN 前，需要安装好 Go 并配置你的 Go 工作区。

1. 安装完 Go（要求 1.14 及以上版本）后，可以使用以下 Go 命令安装 Gen。

```bash
go get -u gorm.io/gen
```

2. 在工程中导入引用 Gen:

```go
import "gorm.io/gen"
```

## <span id="quick-start">快速开始</span>

**注意**：此处所有教程都是在 `WithContext` 模式下写的. 如果你使用的是 `WithoutContext` 模式,则可以删除所有的 `WithContext(ctx)` 代码，这样看起来会更简洁。

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
        //if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
        /* FieldCoverable: true,*/
        // if you want generate field with unsigned integer type, set FieldSignable true
        /* FieldSignable: true,*/
        //if you want to generate index tags from database, set FieldWithIndexTag true
        /* FieldWithIndexTag: true,*/
        //if you want to generate type tags from database, set FieldWithTypeTag true
        /* FieldWithTypeTag: true,*/
        //if you need unit tests for query code, set WithUnitTest true
        /* WithUnitTest: true, */
    })
  
    // reuse the database connection in Project or create a connection here
    // if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
    // db, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
    g.UseDB(db)
  
    // apply basic crud api on structs or table models which is specified by table name with function
    // GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
    // 想对已有的model生成crud等基础方法可以直接指定model struct ，例如model.User{}
    // 如果是想直接生成表的model和crud方法，则可以指定表的名称，例如g.GenerateModel("company")
    // 想自定义某个表生成特性，比如struct的名称/字段类型/tag等，可以指定opt，例如g.GenerateModel("company",gen.FieldIgnore("address")), g.GenerateModelAs("people", "Person", gen.FieldIgnore("address"))
    g.ApplyBasic(model.User{}, g.GenerateModel("company"), g.GenerateModelAs("people", "Person", gen.FieldIgnore("address")))
    
    // apply diy interfaces on structs or table models
    // 如果想给某些表或者model生成自定义方法，可以用ApplyInterface，第一个参数是方法接口，可以参考DIY部分文档定义
    g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))

    // execute the action of code generation
    g.Execute()
}
```

生成 Model：

- `gen.WithoutContext` 可以生成没有 `WithContext` 约束的代码
- `gen.WithDefaultQuery` 使用默认全局变量 `Q` 作为单例生成代码

### <span id="project-directory">项目路径</span>

最佳实践项目模板：

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
|       └── user.gen_test.go # generated unit test
├── biz
│   └── query.go # call function in dal/gorm_generated.go and query databases
├── config
│   └── config.go # DSN for database server
├── generate.sh # a shell to execute cmd/generate
├── go.mod
├── go.sum
└── main.go
```

## <span id="api-examples">API 示例</span>

### <span id="generate">生成</span>

#### <span id="generate-mode">生成模式</span>

```go
 g := gen.NewGenerator(gen.Config{
        ...
        Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface,
        ...
 })
```

- `WithDefaultQuery` 生成默认查询结构体(作为全局变量使用)
- `WithoutContext` 生成没有context调用限制的代码供查询
- `WithQueryInterface` 生成interface形式的查询代码(可导出)

#### <span id="generate-model">模型生成</span>

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

**字段生成选项**

```go
FieldNew           // create new field
FieldIgnore        // ignore field
FieldIgnoreReg     // ignore field (match with regexp)
FieldRename        // rename field in struct
FieldComment       // specify field comment in generated struct
FieldType          // specify field type
FieldTypeReg       // specify field type (match with regexp)
FieldTag           // specify gorm and json tag
FieldJSONTag       // specify json tag
FieldJSONTagWithNS // specify new tag with name strategy
FieldGORMTag       // specify gorm tag
FieldNewTag        // append new tag
FieldNewTagWithNS  // specify new tag with name strategy
FieldTrimPrefix    // trim column prefix
FieldTrimSuffix    // trim column suffix
FieldAddPrefix     // add prefix to struct field's name
FieldAddSuffix     // add suffix to struct field's name
FieldRelate        // specify relationship with other tables
FieldRelateModel   // specify relationship with exist models
```

**生成结构体绑定自定义方法**
```Go
type User struct{
	ID int32
}
func (u *User)IsEmpty()bool{
    if u == nil {
    return true
    }
    return u.ID == 0
}
user := User{}
// 可以直接添加一个绑定了结构体的方法
g.GenerateModel("people", gen.MethodAppend(user.IsEmpty))
// 也可以传入一个结构体，会将这个结构体上绑定的所有方法绑定到新生成的结构体上
g.GenerateModel("people", gen.MethodAppend(user))
```

指定生成的查询结构体字段类型
```Go
//package model
type ITime struct {
    time.Time
}

// 自定义数据结构体字段类型和查询结构体字段类型
g.ApplyBasic(g.GenerateModel("people", gen.FieldType("create_time","model.ITime"), gen.FieldGenType("create_time","Time")))

//package model
type User struct {
  ID int64
  Name string
  CreateTime ITime
}

func (u User) GetFieldGenType(f *schema.Field) string {
  if f.Name == "CreateTime" {
    return "Time"
  }
  return ""
}
// 自定义查询结构体类型
g.ApplyBasic(model.User{})
```

#### <span id="data-mapping">类型映射</span>

指定你期望的数据映射关系，如自定义数据库字段类型和 Go 类型的映射关系。

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

### <span id="field-expression">字段表达式</span>

#### <span id="create-field">创建字段</span>

在实际操作中，你不需要创建一个新的字段变量，这个创建过程将通过生成代码自动完成。

| Field Type | Detail Type           | Create Function               | Supported Query Method                                       |
| ---------- | --------------------- | ------------------------------ | ------------------------------------------------------------ |
| generic    | field                 | NewField                       | IsNull/IsNotNull/Count/Eq/Neq/Gt/Gte/Lt/Lte/Like             |
| int        | int/int8/.../int64    | NewInt/NewInt8/.../NewInt64    | Eq/Neq/Gt/Gte/Lt/Lte/In/NotIn/Between/NotBetween/Like/NotLike/Add/Sub/Mul/Div/Mod/FloorDiv/RightShift/LeftShift/BitXor/BitAnd/BitOr/BitFlip |
| uint       | uint/uint8/.../uint64 | NewUint/NewUint8/.../NewUint64 | same with int                                                |
| float      | float32/float64       | NewFloat32/NewFloat64          | Eq/Neq/Gt/Gte/Lt/Lte/In/NotIn/Between/NotBetween/Like/NotLike/Add/Sub/Mul/Div/FloorDiv |
| string     | string/[]byte         | NewString/NewBytes             | Eq/Neq/Gt/Gte/Lt/Lte/Between/NotBetween/In(val/NotIn(val/Like/NotLike/Regexp/NotRegxp/FindInSet/FindInSetWith |
| bool       | bool                  | NewBool                        | Not/Is/And/Or/Xor/BitXor/BitAnd/BitOr                        |
| time       | time.Time             | NewTime                        | Eq/Neq/Gt/Gte/Lt/Lte/Between/NotBetween/In/NotIn/Add/Sub     |

创建字段示例：

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

### <span id="crud-api">CRUD 接口</span>

以下为一个 `user` 模型和 `DB` 模型的基础结构。

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

#### <span id="create">创建</span>

##### <span id="create-record">创建记录</span>

```go
// u refer to query.user
user := model.User{Name: "Modi", Age: 18, Birthday: time.Now()}

u := query.Use(db).User
err := u.WithContext(ctx).Create(&user) // pass pointer of data to Create

err // returns error
```

##### <span id="create-record-with-selected-fields">选择字段创建</span>

创建记录并为指定的字段赋值。

```go
u := query.Use(db).User
u.WithContext(ctx).Select(u.Name, u.Age).Create(&user)
// INSERT INTO `users` (`name`,`age`) VALUES ("modi", 18)
```

创建记录并通过 `Omit` 方法忽略传递字段的具体值 。

```go
u := query.Use(db).User
u.WithContext(ctx).Omit(u.Name, u.Age).Create(&user)
// INSERT INTO `users` (`Address`, `Birthday`) VALUES ("2021-08-17 20:54:12.000", 18)
```

##### <span id="batch-insert">批量创建</span>

`Create` 方法支持批量创建记录，只需要将对应模型（Model）的切片（slice）类型数据作为参数传入即可。GORM 将生成单个 SQL 语句来插入所有数据并返回对应内容全部主键的值。

```go
var users = []*model.User{{Name: "modi"}, {Name: "zhangqiang"}, {Name: "songyuan"}}
query.Use(db).User.WithContext(ctx).Create(users...)

for _, user := range users {
    user.ID // 1,2,3
}
```

你可以通过 `CreateInBatches` 方法可以指定批量创建记录的大小,如：

```go
var users = []*User{{Name: "modi_1"}, ...., {Name: "modi_10000"}}

// batch size 100
query.Use(db).User.WithContext(ctx).CreateInBatches(users, 100)
```

也可以通过全局配置方式，在初始化 GORM 时设置在 `gorm.Config` / `gorm.Session` 中对应配置 `CreateBatchSize`

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

#### <span id="query">查询</span>

##### <span id="retrieving-a-single-object">单条数据查询</span>

GROM 提供了 `First`、`Take`、`Last` 方法从数据库中查询单条数据，在查询数据库时会自动添加 `LIMIT 1` 条件，如果没有找到记录则返回错误 `ErrRecordNotFound`。

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

##### <span id="retrieving-objects-with-primary-key">根据主键查询数据</span>

```go
u := query.Use(db).User

user, err := u.WithContext(ctx).Where(u.ID.Eq(10)).First()
// SELECT * FROM users WHERE id = 10;

users, err := u.WithContext(ctx).Where(u.ID.In(1,2,3)).Find()
// SELECT * FROM users WHERE id IN (1,2,3);
```

如果主键是一个字符串类型数据（例如，主键是一个 uuid ），查询将写成如下形式：

```go
user, err := u.WithContext(ctx).Where(u.ID.Eq("1b74413f-f3b8-409f-ac47-e8c062e3472a")).First()
// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
```

##### <span id="retrieving-all-objects">查询所有数据</span>

```go
u := query.Use(db).User

// Get all records
users, err := u.WithContext(ctx).Find()
// SELECT * FROM users;
```

##### <span id="conditions">条件查询</span>

###### <span id="string-conditions">字符串和基础查询套件</span>

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

###### <span id="inline-condition">内联条件查询</span>

```go
u := query.Use(db).User

// Get by primary key if it were a non-integer type
user, err := u.WithContext(ctx).Where(u.ID.Eq("string_primary_key")).First()
// SELECT * FROM users WHERE id = 'string_primary_key';

// Plain SQL
users, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).Find()
// SELECT * FROM users WHERE name = "modi";

users, err := u.WithContext(ctx).Where(u.Name.Neq("modi"), u.Age.Gt(17)).Find()
// SELECT * FROM users WHERE name <> "modi" AND age > 17;
```

###### <span id="not-conditions">取反（not）查询</span>

构建取反（not）查询条件，效果类似于 `Where` 查询

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

###### <span id="or-conditions">或（or）查询</span>

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(u.Role.Eq("admin")).Or(u.Role.Eq("super_admin")).Find()
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
```

###### <span id="group-conditions">组合查询</span>

使用 `Where` / `Or` / `Not` 进行组合查询，可以轻松编写复杂的 SQL 查询

```go
p := query.Use(db).Pizza
pd := p.WithContext(ctx)

pizzas, err := pd.Where(
    pd.Where(p.Pizza.Eq("pepperoni")).
        Where(pd.Where(p.Size.Eq("small")).Or(p.Size.Eq("medium"))),
).Or(
    pd.Where(p.Pizza.Eq("hawaiian")).Where(p.Size.Eq("xlarge")),
).Find()

// SELECT * FROM `pizzas` WHERE (pizza = "pepperoni" AND (size = "small" OR size = "medium")) OR (pizza = "hawaiian" AND size = "xlarge")
```

###### <span id="selecting-specific-fields">指定字段查询</span>

`Select` 方法允许你指定要从数据库中查询的字段。否则，GORM 将默认选择所有字段。

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Select(u.Name, u.Age).Find()
// SELECT name, age FROM users;

u.WithContext(ctx).Select(u.Age.Avg()).Rows()
// SELECT Avg(age) FROM users;
```

###### <span id="tuple-query">元组查询</span>

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(u.WithContext(ctx).Columns(u.ID, u.Name).In(field.Values([][]interface{}{{1, "modi"}, {2, "zhangqiang"}}))).Find()
// SELECT * FROM `users` WHERE (`id`, `name`) IN ((1,'humodi'),(2,'tom'));
```

###### <span id="json-query">JSON 查询</span>

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(gen.Cond(datatypes.JSONQuery("attributes").HasKey("role"))...).Find()
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`,'$.role') IS NOT NULL;
```

###### <span id="order">Order 排序</span>

从数据库查询数据时，指定数据排序的方式

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Order(u.Age.Desc(), u.Name).Find()
// SELECT * FROM users ORDER BY age DESC, name;

// Multiple orders
users, err := u.WithContext(ctx).Order(u.Age.Desc()).Order(u.Name).Find()
// SELECT * FROM users ORDER BY age DESC, name;
```

通过字符串获取想要排序的列

```go
u := query.Use(db).User

orderCol, ok := u.GetFieldByName(orderColStr) // orderColStr的值可以是"id"
if !ok {
  // User doesn't contains orderColStr
}

users, err := u.WithContext(ctx).Order(orderCol).Find()
// SELECT * FROM users ORDER BY age;

// OR Desc
users, err := u.WithContext(ctx).Order(orderCol.Desc()).Find()
// SELECT * FROM users ORDER BY age DESC;
```

###### <span id="limit--offset">分页查询（Limit & Offset）</span>

分页查询方法，其中：
* `Limit` 指定要检索的最大记录数
* `Offset` 指定在开始返回记录之前要跳过的记录数（当前分页的位置）

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

###### <span id="group-by--having">分组查询（Group By & Having）</span>

```go
u := query.Use(db).User

var users []struct {
    Name  string
    Total int
}
err := u.WithContext(ctx).Select(u.Name, u.ID.Count().As("total")).Group(u.Name).Scan(&users)
// SELECT name, count(id) as total FROM `users` GROUP BY `name`

err := u.WithContext(ctx).Select(u.Name, u.Age.Sum().As("total")).Where(u.Name.Like("%modi%")).Group(u.Name).Scan(&users)
// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "%modi%" GROUP BY `name`

err := u.WithContext(ctx).Select(u.Name, u.Age.Sum().As("total")).Group(u.Name).Having(u.Name.Eq("group")).Scan(&users)
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

var results []struct {
    Date  time.Time
    Total int
}

o.WithContext(ctx).Select(o.CreateAt.Date().As("date"), o.WithContext(ctx).Amount.Sum().As("total")).Group(o.CreateAt.Date()).Having(u.Amount.Sum().Gt(100)).Scan(&results)
```

###### <span id="distinct">去重（Distinct）</span>

从 Model 中指定字段查询，并对值进行去重。

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Distinct(u.Name, u.Age).Order(u.Name, u.Age.Desc()).Find()
```

`Distinct` works with `Pluck` and `Count` too

###### <span id="joins">联表查询（Joins）</span>

联表查询方法，`Join` 方法对应 `inner join`，此外还有 `LeftJoin` 方法和 `RightJoin` 方法。

```go
q := query.Use(db)
u := q.User
e := q.Email
c := q.CreditCard

type Result struct {
    Name  string
    Email string
    ID    int64
}

var result Result

err := u.WithContext(ctx).Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Scan(&result)
// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

// self join
var result Result
u2 := u.As("u2")
err := u.WithContext(ctx).Select(u.Name, u2.ID).LeftJoin(u2, u2.ID.EqCol(u.ID)).Scan(&result)
// SELECT users.name, u2.id FROM `users` left join `users` u2 on u2.id = users.id

//join with sub query
var result Result
e2 := e.As("e2")
err := u.WithContext(ctx).Select(u.Name, e2.Email).LeftJoin(e.WithContext(ctx).Select(e.Email, e.UserID).Where(e.UserID.Gt(100)).As("e2"), e2.UserID.EqCol(u.ID)).Scan(&result)
// SELECT users.name, e2.email FROM `users` left join (select email,user_id from emails  where user_id > 100) as e2 on e2.user_id = users.id

rows, err := u.WithContext(ctx).Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Rows()
for rows.Next() {
  ...
}

var results []Result

err := u.WithContext(ctx).Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Scan(&results)

// multiple joins with parameter
users := u.WithContext(ctx).Join(e, e.UserID.EqCol(u.id), e.Email.Eq("modi@example.org")).Join(c, c.UserID.EqCol(u.ID)).Where(c.Number.Eq("411111111111")).Find()
```

##### <span id="subquery">子查询</span>

子查询可以嵌套在查询中，GEN 可以在使用 `Dao` 对象作为参数时生成子查询

```go
o := query.Use(db).Order
u := query.Use(db).User

orders, err := o.WithContext(ctx).Where(o.WithContext(ctx).Columns(o.Amount).Gt(o.WithContext(ctx).Select(o.Amount.Avg())).Find()
// SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

subQuery := u.WithContext(ctx).Select(u.Age.Avg()).Where(u.Name.Like("name%"))
users, err := u.WithContext(ctx).Select(u.Age.Avg().As("avgage")).Group(u.Name).Having(u.WithContext(ctx).Columns(u.Age.Avg()).Gt(subQuery).Find()
// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")

// 找出交易数量在100和200之间的用户
subQuery1 := o.WithContext(ctx).Select(o.ID).Where(o.UserID.EqCol(u.ID), o.Amount.Gt(100))
subQuery2 := o.WithContext(ctx).Select(o.ID).Where(o.UserID.EqCol(u.ID), o.Amount.Gt(200))
u.WithContext(ctx).Exists(subQuery1).Not(u.WithContext(ctx).Exists(subQuery2)).Find()
// SELECT * FROM `users` WHERE EXISTS (SELECT `orders`.`id` FROM `orders` WHERE `orders`.`user_id` = `users`.`id` AND `orders`.`amount` > 100 AND `orders`.`deleted_at` IS NULL) AND NOT EXISTS (SELECT `orders`.`id` FROM `orders` WHERE `orders`.`user_id` = `users`.`id` AND `orders`.`amount` > 200 AND `orders`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL
```

###### <span id="from-subquery">From 子查询</span>

通过 `Table` 方法构建出的子查询，可以直接放到 From 语句中:

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

###### <span id="update-from-subquery">从子查询更新</span>

通过子查询更新表字段

```go
u := query.Use(db).User
c := query.Use(db).Company

u.WithContext(ctx).Update(u.CompanyName, c.Select(c.Name).Where(c.ID.EqCol(u.CompanyID)))
// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);

u.WithContext(ctx).Where(u.Name.Eq("modi")).Update(u.CompanyName, c.Select(c.Name).Where(c.ID.EqCol(u.CompanyID)))
```

###### <span id="update-multiple-columns-from-subquery">从子查询更新多个字段</span>

针对 MySQL 提供同时更新多个字段的子查询：

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

##### <span id="transaction">事务</span>

要在事务中执行一组操作，一般的处理流程如下：

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

###### <span id="nested-transactions">嵌套事务</span>

GEN 支持嵌套事务，在一个大事务中嵌套子事务。

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

###### <span id="transactions-by-manual">手动事务</span>

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

###### <span id="savepointrollbackto">保存点/回滚</span>

GEN 提供了 `SavePoint` 和 `RollbackTo` 方法，用于保存和回滚事务点，例如：

```go
tx := q.Begin()
txCtx = tx.WithContext(ctx)

txCtx.User.Create(&user1)

tx.SavePoint("sp1")
txCtx.Create(&user2)
tx.RollbackTo("sp1") // Rollback user2

tx.Commit() // Commit user1
```

##### <span id="advanced-query">高级查询</span>

###### <span id="iteration">迭代</span>

Gen 支持通过 Rows 方法进行迭代（遍历）操作

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

###### <span id="findinbatches">批量查询</span>

FindInBatches 方法支持批量查询并处理记录

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

###### <span id="pluck">Pluck 方法</span>

Pluck 方法支持从数据库中查询单列并扫描成切片。如果要查询多列，请使用 `Select` 和 `Scan` 方法代替 `Pluck` 方法

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

###### <span id="scopes">Scopes 查询</span>

你可以声明一些常用或公共的条件方法，然后使用 `Scopes` 指定调用

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

###### <span id="count">Count 计数</span>

`Count` 方法用于获取查询结果数。

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

###### <span id="firstorinit">首条匹配或指定查询实例初始化条件（FirstOrInit）</span>

获取第一个匹配的记录或在给定条件下初始化一个新实例

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

如果未找到记录，则使用更多属性初始化结构，这些 `Attrs` 将不会用于构建 SQL 查询

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

###### <span id="firstorcreate">首条匹配或指定创建实例初始化条件（FirstOrCreate）</span>

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

如果希望创建的实例包含一些非查询条件的属性，则可以通过 `Attrs` 指定

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
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).Assign(u.Age.Value(20)).FirstOrCreate()
// SELECT * FROM users WHERE name = 'modi' ORDER BY id LIMIT 1;
// UPDATE users SET age=20 WHERE id = 111;
// user -> User{ID: 111, Name: "modi", Age: 20}
```

#### <span id="association">关联关系（Association）</span>

GEN 会像 GORM 一样自动保存关联关系。关联关系 (BelongsTo/HasOne/HasMany/Many2Many) 重用了 GORM 的标签。
此功能目前仅支持现有模型。

##### <span id="relation">关联</span>

There are 4 kind of relationship.

```go
const (
    HasOne    RelationshipType = RelationshipType(schema.HasOne)    // HasOneRel has one relationship
    HasMany   RelationshipType = RelationshipType(schema.HasMany)   // HasManyRel has many relationships
    BelongsTo RelationshipType = RelationshipType(schema.BelongsTo) // BelongsToRel belongs to relationship
    Many2Many RelationshipType = RelationshipType(schema.Many2Many) // Many2ManyRel many to many relationship
)
```

###### <span id="relate-to-exist-model">关联已存在的模型</span>

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

GEN 会检测模型的关联关系：

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

###### <span id="relate-to-table-in-database">和数据库表关联</span>

关联必须由 `gen.FieldRelate` 指定声明。

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

如果是已经存在的关联 model, 则可以用 `gen.FieldRelateModel` 声明。

```go
customer := g.GenerateModel("customers", gen.FieldRelateModel(field.HasMany, "CreditCards", model.CreditCard{}, 
    &field.RelateConfig{
        // RelateSlice: true,
        GORMTag: "foreignKey:CustomerRefer",
    }),
)

g.ApplyBasic(custormer)
```

###### <span id="relate-config">关联配置</span>

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

##### <span id="operation">操作</span>

###### <span id="skip-auto-createupdate">跳过自动创建关联</span>

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

方法 `Field` 会用 ''." 连接一个严谨的字段名，例如：`u.BillingAddress.Field("Address1", "Street")` 等于 `BillingAddress.Address1.Street`

###### <span id="find-associations">查询关联</span>

查询匹配的关联

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

###### <span id="append-associations">添加关联</span>

为 `Many2Many`、`HasMany` 附加新的关联，替换 `has one`、`belongs to` 的当前关联

```go
u := query.Use(db).User

u.Languages.Model(&user).Append(&languageZH, &languageEN)

u.Languages.Model(&user).Append(&Language{Name: "DE"})

u.CreditCards.Model(&user).Append(&CreditCard{Number: "411111111111"})
```

###### <span id="replace-associations">替换关联</span>

使用新关联替换当前关联

```go
u.Languages.Model(&user).Replace(&languageZH, &languageEN)
```

###### <span id="delete-associations">删除关联</span>

删除源表数据和参数之间的关系，只删除引用，不会从数据库中删除这些对象。

```go
u := query.Use(db).User

u.Languages.Model(&user).Delete(&languageZH, &languageEN)

u.Languages.Model(&user).Delete([]*Language{&languageZH, &languageEN}...)
```

###### <span id="clear-associations">清除关联</span>

删除源表和关联表之间的所有引用映射，不会删除这些关联表

```go
u.Languages.Model(&user).Clear()
```

###### <span id="count-associations">统计关联</span>

返回当前关联的计数。

```go
u.Languages.Model(&user).Count()
```

###### <span id="delete-with-select">删除指定关联</span>

删除记录时，允许删除与用 `Select` 方法指定的对象存在 HasOne/HasMany/Many2Many 关系的关联，并删除关联数据，例如：

```go
u := query.Use(db).User

// delete user's account when deleting user
u.Select(u.Account).Delete(&user)

// delete user's Orders, CreditCards relations when deleting user
db.Select(u.Orders.Field(), u.CreditCards.Field()).Delete(&user)

// delete user's has one/many/many2many relations when deleting user
db.Select(field.AssociationFields).Delete(&user)
```

##### <span id="preloading">预加载</span>

此功能目前仅支持现有模型

###### <span id="preload">预加载（Preload）</span>

GEN 允许使用 `Preload` 在其他 SQL 中预先加载关系，例如：

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

###### <span id="preload-all">预加载全部数据（Preload All）</span>

`clause.Associations` 可以和 `Preload` 一起使用，类似于创建/更新时的 `Select`，你可以用它来 `Preload` 预加载所有关联，例如：

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

`clause.Associations` 不会加载嵌套关联, 嵌套关联加载可以用 [Nested Preloading](#nested-preloading) e.g:

```go
users, err := u.WithContext(ctx).Preload(u.Orders.OrderItems.Product).Find()
```

###### <span id="preload-with-select">预加载指定列（Preload with select）</span>

使用`Select`指定要查询的列名, 相应的外键必须被指定。

```go
type User struct {
  gorm.Model
  CreditCards []CreditCard `gorm:"foreignKey:UserRefer"`
}

type CreditCard struct {
  gorm.Model
  Number    string
  UserRefer uint
}

u := q.User
cc := q.CreditCard

// !!! 外键 "cc.UserRefer" 必须被指定
users, err := u.WithContext(ctx).Where(c.ID.Eq(1)).Preload(u.CreditCards.Select(cc.Number, cc.UserRefer)).Find()
// SELECT * FROM `credit_cards` WHERE `credit_cards`.`customer_refer` = 1 AND `credit_cards`.`deleted_at` IS NULL
// SELECT * FROM `customers` WHERE `customers`.`id` = 1 AND `customers`.`deleted_at` IS NULL LIMIT 1
```

###### <span id="nested-preloading">根据条件预加载</span>

GEN 允许预加载与条件关联，它的工作原理类似于内联条件。

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

user, err := u.WithContext(ctx).Where(u.ID.Eq(1)).Preload(u.Orders.Offset(100).Limit(20)).Take()
// SELECT * FROM users WHERE `user_id` = 1 LIMIT 20 OFFSET 100;
// SELECT * FROM `users` WHERE `users`.`id` = 1 LIMIT 1
```

###### <span id="nested-preloading">嵌套预加载</span>

GEN 支持嵌套预加载，例如：

```go
db.Preload(u.Orders.OrderItems.Product).Preload(u.CreditCard).Find(&users)

// Customize Preload conditions for `Orders`
// And GEN won't preload unmatched order's OrderItems then
db.Preload(u.Orders.On(o.State.Eq("paid"))).Preload(u.Orders.OrderItems).Find(&users)
```

#### <span id="update">更新</span>

##### <span id="update-single-column">更新单个字段</span>

使用 `Update` 更新单个列时，必须指定更新条件，否则会引发 `ErrMissingWhereClause` 错误，例如：

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

##### <span id="updates-multiple-columns">更新多个字段</span>

`Updates` 支持使用 `struct` 或 `map[string]interface{}` 更新多个字段，当使用 `struct` 更新时，默认只会更新非零字段

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

> **注意** 当通过 struct 更新的时候，GEN 将只会更新其非零值的字段，你可能需要用 `map` 去更新属性，或者用 `select` 去明确指定哪些字段是需要被更新的

##### <span id="update-selected-fields">更新指定字段</span>

如果你想更新选定的字段或在更新时忽略某些字段，可以使用 `Select`、`Omit`。

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

#### <span id="delete">删除</span>

##### <span id="delete-record">删除记录</span>

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

##### <span id="delete-with-primary-key">根据主键删除</span>

GEN 允许使用具有内联条件的主键来删除对象，这个操作适用于数字类型的主键。

```go
u.WithContext(ctx).Where(u.ID.In(1,2,3)).Delete()
// DELETE FROM users WHERE id IN (1,2,3);
```

##### <span id="batch-delete">批量删除</span>

在进行删除操作时，当指定的值没有具体的值时（如模糊查询），GEN 会进行批量删除，这会删除所有与查询条件匹配的记录。

```go
e := query.Use(db).Email

e.WithContext(ctx).Where(e.Name.Like("%modi%")).Delete()
// DELETE from emails where email LIKE "%modi%";
```

##### <span id="soft-delete">软删除</span>

如果你的 model 中包含有 `gorm.DeletedAt` 字段，则会自动执行软删除。

当使用软删除时，调用 `Delete` 不会将相关记录从数据库中删除，GORM 会将 `gorm.DeletedAt` 对应字段的值设置为当前时间，以表示删除状态和删除的时间。通过软删除的数据，无法使用普通的 Query 方法找到相应的记录。

```go
// Batch Delete
u.WithContext(ctx).Where(u.Age.Eq(20)).Delete()
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// Soft deleted records will be ignored when querying
users, err := u.WithContext(ctx).Where(u.Age.Eq(20)).Find()
// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;
```

如果你不想在实例化的对象中引用包含 `gorm.Model`，可以启用软删除功能，例如：

```go
type User struct {
    ID      int
    Deleted gorm.DeletedAt
    Name    string
}
```

##### <span id="find-soft-deleted-records">查询包含软删除的记录</span>

可以使用 `Unscoped`，你可以查询到软删除的记录。

```go
users, err := db.WithContext(ctx).Unscoped().Where(u.Age.Eq(20)).Find()
// SELECT * FROM users WHERE age = 20;
```

##### <span id="delete-permanently">永久删除</span>

通过 `Unscoped` 可以直接删除数据（物理删除），而不是逻辑删除（软删除）。

```go
o.WithContext(ctx).Unscoped().Where(o.ID.Eq(10)).Delete()
// DELETE FROM orders WHERE id=10;
```

### <span id="diy-method">自定义方法</span>

#### <span id="method-interface">接口定义</span>

自定义方法，需要通过 interface 定义。在方法上通过注释的方式描述具体的 SQL 查询逻辑，简单的 WHERE 查询可以用 `where()` 包住，复杂的查询需要写完整 SQL 可以直接用 `sql()` 包住或者忽略直接写 SQL，太长的 SQL 支持换行，如果有对方法的注释，只需要在前面加一个空行。

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

方法输入参数和返回值支持基础类型（int、string、bool等）、结构体和占位符（`gen.T`/`gen.M`/`gen.RowsAffected`）,类型支持指针和数组，返回值最多返回一个值和一个 error。

用法(完整case见[快速开始](#quick-start))：
```go
// 在表user和company生成的结构体上实现model.Method中包含的方法
g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))
```

##### <span id="syntax-of-template">模板</span>

###### <span id="placeholder">占位符</span>

- `gen.T` 用于返回数据的结构体，会根据生成结构体或者数据库表结构自动生成
- `gen.M` 表示map[string]interface{},用于返回数据
- `gen.RowsAffected` 用于执行SQL进行更新或删除时候,用于返回影响行数(类型为：int64)。
- `@@table` 仅用于SQL语句中，查询的表名，如果没有传参，会根据结构体或者表名自动生成
- `@@<columnName>` 仅用于SQL语句中，当表名或者字段名可控时候，用@@占位，name为可变参数名，需要函数传入。
- `@<name>` 仅用于SQL语句中，当数据可控时候，用@占位，name为可变参数名，需要函数传入

###### <span id="template">模板</span>

逻辑操作必须包裹在 `{{}}` 中，如 `{{if}}`，结束语句必须是 `{{end}}`，所有的语句都可以嵌套。`{{}}` 中的语法除了 `{{end}}` 其它的都是 Golang 语法。

- `if`/`else if`/`else` if 子句通过判断满足条件拼接字符串到SQL。
- `where` where 子句只有在内容不为空时候插入 where，若子句的开头和结尾为 where 语句连接关键字 `AND` 或 `OR`，会将它们去除。
- `Set` set 子句只有在内容不为空时候插入，若子句的开头和结尾为set连接关键字 `,` 会将它们去除。
- `for` for 子句会遍历数组并将其内容插入到 SQL 中,需要注意之前的连接词。
- `...` 未完待续

###### <span id="if-clause">`If` 子句</span>

```
{{if cond1}}
    // do something here
{{else if cond2}}
    // do something here
{{else}}
    // do something here
{{end}}
```

方法使用案例:

```go
// select * from users where 
//  {{if name !=""}} 
//      name=@name
//  {{end}}
methond(name string) (gen.T,error) 
```

SQL模板使用案例:

```
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

###### <span id="where-clause">`Where` 子句</span>

```
{{where}}
    // do something here
{{end}}
```

方法使用案例:

```go
// select * from 
//  {{where}}
//      id=@id
//  {{end}}
methond(id int) error
```

SQL模板使用案例:

```
select * from @@table 
{{where}}
    {{if cond}}id=@id {{end}}
    {{if name != ""}}@@key=@value{{end}}
{{end}}
```

###### <span id="set-clause">`Set` 子句</span>

```
{{set}}
    // sepecify update expression here
{{end}}
```

方法使用案例:

```go
// update users 
//  {{set}}
//      name=@name
//  {{end}}
// where id=@id
methond(name string,id int) error
```

SQL模板使用案例:

```
update @@table 
{{set}}
    {{if name!=""}} name=@name {{end}}
    {{if age>0}} age=@age {{end}}
{{end}}
where id=@id
```

###### <span id="for-clause">`For` 子句</span>

```
{{for _,name:=range names}}
    // do something here
{{end}}
```

方法使用案例:

```go
// select * from users where id>0 
//  {{for _,name:=range names}} 
//      and name=@name
//  {{end}}
methond(names []string) (gen.T,error) 
```

SQL 模板使用案例:

```
select * from @@table where
  {{for index,name:=range names}}
     {{if index >0}} 
        OR
     {{end}}
     name=@name
  {{end}}
```

##### <span id="method-interface-example">方法接口示例</span>

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
	
    // select * from @@table
    //  {{where}}
    //      {{for _,user:=range users}}
    //          {{if user.Age >18}
    //              OR name=@user.Name 
    //          {{end}}
    //      {{end}}
    //  {{end}}
    FindByOrList(users []gen.T) ([]gen.T, error)
}
```

#### <span id="unit-test">单元测试</span>

如果设置了 `WithUnitTest` 方法，就会生成单元测试文件，生成通用查询功能的单元测试代码。

自定义方法的单元测试需要自定义对应的测试用例，它应该和测试文件放在同一个包里。

一个测试用例包含输入和期望结果，输入应和对应的方法参数匹配，期望应和对应的方法返回值相匹配。这将在测试中被断言为 “**Equal（相等）**”。

```go
package query

type Input struct {
  Args []interface{}
}

type Expectation struct {
  Ret []interface{}
}

type TestCase struct {
  Input
  Expectation
}

/* Table student */

var StudentFindByIdTestCase = []TestCase {
  {
    Input{[]interface{}{1}},
    Expectation{[]interface{}{nil, nil}},
  },
}
```

相应测试代码：

```go
//FindById select * from @@table where id = @id
func (s studentDo) FindById(id int64) (result *model.Student, err error) {
    ///
}

func Test_student_FindById(t *testing.T) {
    student := newStudent(db)
    do := student.WithContext(context.Background()).Debug()

    for i, tt := range StudentFindByIdTestCase {
        t.Run("FindById_"+strconv.Itoa(i), func(t *testing.T) {
            res1, res2 := do.FindById(tt.Input.Args[0].(int64))
            assert(t, "FindById", res1, tt.Expectation.Ret[0])
            assert(t, "FindById", res2, tt.Expectation.Ret[1])
        })
    }
}
```

#### <span id="smart-select-fields">智能选择字段</span>

GEN 允许使用 `Select` 选择特定字段，如果你经常在应用程序中使用 ，也许你想为 API 使用定义一个更小的结构，它可以自动选择特定字段，例如：

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

### <span id="advanced-topics">高级教程</span>

#### <span id="hints">查询优化提示（Hints）</span>

优化器选择某个查询执行计划，GORM 支持 `gorm.io/hints`，例如：

```go
import "gorm.io/hints"

u := query.Use(db).User

users, err := u.WithContext(ctx).Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find()
// SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`
```

索引提示允许将索引提示传递给数据库，以避免查询计划器选择劣质的查询计划。

```go
import "gorm.io/hints"

u := query.Use(db).User

users, err := u.WithContext(ctx).Clauses(hints.UseIndex("idx_user_name")).Find()
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

users, err := u.WithContext(ctx).Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find()
// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)"
```

## <span id="binary">二进制命令行工具安装</span>

通过二进制文件安装 gen 命令行工具:

```bash
go install gorm.io/gen/tools/gentool@latest
```

用法：

```bash
$ gentool -h
Usage of gentool:
  -c string
      is path for gen.yml
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
  -onlyModel
      only generate models (without query file)
  -withUnitTest
      generate unit test for query code
```

示例:

``` bash
gentool -dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" -tables "orders,doctor"
gentool -c "./gen.yml"
```

配置文件示例:


config example:

```
version: "0.1"
database:
  # consult[https://gorm.io/docs/connecting_to_the_database.html]"
  dsn : "username:password@tcp(address:port)/db?charset=utf8mb4&parseTime=true&loc=Local"
  # input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html]
  db  : "mysql"
  # enter the required data table or leave it blank.You can input : orders,users,goods
  tables  : "user"
  # specify a directory for output
  outPath :  "./dao/query"
  # query code file name, default: gen.go
  outFile :  ""
  # generate unit test for query code
  withUnitTest  : false
  # generated model code's package name
  modelPkgName  : ""
  # generate with pointer when field is nullable
  fieldNullable : false
  # generate field with gorm index tag
  fieldWithIndexTag : false
  # generate field with gorm column type tag
  fieldWithTypeTag  : false

```

## <span id="maintainers">维护者</span>

[@riverchu](https://github.com/riverchu) [@iDer](https://github.com/idersec) [@qqxhb](https://github.com/qqxhb) [@dino-ma](https://github.com/dino-ma)

[@jinzhu](https://github.com/jinzhu)

## <span id="contributing">如何参与贡献</span>

你可以让 GORM/GEN 变得更好

## <span id="license">开源许可协议</span>

采用 [MIT 许可协议](https://github.com/go-gorm/gen/blob/master/License) 发布
