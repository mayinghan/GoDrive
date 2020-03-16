CREATE TABLE `tbl_userfile` (
  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `hash` varchar(64) NOT NULL DEFAULT '' COMMENT 'hash',
  `size` bigint(20) DEFAULT '0' COMMENT 'file size',
  `filename` varchar(256) NOT NULL DEFAULT '' COMMENT 'filename',
  `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'upload time',
  `last_update` datetime DEFAULT CURRENT_TIMESTAMP 
          ON UPDATE CURRENT_TIMESTAMP COMMENT 'last update time',
  UNIQUE KEY `idx_user_file` (`username`, `hash`, `filename`),
  KEY `idx_user_id` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- //remove status column

-- ALTER TABLE tbl_userfile
-- DROP status;

-- change unique key constraints
-- ALTER TABLE tbl_userfile DROP index idx_user_file;
-- ALTER TABLE tbl_userfile ADD CONSTRAINT idx_user_file UNIQUE (username, hash, filename);