package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// 连接到 PostgreSQL 数据库
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 创建测试表和枚举类型
	setupTestDatabase(db)

	// 初始化 Generator
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./dal",
		ModelPkgPath: "./dal/model",
	})

	// 使用数据库连接
	g.UseDB(db)

	// 生成所有表的模型
	g.GenerateAllTable()

	// 执行代码生成
	g.Execute()

	fmt.Println("代码生成完成，请检查 ./dal/model 目录下的文件是否包含枚举常量")
}

func setupTestDatabase(db *gorm.DB) {
	// 创建枚举类型
	db.Exec("DROP TYPE IF EXISTS user_status CASCADE")
	db.Exec("CREATE TYPE user_status AS ENUM ('active', 'inactive', 'pending')")

	// 创建测试表
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec(`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		status user_status NOT NULL DEFAULT 'pending'
	)`)

	// 插入测试数据
	db.Exec("INSERT INTO users (name, status) VALUES ('User 1', 'active')")
	db.Exec("INSERT INTO users (name, status) VALUES ('User 2', 'inactive')")
	db.Exec("INSERT INTO users (name, status) VALUES ('User 3', 'pending')")

	fmt.Println("测试数据库设置完成")
}