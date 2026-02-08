package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

type JWTConfig struct {
	UserSecret         string
	AdminSecret        string
	ExpireHours        int
	RefreshExpireHours int
}

type Config struct {
	Port string
	DB   DBConfig
	JWT  JWTConfig
}

func LoadEnv() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}
	return Config{
		Port: GetEnvMust("APP_PORT"),
		DB:   loadDBConfig(),
		JWT:  loadJWTConfig(),
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

func loadJWTConfig() JWTConfig {
	expireHours, err := strconv.Atoi(GetEnvMust("JWT_EXPIRE_HOURS"))
	if err != nil {
		log.Fatal("JWT_EXPIRE_HOURS must be a number")
	}

	refreshExpireHours, err := strconv.Atoi(GetEnvMust("JWT_REFRESH_EXPIRE_HOURS"))
	if err != nil {
		log.Fatal("JWT_REFRESH_EXPIRE_HOURS must be a number")
	}

	return JWTConfig{
		UserSecret:         GetEnvMust("JWT_SECRET"),
		AdminSecret:        GetEnvMust("JWT_ADMIN_SECRET"),
		ExpireHours:        expireHours,
		RefreshExpireHours: refreshExpireHours,
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
