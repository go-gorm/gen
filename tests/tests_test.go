package tests_test

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	dsn := os.Getenv("GORM_DSN")
	if dsn == "" {
		dsn = "gen:gen@tcp(localhost:9910)/gen?charset=utf8&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("open mysql fail: %w", err))
	}
	_ = db
}
