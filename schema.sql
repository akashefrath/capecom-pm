-- =========================================================
-- USERS
-- =========================================================
CREATE TABLE users (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       uuid CHAR(36) NOT NULL,
                       name VARCHAR(255) NOT NULL,
                       email VARCHAR(255) NOT NULL,
                       phone VARCHAR(50),
                       country_code INT,
                       password_hash VARCHAR(255) NOT NULL,
                       employee_id VARCHAR(255),
                       group_id BIGINT,
                       designation_id BIGINT,
                       department_id BIGINT,
                       is_admin BOOLEAN DEFAULT FALSE,
                       created_by BIGINT,
                       status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP NULL
);

-- Soft delete friendly unique constraints
CREATE UNIQUE INDEX idx_users_email_active ON users (email, deleted_at);
CREATE UNIQUE INDEX idx_users_uuid_active ON users (uuid, deleted_at);
CREATE UNIQUE INDEX idx_users_employee_active ON users (employee_id, deleted_at);

-- Fast login lookup
CREATE INDEX idx_users_login_status ON users (email, status, deleted_at);



-- =========================================================
-- SESSIONS (HOT TABLE - HIGHLY OPTIMIZED)
-- =========================================================
CREATE TABLE sessions (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT,
                          uuid CHAR(36) NOT NULL,
                          user_id BIGINT NOT NULL,
                          jti CHAR(36) NOT NULL,
                          refresh_token_hash BINARY(32) NOT NULL,
                          refresh_expires_at TIMESTAMP NOT NULL,
                          rotated_at TIMESTAMP NULL,
                          status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                          device_id VARCHAR(100),
                          device_name VARCHAR(100),
                          user_agent TEXT,
                          ip_address VARCHAR(45),
                          last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          deleted_at TIMESTAMP NULL
);

-- Main JWT validation lookup
CREATE UNIQUE INDEX idx_sessions_jti_active
    ON sessions (jti, status);

-- Middleware lookup
CREATE INDEX idx_sessions_middleware
    ON sessions (jti, deleted_at, status);

-- Refresh token rotation lookup
CREATE INDEX idx_sessions_refresh
    ON sessions (refresh_token_hash, status);

-- Logout all devices
CREATE INDEX idx_sessions_user_active
    ON sessions (user_id, status, deleted_at);

-- Cleanup expired tokens
CREATE INDEX idx_sessions_expiry_cleanup
    ON sessions (refresh_expires_at, status);



-- =========================================================
-- ROLES
-- =========================================================
CREATE TABLE roles (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       uuid CHAR(36) NOT NULL,
                       name VARCHAR(50) NOT NULL,
                       status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX idx_roles_uuid_active ON roles (uuid, deleted_at);
CREATE UNIQUE INDEX idx_roles_name_active ON roles (name, deleted_at);
CREATE INDEX idx_roles_status ON roles (status, deleted_at);



-- =========================================================
-- USER ROLES (RBAC LINK TABLE)
-- =========================================================
CREATE TABLE user_roles (
                            id BIGINT PRIMARY KEY AUTO_INCREMENT,
                            uuid CHAR(36) NOT NULL,
                            user_id BIGINT NOT NULL,
                            role_id BIGINT NOT NULL,
                            created_by BIGINT,
                            status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
                            deleted_at TIMESTAMP NULL
);

-- RBAC permission check
CREATE INDEX idx_user_roles_lookup
    ON user_roles (user_id, role_id, status, deleted_at);

-- Prevent duplicate active role
CREATE UNIQUE INDEX idx_user_role_unique
    ON user_roles (user_id, role_id, deleted_at);







-- =========================================================
-- FOREIGN KEYS (ONLY WHERE SAFE)
-- NOTE: NO FK ON SESSIONS (PERFORMANCE REASONS)
-- =========================================================
ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id) REFERENCES roles(id)
            ON DELETE CASCADE;





DELIMITER $$

    CREATE TRIGGER `users_updated_at` BEFORE UPDATE ON `users`
        FOR EACH ROW
    BEGIN
        SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

    CREATE TRIGGER `sessions_updated_at` BEFORE UPDATE ON `sessions`
        FOR EACH ROW
    BEGIN
        SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$


        CREATE TRIGGER `roles_updated_at` BEFORE UPDATE ON `roles`
            FOR EACH ROW
        BEGIN
            SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$


            CREATE TRIGGER `user_roles_updated_at` BEFORE UPDATE ON `user_roles`
                FOR EACH ROW
            BEGIN
                SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

DELIMITER ;


