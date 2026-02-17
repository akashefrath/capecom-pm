package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         int
	DBConfig     DBConfig
	JWT          JWTConfig
	RedisAddress string
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
type JWTConfig struct {
	UserSecret         string
	AdminSecret        string
	UserRefreshSecret  string
	AdminRefreshSecret string
	ExpireHours        int
	RefreshExpireHours int
}

func LoadEnv() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}
	port, _ := strconv.Atoi(GetEnv("APP_PORT"))
	return Config{
		Port:         port,
		DBConfig:     LoadDBConfig(),
		JWT:          loadJWTConfig(),
		RedisAddress: GetEnv("REDIS_ADDRESS"),
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
		UserRefreshSecret:  GetEnvMust("JWT_REFRESH_SECRET"),
		AdminSecret:        GetEnvMust("JWT_ADMIN_SECRET"),
		AdminRefreshSecret: GetEnvMust("JWT_ADMIN_REFRESH_SECRET"),
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
