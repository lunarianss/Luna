-- ----------------------------
-- Table structure for accounts
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE accounts (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255),
    password_salt VARCHAR(255),
    avatar VARCHAR(255),
    interface_language VARCHAR(255),
    interface_theme VARCHAR(255),
    timezone VARCHAR(255),
    last_login_at int(10) NULL,
    last_login_ip VARCHAR(255),
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    initialized_at TIMESTAMP NULL,
    created_at int(10) NOT NULL,
    updated_at int(10)  NOT NULL,
    last_active_at int(10)  NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;


-- ----------------------------
-- Table structure for tenant_account_joins
-- ----------------------------
DROP TABLE IF EXISTS `tenant_account_joins`;
CREATE TABLE tenant_account_joins (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    account_id CHAR(36) NOT NULL,
    role VARCHAR(16) NOT NULL DEFAULT 'normal',
    invited_by CHAR(36),
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL,
    current bit(1) NOT NULL DEFAULT 0
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;


-- ----------------------------
-- Table structure for tenants
-- ----------------------------
DROP TABLE IF EXISTS `tenants`;
CREATE TABLE tenants (
    id CHAR(36) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    encrypt_public_key TEXT,
    plan VARCHAR(255) NOT NULL DEFAULT 'basic',
    status VARCHAR(255) NOT NULL DEFAULT 'normal',
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL,
    custom_config TEXT,
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;
