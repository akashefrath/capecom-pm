package config

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitDB(config Config) *sqlx.DB {
	dbConfig := config.DBConfig

	cfg := mysql.Config{
		User:                 dbConfig.DBUser,
		Passwd:               dbConfig.DBPass,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", dbConfig.DBHost, dbConfig.DBPort),
		DBName:               dbConfig.DBName,
		ParseTime:            true,
		AllowNativePasswords: true,
		Loc:                  time.UTC,
		Params: map[string]string{
			"time_zone": "'+00:00'",
		},
	}

	dbUrl := cfg.FormatDSN()

	// 1. Capture the error! Don't use '_'
	db, err := sqlx.Connect("mysql", dbUrl)
	if err != nil {
		panic(fmt.Sprintf("Invalid connection string: %v", err))
	}

	// 2. Now it is safe to Ping
	err = db.Ping()
	_, err = db.Exec("SET time_zone = '+00:00'")

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
