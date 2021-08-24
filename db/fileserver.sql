/*
Navicat MySQL Data Transfer

Source Server         : 本地
Source Server Version : 50726
Source Host           : localhost:3306
Source Database       : fileserver

Target Server Type    : MYSQL
Target Server Version : 50726
File Encoding         : 65001

Date: 2021-08-24 15:22:29
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for tbl_file
-- ----------------------------
DROP TABLE IF EXISTS `tbl_file`;
CREATE TABLE `tbl_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
  `update_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态 可用 禁用 已删除',
  `ext1` int(11) NOT NULL DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_file_hash` (`file_sha1`) USING BTREE,
  KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_token
-- ----------------------------
DROP TABLE IF EXISTS `tbl_token`;
CREATE TABLE `tbl_token` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_token` varchar(255) NOT NULL DEFAULT '' COMMENT 'token',
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_user
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user`;
CREATE TABLE `tbl_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pw` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `email` varchar(64) NOT NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(255) NOT NULL DEFAULT '' COMMENT '手机',
  `email_validated` tinyint(1) NOT NULL DEFAULT '0' COMMENT '邮箱是否已验证',
  `phone_validated` tinyint(1) NOT NULL DEFAULT '0' COMMENT '手机是否验证了',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃',
  `profile` text COMMENT '用户属性',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '用户状态 ，启用删除 禁用等',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_phone` (`phone`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for tbl_user_file
-- ----------------------------
DROP TABLE IF EXISTS `tbl_user_file`;
CREATE TABLE `tbl_user_file` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '',
  `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小',
  `file_name` varchar(64) NOT NULL DEFAULT '' COMMENT '文件名称',
  `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `last_update` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0正常 1删除 2禁用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_file` (`user_name`,`file_sha1`) USING BTREE,
  KEY `idx_status` (`status`),
  KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
