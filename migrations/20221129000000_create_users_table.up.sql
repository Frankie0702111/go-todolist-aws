CREATE TABLE IF NOT EXISTS `users`
(
  `id`          bigint        NOT NULL  AUTO_INCREMENT  PRIMARY KEY,
  `username`    varchar(30)   NOT NULL  DEFAULT ''      COMMENT '用戶名',
  `email`       varchar(50)   NOT NULL  DEFAULT ''      COMMENT '信箱',
  `password`    varchar(255)  NOT NULL  DEFAULT ''      COMMENT '密碼',
  `status`      bool          NOT NULL  DEFAULT true    COMMENT '狀態',
  `created_at`  timestamp     NOT NULL  DEFAULT NOW()   COMMENT '新增時間',
  `updated_at`  timestamp     NOT NULL  DEFAULT NOW()   COMMENT '更新時間'
);

create unique index `uidx_email` on `users` (`email`) using BTREE;
create index `idx_created_at` on `users` (`created_at` desc) using BTREE;
create index `idx_updated_at` on `users` (`updated_at` desc) using BTREE;

-- Password : 12345678
INSERT INTO `users` VALUES (1, 'admin', 'admin@test.com', '$2a$04$xSVWWmqozt9ir0NY6A7yYeGc/JGG3zASnynjAMPF1YngOCvkd3QqK', 1, '2022-11-29 09:00:00', '2022-11-29 09:00:00');
