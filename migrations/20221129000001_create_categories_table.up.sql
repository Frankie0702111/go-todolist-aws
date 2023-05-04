CREATE TABLE IF NOT EXISTS `categories` (
  `id`          bigint        NOT NULL  AUTO_INCREMENT  PRIMARY KEY,
  `name`        varchar(100)  NOT NULL  DEFAULT ''      COMMENT '類別名稱',
  `created_at`  timestamp     NOT NULL  DEFAULT NOW()   COMMENT '新增時間',
  `updated_at`  timestamp     NOT NULL  DEFAULT NOW()   COMMENT '更新時間'
);

create unique index `uidx_name` on `categories` (`name`) using BTREE;

INSERT INTO `categories` (`name`, `created_at`, `updated_at`)
VALUES
("提醒", '2022-11-29 09:00:00', '2022-11-29 09:00:00'),
("工作", '2022-11-29 09:00:00', '2022-11-29 09:00:00'),
("活動", '2022-11-29 09:00:00', '2022-11-29 09:00:00');
