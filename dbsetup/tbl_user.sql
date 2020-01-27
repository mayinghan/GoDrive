CREATE TABLE `tbl_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(32) NOT NULL DEFAULT '' COMMENT 'username',
  `password` varchar(256) NOT NULL DEFAULT '' COMMENT 'hashed pwd',
  `email` varchar(64) NOT NULL DEFAULT '' COMMENT 'email',
  `email_validated` tinyint(1) DEFAULT 0 COMMENT 'if email is validated',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'signup time',
  `role` varchar(8) COMMENT 'role of this user',
  `profile` text COMMENT 'profile of this user',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT 'active/banned/frozen/removed',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;