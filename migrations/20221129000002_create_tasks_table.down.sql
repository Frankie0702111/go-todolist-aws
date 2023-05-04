ALTER TABLE `tasks` DROP FOREIGN KEY `users_user_id_foreign`;
ALTER TABLE `tasks` DROP FOREIGN KEY `categories_category_id_foreign`;
DROP TABLE IF EXISTS `tasks`;
