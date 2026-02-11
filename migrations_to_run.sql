-- ============================================
-- CRITICAL DATABASE MIGRATIONS
-- Run these immediately for security fixes
-- ============================================

-- 1. Add missing indexes on sessions table for JTI lookups
-- This is CRITICAL for performance when JTI validation is added
CREATE INDEX IF NOT EXISTS `idx_sessions_jti` ON `sessions` (`jti`);
CREATE INDEX IF NOT EXISTS `idx_sessions_jti_status` ON `sessions` (`jti`, `status`);
CREATE INDEX IF NOT EXISTS `idx_sessions_user_status` ON `sessions` (`user_id`, `status`, `deleted_at`);

-- 2. Verify sessions table structure
-- Run this to check if all columns exist
DESCRIBE sessions;

-- Expected output should include:
-- - id, uuid, user_id, jti, refresh_token_hash
-- - refresh_expires_at, rotated_at, status
-- - device_id, device_name, user_agent, ip_address
-- - last_used_at, created_at, updated_at, deleted_at

-- 3. Check existing indexes
SHOW INDEX FROM sessions;

-- 4. Verify foreign key constraint exists
SELECT 
    CONSTRAINT_NAME,
    TABLE_NAME,
    COLUMN_NAME,
    REFERENCED_TABLE_NAME,
    REFERENCED_COLUMN_NAME
FROM
    INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE
    TABLE_NAME = 'sessions'
    AND REFERENCED_TABLE_NAME IS NOT NULL;

-- ============================================
-- OPTIONAL: Performance Optimization Indexes
-- Add these for better query performance
-- ============================================

-- Index for finding sessions by device
CREATE INDEX IF NOT EXISTS `idx_sessions_device` ON `sessions` (`device_id`, `user_id`);

-- Index for finding sessions by IP (security monitoring)
CREATE INDEX IF NOT EXISTS `idx_sessions_ip` ON `sessions` (`ip_address`, `created_at`);

-- Index for cleanup job (finding expired sessions)
CREATE INDEX IF NOT EXISTS `idx_sessions_expired` ON `sessions` (`refresh_expires_at`, `deleted_at`);

-- ============================================
-- VERIFICATION QUERIES
-- Run these to verify everything is working
-- ============================================

-- Check if sessions table is empty or has data
SELECT COUNT(*) as total_sessions FROM sessions;

-- Check active sessions
SELECT COUNT(*) as active_sessions 
FROM sessions 
WHERE status = 'active' AND deleted_at IS NULL;

-- Check for expired sessions that need cleanup
SELECT COUNT(*) as expired_sessions 
FROM sessions 
WHERE refresh_expires_at < NOW() AND deleted_at IS NULL;

-- View recent sessions (last 10)
SELECT 
    uuid,
    user_id,
    status,
    device_name,
    ip_address,
    last_used_at,
    created_at
FROM sessions
ORDER BY created_at DESC
LIMIT 10;

-- ============================================
-- CLEANUP QUERIES (Run periodically)
-- ============================================

-- Soft delete expired sessions (run daily via cron)
-- UPDATE sessions 
-- SET deleted_at = NOW() 
-- WHERE refresh_expires_at < NOW() 
-- AND deleted_at IS NULL;

-- Hard delete old soft-deleted sessions (run monthly)
-- DELETE FROM sessions 
-- WHERE deleted_at IS NOT NULL 
-- AND deleted_at < DATE_SUB(NOW(), INTERVAL 90 DAY);

-- ============================================
-- MONITORING QUERIES
-- ============================================

-- Sessions per user (find users with many sessions)
SELECT 
    user_id,
    COUNT(*) as session_count
FROM sessions
WHERE status = 'active' AND deleted_at IS NULL
GROUP BY user_id
HAVING session_count > 5
ORDER BY session_count DESC;

-- Sessions by status
SELECT 
    status,
    COUNT(*) as count
FROM sessions
WHERE deleted_at IS NULL
GROUP BY status;

-- Recent login activity (last 24 hours)
SELECT 
    DATE_FORMAT(created_at, '%Y-%m-%d %H:00:00') as hour,
    COUNT(*) as login_count
FROM sessions
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
GROUP BY hour
ORDER BY hour DESC;

-- Suspicious activity: Multiple IPs for same user
SELECT 
    user_id,
    COUNT(DISTINCT ip_address) as ip_count,
    GROUP_CONCAT(DISTINCT ip_address) as ip_addresses
FROM sessions
WHERE status = 'active' 
AND deleted_at IS NULL
AND created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
GROUP BY user_id
HAVING ip_count > 3;

-- ============================================
-- ROLLBACK (if needed)
-- ============================================

-- Drop indexes if you need to rollback
-- DROP INDEX `idx_sessions_jti` ON `sessions`;
-- DROP INDEX `idx_sessions_jti_status` ON `sessions`;
-- DROP INDEX `idx_sessions_user_status` ON `sessions`;
-- DROP INDEX `idx_sessions_device` ON `sessions`;
-- DROP INDEX `idx_sessions_ip` ON `sessions`;
-- DROP INDEX `idx_sessions_expired` ON `sessions`;
