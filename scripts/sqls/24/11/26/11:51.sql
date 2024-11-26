
-- ----------------------------
-- Table structure for pinned_conversations
-- ----------------------------
DROP TABLE IF EXISTS `pinned_conversations`;
CREATE TABLE pinned_conversations (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    conversation_id CHAR(36) NOT NULL,
    created_by_role VARCHAR(255) NOT NULL DEFAULT 'end_user',
    created_by CHAR(36) NOT NULL,
    created_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE pinned_conversations ADD INDEX pinned_conversation_conversation_idx (app_id, conversation_id, created_by_role, created_by);