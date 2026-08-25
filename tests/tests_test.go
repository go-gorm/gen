package tests_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	mysqlDSN     = "gen:gen@tcp(localhost:9910)/gen?charset=utf8&parseTime=True&loc=Local"
	postgresDSN  = "user=gen password=gen dbname=gen host=localhost port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	sqlserverDSN = "sqlserver://gen:LoremIpsum86@localhost:9930?database=gen"
)

var DB *gorm.DB

func TestMain(m *testing.M) {
	cleanup, err := setupMySQLContract()
	if err != nil {
		log.Printf("mysql contract setup failed: %v", err)
		os.Exit(1)
	}

	code := m.Run()
	if err := cleanup(); err != nil {
		log.Printf("mysql contract cleanup failed: %v", err)
		code = 1
	}
	os.Exit(code)
}

func setupMySQLContract() (func() error, error) {
	if os.Getenv("GORM_DIALECT") != "mysql" {
		return nil, fmt.Errorf("tests root package requires GORM_DIALECT=mysql")
	}

	var err error
	DB, err = OpenTestConnection()
	if err != nil {
		return nil, err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	cleanup := func() error {
		removeErr := os.RemoveAll(generateDirPrefix)
		closeErr := sqlDB.Close()
		if removeErr != nil {
			return fmt.Errorf("remove generated test output: %w", removeErr)
		}
		if closeErr != nil {
			return fmt.Errorf("close MySQL connection: %w", closeErr)
		}
		return nil
	}
	if err := sqlDB.Ping(); err != nil {
		_ = cleanup()
		return nil, fmt.Errorf("ping MySQL: %w", err)
	}
	if err := os.RemoveAll(generateDirPrefix); err != nil {
		_ = cleanup()
		return nil, fmt.Errorf("remove stale generated test output: %w", err)
	}
	if err := RunMigrations(); err != nil {
		_ = cleanup()
		return nil, err
	}

	var generators []*gen.Generator
	for dir, build := range generateCase {
		generators = append(generators, build(dir))
	}
	if err := RunGenerate(generators...); err != nil {
		_ = cleanup()
		return nil, err
	}
	return cleanup, nil
}

func OpenTestConnection() (db *gorm.DB, err error) {
	dbDSN := os.Getenv("GEN_DSN")
	if os.Getenv("GORM_DIALECT") != "mysql" {
		return nil, fmt.Errorf("unsupported contract dialect %q", os.Getenv("GORM_DIALECT"))
	}
	log.Println("testing mysql...")
	if dbDSN == "" {
		dbDSN = mysqlDSN
	}
	db, err = gorm.Open(mysql.Open(dbDSN), &gorm.Config{})

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

func RunMigrations() error {
	db := DB.Session(&gorm.Session{})
	ddl, err := GetDDL()
	if err != nil {
		return err
	}
	for _, meta := range ddl {
		dropTable, createTable := meta[0], meta[1]
		if err := db.Exec(dropTable).Error; err != nil {
			return fmt.Errorf("drop table: %w", err)
		}
		if err := db.Exec(createTable).Error; err != nil {
			return fmt.Errorf("create table: %w", err)
		}
	}
	return nil
}

func RunGenerate(gs ...*gen.Generator) (err error) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			err = fmt.Errorf("generate fixtures: %v", panicValue)
		}
	}()
	for _, g := range gs {
		g.Execute()
	}
	return nil
}
