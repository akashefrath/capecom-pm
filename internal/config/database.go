package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(config Config) *sql.DB {
	dbConfig := config.DBConfig
	// Ensure this string is exactly: user:pass@tcp(host:port)/dbname
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.DBUser,
		dbConfig.DBPass,
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBName,
	)

	// 1. Capture the error! Don't use '_'
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		panic(fmt.Sprintf("Invalid connection string: %v", err))
	}

	// 2. Now it is safe to Ping
	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Could not connect to DB: %v", err))
	} else {
		fmt.Println("Successfully connected to DB")
	}

	if db == nil {

		panic("failed to connect to database")

	}
	// Optional: Set pool limits here
	db.SetMaxOpenConns(25)

	return db
}
