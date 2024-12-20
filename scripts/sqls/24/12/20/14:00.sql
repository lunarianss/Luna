-- ----------------------------
-- Table structure for message_annotations
-- ----------------------------
DROP TABLE IF EXISTS `message_annotations`;
CREATE TABLE message_annotations (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    conversation_id CHAR(36),
    message_id CHAR(36),
    question TEXT,
    content TEXT NOT NULL,
    hit_count INT NOT NULL DEFAULT 0,
    account_id CHAR(36) NOT NULL,
    created_at int(10) NOT NULL,
    updated_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE message_annotations ADD INDEX message_annotation_app_idx (app_id);
ALTER TABLE message_annotations ADD INDEX message_annotation_conversation_idx (conversation_id);
ALTER TABLE message_annotations ADD INDEX message_annotation_message_idx (message_id);


-- ----------------------------
-- Table structure for app_annotation_settings
-- ----------------------------
DROP TABLE IF EXISTS `app_annotation_settings`;
CREATE TABLE app_annotation_settings (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    score_threshold FLOAT(2,1) NOT NULL DEFAULT 0,
    collection_binding_id CHAR(36) NOT NULL,
    created_user_id CHAR(36) NOT NULL,
    created_at int(10) NOT NULL,
    updated_user_id CHAR(36) NOT NULL,
    updated_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE app_annotation_settings ADD INDEX app_annotation_settings_app_idx (app_id);


-- ----------------------------
-- Table structure for dataset_collection_bindings
-- ----------------------------
DROP TABLE IF EXISTS `dataset_collection_bindings`;
CREATE TABLE dataset_collection_bindings (
    id CHAR(36) NOT NULL PRIMARY KEY,
    provider_name VARCHAR(40) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    type VARCHAR(40) NOT NULL DEFAULT 'dataset',
    collection_name VARCHAR(64) NOT NULL,
    created_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE dataset_collection_bindings ADD INDEX provider_model_name_idx (provider_name, model_name);



-- ----------------------------
-- Table structure for datasets
-- ----------------------------
DROP TABLE IF EXISTS `datasets`;
CREATE TABLE datasets (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    provider VARCHAR(255) NOT NULL DEFAULT 'vendor',
    permission VARCHAR(255) NOT NULL DEFAULT 'only_me',
    data_source_type VARCHAR(255),
    indexing_technique VARCHAR(255),
    index_struct TEXT,
    created_by CHAR(36) NOT NULL,
    created_at int(10) NOT NULL,
    updated_by CHAR(36),
    updated_at int(10) NOT NULL,
    embedding_model VARCHAR(255),
    embedding_model_provider VARCHAR(255),
    collection_binding_id CHAR(36),
    retrieval_model JSON
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE datasets ADD INDEX dataset_tenant_idx (tenant_id);
