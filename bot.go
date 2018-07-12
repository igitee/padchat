package padchat

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Bot struct {
	sync.RWMutex
	ws         *websocket.Conn
	retProcMap *sync.Map
	onQRURL    func(string)
	onScan     func(*ScanResp)
	onMsg      func([]Msg)
	onLogin    func()
}

type WSReq struct {
	Type  string      `json:"type"`
	CMD   string      `json:"cmd"`
	CMDID string      `json:"cmdId"`
	Data  interface{} `json:"data"`
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

func NewBot(url string) (*Bot, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	err = conn.WriteJSON(struct {
		Type  string      `json:"type"`
		CMD   string      `json:"cmd"`
		CMDID string      `json:"cmdId"`
		Data  interface{} `json:"data"`
	}{Type: "user", CMD: "init", CMDID: uuid.New().String()})
	if err != nil {
		return nil, err
	}
	bot := newBot(conn)
	go func() {
		for range time.Tick(time.Millisecond * 500) {
			data := &struct {
				Type   string
				Event  string
				TaskID string
				CMDID  string `json:"cmdId"`
				Data   json.RawMessage
				Msg    string
			}{}
			conn.ReadJSON(data)
			switch data.Type {
			case "userEvent":
				switch data.Event {
				case "qrcode":
					url := &struct {
						URL string
					}{}
					json.Unmarshal(data.Data, url)
					go func() {
						bot.RLock()
						defer bot.RUnlock()
						bot.onQRURL(url.URL)
					}()
				case "scan":
					scan := &ScanResp{}
					json.Unmarshal(data.Data, scan)
					go func() {
						bot.RLock()
						defer bot.RUnlock()
						bot.onScan(scan)
					}()
				case "login":
					go func() {
						bot.RLock()
						defer bot.RUnlock()
						bot.onLogin()
					}()
				case "push":
					push := &PushResp{}
					json.Unmarshal(data.Data, push)
					go func() {
						bot.RLock()
						defer bot.RUnlock()
						bot.onMsg(push.List)
					}()
				default:
					fmt.Println(data.Event, string(data.Data))
					fmt.Println("=============================")
				}
			case "cmdRet":
				proc, ok := bot.retProcMap.Load(data.CMDID)
				if ok {
					var resp CommandResp
					json.Unmarshal(data.Data, &resp)
					proc.(func(CommandResp))(resp)
				}
				bot.retProcMap.Delete(data.CMDID)
			default:
				fmt.Println(data.Type, data.Event, string(data.Data))
			}
		}
	}()
	return bot, nil
}

func newBot(conn *websocket.Conn) *Bot {
	return &Bot{
		RWMutex:    sync.RWMutex{},
		ws:         conn,
		retProcMap: &sync.Map{},
		onQRURL:    func(string) {},
		onScan:     func(*ScanResp) {},
		onMsg:      func([]Msg) {},
		onLogin:    func() {},
	}
}

func (b *Bot) OnQRURL(f func(string)) {
	b.Lock()
	defer b.Unlock()
	b.onQRURL = f
}

func (b *Bot) OnScan(f func(resp *ScanResp)) {
	b.Lock()
	defer b.Unlock()
	b.onScan = f
}

func (b *Bot) OnMsg(f func(msgList []Msg)) {
	b.Lock()
	defer b.Unlock()
	b.onMsg = f
}

func (b *Bot) OnLogin(f func()) {
	b.Lock()
	defer b.Unlock()
	b.onLogin = f
}

type CommandResp struct {
	Success bool
	Data    json.RawMessage
	Msg     string
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

func (b *Bot) sendCommand(cmd string, data interface{}) CommandResp {
	c := make(chan CommandResp)
	id := uuid.New().String()
	b.retProcMap.Store(id, func(resp CommandResp) {
		c <- resp
	})
	b.ws.WriteJSON(WSReq{Type: "user", CMD: cmd, CMDID: id, Data: data})
	return <-c
}

func (b *Bot) Init() CommandResp {
	return b.sendCommand("init", nil)
}

func (b *Bot) QRLogin() CommandResp {
	return b.sendCommand("login", LoginReq{LoginType: "qrcode"})
}

func (b *Bot) RequestLogin(wxData, token string) CommandResp {
	return b.sendCommand("login", LoginReq{
		LoginType: "request",
		WXData:    wxData,
		Token:     token,
	})
}

func (b *Bot) TokenLogin(wxData, token string) CommandResp {
	return b.sendCommand("login", LoginReq{
		LoginType: "token",
		WXData:    wxData,
		Token:     token,
	})
}

func (b *Bot) UserLogin(wxData, username, password string) CommandResp {
	return b.sendCommand("login", LoginReq{
		LoginType: "user",
		WXData:    wxData,
		UserName:  username,
		Password:  password,
	})
}

func (b *Bot) PhoneLogin(wxData, phone, code string) CommandResp {
	return b.sendCommand("login", LoginReq{
		LoginType: "phone",
		WXData:    wxData,
		Phone:     phone,
		Code:      code,
	})
}

func (b *Bot) GetWXData() (string, error) {
	data := &struct {
		WXData string `json:"wx_data"`
	}{}
	resp := b.sendCommand("getWxData", nil)
	if !resp.Success {
		return "", errors.New(resp.Msg)
	}
	err := json.Unmarshal(resp.Data, data)
	if err != nil {
		return "", err
	}
	return data.WXData, nil
}

type LoginTokenResp struct {
	Token string `json:"token"`
	Uin   int64  `json:"uin"`
}

func (b *Bot) GetLoginToken() (*LoginTokenResp, error) {
	data := &LoginTokenResp{}
	resp := b.sendCommand("getLoginToken", nil)
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	err := json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type MyInfoResp struct {
	UserName string `json:"userName"`
	Uin      int64  `json:"uin"`
}

func (b *Bot) GetMyInfo() (*MyInfoResp, error) {
	data := &MyInfoResp{}
	resp := b.sendCommand("getMyInfo", nil)
	if !resp.Success {
		return nil, errors.New(resp.Msg)
	}
	err := json.Unmarshal(resp.Data, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b *Bot) SyncContact() CommandResp {
	return b.sendCommand("syncContact", nil)
}
