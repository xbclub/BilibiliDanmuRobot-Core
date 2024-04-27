package entity

// LotteryRequest 抽奖请求
type LotteryRequest struct {
	Msg      string `json:"msg"`
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	RoomID   int64  `json:"room_id"`
	Version  string `json:"version"`
}

// LotteryResponse 抽奖响应
type LotteryResponse struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	GiftName string `json:"gift_name"`
	Count    int    `json:"count"`
}
