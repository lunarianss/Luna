-- ----------------------------
-- Table structure for tenant_default_models
-- ----------------------------
DROP TABLE IF EXISTS `tenant_default_models`;
CREATE TABLE `tenant_default_models` (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(40) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    created_by CHAR(36),
    updated_by CHAR(36)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;