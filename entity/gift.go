package entity

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
type RedPocketWinnerList struct {
	Cmd  string `json:"cmd"`
	Data struct {
		LotID      int           `json:"lot_id"`
		TotalNum   int           `json:"total_num"`
		AwardNum   int           `json:"award_num"`
		WinnerInfo []interface{} `json:"winner_info"`
		Version    int           `json:"version"`
		RpType     int           `json:"rp_type"`
		Timestamp  int           `json:"timestamp"`
	} `json:"data"`
}
