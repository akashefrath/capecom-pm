# FEATURE CREATION REFERENCE

## QUICK REFERENCE FOR NEW FEATURES

This document provides a quick reference for creating new features following the project architecture.

## EXAMPLE: Role Feature (Get All Active)

### 1. Repository (internal/src/repository/role.go)
```go
package repository

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	models "github.com/akashefrath/capecom-pm/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

type Role struct {
	DB *sqlx.DB
}

func NewRole(db *sqlx.DB) *Role {
	return &Role{DB: db}
}

func (r *Role) GetAllActive() ([]dto.RoleResponse, error) {
	var roles []dto.RoleResponse
	q := `SELECT uuid, name, status FROM roles WHERE deleted_at IS NULL AND status = ?`
	err := r.DB.Select(&roles, q, models.StatusActive)
	return roles, err
}
```

### 2. Service (internal/src/service/role.go)
```go
package service

import (
	"github.com/akashefrath/capecom-pm/internal/domain/dto"
	"github.com/akashefrath/capecom-pm/internal/src/repository"
)

type Role struct {
	RoleRepo *repository.Role
}

func NewRole(roleRepo *repository.Role) *Role {
	return &Role{RoleRepo: roleRepo}
}

func (s *Role) GetAllActive() ([]dto.RoleResponse, error) {
	return s.RoleRepo.GetAllActive()
}
```

### 3. Handler (internal/src/handler/role.go)
```go
package handler

import (
	"github.com/akashefrath/capecom-pm/internal/src/service"
	"github.com/akashefrath/capecom-pm/internal/utils/bind"
	"github.com/akashefrath/capecom-pm/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	Role *service.Role
}

func NewRole(role *service.Role) RoleHandler {
	return RoleHandler{Role: role}
}

// GET endpoint - no request body
func (h *RoleHandler) GetAllActive(c *gin.Context) {
	data, err := h.Role.GetAllActive()
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONOk(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}

// POST endpoint - with request body validation
func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	isValid := bind.AndValidate(c, &req, "role")
	if !isValid {
		return
	}
	
	data, err := h.Role.Create(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
```

### 4. DTO (internal/domain/dto/role.go)
```go
package dto

type RoleResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
```

### 5. Model (internal/domain/model/role.go)
```go
package models

type Role struct {
	BaseModel
	Name string `db:"name" json:"name"`
}
```

### 6. Container - Repositories (internal/container/repositories.go)
```go
type Repositories struct {
	// ... existing
	Role    *repository.Role
}

func NewRepository(db *sqlx.DB, config *config.Config, redis *redis.Client, dbTX *database.Database) *Repositories {
	return &Repositories{
		// ... existing
		Role:    repository.NewRole(db),
	}
}
```

### 7. Container - Service (internal/container/service.go)
```go
type Service struct {
	// ... existing
	Role  *service.Role
}

func NewService(repo *Repositories, jwtManager *jwtutil.Manager) *Service {
	return &Service{
		// ... existing
		Role:  service.NewRole(repo.Role),
	}
}
```

### 8. Container - Handler (internal/container/handler.go)
```go
type Handler struct {
	// ... existing
	RoleHandler  handler.RoleHandler
}

func SetupHandler(service *Service) *Handler {
	return &Handler{
		// ... existing
		RoleHandler:  handler.NewRole(service.Role),
	}
}
```

### 9. Routes (internal/routes/v1/role.go)
```go
package v1

import (
	"github.com/akashefrath/capecom-pm/internal/container"
	"github.com/gin-gonic/gin"
)

func Role(r *gin.RouterGroup, container container.Container) {
	role := r.Group("role")
	handler := container.Handler.RoleHandler
	role.GET("active", handler.GetAllActive)
}
```

### 10. Register Route (internal/routes/v1/v1.go)
```go
func Setup(r *gin.Engine, container container.Container) {
	v1API := r.Group("api/v1")
	// ... existing
	Role(v1API, container)
}
```

## KEY PATTERNS

### Import Path
Always use: `github.com/akashefrath/capecom-pm/internal/...`

### Request Binding & Validation
Use `bind.AndValidate()` for automatic validation with error handling:
```go
var req dto.CreateRequest
isValid := bind.AndValidate(c, &req, "entity_name")
if !isValid {
    return  // Error response sent automatically
}
```

Features:
- Automatically binds JSON/form data to struct
- Validates using struct tags (e.g., `binding:"required"`)
- Returns translated validation errors
- Sends 400 Bad Request with field-specific errors
- No need for manual error response

### SQL Queries
- Use raw SQL with parameterized queries
- Select specific columns only
- Use `r.DB.Select()` for multiple rows
- Use `r.DB.Get()` for single row
- Use `r.DB.Exec()` for insert/update/delete

### Error Handling
- Repository returns error
- Service passes error through
- Handler uses `response.FromError(c, err)`

### Response Pattern
```go
// 200 OK
response.JSONOk(c, response.APIResponse{
    Success: true,
    Data:    data,
})

// 201 Created
response.JSONCreated(c, response.APIResponse{
    Success: true,
    Data:    data,
})
```

## COMMON QUERIES

### Get All Active
```go
q := `SELECT uuid, name FROM table WHERE deleted_at IS NULL AND status = ?`
err := r.DB.Select(&results, q, models.StatusActive)
```

### Get By ID
```go
q := `SELECT uuid, name FROM table WHERE uuid = ? AND deleted_at IS NULL`
err := r.DB.Get(&result, q, id)
```

### Create
```go
q := `INSERT INTO table (uuid, name, status) VALUES (?, ?, ?)`
_, err := r.DB.Exec(q, uuid, name, status)
```

### Update
```go
q := `UPDATE table SET name = ? WHERE uuid = ? AND deleted_at IS NULL`
_, err := r.DB.Exec(q, name, uuid)
```

### Soft Delete
```go
q := `UPDATE table SET deleted_at = NOW() WHERE uuid = ?`
_, err := r.DB.Exec(q, uuid)
```


## REQUEST VALIDATION

### DTO with Validation Tags
```go
package dto

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
	Code string `json:"code" binding:"required,alphanum"`
}

type UpdateRoleRequest struct {
	Name string `json:"name" binding:"omitempty,min=3,max=50"`
}
```

### Common Validation Tags
- `required` - Field must be present
- `omitempty` - Skip validation if empty
- `min=3` - Minimum length/value
- `max=50` - Maximum length/value
- `email` - Valid email format
- `alphanum` - Alphanumeric only
- `uuid` - Valid UUID format
- `oneof=active inactive` - Must be one of values

### Handler with Validation
```go
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateRequest
	isValid := bind.AndValidate(c, &req, "entity")
	if !isValid {
		return  // Validation errors sent automatically
	}
	
	// Continue with validated request
	data, err := h.Service.Create(req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.JSONCreated(c, response.APIResponse{
		Success: true,
		Data:    data,
	})
}
```

### Query Parameters
```go
func (h *Handler) List(c *gin.Context) {
	var query dto.ListQuery
	isValid := bind.AndValidate(c, &query, "entity")
	if !isValid {
		return
	}
	
	data, err := h.Service.List(query)
	// ... rest of handler
}
```

### Pagination
```go
func (h *Handler) List(c *gin.Context) {
	pagination, err := bind.PaginationBinder(c, "entity")
	if err != nil {
		response.FromError(c, err)
		return
	}
	
	data, err := h.Service.List(pagination)
	// ... rest of handler
}
```
