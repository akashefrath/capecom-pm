

-- =========================================================
-- 1. ATTENDANCE POLICIES
-- =========================================================
CREATE TABLE attendance_policies (
                                     id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                     uuid CHAR(36) UNIQUE NOT NULL,
                                     name VARCHAR(50) UNIQUE NOT NULL,
                                     min_work_hours_minutes INT NOT NULL,
                                     half_day_minutes INT NOT NULL,
                                     late_grace_minutes INT NOT NULL,
                                     early_exit_grace_minutes INT NOT NULL,
                                     max_break_minutes INT NOT NULL,
                                     auto_checkout_time INT NOT NULL,
                                     is_default BOOLEAN DEFAULT FALSE,
                                     created_by BIGINT,
                                     status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
                                     deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_policy_default ON attendance_policies (is_default);

-- =========================================================
-- 2. ATTENDANCE POLICY GROUPS
-- =========================================================
CREATE TABLE attendance_policy_groups (
                                          id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                          uuid CHAR(36) UNIQUE NOT NULL,
                                          name VARCHAR(50) UNIQUE NOT NULL,
                                          attendance_policy_id BIGINT NOT NULL,
                                          created_by BIGINT,
                                          status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                          updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
                                          deleted_at TIMESTAMP NULL,
                                          CONSTRAINT fk_policy_group_policy FOREIGN KEY (attendance_policy_id) REFERENCES attendance_policies(id) ON DELETE CASCADE
);

CREATE INDEX idx_policy_groups_policy_id ON attendance_policy_groups (attendance_policy_id);

-- =========================================================
-- 3. SHIFT SYSTEM
-- =========================================================
CREATE TABLE shift_system (
                              id BIGINT PRIMARY KEY AUTO_INCREMENT,
                              uuid CHAR(36) UNIQUE NOT NULL,
                              name VARCHAR(50) UNIQUE NOT NULL,
                              start_time TIME NOT NULL,
                              end_time TIME NOT NULL,
                              checkin_early INT NOT NULL DEFAULT 0,
                              checkin_late INT NOT NULL DEFAULT 0,
                              checkout_early INT NOT NULL DEFAULT 0,
                              checkout_late INT NOT NULL DEFAULT 0,
                              is_overnight BOOLEAN NOT NULL DEFAULT FALSE,
                              is_default BOOLEAN DEFAULT FALSE,
                              created_by BIGINT,
                              status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
                              deleted_at TIMESTAMP NULL
);

-- =========================================================
-- 4. SHIFT SYSTEM GROUPS
-- =========================================================
CREATE TABLE shift_system_groups (
                                     id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                     uuid CHAR(36) UNIQUE NOT NULL,
                                     name VARCHAR(50) UNIQUE NOT NULL,
                                     shift_system_id BIGINT NOT NULL,
                                     created_by BIGINT,
                                     status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
                                     deleted_at TIMESTAMP NULL,
                                     CONSTRAINT fk_shift_groups_system FOREIGN KEY (shift_system_id) REFERENCES shift_system(id) ON DELETE CASCADE
);

CREATE INDEX idx_shift_groups_system_id ON shift_system_groups (shift_system_id);




-- =========================================================
-- 5. USERS
-- =========================================================
CREATE TABLE users (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       uuid CHAR(36) UNIQUE NOT NULL,
                       name VARCHAR(255) UNIQUE NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       phone VARCHAR(50),
                       country_code INT,
                       password_hash VARCHAR(255) NOT NULL,
                       employee_id VARCHAR(255) UNIQUE,
                       group_id BIGINT,
                       designation_id BIGINT,
                       department_id BIGINT,
                       shift_system_group_id BIGINT,
                       attendance_policy_group_id BIGINT,
                       is_admin BOOLEAN DEFAULT FALSE,
                       created_by BIGINT,
                       status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP NULL,
                       CONSTRAINT fk_users_shift_system_group FOREIGN KEY (shift_system_group_id) REFERENCES shift_system_groups(id) ON DELETE CASCADE,
                       CONSTRAINT fk_users_attendance_policy_group FOREIGN KEY (attendance_policy_group_id) REFERENCES attendance_policy_groups(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_users_email_active ON users (email, deleted_at);
CREATE UNIQUE INDEX idx_users_uuid_active ON users (uuid, deleted_at);
CREATE UNIQUE INDEX idx_users_employee_active ON users (employee_id, deleted_at);
CREATE INDEX idx_users_login_status ON users (email, status, deleted_at);

-- =========================================================
-- 6. SESSIONS
-- =========================================================
CREATE TABLE sessions (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT,
                          uuid CHAR(36) UNIQUE NOT NULL,
                          user_id BIGINT NOT NULL,
                          jti CHAR(36) UNIQUE NOT NULL,
                          refresh_token_hash BINARY(32) UNIQUE NOT NULL,
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
                          deleted_at TIMESTAMP NULL,
                          CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_sessions_jti_active ON sessions (jti, status);
CREATE INDEX idx_sessions_middleware ON sessions (jti, deleted_at, status);
CREATE INDEX idx_sessions_refresh ON sessions (refresh_token_hash, status);
CREATE INDEX idx_sessions_user_active ON sessions (user_id, status, deleted_at);
CREATE INDEX idx_sessions_expiry_cleanup ON sessions (refresh_expires_at, status);

-- =========================================================
-- 7. ROLES
-- =========================================================
CREATE TABLE roles (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       uuid CHAR(36) UNIQUE NOT NULL,
                       name VARCHAR(50) UNIQUE NOT NULL,
                       status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP NULL
);

CREATE UNIQUE INDEX idx_roles_uuid_active ON roles (uuid, deleted_at);
CREATE UNIQUE INDEX idx_roles_name_active ON roles (name, deleted_at);
CREATE INDEX idx_roles_status ON roles (status, deleted_at);

-- =========================================================
-- 8. USER ROLES (RBAC)
-- =========================================================
CREATE TABLE user_roles (
                            id BIGINT PRIMARY KEY AUTO_INCREMENT,
                            uuid CHAR(36) UNIQUE NOT NULL,
                            user_id BIGINT NOT NULL,
                            role_id BIGINT NOT NULL,
                            created_by BIGINT,
                            status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
                            deleted_at TIMESTAMP NULL,
                            CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                            CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_user_role_unique ON user_roles (user_id, role_id, deleted_at);
CREATE INDEX idx_user_roles_lookup ON user_roles (user_id, role_id, status, deleted_at);

-- =========================================================
-- 9. ATTENDANCE SUMMARY
-- =========================================================
CREATE TABLE attendance_summary (
                                    id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                    uuid CHAR(36) UNIQUE NOT NULL,
                                    user_id BIGINT NOT NULL,
                                    log_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    total_work_in_sec BIGINT NULL,
                                    total_brake_in_sec BIGINT NULL,
                                    log_status ENUM('PENDING', 'COMPLETED') NOT NULL,
                                    status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
                                    deleted_at TIMESTAMP NULL,
                                    CONSTRAINT fk_attendance_summary_employee FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_attendance_summary_user_time ON attendance_summary (user_id, log_date);

-- =========================================================
-- 10. ATTENDANCE LOGS
-- =========================================================
CREATE TABLE attendance_logs (
                                 id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                 uuid CHAR(36) UNIQUE NOT NULL,
                                 user_id BIGINT NOT NULL,
                                 attendance_summary_id BIGINT NOT NULL,
                                 log_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                 log_type ENUM('IN', 'OUT', 'BREAK_IN', 'BREAK_OUT', 'TIME_OUT') NOT NULL,
                                 source ENUM('mobile', 'biometric', 'admin') NOT NULL,
                                 latitude DECIMAL(10, 7) NULL,
                                 longitude DECIMAL(10, 7) NULL,
                                 device_id VARCHAR(100),
                                 remarks TEXT,
                                 status ENUM ('active','inactive','revoked','blocked','archived') NOT NULL DEFAULT 'active',
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP,
                                 deleted_at TIMESTAMP NULL,
                                 CONSTRAINT fk_logs_employee FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                 CONSTRAINT fk_log_attendance_summary_id FOREIGN KEY (attendance_summary_id) REFERENCES attendance_summary(id) ON DELETE CASCADE
);

CREATE INDEX idx_logs_employee_time ON attendance_logs (user_id, log_time);



