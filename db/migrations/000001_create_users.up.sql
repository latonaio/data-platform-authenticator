DROP TABLE IF EXISTS `Users`;

CREATE TABLE `Users` (
                         `id` int AUTO_INCREMENT,
                         `login_id` varchar(20) NOT NULL DEFAULT '',
                         `password` varchar(70) NOT NULL DEFAULT '',
                         `qos` varchar(20) NOT NULL DEFAULT 'default',
                         `is_encrypt` BOOLEAN NOT NULL DEFAULT TRUE,
                         `last_login_at` datetime DEFAULT NULL,
                         `created_at` datetime DEFAULT NULL,
                         `updated_at` datetime DEFAULT NULL,
                         `deleted_at` datetime DEFAULT NULL,
                         UNIQUE KEY (`login_id`),
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `Users` WRITE;

INSERT INTO `Users` (`login_id`, `password`, `last_login_at`, `created_at`, `updated_at`, `deleted_at`)
VALUES
    ('test_login_id','test_password','2021-01-26 19:03:51','2021-01-26 19:03:51','2021-01-26 19:03:51','2021-01-26 19:03:51');

UNLOCK TABLES;

