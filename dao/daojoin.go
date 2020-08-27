package dao

import (
	"database/sql"
	"strconv"

	"../game"
)

const createUserJoinTable = "CREATE TABLE `joins` (`id` INTEGER NOT NULL PRIMARY KEY,`server` integer,`pkey` TEXT);"

func init() {
	initEnvJoin()
}

func initEnvJoin() {
	db, err := sql.Open("sqlite3", "./user.db")
	checkError(err)
	defer db.Close()
	db.Exec(createUserJoinTable)
}

// FindUserJoinByID 通过 id 查找 User
func FindUserJoinByID(ids string) (user *game.UserJoin) {
	id, err := strconv.Atoi(ids)
	checkError(err)
	sql := "SELECT server, pkey FROM joins WHERE id=?"

	db := getDB()
	defer db.Close()

	rows, err := db.Query(sql, id)
	checkError(err)
	defer rows.Close()

	if rows.Next() {
		var server int
		var key string
		rows.Scan(&server, &key)
		user = &game.UserJoin{
			ID:     id,
			Server: server,
			Key:    key,
		}
	}
	return
}

// FindUserJoinByKey 通过 id 查找 User
func FindUserJoinByKey(key string) (user *game.UserJoin) {
	sql := "SELECT id, server FROM joins WHERE pkey=?"

	db := getDB()
	defer db.Close()

	rows, err := db.Query(sql, key)
	checkError(err)
	defer rows.Close()

	if rows.Next() {
		var id int
		var server int
		rows.Scan(&id, &server)
		user = &game.UserJoin{
			ID:     id,
			Server: server,
			Key:    key,
		}
	}
	return
}

// AddUserJoin 添加新 UserJoin
func AddUserJoin(user *game.UserJoin) {
	sql := "INSERT INTO `joins`(id, server, pkey) VALUES(?,?,?)"
	db := getDB()
	defer db.Close()
	_, err := db.Exec(sql, user.ID, user.Server, user.Key)
	checkError(err)
}

// UpdateUserJoin 更新 UserJoin
func UpdateUserJoin(user *game.UserJoin) {
	sql := "UPDATE joins SET pkey=?, server=? WHERE id=?"

	db := getDB()
	defer db.Close()

	_, err := db.Exec(sql, user.Key, user.Server, user.ID)
	checkError(err)
}
