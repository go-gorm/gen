# GORM/GEN

[![GoVersion](https://img.shields.io/github/go-mod/go-version/go-gorm/gen)](https://github.com/go-gorm/gen/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/go-gorm/gen)](https://github.com/go-gorm/gen/releases)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/gorm.io/gen?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-gorm/gen)](https://goreportcard.com/report/github.com/go-gorm/gen)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![OpenIssue](https://img.shields.io/github/issues/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aopen+is%3Aissue)
[![ClosedIssue](https://img.shields.io/github/issues-closed/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aissue+is%3Aclosed)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/go-gorm/gen)](https://www.tickgit.com/browse?repo=github.com/go-gorm/gen)

Gen: Friendly & Safer [GORM](https://github.com/go-gorm/gorm) powered by Code Generation.

## Multilingual README

[English Version](./README.md) | [中文版本](./README.ZH_CN.md)

## Overview

- CRUD or DIY query method code generation
- Auto migration from database to code
- Transactions, Nested Transactions, Save Point, RollbackTo to Saved Point
- Competely compatible with GORM
- Developer Friendly
- Multiple Generate modes

## Contents

- [GORM/GEN](#gormgen)
  - [Multilingual README](#multilingual-readme)
  - [Overview](#overview)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Quick start](#quick-start)
    - [Project Directory](#project-directory)
  - [API Examples](#api-examples)
    - [Generate](#generate)
      - [Generate Mode](#generate-mode)
      - [Generate Model](#generate-model)
      - [Data Mapping](#data-mapping)
    - [Field Expression](#field-expression)
      - [Create Field](#create-field)
    - [CRUD API](#crud-api)
      - [Create](#create)
        - [Create record](#create-record)
        - [Create record with selected fields](#create-record-with-selected-fields)
        - [Batch Insert](#batch-insert)
      - [Query](#query)
        - [Retrieving a single object](#retrieving-a-single-object)
        - [Retrieving objects with primary key](#retrieving-objects-with-primary-key)
        - [Retrieving all objects](#retrieving-all-objects)
        - [Conditions](#conditions)
          - [String Conditions](#string-conditions)
          - [Inline Condition](#inline-condition)
          - [Not Conditions](#not-conditions)
          - [Or Conditions](#or-conditions)
          - [Group Conditions](#group-conditions)
          - [Selecting Specific Fields](#selecting-specific-fields)
          - [Tuple Query](#tuple-query)
          - [JSON Query](#json-query)
          - [Order](#order)
          - [Limit & Offset](#limit--offset)
          - [Group By & Having](#group-by--having)
          - [Distinct](#distinct)
          - [Joins](#joins)
        - [SubQuery](#subquery)
          - [From SubQuery](#from-subquery)
          - [Update from SubQuery](#update-from-subquery)
          - [Update multiple columns from SubQuery](#update-multiple-columns-from-subquery)
        - [Transaction](#transaction)
          - [Nested Transactions](#nested-transactions)
          - [Transactions by manual](#transactions-by-manual)
          - [SavePoint/RollbackTo](#savepointrollbackto)
        - [Advanced Query](#advanced-query)
          - [Iteration](#iteration)
          - [FindInBatches](#findinbatches)
          - [Pluck](#pluck)
          - [Scopes](#scopes)
          - [Count](#count)
          - [FirstOrInit](#firstorinit)
          - [FirstOrCreate](#firstorcreate)
      - [Association](#association)
        - [Relation](#relation)
          - [Relate to exist model](#relate-to-exist-model)
          - [Relate to table in database](#relate-to-table-in-database)
          - [Relate Config](#relate-config)
        - [Operation](#operation)
          - [Skip Auto Create/Update](#skip-auto-createupdate)
          - [Find Associations](#find-associations)
          - [Append Associations](#append-associations)
          - [Replace Associations](#replace-associations)
          - [Delete Associations](#delete-associations)
          - [Clear Associations](#clear-associations)
          - [Count Associations](#count-associations)
          - [Delete with Select](#delete-with-select)
        - [Preloading](#preloading)
          - [Preload](#preload)
          - [Preload All](#preload-all)
          - [Preload with select](#preload-with-select)
          - [Preload with conditions](#preload-with-conditions)
          - [Nested Preloading](#nested-preloading)
      - [Update](#update)
        - [Update single column](#update-single-column)
        - [Updates multiple columns](#updates-multiple-columns)
        - [Update selected fields](#update-selected-fields)
      - [Delete](#delete)
        - [Delete record](#delete-record)
        - [Delete with primary key](#delete-with-primary-key)
        - [Batch Delete](#batch-delete)
        - [Soft Delete](#soft-delete)
        - [Find soft deleted records](#find-soft-deleted-records)
        - [Delete permanently](#delete-permanently)
    - [DIY method](#diy-method)
      - [Method interface](#method-interface)
        - [Syntax of template](#syntax-of-template)
          - [placeholder](#placeholder)
          - [template](#template)
          - [`If` clause](#if-clause)
          - [`Where` clause](#where-clause)
          - [`Set` clause](#set-clause)
          - [`For` clause](#for-clause)
        - [Method interface example](#method-interface-example)
      - [Unit Test](#unit-test)
      - [Smart select fields](#smart-select-fields)
    - [Advanced Topics](#advanced-topics)
      - [Hints](#hints)
  - [Binary](#binary)
  - [Maintainers](#maintainers)
  - [Contributing](#contributing)
  - [License](#license)

## Installation

To install Gen package, you need to install Go and set your Go workspace first.

1.The first need Go installed(version 1.14,1.14+ is required), then you can use the below Go command to install Gen.

```bash
go get -u gorm.io/gen
```

2.Import it in your code:

```go
import "gorm.io/gen"
```

## Quick start

**Emphasis**: All use cases in this doc are generated under `WithContext` mode. And if you generate code under `WithoutContext` mode, please remove `WithContext(ctx)` before you call any query method, it helps you make code more concise.

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
    g.ApplyBasic(model.User{}, g.GenerateModel("company"), g.GenerateModelAs("people", "Person", gen.FieldIgnore("address")))
    
    // apply diy interfaces on structs or table models
    g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))

    // execute the action of code generation
    g.Execute()
}
```

Generate Mode:

- `gen.WithoutContext` generate code without `WithContext` constraint
- `gen.WithDefaultQuery` generate code with a default global variable `Q` as a singleton

### Project Directory

Here is a template for best practices:

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

## API Examples

### Generate

#### Generate Mode

```go
 g := gen.NewGenerator(gen.Config{
        ...
        Mode: gen.WithoutContext|gen.WithDefaultQuery|gen.WithQueryInterface,
        ...
 })
```

- `WithDefaultQuery` generate default query struct
- `WithoutContext` generate code without context constrain
- `WithQueryInterface` generate interface instead of struct for querying

#### Generate Model

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

Field Generate **Options**

```go
FieldNew           // create new field
FieldIgnore        // ignore field
FieldIgnoreReg     // ignore field (match with regexp)
FieldRename        // rename field in struct
FieldComment       // specify field comment in generated struct
FieldType          // specify field type
FieldTypeReg       // specify field type (match with regexp)
FieldGenType       // specify field gen type
FieldGenTypeReg    // specify field gen type (match with regexp)
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

Generate model bind custom method
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
// add custom method to generated model struct
g.GenerateModel("people", gen.MethodAppend(user.IsEmpty))
// also you can input a struct,will bind all method
g.GenerateModel("people", gen.MethodAppend(user))
```

Generate model with custom gen type
```Go
//package model
type ITime struct {
    time.Time
}

// custom field type and gen type for table
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
// custom field gen type for struct
g.ApplyBasic(model.User{})
```

#### Data Mapping

Specify data mapping relationship to be whatever you want.

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

### Field Expression

#### Create Field

Actually, you're not supposed to create a new field variable, cause it will be accomplished in generated code.

| Field Type | Detail Type           | Create Function               | Supported Query Method                                       |
| ---------- | --------------------- | ------------------------------ | ------------------------------------------------------------ |
| generic    | field                 | NewField                       | IsNull/IsNotNull/Count/Eq/Neq/Gt/Gte/Lt/Lte/Like             |
| int        | int/int8/.../int64    | NewInt/NewInt8/.../NewInt64    | Eq/Neq/Gt/Gte/Lt/Lte/In/NotIn/Between/NotBetween/Like/NotLike/Add/Sub/Mul/Div/Mod/FloorDiv/RightShift/LeftShift/BitXor/BitAnd/BitOr/BitFlip |
| uint       | uint/uint8/.../uint64 | NewUint/NewUint8/.../NewUint64 | same with int                                                |
| float      | float32/float64       | NewFloat32/NewFloat64          | Eq/Neq/Gt/Gte/Lt/Lte/In/NotIn/Between/NotBetween/Like/NotLike/Add/Sub/Mul/Div/FloorDiv |
| string     | string/[]byte         | NewString/NewBytes             | Eq/Neq/Gt/Gte/Lt/Lte/Between/NotBetween/In(val/NotIn(val/Like/NotLike/Regexp/NotRegxp/FindInSet/FindInSetWith |
| bool       | bool                  | NewBool                        | Not/Is/And/Or/Xor/BitXor/BitAnd/BitOr                        |
| time       | time.Time             | NewTime                        | Eq/Neq/Gt/Gte/Lt/Lte/Between/NotBetween/In/NotIn/Add/Sub     |

Create field examples:

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

### CRUD API

Here is a basic struct `user` and struct `DB`.

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

#### Create

##### Create record

```go
// u refer to query.user
user := model.User{Name: "Modi", Age: 18, Birthday: time.Now()}

u := query.Use(db).User
err := u.WithContext(ctx).Create(&user) // pass pointer of data to Create

err // returns error
```

##### Create record with selected fields

Create a record and assign a value to the fields specified.

```go
u := query.Use(db).User
u.WithContext(ctx).Select(u.Name, u.Age).Create(&user)
// INSERT INTO `users` (`name`,`age`) VALUES ("modi", 18)
```

Create a record and ignore the values for fields passed to omit

```go
u := query.Use(db).User
u.WithContext(ctx).Omit(u.Name, u.Age).Create(&user)
// INSERT INTO `users` (`Address`, `Birthday`) VALUES ("2021-08-17 20:54:12.000", 18)
```

##### Batch Insert

To efficiently insert large number of records, pass a slice to the `Create` method. GORM will generate a single SQL statement to insert all the data and backfill primary key values.

```go
var users = []*model.User{{Name: "modi"}, {Name: "zhangqiang"}, {Name: "songyuan"}}
query.Use(db).User.WithContext(ctx).Create(users...)

for _, user := range users {
    user.ID // 1,2,3
}
```

You can specify batch size when creating with `CreateInBatches`, e.g:

```go
var users = []*User{{Name: "modi_1"}, ...., {Name: "modi_10000"}}

// batch size 100
query.Use(db).User.WithContext(ctx).CreateInBatches(users, 100)
```

It will works if you set `CreateBatchSize` in `gorm.Config` / `gorm.Session`

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

#### Query

##### Retrieving a single object

Generated code provides `First`, `Take`, `Last` methods to retrieve a single object from the database, it adds `LIMIT 1` condition when querying the database, and it will return the error `ErrRecordNotFound` if no record is found.

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

##### Retrieving objects with primary key

```go
u := query.Use(db).User

user, err := u.WithContext(ctx).Where(u.ID.Eq(10)).First()
// SELECT * FROM users WHERE id = 10;

users, err := u.WithContext(ctx).Where(u.ID.In(1,2,3)).Find()
// SELECT * FROM users WHERE id IN (1,2,3);
```

If the primary key is a string (for example, like a uuid), the query will be written as follows:

```go
user, err := u.WithContext(ctx).Where(u.ID.Eq("1b74413f-f3b8-409f-ac47-e8c062e3472a")).First()
// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
```

##### Retrieving all objects

```go
u := query.Use(db).User

// Get all records
users, err := u.WithContext(ctx).Find()
// SELECT * FROM users;
```

##### Conditions

###### String Conditions

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

###### Inline Condition

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

###### Not Conditions

Build NOT conditions, works similar to `Where`

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

###### Or Conditions

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(u.Role.Eq("admin")).Or(u.Role.Eq("super_admin")).Find()
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
```

###### Group Conditions

Easier to write complicated SQL query with Group Conditions

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

###### Selecting Specific Fields

`Select` allows you to specify the fields that you want to retrieve from database. Otherwise, GORM will select all fields by default.

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Select(u.Name, u.Age).Find()
// SELECT name, age FROM users;

u.WithContext(ctx).Select(u.Age.Avg()).Rows()
// SELECT Avg(age) FROM users;
```

###### Tuple Query

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(u.WithContext(ctx).Columns(u.ID, u.Name).In(field.Values([][]interface{}{{1, "modi"}, {2, "zhangqiang"}}))).Find()
// SELECT * FROM `users` WHERE (`id`, `name`) IN ((1,'humodi'),(2,'tom'));
```

###### JSON Query

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Where(gen.Cond(datatypes.JSONQuery("attributes").HasKey("role"))...).Find()
// SELECT * FROM `users` WHERE JSON_EXTRACT(`attributes`,'$.role') IS NOT NULL;
```

###### Order

Specify order when retrieving records from the database

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Order(u.Age.Desc(), u.Name).Find()
// SELECT * FROM users ORDER BY age DESC, name;

// Multiple orders
users, err := u.WithContext(ctx).Order(u.Age.Desc()).Order(u.Name).Find()
// SELECT * FROM users ORDER BY age DESC, name;
```

Get field by string

```go
u := query.Use(db).User

orderCol, ok := u.GetFieldByName(orderColStr) // maybe orderColStr == "id"
if !ok {
  // User doesn't contains orderColStr
}

users, err := u.WithContext(ctx).Order(orderCol).Find()
// SELECT * FROM users ORDER BY age;

// OR Desc
users, err := u.WithContext(ctx).Order(orderCol.Desc()).Find()
// SELECT * FROM users ORDER BY age DESC;
```

###### Limit & Offset

`Limit` specify the max number of records to retrieve
`Offset` specify the number of records to skip before starting to return the records

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

###### Distinct

Selecting distinct values from the model

```go
u := query.Use(db).User

users, err := u.WithContext(ctx).Distinct(u.Name, u.Age).Order(u.Name, u.Age.Desc()).Find()
```

`Distinct` works with `Pluck` and `Count` too

###### Joins

Specify Joins conditions

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

##### SubQuery

A subquery can be nested within a query, GEN can generate subquery when using a `Dao` object as param

```go
o := query.Use(db).Order
u := query.Use(db).User

orders, err := o.WithContext(ctx).Where(o.WithContext(ctx).Columns(o.Amount).Gt(o.WithContext(ctx).Select(o.Amount.Avg())).Find()
// SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

subQuery := u.WithContext(ctx).Select(u.Age.Avg()).Where(u.Name.Like("name%"))
users, err := u.WithContext(ctx).Select(u.Age.Avg().As("avgage")).Group(u.Name).Having(u.WithContext(ctx).Columns(u.Age.Avg()).Gt(subQuery).Find()
// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")

// Select users with orders between 100 and 200
subQuery1 := o.WithContext(ctx).Select(o.ID).Where(o.UserID.EqCol(u.ID), o.Amount.Gt(100))
subQuery2 := o.WithContext(ctx).Select(o.ID).Where(o.UserID.EqCol(u.ID), o.Amount.Gt(200))
u.WithContext(ctx).Exists(subQuery1).Not(u.WithContext(ctx).Exists(subQuery2)).Find()
// SELECT * FROM `users` WHERE EXISTS (SELECT `orders`.`id` FROM `orders` WHERE `orders`.`user_id` = `users`.`id` AND `orders`.`amount` > 100 AND `orders`.`deleted_at` IS NULL) AND NOT EXISTS (SELECT `orders`.`id` FROM `orders` WHERE `orders`.`user_id` = `users`.`id` AND `orders`.`amount` > 200 AND `orders`.`deleted_at` IS NULL) AND `users`.`deleted_at` IS NULL
```

###### From SubQuery

GORM allows you using subquery in FROM clause with method `Table`, for example:

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

###### Update from SubQuery

Update a table by using SubQuery

```go
u := query.Use(db).User
c := query.Use(db).Company

u.WithContext(ctx).Update(u.CompanyName, c.Select(c.Name).Where(c.ID.EqCol(u.CompanyID)))
// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);

u.WithContext(ctx).Where(u.Name.Eq("modi")).Update(u.CompanyName, c.Select(c.Name).Where(c.ID.EqCol(u.CompanyID)))
```

###### Update multiple columns from SubQuery

Update multiple columns by using SubQuery (for MySQL):

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

##### Transaction

To perform a set of operations within a transaction, the general flow is as below.

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

###### Nested Transactions

GEN supports nested transactions, you can rollback a subset of operations performed within the scope of a larger transaction, for example:

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

###### Transactions by manual

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

###### SavePoint/RollbackTo

GEN provides `SavePoint`, `RollbackTo` to save points and roll back to a savepoint, for example:

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

###### Iteration

GEN supports iterating through Rows

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

###### FindInBatches

Query and process records in batch

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

Query single column from database and scan into a slice, if you want to query multiple columns, use `Select` with `Scan` instead

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

`Scopes` allows you to specify commonly-used queries which can be referenced as method calls

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

Get matched records count

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

Get first matched record or initialize a new instance with given conditions

```go
u := query.Use(db).User

// User not found, initialize it with give conditions
user, err := u.WithContext(ctx).Where(u.Name.Eq("non_existing")).FirstOrInit()
// user -> User{Name: "non_existing"}

// Found user with `name` = `modi`
user, err := u.WithContext(ctx).Where(u.Name.Eq("modi")).FirstOrInit()
// user -> User{ID: 1, Name: "modi", Age: 17}
```

initialize struct with more attributes if record not found, those `Attrs` won’t be used to build SQL query

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

`Assign` attributes to struct regardless it is found or not, those attributes won’t be used to build SQL query and the final data won’t be saved into database

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

Get first matched record or create a new one with given conditions

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

Create struct with more attributes if record not found, those `Attrs` won’t be used to build SQL query

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

`Assign` attributes to the record regardless it is found or not and save them back to the database.

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

#### Association

GEN will auto-save associations as GORM do. The relationships (BelongsTo/HasOne/HasMany/Many2Many) reuse GORM's tag.
This feature only support exist model for now.

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

###### Relate to exist model

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

GEN will detect model's associations:

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

###### Relate to table in database

The association have to be speified by `gen.FieldRelate`

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

GEN will generate models with associated field:

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

If associated model already exists, `gen.FieldRelateModel` can help you build associations between them.

```go
customer := g.GenerateModel("customers", gen.FieldRelateModel(field.HasMany, "CreditCards", model.CreditCard{}, 
    &field.RelateConfig{
        // RelateSlice: true,
        GORMTag: "foreignKey:CustomerRefer",
    }),
)

g.ApplyBasic(custormer)
```

###### Relate Config

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

##### Operation

###### Skip Auto Create/Update

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

###### Find Associations

Find matched associations

```go
u := query.Use(db).User

languages, err = u.Languages.Model(&user).Find()
```

Find associations with conditions

```go
q := query.Use(db)
u := q.User

languages, err = u.Languages.Where(q.Language.Name.In([]string{"ZH","EN"})).Model(&user).Find()
```

###### Append Associations

Append new associations for `many to many`, `has many`, replace current association for `has one`, `belongs to`

```go
u := query.Use(db).User

u.Languages.Model(&user).Append(&languageZH, &languageEN)

u.Languages.Model(&user).Append(&Language{Name: "DE"})

u.CreditCards.Model(&user).Append(&CreditCard{Number: "411111111111"})
```

###### Replace Associations

Replace current associations with new ones

```go
u.Languages.Model(&user).Replace(&languageZH, &languageEN)
```

###### Delete Associations

Remove the relationship between source & arguments if exists, only delete the reference, won’t delete those objects from DB.

```go
u := query.Use(db).User

u.Languages.Model(&user).Delete(&languageZH, &languageEN)

u.Languages.Model(&user).Delete([]*Language{&languageZH, &languageEN}...)
```

###### Clear Associations

Remove all reference between source & association, won’t delete those associations

```go
u.Languages.Model(&user).Clear()
```

###### Count Associations

Return the count of current associations

```go
u.Languages.Model(&user).Count()
```

###### Delete with Select

You are allowed to delete selected has one/has many/many2many relations with `Select` when deleting records, for example:

```go
u := query.Use(db).User

// delete user's account when deleting user
u.Select(u.Account).Delete(&user)

// delete user's Orders, CreditCards relations when deleting user
db.Select(u.Orders.Field(), u.CreditCards.Field()).Delete(&user)

// delete user's has one/many/many2many relations when deleting user
db.Select(field.AssociationFields).Delete(&user)
```

##### Preloading

This feature only support exist model for now.

###### Preload

GEN allows eager loading relations in other SQL with `Preload`, for example:

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

`clause.Associations` can work with `Preload` similar like `Select` when creating/updating, you can use it to `Preload` all associations, for example:

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

`clause.Associations` won’t preload nested associations, but you can use it with [Nested Preloading](#nested-preloading) together, e.g:

```go
users, err := u.WithContext(ctx).Preload(u.Orders.OrderItems.Product).Find()
```

###### Preload with select

Specify selected columns with method `Select`. Foregin key must be selected.

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

// !!! Foregin key "cc.UserRefer" must be selected
users, err := u.WithContext(ctx).Where(c.ID.Eq(1)).Preload(u.CreditCards.Select(cc.Number, cc.UserRefer)).Find()
// SELECT * FROM `credit_cards` WHERE `credit_cards`.`customer_refer` = 1 AND `credit_cards`.`deleted_at` IS NULL
// SELECT * FROM `customers` WHERE `customers`.`id` = 1 AND `customers`.`deleted_at` IS NULL LIMIT 1
```

###### Preload with conditions

GEN allows Preload associations with conditions, it works similar to Inline Conditions.

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

###### Nested Preloading

GEN supports nested preloading, for example:

```go
db.Preload(u.Orders.OrderItems.Product).Preload(u.CreditCard).Find(&users)

// Customize Preload conditions for `Orders`
// And GEN won't preload unmatched order's OrderItems then
db.Preload(u.Orders.On(o.State.Eq("paid"))).Preload(u.Orders.OrderItems).Find(&users)
```

#### Update

##### Update single column

When updating a single column with `Update`, it needs to have any conditions or it will raise error `ErrMissingWhereClause`, for example:

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

##### Updates multiple columns

`Updates` supports update with `struct` or `map[string]interface{}`, when updating with `struct` it will only update non-zero fields by default

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

##### Update selected fields

If you want to update selected fields or ignore some fields when updating, you can use `Select`, `Omit`

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

#### Delete

##### Delete record

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

##### Delete with primary key

GEN allows to delete objects using primary key(s) with inline condition, it works with numbers.

```go
u.WithContext(ctx).Where(u.ID.In(1,2,3)).Delete()
// DELETE FROM users WHERE id IN (1,2,3);
```

##### Batch Delete

The specified value has no primary value, GEN will perform a batch delete, it will delete all matched records

```go
e := query.Use(db).Email

e.WithContext(ctx).Where(e.Name.Like("%modi%")).Delete()
// DELETE from emails where email LIKE "%modi%";
```

##### Soft Delete

If your model includes a `gorm.DeletedAt` field (which is included in `gorm.Model`), it will get soft delete ability automatically!

When calling `Delete`, the record WON’T be removed from the database, but GORM will set the `DeletedAt`‘s value to the current time, and the data is not findable with normal Query methods anymore.

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

##### Find soft deleted records

You can find soft deleted records with `Unscoped`

```go
users, err := db.WithContext(ctx).Unscoped().Where(u.Age.Eq(20)).Find()
// SELECT * FROM users WHERE age = 20;
```

##### Delete permanently

You can delete matched records permanently with `Unscoped`

```go
o.WithContext(ctx).Unscoped().Where(o.ID.Eq(10)).Delete()
// DELETE FROM orders WHERE id=10;
```

### DIY method

#### Method interface

The DIY method needs to be defined through the interface. In the method, the specific SQL query logic is described in the way of comments. Simple WHERE queries can be wrapped in `where()`. When using complex queries, you need to write complete SQL. You can directly wrap them in `sql()` or write SQL directly. If there are some comments on the method, just add a blank line comment in the middle.

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
Method input parameters and return values support basic types (`int`, `string`, `bool`...), struct and placeholders (`gen.T`/`gen.M`/`gen.RowsAffected`), and types support pointers and arrays. The return value is at most a value and an error.

Usage(complete case on [Quick start](#quick-start)):
```go
// implement model.Method on table "user" and "comany"
g.ApplyInterface(func(method model.Method) {}, model.User{}, g.GenerateModel("company"))
```

##### Syntax of template

###### placeholder

- `gen.T` represents specified `struct` or `table`
- `gen.M` represents `map[string]interface`
- `gen.RowsAffected` represents SQL executed `rowsAffected` (type:int64)
- `@@table`  represents table's name (if method's parameter doesn't contains variable `table`, GEN will generate `table` from model struct)
- `@@<columnName>` represents column's name or table's name
- `@<name>` represents normal query variable

###### template

Logical operations must be wrapped in `{{}}`,and end must used `{{end}}`, All templates support nesting

- `if`/`else if`/`else` the condition accept a bool parameter or operation expression which conforms to Golang syntax.
- `where` The `where` clause will be inserted only if the child elements return something. The key word  `and` or `or`  in front of clause will be removed. And `and` will be added automatically when there is no junction keyword between query condition clause.
- `Set` The  `set` clause will be inserted only if the child elements return something. The `,` in front of columns array will be removed.And `,` will be added automatically when there is no junction keyword between query coulmns.
- `for` The  `for` clause traverses an array according to golang syntax and inserts its contents into SQL,supports array of struct.
- `...` Coming soon

###### `If` clause

```
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
// select * from users where 
//  {{if name !=""}} 
//      name=@name
//  {{end}}
methond(name string) (gen.T,error)
```

Use case in raw SQL template:

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

###### `Where` clause

```
{{where}}
    // do something here
{{end}}
```

Use case in raw SQL

```go
// select * from 
//  {{where}}
//      id=@id
//  {{end}}
methond(id int) error
```

Use case in raw SQL template

```
select * from @@table 
{{where}}
    {{if cond}} id=@id, {{end}}
    {{if name != ""}} @@key=@value, {{end}}
{{end}}
```

###### `Set` clause

```
{{set}}
    // sepecify update expression here
{{end}}
```

Use case in raw SQL

```go
// update users 
//  {{set}}
//      name=@name
//  {{end}}
// where id=@id
methond(name string,id int) error
```

Use case in raw SQL template

```
update @@table 
{{set}}
    {{if name!=""}} name=@name {{end}}
    {{if age>0}} age=@age {{end}}
{{end}}
where id=@id
```
###### `For` clause

```
{{for _,name:=range names}}
    // do something here
{{end}}
```

Use case in raw SQL:

```go
// select * from users where id>0 
//  {{for _,name:=range names}} 
//      and name=@name
//  {{end}}
methond(names []string) (gen.T,error) 
```

Use case in raw SQL template:

```
select * from @@table where
  {{for index,name:=range names}}
     {{if index >0}} 
        OR
     {{end}}
     name=@name
  {{end}}
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

    // select * from @@table
    //  {{where}}
    //      {{for _,user:=range users}}
    //          {{if user.Age >18}
    //              OR name=@user.Name 
    //         {{end}}
    //      {{end}}
    //  {{end}}
    FindByOrList(users []gen.T) ([]gen.T, error)
}
```

#### Unit Test

Unit test file will be generated if `WithUnitTest` is set, which will generate unit test for general query function.

Unit test for DIY method need diy testcase, which should place in the same package with test file.

A testcase contains input and expectation result, input should match the method arguments, expectation should match method return values, which will be asserted **Equal** in test.

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

Corresponding test

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

#### Smart select fields

GEN allows select specific fields with `Select`, if you often use this in your application, maybe you want to define a smaller struct for API usage which can select specific fields automatically, for example:

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

### Advanced Topics

#### Hints

Optimizer hints allow to control the query optimizer to choose a certain query execution plan, GORM supports it with `gorm.io/hints`, e.g:

```go
import "gorm.io/hints"

u := query.Use(db).User

users, err := u.WithContext(ctx).Clauses(hints.New("MAX_EXECUTION_TIME(10000)")).Find()
// SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`
```

Index hints allow passing index hints to the database in case the query planner gets confused.

```go
import "gorm.io/hints"

u := query.Use(db).User

users, err := u.WithContext(ctx).Clauses(hints.UseIndex("idx_user_name")).Find()
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

users, err := u.WithContext(ctx).Clauses(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find()
// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)"
```

## Binary

Install GEN as a binary tool:

```bash
go install gorm.io/gen/tools/gentool@latest
```

usage:

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

example:

``` bash
gentool -dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" -tables "orders,doctor"
gentool -c "./gen.yml"
```

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

## Maintainers

[@riverchu](https://github.com/riverchu) [@iDer](https://github.com/idersec) [@qqxhb](https://github.com/qqxhb) [@dino-ma](https://github.com/dino-ma)

[@jinzhu](https://github.com/jinzhu)

## Contributing

You can help to deliver a better GORM/GEN

## License

Released under the [MIT License](https://github.com/go-gorm/gen/blob/master/License)
