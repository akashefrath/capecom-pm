---
inclusion: auto
description: Routing pattern - modular versioned routes, how to create route modules and RESTful endpoints
---

# Routing Pattern

This project uses a modular, versioned routing structure with dependency injection from the container.

## Routing Architecture

```
routes.go (main setup)
    ↓
version/v1/v1.go (version aggregator)
    ↓
version/v1/auth.go, user.go, project.go, etc. (feature modules)
```

## Current Structure

```
internal/routes/
├── routes.go              # Main router setup, global rate limiter, calls v1.Routes()
└── version/
    └── v1/
        ├── v1.go          # V1 aggregator, registers all v1 modules
        ├── auth.go        # Auth routes (login, refresh, me, logout)
        ├── projects.go    # Project routes (AuthMiddleware + RBAC)
        └── admin/
            ├── admin.go   # Admin route group (AdminMiddleware.VerifyAdminToken)
            └── user.go    # Admin user management routes
```

### Route Groups by Access Level

- `/api/v1/admin/*` — Admin-only (uses `AdminMiddleware.VerifyAdminToken()`)
- `/api/v1/auth/*` — Mixed (some public, some use `AuthMiddleware.VerifyToken()`)
- `/api/v1/project/*` — Manager + Admin only (uses `AuthMiddleware.VerifyToken()` + `RABCMiddleware.IsManagerOrAdmin()`)

See `auth-and-rbac.md` for full middleware documentation.

## How Routes Work

### 1. Main Router Setup (`routes.go`)

```go
package routes

import (
    "capecom-pm/internal/container"
    "capecom-pm/internal/routes/version/v1"
    "github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, c *container.Container) {
    // Global routes (health checks, etc.)
    r.Any("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"msg": "pong"})
    })

    // API version groups
    apiV1 := r.Group("/api/v1")
    v1.Routes(apiV1, c)
}
```

### 2. Version Aggregator (`version/v1/v1.go`)

This file registers all feature modules for v1:

```go
package v1

import (
    "capecom-pm/internal/container"
    "github.com/gin-gonic/gin"
)

func Routes(v1 *gin.RouterGroup, c *container.Container) {
    AuthRoutes(v1, c)
    UserRoutes(v1, c)      // Add new modules here
    ProjectRoutes(v1, c)   // Add new modules here
}
```

### 3. Feature Route Module (`version/v1/auth.go`)

Each feature has its own route file:

```go
package v1

import (
    "capecom-pm/internal/container"
    "github.com/gin-gonic/gin"
)

func AuthRoutes(v1 *gin.RouterGroup, c *container.Container) {
    // Get handler from DI container
    h := c.Handler.AuthHandler
    
    // Create route group
    auth := v1.Group("/auth")
    
    // Register routes
    auth.POST("/login", h.Login)
    auth.POST("/register", h.Register)
    auth.POST("/logout", h.Logout)
}
```

## How to Add a New Route Module

### Step 1: Create Feature Route File

**File:** `internal/routes/version/v1/user.go`

```go
package v1

import (
    "capecom-pm/internal/container"
    "github.com/gin-gonic/gin"
)

func UserRoutes(v1 *gin.RouterGroup, c *container.Container) {
    // Get handler from DI container
    h := c.Handler.UserHandler
    
    // Create route group
    user := v1.Group("/users")
    
    // Register routes
    user.GET("", h.GetAll)           // GET /api/v1/users
    user.GET("/:id", h.GetByID)      // GET /api/v1/users/:id
    user.POST("", h.Create)          // POST /api/v1/users
    user.PUT("/:id", h.Update)       // PUT /api/v1/users/:id
    user.DELETE("/:id", h.Delete)    // DELETE /api/v1/users/:id
}
```

### Step 2: Register in Version Aggregator

**File:** `internal/routes/version/v1/v1.go`

```go
func Routes(v1 *gin.RouterGroup, c *container.Container) {
    AuthRoutes(v1, c)
    UserRoutes(v1, c)      // Add this line
}
```

## Route Patterns

### Basic CRUD Routes

```go
func ResourceRoutes(v1 *gin.RouterGroup, c *container.Container) {
    h := c.Handler.ResourceHandler
    resource := v1.Group("/resources")
    
    resource.GET("", h.List)              // List all
    resource.GET("/:id", h.Get)           // Get one
    resource.POST("", h.Create)           // Create
    resource.PUT("/:id", h.Update)        // Update
    resource.DELETE("/:id", h.Delete)     // Delete
}
```

### Nested Routes

```go
func ProjectRoutes(v1 *gin.RouterGroup, c *container.Container) {
    h := c.Handler.ProjectHandler
    th := c.Handler.TaskHandler
    
    project := v1.Group("/projects")
    project.GET("", h.List)
    project.POST("", h.Create)
    
    // Nested tasks under projects
    project.GET("/:id/tasks", th.GetProjectTasks)
    project.POST("/:id/tasks", th.CreateTask)
}
```

### Protected Routes with Middleware

```go
func UserRoutes(v1 *gin.RouterGroup, c *container.Container) {
    h := c.Handler.UserHandler
    
    user := v1.Group("/users")
    
    // Public routes
    user.GET("", h.List)
    
    // Protected routes (require auth middleware)
    user.Use(middleware.AuthRequired())
    user.POST("", h.Create)
    user.PUT("/:id", h.Update)
    user.DELETE("/:id", h.Delete)
}
```

## Key Principles

1. **One File Per Feature Module**: Each feature (auth, user, project) gets its own route file

2. **DI Container Access**: Always get handlers from `c.Handler.YourHandler`

3. **Route Grouping**: Use `v1.Group("/resource")` to create logical groups

4. **Version Aggregation**: Register all modules in `v1.go`

5. **RESTful Conventions**: Follow REST patterns for consistency
   - GET for retrieval
   - POST for creation
   - PUT/PATCH for updates
   - DELETE for deletion

6. **URL Structure**: `/api/v1/resource` or `/api/v1/resource/:id`

## Complete Example: Adding User Module

### 1. Create Handler (see dependency-injection.md)
```go
// internal/handlers/user_handler.go
type UserHandler struct {
    service *services.UserService
}
```

### 2. Register in Container
```go
// internal/container/handler.go
type Handler struct {
    AuthHandler *handlers.AuthHandler
    UserHandler *handlers.UserHandler  // Add this
}
```

### 3. Create Route Module
```go
// internal/routes/version/v1/user.go
func UserRoutes(v1 *gin.RouterGroup, c *container.Container) {
    h := c.Handler.UserHandler
    user := v1.Group("/users")
    user.GET("", h.GetAll)
    user.POST("", h.Create)
}
```

### 4. Register in v1.go
```go
// internal/routes/version/v1/v1.go
func Routes(v1 *gin.RouterGroup, c *container.Container) {
    AuthRoutes(v1, c)
    UserRoutes(v1, c)  // Add this
}
```

## Testing Routes

You can test routes with curl:

```bash
# Auth routes
curl -X POST http://localhost:8080/api/v1/auth/login

# User routes
curl -X GET http://localhost:8080/api/v1/users
curl -X POST http://localhost:8080/api/v1/users
```

## Future: Adding v2

When you need API v2:

1. Create `internal/routes/version/v2/` directory
2. Create `v2.go` aggregator
3. Add feature modules
4. Register in `routes.go`:

```go
apiV2 := r.Group("/api/v2")
v2.Routes(apiV2, c)
```
