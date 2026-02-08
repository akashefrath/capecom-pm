---
inclusion: manual
---

# Project Steering Files

This directory contains comprehensive documentation for the project's architecture, patterns, and conventions. These files guide development to ensure consistency across the codebase.

## 📚 Available Steering Files

### 1. **dependency-injection.md**
**Purpose:** Explains the DI container pattern used throughout the project.

**Covers:**
- Architecture layers (DB → Repository → Service → Handler)
- How to add new repositories, services, and handlers
- Constructor injection pattern
- Container wiring

**When to use:** Adding any new feature that needs data access, business logic, or HTTP handlers.

---

### 2. **routing-pattern.md**
**Purpose:** Explains the modular, versioned routing structure.

**Covers:**
- Route organization (routes.go → v1/v1.go → v1/feature.go)
- How to create new route modules
- RESTful conventions
- Route grouping and middleware
- API versioning

**When to use:** Adding new API endpoints or route modules.

---

### 3. **domain-structure.md**
**Purpose:** Explains models, DTOs, and domain organization.

**Covers:**
- BaseModel pattern (all models must embed it)
- Model structure and rules
- DTO organization (request/response)
- Validation tags
- Type safety and security

**When to use:** Creating new database models or API request/response structures.

---

### 4. **error-handling-flow.md**
**Purpose:** Explains error handling, translations, and API responses.

**Covers:**
- Complete request flow (Handler → Service → Repository)
- Domain errors pattern
- Error mapper (domain errors → HTTP codes)
- i18n translations (EN/TA)
- Validation errors
- API response format

**When to use:** Handling errors, adding translations, or returning API responses.

---

### 5. **configuration-pattern.md**
**Purpose:** Explains environment variable management and configuration.

**Covers:**
- Config structure (grouped by domain)
- Environment variable loading
- Required vs optional values
- Database connection setup
- How to add new configuration

**When to use:** Adding new configuration values or external service connections.

---

### 6. **common-patterns.md**
**Purpose:** Common utilities, database patterns, and best practices.

**Covers:**
- Database patterns (BaseModel, soft delete, UUID, status enum)
- Utility functions (password hashing, UUID generation)
- Common repository patterns (CRUD, pagination)
- Common service patterns (validation, error handling)
- Common handler patterns (request/response)
- Naming conventions
- Best practices checklist

**When to use:** Implementing CRUD operations, working with database, or following project conventions.

---

## 🚀 Quick Start Guide

### Adding a New Feature (e.g., "Projects")

Follow these steps in order:

#### 1. Create Model
**Reference:** `domain-structure.md`, `common-patterns.md`

```go
// internal/domain/models/project.go
type Project struct {
    BaseModel  // Always embed BaseModel

    Name        string
    Description string
    ClientID    uint64
}
```

#### 2. Create DTOs
**Reference:** `domain-structure.md`

```go
// internal/domain/dto/project/project.go
type CreateProjectRequest struct {
    Name        string `json:"name" binding:"required,min=2"`
    Description string `json:"description" binding:"required"`
}

type ProjectResponse struct {
    UUID        string `json:"uuid"`
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

#### 3. Create Repository
**Reference:** `dependency-injection.md`, `common-patterns.md`

```go
// internal/repositories/project.go
type ProjectRepo struct {
    DB *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo {
    return &ProjectRepo{DB: db}
}

func (r *ProjectRepo) Create(project *models.Project) error {
    return r.DB.Create(project).Error
}
```

#### 4. Register Repository in Container
**Reference:** `dependency-injection.md`

```go
// internal/container/repository.go
type Repository struct {
    AuthRepo    *repositories.AuthRepo
    ProjectRepo *repositories.ProjectRepo  // Add this
}
```

#### 5. Create Service
**Reference:** `dependency-injection.md`, `error-handling-flow.md`

```go
// internal/services/project.go
type ProjectService struct {
    repo *repositories.ProjectRepo
}

func NewProjectService(repo *repositories.ProjectRepo) *ProjectService {
    return &ProjectService{repo: repo}
}
```

#### 6. Register Service in Container
**Reference:** `dependency-injection.md`

```go
// internal/container/service.go
type Service struct {
    AuthService    *services.AuthService
    ProjectService *services.ProjectService  // Add this
}
```

#### 7. Create Handler
**Reference:** `dependency-injection.md`, `error-handling-flow.md`

```go
// internal/handlers/project_handler.go
type ProjectHandler struct {
    service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
    return &ProjectHandler{service: service}
}
```

#### 8. Register Handler in Container
**Reference:** `dependency-injection.md`

```go
// internal/container/handler.go
type Handler struct {
    AuthHandler    *handlers.AuthHandler
    ProjectHandler *handlers.ProjectHandler  // Add this
}
```

#### 9. Create Routes
**Reference:** `routing-pattern.md`

```go
// internal/routes/version/v1/project.go
func ProjectRoutes(v1 *gin.RouterGroup, c *container.Container) {
    h := c.Handler.ProjectHandler
    project := v1.Group("/projects")
    
    project.GET("", h.List)
    project.POST("", h.Create)
}
```

#### 10. Register Routes
**Reference:** `routing-pattern.md`

```go
// internal/routes/version/v1/v1.go
func Routes(v1 *gin.RouterGroup, c *container.Container) {
    AuthRoutes(v1, c)
    ProjectRoutes(v1, c)  // Add this
}
```

#### 11. Add Translations
**Reference:** `error-handling-flow.md`

```go
// internal/utils/i18n/en.go
"project_created_success": "Project created successfully",
"project_not_found": "Project not found",

// internal/utils/i18n/ta.go
"project_created_success": "திட்டம் வெற்றிகரமாக உருவாக்கப்பட்டது",
"project_not_found": "திட்டம் கிடைக்கவில்லை",
```

---

## 🎯 Key Principles

1. **Consistency:** Follow the same patterns across all features
2. **Separation of Concerns:** Each layer has one responsibility
3. **Type Safety:** Always use structs with validation
4. **Security:** Never expose sensitive data or internal IDs
5. **i18n:** All messages must be translatable
6. **Error Handling:** Use domain errors, never raw strings
7. **Documentation:** Code should be self-documenting

---

## 📋 Development Checklist

Before submitting code, ensure:

- [ ] Model embeds `BaseModel`
- [ ] DTOs have proper validation tags
- [ ] Repository returns `nil, nil` for not found
- [ ] Service returns domain errors
- [ ] Handler uses `bind.AndValidate()`
- [ ] Handler uses `response.FromError()`
- [ ] All messages use i18n
- [ ] Routes use UUID, not ID
- [ ] API responses use DTOs, not models
- [ ] Sensitive fields are not exposed
- [ ] Code follows naming conventions
- [ ] Translations added for EN and TA

---

## 🔍 Finding Information

**Need to know how to...**

- Add a new feature? → Start with `dependency-injection.md`
- Create API endpoints? → `routing-pattern.md`
- Define models/DTOs? → `domain-structure.md`
- Handle errors? → `error-handling-flow.md`
- Add config values? → `configuration-pattern.md`
- Implement CRUD? → `common-patterns.md`

---

## 📝 Notes

- All steering files are automatically included in AI context
- These files are living documents - update them as patterns evolve
- When in doubt, check existing code (auth module is the reference)
- Consistency is more important than perfection

---

**Last Updated:** February 8, 2026
**Project:** Capecom PM (Project Management System)
**Tech Stack:** Go, Gin, GORM, MySQL
