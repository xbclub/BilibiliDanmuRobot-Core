package entity

type RoomBlockMsg struct {
	Cmd   string `json:"cmd"`
	Uid   string `json:"uid"`
	UName string `json:"uname"`
	Data  struct {
		Uid      int    `json:"uid"`
		UName    string `json:"uname"`
		DmScore  int    `json:"dmscore"`
		Operator int    `json:"operator"`
	} `json:"data"`
}
