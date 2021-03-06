package game

// ErrCode 错误系列
type ErrCode struct {
	Code    int    ` json:"code"`
	Message string ` json:"message"`
}

var (
	_CodeMessage = map[int]string{
		1000: "登陆成功",
		1001: "账号或密码错误",
		1002: "账号已封停", //石锤外挂，别给我整有的没的，石锤就是石锤，给我重新注册一个账号去！
		1003: "账号已封停", //违禁、可申请更改解封
		1004: "账号已封停", //严重违禁，不是一个客服就能搞定的
		1005: "账号已封停", //证实小孩败家系列账号，不可解封
		1006: "登录已过期",
		1100: "注册成功",
		1101: "注册信息错误",
		1102: "账号名称错误",
		1103: "账号密码错误",
		1104: "手机号错误",
		1105: "账号名称已存在",
		1106: "手机号已存在",
		1107: "两次密码不相符",
		1108: "字段不能为空",
		1109: "游戏名称错误",
		2001: "找不到服务器",
		2002: "服务器已满",
		2003: "服务器已拒绝你加入",
		2004: "服务器已拒绝封停用户",
		2005: "服务器已被清空",
		2006: "记录服务器与当前服务器不符",
		5000: "系统繁忙",
		5001: "奇怪的操作",
	}
)

// SetCode 设置错误代码
func (c *ErrCode) SetCode(code int) {
	c.Code = code
	c.Message = _CodeMessage[code]
}
