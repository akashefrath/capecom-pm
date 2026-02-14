package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

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

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	FolderName      string
}

type Config struct {
	Port         string
	DB           DBConfig
	JWT          JWTConfig
	RedisAddress string
	R2           R2Config
}

func LoadEnv() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}
	return Config{
		Port:         GetEnvMust("APP_PORT"),
		DB:           loadDBConfig(),
		JWT:          loadJWTConfig(),
		RedisAddress: GetEnvMust("REDIS_ADDRESS"),
		R2:           loadR2Config(),
	}

}

func loadDBConfig() DBConfig {
	maxIdle, _ := strconv.Atoi(GetEnv("DB_MAX_IDLE"))
	maxOpen, _ := strconv.Atoi(GetEnv("DB_MAX_OPEN"))
	life, _ := strconv.Atoi(GetEnv("DB_MAX_LIFETIME"))

	return DBConfig{
		DBUser:        GetEnvMust("DB_USER"),
		DBPass:        GetEnv("DB_PASS"),
		DBHost:        GetEnvMust("DB_HOST"),
		DBPort:        GetEnvMust("DB_PORT"),
		DBName:        GetEnvMust("DB_NAME"),
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

func loadR2Config() R2Config {
	return R2Config{
		AccountID:       GetEnvMust("R2_ACCOUNT_ID"),
		AccessKeyID:     GetEnvMust("R2_ACCESS_KEY_ID"),
		AccessKeySecret: GetEnvMust("R2_ACCESS_KEY_SECRET"),
		BucketName:      GetEnvMust("R2_BUCKET_NAME"),
		FolderName:      GetEnv("R2_FOLDER_NAME"),
	}
}
