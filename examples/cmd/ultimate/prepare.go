package main

import (
	"gorm.io/gorm"
)

// prepare table for test

const mytableSQL = "CREATE TABLE IF NOT EXISTS `mytables` (" +
	"    `ID` int(11) NOT NULL," +
	"    `username` varchar(16) DEFAULT NULL," +
	"    `age` int(8) NOT NULL," +
	"    `phone` varchar(11) NOT NULL," +
	"    INDEX `idx_username` (`username`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;"

func prepare(db *gorm.DB) {
	db.Exec(mytableSQL)
}
