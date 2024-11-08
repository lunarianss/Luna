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

