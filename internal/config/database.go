package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(dbConfig DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbConfig.DBUser,
		dbConfig.DBPass,
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	println("opened")

	sqlDB.SetMaxIdleConns(dbConfig.DBMaxIdle)
	sqlDB.SetMaxOpenConns(dbConfig.DBMaxOpen)
	sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(dbConfig.DBMaxLifetime))

	return db
}
