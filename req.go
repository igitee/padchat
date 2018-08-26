package padchat

type WSReq struct {
	Type  string      `json:"type"`
	Cmd   string      `json:"cmd"`
	CmdID string      `json:"cmdId"`
	Data  interface{} `json:"data"`
}

type LoginReq struct {
	LoginType string `json:"loginType"`
	WXData    string `json:"wxData"`
	Token     string `json:"token"`
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	Code      string `json:"code"`
}

type SendMsgReq struct {
	ToUserName string   `json:"toUserName"`
	Content    string   `json:"content"`
	AtList     []string `json:"atList"`
	File       string   `json:"file"` // base64
}
