CREATE DATABASE if not exists `short` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `tbl_url_code`;

CREATE TABLE `tbl_url_code`
(
  `id`         int(11) unsigned                         NOT NULL AUTO_INCREMENT,
  `url`        varchar(1000) COLLATE utf8mb4_unicode_ci NOT NULL,
  `md5`        varchar(32) COLLATE utf8mb4_unicode_ci   NOT NULL,
  `code`       varchar(12) COLLATE utf8mb4_unicode_ci   NOT NULL,
  `created_at` int(11) unsigned                         NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_md5` (`md5`) USING HASH
) ENGINE = InnoDB
  AUTO_INCREMENT = 10000
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci

