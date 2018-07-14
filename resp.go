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

type Contact struct {
	BigHead         string `json:"big_head"`
	BitMask         int64  `json:"bit_mask"`
	BitValue        int    `json:"bit_value"`
	ChatroomID      int    `json:"chatroom_id"`
	ChatroomOwner   string `json:"chatroom_owner"`
	City            string `json:"city"`
	Continue        int    `json:"continue"`
	Country         string `json:"country"`
	ID              int    `json:"id"`
	ImgFlag         int    `json:"img_flag"`
	Intro           string `json:"intro"`
	Label           string `json:"label"`
	Level           int    `json:"level"`
	Member          string `json:"member"`
	MaxMemberCount  int    `json:"max_member_count"`
	MemberCount     int    `json:"member_count"`
	MsgType         int    `json:"msg_type"`
	NickName        string `json:"nick_name"`
	Provincia       string `json:"provincia"`
	PyInitial       string `json:"py_initial"`
	QuanPin         string `json:"quan_pin"`
	Remark          string `json:"remark"`
	RemarkPyInitial string `json:"remark_py_initial"`
	RemarkQuanPin   string `json:"remark_quan_pin"`
	Sex             int    `json:"sex"`
	Signature       string `json:"signature"`
	SmallHead       string `json:"small_head"`
	Source          int    `json:"source"`
	Status          int    `json:"status"`
	Stranger        string `json:"stranger"`
	Uin             int    `json:"uin"`
	UserName        string `json:"user_name"`
}

type ChatroomInfo struct {
	ChatroomID int              `json:"chatroom_id"`
	Count      int              `json:"count"`
	Member     string           `json:"member"`
	Members    []ChatMemberInfo `json:"-"`
	Message    string           `json:"message"`
	Status     int              `json:"status"`
	UserName   string           `json:"user_name"`
}

type ChatMemberInfo struct {
	BigHead          string `json:"big_head"`
	ChatroomNickName string `json:"chatroom_nick_name"`
	InvitedBy        string `json:"invited_by"`
	NickName         string `json:"nick_name"`
	SmallHead        string `json:"small_head"`
	UserName         string `json:"user_name"`
}
