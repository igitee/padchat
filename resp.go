package padchat

import "encoding/json"

type ServerData struct {
	Type   string
	Event  string
	TaskID string
	CMDID  string `json:"cmdId"`
	Data   json.RawMessage
	Msg    string
}

type ScanResp struct {
	DeviceType  string `json:"device_type"`
	ExpiredTime int    `json:"expired_time"`
	HeadURL     string `json:"head_url"`
	NickName    string `json:"nick_name"`
	Password    string `json:"password"`
	Status      int    `json:"status"`
	UserName    string `json:"user_name"`
	External    string `json:"external"`
	Email       string `json:"email"`
	Uin         int    `json:"uin"`
	PhoneNumber string `json:"phone_number"`
	SubStatus   int    `json:"sub_status"`
}

type PushResp struct {
	List []json.RawMessage `json:"list"`
}

type Msg struct {
	Data        string          `json:"data,omitempty"`
	Content     json.RawMessage `json:"content"`
	Continue    int             `json:"continue"`
	Description string          `json:"description"`
	Status      int             `json:"status"`
	Timestamp   int             `json:"timestamp"`
	Uin         int             `json:"uin"`
	FromUser    string          `json:"from_user"`
	MsgID       string          `json:"msg_id"`
	MsgSource   string          `json:"msg_source"`
	MsgType     int             `json:"msg_type"`
	SubType     int             `json:"sub_type"`
	ToUser      string          `json:"to_user"`
	MType       int             `json:"m_type"`
}

type CommandResp struct {
	Success bool
	Data    json.RawMessage
	Msg     string
}

type LoginTokenResp struct {
	Token string `json:"token"`
	Uin   int64  `json:"uin"`
}

type MyInfoResp struct {
	UserName string `json:"userName"`
	Uin      int64  `json:"uin"`
}

type SendMsgResp struct {
	Message string `json:"message"`
	MsgID   string `json:"msg_id"`
	Status  int    `json:"status"`
}

type MsgImageResp struct {
	Image   string `json:"image"`
	Message string `json:"message"`
	Size    int    `json:"size"`
	Status  int    `json:"status"`
}
