---
inclusion: auto
---

# Domain Structure Pattern

This project follows a clean domain-driven design with strict separation between models (database entities) and DTOs (request/response objects).

## Domain Architecture

```
internal/domain/
├── models/           # Database entities (GORM models)
│   ├── base_model.go # Common fields for all models
│   └── user.go       # User entity
└── dto/              # Data Transfer Objects (request/response)
    └── auth/         # Auth-related DTOs
        └── login.go  # Login request/response
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

```
internal/domain/dto/
├── auth/
│   └── login.go          # LoginRequest, LoginResponse
├── user/
│   └── user.go           # CreateUserRequest, UpdateUserRequest, UserResponse
└── project/
    └── project.go        # CreateProjectRequest, ProjectResponse
```

### Request DTOs

Request DTOs define what data the API accepts with validation rules.

**File:** `internal/domain/dto/auth/login.go`

```go
package authdto

type LoginRequest struct {
    Email    string `json:"email" form:"email" binding:"required,email"`
    Password string `json:"password" form:"password" binding:"required,min=6"`
}
```

**Validation Tags:**
- `binding:"required"` - Field is mandatory
- `binding:"email"` - Must be valid email
- `binding:"min=6"` - Minimum length
- `binding:"max=100"` - Maximum length
- `binding:"oneof=active inactive"` - Must be one of values
- `binding:"uuid"` - Must be valid UUID

**Example User Request DTOs:**
```go
package userdto

type CreateUserRequest struct {
    Name        string  `json:"name" binding:"required,min=2,max=100"`
    Email       string  `json:"email" binding:"required,email"`
    Password    string  `json:"password" binding:"required,min=8"`
    Phone       *string `json:"phone" binding:"omitempty,e164"`
    CountryCode *int    `json:"country_code" binding:"omitempty"`
}

type UpdateUserRequest struct {
    Name        *string `json:"name" binding:"omitempty,min=2,max=100"`
    Email       *string `json:"email" binding:"omitempty,email"`
    Phone       *string `json:"phone" binding:"omitempty,e164"`
    CountryCode *int    `json:"country_code" binding:"omitempty"`
}
```

### Response DTOs

Response DTOs define what data the API returns. They should NEVER include sensitive fields.

**File:** `internal/domain/dto/auth/login.go`

```go
package authdto

type LoginResponse struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int    `json:"expires_in"`
}
```

**Example User Response DTOs:**
```go
package userdto

import "time"

type UserResponse struct {
    ID          uint64     `json:"id"`
    UUID        string     `json:"uuid"`
    Name        string     `json:"name"`
    Email       string     `json:"email"`
    Phone       *string    `json:"phone,omitempty"`
    CountryCode *int       `json:"country_code,omitempty"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type UserListResponse struct {
    Users      []UserResponse `json:"users"`
    Total      int            `json:"total"`
    Page       int            `json:"page"`
    PerPage    int            `json:"per_page"`
}
```

### DTO Rules

1. ✅ ALWAYS use structs for requests and responses
2. ✅ Add validation tags to request DTOs
3. ✅ Use JSON tags for all DTO fields
4. ✅ Use `omitempty` for optional response fields
5. ✅ Create separate DTOs for Create, Update, and Response
6. ✅ Group related DTOs in the same package (e.g., `userdto`)
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

**File:** `internal/domain/dto/project/project.go`

```go
package projectdto

import "time"

// Request DTOs
type CreateProjectRequest struct {
    Name        string     `json:"name" binding:"required,min=2,max=200"`
    Description string     `json:"description" binding:"required"`
    StartDate   *time.Time `json:"start_date" binding:"omitempty"`
    EndDate     *time.Time `json:"end_date" binding:"omitempty"`
    ClientID    uint64     `json:"client_id" binding:"required"`
}

type UpdateProjectRequest struct {
    Name        *string    `json:"name" binding:"omitempty,min=2,max=200"`
    Description *string    `json:"description" binding:"omitempty"`
    StartDate   *time.Time `json:"start_date" binding:"omitempty"`
    EndDate     *time.Time `json:"end_date" binding:"omitempty"`
    Status      *string    `json:"status" binding:"omitempty,oneof=active inactive completed"`
}

// Response DTOs
type ProjectResponse struct {
    ID          uint64     `json:"id"`
    UUID        string     `json:"uuid"`
    Name        string     `json:"name"`
    Description string     `json:"description"`
    StartDate   *time.Time `json:"start_date,omitempty"`
    EndDate     *time.Time `json:"end_date,omitempty"`
    Status      string     `json:"status"`
    OwnerID     uint64     `json:"owner_id"`
    ClientID    uint64     `json:"client_id"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type ProjectListResponse struct {
    Projects []ProjectResponse `json:"projects"`
    Total    int               `json:"total"`
    Page     int               `json:"page"`
    PerPage  int               `json:"per_page"`
}
```

### Step 3: Use in Handler

```go
package handlers

import (
    "capecom-pm/internal/domain/dto/project"
    "github.com/gin-gonic/gin"
)

func (h *ProjectHandler) Create(c *gin.Context) {
    var req projectdto.CreateProjectRequest
    
    // Bind and validate request
    if err := c.ShouldBindJSON(&req); err != nil {
        // Handle validation error
        return
    }
    
    // Call service
    project, err := h.service.Create(&req)
    if err != nil {
        // Handle error
        return
    }
    
    // Return response DTO
    response := projectdto.ProjectResponse{
        ID:          project.ID,
        UUID:        project.UUID,
        Name:        project.Name,
        Description: project.Description,
        // ... map other fields
    }
    
    c.JSON(201, response)
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
- [ ] Create DTO package in `internal/domain/dto/feature/`
- [ ] Create request DTOs with validation tags
- [ ] Create response DTOs with JSON tags
- [ ] Never expose models directly in API
- [ ] Never include sensitive fields in responses
- [ ] Use pointers for optional fields
- [ ] Add proper validation rules
