create schema if not exists short collate utf8mb4_unicode_ci;

create table if not exists tbl_url_code
(
  id         int(11) unsigned auto_increment primary key,
  url        varchar(1000)    not null,
  md5        varchar(32)      not null,
  code       varchar(12)      not null,
  created_at int(11) unsigned not null
);