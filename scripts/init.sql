CREATE DATABASE IF NOT EXISTS betxin;

use betxin;

CREATE TABLE
    IF NOT EXISTS `user` (
        `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
        `identity_number` varchar(36) NOT NULL DEFAULT "",
        `uid` varchar(36) NOT NULL,
        `full_name` varchar(255) DEFAULT NULL,
        `avatar_url` varchar(255) DEFAULT NULL,
        `session_id` varchar(255) DEFAULT NULL,
        `biography` varchar(255) DEFAULT NULL,
        `private_key` varchar(255) DEFAULT NULL,
        `client_id` VARCHAR(36) DEFAULT NULL,
        `contract` VARCHAR(255) DEFAULT NULL,
        `is_mvm_user` TINYINT(1) DEFAULT 0,
        `created_at` BIGINT(13) NOT NULL DEFAULT 0,
        `updated_at` BIGINT(13) NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_identity_number` (`identity_number`),
        UNIQUE KEY `idx_uid` (`uid`),
        KEY `idx_full_name` (`full_name`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `topic` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '话题自增ID',
        `tid` BIGINT(20) NOT NULL,
        `cid` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类ID',
        `title` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '标题',
        `intro` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '概述',
        `content` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '内容',
        `yes_ratio` VARCHAR(10) NOT NULL DEFAULT '50.00' COMMENT '赞成率',
        `no_ratio` VARCHAR(10) NOT NULL DEFAULT '50.00' COMMENT '反对率',
        `yes_count` VARCHAR(40) NOT NULL DEFAULT '0.00000000' COMMENT '赞成计数',
        `no_count` VARCHAR(40) NOT NULL DEFAULT '0.00000000' COMMENT '反对计数',
        `total_count` VARCHAR(41) NOT NULL DEFAULT '0.00000000' COMMENT '总计数',
        `collect_count` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '收藏数',
        `read_count` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '阅读数',
        `img_url` TEXT NOT NULL COMMENT '图片URL',
        `is_stop` TINYINT(1) DEFAULT 0 COMMENT '是否结束',
        `refund_end_time` BIGINT NOT NULL DEFAULT 0 COMMENT '退款截止时间',
        `end_time` BIGINT NOT NULL DEFAULT 0 COMMENT '话题结束时间',
        `is_deleted` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否删除',
        `created_at` BIGINT(13) NOT NULL DEFAULT 0 COMMENT '创建时间',
        `updated_at` BIGINT(13) NOT NULL DEFAULT 0 COMMENT '更新时间',
        `deleted_at` BIGINT(13) DEFAULT NULL COMMENT '删除时间',
        PRIMARY KEY (`id`),
        INDEX idx_cid (cid) COMMENT '分类索引',
        UNIQUE idx_tid (tid) COMMENT '唯一索引',
        INDEX title_intro_content_index (title, intro, content) COMMENT '全文索引'
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 种类系统

CREATE TABLE
    IF NOT EXISTS `category` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `category_name` VARCHAR(20) NOT NULL DEFAULT '',
        PRIMARY KEY (`id`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

INSERT INTO category (id, category_name) VALUES (1, 'Buisiness');

INSERT INTO category (id, category_name) VALUES (2, 'Crypto');

INSERT INTO category (id, category_name) VALUES (3, 'Sports');

INSERT INTO category (id, category_name) VALUES (4, 'Politics');

INSERT INTO category (id, category_name) VALUES (5, 'New');

INSERT INTO category (id, category_name) VALUES (6, 'Trending');

-- 转账信息系统

CREATE TABLE
    IF NOT EXISTS `snapshot`(
        `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
        `trace_id` VARCHAR(36) NOT NULL DEFAULT '',
        `memo` VARCHAR(256) NOT NULL DEFAULT '',
        `type` VARCHAR(10) NOT NULL DEFAULT '',
        `snapshot_id` VARCHAR(36) NOT NULL DEFAULT '',
        `opponent_id` VARCHAR(36) NOT NULL DEFAULT '',
        `asset_id` VARCHAR(36) NOT NULL DEFAULT '',
        `amount` VARCHAR(36) NOT NULL DEFAULT '',
        `opening_balance` VARCHAR(40) NOT NULL DEFAULT '',
        `closing_balance` VARCHAR(40) NOT NULL DEFAULT '',
        `created_at` VARCHAR(50) NOT NULL DEFAULT '',
        PRIMARY KEY (`id`),
        INDEX `idx_trace_id` (`trace_id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 话题购买系统

CREATE TABLE
    IF NOT EXISTS `topic_purchases`(
        `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
        `trace_id` VARCHAR(36) NOT NULL DEFAULT '' COMMENT '话题购买的trace_id',
        `uid` varchar(36) NOT NULL DEFAULT '',
        `tid` BIGINT(20) NOT NULL,
        `yes_price` VARCHAR(40) NOT NULL DEFAULT '0.00000000' COMMENT '支持金额',
        `no_price` VARCHAR(40) NOT NULL DEFAULT '0.00000000' COMMENT '反对金额',
        `created_at` BIGINT(13) NOT NULL DEFAULT 0,
        `updated_at` BIGINT(13) NOT NULL DEFAULT 0,
        `deleted_at` BIGINT(13) DEFAULT NULL,
        PRIMARY KEY (id),
        UNIQUE KEY idx_uid_tid (uid, tid),
        KEY idx_tid (tid),
        KEY idx_yes_price (yes_price),
        KEY idx_no_price (no_price)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 退款系统

CREATE TABLE
    IF NOT EXISTS `refund`(
        `id` int NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL DEFAULT '',
        `asset_id` VARCHAR(36) NOT NULL DEFAULT '',
        `trace_id` VARCHAR(36) NOT NULL DEFAULT '',
        `price` VARCHAR(40) NOT NULL DEFAULT '0.00000000' COMMENT '退款金额',
        `select` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '选择',
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
        `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态',
        `created_at` BIGINT(13) NOT NULL DEFAULT 0,
        `updated_at` BIGINT(13) NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_uid_tid` (`uid`, `tid`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `feedback`(
        `id` int NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL DEFAULT '',
        `title` varchar(150) NOT NULL DEFAULT '',
        `content` TEXT NOT NULL,
        `created_at` BIGINT(13) NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `message`(
        `id` int NOT NULL AUTO_INCREMENT,
        `uid` varchar(36) NOT NULL DEFAULT '',
        `data` LONGTEXT NOT NULL,
        `conversation_id` VARCHAR(36) NOT NULL DEFAULT '',
        `recipient_id` VARCHAR(36) NOT NULL DEFAULT '',
        `message_id` VARCHAR(36) NOT NULL DEFAULT '',
        `category` VARCHAR(20) NOT NULL DEFAULT '',
        `created_at` BIGINT(13) NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;