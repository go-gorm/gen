module gorm.io/gen

go 1.14

require (
	github.com/jackc/pgx/v4 v4.16.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.0.0-20220222200937-f2425489ef4c // indirect
	golang.org/x/tools v0.1.11
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/datatypes v1.0.7
	gorm.io/driver/mysql v1.3.2
	gorm.io/driver/postgres v1.3.4
	gorm.io/driver/sqlite v1.3.5
	gorm.io/driver/sqlserver v1.3.1
	gorm.io/gorm v1.23.7
	gorm.io/hints v1.1.0
	gorm.io/plugin/dbresolver v1.2.1
)

// replace golang.org/x/net => golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f
replace gorm.io/gorm => github.com/go-gorm/gorm v1.23.7-0.20220617030057-1305f637f834
