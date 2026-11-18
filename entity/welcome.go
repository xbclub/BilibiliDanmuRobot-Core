package entity

type EntryEffectText struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Uid            int64  `json:"uid"`
		PrivilegeType  int    `json:"privilege_type"` // 0 无身份 1 总督? 2 提督 3 舰长
		CopyWriting    string `json:"copy_writing"`
		CopyColor      string `json:"copy_color"`
		HighlightColor string `json:"highlight_color"`
		Priority       int    `json:"priority"`
		Business       int    `json:"business"`
		CopyWritingV2  string `json:"copy_writing_v2"`
		Uinfo          struct {
			Uid  int64 `json:"uid"`
			Base struct {
				Name         string `json:"name"`
				NameColorStr string `json:"name_color_str"`
			} `json:"base"`
			Wealth struct { // 财富信息
				Level int `json:"level"` // 财富等级
			} `json:"wealth"`
			Guard struct {
				Level      int    `json:"level"`
				ExpiredStr string `json:"expired_str"`
			} `json:"guard"`
		} `json:"uinfo"`
	} `json:"data"`
}

type InteractWordText struct {
	Data struct {
		Uname   string `json:"uname"`
		Uid     int64  `json:"uid"`
		MsgType int32  `json:"msg_type"`
	} `json:"data"`
}
