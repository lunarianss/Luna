-- Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file.


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



-- ----------------------------
-- Table structure for api_tokens
-- ----------------------------
DROP TABLE IF EXISTS `api_tokens`;
CREATE TABLE api_tokens (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    type VARCHAR(16) NOT NULL,
    token VARCHAR(255) NOT NULL,
    last_used_at int(10),
    created_at int(10) NOT NULL,
    UNIQUE (token)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE api_tokens ADD INDEX api_token_app_id_type_idx (app_id, type);
ALTER TABLE api_tokens ADD INDEX api_token_token_idx (token, type);
ALTER TABLE api_tokens ADD INDEX api_token_tenant_idx (tenant_id, type);