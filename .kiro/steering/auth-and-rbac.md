---
inclusion: auto
description: Authentication, authorization, RBAC - token types, middleware usage, role system, context values, route protection patterns
---

# Authentication & Role-Based Access Control (RBAC)

This project uses JWT-based authentication with separate token types for admin and user roles, plus RBAC middleware for fine-grained access control.

---

## Token System

There are two token types with separate secrets:

- **User Token** (`TokenTypeUser`) — signed with `JWT_SECRET`. Issued to employees (role_id=3 Manager, role_id=4 Employee).
- **Admin Token** (`TokenTypeAdmin`) — signed with `JWT_ADMIN_SECRET`. Issued to users with role_id=1 (Super Admin) or role_id=2 (Admin).

Token type is determined at login time based on the user's roles in `user_roles` table. See `internal/services/auth.go` → `CreateAndReturnToken()`.

Each token also has a corresponding refresh token type (`TokenTypeRefresh`, `TokenTypeAdminRefresh`) with separate secrets.

---

## Role System (Database)

Roles are seeded in `internal/bootstrap/seeder.go`:

| role_id | Name         | Token Type  |
|---------|-------------|-------------|
| 1       | Super Admin | Admin token |
| 2       | Admin       | Admin token |
| 3       | Manager     | User token  |
| 4       | Employee    | User token  |

A user can have multiple roles via the `user_roles` join table. The `IsAdmin` check looks for role_id IN (1, 2). The `IsManager` check looks for role_id = 3. The `IsManagerOrAdmin` check looks for role_id IN (1, 2, 3).

---

## Middleware Overview

All middleware lives in `internal/middleware/` and is wired through `internal/container/middleware.go`.

### Available Middleware

| Middleware | Method | Purpose |
|---|---|---|
| `AdminMiddleware` | `VerifyAdminToken()` | Accepts ONLY admin tokens (role 1, 2). Rejects user tokens. |
| `UserMiddleware` | `VerifyUserToken()` | Accepts ONLY user tokens (role 3, 4). Rejects admin tokens. |
| `AuthMiddleware` | `VerifyToken()` | Accepts BOTH admin and user tokens. Sets `isAdmin` flag in context. |
| `RABCMiddleware` | `IsManagerOrAdmin()` | Must be chained AFTER `AuthMiddleware.VerifyToken()`. Checks if the authenticated user has role 1, 2, or 3. |

### What Each Middleware Sets in Context

**AdminMiddleware.VerifyAdminToken():**
```go
c.Set("userID", claims.UserID)   // user UUID string
```

**UserMiddleware.VerifyUserToken():**
```go
c.Set("userID", claims.UserID)   // user UUID string
c.Set("jti", claims.ID)          // JWT ID for session tracking
```

**AuthMiddleware.VerifyToken():**
```go
c.Set("userID", claims.UserID)   // user UUID string
c.Set("isAdmin", isAdmin)        // bool — true if admin token, false if user token
c.Set("jti", claims.ID)          // JWT ID for session tracking
```

**RABCMiddleware.IsManagerOrAdmin():**
Sets nothing extra. Only gates access — aborts with 401 if user doesn't have role 1, 2, or 3.

---

## Reading Context Values

Use the utility functions in `internal/utils/utils.go`:

```go
userID := utils.GetUserID(c)   // returns string (user UUID from token)
jti := utils.GetJTI(c)         // returns string (JWT ID)
```

For checking admin status (only available after `AuthMiddleware.VerifyToken()`):
```go
isAdmin := c.GetBool("isAdmin")
```

---

## Route Protection Patterns

### Pattern 1: Admin-Only Routes

For routes that ONLY admins (Super Admin + Admin) can access. Uses a separate admin route group.

```go
// internal/routes/version/v1/admin/admin.go
func SetAdminRoute(v1 *gin.RouterGroup, c *container.Container) {
    admin := v1.Group("/admin")
    admin.Use(c.Middleware.AdminMiddleware.VerifyAdminToken())
    
    UserRoutes(admin, c)  // All routes under /admin require admin token
}
```

Result: `POST /api/v1/admin/users`, `GET /api/v1/admin/users`, etc.

### Pattern 2: Any Authenticated User (Admin OR User)

For routes where both admin and regular users can access. Uses `AuthMiddleware.VerifyToken()`.

```go
// internal/routes/version/v1/auth.go
func AuthRoutes(v1 *gin.RouterGroup, c *container.Container) {
    auth := v1.Group("/auth")
    auth.GET("/me", c.Middleware.AuthMiddleware.VerifyToken(), h.Me)
    auth.POST("/logout", c.Middleware.AuthMiddleware.VerifyToken(), h.Logout)
}
```

Inside the handler, check `c.GetBool("isAdmin")` if you need to differentiate behavior.

### Pattern 3: Manager + Admin Only (No Employee)

For routes where managers and admins can access, but regular employees cannot. Chain `AuthMiddleware.VerifyToken()` then `RABCMiddleware.IsManagerOrAdmin()`.

```go
// internal/routes/version/v1/projects.go
func Projects(v1 *gin.RouterGroup, c *container.Container) {
    project := v1.Group("/project")
    project.Use(c.Middleware.AuthMiddleware.VerifyToken())
    project.Use(c.Middleware.RABCMiddleware.IsManagerOrAdmin())
    
    // Only role 1, 2, 3 can access these routes
}
```

### Pattern 4: User-Only Routes (No Admin)

For routes that only regular users (Manager + Employee) can access:

```go
func SomeRoutes(v1 *gin.RouterGroup, c *container.Container) {
    group := v1.Group("/something")
    group.Use(c.Middleware.UserMiddleware.VerifyUserToken())
    
    // Only user token holders can access
}
```

---

## How to Choose the Right Middleware

| Who should access? | Middleware to use |
|---|---|
| Only Super Admin + Admin | `AdminMiddleware.VerifyAdminToken()` |
| Only Manager + Employee | `UserMiddleware.VerifyUserToken()` |
| Everyone authenticated | `AuthMiddleware.VerifyToken()` |
| Admin + Manager (no Employee) | `AuthMiddleware.VerifyToken()` + `RABCMiddleware.IsManagerOrAdmin()` |

---

## Adding a New RBAC Check

If you need a new role-based check (e.g., "only managers"):

1. Add the repo method in `internal/repositories/user.go`:
```go
func (r *UserRepo) IsManager(uuid string) (bool, error) {
    result, err := r.FindByUuidAndMailHasRoles(uuid, []int{3})
    return result, err
}
```

2. Add the middleware method in `internal/middleware/rbac.go`:
```go
func (m *RABCMiddleware) IsManager() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := utils.GetUserID(c)
        is, err := m.UserRepo.IsManager(userID)
        if err != nil || !is {
            response.FromError(c, domainerrors.NewWithCode(http.StatusUnauthorized, domainerrors.ErrUnauthorized.Error(), "check_role", "manager_only"))
            c.Abort()
            return
        }
        c.Next()
    }
}
```

3. Use in routes (always AFTER `AuthMiddleware.VerifyToken()`):
```go
group.Use(c.Middleware.AuthMiddleware.VerifyToken())
group.Use(c.Middleware.RABCMiddleware.IsManager())
```

---

## Session Verification Flow

Every middleware (Admin, User, Auth) performs these checks in order:

1. Extract Bearer token from `Authorization` header
2. Validate JWT signature and expiry using the appropriate secret
3. Check user status via `verifyUserStatus()` — uses Redis cache with key `user_status_{uuid}`
4. Check session is active via `session:jti:{jti}` Redis cache key, falls back to DB `sessions` table
5. Set context values (`userID`, `isAdmin`, `jti`)

If any check fails, the middleware returns 401 and aborts the request chain.

---

## Token Claims Structure

```go
type Claims struct {
    UserID string `json:"uid"`    // User UUID (NOT the integer ID)
    jwt.RegisteredClaims          // Includes ID (JTI), ExpiresAt, IssuedAt
}
```

The `UserID` in claims is always the user's UUID string, not the database integer ID. The `ID` field (from RegisteredClaims) is the JTI used for session tracking.
