CREATE TABLE `users` (
                         `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
                         `created_at` datetime(3) DEFAULT NULL,
                         `name` varchar(255) DEFAULT NULL COMMENT 'oneline',
                         `address` varchar(255) DEFAULT '',
                         `register_time` datetime(3) DEFAULT NULL,
                         `alive` tinyint(1) DEFAULT NULL COMMENT 'multiline\nline1\nline2',
                         `company_id` bigint(20) unsigned DEFAULT '666',
                         `private_url` varchar(255) DEFAULT 'https://a.b.c ',
                         PRIMARY KEY (`id`),
                         KEY `idx_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
