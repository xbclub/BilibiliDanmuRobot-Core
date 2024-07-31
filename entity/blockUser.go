package entity

type RoomBlockMsg struct {
	Cmd   string `json:"cmd"`
	Uid   int    `json:"uid"`
	UName string `json:"uname"`
	Data  struct {
		BlockExpired int    `json:"block_expired"`
		Dmscore      int    `json:"dmscore"`
		Operator     int    `json:"operator"`
		Uid          int    `json:"uid"`
		Uname        string `json:"uname"`
		VaildPeriod  string `json:"vaild_period"`
	} `json:"data"`
}
