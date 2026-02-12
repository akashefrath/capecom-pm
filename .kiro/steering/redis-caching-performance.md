---
inclusion: auto
description: Redis caching, GetOrSet pattern, rate limiting, cache keys, performance goals for low-latency APIs
---

# Redis Caching & Performance

This project uses Redis for caching, rate limiting, and session management. The goal is high-efficiency APIs with low latency.

---

## Redis Architecture

```
internal/cache/redis.go          — Redis client connection
internal/cache/rate_cache.go     — Sliding window rate limiter middleware
internal/repositories/cache/redis.go — RedisRepo with GetOrSet generic pattern
```

Redis client is initialized in `cmd/main.go` and passed through the DI container. The `cache.IsRedisConnected` flag allows graceful degradation — if Redis is down, the app continues without caching.

---

## GetOrSet Pattern (Cache-Aside)

The core caching pattern is the generic `GetOrSet[T]` function in `internal/repositories/cache/redis.go`. It implements cache-aside: try cache first, fall back to DB, then populate cache.

```go
func GetOrSet[T any](
    ctx context.Context,
    cache *RedisRepo,
    key string,
    ttl time.Duration,       // 0 = default 15 minutes
    dbFunc func() (*T, error),
) (*T, error)
```

**Usage:**
```go
// Cache a user's status for 15 min (ttl=0 uses default)
status, err := cacherepo.GetOrSet(
    context.Background(),
    m.redisRepo,
    fmt.Sprintf("user_status_%s", claims.UserID),
    0,
    func() (*string, error) {
        return m.UserRepo.FindUserStatus(claims.UserID)
    },
)

// Cache a session for 5 minutes
session, err := cacherepo.GetOrSet(
    context.Background(),
    m.redisRepo,
    fmt.Sprintf("session:jti:%s", claims.ID),
    5*time.Minute,
    func() (*models.Session, error) {
        return m.SessionRepo.GetByJTI(claims.ID)
    },
)
```

**How it works:**
1. Try to get value from Redis by key
2. If found, JSON unmarshal and return
3. If not found, call `dbFunc()` to fetch from database
4. JSON marshal the result and store in Redis with TTL
5. Return the value

---

## Current Cache Keys

| Key Pattern | TTL | Used In | Purpose |
|---|---|---|---|
| `user_status_{uuid}` | 15 min (default) | Auth/Admin/User middleware | Cache user active/inactive status |
| `session:jti:{jti}` | 5 min | Auth/Admin/User middleware | Cache session lookup by JWT ID |
| `uuid_by_id:{user_id}` | 15 min (default) | RedisRepo.GetUserUuidById | Map integer user ID → UUID |
| `id_by_uuid:{uuid}` | 15 min (default) | RedisRepo.GetUserIdByUuid | Map UUID → integer user ID |
| `rate_limit:{ip}` | window duration | Rate limiter middleware | Sliding window rate limit per IP |

---

## Rate Limiting

The project uses a Redis-based sliding window rate limiter in `internal/cache/rate_cache.go`.

**Global rate limit** (applied to all `/api/v1` routes):
```go
// internal/routes/routes.go
apiV1.Use(cache.SlidingWindowRateLimiter(c.RedisClient, 100, time.Minute))
```

**Per-route rate limit** (e.g., auth routes have stricter limits):
```go
// internal/routes/version/v1/auth.go
auth.Use(cache.SlidingWindowRateLimiter(c.RedisClient, 30, time.Minute))
```

The rate limiter uses Redis sorted sets with timestamps as scores. It fails open — if Redis is down, requests pass through.

---

## How to Use Redis Cache in New Features

### In Middleware

Middleware already has `redisRepo` injected. Use `GetOrSet` directly:

```go
func (m *YourMiddleware) SomeCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        result, err := cacherepo.GetOrSet(
            context.Background(),
            m.redisRepo,
            "your_cache_key",
            5*time.Minute,
            func() (*YourType, error) {
                return m.SomeRepo.FetchFromDB()
            },
        )
        // use result...
    }
}
```

### In Services

To use Redis in a service, inject `RedisRepo` through the container:

1. Add `redisRepo` to the service struct:
```go
type YourService struct {
    yourRepo  *repositories.YourRepo
    redisRepo *cacherepo.RedisRepo
}
```

2. Pass it from `internal/container/service.go`:
```go
YourService: services.NewYourService(repository.YourRepo, repository.CacheRepo),
```

3. Use `GetOrSet` in service methods:
```go
func (s *YourService) GetSomething(id string) (*YourDTO, error) {
    result, err := cacherepo.GetOrSet(
        context.Background(),
        s.redisRepo,
        fmt.Sprintf("something:%s", id),
        10*time.Minute,
        func() (*YourDTO, error) {
            return s.yourRepo.FindByID(id)
        },
    )
    return result, err
}
```

### Cache Invalidation

When data changes, delete the cache key:

```go
cacheKey := fmt.Sprintf("session:jti:%s", jti)
_ = s.redisRepo.Delete(context.Background(), cacheKey)
```

See `internal/services/auth.go` → `LogoutUserByJTI()` for a real example.

---

## RedisRepo Methods

Available on `*cacherepo.RedisRepo`:

```go
SetString(ctx, key, value, ttl)   // Store a string value
GetString(ctx, key)                // Get a string value
Delete(ctx, key)                   // Delete a key
Exists(ctx, key)                   // Check if key exists
GetUserUuidById(userID, userRepo)  // Cached user ID → UUID lookup
GetUserIdByUuid(uuid, userRepo)    // Cached UUID → user ID lookup
```

Plus the standalone generic function:
```go
cacherepo.GetOrSet[T](ctx, redisRepo, key, ttl, dbFunc)
```

---

## Performance Guidelines

### Cache TTL Strategy

- **Hot path data** (middleware checks): 5 min — user status, session validation
- **Lookup mappings** (ID ↔ UUID): 15 min — rarely changes
- **Business data**: 5–10 min depending on freshness requirements
- **Rate limit windows**: Match the rate limit window duration

### Low-Latency API Checklist

- Use `GetOrSet` for any data fetched on every request (user status, sessions, permissions)
- Cache role/permission checks that hit the DB on every request
- Use Redis pipeline for multiple cache operations in a single round-trip
- Keep cache keys short and predictable
- Always set a TTL — never cache indefinitely
- Invalidate cache on write operations (create, update, delete)
- Use `context.Background()` for cache operations (don't tie to request context for writes)
- DB connection pooling: uncomment and configure `SetMaxIdleConns`, `SetMaxOpenConns`, `SetConnMaxLifetime` in `internal/config/database.go`
- Use raw SQL queries with joins instead of multiple GORM queries where possible (see `UserRepo.UserJoinRaw` pattern)
- Use `UNION ALL` for fetching multiple related IDs in one query (see `MasterDataRepo.GetUserRelatedIDs`)

### What to Cache

- User status checks (already cached)
- Session lookups (already cached)
- User ID ↔ UUID mappings (already cached)
- Role/permission checks (consider caching `IsManagerOrAdmin` results)
- Master data (roles, groups, departments, designations) — rarely changes
- Frequently accessed read-heavy data

### What NOT to Cache

- Write operations
- Data that changes on every request
- Large result sets (paginated lists)
- Sensitive data that shouldn't persist (raw tokens, passwords)
