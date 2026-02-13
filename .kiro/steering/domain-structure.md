---
inclusion: auto
description: Domain structure - models with BaseModel, DTOs for requests/responses, validation patterns
---

# Domain Structure Pattern

This project follows a clean domain-driven design with strict separation between models (database entities) and DTOs (request/response objects).

## Domain Architecture

```
internal/domain/
├── models/           # Database entities (GORM models)
│   ├── base_model.go # Common fields for all models
│   ├── client.go
│   ├── project.go
│   └── project_asset.go
├── dto/              # Data Transfer Objects (flat, one file per feature)
│   ├── auth/         # Auth DTOs (separate package)
│   │   └── auth.go
│   ├── client.go
│   ├── project.go
│   ├── project_asset.go
│   └── user.go
├── common/           # Shared types (pagination)
│   └── pagination.go
└── error/            # Domain errors
    ├── app_error.go
    └── error.go
```

## Key Principles

1. **Models = Database Entities**: All database tables are represented as models
2. **DTOs = API Contract**: All requests and responses use DTOs
3. **BaseModel = Common Fields**: Every model embeds BaseModel for consistency
4. **Type Safety**: Always use structs with validation tags
5. **Separation**: Never expose models directly in API responses

---

## Models (Database Entities)

### BaseModel

All models MUST embed `BaseModel` which provides common fields:

```go
// internal/domain/models/base_model.go
package models

import "time"

type BaseModel struct {
    ID uint64       // Primary key

    UUID string     // Public identifier

    Status string   // Record status (active, inactive, deleted, etc.)

    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt time.Time

    CreatedBy *uint64  // User ID who created this record
}
```

### Creating a New Model

**File:** `internal/domain/models/your_model.go`

```go
package models

type YourModel struct {
    BaseModel  // ALWAYS embed BaseModel first

    // Your specific fields
    Name        string
    Description string
    
    // Optional fields use pointers
    Email       *string
    Phone       *string
    
    // Foreign keys
    UserID      uint64
    CategoryID  uint64
}
```

### Model Examples

**Current User Model:**
```go
type User struct {
    BaseModel

    Name  string
    Email string

    Phone       *string
    CountryCode *int

    PasswordHash string

    EmployeeID *string

    GroupID       uint64
    DesignationID uint64
    DepartmentID  uint64
}
```

**Example Project Model:**
```go
type Project struct {
    BaseModel

    Name        string
    Description string
    
    StartDate   *time.Time
    EndDate     *time.Time
    
    OwnerID     uint64
    ClientID    uint64
}
```

### Model Rules

1. ✅ ALWAYS embed `BaseModel` first
2. ✅ Use pointers for optional fields (`*string`, `*int`, `*time.Time`)
3. ✅ Use `uint64` for IDs and foreign keys
4. ✅ Use descriptive field names (PascalCase)
5. ✅ Add GORM tags if needed: `gorm:"column:email;unique"`
6. ❌ NEVER expose models directly in API responses
7. ❌ NEVER add JSON tags to models (use DTOs instead)

---

## DTOs (Data Transfer Objects)

DTOs are used for ALL API requests and responses. They provide type safety, validation, and API contract definition.

### DTO Organization

DTOs are kept flat in `internal/domain/dto/` with one file per feature. Auth is the only exception with its own subfolder due to separate package needs.

```
internal/domain/dto/
├── auth/
│   └── auth.go           # LoginRequest, LoginResponse (separate package)
├── client.go             # CreateClientRequest, ClientResponse
├── project.go            # CreateProjectRequest, ProjectResponse
├── project_asset.go      # CreateProjectAssetRequest, ProjectAssetResponse
├── user.go               # CreateUserRequest, UserResponse
├── file.go               # CreateFileRequest, CreateFileResponse
├── list_with_pagination.go # ListWithMeta, PaginationMeta
└── utils.go              # Shared DTO types
```

### Request DTOs

Request DTOs define what data the API accepts with validation rules. All DTOs live in the `dto` package (not subpackages), except auth which has its own package.

**File:** `internal/domain/dto/client.go`

```go
package dto

type CreateClientRequest struct {
    Name    string  `json:"name" form:"name" binding:"required,min=2,max=150"`
    Email   *string `json:"email" form:"email" binding:"omitempty,email,max=255"`
    Phone   *string `json:"phone" form:"phone" binding:"omitempty,max=20"`
    Address *string `json:"address" form:"address" binding:"omitempty"`
}

type UpdateClientRequest struct {
    Name    *string `json:"name" form:"name" binding:"omitempty,min=2,max=150"`
    Email   *string `json:"email" form:"email" binding:"omitempty,email,max=255"`
    Phone   *string `json:"phone" form:"phone" binding:"omitempty,max=20"`
    Address *string `json:"address" form:"address" binding:"omitempty"`
}
```

**Validation Tags:**
- `binding:"required"` - Field is mandatory
- `binding:"email"` - Must be valid email
- `binding:"min=6"` - Minimum length
- `binding:"max=100"` - Maximum length
- `binding:"oneof=active inactive"` - Must be one of values
- `binding:"uuid"` - Must be valid UUID
- `binding:"omitempty"` - Skip validation if empty (used for optional and update fields)

### Response DTOs

Response DTOs define what data the API returns. They embed `models.BaseResponse` for timestamps and should NEVER include sensitive fields.

**File:** `internal/domain/dto/client.go`

```go
package dto

import "capecom-pm/internal/domain/models"

type ClientResponse struct {
    Id      string  `json:"id"`
    Name    string  `json:"name"`
    Email   *string `json:"email"`
    Phone   *string `json:"phone"`
    Address *string `json:"address"`

    models.BaseResponse
}
```

### DTO Rules

1. ✅ ALWAYS use structs for requests and responses
2. ✅ Add validation tags to request DTOs
3. ✅ Use JSON tags for all DTO fields
4. ✅ Use `omitempty` for optional response fields
5. ✅ Create separate DTOs for Create, Update, and Response
6. ✅ Keep DTOs flat in `internal/domain/dto/` with one file per feature (e.g., `dto/client.go`, `dto/project.go`)
7. ❌ NEVER expose sensitive fields (passwords, hashes, tokens)
8. ❌ NEVER use models directly as responses
9. ❌ NEVER add business logic to DTOs

---

## How to Add a New Feature

### Step 1: Create the Model

**File:** `internal/domain/models/project.go`

```go
package models

import "time"

type Project struct {
    BaseModel  // Always embed BaseModel

    Name        string
    Description string
    
    StartDate   *time.Time
    EndDate     *time.Time
    
    OwnerID     uint64
    ClientID    uint64
}
```

### Step 2: Create DTOs

**File:** `internal/domain/dto/project.go`

```go
package dto

import (
    "capecom-pm/internal/domain/models"
    "time"
)

type CreateProjectRequest struct {
    ProjectName string  `json:"project_name" form:"project_name" binding:"required,min=2,max=120"`
    ProjectCode string  `json:"project_code" form:"project_code" binding:"required,min=2,max=120"`
    ClientUUID  *string `json:"client_id" form:"client_id" binding:"omitempty,uuid"`
}

type ProjectResponse struct {
    Id          string     `json:"id"`
    ProjectName string     `json:"project_name"`
    ProjectCode string     `json:"project_code"`
    ClientID    *string    `json:"client_id"`
    ClientName  *string    `json:"client_name"`
    Status      string     `json:"status"`

    models.BaseResponse
}
```

### Step 3: Use in Handler

```go
package handlers

import (
    "capecom-pm/internal/domain/dto"
    "capecom-pm/internal/utils"
    "capecom-pm/internal/utils/bind"
    "capecom-pm/internal/utils/response"
    "net/http"
    "github.com/gin-gonic/gin"
)

func (h *ProjectHandler) CreateProject(c *gin.Context) {
    var req dto.CreateProjectRequest
    if !bind.AndValidate(c, &req, "create_project") {
        return
    }

    userID := utils.GetUserID(c)

    project, err := h.ProjectService.Create(req, userID)
    if err != nil {
        response.FromError(c, err)
        return
    }

    response.JSON(c, http.StatusCreated, response.APIResponse{
        Success: true,
        Data:    project,
    })
}
```

---

## Common Patterns

### Pagination Request/Response

```go
type PaginationRequest struct {
    Page    int `form:"page" binding:"omitempty,min=1"`
    PerPage int `form:"per_page" binding:"omitempty,min=1,max=100"`
}

type PaginatedResponse struct {
    Data    interface{} `json:"data"`
    Total   int         `json:"total"`
    Page    int         `json:"page"`
    PerPage int         `json:"per_page"`
}
```

### Filter/Search Request

```go
type UserFilterRequest struct {
    Search      *string `form:"search"`
    Status      *string `form:"status" binding:"omitempty,oneof=active inactive"`
    DepartmentID *uint64 `form:"department_id"`
    Page        int     `form:"page" binding:"omitempty,min=1"`
    PerPage     int     `form:"per_page" binding:"omitempty,min=1,max=100"`
}
```

### Nested Response

```go
type ProjectDetailResponse struct {
    ID          uint64        `json:"id"`
    Name        string        `json:"name"`
    Owner       UserResponse  `json:"owner"`      // Nested user
    Tasks       []TaskResponse `json:"tasks"`     // Nested tasks
    CreatedAt   time.Time     `json:"created_at"`
}
```

---

## Validation Reference

Common Gin validation tags:

- `required` - Field must be present
- `omitempty` - Skip validation if empty
- `email` - Valid email format
- `min=X` - Minimum value/length
- `max=X` - Maximum value/length
- `len=X` - Exact length
- `oneof=a b c` - Must be one of values
- `uuid` - Valid UUID format
- `e164` - Valid phone number (E.164 format)
- `url` - Valid URL
- `datetime=2006-01-02` - Valid date format

---

## Summary Checklist

When adding a new feature:

- [ ] Create model in `internal/domain/models/`
- [ ] Embed `BaseModel` in the model
- [ ] Create DTO file in `internal/domain/dto/feature.go`
- [ ] Create request DTOs with validation tags
- [ ] Create response DTOs with JSON tags
- [ ] Never expose models directly in API
- [ ] Never include sensitive fields in responses
- [ ] Use pointers for optional fields
- [ ] Add proper validation rules
