CREATE TABLE `tbl_userfile` (
  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `hash` varchar(64) NOT NULL DEFAULT '' COMMENT 'hash',
  `size` bigint(20) DEFAULT '0' COMMENT 'file size',
  `filename` varchar(256) NOT NULL DEFAULT '' COMMENT 'filename',
  `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'upload time',
  `last_update` datetime DEFAULT CURRENT_TIMESTAMP 
          ON UPDATE CURRENT_TIMESTAMP COMMENT 'last update time',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT 'file status (0:available 1:deleted 2: banned)', //removed
  UNIQUE KEY `idx_user_file` (`username`, `hash`),
  KEY `idx_status` (`status`),
  KEY `idx_user_id` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


//remove status column

ALTER TABLE tbl_userfile
DROP status;