/* sql to create table for files */
CREATE TABLE `tbl_file`
(
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sha1` char(40) NOT NULL DEFAULT '' COMMENT 'file''s hash',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT 'file name',
  `size` bigint(20) DEFAULT 0 COMMENT 'file size',
  `location` varchar(512) NOT NULL DEFAULT '' COMMENT 'file location',
  `create_at` datetime DEFAULT NOW() COMMENT 'create date',
  `update_at` datetime DEFAULT NOW() COMMENT 'update date',
  `status` int NOT NULL DEFAULT 0 COMMENT 'status code(available/delete/banned)',
  `ext1` text COMMENT 'backup info, not neccessarilly gonna be used',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_sha1`),
  key `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;