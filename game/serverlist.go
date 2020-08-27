package game

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

// ServerList 服务器列表
type ServerList struct {
	List []ServerInfo ` json:"serverlist"`
}

// Server 服务器信息
type Server struct {
	ID   int    ` json:"id"`
	Name string ` json:"name"`
	IP   string ` json:"ip"`
	Time string ` json:"time"`
	Key  string ` json:"key"`
	P1   string ` json:"p1"`
	P2   string ` json:"p2"`
	P3   string ` json:"p3"`
	P4   string ` json:"p4"`
}

// ServerInfo 服务器信息（公开显示的
type ServerInfo struct {
	ID    int    ` json:"id"`
	Name  string ` json:"name"`
	Count int    ` json:"count"`
}

// JoinServer 提供key以及IP给玩家加入服务器
type JoinServer struct {
	IP  string ` json:"ip"`
	Key string ` json:"key"`
}

// SetServerKey 设置服务器key
func (s *Server) SetServerKey(user *User) string {
	h := md5.New()
	h.Write([]byte(user.Username + s.IP + strconv.Itoa(time.Now().Nanosecond())))
	s.Key = hex.EncodeToString(h.Sum(nil))[8:24]
	return s.Key
}

// KickPlayers 踢出全部玩家
func (s *Server) KickPlayers() {
	s.P1 = "NULL"
	s.P2 = "NULL"
	s.P3 = "NULL"
	s.P4 = "NULL"
}

// KickPlayer 踢出部分玩家
func (s *Server) KickPlayer(player int) {
	switch player {
	case 0:
		s.P1 = "NULL"
		break
	case 1:
		s.P2 = "NULL"
		break
	case 2:
		s.P3 = "NULL"
		break
	case 3:
		s.P4 = "NULL"
		break
	}
}

// AddPlayer 添加玩家
func (s *Server) AddPlayer(player int, key string) {
	switch player {
	case 0:
		s.P1 = key
		break
	case 1:
		s.P2 = key
		break
	case 2:
		s.P3 = key
		break
	case 3:
		s.P4 = key
		break
	}
}
