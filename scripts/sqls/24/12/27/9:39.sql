
-- ----------------------------
-- Table structure for app_annotation_hit_histories
-- ----------------------------
DROP TABLE IF EXISTS `app_annotation_hit_histories`;
CREATE TABLE app_annotation_hit_histories (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    annotation_id CHAR(36) NOT NULL,
    source VARCHAR(64) NOT NULL,
    question text NOT NULL,
    account_id CHAR(36) NOT NULL,
    created_at int(10) NOT NULL,
    score FLOAT(10, 9) NOT NULL,
    message_id CHAR(36) NOT NULL,
    annotation_question text NOT NULL,
    annotation_content text NOT NULL
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

ALTER TABLE app_annotation_hit_histories ADD INDEX app_annotation_hit_histories_app_idx (app_id);
ALTER TABLE app_annotation_hit_histories ADD INDEX app_annotation_hit_histories_account_idx (account_id);
ALTER TABLE app_annotation_hit_histories ADD INDEX app_annotation_hit_histories_annotation_idx (annotation_id);
ALTER TABLE app_annotation_hit_histories ADD INDEX app_annotation_hit_histories_message_idx (message_id);