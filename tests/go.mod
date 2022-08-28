module gorm.io/gen/tests

go 1.16

require (
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	golang.org/x/sys v0.0.0-20220825204002-c680a09ffe64 // indirect
	gorm.io/driver/mysql v1.3.6
	gorm.io/driver/sqlite v1.3.6
	gorm.io/gen v0.3.15
	gorm.io/gorm v1.23.9-0.20220713102635-3262daf8d468
	gorm.io/plugin/dbresolver v1.2.3 // indirect
)

replace gorm.io/gen => ../
