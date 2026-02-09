CREATE TABLE `users` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `name` varchar(120) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `phone` varchar(15),
  `country_code` int,
  `password_hash` varchar(255) NOT NULL,
  `employee_id` varchar(50) UNIQUE,
  `group_id` bigint NOT NULL,
  `designation_id` bigint NOT NULL,
  `department_id` bigint NOT NULL,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `user_roles` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `user_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `roles` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `name` varchar(120) NOT NULL,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `user_groups` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `name` varchar(120) NOT NULL,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `designations` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `designation` varchar(120) NOT NULL,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `departments` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `department` varchar(120) NOT NULL,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `password_resets` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `user_id` bigint NOT NULL,
  `token_hash` varchar(255) NOT NULL,
  `expires_at` timestamp NOT NULL,
  `used_at` timestamp,
  `requested_ip` varchar(45),
  `user_agent` varchar(255),
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `clients` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `name` varchar(150) NOT NULL,
  `email` varchar(255),
  `phone` varchar(20),
  `address` text,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_by` bigint,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `projects` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `project_name` varchar(120) NOT NULL,
  `project_code` varchar(120) UNIQUE NOT NULL,
  `client_id` bigint,
  `client_name_snapshot` varchar(150),
  `lifecycle_status` ENUM ('todo', 'in_progress', 'in_review', 'done', 'closed', 'on_hold') NOT NULL DEFAULT 'todo',
  `start_date` date,
  `internal_start_date` date,
  `end_date` date,
  `internal_end_date` date,
  `estimated_hours` decimal(7,2) NOT NULL COMMENT '>= 0',
  `internal_estimated_hours` decimal(7,2) NOT NULL COMMENT '>= 0',
  `primary_repo_url` varchar(500),
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_by` bigint,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `project_assets` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `project_id` bigint NOT NULL,
  `title` varchar(150) NOT NULL,
  `asset_type` varchar(20) NOT NULL COMMENT 'text|file|image|link|secret',
  `description` text,
  `file_id` bigint,
  `content` text,
  `is_private` boolean DEFAULT false,
  `created_by` bigint NOT NULL,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `files` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `storage_key` varchar(500) UNIQUE NOT NULL,
  `file_status` varchar(20) NOT NULL DEFAULT 'uploaded',
  `original_name` varchar(255),
  `mime_type` varchar(120),
  `size_bytes` bigint,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `ticket_types` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `code` varchar(30) UNIQUE NOT NULL,
  `name` varchar(120) NOT NULL,
  `description` text,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `project_managers` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `user_id` bigint NOT NULL,
  `project_id` bigint NOT NULL,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `project_members` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `user_id` bigint NOT NULL,
  `project_id` bigint NOT NULL,
  `allocated_hours` decimal(7,2) NOT NULL,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `tickets` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `project_id` bigint NOT NULL,
  `code` varchar(30) UNIQUE NOT NULL,
  `title` varchar(120) NOT NULL,
  `description` text,
  `branch` varchar(120),
  `ticket_type_id` bigint NOT NULL,
  `assigned_to` bigint,
  `assigned_by` bigint,
  `start_date` date,
  `internal_start_date` date,
  `end_date` date,
  `internal_end_date` date,
  `estimated_hours` decimal(7,2) NOT NULL COMMENT '>= 0',
  `internal_estimated_hours` decimal(7,2) NOT NULL COMMENT '>= 0',
  `lifecycle_status` ENUM ('todo', 'in_progress', 'in_review', 'testing', 'done', 'closed', 'reopened') NOT NULL DEFAULT 'todo',
  `priority` ENUM ('low', 'medium', 'high', 'urgent') DEFAULT 'medium',
  `parent_id` bigint,
  `due_date` date,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `time_entries` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `ticket_id` bigint NOT NULL,
  `project_id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `work_date` date NOT NULL,
  `hours` decimal(5,2) NOT NULL COMMENT '>= 0',
  `description` text,
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);

CREATE TABLE `ticket_history` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `ticket_id` bigint NOT NULL,
  `changed_by` bigint NOT NULL,
  `field_name` varchar(50) NOT NULL,
  `old_value` text,
  `new_value` text,
  `note` text,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `activity_logs` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `uuid` varchar(36) UNIQUE NOT NULL,
  `user_id` bigint,
  `action` varchar(50) NOT NULL,
  `entity_type` varchar(50),
  `entity_id` bigint,
  `metadata` json,
  `ip_address` varchar(45),
  `user_agent` varchar(255),
  `created_by` bigint,
  `status` ENUM ('active', 'inactive', 'blocked', 'archived') NOT NULL DEFAULT 'active',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp,
  `deleted_at` timestamp NULL DEFAULT NULL
);
CREATE TABLE `sessions` (
                            `id` bigint PRIMARY KEY AUTO_INCREMENT,
                            `uuid` varchar(36) UNIQUE NOT NULL,
                            `user_id` bigint NOT NULL,
                            `jti` varchar(36) UNIQUE NOT NULL,
                            `refresh_token_hash` varchar(255) NOT NULL,
                            `refresh_expires_at` timestamp NOT NULL,
                            `rotated_at` timestamp NULL DEFAULT NULL,
                            `status` ENUM ('active', 'inactive', 'blocked', 'revoked') NOT NULL DEFAULT 'active',
                            `device_id` varchar(100),
                            `device_name` varchar(100),
                            `user_agent` text,
                            `ip_address` varchar(45),
                            `last_used_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
                            `created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
                            `updated_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
                            `deleted_at` timestamp NULL DEFAULT NULL
);
CREATE INDEX `users_index_0` ON `users` (`uuid`);

CREATE INDEX `users_index_1` ON `users` (`email`);

CREATE INDEX `users_index_2` ON `users` (`employee_id`);

CREATE INDEX `users_index_3` ON `users` (`group_id`);

CREATE INDEX `users_index_4` ON `users` (`department_id`);

CREATE INDEX `users_index_5` ON `users` (`designation_id`);

CREATE INDEX `users_index_6` ON `users` (`deleted_at`);

CREATE UNIQUE INDEX `user_roles_index_7` ON `user_roles` (`user_id`, `role_id`);

CREATE INDEX `user_roles_index_8` ON `user_roles` (`role_id`);

CREATE INDEX `password_resets_index_9` ON `password_resets` (`user_id`);

CREATE UNIQUE INDEX `password_resets_index_10` ON `password_resets` (`token_hash`);

CREATE INDEX `clients_index_11` ON `clients` (`name`);

CREATE INDEX `projects_index_12` ON `projects` (`project_code`);

CREATE INDEX `projects_index_13` ON `projects` (`lifecycle_status`);

CREATE INDEX `projects_index_14` ON `projects` (`created_by`);

CREATE INDEX `projects_index_15` ON `projects` (`deleted_at`);

CREATE INDEX `project_assets_index_16` ON `project_assets` (`project_id`);

CREATE INDEX `project_assets_index_17` ON `project_assets` (`asset_type`);

CREATE INDEX `project_assets_index_18` ON `project_assets` (`file_id`);

CREATE INDEX `project_assets_index_19` ON `project_assets` (`created_by`);

CREATE INDEX `files_index_20` ON `files` (`file_status`);

CREATE INDEX `files_index_21` ON `files` (`storage_key`);

CREATE INDEX `project_managers_index_22` ON `project_managers` (`user_id`);

CREATE UNIQUE INDEX `project_managers_index_23` ON `project_managers` (`user_id`, `project_id`);

CREATE INDEX `project_managers_index_24` ON `project_managers` (`project_id`);

CREATE INDEX `project_members_index_25` ON `project_members` (`user_id`);

CREATE UNIQUE INDEX `project_members_index_26` ON `project_members` (`user_id`, `project_id`);

CREATE INDEX `project_members_index_27` ON `project_members` (`project_id`);

CREATE INDEX `tickets_index_28` ON `tickets` (`project_id`);

CREATE INDEX `tickets_index_29` ON `tickets` (`assigned_to`);

CREATE INDEX `tickets_index_30` ON `tickets` (`lifecycle_status`);

CREATE INDEX `tickets_index_31` ON `tickets` (`ticket_type_id`);

CREATE INDEX `tickets_index_32` ON `tickets` (`priority`);

CREATE INDEX `tickets_index_33` ON `tickets` (`project_id`, `lifecycle_status`);

CREATE INDEX `time_entries_index_34` ON `time_entries` (`ticket_id`);

CREATE INDEX `time_entries_index_35` ON `time_entries` (`project_id`);

CREATE INDEX `time_entries_index_36` ON `time_entries` (`user_id`);

CREATE INDEX `time_entries_index_37` ON `time_entries` (`work_date`);

CREATE INDEX `time_entries_index_38` ON `time_entries` (`user_id`, `work_date`);

CREATE INDEX `time_entries_index_39` ON `time_entries` (`project_id`, `work_date`);

CREATE INDEX `ticket_history_index_40` ON `ticket_history` (`ticket_id`);

CREATE INDEX `ticket_history_index_41` ON `ticket_history` (`changed_by`);

CREATE INDEX `ticket_history_index_42` ON `ticket_history` (`created_at`);

CREATE INDEX `activity_logs_index_43` ON `activity_logs` (`user_id`);

CREATE INDEX `activity_logs_index_44` ON `activity_logs` (`entity_type`, `entity_id`);

CREATE INDEX `activity_logs_index_45` ON `activity_logs` (`created_at`);

-- Additional performance indexes

CREATE INDEX `users_index_46` ON `users` (`status`);

CREATE INDEX `users_index_47` ON `users` (`status`, `deleted_at`);

CREATE INDEX `users_index_48` ON `users` (`phone`);

CREATE INDEX `user_roles_index_49` ON `user_roles` (`status`);

CREATE INDEX `user_roles_index_50` ON `user_roles` (`deleted_at`);

CREATE INDEX `user_roles_index_51` ON `user_roles` (`status`, `deleted_at`);

CREATE INDEX `roles_index_52` ON `roles` (`name`);

CREATE INDEX `roles_index_53` ON `roles` (`status`);

CREATE INDEX `groups_index_54` ON `user_groups` (`name`);

CREATE INDEX `groups_index_55` ON `user_groups` (`status`);

CREATE INDEX `designations_index_56` ON `designations` (`designation`);

CREATE INDEX `designations_index_57` ON `designations` (`status`);

CREATE INDEX `departments_index_58` ON `departments` (`department`);

CREATE INDEX `departments_index_59` ON `departments` (`status`);

CREATE INDEX `password_resets_index_60` ON `password_resets` (`expires_at`);

CREATE INDEX `password_resets_index_61` ON `password_resets` (`status`);

CREATE INDEX `clients_index_62` ON `clients` (`status`);

CREATE INDEX `clients_index_63` ON `clients` (`deleted_at`);

CREATE INDEX `clients_index_64` ON `clients` (`status`, `deleted_at`);

CREATE INDEX `projects_index_65` ON `projects` (`client_id`);

CREATE INDEX `projects_index_66` ON `projects` (`status`);

CREATE INDEX `projects_index_67` ON `projects` (`status`, `deleted_at`);

CREATE INDEX `projects_index_68` ON `projects` (`client_id`, `status`);

CREATE INDEX `project_assets_index_69` ON `project_assets` (`project_id`, `is_private`, `status`);

CREATE INDEX `project_assets_index_70` ON `project_assets` (`deleted_at`);

CREATE INDEX `files_index_71` ON `files` (`deleted_at`);

CREATE INDEX `ticket_types_index_72` ON `ticket_types` (`status`);

CREATE INDEX `project_managers_index_73` ON `project_managers` (`status`);

CREATE INDEX `project_managers_index_74` ON `project_managers` (`deleted_at`);

CREATE INDEX `project_members_index_75` ON `project_members` (`status`);

CREATE INDEX `project_members_index_76` ON `project_members` (`deleted_at`);

CREATE INDEX `tickets_index_77` ON `tickets` (`code`);

CREATE INDEX `tickets_index_78` ON `tickets` (`due_date`);

CREATE INDEX `tickets_index_79` ON `tickets` (`parent_id`);

CREATE INDEX `tickets_index_80` ON `tickets` (`deleted_at`);

CREATE INDEX `tickets_index_81` ON `tickets` (`status`);

CREATE INDEX `tickets_index_82` ON `tickets` (`assigned_to`, `lifecycle_status`, `deleted_at`);

CREATE INDEX `tickets_index_83` ON `tickets` (`project_id`, `assigned_to`);

CREATE INDEX `tickets_index_84` ON `tickets` (`status`, `deleted_at`);

CREATE INDEX `time_entries_index_85` ON `time_entries` (`user_id`, `project_id`, `work_date`);

CREATE INDEX `time_entries_index_86` ON `time_entries` (`deleted_at`);

CREATE INDEX `time_entries_index_87` ON `time_entries` (`status`);

CREATE INDEX `activity_logs_index_88` ON `activity_logs` (`user_id`, `created_at`);

CREATE INDEX `activity_logs_index_89` ON `activity_logs` (`deleted_at`);

-- Additional missing indexes

CREATE INDEX `users_index_90` ON `users` (`name`);

CREATE INDEX `users_index_91` ON `users` (`created_by`);

CREATE INDEX `users_index_92` ON `users` (`country_code`);

CREATE INDEX `roles_index_93` ON `roles` (`deleted_at`);

CREATE INDEX `groups_index_94` ON `user_groups` (`deleted_at`);

CREATE INDEX `designations_index_95` ON `designations` (`deleted_at`);

CREATE INDEX `departments_index_96` ON `departments` (`deleted_at`);

CREATE INDEX `clients_index_97` ON `clients` (`email`);

CREATE INDEX `projects_index_98` ON `projects` (`project_name`);

CREATE INDEX `projects_index_99` ON `projects` (`start_date`);

CREATE INDEX `projects_index_100` ON `projects` (`end_date`);

CREATE INDEX `projects_index_101` ON `projects` (`start_date`, `end_date`);

CREATE INDEX `files_index_102` ON `files` (`status`);

CREATE INDEX `tickets_index_103` ON `tickets` (`assigned_by`);

CREATE INDEX `tickets_index_104` ON `tickets` (`created_by`);

CREATE INDEX `tickets_index_105` ON `tickets` (`start_date`);

CREATE INDEX `tickets_index_106` ON `tickets` (`end_date`);

CREATE INDEX `tickets_index_107` ON `tickets` (`branch`);

CREATE INDEX `tickets_index_108` ON `tickets` (`project_id`, `lifecycle_status`, `priority`);

CREATE INDEX `tickets_index_109` ON `tickets` (`assigned_to`, `due_date`);

CREATE INDEX `ticket_history_index_110` ON `ticket_history` (`field_name`);

CREATE INDEX `time_entries_index_111` ON `time_entries` (`ticket_id`, `work_date`);

CREATE INDEX `activity_logs_index_112` ON `activity_logs` (`action`);

CREATE INDEX `password_resets_index_113` ON `password_resets` (`used_at`);

CREATE INDEX `password_resets_index_114` ON `password_resets` (`expires_at`, `used_at`);

CREATE INDEX `project_members_index_115` ON `project_members` (`project_id`, `status`);

CREATE INDEX `project_managers_index_116` ON `project_managers` (`project_id`, `status`);

ALTER TABLE `users` ADD FOREIGN KEY (`group_id`) REFERENCES `user_groups` (`id`);

ALTER TABLE `users` ADD FOREIGN KEY (`designation_id`) REFERENCES `designations` (`id`);

ALTER TABLE `users` ADD FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`);

ALTER TABLE `users` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`);

ALTER TABLE `user_roles` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `roles` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `user_groups` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `designations` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `departments` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `password_resets` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `password_resets` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `clients` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `projects` ADD FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`);

ALTER TABLE `projects` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `project_assets` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `project_assets` ADD FOREIGN KEY (`file_id`) REFERENCES `files` (`id`);

ALTER TABLE `project_assets` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `files` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `ticket_types` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `project_managers` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `project_managers` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `project_managers` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `project_members` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `project_members` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `project_members` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `tickets` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `tickets` ADD FOREIGN KEY (`ticket_type_id`) REFERENCES `ticket_types` (`id`);

ALTER TABLE `tickets` ADD FOREIGN KEY (`assigned_to`) REFERENCES `users` (`id`);

ALTER TABLE `tickets` ADD FOREIGN KEY (`assigned_by`) REFERENCES `users` (`id`);

ALTER TABLE `tickets` ADD FOREIGN KEY (`parent_id`) REFERENCES `tickets` (`id`);

ALTER TABLE `tickets` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `time_entries` ADD FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`);

ALTER TABLE `time_entries` ADD FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

ALTER TABLE `time_entries` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `time_entries` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);

ALTER TABLE `ticket_history` ADD FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`);

ALTER TABLE `ticket_history` ADD FOREIGN KEY (`changed_by`) REFERENCES `users` (`id`);

ALTER TABLE `activity_logs` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `activity_logs` ADD FOREIGN KEY (`created_by`) REFERENCES `users` (`id`);



CREATE INDEX `idx_sessions_user` ON `sessions` (`user_id`);

CREATE INDEX `idx_sessions_refresh_hash` ON `sessions` (`refresh_token_hash`);

CREATE INDEX `idx_sessions_refresh_exp` ON `sessions` (`refresh_expires_at`);

CREATE INDEX `idx_sessions_status` ON `sessions` (`status`);

CREATE INDEX `idx_sessions_deleted` ON `sessions` (`deleted_at`);

CREATE INDEX `idx_sessions_refresh_status` ON `sessions` (`refresh_token_hash`, `status`);

ALTER TABLE `sessions` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);


-- Triggers for auto-updating updated_at timestamps

DELIMITER $$

CREATE TRIGGER `users_updated_at` BEFORE UPDATE ON `users`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `user_roles_updated_at` BEFORE UPDATE ON `user_roles`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `roles_updated_at` BEFORE UPDATE ON `roles`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `groups_updated_at` BEFORE UPDATE ON `user_groups`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `designations_updated_at` BEFORE UPDATE ON `designations`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `departments_updated_at` BEFORE UPDATE ON `departments`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `password_resets_updated_at` BEFORE UPDATE ON `password_resets`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `clients_updated_at` BEFORE UPDATE ON `clients`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `projects_updated_at` BEFORE UPDATE ON `projects`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `project_assets_updated_at` BEFORE UPDATE ON `project_assets`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `files_updated_at` BEFORE UPDATE ON `files`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `ticket_types_updated_at` BEFORE UPDATE ON `ticket_types`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `project_managers_updated_at` BEFORE UPDATE ON `project_managers`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `project_members_updated_at` BEFORE UPDATE ON `project_members`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `tickets_updated_at` BEFORE UPDATE ON `tickets`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `time_entries_updated_at` BEFORE UPDATE ON `time_entries`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$

CREATE TRIGGER `activity_logs_updated_at` BEFORE UPDATE ON `activity_logs`
FOR EACH ROW
BEGIN
  SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$



  CREATE TRIGGER `session_updated_at` BEFORE UPDATE ON `sessions`
      FOR EACH ROW
  BEGIN
      SET NEW.updated_at = CURRENT_TIMESTAMP;
END$$


      DELIMITER ;

-- ============================================
-- PERFORMANCE OPTIMIZATION RECOMMENDATIONS
-- ============================================

-- 1. PARTITIONING (For large-scale data)
-- Uncomment and adjust based on your data volume

/*
-- Partition activity_logs by month
ALTER TABLE `activity_logs` 
PARTITION BY RANGE (YEAR(created_at) * 100 + MONTH(created_at)) (
  PARTITION p202401 VALUES LESS THAN (202402),
  PARTITION p202402 VALUES LESS THAN (202403),
  PARTITION p202403 VALUES LESS THAN (202404),
  -- Add more partitions as needed
  PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- Partition ticket_history by year
ALTER TABLE `ticket_history`
PARTITION BY RANGE (YEAR(created_at)) (
  PARTITION p2024 VALUES LESS THAN (2025),
  PARTITION p2025 VALUES LESS THAN (2026),
  PARTITION p2026 VALUES LESS THAN (2027),
  PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- Partition time_entries by year
ALTER TABLE `time_entries`
PARTITION BY RANGE (YEAR(work_date)) (
  PARTITION p2024 VALUES LESS THAN (2025),
  PARTITION p2025 VALUES LESS THAN (2026),
  PARTITION p2026 VALUES LESS THAN (2027),
  PARTITION p_future VALUES LESS THAN MAXVALUE
);
*/

-- 2. ARCHIVAL TABLES (For historical data)

/*
CREATE TABLE `activity_logs_archive` LIKE `activity_logs`;
CREATE TABLE `ticket_history_archive` LIKE `ticket_history`;
CREATE TABLE `time_entries_archive` LIKE `time_entries`;

-- Archive old activity logs (older than 1 year)
-- Run this periodically via cron job
INSERT INTO `activity_logs_archive` 
SELECT * FROM `activity_logs` 
WHERE created_at < DATE_SUB(NOW(), INTERVAL 1 YEAR);

DELETE FROM `activity_logs` 
WHERE created_at < DATE_SUB(NOW(), INTERVAL 1 YEAR);
*/

-- 3. CLEANUP QUERIES (Run periodically)

/*
-- Delete expired password reset tokens (older than 7 days)
DELETE FROM `password_resets` 
WHERE expires_at < DATE_SUB(NOW(), INTERVAL 7 DAY) 
AND used_at IS NOT NULL;

-- Archive completed projects older than 2 years
UPDATE `projects` 
SET status = 'archived' 
WHERE lifecycle_status = 'closed' 
AND updated_at < DATE_SUB(NOW(), INTERVAL 2 YEAR)
AND status != 'archived';
*/

-- 4. DATA TYPE OPTIMIZATIONS (Breaking changes - use with caution)

/*
-- Convert UUID from varchar(36) to binary(16) for better performance
-- This requires application-level changes to convert UUIDs

ALTER TABLE `users` MODIFY `uuid` binary(16) NOT NULL;
ALTER TABLE `user_roles` MODIFY `uuid` binary(16) NOT NULL;
-- Repeat for all tables with uuid columns

-- Convert ENUM status to TINYINT for more flexibility
-- 1=active, 2=inactive, 3=blocked, 4=archived

ALTER TABLE `users` MODIFY `status` tinyint NOT NULL DEFAULT 1;
-- Repeat for all tables with status columns
*/

-- 5. COVERING INDEXES (For specific frequent queries)
-- Add these based on your actual query patterns

/*
-- Example: If you frequently query user details with role info
CREATE INDEX `users_covering_idx` ON `users` 
  (`id`, `email`, `name`, `status`, `group_id`, `department_id`, `designation_id`);

-- Example: For ticket listing with assignment info
CREATE INDEX `tickets_covering_idx` ON `tickets` 
  (`project_id`, `lifecycle_status`, `assigned_to`, `priority`, `due_date`, `status`);
*/

-- 6. QUERY OPTIMIZATION TIPS
-- - Use EXPLAIN to analyze slow queries
-- - Avoid SELECT * - specify only needed columns
-- - Use LIMIT for pagination
-- - Consider read replicas for reporting queries
-- - Use connection pooling in application
-- - Enable query cache if using MySQL 5.7 or lower
-- - Monitor slow query log regularly
