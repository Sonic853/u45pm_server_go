package dao

import (
	"database/sql"
	"strconv"
	"time"

	"../game"
)

const createServerTable = "CREATE TABLE `server`(`id` INTEGER NOT NULL PRIMARY KEY UNIQUE,`name` TEXT NOT NULL,`ip` TEXT(15) NOT NULL,`time` TEXT NOT NULL,`key` TEXT NOT NULL,`p1` TEXT,`p2` TEXT,`p3` TEXT,`p4` TEXT);"

func init() {
	initEnvServer()
}

func initEnvServer() {
	db, err := sql.Open("sqlite3", "./user.db")
	checkError(err)
	defer db.Close()

	db.Exec(createServerTable)
}

// ListServer 列出服务器
func ListServer() (lists *game.ServerList) {
	sql := "SELECT * FROM server"
	db := getDB()
	defer db.Close()
	rows, err := db.Query(sql)
	checkError(err)
	defer rows.Close()
	var listz = &game.ServerList{}
	timeinow, err := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	checkError(err)
	timeinow = timeinow - 30
	for rows.Next() {
		var id int
		var name string
		var ip string
		var times string
		var key string
		var p1 string
		var p2 string
		var p3 string
		var p4 string
		var count int
		rows.Scan(&id, &name, &ip, &times, &key, &p1, &p2, &p3, &p4)
		timei, err := strconv.Atoi(times)
		checkError(err)
		if timei < timeinow {
			continue
		}
		if p1 != "NULL" {
			count = count + 1
		}
		if p2 != "NULL" {
			count = count + 1
		}
		if p3 != "NULL" {
			count = count + 1
		}
		if p4 != "NULL" {
			count = count + 1
		}
		listz.List = append(listz.List, game.ServerInfo{ID: id, Name: name, Count: count})
	}
	return listz
}

// FindServerByID 通过ID寻找服务器记录
func FindServerByID(ids string) (server *game.Server) {
	id, err := strconv.Atoi(ids)
	checkError(err)
	sql := "SELECT name, ip, time, key, p1, p2, p3, p4 FROM server WHERE id=?"
	db := getDB()
	defer db.Close()
	rows, err := db.Query(sql, id)
	checkError(err)
	defer rows.Close()
	if rows.Next() {
		var name string
		var ip string
		var time string
		var key string
		var p1 string
		var p2 string
		var p3 string
		var p4 string
		rows.Scan(&name, &ip, &time, &key, &p1, &p2, &p3, &p4)
		server = &game.Server{
			ID:   id,
			Name: name,
			IP:   ip,
			Time: time,
			Key:  key,
			P1:   p1,
			P2:   p2,
			P3:   p3,
			P4:   p4,
		}
	}
	return
}

// AddServer 添加一个服务器
func AddServer(server *game.Server) {
	sql := "INSERT into server(id, name, ip, time, key, p1, p2, p3, p4) values(?,?,?,?,?,?,?,?,?)"
	db := getDB()
	defer db.Close()
	_, err := db.Exec(sql, server.ID, server.Name, server.IP, server.Time, server.Key, server.P1, server.P2, server.P3, server.P4)
	checkError(err)
}

// UpdateServer 更新服务器信息
func UpdateServer(server *game.Server) {
	sql := "UPDATE server SET name=?, ip=?, time=?, key=?, p1=?, p2=?, p3=?, p4=? WHERE id=?"

	db := getDB()
	defer db.Close()

	_, err := db.Exec(sql, server.Name, server.IP, server.Time, server.Key, server.P1, server.P2, server.P3, server.P4, server.ID)
	checkError(err)
}
