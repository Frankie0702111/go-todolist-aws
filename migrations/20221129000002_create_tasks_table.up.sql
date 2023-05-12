CREATE TABLE IF NOT EXISTS `tasks` (
  `id`                bigint        NOT NULL  AUTO_INCREMENT  PRIMARY KEY,
  `user_id`           bigint        NOT NULL,
  `category_id`       bigint        NOT NULL,
  `title`             varchar(100)  NOT NULL  DEFAULT ''      COMMENT '標題',
  `note`              text          NULL                      COMMENT '備註',
  `url`               text          NULL                      COMMENT '網址',
  `img`               varchar(100)  NULL                      COMMENT '影像名稱',
  `img_link`          varchar(255)  NULL                      COMMENT '影像路徑',
  `img_uuid`          varchar(128)  NULL      DEFAULT NULL    COMMENT '影像唯一碼',
  `specify_datetime`  datetime      NULL      DEFAULT NULL    COMMENT '指定日期時間(Y-m-d H:i:s)',
  `is_specify_time`   bool          NULL      DEFAULT false   COMMENT '是否指定時間',
  `priority`          tinyint       NOT NULL  DEFAULT 1       COMMENT '1:低 2:中 3:高',
  `is_complete`       bool          NOT NULL  DEFAULT false   COMMENT '是否完成',
  `is_notify`         tinyint       NOT NULL  DEFAULT 0       COMMENT '0:未通知, 1:已通知',
  `created_at`        timestamp     NOT NULL  DEFAULT NOW()   COMMENT '新增時間',
  `updated_at`        timestamp     NOT NULL  DEFAULT NOW()   COMMENT '更新時間'
);

create unique index `unique_title` on `tasks` (`title`) using BTREE;
create index `idx_user_id_title` on `tasks` (`user_id`, `title`(15)) using BTREE;
create index `idx_created_at` on `tasks` (`created_at` desc) using BTREE;
create index `idx_updated_at` on `tasks` (`updated_at` desc) using BTREE;
ALTER TABLE `tasks` ADD CONSTRAINT `users_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE;
ALTER TABLE `tasks` ADD CONSTRAINT `categories_category_id_foreign` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE CASCADE;

INSERT INTO `tasks` (`user_id`, `category_id`, `title`, `note`, `url`, `img`, `img_link`, `img_uuid`, `specify_datetime`, `is_specify_time`, `priority`, `is_complete`, `is_notify`, `created_at`, `updated_at`)
VALUES
(1, 2, 'Read book', '10 pag to 15 pag', NULL, NULL, NULL, NULL, '2023-03-23 07:13:00', 1, 1, 0, 0, '2022-11-29 09:00:00', '2022-11-29 09:00:00'),
(1, 1, 'Go to super market', 'apple', NULL, NULL, NULL, NULL, '2023-03-23 00:00:00', 0, 2, 0, 0, '2022-11-29 09:00:00', '2022-11-29 09:00:00'),
(1, 3, 'dating with friends', 'Frankie, Daisy', 'https://goo.gl/maps/VcDMEzRLKpqPtm697', NULL, NULL, 'a44c445c-0060-4c8a-8409-17e9479f91df', '2022-11-29 14:00:00', 1, 3, 1, 1, '2022-11-29 09:00:00', '2022-11-29 09:00:00');
