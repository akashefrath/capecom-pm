---
inclusion: auto
description: Common patterns - database patterns (BaseModel, soft delete, UUID), utilities, CRUD operations, naming conventions
---

# Common Patterns & Utilities

This document covers common patterns, utilities, and conventions used throughout the project.

---

## Database Patterns

### BaseModel Structure

Every database table MUST follow the BaseModel pattern:

```sql
CREATE TABLE `table_name` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  
  -- Your specific fields here
  `name` varchar(120) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  
  -- Standard fields (always at the end)
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp
);
```

**Corresponding GORM Model:**
```go
type BaseModel struct {
    ID uint64

    UUID string

    Status string

    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt time.Time

    CreatedBy *uint64
}
```

### Status Enum

All tables use the same status enum:
- `active` - Default, record is active
- `inactive` - Temporarily disabled
- `blocked` - Blocked by admin/system
- `archived` - Soft archived, not deleted

**Usage in queries:**
```go
// Get only active records
db.Where("status = ?", "active").Find(&users)

// Exclude archived
db.Where("status != ?", "archived").Find(&users)

// Get all except deleted
db.Where("deleted_at IS NULL").Find(&users)
```

### Soft Delete Pattern

Use `deleted_at` timestamp for soft deletes:

```go
// Soft delete (set deleted_at)
db.Model(&user).Update("deleted_at", time.Now())

// Or use GORM's soft delete
db.Delete(&user)  // Sets deleted_at automatically if using gorm.DeletedAt

// Restore soft deleted record
db.Model(&user).Update("deleted_at", nil)

// Permanently delete
db.Unscoped().Delete(&user)

// Query including soft deleted
db.Unscoped().Find(&users)
```

### UUID Pattern

Every record has both `id` (internal) and `uuid` (public):

```go
// Use ID for internal operations (foreign keys, joins)
db.Where("user_id = ?", user.ID).Find(&orders)

// Use UUID for external API responses
type UserResponse struct {
    UUID  string `json:"uuid"`  // Expose this
    Name  string `json:"name"`
    // Never expose ID in API
}

// Generate UUID on create
import "github.com/google/uuid"

user := &models.User{
    UUID: uuid.New().String(),
    Name: "John Doe",
}
db.Create(user)
```

### Audit Trail Pattern — Resolving JWT UUID to Internal User ID for `CreatedBy`

The auth middleware sets `userID` in the gin context as a UUID string (from the JWT `uid` claim). But `BaseModel.CreatedBy` is `*uint64` (the internal auto-increment ID). You must resolve UUID → internal ID before using it.

**Handler layer** — extract UUID using `utils.GetUserID(c)` and pass it to the service as a string:

```go
func (h *ClientHandler) CreateClient(c *gin.Context) {
    var req dto.CreateClientRequest
    if !bind.AndValidate(c, &req, "create_client") {
        return
    }

    userID := utils.GetUserID(c) // returns UUID string from context

    client, err := h.ClientService.Create(req, userID)
    // ...
}
```

**Service layer** — resolve UUID → internal `int64` ID via Redis cache (falls back to DB), then convert to `uint64` for `CreatedBy`:

```go
func (s *ClientService) Create(req dto.CreateClientRequest, userUUID string) (*dto.ClientResponse, error) {
    var createdBy *uint64

    if userUUID != "" {
        userID, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
        if err != nil || userID == nil {
            return nil, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "client_service", "get_user_id")
        }
        uid := uint64(*userID)
        createdBy = &uid
    }

    client := &models.Client{
        Name:      req.Name,
        BaseModel: models.NewBase(createdBy),
    }
    // ...
}
```

**Service must have `UserRepo` and `RedisRepo` injected** — wire them in `internal/container/service.go`:

```go
ClientService: services.NewClientService(repository.ClientRepo, repository.UserRepo, repository.CacheRepo),
```

**Key rules:**
- NEVER cast `c.Get("userID")` directly to `uint64` — it's a UUID string, not an integer.
- NEVER use `c.Get("user_id")` — the context key is `"userID"` (set by auth middleware).
- ALWAYS resolve UUID → ID through `redisRepo.GetUserIdByUuid()` for cache-first performance.
- See `internal/services/file.go` → `CreateFileAndGetUploadURL()` for the reference implementation of this pattern.

### Auto-Update Timestamps

Database triggers handle `updated_at` automatically:

```sql
CREATE TRIGGER `users_updated_at` BEFORE UPDATE ON `users`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END;
```

**In GORM models:**
```go
type BaseModel struct {
    CreatedAt time.Time  // Auto-set on create
    UpdatedAt time.Time  // Auto-updated by trigger
    DeletedAt time.Time  // Set manually or by GORM soft delete
}
```

---

## Utility Functions

### Password Hashing (`internal/utils/hash.go`)

```go
package utils

import "golang.org/x/crypto/bcrypt"

// Hash a password
func HashPassword(p string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(p), 14)
    return string(hash)
}

// Verify password
func CheckPassword(hash, p string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}
```

**Usage in service:**
```go
// On user registration
func (s *UserService) Register(req RegisterRequest) error {
    hash := utils.HashPassword(req.Password)
    
    user := &models.User{
        Email:        req.Email,
        PasswordHash: hash,
    }
    
    return s.repo.Create(user)
}

// On login
func (s *AuthService) Login(req LoginRequest) error {
    user, _ := s.repo.FindByEmail(req.Email)
    
    if !utils.CheckPassword(user.PasswordHash, req.Password) {
        return domainerrors.ErrInvalidCredentials
    }
    
    // Generate token...
}
```

### UUID Generation

```go
import "github.com/google/uuid"

// Generate new UUID
newUUID := uuid.New().String()

// Use in model creation
user := &models.User{
    UUID: uuid.New().String(),
    Name: "John Doe",
}
```

---

## Common Repository Patterns

### Find by UUID (Public API)

```go
func (r *UserRepo) FindByUUID(uuid string) (*models.User, error) {
    var user models.User
    err := r.DB.Where("uuid = ? AND deleted_at IS NULL", uuid).First(&user).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    
    return &user, err
}
```

### Find by ID (Internal)

```go
func (r *UserRepo) FindByID(id uint64) (*models.User, error) {
    var user models.User
    err := r.DB.Where("deleted_at IS NULL").First(&user, id).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    
    return &user, err
}
```

### List with Pagination

```go
func (r *UserRepo) List(page, limit int, status string) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    
    query := r.DB.Where("deleted_at IS NULL")
    
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Get total count
    query.Model(&models.User{}).Count(&total)
    
    // Get paginated results
    offset := (page - 1) * limit
    err := query.Limit(limit).Offset(offset).Find(&users).Error
    
    return users, total, err
}
```

### Create with UUID

```go
func (r *UserRepo) Create(user *models.User) error {
    // Generate UUID if not set
    if user.UUID == "" {
        user.UUID = uuid.New().String()
    }
    
    // Set default status
    if user.Status == "" {
        user.Status = "active"
    }
    
    return r.DB.Create(user).Error
}
```

### Update (Partial)

```go
func (r *UserRepo) Update(id uint64, updates map[string]interface{}) error {
    return r.DB.Model(&models.User{}).
        Where("id = ? AND deleted_at IS NULL", id).
        Updates(updates).Error
}
```

### Soft Delete

```go
func (r *UserRepo) Delete(id uint64) error {
    return r.DB.Model(&models.User{}).
        Where("id = ?", id).
        Update("deleted_at", time.Now()).Error
}
```

### Check Exists

```go
func (r *UserRepo) ExistsByEmail(email string) (bool, error) {
    var count int64
    err := r.DB.Model(&models.User{}).
        Where("email = ? AND deleted_at IS NULL", email).
        Count(&count).Error
    
    return count > 0, err
}
```

---

## Common Service Patterns

### Create with Validation

```go
func (s *UserService) Create(req CreateUserRequest) (*UserResponse, error) {
    // Check if email exists
    exists, err := s.repo.ExistsByEmail(req.Email)
    if err != nil {
        return nil, domainerrors.ErrInternal
    }
    if exists {
        return nil, domainerrors.ErrDuplicateEmail
    }
    
    // Hash password
    hash := utils.HashPassword(req.Password)
    
    // Create user
    user := &models.User{
        UUID:         uuid.New().String(),
        Name:         req.Name,
        Email:        req.Email,
        PasswordHash: hash,
        Status:       "active",
    }
    
    if err := s.repo.Create(user); err != nil {
        return nil, domainerrors.ErrInternal
    }
    
    // Return response DTO
    return &UserResponse{
        UUID:      user.UUID,
        Name:      user.Name,
        Email:     user.Email,
        Status:    user.Status,
        CreatedAt: user.CreatedAt,
    }, nil
}
```

### Update with Validation

```go
func (s *UserService) Update(uuid string, req UpdateUserRequest) (*UserResponse, error) {
    // Find user
    user, err := s.repo.FindByUUID(uuid)
    if err != nil {
        return nil, domainerrors.ErrInternal
    }
    if user == nil {
        return nil, domainerrors.ErrUserNotFound
    }
    
    // Check email uniqueness if changing
    if req.Email != nil && *req.Email != user.Email {
        exists, _ := s.repo.ExistsByEmail(*req.Email)
        if exists {
            return nil, domainerrors.ErrDuplicateEmail
        }
    }
    
    // Build updates map
    updates := make(map[string]interface{})
    if req.Name != nil {
        updates["name"] = *req.Name
    }
    if req.Email != nil {
        updates["email"] = *req.Email
    }
    
    // Update
    if err := s.repo.Update(user.ID, updates); err != nil {
        return nil, domainerrors.ErrInternal
    }
    
    // Fetch updated user
    user, _ = s.repo.FindByID(user.ID)
    
    return &UserResponse{
        UUID:      user.UUID,
        Name:      user.Name,
        Email:     user.Email,
        Status:    user.Status,
        UpdatedAt: user.UpdatedAt,
    }, nil
}
```

### Get by UUID

```go
func (s *UserService) GetByUUID(uuid string) (*UserResponse, error) {
    user, err := s.repo.FindByUUID(uuid)
    if err != nil {
        return nil, domainerrors.ErrInternal
    }
    if user == nil {
        return nil, domainerrors.ErrUserNotFound
    }
    
    return &UserResponse{
        UUID:      user.UUID,
        Name:      user.Name,
        Email:     user.Email,
        Status:    user.Status,
        CreatedAt: user.CreatedAt,
    }, nil
}
```

### List with Pagination

```go
func (s *UserService) List(page, limit int, status string) (*UserListResponse, error) {
    users, total, err := s.repo.List(page, limit, status)
    if err != nil {
        return nil, domainerrors.ErrInternal
    }
    
    // Convert to response DTOs
    userResponses := make([]UserResponse, len(users))
    for i, user := range users {
        userResponses[i] = UserResponse{
            UUID:      user.UUID,
            Name:      user.Name,
            Email:     user.Email,
            Status:    user.Status,
            CreatedAt: user.CreatedAt,
        }
    }
    
    return &UserListResponse{
        Users: userResponses,
        Meta: response.PageMeta{
            Page:    page,
            Limit:   limit,
            Total:   total,
            Pages:   int(math.Ceil(float64(total) / float64(limit))),
            HasNext: page*limit < int(total),
        },
    }, nil
}
```

### Delete (Soft)

```go
func (s *UserService) Delete(uuid string) error {
    user, err := s.repo.FindByUUID(uuid)
    if err != nil {
        return domainerrors.ErrInternal
    }
    if user == nil {
        return domainerrors.ErrUserNotFound
    }
    
    return s.repo.Delete(user.ID)
}
```

---

## Common Handler Patterns

### Create

```go
func (h *UserHandler) Create(c *gin.Context) {
    var req userdto.CreateUserRequest
    if !bind.AndValidate(c, &req, "user") {
        return
    }
    
    user, err := h.service.Create(req)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    lang := c.GetHeader("Accept-Language")
    msgs := i18n.GetMessages(lang)
    
    response.JSON(c, 201, response.APIResponse{
        Success: true,
        Message: msgs["user_created_success"],
        Data:    user,
    })
}
```

### Get by UUID

```go
func (h *UserHandler) GetByUUID(c *gin.Context) {
    uuid := c.Param("uuid")
    
    user, err := h.service.GetByUUID(uuid)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    response.JSON(c, 200, response.APIResponse{
        Success: true,
        Data:    user,
    })
}
```

### List with Pagination

```go
func (h *UserHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    status := c.Query("status")
    
    result, err := h.service.List(page, limit, status)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    response.JSON(c, 200, response.APIResponse{
        Success: true,
        Data:    result.Users,
        Meta:    result.Meta,
    })
}
```

### Update

```go
func (h *UserHandler) Update(c *gin.Context) {
    uuid := c.Param("uuid")
    
    var req userdto.UpdateUserRequest
    if !bind.AndValidate(c, &req, "user") {
        return
    }
    
    user, err := h.service.Update(uuid, req)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    lang := c.GetHeader("Accept-Language")
    msgs := i18n.GetMessages(lang)
    
    response.JSON(c, 200, response.APIResponse{
        Success: true,
        Message: msgs["user_updated_success"],
        Data:    user,
    })
}
```

### Delete

```go
func (h *UserHandler) Delete(c *gin.Context) {
    uuid := c.Param("uuid")
    
    err := h.service.Delete(uuid)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    lang := c.GetHeader("Accept-Language")
    msgs := i18n.GetMessages(lang)
    
    response.JSON(c, 200, response.APIResponse{
        Success: true,
        Message: msgs["user_deleted_success"],
    })
}
```

---

## Naming Conventions

### Files
- Models: `internal/domain/models/user.go`
- DTOs: `internal/domain/dto/user/user.go`
- Repositories: `internal/repositories/user.go`
- Services: `internal/services/user.go`
- Handlers: `internal/handlers/user_handler.go`
- Routes: `internal/routes/version/v1/user.go`

### Structs
- Model: `User`
- Repository: `UserRepo`
- Service: `UserService`
- Handler: `UserHandler`
- Request DTO: `CreateUserRequest`, `UpdateUserRequest`
- Response DTO: `UserResponse`, `UserListResponse`

### Functions
- Repository: `FindByID`, `FindByUUID`, `Create`, `Update`, `Delete`, `List`
- Service: `Create`, `GetByUUID`, `Update`, `Delete`, `List`
- Handler: `Create`, `GetByUUID`, `Update`, `Delete`, `List`

### Routes
- GET `/api/v1/users` - List
- GET `/api/v1/users/:uuid` - Get one
- POST `/api/v1/users` - Create
- PUT `/api/v1/users/:uuid` - Update
- DELETE `/api/v1/users/:uuid` - Delete

---

## Best Practices Checklist

- [ ] All models embed `BaseModel`
- [ ] All tables have `id`, `uuid`, `status`, timestamps, `created_by`
- [ ] Use UUID in API routes, not ID
- [ ] Use soft delete (`deleted_at`), not hard delete
- [ ] Hash passwords with bcrypt (cost 14)
- [ ] Validate uniqueness before create/update
- [ ] Return domain errors from services
- [ ] Use `bind.AndValidate()` in handlers
- [ ] Use `response.FromError()` for errors
- [ ] Use i18n for all messages
- [ ] Check `deleted_at IS NULL` in queries
- [ ] Generate UUID on create
- [ ] Set `created_by` from auth context
- [ ] Use pagination for list endpoints
- [ ] Return DTOs, never models
- [ ] Never expose sensitive fields (passwords, hashes)
