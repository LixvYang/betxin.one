CREATE DATABASE IF NOT EXISTS betxin;

use betxin;

CREATE TABLE
    IF NOT EXISTS `user` (
        `id` int NOT NULL AUTO_INCREMENT,
        `identity_number` varchar(36) NOT NULL,
        `uid` varchar(36) NOT NULL,
        `full_name` varchar(255) DEFAULT NULL,
        `avatar_url` varchar(255) DEFAULT NULL,
        `session_id` varchar(255) DEFAULT NULL,
        `biography` varchar(255) DEFAULT NULL,
        `private_key` varchar(255) DEFAULT NULL,
        `client_id` VARCHAR(36) DEFAULT NULL,
        `contract` VARCHAR(255) DEFAULT NULL,
        `is_mvm_user` TINYINT(1) DEFAULT 0,
        `created_at` BIGINT NOT NULL DEFAULT 0,
        `updated_at` BIGINT NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`),
        UNIQUE KEY `idx_identity_number` (`identity_number`),
        UNIQUE KEY `idx_uid` (`uid`),
        KEY `idx_full_name` (`full_name`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS topic (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '话题自增ID',
        `tid` VARCHAR(36) NOT NULL DEFAULT '' COMMENT '话题唯一标识',
        `cid` BIGINT(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '分类ID',
        `title` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '标题',
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
        `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间',
        `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间',
        `deleted_at` BIGINT(20) DEFAULT NULL COMMENT '删除时间',
        PRIMARY KEY (`id`),
        INDEX idx_cid (cid) COMMENT '分类索引',
        UNIQUE idx_tid (tid) COMMENT '唯一索引',
        INDEX title_intro_content_index (title, intro, content) COMMENT '全文索引'
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE
    IF NOT EXISTS `category` (
        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
        `category_name` VARCHAR(20) NOT NULL DEFAULT '',
        PRIMARY KEY (`id`)
    ) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;