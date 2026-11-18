package entity

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
				TargetID         int    `json:"targetId"`
				IsLight          int    `json:"isLight"`
			} `json:"medalInfo,omitempty"`
			GuardLevel  int  `json:"guard_level"`
			WealthLevel int  `json:"wealth_level"`
			IsMystery   bool `json:"is_mystery"`
			Uinfo       struct {
				Uid  int `json:"uid"`
				Base struct {
					Name         string `json:"name"`
					Face         string `json:"face"`
					NameColor    int    `json:"name_color"`
					IsMystery    bool   `json:"is_mystery"`
					RiskCtrlInfo struct {
						Name string `json:"name"`
						Face string `json:"face"`
					} `json:"risk_ctrl_info"`
					OriginInfo struct {
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
				Medal *struct {
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
				} `json:"medal,omitempty"`
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
			} `json:"uinfo"`
		} `json:"OnlineRankItem"`
		OwnInfo struct {
			Uid         int         `json:"uid"`
			Name        string      `json:"name"`
			Face        string      `json:"face"`
			Rank        int         `json:"rank"`
			NeedScore   int         `json:"needScore"`
			Score       int         `json:"score"`
			GuardLevel  int         `json:"guard_level"`
			WealthLevel int         `json:"wealth_level"`
			IsMystery   bool        `json:"is_mystery"`
			Uinfo       interface{} `json:"uinfo"`
		} `json:"ownInfo"`
		TipsText  string `json:"tips_text"`
		ValueText string `json:"value_text"`
		Ab        struct {
			GuardAccompanyList int `json:"guard_accompany_list"`
		} `json:"ab"`
		OnlineNumText string `json:"onlineNumText"`
	} `json:"data"`
}
