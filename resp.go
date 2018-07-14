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
	List []Msg `json:"list"`
}

type Msg struct {
	Data            string `json:"data"`
	Content         string `json:"content"`
	Description     string `json:"description"`
	FromUser        string `json:"from_user"`
	MsgID           string `json:"msg_id"`
	MsgSource       string `json:"msg_source"`
	SubType         int    `json:"sub_type"`
	Timestamp       int    `json:"timestamp"`
	ToUser          string `json:"to_user"`
	Continue        int    `json:"continue"`
	MsgType         int    `json:"msg_type"`
	Status          int    `json:"status"`
	Uin             int    `json:"uin"`
	BigHead         string `json:"big_head,omitempty"`
	BitMask         int64  `json:"bit_mask,omitempty"`
	BitValue        int    `json:"bit_value,omitempty"`
	ChatroomID      int    `json:"chatroom_id,omitempty"`
	ChatroomOwner   string `json:"chatroom_owner,omitempty"`
	City            string `json:"city,omitempty"`
	Country         string `json:"country,omitempty"`
	ID              int    `json:"id,omitempty"`
	ImgFlag         int    `json:"img_flag,omitempty"`
	Intro           string `json:"intro,omitempty"`
	Label           string `json:"label,omitempty"`
	Level           int    `json:"level,omitempty"`
	MaxMemberCount  int    `json:"max_member_count,omitempty"`
	MemberCount     int    `json:"member_count,omitempty"`
	NickName        string `json:"nick_name,omitempty"`
	Provincia       string `json:"provincia,omitempty"`
	PyInitial       string `json:"py_initial,omitempty"`
	QuanPin         string `json:"quan_pin,omitempty"`
	Remark          string `json:"remark,omitempty"`
	RemarkPyInitial string `json:"remark_py_initial,omitempty"`
	RemarkQuanPin   string `json:"remark_quan_pin,omitempty"`
	Sex             int    `json:"sex,omitempty"`
	Signature       string `json:"signature,omitempty"`
	SmallHead       string `json:"small_head,omitempty"`
	Source          int    `json:"source,omitempty"`
	Stranger        string `json:"stranger,omitempty"`
	UserName        string `json:"user_name,omitempty"`
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
