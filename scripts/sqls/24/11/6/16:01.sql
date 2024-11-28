-- Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file.

-- ----------------------------
-- Table structure for tenant_default_models
-- ----------------------------
DROP TABLE IF EXISTS `tenant_default_models`;
CREATE TABLE `tenant_default_models` (
    id CHAR(36) NOT NULL PRIMARY KEY,
    tenant_id CHAR(36) NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(40) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    created_by CHAR(36),
    updated_by CHAR(36)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;