package game

// UserInfo 用户信息
type UserInfo struct {
	ID       int    ` json:"id"`
	Username string ` json:"username"`
	Name     string ` json:"name"`
	Phone    string ` json:"phone"`
	Password string ` json:"password"`
	LoginKey string ` json:"loginkey"`
}

// SetPhone 设置手机号（显示用
func (u *UserInfo) SetPhone(phone string) {
	u.Phone = string([]rune(phone)[0:3]) + "xxxx" + string([]rune(phone)[len([]rune(phone))-4:])
}
