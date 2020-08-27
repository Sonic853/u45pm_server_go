package dao

import (
	"database/sql"
	"log"
	"strconv"

	"../game"
)

const createUserTable = "CREATE TABLE `user` ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, `username` TEXT NOT NULL, `name` TEXT, `phone` TEXT NOT NULL, `password` TEXT NOT NULL, `salt` TEXT NOT NULL, `loginkey` TEXT NOT NULL );"

const createAdmin = "INSERT INTO user(id, username, name, phone, password, salt, loginkey) VALUES(1, 'admin', '管理员', '13322224444', '20c136fb22a0ad88e280d2b0c63c152e361309f8', '54f7c1be92f64550', 'loginkey');"

func init() {
	initEnv()
	initAdmin()
}

func initEnv() {
	db, err := sql.Open("sqlite3", "./user.db")
	checkError(err)
	defer db.Close()
	db.Exec(createUserTable)
}

func initAdmin() {
	db, err := sql.Open("sqlite3", "./user.db")
	checkError(err)
	defer db.Close()
	db.Exec(createAdmin)
}

func getDB() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./user.db")
	checkError(err)
	return
}

// FindUserByUsername 通过 username 查找 User
func FindUserByUsername(username string) (user *game.User) {
	sql := `SELECT id, name, phone, password, salt, loginkey FROM user WHERE username=?`

	db := getDB()

	rows, err := db.Query(sql, username)
	checkError(err)
	var id int
	var name string
	var phone string
	var password string
	var salt string
	var loginkey string

	if rows.Next() {
		// log.Println(rows.Scan(&id, &name, &phone, &password, &salt, &loginkey))
		rows.Scan(&id, &name, &phone, &password, &salt, &loginkey)

		user = &game.User{
			ID:       id,
			Username: username,
			Name:     name,
			Phone:    phone,
			Password: password,
			Salt:     salt,
			LoginKey: loginkey,
		}
	}
	defer db.Close()
	defer rows.Close()
	return
}

// FindUserByPhone 通过 phone 查找 User
func FindUserByPhone(phone string) (user *game.User) {
	sql := "SELECT id, username, name, password, salt, loginkey FROM user WHERE phone=?"

	db := getDB()
	defer db.Close()

	rows, err := db.Query(sql, phone)
	checkError(err)
	defer rows.Close()

	if rows.Next() {
		var id int
		var username string
		var name string
		var password string
		var salt string
		var loginkey string
		rows.Scan(&id, &username, &name, &password, &salt, &loginkey)

		user = &game.User{
			ID:       id,
			Username: username,
			Name:     name,
			Phone:    phone,
			Password: password,
			Salt:     salt,
			LoginKey: loginkey,
		}
	}
	return
}

// FindUserByID 通过 id 查找 User
func FindUserByID(ids string) (user *game.User) {
	id, err := strconv.Atoi(ids)
	checkError(err)
	sql := "SELECT username, name, phone, password, salt, loginkey FROM user WHERE id=?"

	db := getDB()
	defer db.Close()

	rows, err := db.Query(sql, id)
	checkError(err)
	defer rows.Close()

	if rows.Next() {
		var username string
		var name string
		var phone string
		var password string
		var salt string
		var loginkey string
		rows.Scan(&username, &name, &phone, &password, &salt, &loginkey)

		user = &game.User{
			ID:       id,
			Username: username,
			Name:     name,
			Phone:    phone,
			Password: password,
			Salt:     salt,
			LoginKey: loginkey,
		}
	}
	return
}

// FindUserByUsernameAndPassword 通过 username 和 password 查找 User
func FindUserByUsernameAndPassword(username string, password string) (user *game.User) {
	sql := "SELECT id, phone FROM user WHERE username=? AND password=?"

	db := getDB()
	defer db.Close()

	rows, err := db.Query(sql, username, password)
	checkError(err)
	defer rows.Close()

	if rows.Next() {
		var id int
		var phone string
		rows.Scan(&id, &phone)

		user = &game.User{
			ID:       id,
			Username: username,
			Phone:    phone,
			Password: password,
		}
	}
	return
}

// FindUserByPhoneAndPassword 通过 phone 和 password 查找 User
func FindUserByPhoneAndPassword(phone string, password string) (user *game.User) {
	sql := "SELECT id, username FROM user WHERE phone=? AND password=?"

	db := getDB()
	defer db.Close()

	rows, err := db.Query(sql, phone, password)
	checkError(err)
	defer rows.Close()

	if rows.Next() {
		var id int
		var username string
		rows.Scan(&id, &username)

		user = &game.User{
			ID:       id,
			Username: username,
			Phone:    phone,
			Password: password,
		}
	}
	return
}

// AddUser 添加新 User
func AddUser(user *game.User) {
	sql := "INSERT INTO `user`(username, name, phone, password, salt, loginkey) VALUES(?,?,?,?,?,?)"

	db := getDB()
	defer db.Close()

	_, err := db.Exec(sql, user.Username, user.Name, user.Phone, user.Password, user.Salt, user.LoginKey)
	checkError(err)
}

// UpdateUser 更新 User
func UpdateUser(user *game.User) {
	sql := "UPDATE user SET username=?, name=?, phone=?, password=?, salt=? WHERE id=?"

	db := getDB()
	defer db.Close()

	_, err := db.Exec(sql, user.Username, user.Name, user.Phone, user.Password, user.Salt, user.ID)
	checkError(err)
}

// DeleteUser 删除 User
func DeleteUser(id int) {
	sql := "DELETE FROM user WHERE id=?"

	db := getDB()
	defer db.Close()

	_, err := db.Exec(sql, id)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// UpdateLoginKey 更新 LoginKey
func UpdateLoginKey(user *game.User) {
	sql := "UPDATE user SET loginkey=? WHERE id=?"

	db := getDB()
	defer db.Close()

	_, err := db.Exec(sql, user.LoginKey, user.ID)
	checkError(err)
}
