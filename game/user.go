package game

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"time"
)

// User 用户
type User struct {
	ID       int    ` json:"id"`
	Username string ` json:"username"`
	Name     string ` json:"name"`
	Phone    string ` json:"phone"`
	Password string ` json:"password"`
	Salt     string ` json:"salt"`
	LoginKey string ` json:"loginkey"`
}

// UserJoin 加入服务器的玩家
type UserJoin struct {
	ID     int    ` json:"id"`
	Server int    ` json:"server"`
	Key    string ` json:"key"`
}

// SetPassword 更换密码
func (u *User) SetPassword(password string) string {
	var extrakey string // extrakey := "yourkey"
	u.Salt = SetSalt(extrakey)
	s := sha512.New()
	s.Write([]byte(password + u.Salt))
	h := sha1.New()
	h.Write(s.Sum(nil))
	u.Password = hex.EncodeToString(h.Sum(nil))
	return u.Password
}

// SetSalt 加盐（致死量
func SetSalt(extrakey string) string {
	nowTime := strconv.Itoa(time.Now().Year()) + strconv.Itoa(int(time.Now().Month())) + strconv.Itoa(time.Now().Day()) + strconv.Itoa(time.Now().Nanosecond())
	h := md5.New()
	h.Write([]byte(nowTime + extrakey))
	return hex.EncodeToString(h.Sum(nil))[8:24]
}

// GetSaltedPassword 获取加盐后的密码
func (u *User) GetSaltedPassword(password string, salt string) string {
	s := sha512.New()
	s.Write([]byte(password + salt))
	h := sha1.New()
	h.Write(s.Sum(nil))
	return hex.EncodeToString(h.Sum(nil))
}

// SetLoginKey 更换登录Key
func (u *User) SetLoginKey(password string) string {
	var loginkey string
	if password == u.Password {
		// nowTime := strconv.Itoa(time.Now().Year()) + strconv.Itoa(int(time.Now().Month())) + strconv.Itoa(time.Now().Day()) + strconv.Itoa(time.Now().Nanosecond())
		nowTime := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(time.Now().Nanosecond())
		// log.Println(nowTime)
		h := sha1.New()
		h.Write([]byte(password + nowTime))
		u.LoginKey = hex.EncodeToString(h.Sum(nil))
		loginkey = u.LoginKey
	}
	return loginkey
}
