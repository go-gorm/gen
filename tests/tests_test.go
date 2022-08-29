package tests_test

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var (
	mysqlDSN     = "gen:gen@tcp(localhost:9910)/gen?charset=utf8&parseTime=True&loc=Local"
	postgresDSN  = "user=gen password=gen dbname=gen host=localhost port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	sqlserverDSN = "sqlserver://gen:LoremIpsum86@localhost:9930?database=gen"
)

func init() {
	log.Print("initing...")
	var err error
	if DB, err = OpenTestConnection(); err != nil {
		log.Printf("failed to connect database, got error %v", err)
		os.Exit(1)
	} else {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("failed to connect database, got error %v", err)
			os.Exit(1)
		}

		err = sqlDB.Ping()
		if err != nil {
			log.Printf("failed to ping sqlDB, got error %v", err)
			os.Exit(1)
		}

		// RunMigrations()
		if DB.Dialector.Name() == "sqlite" {
			DB.Exec("PRAGMA foreign_keys = ON")
		}
	}
	RunMigrations()
}

func OpenTestConnection() (db *gorm.DB, err error) {
	dbDSN := os.Getenv("GEN_DSN")
	switch os.Getenv("GORM_DIALECT") {
	case "mysql":
		log.Println("testing mysql...")
		if dbDSN == "" {
			dbDSN = mysqlDSN
		}
		db, err = gorm.Open(mysql.Open(dbDSN), &gorm.Config{})
	default:
		log.Println("testing sqlite3...")
		db, err = gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	}

	if err != nil {
		return
	}

	if debug := os.Getenv("DEBUG"); debug == "true" {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else if debug == "false" {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}

	return
}

func RunMigrations() {
	db := DB.Session(&gorm.Session{})
	for _, meta := range GetDDL() {
		dropTable, createTable := meta[0], meta[1]
		if err := db.Exec(dropTable).Error; err != nil {
			log.Printf("drop table fail: %s", err)
		}
		if err := db.Exec(createTable).Error; err != nil {
			log.Printf("create table fail: %s", err)
		}
	}
}
