/*
 Navicat Premium Data Transfer

 Source Server         : user
 Source Server Type    : SQLite
 Source Server Version : 3030001
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3030001
 File Encoding         : 65001

 Date: 28/08/2020 01:53:29
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for joins
-- ----------------------------
DROP TABLE IF EXISTS "joins";
CREATE TABLE "joins" (
  "id" INTEGER NOT NULL,
  "server" integer,
  "pkey" TEXT,
  PRIMARY KEY ("id")
);

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

-- ----------------------------
-- Table structure for sqlite_sequence
-- ----------------------------
DROP TABLE IF EXISTS "sqlite_sequence";
CREATE TABLE "sqlite_sequence" (
  "name" ,
  "seq" 
);

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS "user";
CREATE TABLE "user" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "username" TEXT NOT NULL,
  "name" TEXT,
  "phone" TEXT NOT NULL,
  "password" TEXT NOT NULL,
  "salt" TEXT NOT NULL,
  "loginkey" TEXT NOT NULL,
  UNIQUE ("id" ASC)
);

-- ----------------------------
-- Auto increment value for user
-- ----------------------------
UPDATE "sqlite_sequence" SET seq = 15 WHERE name = 'user';

PRAGMA foreign_keys = true;
