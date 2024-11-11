-- ----------------------------
-- Table structure for conversation
-- ----------------------------
CREATE TABLE `conversation` (
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


-- ----------------------------
-- Table structure for message
-- ----------------------------
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
