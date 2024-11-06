
-- ----------------------------
-- Table structure for apps
-- ----------------------------
DROP TABLE IF EXISTS `apps`;
CREATE TABLE apps (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    mode VARCHAR(255) NOT NULL,
    icon VARCHAR(255),
    icon_background VARCHAR(255),
    app_model_config_id CHAR(36),
    status VARCHAR(255) DEFAULT 'normal' NOT NULL,
    enable_site BIT(1) NOT NULL,
    enable_api BIT(1) NOT NULL,
    api_rpm int DEFAULT 0 NOT NULL,
    api_rph int DEFAULT 0 NOT NULL,
    is_demo bit(1) DEFAULT 0 NOT NULL,
    is_public bit(1) DEFAULT 0 NOT NULL,
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL,
    is_universal BIT(1) DEFAULT 0 NOT NULL,
    workflow_id CHAR(36),
    description TEXT,
    tracing TEXT,
    max_active_requests int,
    icon_type VARCHAR(255),
    created_by CHAR(36),
    updated_by CHAR(36),
    use_icon_as_answer_icon BIT(1) DEFAULT 0 NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;