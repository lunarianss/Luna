-- Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file.


DROP TABLE IF EXISTS `app_model_configs`;
CREATE TABLE `app_model_configs` (
    id CHAR(36) NOT NULL PRIMARY KEY,
    app_id CHAR(36) NOT NULL,
    provider VARCHAR(255),
    model_id VARCHAR(255),
    configs JSON,
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
