CREATE TABLE `users` (
                         `id` bigint PRIMARY KEY AUTO_INCREMENT,
                         `uuid` varchar(255) UNIQUE NOT NULL,
                         `name` varchar(255) NOT NULL,
                         `email` varchar(255) UNIQUE NOT NULL,
                         `phone` varchar(255),
                         `country_code` integer,
                         `password_hash` varchar(255) NOT NULL,
                         `employee_id` varchar(255) UNIQUE,
                         `group_id` bigint,
                         `designation_id` bigint,
                         `department_id` bigint,
                         `is_admin` boolean DEFAULT false,
                         `created_by` bigint,
                         `status` ENUM ('active', 'inactive','revoked', 'blocked', 'archived') NOT NULL DEFAULT 'active',
                         `created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
                         `updated_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
                         `deleted_at` timestamp NULL DEFAULT null
);

CREATE TABLE `sessions` (
                            `id` bigint PRIMARY KEY AUTO_INCREMENT,
                            `uuid` varchar(36) UNIQUE NOT NULL,
                            `user_id` bigint NOT NULL,
                            `jti` varchar(36) UNIQUE NOT NULL,
                            `refresh_token_hash` varchar(255) NOT NULL,
                            `refresh_expires_at` timestamp NOT NULL,
                            `rotated_at` timestamp NULL DEFAULT null,
                            `status` ENUM ('active', 'inactive', 'revoked','blocked', 'archived') NOT NULL DEFAULT 'active',
                            `device_id` varchar(100),
                            `device_name` varchar(100),
                            `user_agent` text,
                            `ip_address` varchar(45),
                            `last_used_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
                            `created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
                            `updated_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
                            `deleted_at` timestamp NULL DEFAULT null
);

-- INDEX
CREATE INDEX `idx_users_email` ON `users` (`email`);

CREATE INDEX `idx_users_uuid` ON `users` (`uuid`);

CREATE INDEX `idx_users_group` ON `users` (`group_id`);

CREATE INDEX `idx_users_dept` ON `users` (`department_id`);

CREATE INDEX `idx_users_active_status` ON `users` (`status`, `deleted_at`);

CREATE INDEX `idx_users_employee_id` ON `users` (`employee_id`);

CREATE UNIQUE INDEX `idx_sessions_jti` ON `sessions` (`jti`);

CREATE UNIQUE INDEX `idx_sessions_rt_hash` ON `sessions` (`refresh_token_hash`);

CREATE INDEX `idx_sessions_user_status` ON `sessions` (`user_id`, `status`);

CREATE INDEX `idx_sessions_user_device` ON `sessions` (`user_id`, `device_id`);

CREATE INDEX `idx_sessions_expiry` ON `sessions` (`refresh_expires_at`);

--  FG

ALTER TABLE `sessions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);