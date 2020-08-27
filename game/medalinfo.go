package game

// MedalInfo 奖牌信息
type MedalInfo struct {
	//奖牌，只能获得一次，可升级
	ID      int    ` json:"id"`
	GetDate string ` json:"getdate"`
}
