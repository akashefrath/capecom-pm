package config

import (
	"fmt"
	"log"

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

	//maxIdle, _ := strconv.Atoi(GetEnv("DB_MAX_IDLE"))
	//maxOpen, _ := strconv.Atoi(GetEnv("DB_MAX_OPEN"))
	//life, _ := strconv.Atoi(GetEnv("DB_MAX_LIFETIME"))
	//
	//sqlDB.SetMaxIdleConns(maxIdle)
	//sqlDB.SetMaxOpenConns(maxOpen)
	//sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(life))

	return db
}
