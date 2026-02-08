---
inclusion: auto
---

# Configuration Pattern

This project uses environment variables with structured config loading. All configuration is centralized in `internal/config/` and loaded at application startup.

## Configuration Architecture

```
.env file
    ↓
config.LoadEnv() - Loads all env vars into structs
    ↓
main.go - Initializes app with config
    ↓
Pass config values to needed components
```

---

## Configuration Structure

### Main Config (`internal/config/config.go`)

```go
package config

import (
    "log"
    "os"
    "github.com/joho/godotenv"
)

// Main config struct - add new config groups here
type Config struct {
    Port string
    DB   DBConfig
    JWT  JWTConfig    // Example: add new config groups
    AWS  AWSConfig    // Example: add new config groups
}

// Database configuration
type DBConfig struct {
    DBUser string
    DBPass string
    DBHost string
    DBPort string
    DBName string
}

// Load all environment variables
func LoadEnv() Config {
    err := godotenv.Load()
    if err != nil {
        log.Println(".env not found, using system env")
    }
    return Config{
        Port: GetEnvMust("APP_PORT"),
        DB:   loadDBConfig(),
        // JWT:  loadJWTConfig(),    // Add new loaders
        // AWS:  loadAWSConfig(),    // Add new loaders
    }
}

// Load database config
func loadDBConfig() DBConfig {
    return DBConfig{
        DBUser: GetEnvMust("DB_USER"),
        DBPass: GetEnv("DB_PASS"),        // Optional
        DBHost: GetEnvMust("DB_HOST"),
        DBPort: GetEnvMust("DB_PORT"),
        DBName: GetEnvMust("DB_NAME"),
    }
}

// Get optional env var
func GetEnv(key string) string {
    return os.Getenv(key)
}

// Get required env var (fails if not set)
func GetEnvMust(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatal("env var " + key + " is not set")
    }
    return value
}
```

### Database Connection (`internal/config/database.go`)

```go
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
    
    println("Database connected")

    // Optional: Configure connection pool
    // sqlDB.SetMaxIdleConns(10)
    // sqlDB.SetMaxOpenConns(100)
    // sqlDB.SetConnMaxLifetime(time.Minute * 60)

    return db
}
```

### Main Application (`cmd/main.go`)

```go
package main

import (
    "capecom-pm/internal/config"
    "capecom-pm/internal/container"
    "capecom-pm/internal/routes"
    "fmt"
    "log"
    "github.com/gin-gonic/gin"
)

func main() {
    // Step 1: Load configuration
    appConfig := config.LoadEnv()
    
    // Step 2: Connect to database
    db := config.ConnectDB(appConfig.DB)
    if db == nil {
        log.Fatal("failed to connect to database")
    }
    
    // Step 3: Initialize DI container
    c := container.NewContainer(db)
    
    // Step 4: Setup router
    r := gin.Default()
    routes.Setup(r, c)
    
    // Step 5: Start server
    fmt.Println("http://localhost:" + appConfig.Port)
    err := r.Run(":" + appConfig.Port)
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## Environment Variables (.env)

```env
# Application
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=
DB_NAME=pm_db

# Database Connection Pool (optional)
DB_MAX_IDLE=10
DB_MAX_OPEN=100
DB_MAX_LIFETIME=60

# JWT (example)
JWT_SECRET=supersecretkey
JWT_EXPIRE_HOURS=24

# AWS (example)
AWS_REGION=us-east-1
AWS_ACCESS_KEY=your-access-key
AWS_SECRET_KEY=your-secret-key

# Email (example)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-password
```

---

## How to Add New Configuration

### Step 1: Add to .env

```env
# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASS=
REDIS_DB=0
```

### Step 2: Create Config Struct

**File:** `internal/config/config.go`

```go
// Add new config struct
type RedisConfig struct {
    Host string
    Port string
    Pass string
    DB   string
}

// Add to main Config struct
type Config struct {
    Port  string
    DB    DBConfig
    Redis RedisConfig  // Add this
}

// Create loader function
func loadRedisConfig() RedisConfig {
    return RedisConfig{
        Host: GetEnvMust("REDIS_HOST"),
        Port: GetEnvMust("REDIS_PORT"),
        Pass: GetEnv("REDIS_PASS"),      // Optional
        DB:   GetEnv("REDIS_DB"),        // Optional
    }
}

// Add to LoadEnv()
func LoadEnv() Config {
    err := godotenv.Load()
    if err != nil {
        log.Println(".env not found, using system env")
    }
    return Config{
        Port:  GetEnvMust("APP_PORT"),
        DB:    loadDBConfig(),
        Redis: loadRedisConfig(),  // Add this
    }
}
```

### Step 3: Create Connection Function (if needed)

**File:** `internal/config/redis.go`

```go
package config

import (
    "context"
    "log"
    "github.com/redis/go-redis/v9"
)

func ConnectRedis(redisConfig RedisConfig) *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     redisConfig.Host + ":" + redisConfig.Port,
        Password: redisConfig.Pass,
        DB:       0,
    })

    ctx := context.Background()
    _, err := client.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }

    println("Redis connected")
    return client
}
```

### Step 4: Use in main.go

```go
func main() {
    appConfig := config.LoadEnv()
    
    db := config.ConnectDB(appConfig.DB)
    if db == nil {
        log.Fatal("failed to connect to database")
    }
    
    // Add new connection
    redisClient := config.ConnectRedis(appConfig.Redis)
    
    // Pass to container if needed
    c := container.NewContainer(db, redisClient)
    
    // ... rest of setup
}
```

### Step 5: Pass to Services (if needed)

If a service needs config values, pass them through the container:

```go
// internal/container/service.go
func NewService(db *gorm.DB, repository *Repository, jwtSecret string) *Service {
    return &Service{
        AuthService: services.NewAuthService(repository.AuthRepo, jwtSecret),
    }
}

// internal/services/auth.go
type AuthService struct {
    repo      *repositories.AuthRepo
    jwtSecret string
}

func NewAuthService(repo *repositories.AuthRepo, jwtSecret string) *AuthService {
    return &AuthService{
        repo:      repo,
        jwtSecret: jwtSecret,
    }
}
```

---

## Configuration Best Practices

### 1. Separate Config by Domain

Group related config values into separate structs:

```go
type Config struct {
    Port  string
    DB    DBConfig      // Database
    JWT   JWTConfig     // Authentication
    AWS   AWSConfig     // Cloud services
    SMTP  SMTPConfig    // Email
    Redis RedisConfig   // Cache
}
```

### 2. Use GetEnvMust for Required Values

```go
// Application will fail to start if not set
Port: GetEnvMust("APP_PORT"),
DBHost: GetEnvMust("DB_HOST"),
```

### 3. Use GetEnv for Optional Values

```go
// Application will continue if not set
DBPass: GetEnv("DB_PASS"),
RedisPass: GetEnv("REDIS_PASS"),
```

### 4. Validate Config Values

```go
func loadJWTConfig() JWTConfig {
    secret := GetEnvMust("JWT_SECRET")
    if len(secret) < 32 {
        log.Fatal("JWT_SECRET must be at least 32 characters")
    }
    
    return JWTConfig{
        Secret:      secret,
        ExpireHours: GetEnvMust("JWT_EXPIRE_HOURS"),
    }
}
```

### 5. Use Default Values

```go
func GetEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

// Usage
Port: GetEnvWithDefault("APP_PORT", "8080"),
```

### 6. Convert Types When Needed

```go
import "strconv"

func loadDBConfig() DBConfig {
    maxIdle, _ := strconv.Atoi(GetEnvWithDefault("DB_MAX_IDLE", "10"))
    maxOpen, _ := strconv.Atoi(GetEnvWithDefault("DB_MAX_OPEN", "100"))
    
    return DBConfig{
        // ... other fields
        MaxIdle: maxIdle,
        MaxOpen: maxOpen,
    }
}
```

---

## Common Configuration Examples

### JWT Configuration

```go
type JWTConfig struct {
    Secret      string
    ExpireHours int
}

func loadJWTConfig() JWTConfig {
    expireHours, _ := strconv.Atoi(GetEnvWithDefault("JWT_EXPIRE_HOURS", "24"))
    return JWTConfig{
        Secret:      GetEnvMust("JWT_SECRET"),
        ExpireHours: expireHours,
    }
}
```

### AWS Configuration

```go
type AWSConfig struct {
    Region    string
    AccessKey string
    SecretKey string
    Bucket    string
}

func loadAWSConfig() AWSConfig {
    return AWSConfig{
        Region:    GetEnvMust("AWS_REGION"),
        AccessKey: GetEnvMust("AWS_ACCESS_KEY"),
        SecretKey: GetEnvMust("AWS_SECRET_KEY"),
        Bucket:    GetEnvMust("AWS_S3_BUCKET"),
    }
}
```

### SMTP Configuration

```go
type SMTPConfig struct {
    Host string
    Port string
    User string
    Pass string
    From string
}

func loadSMTPConfig() SMTPConfig {
    return SMTPConfig{
        Host: GetEnvMust("SMTP_HOST"),
        Port: GetEnvMust("SMTP_PORT"),
        User: GetEnvMust("SMTP_USER"),
        Pass: GetEnvMust("SMTP_PASS"),
        From: GetEnvMust("SMTP_FROM"),
    }
}
```

---

## Configuration Checklist

When adding new configuration:

- [ ] Add environment variables to `.env`
- [ ] Create config struct in `internal/config/config.go`
- [ ] Create loader function (e.g., `loadRedisConfig()`)
- [ ] Add to main `Config` struct
- [ ] Call loader in `LoadEnv()`
- [ ] Use `GetEnvMust()` for required values
- [ ] Use `GetEnv()` for optional values
- [ ] Create connection function if needed (e.g., `ConnectRedis()`)
- [ ] Update `main.go` to initialize connection
- [ ] Pass config values to services through container
- [ ] Add validation if needed
- [ ] Document in `.env.example` file

---

## Environment Files

### Development (.env)
```env
APP_PORT=8080
DB_HOST=localhost
DB_NAME=pm_db_dev
```

### Production (.env.production)
```env
APP_PORT=80
DB_HOST=prod-db.example.com
DB_NAME=pm_db_prod
```

### Example (.env.example)
```env
# Application
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=
DB_NAME=pm_db

# Add all required env vars here for documentation
```

**Note:** Never commit `.env` to git! Add to `.gitignore`:
```
.env
.env.local
.env.production
```
