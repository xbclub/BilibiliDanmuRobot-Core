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

type GuardBuyText struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Uid        int    `json:"uid"`
		Username   string `json:"username"`
		GuardLevel int    `json:"guard_level"`
		Num        int    `json:"num"`
		Price      int    `json:"price"`
		GiftID     int    `json:"gift_id"`
		GiftName   string `json:"gift_name"`
		StartTime  int    `json:"start_time"`
		EndTime    int    `json:"end_time"`
	} `json:"data"`
}

type CommonNoticeDanmaku struct {
	Cmd  string `json:"cmd"`
	Data struct {
		ContentSegments []struct {
			FontColor string `json:"font_color"`
			Text      string `json:"text"`
			Type      int    `json:"type"`
		} `json:"content_segments"`
		Dmscore   int   `json:"dmscore"`
		Terminals []int `json:"terminals"`
	} `json:"data"`
}

type SendGiftText struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Action         string `json:"action"`
		BatchComboID   string `json:"batch_combo_id"`
		BatchComboSend struct {
			Action        string      `json:"action"`
			BatchComboID  string      `json:"batch_combo_id"`
			BatchComboNum int         `json:"batch_combo_num"`
			BlindGift     interface{} `json:"blind_gift"`
			GiftID        int         `json:"gift_id"`
			GiftName      string      `json:"gift_name"`
			GiftNum       int         `json:"gift_num"`
			SendMaster    interface{} `json:"send_master"`
			UID           int         `json:"uid"`
			Uname         string      `json:"uname"`
		} `json:"batch_combo_send"`
		Beatid    string `json:"beatId"`
		BizSource string `json:"biz_source"`
		BlindGift struct {
			BlindGiftConfigId string      `json:"blind_gift_config_id"`
			From              string      `json:"from"`
			GiftAction        int         `json:"gift_action"`
			GiftTipPrice      interface{} `json:"gift_tip_price"`
			OriginalGiftId    int         `json:"original_gift_id"`
			OriginalGiftName  string      `json:"original_gift_name"`
			OriginalGiftPrice int         `json:"original_gift_price"`
		} `json:"blind_gift"`
		BroadcastID      int    `json:"broadcast_id"`
		CoinType         string `json:"coin_type"`
		ComboResourcesID int    `json:"combo_resources_id"`
		ComboSend        struct {
			Action     string      `json:"action"`
			ComboID    string      `json:"combo_id"`
			ComboNum   int         `json:"combo_num"`
			GiftID     int         `json:"gift_id"`
			GiftName   string      `json:"gift_name"`
			GiftNum    int         `json:"gift_num"`
			SendMaster interface{} `json:"send_master"`
			UID        int         `json:"uid"`
			Uname      string      `json:"uname"`
		} `json:"combo_send"`
		ComboStayTime  int     `json:"combo_stay_time"`
		ComboTotalCoin int     `json:"combo_total_coin"`
		CritProb       int     `json:"crit_prob"`
		Demarcation    int     `json:"demarcation"`
		Dmscore        int     `json:"dmscore"`
		Draw           int     `json:"draw"`
		Effect         int     `json:"effect"`
		EffectBlock    int     `json:"effect_block"`
		Face           string  `json:"face"`
		Giftid         int     `json:"giftId"`
		GiftName       string  `json:"giftName"`
		Gifttype       int     `json:"giftType"`
		Gold           int     `json:"gold"`
		GuardLevel     int     `json:"guard_level"`
		IsFirst        bool    `json:"is_first"`
		IsSpecialBatch int     `json:"is_special_batch"`
		Magnification  float64 `json:"magnification"`
		MedalInfo      struct {
			AnchorRoomid     int    `json:"anchor_roomid"`
			AnchorUname      string `json:"anchor_uname"`
			GuardLevel       int    `json:"guard_level"`
			IconID           int    `json:"icon_id"`
			IsLighted        int    `json:"is_lighted"`
			MedalColor       int    `json:"medal_color"`
			MedalColorBorder int    `json:"medal_color_border"`
			MedalColorEnd    int    `json:"medal_color_end"`
			MedalColorStart  int    `json:"medal_color_start"`
			MedalLevel       int    `json:"medal_level"`
			MedalName        string `json:"medal_name"`
			Special          string `json:"special"`
			TargetID         int    `json:"target_id"`
		} `json:"medal_info"`
		NameColor         string      `json:"name_color"`
		Num               int         `json:"num"`
		OriginalGiftName  string      `json:"original_gift_name"`
		Price             int         `json:"price"`
		Rcost             int64       `json:"rcost"`
		Remain            int         `json:"remain"`
		Rnd               string      `json:"rnd"`
		SendMaster        interface{} `json:"send_master"`
		Silver            int         `json:"silver"`
		Super             int         `json:"super"`
		SuperBatchGiftNum int         `json:"super_batch_gift_num"`
		SuperGiftNum      int         `json:"super_gift_num"`
		SvgaBlock         int         `json:"svga_block"`
		TagImage          string      `json:"tag_image"`
		Tid               string      `json:"tid"`
		Timestamp         int         `json:"timestamp"`
		TopList           interface{} `json:"top_list"`
		TotalCoin         int         `json:"total_coin"`
		UID               int         `json:"uid"`
		Uname             string      `json:"uname"`
	} `json:"data"`
}

// 某人发了红包
// {"cmd":"POPULARITY_RED_POCKET_NEW","data":{"lot_id":20220211,"start_time":1718850902,"current_time":1718850736,"wait_num":2,"wait_num_v2":0,"uname":"一颗困苏苏丶","uid":651041160,"action":"送出","num":1,"gift_name":"红包","gift_id":13000,"price":20,"name_color":"","medal_info":{"target_id":505986955,"special":"","icon_id":0,"anchor_uname":"","anchor_roomid":0,"medal_level":20,"medal_name":"好困ya","medal_color":13081892,"medal_color_start":13081892,"medal_color_end":13081892,"medal_color_border":13081892,"is_lighted":1,"guard_level":0},"wealth_level":24,"group_medal":null,"is_mystery":false,"sender_info":{"uid":651041160,"base":{"name":"一颗困苏苏丶","face":"https://i0.hdslb.com/bfs/face/114bfd4db9cdbcf57cd36cd3d30041330d6a97b3.jpg","name_color":0,"is_mystery":false,"risk_ctrl_info":null,"origin_info":{"name":"一颗困苏苏丶","face":"https://i0.hdslb.com/bfs/face/114bfd4db9cdbcf57cd36cd3d30041330d6a97b3.jpg"},"official_info":{"role":0,"title":"","desc":"","type":-1},"name_color_str":""},"medal":{"name":"好困ya","level":20,"color_start":13081892,"color_end":13081892,"color_border":13081892,"color":13081892,"id":0,"typ":0,"is_light":1,"ruid":505986955,"guard_level":0,"score":1550451,"guard_icon":"","honor_icon":"","v2_medal_color_start":"#DC6B6B99","v2_medal_color_end":"#DC6B6B99","v2_medal_color_border":"#DC6B6B99","v2_medal_color_text":"#FFFFFFFF","v2_medal_color_level":"#81001F99","user_receive_count":0},"wealth":{"level":24,"dm_icon_key":""},"title":null,"guard":{"level":0,"expired_str":""},"uhead_frame":null,"guard_leader":null},"gift_icon":"","rp_type":0}}   caller=handler/redPocket.go:22
type RedPocketNew struct {
	Cmd  string `json:"cmd"`
	Data struct {
		LotID       int    `json:"lot_id"`
		StartTime   int    `json:"start_time"`
		CurrentTime int    `json:"current_time"`
		WaitNum     int    `json:"wait_num"`
		WaitNumV2   int    `json:"wait_num_v2"`
		Uname       string `json:"uname"`
		Uid         int    `json:"uid"`
		Action      string `json:"action"`
		Num         int    `json:"num"`
		GiftName    string `json:"gift_name"`
		GiftID      int    `json:"gift_id"`
		Price       int    `json:"price"`
		NameColor   string `json:"name_color"`
		MedalInfo   struct {
			TargetID         int    `json:"target_id"`
			Special          string `json:"special"`
			IconID           int    `json:"icon_id"`
			AnchorUname      string `json:"anchor_uname"`
			AnchorRoomid     int    `json:"anchor_roomid"`
			MedalLevel       int    `json:"medal_level"`
			MedalName        string `json:"medal_name"`
			MedalColor       int    `json:"medal_color"`
			MedalColorStart  int    `json:"medal_color_start"`
			MedalColorEnd    int    `json:"medal_color_end"`
			MedalColorBorder int    `json:"medal_color_border"`
			IsLighted        int    `json:"is_lighted"`
			GuardLevel       int    `json:"guard_level"`
		} `json:"medal_info"`
		WealthLevel int         `json:"wealth_level"`
		GroupMedal  interface{} `json:"group_medal"`
		IsMystery   bool        `json:"is_mystery"`
		SenderInfo  struct {
			Uid  int `json:"uid"`
			Base struct {
				Name         string      `json:"name"`
				Face         string      `json:"face"`
				NameColor    int         `json:"name_color"`
				IsMystery    bool        `json:"is_mystery"`
				RiskCtrlInfo interface{} `json:"risk_ctrl_info"`
				OriginInfo   struct {
					Name string `json:"name"`
					Face string `json:"face"`
				} `json:"origin_info"`
				OfficialInfo struct {
					Role  int    `json:"role"`
					Title string `json:"title"`
					Desc  string `json:"desc"`
					Type  int    `json:"type"`
				} `json:"official_info"`
				NameColorStr string `json:"name_color_str"`
			} `json:"base"`
			Medal struct {
				Name               string `json:"name"`
				Level              int    `json:"level"`
				ColorStart         int    `json:"color_start"`
				ColorEnd           int    `json:"color_end"`
				ColorBorder        int    `json:"color_border"`
				Color              int    `json:"color"`
				ID                 int    `json:"id"`
				Typ                int    `json:"typ"`
				IsLight            int    `json:"is_light"`
				Ruid               int    `json:"ruid"`
				GuardLevel         int    `json:"guard_level"`
				Score              int    `json:"score"`
				GuardIcon          string `json:"guard_icon"`
				HonorIcon          string `json:"honor_icon"`
				V2MedalColorStart  string `json:"v2_medal_color_start"`
				V2MedalColorEnd    string `json:"v2_medal_color_end"`
				V2MedalColorBorder string `json:"v2_medal_color_border"`
				V2MedalColorText   string `json:"v2_medal_color_text"`
				V2MedalColorLevel  string `json:"v2_medal_color_level"`
				UserReceiveCount   int    `json:"user_receive_count"`
			} `json:"medal"`
			Wealth struct {
				Level     int    `json:"level"`
				DmIconKey string `json:"dm_icon_key"`
			} `json:"wealth"`
			Title interface{} `json:"title"`
			Guard struct {
				Level      int    `json:"level"`
				ExpiredStr string `json:"expired_str"`
			} `json:"guard"`
			UheadFrame  interface{} `json:"uhead_frame"`
			GuardLeader interface{} `json:"guard_leader"`
		} `json:"sender_info"`
		GiftIcon string `json:"gift_icon"`
		RpType   int    `json:"rp_type"`
	} `json:"data"`
}

// 红包开始
// {"cmd":"POPULARITY_RED_POCKET_START","data":{"lot_id":20220211,"sender_uid":651041160,"sender_name":"一颗困苏苏丶","sender_face":"https://i0.hdslb.com/bfs/face/114bfd4db9cdbcf57cd36cd3d30041330d6a97b3.jpg","join_requirement":1,"danmu":"老板大气！点点红包抽礼物","current_time":1718850902,"start_time":1718850902,"end_time":1718851082,"last_time":180,"remove_time":1718851097,"replace_time":1718851092,"lot_status":1,"h5_url":"https://live.bilibili.com/p/html/live-app-red-envelope/popularity.html?is_live_half_webview=1\u0026hybrid_half_ui=1,5,100p,100p,000000,0,50,0,0,1;2,5,100p,100p,000000,0,50,0,0,1;3,5,100p,100p,000000,0,50,0,0,1;4,5,100p,100p,000000,0,50,0,0,1;5,5,100p,100p,000000,0,50,0,0,1;6,5,100p,100p,000000,0,50,0,0,1;7,5,100p,100p,000000,0,50,0,0,1;8,5,100p,100p,000000,0,50,0,0,1\u0026hybrid_rotate_d=1\u0026hybrid_biz=popularityRedPacket\u0026lotteryId=20220211","user_status":2,"awards":[{"gift_id":31212,"gift_name":"打call","gift_pic":"https://s1.hdslb.com/bfs/live/461be640f60788c1d159ec8d6c5d5cf1ef3d1830.png","num":2},{"gift_id":34003,"gift_name":"人气票","gift_pic":"https://s1.hdslb.com/bfs/live/7164c955ec0ed7537491d189b821cc68f1bea20d.png","num":3},{"gift_id":31216,"gift_name":"小花花","gift_pic":"https://s1.hdslb.com/bfs/live/5126973892625f3a43a8290be6b625b5e54261a5.png","num":3}],"lot_config_id":3,"total_price":1600,"wait_num":0,"wait_num_v2":0,"is_mystery":false,"rp_type":0,"sender_uinfo":{"uid":651041160,"base":{"name":"一颗困苏苏丶","face":"https://i0.hdslb.com/bfs/face/114bfd4db9cdbcf57cd36cd3d30041330d6a97b3.jpg","name_color":0,"is_mystery":false,"risk_ctrl_info":null,"origin_info":{"name":"一颗困苏苏丶","face":"https://i0.hdslb.com/bfs/face/114bfd4db9cdbcf57cd36cd3d30041330d6a97b3.jpg"},"official_info":{"role":0,"title":"","desc":"","type":-1},"name_color_str":""},"medal":null,"wealth":null,"title":null,"guard":null,"uhead_frame":null,"guard_leader":null},"icon_url":"","animation_icon_url":"","rp_guard_info":null}} caller=handler/redPocket.go:13
type RedPocketStart struct {
	Cmd  string `json:"cmd"`
	Data struct {
		LotID           int    `json:"lot_id"`
		SenderUid       int    `json:"sender_uid"`
		SenderName      string `json:"sender_name"`
		SenderFace      string `json:"sender_face"`
		JoinRequirement int    `json:"join_requirement"`
		Danmu           string `json:"danmu"`
		CurrentTime     int    `json:"current_time"`
		StartTime       int    `json:"start_time"`
		EndTime         int    `json:"end_time"`
		LastTime        int    `json:"last_time"`
		RemoveTime      int    `json:"remove_time"`
		ReplaceTime     int    `json:"replace_time"`
		LotStatus       int    `json:"lot_status"`
		H5URL           string `json:"h5_url"`
		UserStatus      int    `json:"user_status"`
		Awards          []struct {
			GiftID   int    `json:"gift_id"`
			GiftName string `json:"gift_name"`
			GiftPic  string `json:"gift_pic"`
			Num      int    `json:"num"`
		} `json:"awards"`
		LotConfigID int  `json:"lot_config_id"`
		TotalPrice  int  `json:"total_price"`
		WaitNum     int  `json:"wait_num"`
		WaitNumV2   int  `json:"wait_num_v2"`
		IsMystery   bool `json:"is_mystery"`
		RpType      int  `json:"rp_type"`
		SenderUinfo struct {
			Uid  int `json:"uid"`
			Base struct {
				Name         string      `json:"name"`
				Face         string      `json:"face"`
				NameColor    int         `json:"name_color"`
				IsMystery    bool        `json:"is_mystery"`
				RiskCtrlInfo interface{} `json:"risk_ctrl_info"`
				OriginInfo   struct {
					Name string `json:"name"`
					Face string `json:"face"`
				} `json:"origin_info"`
				OfficialInfo struct {
					Role  int    `json:"role"`
					Title string `json:"title"`
					Desc  string `json:"desc"`
					Type  int    `json:"type"`
				} `json:"official_info"`
				NameColorStr string `json:"name_color_str"`
			} `json:"base"`
			Medal       interface{} `json:"medal"`
			Wealth      interface{} `json:"wealth"`
			Title       interface{} `json:"title"`
			Guard       interface{} `json:"guard"`
			UheadFrame  interface{} `json:"uhead_frame"`
			GuardLeader interface{} `json:"guard_leader"`
		} `json:"sender_uinfo"`
		IconURL          string      `json:"icon_url"`
		AnimationIconURL string      `json:"animation_icon_url"`
		RpGuardInfo      interface{} `json:"rp_guard_info"`
	} `json:"data"`
}

// 红包结束(中奖信息)
// {"cmd":"POPULARITY_RED_POCKET_WINNER_LIST","data":{"lot_id":20220211,"total_num":8,"award_num":5,"winner_info":[[290386726,"人革联屹立于大地之上",9361206,31212,false,null,1718851084,505986955],[28848915,"本尊-主播真漂亮",9371617,31212,false,null,1718851084,505986955],[3546375218268969,"芽の可愛小迷弟",9307426,34003,false,null,1718851084,505986955],[651041160,"一颗困苏苏丶",9352235,34003,false,null,1718851084,505986955],[30812176,"阿橙嘎嘎猛",9360461,34003,false,null,1718851084,505986955]],"awards":{"31212":{"award_type":1,"award_name":"打call","award_pic":"https://s1.hdslb.com/bfs/live/461be640f60788c1d159ec8d6c5d5cf1ef3d1830.png","award_big_pic":"https://i0.hdslb.com/bfs/live/9e6521c57f24c7149c054d265818d4b82059f2ef.png","award_price":500},"31216":{"award_type":1,"award_name":"小花花","award_pic":"https://s1.hdslb.com/bfs/live/5126973892625f3a43a8290be6b625b5e54261a5.png","award_big_pic":"https://i0.hdslb.com/bfs/live/cf90eac49ac0df5c26312f457e92edfff266f3f1.png","award_price":100},"34003":{"award_type":1,"award_name":"人气票","award_pic":"https://s1.hdslb.com/bfs/live/7164c955ec0ed7537491d189b821cc68f1bea20d.png","award_big_pic":"https://i0.hdslb.com/bfs/live/5bfaddf9a78e677501bb6d440f4d690668136496.png","award_price":100}},"version":1,"rp_type":0,"timestamp":1718851084}}  caller=handler/redPocket.go:32
type RedPocketWinnerList struct {
	Cmd  string `json:"cmd"`
	Data struct {
		LotID      int           `json:"lot_id"`
		TotalNum   int           `json:"total_num"`
		AwardNum   int           `json:"award_num"`
		WinnerInfo []interface{} `json:"winner_info"`
		// 礼物列表 没什么用 去掉了
		// Awards     struct {
		// 	_31212 struct {
		// 		AwardType   int    `json:"award_type"`
		// 		AwardName   string `json:"award_name"`
		// 		AwardPic    string `json:"award_pic"`
		// 		AwardBigPic string `json:"award_big_pic"`
		// 		AwardPrice  int    `json:"award_price"`
		// 	} `json:"31212"`
		// 	_31216 struct {
		// 		AwardType   int    `json:"award_type"`
		// 		AwardName   string `json:"award_name"`
		// 		AwardPic    string `json:"award_pic"`
		// 		AwardBigPic string `json:"award_big_pic"`
		// 		AwardPrice  int    `json:"award_price"`
		// 	} `json:"31216"`
		// 	_34003 struct {
		// 		AwardType   int    `json:"award_type"`
		// 		AwardName   string `json:"award_name"`
		// 		AwardPic    string `json:"award_pic"`
		// 		AwardBigPic string `json:"award_big_pic"`
		// 		AwardPrice  int    `json:"award_price"`
		// 	} `json:"34003"`
		// } `json:"awards"`
		Version   int `json:"version"`
		RpType    int `json:"rp_type"`
		Timestamp int `json:"timestamp"`
	} `json:"data"`
}

type PKProcessInfo struct {
	Cmd  string `json:"cmd"`
	Data struct {
		BattleType int `json:"battle_type"`
		InitInfo   struct {
			RoomId     int    `json:"room_id"`
			Votes      int    `json:"votes"`
			BestUname  string `json:"best_uname"`
			VisionDesc int    `json:"vision_desc"`
		} `json:"init_info"`
		MatchInfo struct {
			RoomId     int    `json:"room_id"`
			Votes      int    `json:"votes"`
			BestUname  string `json:"best_uname"`
			VisionDesc int    `json:"vision_desc"`
		} `json:"match_info"`
	} `json:"data"`
	PkId      int `json:"pk_id"`
	PkStatus  int `json:"pk_status"`
	Timestamp int `json:"timestamp"`
}
type PKStartInfo struct {
	Cmd       string `json:"cmd"`
	PkId      int    `json:"pk_id"`
	PkStatus  int    `json:"pk_status"`
	Timestamp int    `json:"timestamp"`
	Data      struct {
		BattleType    int    `json:"battle_type"`
		FinalHitVotes int    `json:"final_hit_votes"`
		PkStartTime   int    `json:"pk_start_time"`
		PkFrozenTime  int    `json:"pk_frozen_time"`
		PkEndTime     int    `json:"pk_end_time"`
		PkVotesType   int    `json:"pk_votes_type"`
		PkVotesAdd    int    `json:"pk_votes_add"`
		PkVotesName   string `json:"pk_votes_name"`
		StarLightMsg  string `json:"star_light_msg"`
		PkCountdown   int    `json:"pk_countdown"`
		FinalConf     struct {
			Switch    int `json:"switch"`
			StartTime int `json:"start_time"`
			EndTime   int `json:"end_time"`
		} `json:"final_conf"`
		InitInfo struct {
			RoomId     int `json:"room_id"`
			DateStreak int `json:"date_streak"`
		} `json:"init_info"`
		MatchInfo struct {
			RoomId     int `json:"room_id"`
			DateStreak int `json:"date_streak"`
		} `json:"match_info"`
	} `json:"data"`
	Roomid int `json:"roomid,string"`
}
type RankListInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ttl     int    `json:"ttl"`
	Data    struct {
		OnlineNum      int `json:"onlineNum"`
		OnlineRankItem []struct {
			UserRank  int    `json:"userRank"`
			Uid       int64  `json:"uid"`
			Name      string `json:"name"`
			Face      string `json:"face"`
			Score     int    `json:"score"`
			MedalInfo *struct {
				GuardLevel       int    `json:"guardLevel"`
				MedalColorStart  int    `json:"medalColorStart"`
				MedalColorEnd    int    `json:"medalColorEnd"`
				MedalColorBorder int    `json:"medalColorBorder"`
				MedalName        string `json:"medalName"`
				Level            int    `json:"level"`
				TargetId         int64  `json:"targetId"`
				IsLight          int    `json:"isLight"`
			} `json:"medalInfo"`
			GuardLevel int `json:"guard_level"`
		} `json:"OnlineRankItem"`
		OwnInfo struct {
			Uid        int    `json:"uid"`
			Name       string `json:"name"`
			Face       string `json:"face"`
			Rank       int    `json:"rank"`
			NeedScore  int    `json:"needScore"`
			Score      int    `json:"score"`
			GuardLevel int    `json:"guard_level"`
		} `json:"ownInfo"`
		TipsText  string `json:"tips_text"`
		ValueText string `json:"value_text"`
	} `json:"data"`
}
