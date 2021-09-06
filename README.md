# GORM/GEN

[![GoVersion](https://img.shields.io/github/go-mod/go-version/go-gorm/gen)](https://github.com/go-gorm/gen/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/go-gorm/gen)](https://github.com/go-gorm/gen/releases)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/gorm.io/gen?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-gorm/gen)](https://goreportcard.com/report/github.com/go-gorm/gen)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![OpenIssue](https://img.shields.io/github/issues/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aopen+is%3Aissue)
[![ClosedIssue](https://img.shields.io/github/issues-closed/go-gorm/gen)](https://github.com/go-gorm/gen/issues?q=is%3Aissue+is%3Aclosed)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/go-gorm/gen)](https://www.tickgit.com/browse?repo=github.com/go-gorm/gen)

The code generator base on [GORM](https://github.com/go-gorm/gorm), aims to be developer friendly.

## Overview

- CRUD or DIY query method code generation
- Auto migration from database to code
- Transactions, Nested Transactions, Save Point, RollbackTo to Saved Point
- Competely compatible with GORM
- Developer Friendly

## Contents

- [GORM/GEN](#gormgen)
  - [Overview](#overview)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Quick start](#quick-start)
    - [Project Directory](#project-directory)
  - [API Examples](#api-examples)
    - [Generate](#generate)
      - [Generate Model](#generate-model)
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
          - [Order](#order)
          - [Limit & Offset](#limit--offset)
          - [Group By & Having](#group-by--having)
          - [Distinct](#distinct)
          - [Joins](#joins)
        - [SubQuery](#subquery)
          - [From SubQuery](#from-subquery)
          - [Update from SubQuery](#update-from-subquery)
        - [Advanced Query](#advanced-query)
          - [Iteration](#iteration)
          - [FindInBatches](#findinbatches)
          - [Pluck](#pluck)
          - [Scopes](#scopes)
          - [Count](#count)
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
        - [Method interface example](#method-interface-example)
      - [Smart select fields](#smart-select-fields)
    - [Advanced Topics](#advanced-topics)
      - [Hints](#hints)
  - [Maintainers](#maintainers)
  - [Contributing](#contributing)
  - [License](#license)

## Installation

To install Gen package, you need to install Go and set your Go workspace first.

1.The first need Go installed(version 1.14+ is required), then you can use the below Go command to install Gen.

```bash
go get -u gorm.io/gen
```

2.Import it in your code:

```go
import "gorm.io/gen"
```

## Quick start

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
    g := gen.NewGenerator(gen.Config{OutPath: "../dal/query"})
  
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
│       └── gorm_generated.go # generated code
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

#### Generate Model

```go
// generate a model struct map to table `people` in database
g.GenerateModel("people")

// generate a struct and specify struct's name
g.GenerateModelAs("people", "People")

// add option to ignore field
g.GenerateModel("people", gen.FieldIgnore("address"))
```

**Options**

```go
FieldIgnore // ignore field
FieldRename // rename field in struct
FieldType // specify field type
FieldTag // specify gorm and json tag
FieldJSONTag // specify json tag
FieldGORMTag // specify gorm tag
FieldNewTag // append new tag
```

### Field Expression

#### Create Field

Actually, you're not supposed to create a new field variable, cause it will be accomplished in generated code.

| Field Type | Detail Type           | Crerate Function               | Supported Query Method                                       |
| ---------- | --------------------- | ------------------------------ | ------------------------------------------------------------ |
| generic    | field                 | NewField                       | IsNull/IsNotNull/Count/Eq/Neq/Gt/Gte/Lt/Lte/Like             |
| int        | int/int8/.../int64    | NewInt/NewInt8/.../NewInt64    | Eq/Neq/Gt/Gte/Lt/Lte/In/NotIn/Between/NotBetween/Like/NotLike/Add/Sub/Mul/Div/Mod/FloorDiv/RightShift/LeftShift/BitXor/BitAnd/BitOr/BitFlip |
| uint       | uint/uint8/.../uint64 | NewUint/NewUint8/.../NewUint64 | same with int                                                |
| float      | float32/float64       | NewFloat32/NewFloat64          | Eq/Neq/Gt/Gte/Lt/Lte/In/NotIn/Between/NotBetween/Like/NotLike/Add/Sub/Mul/Div/FloorDiv |
| string     | string/[]byte         | NewString/NewBytes             | Eq/Neq/Gt/Gte/Lt/Lte/Between/NotBetween/In(val/NotIn(val/Like/NotLike/Regexp/NotRegxp |
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
err := u.Create(&user) // pass pointer of data to Create

err // returns error
```

##### Create record with selected fields

Create a record and assgin a value to the fields specified.

```go
u := query.Use(db).User
u.Select(u.Name, u.Age).Create(&user)
// INSERT INTO `users` (`name`,`age`) VALUES ("modi", 18)
```

Create a record and ignore the values for fields passed to omit

```go
u := query.Use(db).User
u.Omit(u.Name, u.Age).Create(&user)
// INSERT INTO `users` (`Address`, `Birthday`) VALUES ("2021-08-17 20:54:12.000", 18)
```

##### Batch Insert

To efficiently insert large number of records, pass a slice to the `Create` method. GORM will generate a single SQL statement to insert all the data and backfill primary key values.

```go
var users = []model.User{{Name: "modi"}, {Name: "zhangqiang"}, {Name: "songyuan"}}
query.Use(db).User.Create(&users)

for _, user := range users {
    user.ID // 1,2,3
}
```

You can specify batch size when creating with `CreateInBatches`, e.g:

```go
var users = []User{{Name: "modi_1"}, ...., {Name: "modi_10000"}}

// batch size 100
query.Use(db).User.CreateInBatches(users, 100)
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

u.Create(&users)
// INSERT INTO users xxx (5 batches)
```

#### Query

##### Retrieving a single object

Generated code provides `First`, `Take`, `Last` methods to retrieve a single object from the database, it adds `LIMIT 1` condition when querying the database, and it will return the error `ErrRecordNotFound` if no record is found.

```go
u := query.Use(db).User

// Get the first record ordered by primary key
user, err := u.First()
// SELECT * FROM users ORDER BY id LIMIT 1;

// Get one record, no specified order
user, err := u.Take()
// SELECT * FROM users LIMIT 1;

// Get last record, ordered by primary key desc
user, err := u.Last()
// SELECT * FROM users ORDER BY id DESC LIMIT 1;

// check error ErrRecordNotFound
errors.Is(err, gorm.ErrRecordNotFound)
```

##### Retrieving objects with primary key

```go
u := query.Use(db).User

user, err := u.Where(u.ID.Eq(10)).First()
// SELECT * FROM users WHERE id = 10;

users, err := u.Where(u.ID.In(1,2,3)).Find()
// SELECT * FROM users WHERE id IN (1,2,3);
```

If the primary key is a string (for example, like a uuid), the query will be written as follows:

```go
user, err := u.Where(u.ID.Eq("1b74413f-f3b8-409f-ac47-e8c062e3472a")).First()
// SELECT * FROM users WHERE id = "1b74413f-f3b8-409f-ac47-e8c062e3472a";
```

##### Retrieving all objects

```go
u := query.Use(db).User

// Get all records
users, err := u.Find()
// SELECT * FROM users;
```

##### Conditions

###### String Conditions

```go
u := query.Use(db).User

// Get first matched record
user, err := u.Where(u.Name.Eq("modi")).First()
// SELECT * FROM users WHERE name = 'modi' ORDER BY id LIMIT 1;

// Get all matched records
users, err := u.Where(u.Name.Neq("modi")).Find()
// SELECT * FROM users WHERE name <> 'modi';

// IN
users, err := u.Where(u.Name.In("modi", "zhangqiang")).Find()
// SELECT * FROM users WHERE name IN ('modi','zhangqiang');

// LIKE
users, err := u.Where(u.Name.Like("%modi%")).Find()
// SELECT * FROM users WHERE name LIKE '%modi%';

// AND
users, err := u.Where(u.Name.Eq("modi"), u.Age.Gte(17)).Find()
// SELECT * FROM users WHERE name = 'modi' AND age >= 17;

// Time
users, err := u.Where(u.Birthday.Gt(birthTime).Find()
// SELECT * FROM users WHERE birthday > '2000-01-01 00:00:00';

// BETWEEN
users, err := u.Where(u.Birthday.Between(lastWeek, today)).Find()
// SELECT * FROM users WHERE birthday BETWEEN '2000-01-01 00:00:00' AND '2000-01-08 00:00:00';
```

###### Inline Condition

```go
u := query.Use(db).User

// Get by primary key if it were a non-integer type
user, err := u.Where(u.ID.Eq("string_primary_key")).First()
// SELECT * FROM users WHERE id = 'string_primary_key';

// Plain SQL
users, err := u.Where(u.Name.Eq("modi")).Find()
// SELECT * FROM users WHERE name = "modi";

users, err := u.Where(u.Name.Neq("modi"), u.Age.Gt(17)).Find()
// SELECT * FROM users WHERE name <> "modi" AND age > 17;
```

###### Not Conditions

Build NOT conditions, works similar to `Where`

```go
u := query.Use(db).User

user, err := u.Not(u.Name.Eq("modi")).First()
// SELECT * FROM users WHERE NOT name = "modi" ORDER BY id LIMIT 1;

// Not In
users, err := u.Not(u.Name.In("modi", "zhangqiang")).Find()
// SELECT * FROM users WHERE name NOT IN ("modi", "zhangqiang");

// Not In slice of primary keys
user, err := u.Not(u.ID.In(1,2,3)).First()
// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;
```

###### Or Conditions

```go
u := query.Use(db).User

users, err := u.Where(u.Role.Eq("admin")).Or(u.Role.Eq("super_admin")).Find()
// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
```

###### Group Conditions

Easier to write complicated SQL query with Group Conditions

```go
p := query.Use(db).Pizza

pizzas, err := p.Where(
    p.Where(p.Pizza.Eq("pepperoni")).Where(p.Where(p.Size.Eq("small")).Or(p.Size.Eq("medium"))),
).Or(
    p.Where(p.Pizza.Eq("hawaiian")).Where(p.Size.Eq("xlarge")),
).Find()

// SELECT * FROM `pizzas` WHERE (pizza = "pepperoni" AND (size = "small" OR size = "medium")) OR (pizza = "hawaiian" AND size = "xlarge")
```

###### Selecting Specific Fields

`Select` allows you to specify the fields that you want to retrieve from database. Otherwise, GORM will select all fields by default.

```go
u := query.Use(db).User

users, err := u.Select(u.Name, u.Age).Find()
// SELECT name, age FROM users;

u.Select(u.Age.Avg()).Rows()
// SELECT Avg(age) FROM users;
```

###### Tuple Query

```go
u := query.Use(db).User

users, err := u.Where(u.Columns(u.ID, u.Name).In(field.Values([][]inferface{}{{1, "modi"}, {2, "zhangqiang"}}))).Find()
// SELECT * FROM `users` WHERE (`id`, `name`) IN ((1,'humodi'),(2,'tom'));
```

###### Order

Specify order when retrieving records from the database

```go
u := query.Use(db).User

users, err := u.Order(u.Age.Desc(), u.Name).Find()
// SELECT * FROM users ORDER BY age DESC, name;

// Multiple orders
users, err := u.Order(u.Age.Desc()).Order(u.Name).Find()
// SELECT * FROM users ORDER BY age DESC, name;
```

###### Limit & Offset

`Limit` specify the max number of records to retrieve
`Offset` specify the number of records to skip before starting to return the records

```go
u := query.Use(db).User

urers, err := u.Limit(3).Find()
// SELECT * FROM users LIMIT 3;

// Cancel limit condition with -1
users, err := u.Limit(10).Limit(-1).Find()
// SELECT * FROM users;

users, err := u.Offset(3).Find()
// SELECT * FROM users OFFSET 3;

users, err := u.Limit(10).Offset(5).Find()
// SELECT * FROM users OFFSET 5 LIMIT 10;

// Cancel offset condition with -1
users, err := u.Offset(10).Offset(-1).Find()
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

err := u.Select(u.Name, u.Age.Sum().As("total")).Where(u.Name.Like("%modi%")).Group(u.Name).Scan(&result)
// SELECT name, sum(age) as total FROM `users` WHERE name LIKE "group%" GROUP BY `name`

err := u.Select(u.Name, u.Age.Sum().As("total")).Group(u.Name).Having(u.Name.Eq("group")).Scan(&result)
// SELECT name, sum(age) as total FROM `users` GROUP BY `name` HAVING name = "group"

rows, err := u.Select(u.Birthday.As("date"), u.Age.Sum().As("total")).Group(u.Birthday).Rows()
for rows.Next() {
  ...
}

o := query.Use(db).Order

rows, err := o.Select(o.CreateAt.Date().As("date"), o.Amount.Sum().As("total")).Group(o.CreateAt.Date()).Having(u.Amount.Sum().Gt(100)).Rows()
for rows.Next() {
  ...
}

var results []Result

o.Select(o.CreateAt.Date().As("date"), o.Amount.Sum().As("total")).Group(o.CreateAt.Date()).Having(u.Amount.Sum().Gt(100)).Scan(&results)
```

###### Distinct

Selecting distinct values from the model

```go
u := query.Use(db).User

users, err := u.Distinct(u.Name, u.Age).Order(u.Name, u.Age.Desc()).Find()
```

`Distinct` works with `Pluck` and `Count` too

###### Joins

Specify Joins conditions

```go
u := query.Use(db).User
e := query.Use(db).Email
c := query.Use(db).CreditCard

type Result struct {
    Name  string
    Email string
}

var result Result

err := u.Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Scan(&result)
// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

rows, err := u.Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Rows()
for rows.Next() {
  ...
}

var results []Result

err := u.Select(u.Name, e.Email).LeftJoin(e, e.UserID.EqCol(u.ID)).Scan(&results)

// multiple joins with parameter
users := u.Join(e, e.UserID.EqCol(u.id), e.Email.Eq("modi@example.org")).Join(c, c.UserID.EqCol(u.ID)).Where(c.Number.Eq("411111111111")).Find()
```

##### SubQuery

A subquery can be nested within a query, GEN can generate subquery when using a `Dao` object as param

```go
o := query.Use(db).Order
u := query.Use(db).User

orders, err := o.Where(u.Columns(o.Amount).Gt(o.Select(u.Amount.Avg())).Find()
// SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

subQuery := u.Select(u.Age.Avg()).Where(u.Name.Like("name%"))
users, err := u.Select(u.Age.Avg().As("avgage")).Group(u.Name).Having(u.Columns(u.Age.Avg()).Gt(subQuery).Find()
// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")
```

###### From SubQuery

GORM allows you using subquery in FROM clause with method `Table`, for example:

```go
u := query.Use(db).User
p := query.Use(db).Pet

users, err := gen.Table(u.Select(u.Name, u.Age).As("u")).Where(u.Age.Eq(18)).Find()
// SELECT * FROM (SELECT `name`,`age` FROM `users`) as u WHERE `age` = 18

subQuery1 := u.Select(u.Name)
subQuery2 := p.Select(p.Name)
users, err := gen.Table(subQuery1.As("u"), subQuery2.As("p")).Find()
db.Table("(?) as u, (?) as p", subQuery1, subQuery2).Find(&User{})
// SELECT * FROM (SELECT `name` FROM `users`) as u, (SELECT `name` FROM `pets`) as p
```

###### Update from SubQuery

Update a table by using SubQuery

```go
u := query.Use(db).User
c := query.Use(db).Company

u.Update(u.CompanyName, c.Select(c.Name).Where(c.ID.EqCol(u.CompanyID)))
// UPDATE "users" SET "company_name" = (SELECT name FROM companies WHERE companies.id = users.company_id);

u.Where(u.Name.Eq("modi")).Update(u.CompanyName, c.Select(c.Name).Where(c.ID.EqCol(u.CompanyID)))
```

##### Advanced Query

###### Iteration

GEN supports iterating through Rows

```go
rows, err := query.Use(db).User.Where(u.Name.Eq("modi")).Rows()
defer rows.Close()

for rows.Next() {
    var user User
    // ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
    db.ScanRows(rows, &user)

    // do something
}
```

###### FindInBatches

Query and process records in batch

```go
u := query.Use(db).User

// batch size 100
err := u.Where(u.ID.Gt(9)).FindInBatches(&results, 100, func(tx gen.Dao, batch int) error {
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
u.Pluck(u.Age, &ages)

var names []string
u.Pluck(u.Name, &names)

// Distinct Pluck
u.Distinct().Pluck(u.Name, &names)
// SELECT DISTINCT `name` FROM `users`

// Requesting more than one column, use `Scan` or `Find` like this:
db.Select(u.Name, u.Age).Scan(&users)
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

orders, err := o.Scopes(AmountGreaterThan1000, PaidWithCreditCard).Find()
// Find all credit card orders and amount greater than 1000

orders, err := o.Scopes(AmountGreaterThan1000, PaidWithCod).Find()
// Find all COD orders and amount greater than 1000

orders, err := o.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find()
// Find all paid, shipped orders that amount greater than 1000
```

###### Count

Get matched records count

```go
u := query.Use(db).User

count, err := u.Where(u.Name.Eq("modi")).Or(u.Name.Eq("zhangqiang")).Count()
// SELECT count(1) FROM users WHERE name = 'modi' OR name = 'zhangqiang'

count, err := u.Where(u.Name.Eq("modi")).Count()
// SELECT count(1) FROM users WHERE name = 'modi'; (count)

// Count with Distinct
u.Distinct(u.Name).Count()
// SELECT COUNT(DISTINCT(`name`)) FROM `users`
```

#### Update

##### Update single column

When updating a single column with `Update`, it needs to have any conditions or it will raise error `ErrMissingWhereClause`, for example:

```go
u := query.Use(db).User

// Update with conditions
u.Where(u.Activate.Is(true)).Update(u.Name, "hello")
// UPDATE users SET name='hello', updated_at='2013-11-17 21:34:10' WHERE active=true;

// Update with conditions
u.Where(u.Activate.Is(true)).Update(u.Age, u.Age.Add(1))
// or
u.Where(u.Activate.Is(true)).UpdateSimple(u.Age.Add(1))
// UPDATE users SET age=age+1, updated_at='2013-11-17 21:34:10' WHERE active=true;
```

##### Updates multiple columns

`Updates` supports update with `struct` or `map[string]interface{}`, when updating with `struct` it will only update non-zero fields by default

```go
u := query.Use(db).User

// Update attributes with `map`
u.Model(&model.User{ID: 111}).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello', age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```

> **NOTE** When update with struct, GEN will only update non-zero fields, you might want to use `map` to update attributes or use `Select` to specify fields to update

##### Update selected fields

If you want to update selected fields or ignore some fields when updating, you can use `Select`, `Omit`

```go
u := query.Use(db).User

// Select with Map
// User's ID is `111`:
u.Select(u.Name).Where(u.ID.Eq(111)).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET name='hello' WHERE id=111;

u.Omit(u.Name).Where(u.ID.Eq(111)).Updates(map[string]interface{}{"name": "hello", "age": 18, "active": false})
// UPDATE users SET age=18, active=false, updated_at='2013-11-17 21:34:10' WHERE id=111;
```

#### Delete

##### Delete record

```go
e := query.Use(db).Email

// Email's ID is `10`
e.Where(e.ID.Eq(10)).Delete()
// DELETE from emails where id = 10;

// Delete with additional conditions
e.Where(e.ID.Eq(10), e.Name.Eq("modi")).Delete()
// DELETE from emails where id = 10 AND name = "modi";
```

##### Delete with primary key

GEN allows to delete objects using primary key(s) with inline condition, it works with numbers.

```go
u.Where(u.ID.In(1,2,3)).Delete()
// DELETE FROM users WHERE id IN (1,2,3);
```

##### Batch Delete

The specified value has no primary value, GEN will perform a batch delete, it will delete all matched records

```go
e := query.Use(db).Email

err := e.Where(e.Name.Like("%modi%")).Delete()
// DELETE from emails where email LIKE "%modi%";
```

##### Soft Delete

If your model includes a `gorm.DeletedAt` field (which is included in `gorm.Model`), it will get soft delete ability automatically!

When calling `Delete`, the record WON’T be removed from the database, but GORM will set the `DeletedAt`‘s value to the current time, and the data is not findable with normal Query methods anymore.

```go
// Batch Delete
err := u.Where(u.Age.Eq(20)).Delete()
// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

// Soft deleted records will be ignored when querying
users, err := u.Where(u.Age.Eq(20)).Find()
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
users, err := db.Unscoped().Where(u.Age.Eq(20)).Find()
// SELECT * FROM users WHERE age = 20;
```

##### Delete permanently

You can delete matched records permanently with `Unscoped`

```go
o.Unscoped().Where(o.ID.Eq(10)).Delete()
// DELETE FROM orders WHERE id=10;
```

### DIY method

#### Method interface

Method interface is an abstraction of query methods, all functions it contains are query methods and above comments describe the specific query conditions or logic.
SQL supports simple `where` query or execute raw SQL. Simple query conditions wrapped by `where()`, and raw SQL wrapped by `sql()`（not required）

```go
type Method interface {
    // where("name=@name and age=@age")
    SimpleFindByNameAndAge(name string, age int) (gen.T, error)
    
    // sql(select * from users where id=@id)
    FindUserToMap(id int) (gen.M, error)
    
    // insert into users (name,age) values (@name,@age)
    InsertValue(age int, name string) error
}
```

Return values must contains less than 1 `gen.T`/`gen.M` and less than 1 error. You can also use bulitin type (like `string`/ `int`) as the return parameter，`gen.T` represents return a single result struct's pointer, `[]gen.T` represents return an array of result structs' pointer，

##### Syntax of template

###### placeholder

- `gen.T` represents specified `struct` or `table`
- `gen.M` represents `map[string]interface`
- `@@table`  represents table's name (if method's parameter doesn't contains variable `table`, GEN will generate `table` from model struct)
- `@@<columnName>` represents column's name or table's name
- `@<name>` represents normal query variable

###### template

Logical operations must be wrapped in `{{}}`,and end must used `{{end}}`, All templates support nesting

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
    UpdateName(name string, id int) error
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

apiusers, err := u.Limit(10).FindSome()
// SELECT `id`, `name` FROM `users` LIMIT 10
```

### Advanced Topics

#### Hints

Optimizer hints allow to control the query optimizer to choose a certain query execution plan, GORM supports it with `gorm.io/hints`, e.g:

```go
import "gorm.io/hints"

u := query.Use(db).User

users, err := u.Hints(hints.New("MAX_EXECUTION_TIME(10000)")).Find()
// SELECT * /*+ MAX_EXECUTION_TIME(10000) */ FROM `users`
```

Index hints allow passing index hints to the database in case the query planner gets confused.

```go
import "gorm.io/hints"

u := query.Use(db).User

users, err := u.Hints(hints.UseIndex("idx_user_name")).Find()
// SELECT * FROM `users` USE INDEX (`idx_user_name`)

users, err := u.Hints(hints.ForceIndex("idx_user_name", "idx_user_id").ForJoin()).Find()
// SELECT * FROM `users` FORCE INDEX FOR JOIN (`idx_user_name`,`idx_user_id`)"
```

## Maintainers

[@riverchu](https://github.com/riverchu) [@idersec](https://github.com/idersec) [@qqxhb](https://github.com/qqxhb)

[@jinzhu](https://github.com/jinzhu)

## Contributing

You can help to deliver a better GORM/GEN

## License

Released under the [MIT License](https://github.com/go-gorm/gen/blob/master/License)
