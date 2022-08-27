package tests_test

var tableMetas = [][2]string{
	{
		"DROP TABLE IF EXISTS `demo`;",
		"CREATE TABLE `demo` (" +
			" `id` int(11) NOT NULL AUTO_INCREMENT, " +
			"PRIMARY KEY (`id`) " +
			") ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;",
	},
	{
		"DROP TABLE IF EXISTS `banks`;",
		"CREATE TABLE `banks` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`name` longtext," +
			"`address` longtext," +
			"`scale` bigint(20) DEFAULT NULL," +
			"PRIMARY KEY (`id`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	},
	{
		"DROP TABLE IF EXISTS `credit_cards`;",
		"CREATE TABLE `credit_cards` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`created_at` datetime(3) DEFAULT NULL," +
			"`updated_at` datetime(3) DEFAULT NULL," +
			"`deleted_at` datetime(3) DEFAULT NULL," +
			"`number` longtext," +
			"`customer_refer` bigint(20) unsigned DEFAULT NULL," +
			"`bank_id` bigint(20) unsigned DEFAULT NULL," +
			"PRIMARY KEY (`id`)," +
			"KEY `idx_credit_cards_deleted_at` (`deleted_at`)," +
			"KEY `fk_credit_cards_bank` (`bank_id`)," +
			"KEY `fk_customers_credit_cards` (`customer_refer`)," +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	},
	{
		"DROP TABLE IF EXISTS `customers`;",
		"CREATE TABLE `customers` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`created_at` datetime(3) DEFAULT NULL," +
			"`updated_at` datetime(3) DEFAULT NULL," +
			"`deleted_at` datetime(3) DEFAULT NULL," +
			"`bank_id` bigint(20) unsigned DEFAULT NULL," +
			"PRIMARY KEY (`id`)," +
			"KEY `idx_customers_deleted_at` (`deleted_at`)," +
			"KEY `fk_banks_c` (`bank_id`)," +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	},
	{
		"DROP TABLE IF EXISTS `languages`;",
		"CREATE TABLE `languages` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`created_at` datetime(3) DEFAULT NULL," +
			"`updated_at` datetime(3) DEFAULT NULL," +
			"`deleted_at` datetime(3) DEFAULT NULL," +
			"`name` longtext," +
			"PRIMARY KEY (`id`)," +
			"KEY `idx_languages_deleted_at` (`deleted_at`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	},
	{
		"DROP TABLE IF EXISTS `people`;",
		"CREATE TABLE `people` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`name` varchar(255) DEFAULT NULL," +
			"`age` int(11) unsigned DEFAULT NULL," +
			"`flag` tinyint(1) unsigned DEFAULT NULL," +
			"`commit` varchar(255) DEFAULT NULL," +
			"`First` tinyint(1) DEFAULT NULL," +
			"`flag_another` tinyint(4) DEFAULT NULL," +
			"`bit` bit(1) DEFAULT NULL," +
			"`small` smallint(5) unsigned DEFAULT NULL," +
			"`deleted_at` datetime(3) DEFAULT NULL," +
			"`score` decimal(19,0) DEFAULT NULL," +
			"`type` int(11) DEFAULT NULL," +
			"`birth` datetime DEFAULT CURRENT_TIMESTAMP," +
			"PRIMARY KEY (`id`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	},
	{
		"DROP TABLE IF EXISTS `players`;",
		"CREATE TABLE `players` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`created_at` datetime(3) DEFAULT NULL," +
			"`updated_at` datetime(3) DEFAULT NULL," +
			"`deleted_at` datetime(3) DEFAULT NULL," +
			"PRIMARY KEY (`id`)," +
			"KEY `idx_players_deleted_at` (`deleted_at`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	},
	{
		"DROP TABLE IF EXISTS `user`;",
		"CREATE TABLE `user` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`name` varchar(255) DEFAULT NULL COMMENT 'oneline'," +
			"`address` varchar(255) DEFAULT ''," +
			"`register_time` time DEFAULT NULL," +
			"`alive` tinyint(1) DEFAULT NULL COMMENT 'multiline\nline1\nline2'," +
			"`created_at` time DEFAULT NULL," +
			"`company_id` bigint(20) unsigned DEFAULT '666'," +
			"`private_url` varchar(255) DEFAULT 'https://a.b.c '," +
			"`xmlHTTPRequest` varchar(255) DEFAULT ' '," +
			"`jStr` json DEFAULT NULL," +
			"`geo` geometry DEFAULT NULL," +
			"`mint` mediumint(9) DEFAULT NULL," +
			"`blank` varchar(64) DEFAULT ' '," +
			"PRIMARY KEY (`id`)," +
			"KEY `idx_name` (`name`) USING BTREE" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	},
	{
		"DROP TABLE IF EXISTS `users`;",
		"CREATE TABLE `users` (" +
			"`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
			"`name` longtext," +
			"`age` varchar(64) DEFAULT NULL," +
			"`address` varchar(255) DEFAULT NULL," +
			"`role` varchar(64) DEFAULT NULL," +
			"`created_at` datetime(3) DEFAULT NULL," +
			"`updated_at` datetime(3) DEFAULT NULL," +
			"`deleted_at` datetime(3) DEFAULT NULL," +
			"`remark` text," +
			"PRIMARY KEY (`id`)," +
			"KEY `idx_users_deleted_at` (`deleted_at`)," +
			"KEY `idx_age` (`age`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	},
}

func GetDDL() [][2]string {
	return tableMetas
}
