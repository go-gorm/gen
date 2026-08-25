module gorm.io/gen/examples

go 1.22.0

toolchain go1.24.4

require (
	gorm.io/driver/mysql v1.5.7
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gen v0.3.25
	gorm.io/gorm v1.30.0
)

require (
	filippo.io/edwards25519 v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/exp v0.0.0-20240112132812-db7319d0e0e3 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	gorm.io/datatypes v1.2.4 // indirect
	gorm.io/hints v1.1.0 // indirect
	gorm.io/plugin/dbresolver v1.5.3 // indirect
)

replace gorm.io/gen => ../
