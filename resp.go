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

type MsgVideoResp struct {
	Video   string `json:"video"`
	Message string `json:"message"`
	Size    int    `json:"size"`
	Status  int    `json:"status"`
}

type MsgVoiceResp struct {
	Voice   string `json:"voice"`
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

type CreateRoomResp struct {
	Message  string `json:"message"`
	Status   int    `json:"status"`
	UserName string `json:"user_name"`
}

type MsgAndStatus struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type QRCodeResp struct {
	Footer  string `json:"footer"`
	Message string `json:"message"`
	QRCode  string `json:"qr_code"`
	Status  int    `json:"status"`
}

type ImgResp struct {
	BigHead   string `json:"big_url"`
	SmallHead string `json:"small_url"`
	Status    int    `json:"status"`
	Size      int    `json:"size"`
	Message   string `json:"message"`
	Data      int    `json:"data"`
}

type Moment struct {
	CreateTime  int    `json:"create_time"`
	Description string `json:"description"`
	ID          string `json:"id"`
	NickName    string `json:"nick_name"`
	UserName    string `json:"user_name"`
}

type MomentResp struct {
	Data    Moment `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type MomentListResp struct {
	Bgi     string   `json:"bgi"`
	Data    []Moment `json:"data"`
	Message string   `json:"message"`
	Page    string   `json:"page"`
	Status  int      `json:"status"`
}

type MomentComment struct {
	CommentFlag   int    `json:"comment_flag"`
	Content       string `json:"content"`
	CreateTime    int    `json:"create_time"`
	DeleteFlag    int    `json:"delete_flag"`
	ID            int    `json:"id"`
	NickName      string `json:"nick_name"`
	ReplyID       int    `json:"reply_id"`
	ReplyUserName string `json:"reply_user_name"`
	Source        int    `json:"source"`
	Type          int    `json:"type"`
	UserName      string `json:"user_name"`
}

type MomentLike struct {
	Content    string `json:"content"`
	CreateTime int    `json:"create_time"`
	ID         int    `json:"id"`
	NickName   string `json:"nick_name"`
	Type       int    `json:"type"`
	UserName   string `json:"user_name"`
}

type MomentDetail struct {
	Comment     []MomentComment `json:"comment"`
	CreateTime  int             `json:"create_time"`
	Description string          `json:"description"`
	ID          string          `json:"id"`
	Like        []MomentLike    `json:"like"`
	NickName    string          `json:"nick_name"`
	UserName    string          `json:"user_name"`
}

type MomentDetailResp struct {
	Data    MomentDetail `json:"data"`
	Message string       `json:"message"`
	Status  int          `json:"status"`
}

type Fav struct {
	Flag int `json:"flag"`
	ID   int `json:"id"`
	Seq  int `json:"seq"`
	Time int `json:"time"`
	Type int `json:"type"`
}

type FavListResp struct {
	Continue int    `json:"continue"`
	Data     []Fav  `json:"data"`
	Key      string `json:"key"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
}

type AddFavResp struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Seq     int    `json:"seq"`
	Status  int    `json:"status"`
}

type FavDetail struct {
	Flag   int    `json:"flag"`
	ID     int    `json:"id"`
	Object string `json:"object"`
	Seq    int    `json:"seq"`
	Status int    `json:"status"`
	Time   int    `json:"time"`
}

type FavResp struct {
	Data    []FavDetail `json:"data"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
}

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LabelListResp struct {
	Label   []Label `json:"label"`
	Message string  `json:"message"`
	Status  int     `json:"status"`
}

type RedPacketResp struct {
	External string `json:"external"`
	Key      string `json:"key"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
}

type SearchMPResp struct {
	Code    int    `json:"code"`
	Info    string `json:"info"`
	Message string `json:"message"`
	Offset  int    `json:"offset"`
	Status  int    `json:"status"`
}

type RequestTokenResp struct {
	FullURL  string `json:"full_url"`
	Info     string `json:"info"`
	Message  string `json:"message"`
	ShareURL string `json:"share_url"`
	Status   int    `json:"status"`
}

type RequestUrlResp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	// 完整的访问结果原始数据文本（包含http头数据）
	Response string `json:"response"`
}
