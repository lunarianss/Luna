-- ----------------------------
-- Table structure for document_segments
-- ----------------------------
DROP TABLE IF EXISTS `document_segments`;
CREATE TABLE document_segments (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    dataset_id CHAR(36) NOT NULL,
    document_id CHAR(36) NOT NULL,
    position INT NOT NULL,
    content TEXT NOT NULL,
    answer TEXT,
    word_count INT NOT NULL,
    tokens INT NOT NULL,
    keywords JSON,
    index_node_id VARCHAR(255),
    index_node_hash VARCHAR(255),
    hit_count INT NOT NULL DEFAULT 0,
    enabled bit(1) NOT NULL DEFAULT 1,
    disabled_at INT(10) DEFAULT NULL,
    disabled_by CHAR(36) DEFAULT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'waiting',
    created_by CHAR(36) NOT NULL,
    created_at int(10) NOT NULL,
    updated_by CHAR(36) DEFAULT NULL,
    updated_at int(10) NOT NULL,
    indexing_at int(10) DEFAULT NULL,
    completed_at int(10) DEFAULT NULL,
    error TEXT,
    stopped_at int(10) DEFAULT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Indexes for document_segments
-- ----------------------------
ALTER TABLE document_segments ADD INDEX document_segment_dataset_id_idx (dataset_id);
ALTER TABLE document_segments ADD INDEX document_segment_document_id_idx (document_id);
ALTER TABLE document_segments ADD INDEX document_segment_tenant_dataset_idx (dataset_id, tenant_id);
ALTER TABLE document_segments ADD INDEX document_segment_tenant_document_idx (document_id, tenant_id);
ALTER TABLE document_segments ADD INDEX document_segment_dataset_node_idx (dataset_id, index_node_id);
ALTER TABLE document_segments ADD INDEX document_segment_tenant_idx (tenant_id);




-- ----------------------------
-- Table structure for documents
-- ----------------------------
DROP TABLE IF EXISTS `documents`;
CREATE TABLE documents (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    dataset_id CHAR(36) NOT NULL,
    position INT NOT NULL,
    data_source_type VARCHAR(255) NOT NULL,
    data_source_info TEXT,
    dataset_process_rule_id CHAR(36),
    batch VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_from VARCHAR(255) NOT NULL,
    created_by CHAR(36) NOT NULL,
    created_api_request_id CHAR(36),
    created_at int(10) NOT NULL,
    processing_started_at int(10),
    file_id TEXT,
    word_count INT,
    parsing_completed_at int(10),
    cleaning_completed_at int(10),
    splitting_completed_at int(10),
    tokens INT,
    indexing_latency FLOAT,
    completed_at int(10),
    is_paused bit(1) DEFAULT 0,
    paused_by CHAR(36),
    paused_at int(10),
    error TEXT,
    stopped_at int(10),
    indexing_status VARCHAR(255) NOT NULL DEFAULT 'waiting',
    enabled bit(1) NOT NULL DEFAULT 1,
    disabled_at int(10),
    disabled_by CHAR(36),
    archived bit(1) NOT NULL DEFAULT 0,
    archived_reason VARCHAR(255),
    archived_by CHAR(36),
    archived_at int(10),
    updated_at int(10) NOT NULL,
    doc_type VARCHAR(40),
    doc_metadata JSON,
    doc_form VARCHAR(255) NOT NULL DEFAULT 'text_model',
    doc_language VARCHAR(255)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

CREATE INDEX document_dataset_id_idx ON documents (dataset_id);
CREATE INDEX document_is_paused_idx ON documents (is_paused);
CREATE INDEX document_tenant_idx ON documents (tenant_id);




-- ----------------------------
-- Table structure for dataset_permissions
-- ----------------------------
DROP TABLE IF EXISTS `dataset_permissions`;
CREATE TABLE dataset_permissions (
    id CHAR(36) NOT NULL PRIMARY KEY,
    dataset_id CHAR(36) NOT NULL,
    account_id CHAR(36) NOT NULL,
    tenant_id CHAR(36) NOT NULL,
    has_permission bit(1) NOT NULL DEFAULT 1,
    created_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE INDEX idx_dataset_permissions_dataset_id ON dataset_permissions (dataset_id);
CREATE INDEX idx_dataset_permissions_account_id ON dataset_permissions (account_id);
CREATE INDEX idx_dataset_permissions_tenant_id ON dataset_permissions (tenant_id);


-- ----------------------------
-- Table structure for dataset_process_rules
-- ----------------------------
DROP TABLE IF EXISTS `dataset_process_rules`;
CREATE TABLE dataset_process_rules (
    id CHAR(36) NOT NULL PRIMARY KEY,
    dataset_id CHAR(36) NOT NULL,
    mode VARCHAR(255) NOT NULL DEFAULT 'automatic',
    rules TEXT,
    created_by CHAR(36) NOT NULL,
    created_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE INDEX dataset_process_rule_dataset_id_idx ON dataset_process_rules (dataset_id);



-- ----------------------------
-- Table structure for app_dataset_joins
-- ----------------------------
DROP TABLE IF EXISTS `app_dataset_joins`;
CREATE TABLE app_dataset_joins (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    dataset_id CHAR(36) NOT NULL,
    created_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE INDEX app_dataset_join_app_dataset_idx ON app_dataset_joins (dataset_id, app_id);


-- ----------------------------
-- Table structure for dataset_queries
-- ----------------------------
DROP TABLE IF EXISTS `dataset_queries`;
CREATE TABLE dataset_queries (
    id CHAR(36) NOT NULL PRIMARY KEY,
    dataset_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    source VARCHAR(255) NOT NULL,
    source_app_id CHAR(36),
    created_by_role VARCHAR(255) NOT NULL,
    created_by CHAR(36) NOT NULL,
    created_at int(10) NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE INDEX dataset_query_dataset_id_idx ON dataset_queries (dataset_id);



-- ----------------------------
-- Table structure for dataset_keyword_tables
-- ----------------------------
DROP TABLE IF EXISTS `dataset_keyword_tables`;
CREATE TABLE dataset_keyword_tables (
    id CHAR(36) NOT NULL PRIMARY KEY,
    dataset_id CHAR(36) NOT NULL UNIQUE,
    keyword_table TEXT NOT NULL,
    data_source_type VARCHAR(255) NOT NULL DEFAULT 'database'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE INDEX dataset_keyword_table_dataset_id_idx ON dataset_keyword_tables (dataset_id);




-- ----------------------------
-- Table structure for dataset_keyword_tables
-- ----------------------------
DROP TABLE IF EXISTS `dataset_keyword_tables`;
CREATE TABLE dataset_keyword_tables (
    id CHAR(36) NOT NULL PRIMARY KEY,
    dataset_id CHAR(36) NOT NULL UNIQUE,
    keyword_table TEXT NOT NULL,
    data_source_type VARCHAR(255) NOT NULL DEFAULT 'database'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE INDEX dataset_keyword_table_dataset_id_idx ON dataset_keyword_tables (dataset_id);

