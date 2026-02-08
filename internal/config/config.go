package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

type Config struct {
	Port string
	DB   DBConfig
}

func LoadEnv() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}
	return Config{
		Port: GetEnvMust("APP_PORT"),
		DB:   loadDBConfig(),
	}

}

func loadDBConfig() DBConfig {
	return DBConfig{
		DBUser: GetEnvMust("DB_USER"),
		DBPass: GetEnv("DB_PASS"),
		DBHost: GetEnvMust("DB_HOST"),
		DBPort: GetEnvMust("DB_PORT"),
		DBName: GetEnvMust("DB_NAME"),
	}

}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvMust(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal("env var " + key + " is not set")
	}
	return value
}
