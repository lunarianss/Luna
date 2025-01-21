-- ----------------------------
-- Table structure for upload_files
-- ----------------------------
DROP TABLE IF EXISTS `upload_files`;
CREATE TABLE upload_files (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    storage_type VARCHAR(255) NOT NULL,
    `key` VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    size INT NOT NULL,
    extension VARCHAR(255) NOT NULL,
    mime_type VARCHAR(255),
    created_by_role VARCHAR(255) NOT NULL DEFAULT 'account',
    created_by CHAR(36) NOT NULL,
    created_at INT(10) NOT NULL,
    used BIT(1) NOT NULL DEFAULT 0,
    used_by CHAR(36),
    used_at INT(10),
    hash VARCHAR(255)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE INDEX upload_file_tenant_idx ON upload_files (tenant_id);