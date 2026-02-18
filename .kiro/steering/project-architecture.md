---
inclusion: always
---

# PROJECT ARCHITECTURE FLOW

## EXECUTION FLOW
cmd/main.go → internal/routes → internal/routes/v1 → handler → service → repository

## 1. ENTRY POINT (cmd/main.go)
- Initialize container with DB, config, redis
- Setup routes with `routes.Setup(r, appContainer)`
- Start server

## 2. ROUTES (internal/routes/routes.go)
- Call v1.Setup for API v1 routes

## 3. V1 ROUTES (internal/routes/v1/)
- Create route groups (auth, utils, etc.)
- Access handlers via `container.Handler.{HandlerName}`
- Register endpoints: `group.POST("endpoint", handler.Method)`
- Apply middleware: `group.Use(container.Middleware.Auth.VerifyToken(...))`

Example:
```go
func Feature(r *gin.RouterGroup, container container.Container) {
    feature := r.Group("feature")
    handler := container.Handler.FeatureHandler
    feature.POST("create", handler.Create)
}
```

## 4. HANDLER (internal/src/handler/)
- Bind request with `bind.AndValidate(c, &req, "key")`
- Call service method
- Send response with `response.JSONOk(c, response.APIResponse{...})`
- Handle errors with `response.FromError(c, err)`

Example:
```go
type FeatureHandler struct {
    Feature *service.Feature
}

func NewFeature(feature *service.Feature) FeatureHandler {
    return FeatureHandler{Feature: feature}
}

func (h *FeatureHandler) Create(c *gin.Context) {
    var req dto.CreateRequest
    isValid := bind.AndValidate(c, &req, "key")
    if !isValid {
        return
    }
    data, err := h.Feature.Create(req)
    if err != nil {
        response.FromError(c, err)
        return
    }
    response.JSONOk(c, response.APIResponse{
        Success: true,
        Data: data,
    })
}
```

## 5. SERVICE (internal/src/service/)
- Business logic
- Call repository methods
- Return data or error

Example:
```go
type Feature struct {
    FeatureRepo *repository.Feature
}

func NewFeature(featureRepo *repository.Feature) *Feature {
    return &Feature{FeatureRepo: featureRepo}
}

func (s *Feature) Create(req dto.CreateRequest) (*dto.Response, error) {
    return s.FeatureRepo.Create(req)
}
```

## 6. REPOSITORY (internal/src/repository/)
- Use RAW SQL queries for performance and security
- Use `r.DB.Get()` for single row
- Use `r.DB.Select()` for multiple rows
- Use `r.DB.Exec()` for insert/update/delete

Example:
```go
type Feature struct {
    DB *sqlx.DB
}

func NewFeature(db *sqlx.DB) *Feature {
    return &Feature{DB: db}
}

func (r *Feature) Create(data dto.CreateRequest) (*dto.Response, error) {
    q := `INSERT INTO table (field) VALUES (?)`
    _, err := r.DB.Exec(q, data.Field)
    return nil, err
}
```

## 7. DEPENDENCY INJECTION (internal/container/)

### container.go
- Main container struct with all dependencies

### repositories.go
- Initialize all repositories
- Add to Repositories struct

Example:
```go
type Repositories struct {
    Feature *repository.Feature
}

func NewRepository(db *sqlx.DB, config *config.Config, redis *redis.Client, dbTX *database.Database) *Repositories {
    return &Repositories{
        Feature: repository.NewFeature(db),
    }
}
```

### service.go
- Initialize all services with repository dependencies
- Add to Service struct

Example:
```go
type Service struct {
    Feature *service.Feature
}

func NewService(repo *Repositories, jwtManager *jwtutil.Manager) *Service {
    return &Service{
        Feature: service.NewFeature(repo.Feature),
    }
}
```

### handler.go
- Initialize all handlers with service dependencies
- Add to Handler struct

Example:
```go
type Handler struct {
    FeatureHandler handler.FeatureHandler
}

func SetupHandler(service *Service) *Handler {
    return &Handler{
        FeatureHandler: handler.NewFeature(service.Feature),
    }
}
```

## RESPONSE UTILITIES (internal/utils/response/)
- `response.JSONOk(c, response.APIResponse{...})` - 200 status
- `response.JSONCreated(c, response.APIResponse{...})` - 201 status
- `response.FromError(c, err)` - error response

## CREATION ORDER FOR NEW FEATURE
1. Create repository in internal/src/repository/
2. Create service in internal/src/service/
3. Create handler in internal/src/handler/
4. Add repository to internal/container/repositories.go
5. Add service to internal/container/service.go
6. Add handler to internal/container/handler.go
7. Create route file in internal/routes/v1/
8. Register route in internal/routes/v1/v1.go
