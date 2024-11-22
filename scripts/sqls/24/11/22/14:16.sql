
-- ----------------------------
-- Table structure for end_users
-- ----------------------------
DROP TABLE IF EXISTS `end_users`;
CREATE TABLE end_users (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    app_id CHAR(36),
    type VARCHAR(255) NOT NULL,
    external_user_id VARCHAR(255),
    name VARCHAR(255),
    is_anonymous bit(1) NOT NULL DEFAULT 1,
    session_id VARCHAR(255) NOT NULL,
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE end_users ADD INDEX end_user_session_id_idx (session_id, type);
ALTER TABLE end_users ADD INDEX end_user_tenant_session_id_idx (tenant_id, session_id, type);


-- ----------------------------
-- Table structure for installed_apps
-- ----------------------------
DROP TABLE IF EXISTS `installed_apps`;
CREATE TABLE installed_apps (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    app_id CHAR(36) NOT NULL,
    app_owner_tenant_id CHAR(36) NOT NULL,
    position int NOT NULL DEFAULT 0,
    is_pinned bit(1) NOT NULL DEFAULT 0,
    last_used_at int(10) NOT NULL,
    created_at int(10) NOT NULL,
    UNIQUE (tenant_id, app_id)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE installed_apps ADD INDEX installed_app_tenant_id_idx (tenant_id);
ALTER TABLE installed_apps ADD INDEX installed_app_app_id_idx (app_id);



-- ----------------------------
-- Table structure for sites
-- ----------------------------
DROP TABLE IF EXISTS `sites`;
CREATE TABLE sites (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    icon_type VARCHAR(255),
    icon VARCHAR(255),
    icon_background VARCHAR(255),
    description TEXT,
    default_language VARCHAR(255) NOT NULL,
    chat_color_theme VARCHAR(255),
    chat_color_theme_inverted bit(1) NOT NULL DEFAULT 0,
    copyright VARCHAR(255),
    privacy_policy VARCHAR(255),
    show_workflow_steps bit(1) NOT NULL DEFAULT 1,
    use_icon_as_answer_icon bit(1) NOT NULL DEFAULT 0,
    custom_disclaimer VARCHAR(255),
    customize_domain VARCHAR(255),
    customize_token_strategy VARCHAR(255) NOT NULL,
    prompt_public bit(1) NOT NULL DEFAULT 0,
    status VARCHAR(255) NOT NULL DEFAULT 'normal',
    created_by CHAR(36),
    created_at int(10) NOT NULL,
    updated_by CHAR(36),
    updated_at int(10) NOT NULL,
    code VARCHAR(255),
    UNIQUE (code)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE sites ADD INDEX site_app_id_idx (app_id);
ALTER TABLE sites ADD INDEX site_code_idx (code, status);
