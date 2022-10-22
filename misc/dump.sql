DROP TABLE IF EXISTS `data_platform_authenticator_business_user_data`;

CREATE TABLE `data_platform_authenticator_business_user_data`
(
    `ID`              int AUTO_INCREMENT,
    `LoginID`         varchar(20) NOT NULL DEFAULT '',
    `BusinessPartner` int(10) NOT NULL,
    `Password`        varchar(70) NOT NULL DEFAULT '',
    `Qos`             varchar(20) NOT NULL DEFAULT 'default',
    `IsEncrypt`       BOOLEAN     NOT NULL DEFAULT TRUE,
    `LastLoginAt`     datetime             DEFAULT NULL,
    `CreatedAt`       datetime             DEFAULT NULL,
    `UpdatedAt`       datetime             DEFAULT NULL,
    `DeletedAt`       datetime             DEFAULT NULL,
    UNIQUE KEY (`LoginId`),
    PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK
TABLES `data_platform_authenticator_business_user_data` WRITE;
INSERT INTO `data_platform_authenticator_business_user_data` (`LoginID`, `BusinessPartner`, `Password`,
                                                              `LastLoginAt`, `CreatedAt`,
                                                              `UpdatedAt`, `DeletedAt`)
VALUES ('test_login_id', 1, 'test_password', '2021-01-26 19:03:51', '2021-01-26 19:03:51', '2021-01-26 19:03:51',
        '2021-01-26 19:03:51');

/*!40000 ALTER TABLE `data_platform_authenticator_business_user_data` ENABLE KEYS */;
UNLOCK
TABLES;
