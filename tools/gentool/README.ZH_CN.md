# GenTool

将Gen作为二进制的方式进行安装



## 安装

```shell
 go install gorm.io/gen/tools/gentool@latest
```

## 使用方式

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
  -onlyModel
        only generate models (without query file)
  -withUnitTest
        generate unit test for query code
  -fieldSignable
        detect integer field's unsigned type, adjust generated data type

```

#### c
default ""
可以指定配置文件gen.yml的路径。
用配置文件来代替命令行。
命令行是最高优先级。

#### db

默认值：mysql

可以输入： mysql、 postgres、 sqlite 、 sqlserve

参考：https://gorm.io/docs/connecting_to_the_database.html

#### dsn

你可以使用GORM所有的连接。

 参考：https://gorm.io/docs/connecting_to_the_database.html

#### fieldNullable

字段可为空时使用指针生成

#### fieldWithIndexTag

使用GROM索引标记生成字段

#### fieldWithTypeTag

使用gorm列类型标记生成字段

#### modelPkgName

默认值是数据表名称。

 生成的model代码的包名称。

#### outFile

默认为：gen.go

查询代码文件名。

#### outPath

默认为：/dao/query

指定输出目录

#### tables

值为 : 输入所需的数据表或将其留空

eg :

​       --tables="orders" #orders 数据表

​       --tables="orders,users" #orders 数据表和 users数据表

​       --tables=""          # 数据库中所有的数据表

基于数据表生成对应的代码。

#### withUnitTest

值为 : False / True

生成单元测试。

#### fieldSignable

Value : False / True

基于数据表定义的数据类型，生成对应的数据类型


### 使用示例

```shell
gentool -dsn "user:pwd@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" -tables "orders,doctor"
```