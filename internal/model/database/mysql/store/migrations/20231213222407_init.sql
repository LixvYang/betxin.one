-- +goose Up

CREATE TABLE
    `user` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL DEFAULT '',
        `identity_number` VARCHAR(36) DEFAULT '',
        `full_name` varchar(255) DEFAULT '',
        `avatar_url` varchar(255) DEFAULT '',
        `session_id` varchar(255) DEFAULT '',
        `biography` varchar(255) DEFAULT '',
        `private_key` varchar(255) DEFAULT '',
        `client_id` VARCHAR(36) DEFAULT '',
        `contract` VARCHAR(255) DEFAULT '',
        `is_mvm_user` TINYINT(1) DEFAULT 0,
        `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_identity_number` (`identity_number`),
        UNIQUE KEY `idx_uid` (`uid`),
        KEY `idx_full_name` (`full_name`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `topic` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `tid` VARCHAR(36) NOT NULL DEFAULT '',
        `cid` BIGINT(5) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'categyro id',
        `title` VARCHAR(20) NOT NULL DEFAULT '',
        `intro` VARCHAR(255) NOT NULL DEFAULT '',
        `content` VARCHAR(255) NOT NULL DEFAULT '',
        `yes_ratio` VARCHAR(10) NOT NULL DEFAULT '50.00',
        `no_ratio` VARCHAR(10) NOT NULL DEFAULT '50.00',
        `yes_count` VARCHAR(40) NOT NULL DEFAULT '0.00000000',
        `no_count` VARCHAR(40) NOT NULL DEFAULT '0.00000000',
        `total_count` VARCHAR(41) NOT NULL DEFAULT '0.00000000',
        `collect_count` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
        `read_count` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0,
        `img_url` TEXT NOT NULL COMMENT 'image url',
        `is_stop` TINYINT(1) DEFAULT 0,
        `refund_end_time` TIMESTAMP DEFAULT NULL,
        `end_time` TIMESTAMP DEFAULT NULL,
        `is_deleted` TINYINT(1) NOT NULL DEFAULT 0,
        `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP DEFAULT NULL,
        PRIMARY KEY (`id`),
        INDEX `idx_cid` (cid),
        UNIQUE `idx_tid` (tid),
        INDEX `title_intro_content_index` (title, intro, content)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
-- 种类系统

CREATE TABLE
    IF NOT EXISTS `category` (
        `id` BIGINT(5) UNSIGNED NOT NULL AUTO_INCREMENT,
        `name` VARCHAR(20) NOT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

INSERT INTO category (id, name) VALUES (1, 'Buisiness');

INSERT INTO category (id, name) VALUES (2, 'Crypto');

INSERT INTO category (id, name) VALUES (3, 'Sports');

INSERT INTO category (id, name) VALUES (4, 'Games');

INSERT INTO category (id, name) VALUES (5, 'News');

INSERT INTO category (id, name) VALUES (6, 'Trending');

INSERT INTO category (id, name) VALUES (7, 'Others');

-- 转账信息系统

CREATE TABLE
    IF NOT EXISTS `snapshot`(
        `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
        `trace_id` VARCHAR(36) NOT NULL,
        `memo` VARCHAR(256) NOT NULL,
        `type` VARCHAR(10) NOT NULL,
        `snapshot_id` VARCHAR(36) NOT NULL,
        `opponent_id` VARCHAR(36) NOT NULL,
        `asset_id` VARCHAR(36) NOT NULL,
        `amount` VARCHAR(36) NOT NULL,
        `opening_balance` VARCHAR(40) NOT NULL,
        `closing_balance` VARCHAR(40) NOT NULL,
        `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        INDEX `idx_trace_id` (`trace_id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 话题购买系统

CREATE TABLE
    IF NOT EXISTS `topic_purchases`(
        `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL DEFAULT '',
        `tid` BIGINT(20) NOT NULL,
        `yes_price` VARCHAR(40) NOT NULL DEFAULT '0.00000000',
        `no_price` VARCHAR(40) NOT NULL DEFAULT '0.00000000',
        `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY idx_uid_tid (uid, tid),
        KEY idx_tid (tid),
        KEY idx_yes_price (yes_price),
        KEY idx_no_price (no_price)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `bonuse` (
        `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
        `uid` VARCHAR(36) NOT NULL DEFAULT '',
        `tid` BIGINT(20) NOT NULL,
        `asset_id` VARCHAR(36) NOT NULL,
        `amount` VARCHAR(36) NOT NULL,
        `memo` VARCHAR(256) NOT NULL,
        `trace_id` VARCHAR(36) NOT NULL,
        `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP DEFAULT NULL,
        PRIMARY KEY (id)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 退款系统

CREATE TABLE
    IF NOT EXISTS `refund`(
        `id` int NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL,
        `asset_id` VARCHAR(36) NOT NULL,
        `trace_id` VARCHAR(36) NOT NULL,
        `price` VARCHAR(40) NOT NULL DEFAULT '0.00000000',
        `select` TINYINT(1) NOT NULL DEFAULT 0,
        `memo` VARCHAR(256) NOT NULL DEFAULT '',
        `created_at` BIGINT(13) NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_trace_id` (`trace_id`),
        UNIQUE KEY `idx_uid` (`uid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

--collect

CREATE TABLE
    IF NOT EXISTS `collect`(
        `id` int NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL DEFAULT '',
        `tid` BIGINT(20) NOT NULL,
        `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT 'status',
        `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_uid_tid` (`uid`, `tid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `feedback`(
        `id` int NOT NULL AUTO_INCREMENT,
        `fid` VARCHAR(36) NOT NULL DEFAULT '',
        `uid` varchar(36) NOT NULL DEFAULT '',
        `title` varchar(150) NOT NULL DEFAULT '',
        `content` VARCHAR(512) NOT NULL,
        `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP DEFAULT NULL,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `message`(
        `id` int NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL DEFAULT '',
        `data` VARCHAR(512) NOT NULL,
        `conversation_id` VARCHAR(36) NOT NULL DEFAULT '',
        `recipient_id` VARCHAR(36) NOT NULL DEFAULT '',
        `message_id` VARCHAR(36) NOT NULL DEFAULT '',
        `category` VARCHAR(20) NOT NULL DEFAULT '',
        `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- +goose Down

SELECT 'down SQL query';

DROP TABLE IF EXISTS `user`;

DROP TABLE IF EXISTS `bonuse`;

DROP TABLE IF EXISTS `category`;

DROP TABLE IF EXISTS `collect`;

DROP TABLE IF EXISTS feedback;

DROP TABLE IF EXISTS message;

DROP TABLE IF EXISTS refund;

DROP TABLE IF EXISTS snapshot;

DROP TABLE IF EXISTS topic;

DROP TABLE IF EXISTS topic_purchases;