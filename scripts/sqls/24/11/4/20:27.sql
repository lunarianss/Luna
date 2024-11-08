
-- ----------------------------
-- Table structure for provider_models
-- ----------------------------
DROP TABLE IF EXISTS `provider_models`;
CREATE TABLE `provider_models`  (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(40) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    encrypted_config TEXT,
    is_valid bit(1) NOT NULL DEFAULT 0,
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

