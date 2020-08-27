/*
 Navicat Premium Data Transfer

 Source Server         : server
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 28/08/2020 01:51:30
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for server
-- ----------------------------
DROP TABLE IF EXISTS "server";
CREATE TABLE "server" (
  "id" INTEGER NOT NULL,
  "name" TEXT NOT NULL,
  "ip" TEXT(15) NOT NULL,
  "time" TEXT NOT NULL,
  "key" TEXT NOT NULL,
  "p1" TEXT,
  "p2" TEXT,
  "p3" TEXT,
  "p4" TEXT,
  PRIMARY KEY ("id"),
  UNIQUE ("id" ASC)
);

PRAGMA foreign_keys = true;
