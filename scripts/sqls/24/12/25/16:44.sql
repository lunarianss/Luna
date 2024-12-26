
-- ----------------------------
-- Table structure for embeddings
-- ----------------------------
DROP TABLE IF EXISTS `embeddings`;
CREATE TABLE embeddings (
    id CHAR(36) NOT NULL PRIMARY KEY,
    provider_name VARCHAR(40) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    hash VARCHAR(64) NOT NULL,
    embedding LONGBLOB NOT NULL,
    created_at int(10) NOT NULL,
    UNIQUE KEY embedding_hash_idx (model_name,hash,provider_name)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE embeddings ADD INDEX created_at_idx (create_at);