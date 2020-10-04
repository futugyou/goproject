CREATE DATABASE if NOT EXISTS blog_service DEFAULT CHARACTER
SET UTF8MB4 DEFAULT  collate UTF8MB4_GENERAL_CI;

CREATE TABLE `blog_tag`(
`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
`name` VARCHAR(100) DEFAULT '' COMMENT 'tag name',
`state` TINYINT(3) UNSIGNED DEFAULT '1' COMMENT '0not use ,1canuse',
PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=UTF8MB4 COMMENT='tag management';


CREATE TABLE `blog_article`(
`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
`title` VARCHAR(100) DEFAULT '' COMMENT 'article title',
`desc` VARCHAR(255) DEFAULT '' COMMENT 'article desc',
`cover_image_url` VARCHAR(255) DEFAULT '' COMMENT 'article image',
`content` LONGTEXT COMMENT 'article content',
`created_on` INT(10) UNSIGNED DEFAULT '0' COMMENT 'create timme',
`create_by` VARCHAR(100) DEFAULT '' COMMENT 'create by',
`modified_on` INT(10) UNSIGNED DEFAULT '0' COMMENT 'modified timme',
`modified_by` VARCHAR(100) DEFAULT '' COMMENT 'modified by',
`deleted_on` INT(10) UNSIGNED DEFAULT '0' COMMENT 'deleted time',
`is_del` TINYINT(3) UNSIGNED  DEFAULT '0' COMMENT '0not delete ,1deleted',
`state` TINYINT(3) UNSIGNED DEFAULT '1' COMMENT '0not use ,1canuse',
PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=UTF8MB4 COMMENT='article management';

CREATE TABLE `blog_article_tag`(
`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
`article_id` INT(11)    NOT NULL   COMMENT 'article id',
`tag_id` INT(10) UNSIGNED  NOT NULL  DEFAULT '0' COMMENT 'tag id',
`created_on` INT(10) UNSIGNED DEFAULT '0' COMMENT 'create timme',
`create_by` VARCHAR(100) DEFAULT '' COMMENT 'create by',
`modified_on` INT(10) UNSIGNED DEFAULT '0' COMMENT 'modified timme',
`modified_by` VARCHAR(100) DEFAULT '' COMMENT 'modified by',
`deleted_on` INT(10) UNSIGNED DEFAULT '0' COMMENT 'deleted time',
`is_del` TINYINT(3)UNSIGNED  DEFAULT '0' COMMENT '0not delete ,1deleted',
PRIMARY KEY (`id`)
) ENGINE=INNODB DEFAULT CHARSET=UTF8MB4 COMMENT='article tag rel management';



