# GenTool

Install GEN as a binary tool

## install

```shell
 go install gorm.io/gen/tools/gentool@latest
```

## usage

```shell
gentool -h  

Usage of gentool:
  -c string
        is path for gen.yml
  -db string
        input mysql|postgres|sqlite|sqlserver|clickhouse. consult[https://gorm.io/docs/connecting_to_the_database.html] (default "mysql")
  -dsn string
        consult[https://gorm.io/docs/connecting_to_the_database.html]
  -fieldCoverable
        generate with pointer when field has default value
  -fieldNullable
        generate with pointer when field is nullable
  -fieldSignable
        detect integer field's unsigned type, adjust generated data type
  -fieldWithIndexTag
        generate field with gorm index tag
  -fieldWithTypeTag
        generate field with gorm column type tag
  -modelPkgName string
        generated model code's package name
  -onlyModel
        only generate models (without query file)
  -outFile string
        query code file name, default: gen.go
  -outPath string
        specify a directory for output (default "./dao/query")
  -tables string
        enter the required data table or leave it blank
  -withDefaultQuery
        create default query in generated code
  -withQueryInterface
        generate code with exported interface object
  -withUnitTest
        generate unit test for query code
  -withoutContext
        generate code without context constrain
```

#### c

default ""
Is path for gen.yml
Replace the command line with a configuration file
The command line is the highest priority

#### db

default:mysql

input mysql or postgres or sqlite or sqlserver.

consult : https://gorm.io/docs/connecting_to_the_database.html

#### dsn

You can use all gorm's dsn.

 consult : https://gorm.io/docs/connecting_to_the_database.html

#### fieldNullable

generate with pointer when field is nullable

#### fieldCoverable

generate with pointer when field has default value

#### fieldWithIndexTag

generate field with gorm index tag

#### fieldWithTypeTag

generate field with gorm column type tag

#### modelPkgName

defalut table name.

 generated model code's package name.

#### outFile

 query code file name, default: gen.go

#### outPath

specify a directory for output (default "./dao/query")

#### tables

Value : enter the required data table or leave it blank.

eg :

​       --tables="orders" #orders table

​       --tables="orders,users" #orders table and users table

​       --tables=""          # All data tables in the database.

Generate some tables code.

#### withUnitTest

Value : False / True

Generate unit test.

#### fieldSignable

Value : False / True

detect integer field's unsigned type, adjust generated data type

#### withDefaultQuery

Default: false

create default query in generated code

#### withQueryInterface

Default: false

generate code with exported interface object

#### withoutContext

Default: false

generate code without context constrain

### example

```shell
gentool -dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" -tables "orders,doctor"
```