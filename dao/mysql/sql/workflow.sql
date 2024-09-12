/*
Navicat MySQL Data Transfer

Source Server         : 10.0.0.180
Source Server Version : 50744
Source Host           : 10.0.0.180:3306
Source Database       : k8s_platform

Target Server Type    : MYSQL
Target Server Version : 50744
File Encoding         : 65001

Date: 2024-09-09 15:20:18
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for workflow
-- ----------------------------
DROP TABLE IF EXISTS `workflow`;
CREATE TABLE `workflow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `namespace` varchar(32) DEFAULT NULL,
  `replicas` int(11) DEFAULT NULL,
  `deployment` varchar(32) DEFAULT NULL,
  `service` varchar(32) DEFAULT NULL,
  `ingress` varchar(32) DEFAULT NULL,
  `type` varchar(32) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
