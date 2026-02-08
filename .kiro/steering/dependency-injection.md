---
inclusion: auto
---

# Dependency Injection Pattern

This project uses a clean DI pattern with a container-based approach. Follow this guide when adding new features.

## Architecture Layers

```
Database (gorm.DB)
    ↓
Repository Layer (data access)
    ↓
Service Layer (business logic)
    ↓
Handler Layer (HTTP handlers)
```

## How to Add a New Feature

### 1. Create the Repository

**File:** `internal/repositories/your_feature.go`

```go
package repositories

import "gorm.io/gorm"

type YourFeatureRepo struct {
    db *gorm.DB
}

func NewYourFeatureRepo(db *gorm.DB) *YourFeatureRepo {
    return &YourFeatureRepo{db: db}
}

// Add your methods here
func (r *YourFeatureRepo) GetSomething(id uint) (*Model, error) {
    // Implementation
}
```

### 2. Register Repository in Container

**File:** `internal/container/repository.go`

```go
type Repository struct {
    AuthRepo        *repositories.AuthRepo
    YourFeatureRepo *repositories.YourFeatureRepo  // Add this
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{
        AuthRepo:        repositories.NewAuthRepo(db),
        YourFeatureRepo: repositories.NewYourFeatureRepo(db),  // Add this
    }
}
```

### 3. Create the Service

**File:** `internal/services/your_feature.go`

```go
package services

import "capecom-pm/internal/repositories"

type YourFeatureService struct {
    repo *repositories.YourFeatureRepo
}

func NewYourFeatureService(repo *repositories.YourFeatureRepo) *YourFeatureService {
    return &YourFeatureService{repo: repo}
}

// Add your business logic methods here
func (s *YourFeatureService) DoSomething(id uint) error {
    // Business logic
    return s.repo.GetSomething(id)
}
```

### 4. Register Service in Container

**File:** `internal/container/service.go`

```go
type Service struct {
    AuthService        *services.AuthService
    YourFeatureService *services.YourFeatureService  // Add this
}

func NewService(db *gorm.DB, repository *Repository) *Service {
    return &Service{
        AuthService:        services.NewAuthService(repository.AuthRepo),
        YourFeatureService: services.NewYourFeatureService(repository.YourFeatureRepo),  // Add this
    }
}
```

### 5. Create the Handler

**File:** `internal/handlers/your_feature_handler.go`

```go
package handlers

import (
    "capecom-pm/internal/services"
    "github.com/gin-gonic/gin"
)

type YourFeatureHandler struct {
    service *services.YourFeatureService
}

func NewYourFeatureHandler(service *services.YourFeatureService) *YourFeatureHandler {
    return &YourFeatureHandler{service: service}
}

// Add your HTTP handlers here
func (h *YourFeatureHandler) HandleSomething(c *gin.Context) {
    // Handler logic
    h.service.DoSomething(id)
}
```

### 6. Register Handler in Container

**File:** `internal/container/handler.go`

```go
type Handler struct {
    AuthHandler        *handlers.AuthHandler
    YourFeatureHandler *handlers.YourFeatureHandler  // Add this
}

func NewHandler(service *Service) *Handler {
    return &Handler{
        AuthHandler:        handlers.NewAuthHandler(service.AuthService),
        YourFeatureHandler: handlers.NewYourFeatureHandler(service.YourFeatureService),  // Add this
    }
}
```

## Key Principles

1. **Single Responsibility**: Each layer has one job
   - Repository: Database operations
   - Service: Business logic
   - Handler: HTTP request/response

2. **Dependency Flow**: Always flows downward
   - Handlers depend on Services
   - Services depend on Repositories
   - Repositories depend on DB

3. **Constructor Injection**: Use `New*()` functions to inject dependencies

4. **No Global State**: Everything is passed through constructors

5. **Testability**: Easy to mock dependencies for testing

## Example: Current Auth Implementation

See these files for reference:
- `internal/repositories/auth.go`
- `internal/services/auth.go`
- `internal/handlers/auth_handler.go`

The container wires them all together in `internal/container/container.go`
