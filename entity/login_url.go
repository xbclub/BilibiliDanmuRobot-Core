package entity

type LoginUrl struct {
	Data struct {
		Url      string `json:"url"`
		OauthKey string `json:"qrcode_key"`
	} `json:"data"`
}

//	type LoginInfoPre struct {
//		Status bool `json:"status"`
//	}
type LoginInfoPre struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type LoginInfoData struct {
	Data struct {
		Url          string `json:"url"`
		RefreshToken string `json:"refresh_token"`
		Code         int    `json:"code"`
		Message      string `json:"message"`
	} `json:"data"`
}
type LoginInfoCookies struct {
	SetCookie []string `json:"Set-Cookie"`
}
