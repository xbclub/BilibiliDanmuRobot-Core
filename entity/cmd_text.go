package entity

type CmdText struct {
	Cmd string `json:"cmd"`
}

type DanmuMsgText struct {
	Info []interface{} `json:"info"`
}

type EntryEffectText struct {
	Data struct {
		Uid         int64  `json:"uid"`
		CopyWriting string `json:"copy_writing"`
	} `json:"data"`
}

type InteractWordText struct {
	Data struct {
		Uname   string `json:"uname"`
		Uid     int64  `json:"uid"`
		MsgType int32  `json:"msg_type"`
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
		Beatid           string      `json:"beatId"`
		BizSource        string      `json:"biz_source"`
		BlindGift        interface{} `json:"blind_gift"`
		BroadcastID      int         `json:"broadcast_id"`
		CoinType         string      `json:"coin_type"`
		ComboResourcesID int         `json:"combo_resources_id"`
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
