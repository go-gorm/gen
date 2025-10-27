package main

import (
	"gorm.io/gorm"
)

// prepare table for test

const mytableSQL = "CREATE TABLE IF NOT EXISTS `mytables` (" +
	"    `ID` INTEGER NOT NULL PRIMARY KEY," +
	"    `username` TEXT," +
	"    `age` INTEGER NOT NULL," +
	"    `phone` TEXT NOT NULL" +
	");"

const indexSQL = "CREATE INDEX IF NOT EXISTS `idx_username` ON `mytables` (`username`);"

func prepare(db *gorm.DB) {
	db.Exec(mytableSQL)
	db.Exec(indexSQL)
}
