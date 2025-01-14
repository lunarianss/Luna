-- ----------------------------
-- Table structure for message_agent_thoughts
-- ----------------------------
DROP TABLE IF EXISTS `message_agent_thoughts`;
CREATE TABLE message_agent_thoughts (
    id CHAR(36) NOT NULL PRIMARY KEY,
    message_id CHAR(36) NOT NULL,
    message_chain_id CHAR(36),
    position INT NOT NULL,
    thought TEXT,
    tool TEXT,
    tool_labels_str TEXT NOT NULL DEFAULT '{}',
    tool_meta_str TEXT NOT NULL DEFAULT '{}',
    tool_input TEXT,
    observation TEXT,
    tool_process_data TEXT,
    message TEXT,
    message_token INT,
    message_unit_price DECIMAL(10, 4),
    message_price_unit DECIMAL(10, 7) NOT NULL DEFAULT 0.001,
    message_files TEXT,
    answer TEXT,
    answer_token INT,
    answer_unit_price DECIMAL(10, 4),
    answer_price_unit DECIMAL(10, 7) NOT NULL DEFAULT 0.001,
    tokens INT,
    total_price DECIMAL(15, 7),
    currency VARCHAR(255),
    latency FLOAT,
    created_by_role VARCHAR(255) NOT NULL,
    created_by CHAR(36) NOT NULL,
    created_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

CREATE INDEX message_agent_thought_message_id_idx ON message_agent_thoughts (message_id);
CREATE INDEX message_agent_thought_message_chain_id_idx ON message_agent_thoughts (message_chain_id);
