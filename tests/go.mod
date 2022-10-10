module gorm.io/gen/tests

go 1.16

require (
	gorm.io/driver/mysql v1.4.0
	gorm.io/driver/sqlite v1.4.1
	gorm.io/gen v0.3.16
	gorm.io/gorm v1.24.0
	gorm.io/plugin/dbresolver v1.3.0
)

replace gorm.io/gen => ../
