module gorm.io/gen/tests

go 1.16

require (
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	gorm.io/driver/mysql v1.5.4
	gorm.io/driver/sqlite v1.4.4
	gorm.io/gen v0.3.19
	gorm.io/gorm v1.25.7-0.20240204074919-46816ad31dde
	gorm.io/hints v1.1.1 // indirect
	gorm.io/plugin/dbresolver v1.5.0
)

replace gorm.io/gen => ../
