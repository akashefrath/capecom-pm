---
inclusion: auto
---

# Error Handling & Flow Pattern

This project uses a structured error handling system with i18n support, domain errors, and consistent API responses.

## Error Flow Architecture

```
Handler (validates request)
    ↓
Service (business logic, returns domain errors)
    ↓
Repository (data access, handles DB errors)
    ↓
Error Mapper (maps to HTTP codes + translates)
    ↓
API Response (consistent JSON format)
```

---

## Complete Request Flow

### 1. Handler Layer (Entry Point)

**Responsibilities:**
- Bind and validate request DTO
- Call service layer
- Handle errors using `response.FromError()`
- Return success response

**Pattern:**
```go
package handlers

import (
    authdto "capecom-pm/internal/domain/dto/auth"
    "capecom-pm/internal/services"
    "capecom-pm/internal/utils/bind"
    "capecom-pm/internal/utils/response"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    AuthService *services.AuthService
}

func NewAuthHandler(AuthService *services.AuthService) *AuthHandler {
    return &AuthHandler{
        AuthService: AuthService,
    }
}

func (h *AuthHandler) Login(c *gin.Context) {
    // Step 1: Bind and validate request
    var req authdto.LoginRequest
    if !bind.AndValidate(c, &req, "auth") {
        return  // Validation errors already sent
    }

    // Step 2: Call service
    result, err := h.AuthService.Login(c, req)
    if err != nil {
        // Step 3: Handle error (auto-translates and maps to HTTP code)
        response.FromError(c, err)
        return
    }

    // Step 4: Return success response with translation
    lang := c.GetHeader("Accept-Language")
    msgs := i18n.GetMessages(lang)
    
    response.JSON(c, 200, response.APIResponse{
        Success: true,
        Message: msgs["login_success"],
        Data:    result,
    })
}
```

### 2. Service Layer (Business Logic)

**Responsibilities:**
- Implement business logic
- Call repository methods
- Return domain errors (from `internal/domain/error`)
- NEVER return raw database errors

**Pattern:**
```go
package services

import (
    authdto "capecom-pm/internal/domain/dto/auth"
    domainerrors "capecom-pm/internal/domain/error"
    "capecom-pm/internal/repositories"
    "github.com/gin-gonic/gin"
)

type AuthService struct {
    AuthRepo *repositories.AuthRepo
}

func NewAuthService(AuthRepo *repositories.AuthRepo) *AuthService {
    return &AuthService{
        AuthRepo: AuthRepo,
    }
}

func (s *AuthService) Login(c *gin.Context, req authdto.LoginRequest) (*authdto.LoginResponse, error) {
    // Step 1: Get user from repository
    user, err := s.AuthRepo.FindUserByEmailAndReturnPassword(req.Email)
    if err != nil {
        // Database error - return internal error
        return nil, domainerrors.ErrInternal
    }
    
    // Step 2: Check if user exists
    if user == nil {
        return nil, domainerrors.ErrInvalidCredentials
    }

    // Step 3: Verify password
    if !utils.VerifyPassword(user.PasswordHash, req.Password) {
        return nil, domainerrors.ErrInvalidCredentials
    }

    // Step 4: Generate token
    token, err := generateToken(user)
    if err != nil {
        return nil, domainerrors.ErrInternal
    }

    // Step 5: Return success response
    return &authdto.LoginResponse{
        AccessToken: token,
        TokenType:   "Bearer",
        ExpiresIn:   3600,
    }, nil
}
```

**Service Rules:**
1. ✅ ALWAYS return domain errors (from `internal/domain/error`)
2. ✅ Handle repository errors and convert to domain errors
3. ✅ Implement business logic validation
4. ✅ Return DTOs, never models
5. ❌ NEVER return raw database errors
6. ❌ NEVER return strings as errors
7. ❌ NEVER handle HTTP responses

### 3. Repository Layer (Data Access)

**Responsibilities:**
- Execute database queries
- Handle GORM errors
- Return nil for "not found" (not an error)
- Return actual errors for DB failures

**Pattern:**
```go
package repositories

import (
    "capecom-pm/internal/domain/models"
    "errors"
    "gorm.io/gorm"
)

type AuthRepo struct {
    DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
    return &AuthRepo{
        DB: db,
    }
}

func (r *AuthRepo) FindUserByEmailAndReturnPassword(email string) (*models.User, error) {
    var user models.User
    err := r.DB.Where("email = ?", email).First(&user).Error

    // Not found is NOT an error - return nil, nil
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }

    // Return user and any actual error
    return &user, err
}

func (r *AuthRepo) CreateUser(user *models.User) error {
    return r.DB.Create(user).Error
}

func (r *AuthRepo) UpdateUser(user *models.User) error {
    return r.DB.Save(user).Error
}

func (r *AuthRepo) DeleteUser(id uint64) error {
    return r.DB.Delete(&models.User{}, id).Error
}

func (r *AuthRepo) FindByID(id uint64) (*models.User, error) {
    var user models.User
    err := r.DB.First(&user, id).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    
    return &user, err
}

func (r *AuthRepo) FindAll(limit, offset int) ([]models.User, error) {
    var users []models.User
    err := r.DB.Limit(limit).Offset(offset).Find(&users).Error
    return users, err
}
```

**Repository Rules:**
1. ✅ Return `nil, nil` for "not found" (let service decide if it's an error)
2. ✅ Return actual database errors
3. ✅ Use GORM methods correctly
4. ✅ Return models, not DTOs
5. ❌ NEVER return domain errors
6. ❌ NEVER implement business logic
7. ❌ NEVER handle HTTP responses

---

## Domain Errors

All errors are defined in `internal/domain/error/error.go`:

```go
package domainerrors

import "errors"

var (
    // Auth errors
    ErrInvalidCredentials = errors.New("invalid_login_credentials")
    ErrUnauthorized       = errors.New("unauthorized")

    // User errors
    ErrUserNotFound   = errors.New("user_not_found")
    ErrDuplicateEmail = errors.New("duplicate_email")

    // Common errors
    ErrBadRequest = errors.New("bad_request")
    ErrInternal   = errors.New("internal_error")
    ErrNotFound   = errors.New("not_found")
    ErrForbidden  = errors.New("forbidden")
)
```

**How to Add New Domain Error:**

1. Add to `internal/domain/error/error.go`:
```go
var (
    ErrProjectNotFound = errors.New("project_not_found")
)
```

2. Add to error mapper `internal/utils/response/error_mapper.go`:
```go
case errors.Is(err, domainerrors.ErrProjectNotFound):
    code = http.StatusNotFound
```

3. Add translations to `internal/utils/i18n/en.go`:
```go
var EN = map[string]string{
    "project_not_found": "Project not found",
}
```

4. Add translations to `internal/utils/i18n/ta.go`:
```go
var TA = map[string]string{
    "project_not_found": "திட்டம் கிடைக்கவில்லை",
}
```

---

## Error Mapper

The error mapper (`internal/utils/response/error_mapper.go`) converts domain errors to HTTP status codes and translates messages:

```go
package response

import (
    domainerrors "capecom-pm/internal/domain/error"
    "capecom-pm/internal/utils/i18n"
    "errors"
    "net/http"
    "github.com/gin-gonic/gin"
)

func FromError(c *gin.Context, err error) {
    // Get language from header
    lang := c.GetHeader("Accept-Language")
    
    // Default to internal server error
    code := http.StatusInternalServerError
    
    // Map domain errors to HTTP codes
    switch {
    case errors.Is(err, domainerrors.ErrInvalidCredentials):
        code = http.StatusUnauthorized
    case errors.Is(err, domainerrors.ErrDuplicateEmail):
        code = http.StatusConflict
    case errors.Is(err, domainerrors.ErrUserNotFound):
        code = http.StatusNotFound
    case errors.Is(err, domainerrors.ErrBadRequest):
        code = http.StatusBadRequest
    case errors.Is(err, domainerrors.ErrUnauthorized):
        code = http.StatusUnauthorized
    case errors.Is(err, domainerrors.ErrForbidden):
        code = http.StatusForbidden
    }

    // Get translated message
    message := i18n.GetMessages(lang)
    
    // Return error response
    JSON(c, code, APIResponse{
        Success: false,
        Message: message[err.Error()],
    })
}
```

**HTTP Status Code Mapping:**
- `400 Bad Request` - Invalid input, validation errors
- `401 Unauthorized` - Authentication required or failed
- `403 Forbidden` - Authenticated but not authorized
- `404 Not Found` - Resource doesn't exist
- `409 Conflict` - Duplicate resource (e.g., email exists)
- `500 Internal Server Error` - Unexpected errors

---

## Validation Errors

Validation is handled by `internal/utils/bind/bind.go`:

```go
package bind

import (
    "capecom-pm/internal/utils/i18n"
    "capecom-pm/internal/utils/response"
    "errors"
    "fmt"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type FieldErrors map[string][]string

func AndValidate(c *gin.Context, req any, entity string) bool {
    if err := c.ShouldBind(req); err != nil {
        var ve validator.ValidationErrors

        lang := c.GetHeader("Accept-Language")
        msgs := i18n.GetMessages(lang)

        errs := FieldErrors{}

        if errors.As(err, &ve) {
            // Field validation errors
            for _, fe := range ve {
                field := strings.ToLower(fe.Field())
                tag := fe.Tag()
                param := fe.Param()

                msg := msgs[tag]
                if param != "" {
                    msg = fmt.Sprintf(msg, param)
                }

                errs[field] = append(errs[field], msg)
            }
        } else {
            // Body parsing error
            errs["body"] = []string{msgs["invalid_body"]}
        }

        response.JSON(c, http.StatusBadRequest, response.APIResponse{
            Success: false,
            Errors:  errs,
            Func:    "validate",
            Entity:  entity,
        })

        return false
    }

    return true
}
```

**Validation Error Response:**
```json
{
    "success": false,
    "func": "validate",
    "entity": "auth",
    "errors": {
        "email": ["Invalid email format"],
        "password": ["Minimum 6 characters required"]
    }
}
```

---

## i18n (Internationalization)

Translation files are in `internal/utils/i18n/`:

**English (`en.go`):**
```go
var EN = map[string]string{
    // Validation
    "required":     "This field is required",
    "email":        "Invalid email format",
    "min":          "Minimum %s characters required",
    "max":          "Maximum %s characters allowed",
    "invalid_body": "Invalid request body",

    // Common
    "bad_request":    "Bad request",
    "internal_error": "Internal server error",
    "unauthorized":   "Unauthorized access",

    // User
    "user_not_found":  "User not found",
    "duplicate_email": "Email already exists",

    // Auth
    "invalid_login_credentials": "Invalid login credentials",
    
    // Success messages
    "login_success":        "Login successful",
    "user_created_success": "User created successfully",
    "user_updated_success": "User updated successfully",
    "user_deleted_success": "User deleted successfully",
}
```

**Tamil (`ta.go`):**
```go
var TA = map[string]string{
    "required":     "இந்த புலம் அவசியம்",
    "email":        "தவறான மின்னஞ்சல் வடிவம்",
    // ... other translations
}
```

**How to Add Translation:**
1. Add key-value to `en.go` (both errors and success messages)
2. Add same key with translated value to `ta.go`
3. Use the key in domain errors or in handler responses

**Using Translations in Handlers:**
```go
// Get language and messages
lang := c.GetHeader("Accept-Language")
msgs := i18n.GetMessages(lang)

// Use in response
response.JSON(c, 200, response.APIResponse{
    Success: true,
    Message: msgs["operation_success"],
    Data:    result,
})
```

**Language Detection:**
- Client sends `Accept-Language: ta` header
- System defaults to English if not specified
- Add more languages by creating new files (e.g., `fr.go`, `es.go`)

---

## API Response Format

All responses use `internal/utils/response/api_response.go`:

```go
type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Func    string      `json:"func,omitempty"`
    Entity  string      `json:"entity,omitempty"`
    Errors  interface{} `json:"errors,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Meta    interface{} `json:"meta,omitempty"`
}

type PageMeta struct {
    Page    int     `json:"page,omitempty"`
    Limit   int     `json:"limit,omitempty"`
    Total   int64   `json:"total,omitempty"`
    Pages   int     `json:"pages,omitempty"`
    NextID  *string `json:"next_id,omitempty"`
    HasNext bool    `json:"has_next,omitempty"`
}
```

**Success Response:**
```json
{
    "success": true,
    "message": "Login successful",  // Translated based on Accept-Language header
    "data": {
        "access_token": "eyJhbGc...",
        "token_type": "Bearer",
        "expires_in": 3600
    }
}
```

**Error Response:**
```json
{
    "success": false,
    "message": "Invalid login credentials"
}
```

**Validation Error Response:**
```json
{
    "success": false,
    "func": "validate",
    "entity": "user",
    "errors": {
        "email": ["Invalid email format"],
        "password": ["Minimum 8 characters required"]
    }
}
```

**Paginated Response:**
```json
{
    "success": true,
    "data": [...],
    "meta": {
        "page": 1,
        "limit": 10,
        "total": 100,
        "pages": 10,
        "has_next": true
    }
}
```

---

## Complete Example: User CRUD

### Handler
```go
func (h *UserHandler) Create(c *gin.Context) {
    var req userdto.CreateUserRequest
    if !bind.AndValidate(c, &req, "user") {
        return
    }

    user, err := h.service.Create(c, req)
    if err != nil {
        response.FromError(c, err)
        return
    }

    // Get translated success message
    lang := c.GetHeader("Accept-Language")
    msgs := i18n.GetMessages(lang)

    response.JSON(c, 201, response.APIResponse{
        Success: true,
        Message: msgs["user_created_success"],
        Data:    user,
    })
}
```

### Service
```go
func (s *UserService) Create(c *gin.Context, req userdto.CreateUserRequest) (*userdto.UserResponse, error) {
    // Check if email exists
    existing, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        return nil, domainerrors.ErrInternal
    }
    if existing != nil {
        return nil, domainerrors.ErrDuplicateEmail
    }

    // Hash password
    hash, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, domainerrors.ErrInternal
    }

    // Create user model
    user := &models.User{
        Name:         req.Name,
        Email:        req.Email,
        PasswordHash: hash,
    }

    // Save to database
    if err := s.repo.Create(user); err != nil {
        return nil, domainerrors.ErrInternal
    }

    // Return response DTO
    return &userdto.UserResponse{
        ID:        user.ID,
        UUID:      user.UUID,
        Name:      user.Name,
        Email:     user.Email,
        Status:    user.Status,
        CreatedAt: user.CreatedAt,
    }, nil
}
```

### Repository
```go
func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.DB.Where("email = ?", email).First(&user).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    
    return &user, err
}

func (r *UserRepo) Create(user *models.User) error {
    return r.DB.Create(user).Error
}
```

---

## Error Handling Checklist

When implementing a new feature:

- [ ] Define domain errors in `internal/domain/error/error.go`
- [ ] Add error mapping in `internal/utils/response/error_mapper.go`
- [ ] Add error translations in `internal/utils/i18n/en.go`
- [ ] Add error translations in `internal/utils/i18n/ta.go`
- [ ] Add success message translations in `internal/utils/i18n/en.go`
- [ ] Add success message translations in `internal/utils/i18n/ta.go`
- [ ] Handler: Use `bind.AndValidate()` for validation
- [ ] Handler: Use `response.FromError()` for errors
- [ ] Handler: Get `i18n.GetMessages(lang)` for success messages
- [ ] Handler: Use `response.JSON()` with translated messages
- [ ] Service: Return domain errors only
- [ ] Service: Convert repo errors to domain errors
- [ ] Repository: Return `nil, nil` for not found
- [ ] Repository: Return actual DB errors
