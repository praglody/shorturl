CREATE DATABASE if not exists `short` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `url_codes`;

CREATE TABLE `url_codes`
(
    `id`         int(11) unsigned                         NOT NULL AUTO_INCREMENT,
    `url`        varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '长url',
    `md5`        varchar(32) COLLATE utf8mb4_unicode_ci   NOT NULL COMMENT '长url的md5值',
    `code`       varchar(12) COLLATE utf8mb4_bin   NOT NULL COMMENT '生成的短码',
    `click`      int(11) unsigned                         NOT NULL DEFAULT '0' COMMENT '点击数',
    `user_id`    int(11) unsigned                         NOT NULL DEFAULT '0' COMMENT '用户id',
    `created_at` int(10) unsigned                         NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_md5` (`md5`) USING HASH
) ENGINE = InnoDB
  AUTO_INCREMENT = 20000000
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users`
(
    `id`         int(11) unsigned                       NOT NULL AUTO_INCREMENT,
    `name`       varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '应用名',
    `app_id`     varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
    `app_secret` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户秘钥',
    `status`     tinyint(1)                             NOT NULL DEFAULT '0' COMMENT '是否有效，0有效，1无效，默认0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_name` (`name`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户授权表'



