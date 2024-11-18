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

ALTER TABLE providers ADD INDEX idx_tenant_id (tenant_id);


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

ALTER TABLE provider_models ADD INDEX idx_tenant_id (tenant_id);

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


ALTER TABLE apps ADD INDEX idx_tenant_id (tenant_id);

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


ALTER TABLE tenant_default_models ADD INDEX idx_tenant_id (tenant_id);
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

ALTER TABLE app_model_configs ADD INDEX idx_app_id (app_id);

-- ----------------------------
-- Table structure for conversation
-- ----------------------------
DROP TABLE IF EXISTS `conversations`;
CREATE TABLE `conversations` (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    app_model_config_id CHAR(36),
    model_provider VARCHAR(255),
    override_model_configs TEXT,
    model_id VARCHAR(255),
    mode VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    summary TEXT,
    inputs TEXT,
    introduction TEXT,
    system_instruction TEXT,
    system_instruction_tokens INT NOT NULL DEFAULT 0,
    status VARCHAR(255) NOT NULL,
    invoke_from VARCHAR(255),
    from_source VARCHAR(255) NOT NULL,
    from_end_user_id CHAR(36),
    from_account_id CHAR(36),
    read_at int(10),
    read_account_id CHAR(36),
    dialogue_count INT DEFAULT 0 NOT NULL,
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL,
    is_deleted bit(1) NOT NULL DEFAULT 0
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE conversations ADD INDEX conversation_app_from_user_idx (app_id, from_source, from_end_user_id);

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    model_provider VARCHAR(255),
    model_id VARCHAR(255),
    override_model_configs TEXT,
    conversation_id CHAR(36) NOT NULL,
    inputs TEXT,
    query TEXT NOT NULL,
    message TEXT NOT NULL,
    answer TEXT NOT NULL,
    message_tokens INT DEFAULT 0 NOT NULL,
    message_price_unit DECIMAL(10,7) DEFAULT 0.001 NOT NULL,
    message_unit_price DECIMAL(10,4) NOT NULL,
    answer_price_unit DECIMAL(10,7) DEFAULT 0.001 NOT NULL,
    answer_tokens INT DEFAULT 0 NOT NULL,
    answer_unit_price DECIMAL(10,4) NOT NULL,
    provider_response_latency DOUBLE DEFAULT 0 NOT NULL,
    total_price DECIMAL(15,7),
    currency VARCHAR(255) NOT NULL,
    from_source VARCHAR(255) NOT NULL,
    from_end_user_id CHAR(36) NOT NULL,
    from_account_id CHAR(36) NOT NULL,
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL,
    agent_based bit(1) DEFAULT 0 NOT NULL,
    workflow_run_id  CHAR(36) NOT NULL,
    status VARCHAR(255) DEFAULT 'normal' NOT NULL,
    error TEXT,
    message_metadata TEXT,
    invoke_from VARCHAR(255),
    parent_message_id  CHAR(36) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;


ALTER TABLE messages ADD INDEX message_app_id_idx (app_id, created_at);
ALTER TABLE messages ADD INDEX message_conversation_id_idx (conversation_id);
ALTER TABLE messages ADD INDEX message_end_user_idx (app_id, from_source, from_end_user_id);
ALTER TABLE messages ADD INDEX message_account_idx (app_id, from_source, from_account_id);
ALTER TABLE messages ADD INDEX message_workflow_run_id_idx (conversation_id, from_source, workflow_run_id);


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

/* 
-- 为 tenant_id 创建普通索引
CREATE INDEX idx_tenant_id ON tenant_account_joins (tenant_id);

-- 为 account_id 创建普通索引
CREATE INDEX idx_account_id ON tenant_account_joins (account_id); */

ALTER TABLE tenant_account_joins ADD INDEX idx_tenant_id (tenant_id);
ALTER TABLE tenant_account_joins ADD INDEX idx_account_id (account_id);




/* 
CREATE UNIQUE INDEX unique_tenant_account ON tenant_account_joins (tenant_id, account_id); */

ALTER TABLE tenant_account_joins ADD UNIQUE INDEX unique_tenant_account (tenant_id, account_id);



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
    custom_config TEXT
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;


