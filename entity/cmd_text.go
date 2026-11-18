package entity

type CmdText struct {
	Cmd string `json:"cmd"`
}

type DanmuMsgText struct {
	Info []interface{} `json:"info"`
}

type DanmuMsgTextReplyInfo struct {
	ReplyUid   string
	ReplyMsgId string
}

type Bullet struct {
	Msg   string
	Reply []*DanmuMsgTextReplyInfo
}

type DanmuMsgTextInfo0Extra struct {
	SendFromMe            bool        `json:"send_from_me"`
	Mode                  int         `json:"mode"`
	Color                 int         `json:"color"`
	DmType                int         `json:"dm_type"`
	FontSize              int         `json:"font_size"`
	PlayerMode            int         `json:"player_mode"`
	ShowPlayerType        int         `json:"show_player_type"`
	Content               string      `json:"content"`
	UserHash              string      `json:"user_hash"`
	EmoticonUnique        string      `json:"emoticon_unique"`
	BulgeDisplay          int         `json:"bulge_display"`
	RecommendScore        int         `json:"recommend_score"`
	MainStateDmColor      string      `json:"main_state_dm_color"`
	ObjectiveStateDmColor string      `json:"objective_state_dm_color"`
	Direction             int         `json:"direction"`
	PkDirection           int         `json:"pk_direction"`
	QuartetDirection      int         `json:"quartet_direction"`
	AnniversaryCrowd      int         `json:"anniversary_crowd"`
	YeahSpaceType         string      `json:"yeah_space_type"`
	YeahSpaceURL          string      `json:"yeah_space_url"`
	JumpToURL             string      `json:"jump_to_url"`
	SpaceType             string      `json:"space_type"`
	SpaceURL              string      `json:"space_url"`
	Animation             interface{} `json:"animation"`
	Emots                 interface{} `json:"emots"`
	IsAudited             bool        `json:"is_audited"`
	IdStr                 string      `json:"id_str"`
	Icon                  interface{} `json:"icon"`
	ShowReply             bool        `json:"show_reply"`
	ReplyMid              int         `json:"reply_mid"`
	ReplyUname            string      `json:"reply_uname"`
	ReplyUnameColor       string      `json:"reply_uname_color"`
	ReplyIsMystery        bool        `json:"reply_is_mystery"`
	HitCombo              int         `json:"hit_combo"`
}
