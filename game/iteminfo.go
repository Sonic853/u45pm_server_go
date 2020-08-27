package game

// ItemInfo 皮肤信息
type ItemInfo struct {
	/*
		有想过单独ID，但是估计会耗性能，故加双ID。
		一个是单独ID保证唯一性，另一个ID方便分类皮肤。
	*/
	ID           int    ` json:"id"`
	SkinID       int    ` json:"skinid"`
	SkinAbrasion string ` json:"abrasion"`
	TagName      string ` json:"tagname"`
}
