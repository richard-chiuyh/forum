DROP TABLE IF EXISTS `user_info`;
CREATE TABLE IF NOT EXISTS `user_info` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint NOT NULL DEFAULT 0 COMMENT '用户ID',
  `account` varchar(32) NOT NULL DEFAULT '' COMMENT '用户账号名',
  `last_seen_post` bigint NOT NULL DEFAULT 0 COMMENT '最后查看的帖子时间',
  `post_count` int NOT NULL DEFAULT 0 COMMENT '帖子数量',
  `create_time` bigint NOT NULL DEFAULT 0 COMMENT '创建时间',
  `update_time` bigint NOT NULL DEFAULT 0 COMMENT '更新时间',
  `delete_time` bigint NOT NULL DEFAULT 0 COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';

DROP TABLE IF EXISTS `post`;
CREATE TABLE IF NOT EXISTS `post` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `parent_id` bigint NOT NULL DEFAULT 0 COMMENT '父节点ID 0:帖子 帖子id:一级评论 评论id:回复',
  `user_id` bigint NOT NULL DEFAULT 0 COMMENT '用户ID',
  `account` varchar(100) NOT NULL DEFAULT '' COMMENT '用户账号名',
  `content` text NOT NULL COMMENT '内容',
  `images` json NOT NULL DEFAULT (json_array()) COMMENT '图片信息',
  `word_count` int NOT NULL DEFAULT 0 COMMENT '字数',
  `status` tinyint NOT NULL DEFAULT 30 COMMENT '状态 30:已发布 90:已删除',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `view_count` int NOT NULL DEFAULT 0 COMMENT '浏览量',
  `like_count` int NOT NULL DEFAULT 0 COMMENT '点赞量',
  `comment_count` int NOT NULL DEFAULT 0 COMMENT '评论量',
  `fav_count` int NOT NULL DEFAULT 0 COMMENT '收藏量',
  `ip_address` varchar(256) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `publish_time` bigint NOT NULL DEFAULT 0 COMMENT '发布时间',
  `status_update_time` bigint NOT NULL DEFAULT 0 COMMENT '状态更新时间',
  `create_time` bigint NOT NULL DEFAULT 0 COMMENT '创建时间',
  `update_time` bigint NOT NULL DEFAULT 0 COMMENT '更新时间',
  `delete_time` bigint NOT NULL DEFAULT 0 COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_status_id` (`parent_id`, `status`, `id`),
  KEY `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='内容表';

DROP TABLE IF EXISTS `like`;
CREATE TABLE IF NOT EXISTS `like` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `target_id` bigint NOT NULL DEFAULT 0 COMMENT '目标内容ID(帖子或评论)',
  `user_id` bigint NOT NULL DEFAULT 0 COMMENT '用户ID',
  `account` varchar(32) NOT NULL DEFAULT '' COMMENT '用户账号名',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态 30:正常 90:已删除',
  `ip_address` varchar(256) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `create_time` bigint NOT NULL DEFAULT 0 COMMENT '创建时间',
  `update_time` bigint NOT NULL DEFAULT 0 COMMENT '更新时间',
  `delete_time` bigint NOT NULL DEFAULT 0 COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_target_user_status` (`target_id`, `user_id`, `status`),
  KEY `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='点赞表';

DROP TABLE IF EXISTS `fav`;
CREATE TABLE IF NOT EXISTS `fav` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `target_id` bigint NOT NULL DEFAULT 0 COMMENT '目标内容ID(帖子或评论)',
  `user_id` bigint NOT NULL DEFAULT 0 COMMENT '用户ID',
  `account` varchar(32) NOT NULL DEFAULT '' COMMENT '用户账号名',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '状态 30:正常 90:已删除',
  `create_time` bigint NOT NULL DEFAULT 0 COMMENT '创建时间',
  `update_time` bigint NOT NULL DEFAULT 0 COMMENT '更新时间',
  `delete_time` bigint NOT NULL DEFAULT 0 COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_target_user_status` (`target_id`, `user_id`, `status`),
  KEY `idx_user_status` (`user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏表';
