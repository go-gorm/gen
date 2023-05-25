CREATE TABLE `banks` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext,
  `address` longtext,
  `scale` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `credit_cards` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `number` longtext,
  `customer_refer` bigint(20) unsigned DEFAULT NULL,
  `bank_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_credit_cards_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `customers` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `bank_id` bigint(20) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_customers_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `people` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `alias` varchar(255) DEFAULT NULL,
  `age` int(11) unsigned DEFAULT NULL,
  `flag` tinyint(1) DEFAULT NULL,
  `another_flag` tinyint(4) DEFAULT NULL,
  `commit` varchar(255) DEFAULT NULL,
  `First` tinyint(1) DEFAULT NULL,
  `bit` bit(1) DEFAULT NULL,
  `small` smallint(5) unsigned DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `score` decimal(19,0) DEFAULT NULL,
  `number` int(11) DEFAULT NULL,
  `birth` datetime DEFAULT CURRENT_TIMESTAMP,
  `xmlHTTPRequest` varchar(255) DEFAULT ' ',
  `jStr` json DEFAULT NULL,
  `geo` json DEFAULT NULL,
  `mint` mediumint(9) DEFAULT NULL,
  `blank` varchar(64) DEFAULT ' ',
  `remark` text,
  `long_remark` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
