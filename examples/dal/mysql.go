package dal

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(dsn string) (db *gorm.DB) {
	var err error

	if strings.HasSuffix(dsn, "sqlite.db") {
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	} else {
		db, err = gorm.Open(mysql.Open(dsn))
	}

	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}

	return db
}
