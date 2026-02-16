package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	DBConfig DBConfig
}
type DBConfig struct {
	DBUser        string
	DBPass        string
	DBHost        string
	DBPort        string
	DBName        string
	DBMaxIdle     int
	DBMaxOpen     int
	DBMaxLifetime int
}

func LoadEnv() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}
	port, _ := strconv.Atoi(GetEnv("APP_PORT"))
	return Config{
		Port:     port,
		DBConfig: LoadDBConfig(),
	}
}

func LoadDBConfig() DBConfig {
	maxIdle, _ := strconv.Atoi(GetEnv("DB_MAX_IDLE"))
	maxOpen, _ := strconv.Atoi(GetEnv("DB_MAX_OPEN"))
	life, _ := strconv.Atoi(GetEnv("DB_MAX_LIFETIME"))

	return DBConfig{

		DBUser:        os.Getenv("DB_USER"),
		DBPass:        os.Getenv("DB_PASS"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBName:        os.Getenv("DB_NAME"),
		DBMaxIdle:     maxIdle,
		DBMaxOpen:     maxOpen,
		DBMaxLifetime: life,
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
