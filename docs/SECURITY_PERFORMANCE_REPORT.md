# Security & Performance Analysis Report
**Project:** Capecom PM - Authentication & Session Management
**Date:** February 11, 2026
**Status:** ✅ Good Foundation with Critical Improvements Needed

---

## Executive Summary

Your authentication system has a **solid foundation** with JWT-based auth, refresh token rotation, and session management. However, there are **critical security vulnerabilities** and **performance bottlenecks** that need immediate attention.

**Overall Security Score:** 6.5/10
**Overall Performance Score:** 5/10

---

## 🔴 CRITICAL SECURITY ISSUES

### 1. **Missing JTI Validation in Middleware** ⚠️ HIGH PRIORITY
**Issue:** Your middleware validates JWT tokens but doesn't check if the JTI (JWT ID) matches an active session.

**Risk:** 
- Stolen access tokens remain valid even after logout
- No way to revoke individual sessions
- Token replay attacks possible

**Current Flow:**
```
User logs in → Access token issued with JTI
User logs out → Session revoked in DB
BUT: Access token still works until expiry!
```

**Fix Required:**
```go
// In middleware/user.go - Add JTI validation
claims, err := m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeUser)
if err != nil {
    response.FromError(c, err)
    c.Abort()
    return
}

// ⚠️ MISSING: Verify JTI exists in active session
session, err := m.SessionRepo.GetByJTI(claims.ID) // Need to add this method
if err != nil || session == nil || session.Status != models.SessionStatusActive {
    response.FromError(c, domainerrors.ErrUnauthorized)
    c.Abort()
    return
}
```

### 2. **No Device/IP Tracking on Session Creation** ⚠️ MEDIUM PRIORITY
**Issue:** Session table has device_id, device_name, user_agent, ip_address fields but they're never populated.

**Risk:**
- Cannot detect suspicious login locations
- No audit trail for security incidents
- Cannot implement "new device" alerts

**Fix Required:**
```go
// In services/auth.go - CreateAndReturnToken
session := &models.Session{
    UserID:           userID,
    JTI:              jti,
    RefreshTokenHash: hashedToken,
    RefreshExpiresAt: s.jwt.GetExpireTime(),
    // ADD THESE:
    DeviceID:   extractDeviceID(c),
    DeviceName: extractDeviceName(c),
    UserAgent:  &c.Request.UserAgent,
    IPAddress:  &c.ClientIP(),
}
```

### 3. **Weak Password Hashing Cost** ⚠️ MEDIUM PRIORITY
**Issue:** `bcrypt.GenerateFromPassword([]byte(p), 6)` uses cost factor 6 (very low).

**Risk:**
- Passwords can be brute-forced quickly
- OWASP recommends cost 12-14

**Fix:**
```go
// In utils/hash.go
func HashPassword(p string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(p), 12) // Change to 12
    return string(hash)
}
```

### 4. **No Rate Limiting on Login/Refresh** ⚠️ HIGH PRIORITY
**Issue:** No rate limiting on `/auth/login` or `/auth/refresh-token`.

**Risk:**
- Brute force attacks on passwords
- Token enumeration attacks
- DDoS vulnerability

**Fix:** Add rate limiting middleware (use `github.com/ulule/limiter/v3`)

### 5. **Session Cleanup Not Automated** ⚠️ LOW PRIORITY
**Issue:** `DeleteExpiredSessions()` exists but is never called.

**Risk:**
- Database bloat
- Potential performance degradation

**Fix:** Add cron job or background worker to call this periodically.

---

## 🟡 SECURITY IMPROVEMENTS NEEDED

### 6. **Missing Logout Endpoint**
You have login and refresh, but no logout to revoke sessions.

**Add:**
```go
// In handlers/auth_handler.go
func (h *AuthHandler) Logout(c *gin.Context) {
    userID := utils.GetUserID(c)
    jti := c.GetString("jti") // Extract from middleware
    
    err := h.AuthService.Logout(userID, jti)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    response.JSON(c, http.StatusOK, response.APIResponse{
        Success: true,
        Message: "Logged out successfully",
    })
}
```

### 7. **No "Logout All Devices" Feature**
Users should be able to revoke all sessions.

**Add:**
```go
func (s *AuthService) LogoutAllDevices(userUUID string) error {
    userID, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
    if err != nil || userID == nil {
        return domainerrors.ErrUserNotFound
    }
    return s.sessionRepo.RevokeAllUserSessions(*userID)
}
```

### 8. **Refresh Token Not Validated in Middleware**
The refresh endpoint doesn't use middleware to validate the refresh token structure.

### 9. **No HTTPS Enforcement**
Ensure all auth endpoints require HTTPS in production.

### 10. **Secrets in Environment Variables**
Good practice, but ensure `.env` is in `.gitignore` (it is ✅).

---

## 🔵 PERFORMANCE ISSUES

### 1. **N+1 Query Problem in Session Lookup** ⚠️ HIGH IMPACT
**Issue:** `GetByHashedToken` uses raw SQL with `Count()` after `Scan()`.

**Current:**
```go
err := r.DB.Raw("SELECT * FROM sessions WHERE refresh_token_hash = ? AND deleted_at IS NULL LIMIT 1", refreshToken).Scan(&session).Count(&count).Error
```

**Problem:** This executes TWO queries:
1. SELECT * FROM sessions...
2. SELECT COUNT(*) FROM sessions...

**Fix:**
```go
func (r *SessionRepo) GetByHashedToken(refreshToken string) (*models.Session, error) {
    var session models.Session
    err := r.DB.Where("refresh_token_hash = ? AND deleted_at IS NULL", refreshToken).
        First(&session).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    
    return &session, err
}
```

### 2. **Missing Index on JTI Column** ⚠️ HIGH IMPACT
**Issue:** Schema has index on `refresh_token_hash` but NOT on `jti`.

**Impact:** When you add JTI validation (critical fix #1), lookups will be slow.

**Fix in schema.sql:**
```sql
CREATE INDEX `idx_sessions_jti` ON `sessions` (`jti`);
CREATE INDEX `idx_sessions_jti_status` ON `sessions` (`jti`, `status`);
```

### 3. **Redis Cache Not Used for Session Validation** ⚠️ MEDIUM IMPACT
**Issue:** Every request hits the database to check user status, but sessions aren't cached.

**Current Flow:**
```
Request → Validate JWT → DB query for user status → Response
```

**Optimized Flow:**
```
Request → Validate JWT → Check Redis for session → (Cache miss: DB query) → Response
```

**Fix:**
```go
// Cache active sessions in Redis with TTL matching access token expiry
func (r *SessionRepo) GetByJTICached(jti string, cache *cacherepo.RedisRepo) (*models.Session, error) {
    cacheKey := fmt.Sprintf("session:jti:%s", jti)
    
    return cacherepo.GetOrSet(
        context.Background(),
        cache,
        cacheKey,
        5*time.Minute, // Match access token TTL
        func() (*models.Session, error) {
            return r.GetByJTI(jti)
        },
    )
}
```

### 4. **Unnecessary Transaction in Login** ⚠️ LOW IMPACT
**Issue:** `CreateAndReturnToken` wraps everything in a transaction, but only session creation needs atomicity.

**Impact:** Holds DB connection longer than needed.

**Fix:** Only wrap session creation/update in transaction, not Redis/JWT operations.

### 5. **User Status Check on Every Request** ⚠️ MEDIUM IMPACT
**Issue:** Middleware queries user status on every authenticated request.

**Current:**
```go
err = m.UserRepo.DB.Model(models.User{}).Where("uuid = ?", claims.UserID).Select("status").Scan(&status).Error
```

**Fix:** Cache user status in Redis with short TTL (already partially implemented in admin middleware).

### 6. **Missing Connection Pooling Configuration**
Ensure your database connection pool is configured:

```go
// In config/database.go
sqlDB, err := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

---

## 🟢 WHAT'S GOOD

### ✅ Strong Points

1. **Separate Access & Refresh Tokens** - Good security practice
2. **Refresh Token Rotation** - Prevents token reuse attacks
3. **SHA-256 Hashing for Refresh Tokens** - Correct approach
4. **Separate Admin/User JWT Secrets** - Good separation of concerns
5. **UUID for Public IDs** - Prevents enumeration attacks
6. **Soft Deletes** - Maintains audit trail
7. **Redis Caching Layer** - Good foundation (needs more usage)
8. **Environment-based Configuration** - Secure secret management
9. **GORM with Transactions** - Ensures data consistency
10. **Comprehensive Indexes** - Schema has good index coverage

---

## 📊 PERFORMANCE BENCHMARKS (Estimated)

### Current Performance
- **Login:** ~150-200ms (DB + Redis + JWT generation)
- **Refresh:** ~100-150ms (DB lookup + update + JWT generation)
- **Authenticated Request:** ~50-80ms (JWT validation + DB user status check)

### After Optimizations
- **Login:** ~120-150ms (-20% improvement)
- **Refresh:** ~60-80ms (-40% improvement with cached sessions)
- **Authenticated Request:** ~10-20ms (-75% improvement with cached validation)

---

## 🎯 PRIORITY ACTION ITEMS

### Immediate (This Week)
1. ✅ Add `GetByJTI` method to SessionRepo
2. ✅ Implement JTI validation in middleware
3. ✅ Add logout endpoint
4. ✅ Increase bcrypt cost to 12
5. ✅ Add rate limiting to auth endpoints

### Short Term (Next 2 Weeks)
6. ✅ Populate device/IP fields on session creation
7. ✅ Add JTI index to database
8. ✅ Implement session caching in Redis
9. ✅ Add "logout all devices" endpoint
10. ✅ Add automated session cleanup job

### Medium Term (Next Month)
11. ✅ Implement suspicious login detection
12. ✅ Add "new device" email notifications
13. ✅ Add session management UI (list active sessions)
14. ✅ Implement token blacklist for immediate revocation
15. ✅ Add comprehensive audit logging

---

## 📝 RECOMMENDED ARCHITECTURE CHANGES

### Session Validation Flow (Recommended)

```
┌─────────────┐
│   Request   │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  Extract JWT    │
└──────┬──────────┘
       │
       ▼
┌─────────────────┐
│  Validate JWT   │
│   Signature     │
└──────┬──────────┘
       │
       ▼
┌─────────────────┐
│  Check Redis    │
│  for JTI        │
└──────┬──────────┘
       │
   ┌───┴───┐
   │ Hit?  │
   └───┬───┘
       │
    No │  Yes
       │   │
       ▼   ▼
   ┌────┐ ┌────┐
   │ DB │ │ OK │
   └─┬──┘ └────┘
     │
     ▼
  ┌────────┐
  │ Cache  │
  │ Result │
  └────────┘
```

---

## 🔧 CODE FIXES TO IMPLEMENT

### Fix 1: Add GetByJTI Method

```go
// Add to internal/repositories/session.go
func (r *SessionRepo) GetByJTI(jti string) (*models.Session, error) {
    var session models.Session
    err := r.DB.Where("jti = ? AND deleted_at IS NULL", jti).First(&session).Error
    
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    
    return &session, err
}
```

### Fix 2: Update Middleware with JTI Validation

```go
// Update internal/middleware/user.go
func (m *UserMiddleware) VerifyUserToken() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.FromError(c, domainerrors.ErrUnauthorized)
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            response.FromError(c, domainerrors.ErrInvalidToken)
            c.Abort()
            return
        }

        tokenString := parts[1]
        claims, err := m.JWTManager.ValidateToken(tokenString, jwtutil.TokenTypeUser)
        if err != nil {
            response.FromError(c, err)
            c.Abort()
            return
        }

        // ✅ NEW: Validate JTI against active session
        session, err := m.SessionRepo.GetByJTI(claims.ID)
        if err != nil || session == nil || session.Status != models.SessionStatusActive {
            response.FromError(c, domainerrors.ErrUnauthorized)
            c.Abort()
            return
        }

        // ✅ NEW: Check session expiry
        if session.RefreshExpiresAt.Before(time.Now()) {
            response.FromError(c, domainerrors.ErrUnauthorized)
            c.Abort()
            return
        }

        // Verify user exists and is active
        var status string
        err = m.UserRepo.DB.Model(models.User{}).Where("uuid = ?", claims.UserID).Select("status").Scan(&status).Error
        if err != nil || status != "active" {
            response.FromError(c, domainerrors.ErrUnauthorized)
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Set("jti", claims.ID) // Store JTI for logout
        c.Next()
    }
}
```

### Fix 3: Add Logout Functionality

```go
// Add to internal/services/auth.go
func (s *AuthService) Logout(jti string) error {
    session, err := s.sessionRepo.GetByJTI(jti)
    if err != nil || session == nil {
        return domainerrors.ErrUnauthorized
    }
    
    return s.sessionRepo.UpdateStatus(session.UUID, models.SessionStatusRevoked)
}

func (s *AuthService) LogoutAllDevices(userUUID string) error {
    userID, err := s.redisRepo.GetUserIdByUuid(userUUID, *s.userRepo)
    if err != nil || userID == nil {
        return domainerrors.ErrUserNotFound
    }
    
    return s.sessionRepo.RevokeAllUserSessions(uint64(*userID))
}
```

```go
// Add to internal/handlers/auth_handler.go
func (h *AuthHandler) Logout(c *gin.Context) {
    jti := c.GetString("jti")
    
    err := h.AuthService.Logout(jti)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    response.JSON(c, http.StatusOK, response.APIResponse{
        Success: true,
        Message: "Logged out successfully",
    })
}

func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
    userID := utils.GetUserID(c)
    
    err := h.AuthService.LogoutAllDevices(userID)
    if err != nil {
        response.FromError(c, err)
        return
    }
    
    response.JSON(c, http.StatusOK, response.APIResponse{
        Success: true,
        Message: "All sessions revoked successfully",
    })
}
```

```go
// Update internal/routes/version/v1/auth.go
func AuthRoutes(v1 *gin.RouterGroup, c *container.Container) {
    h := c.Handler.AuthHandler
    auth := v1.Group("/auth")

    auth.POST("/login", h.Login)
    auth.POST("/refresh-token", h.Refresh)
    auth.GET("/me", c.Middleware.UserMiddleware.VerifyUserToken(), h.Me)
    auth.POST("/logout", c.Middleware.UserMiddleware.VerifyUserToken(), h.Logout) // NEW
    auth.POST("/logout-all", c.Middleware.UserMiddleware.VerifyUserToken(), h.LogoutAllDevices) // NEW
}
```

### Fix 4: Add Device/IP Tracking

```go
// Add helper functions to internal/utils/utils.go
func GetClientIP(c *gin.Context) string {
    // Check X-Forwarded-For header first (for proxies/load balancers)
    forwarded := c.GetHeader("X-Forwarded-For")
    if forwarded != "" {
        // Take the first IP if multiple
        ips := strings.Split(forwarded, ",")
        return strings.TrimSpace(ips[0])
    }
    
    // Check X-Real-IP header
    realIP := c.GetHeader("X-Real-IP")
    if realIP != "" {
        return realIP
    }
    
    // Fall back to RemoteAddr
    return c.ClientIP()
}

func GetDeviceInfo(userAgent string) (deviceID, deviceName string) {
    // Simple device fingerprinting (enhance with proper library)
    hash := sha256.Sum256([]byte(userAgent))
    deviceID = hex.EncodeToString(hash[:16])
    
    // Parse user agent for device name
    ua := strings.ToLower(userAgent)
    switch {
    case strings.Contains(ua, "mobile"):
        deviceName = "Mobile Device"
    case strings.Contains(ua, "tablet"):
        deviceName = "Tablet"
    case strings.Contains(ua, "windows"):
        deviceName = "Windows PC"
    case strings.Contains(ua, "mac"):
        deviceName = "Mac"
    case strings.Contains(ua, "linux"):
        deviceName = "Linux PC"
    default:
        deviceName = "Unknown Device"
    }
    
    return deviceID, deviceName
}
```

```go
// Update internal/services/auth.go - Login method
func (s AuthService) Login(req authdto.LoginRequest, c *gin.Context) (*authdto.LoginResponse, error) {
    if usr, err := s.userRepo.FindByEmail(req.Email); err != nil {
        return nil, err
    } else if usr == nil {
        return nil, domainerrors.ErrInvalidCredentials
    } else {
        if !utils.CheckPassword(usr.PasswordHash, req.Password) {
            return nil, domainerrors.ErrInvalidCredentials
        }
        return s.CreateAndReturnToken(usr.UUID, "", usr.IsAdmin, c)
    }
}

// Update CreateAndReturnToken signature
func (s AuthService) CreateAndReturnToken(userUuid, oldToken string, isAdmin bool, c *gin.Context) (*authdto.LoginResponse, error) {
    jti := uuid.NewString()
    refreshToken, err := utils.GenerateRefreshToken()
    if err != nil {
        return nil, err
    }

    hashedToken := utils.HashToken(refreshToken)
    var accessToken string

    err = s.sessionRepo.DB.Transaction(func(tx *gorm.DB) error {
        if oldToken == "" {
            var userID int64 = 0
            userIdIn, err := s.redisRepo.GetUserIdByUuid(userUuid, *s.userRepo)
            if err != nil || userIdIn == nil {
                return domainerrors.ErrUnauthorized
            }
            if userIdIn != nil {
                userID = *userIdIn
            }

            if err != nil {
                return err
            }
            if userID == int64(0) {
                return domainerrors.ErrUnauthorized
            }

            // ✅ NEW: Extract device and IP info
            userAgent := c.Request.UserAgent()
            ipAddress := utils.GetClientIP(c)
            deviceID, deviceName := utils.GetDeviceInfo(userAgent)

            session := &models.Session{
                UserID:           uint64(userID),
                JTI:              jti,
                RefreshTokenHash: hashedToken,
                RefreshExpiresAt: s.jwt.GetExpireTime(),
                DeviceID:         &deviceID,
                DeviceName:       &deviceName,
                UserAgent:        &userAgent,
                IPAddress:        &ipAddress,
            }

            if err := s.sessionRepo.Create(session); err != nil {
                return err
            }
        } else {
            // ... existing refresh token logic
        }
        
        // ... rest of the function
    })

    // ... rest of the function
}
```

### Fix 5: Add Database Indexes

```sql
-- Add to schema.sql after sessions table creation
CREATE INDEX `idx_sessions_jti` ON `sessions` (`jti`);
CREATE INDEX `idx_sessions_jti_status` ON `sessions` (`jti`, `status`);
CREATE INDEX `idx_sessions_user_status` ON `sessions` (`user_id`, `status`, `deleted_at`);
```

### Fix 6: Increase Bcrypt Cost

```go
// Update internal/utils/hash.go
func HashPassword(p string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(p), 12) // Changed from 6 to 12
    return string(hash)
}
```

### Fix 7: Add Rate Limiting

```go
// Create internal/middleware/rate_limit.go
package middleware

import (
    "capecom-pm/internal/utils/response"
    domainerrors "capecom-pm/internal/domain/error"
    "net/http"
    "sync"
    "time"
    
    "github.com/gin-gonic/gin"
)

type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.RWMutex
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
    
    // Cleanup old entries every minute
    go rl.cleanup()
    
    return rl
}

func (rl *RateLimiter) cleanup() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        rl.mu.Lock()
        now := time.Now()
        for key, times := range rl.requests {
            // Remove entries older than window
            valid := []time.Time{}
            for _, t := range times {
                if now.Sub(t) < rl.window {
                    valid = append(valid, t)
                }
            }
            if len(valid) == 0 {
                delete(rl.requests, key)
            } else {
                rl.requests[key] = valid
            }
        }
        rl.mu.Unlock()
    }
}

func (rl *RateLimiter) Limit() gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.ClientIP() // Or use email from request body for login
        
        rl.mu.Lock()
        defer rl.mu.Unlock()
        
        now := time.Now()
        
        // Get request times for this key
        times, exists := rl.requests[key]
        if !exists {
            times = []time.Time{}
        }
        
        // Remove old requests outside the window
        valid := []time.Time{}
        for _, t := range times {
            if now.Sub(t) < rl.window {
                valid = append(valid, t)
            }
        }
        
        // Check if limit exceeded
        if len(valid) >= rl.limit {
            response.FromError(c, domainerrors.NewWithCode(
                http.StatusTooManyRequests,
                "Too many requests. Please try again later.",
                "rate_limiter",
                "Limit",
            ))
            c.Abort()
            return
        }
        
        // Add current request
        valid = append(valid, now)
        rl.requests[key] = valid
        
        c.Next()
    }
}
```

```go
// Update internal/routes/version/v1/auth.go
func AuthRoutes(v1 *gin.RouterGroup, c *container.Container) {
    h := c.Handler.AuthHandler
    auth := v1.Group("/auth")
    
    // Rate limiting: 5 requests per minute for auth endpoints
    rateLimiter := middleware.NewRateLimiter(5, 1*time.Minute)

    auth.POST("/login", rateLimiter.Limit(), h.Login)
    auth.POST("/refresh-token", rateLimiter.Limit(), h.Refresh)
    auth.GET("/me", c.Middleware.UserMiddleware.VerifyUserToken(), h.Me)
    auth.POST("/logout", c.Middleware.UserMiddleware.VerifyUserToken(), h.Logout)
    auth.POST("/logout-all", c.Middleware.UserMiddleware.VerifyUserToken(), h.LogoutAllDevices)
}
```

### Fix 8: Add Session Cleanup Job

```go
// Create internal/jobs/session_cleanup.go
package jobs

import (
    "capecom-pm/internal/repositories"
    "log"
    "time"
)

type SessionCleanupJob struct {
    sessionRepo *repositories.SessionRepo
}

func NewSessionCleanupJob(sessionRepo *repositories.SessionRepo) *SessionCleanupJob {
    return &SessionCleanupJob{
        sessionRepo: sessionRepo,
    }
}

func (j *SessionCleanupJob) Start() {
    ticker := time.NewTicker(24 * time.Hour) // Run daily
    
    // Run immediately on start
    j.run()
    
    // Then run on schedule
    go func() {
        for range ticker.C {
            j.run()
        }
    }()
}

func (j *SessionCleanupJob) run() {
    log.Println("Running session cleanup job...")
    
    err := j.sessionRepo.DeleteExpiredSessions()
    if err != nil {
        log.Printf("Session cleanup failed: %v", err)
    } else {
        log.Println("Session cleanup completed successfully")
    }
}
```

```go
// Update cmd/main.go to start the job
func main() {
    // ... existing setup ...
    
    // Start background jobs
    sessionCleanupJob := jobs.NewSessionCleanupJob(container.Repository.SessionRepo)
    sessionCleanupJob.Start()
    
    // ... rest of main ...
}
```

---

## 📈 MONITORING & METRICS TO ADD

### Key Metrics to Track

1. **Authentication Metrics**
   - Login success/failure rate
   - Average login time
   - Failed login attempts per IP
   - Token refresh rate

2. **Session Metrics**
   - Active sessions per user
   - Average session duration
   - Session revocation rate
   - Expired sessions cleaned up

3. **Performance Metrics**
   - Database query time (P50, P95, P99)
   - Redis cache hit rate
   - Middleware execution time
   - API response times

4. **Security Metrics**
   - Rate limit violations
   - Invalid token attempts
   - Suspicious login patterns
   - Concurrent sessions per user

---

## 🧪 TESTING RECOMMENDATIONS

### Unit Tests Needed

```go
// Test session creation
func TestSessionRepo_Create(t *testing.T) {
    // Test UUID generation
    // Test default status
    // Test timestamp population
}

// Test JTI validation
func TestSessionRepo_GetByJTI(t *testing.T) {
    // Test valid JTI
    // Test invalid JTI
    // Test deleted session
}

// Test token rotation
func TestAuthService_RefreshToken(t *testing.T) {
    // Test valid refresh
    // Test expired token
    // Test revoked session
    // Test token reuse
}

// Test rate limiting
func TestRateLimiter_Limit(t *testing.T) {
    // Test within limit
    // Test exceeding limit
    // Test window expiry
}
```

### Integration Tests Needed

1. Full login flow with session creation
2. Token refresh with session rotation
3. Logout with session revocation
4. Concurrent login attempts
5. Session expiry handling

---

## 📚 ADDITIONAL RECOMMENDATIONS

### 1. Add Session Management Endpoints

```go
// List active sessions
GET /auth/sessions

// Revoke specific session
DELETE /auth/sessions/:sessionId
```

### 2. Implement Token Blacklist for Immediate Revocation

Use Redis SET with TTL matching token expiry:

```go
func (r *RedisRepo) BlacklistToken(jti string, ttl time.Duration) error {
    return r.SetString(context.Background(), fmt.Sprintf("blacklist:%s", jti), "1", ttl)
}

func (r *RedisRepo) IsTokenBlacklisted(jti string) bool {
    _, err := r.GetString(context.Background(), fmt.Sprintf("blacklist:%s", jti))
    return err == nil
}
```

### 3. Add Audit Logging

Log all authentication events:
- Login attempts (success/failure)
- Token refreshes
- Logouts
- Session revocations
- Suspicious activities

### 4. Implement CSRF Protection

For web clients, add CSRF tokens to prevent cross-site request forgery.

### 5. Add 2FA Support

Consider adding two-factor authentication for enhanced security.

---

## 🎓 SECURITY BEST PRACTICES CHECKLIST

- [x] Passwords hashed with bcrypt
- [ ] Bcrypt cost factor >= 12 (currently 6)
- [x] Refresh tokens hashed before storage
- [x] Separate secrets for different token types
- [ ] JTI validation in middleware
- [ ] Rate limiting on auth endpoints
- [ ] Device/IP tracking
- [ ] Logout functionality
- [ ] Session expiry enforcement
- [ ] HTTPS enforcement (production)
- [x] Environment-based secrets
- [ ] Audit logging
- [ ] Suspicious activity detection
- [ ] Token blacklist for immediate revocation

---

## 💰 ESTIMATED IMPLEMENTATION TIME

| Task | Priority | Effort | Impact |
|------|----------|--------|--------|
| Add JTI validation | Critical | 2 hours | High |
| Add logout endpoints | Critical | 1 hour | High |
| Increase bcrypt cost | Critical | 15 mins | High |
| Add rate limiting | High | 3 hours | High |
| Add device/IP tracking | High | 2 hours | Medium |
| Add JTI indexes | High | 15 mins | High |
| Implement session caching | Medium | 4 hours | High |
| Add session cleanup job | Medium | 1 hour | Low |
| Add session management UI | Low | 8 hours | Medium |
| Implement audit logging | Low | 6 hours | Medium |

**Total Critical Path:** ~8 hours
**Total Recommended:** ~27 hours

---

## 🏁 CONCLUSION

Your authentication system has a **solid architectural foundation** but needs **critical security hardening** before production deployment. The main issues are:

1. **Missing JTI validation** - Tokens can't be revoked effectively
2. **No rate limiting** - Vulnerable to brute force
3. **Weak password hashing** - Easy to crack
4. **Missing logout** - No way to end sessions

The good news: These are all fixable within a week of focused work. The performance optimizations (caching, indexing) will provide significant speed improvements with minimal code changes.

**Recommendation:** Implement the critical fixes immediately, then roll out the high-priority improvements before production launch.

---

**Report Generated:** February 11, 2026
**Next Review:** After implementing critical fixes
