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

#### db

default:mysql

input mysql or postgres or sqlite or sqlserver.

consult : https://gorm.io/docs/connecting_to_the_database.html

#### dsn

You can use all gorm's dsn.

 consult : https://gorm.io/docs/connecting_to_the_database.html

#### fieldNullable

generate with pointer when field is nullable

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



### example

```shell
gentool -dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" -tables "orders,doctor"
```