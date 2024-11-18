-- Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file.

-- ----------------------------
-- Table structure for providers
-- ----------------------------
DROP TABLE IF EXISTS `providers`;
CREATE TABLE `providers`  (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    provider_type VARCHAR(40) NOT NULL DEFAULT 'custom',
    encrypted_config TEXT,
    is_valid bit(1) NOT NULL DEFAULT 0,
    last_used int(10),
    quota_type VARCHAR(40) DEFAULT '',
    quota_limit BIGINT,
    quota_used BIGINT,
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;



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
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

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
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;



-- ----------------------------
-- Table structure for app_model_configs
-- ----------------------------
DROP TABLE IF EXISTS `app_model_configs`;
CREATE TABLE `app_model_configs` (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    provider VARCHAR(255),
    model_id VARCHAR(255),
    configs TEXT,
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL,
    opening_statement TEXT,
    suggested_questions TEXT,
    suggested_questions_after_answer TEXT,
    more_like_this TEXT,
    model TEXT,
    user_input_form TEXT,
    pre_prompt TEXT,
    agent_mode TEXT,
    speech_to_text TEXT,
    sensitive_word_avoidance TEXT,
    retriever_resource TEXT,
    dataset_query_variable VARCHAR(255),
    prompt_type VARCHAR(255) DEFAULT 'simple' NOT NULL,
    chat_prompt_config TEXT,
    completion_prompt_config TEXT,
    dataset_configs TEXT,
    external_data_tools TEXT,
    file_upload TEXT,
    text_to_speech TEXT,
    created_by CHAR(36),
    updated_by CHAR(36)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE apps ADD INDEX idx_tenant_id (tenant_id);