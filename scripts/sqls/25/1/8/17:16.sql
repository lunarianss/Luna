-- ----------------------------
-- Table structure for tool_files
-- ----------------------------
DROP TABLE IF EXISTS `tool_files`;
CREATE TABLE `tool_files` (
    `id` CHAR(36) NOT NULL PRIMARY KEY,
    `user_id` CHAR(36) NOT NULL,
    `tenant_id` CHAR(36) NOT NULL,
    `conversation_id` CHAR(36) DEFAULT '',
    `file_key` VARCHAR(255) NOT NULL,
    `mimetype` VARCHAR(255) NOT NULL,
    `original_url` VARCHAR(2048) DEFAULT '',
    `name` VARCHAR(255) NOT NULL DEFAULT '',
    `size` INT NOT NULL DEFAULT -1
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

ALTER TABLE `tool_files` ADD INDEX `tool_file_conversation_id_idx` (`conversation_id`);


-- ----------------------------
-- Table structure for message_files
-- ----------------------------
DROP TABLE IF EXISTS `message_files`;
CREATE TABLE `message_files` (
    `id` CHAR(36) NOT NULL PRIMARY KEY,
    `message_id` CHAR(36) NOT NULL,
    `type` VARCHAR(255) NOT NULL,
    `transfer_method` VARCHAR(255) NOT NULL,
    `url` TEXT,
    `belongs_to` VARCHAR(255) DEFAULT '',
    `upload_file_id` CHAR(36) DEFAULT '',
    `created_by_role` VARCHAR(255) NOT NULL,
    `created_by` CHAR(36) NOT NULL,
    `created_at` int(10) NOT NULL
) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- Indexes
CREATE INDEX `message_file_message_idx` ON `message_files` (`message_id`);
CREATE INDEX `message_file_created_by_idx` ON `message_files` (`created_by`);
