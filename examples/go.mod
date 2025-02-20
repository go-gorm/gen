module examples

go 1.19

require (
	gorm.io/driver/mysql v1.5.6
	gorm.io/driver/sqlite v1.5.5
	gorm.io/gen v0.3.25
	gorm.io/gorm v1.25.11
)

require (
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	gorm.io/datatypes v1.2.4 // indirect
	gorm.io/hints v1.1.0 // indirect
	gorm.io/plugin/dbresolver v1.5.0 // indirect
)

replace gorm.io/gen => ../
