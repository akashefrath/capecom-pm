package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(dbConfig DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		dbConfig.DBUser,
		dbConfig.DBPass,
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBName,
	)

	// 1. Use a more robust Logger to see slow queries caused by latency
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),

		PrepareStmt: true, // Cache prepared statements to reduce round trips
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Optimized Pool Settings
	// Keep more idle connections than usual because opening new ones over 200ms is slow.
	sqlDB.SetMaxIdleConns(dbConfig.DBMaxIdle)
	sqlDB.SetMaxOpenConns(dbConfig.DBMaxOpen)

	// Prevents "Connection Reset by Peer" common in long-distance networking
	sqlDB.SetConnMaxLifetime(time.Minute * 30)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	return db
}
