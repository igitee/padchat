package padchat

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
)

type Bot struct {
	sync.RWMutex
	ws            WSConn
	retProcMap    *sync.Map
	reqTimeout    time.Duration
	onQRURL       func(string)
	onScan        func(ScanResp)
	onMsg         func(Msg)
	onLogin       func()
	onLoaded      func()
	onContactSync func(Contact)
}

// NewBot 乃万物之始
// 新建 Bot 实例, 传入 PadChat 服务端地址
func NewBot(url string) (*Bot, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	err = conn.WriteJSON(struct {
		Type  string      `json:"type"`
		Cmd   string      `json:"cmd"`
		CmdID string      `json:"cmdId"`
		Data  interface{} `json:"data"`
	}{Type: "user", Cmd: "init", CmdID: uuid.New().String()})
	if err != nil {
		return nil, err
	}
	bot := newBot(conn)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 500)
		defer ticker.Stop()
		for range ticker.C {
			data := &ServerData{}
			conn.ReadJSON(data)
			switch data.Type {
			case "userEvent":
				bot.processUserEvent(data)
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
				fmt.Println(strings.Repeat("*", 100))
			}
		}
	}()
	return bot, nil
}

// OnClose ws 断开回调
func (bot *Bot) OnClose(f func(code int, text string) error) {
	bot.ws.SetCloseHandler(f)
}

func (bot *Bot) processUserEvent(data *ServerData) {
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
		var scan ScanResp
		json.Unmarshal(data.Data, &scan)
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
		jsoniter.Unmarshal(data.Data, push)
		for _, v := range push.List {
			msgType := jsoniter.Get(v, "msg_type").ToInt()
			switch msgType {
			case 5:
				var msg Msg
				jsoniter.Unmarshal(v, &msg)
				msg.MType = msg.SubType
				go func() {
					bot.RLock()
					defer bot.RUnlock()
					bot.onMsg(msg)
				}()
			case 2:
				var contact Contact
				jsoniter.Unmarshal(v, &contact)
				go func() {
					bot.RLock()
					defer bot.RUnlock()
					bot.onContactSync(contact)
				}()
			case 2048, 32768:
			default:
				fmt.Println(string(v))
				fmt.Println(strings.Repeat("=", 100))
			}
		}

	case "loaded":
		go func() {
			bot.RLock()
			defer bot.RUnlock()
			bot.onLoaded()
		}()
	default:
		fmt.Println(data.Event, string(data.Data))
		fmt.Println(strings.Repeat("~", 100))
	}
}

func newBot(conn *websocket.Conn) *Bot {
	return &Bot{
		RWMutex: sync.RWMutex{},
		ws: WSConn{
			Mutex: sync.Mutex{},
			Conn:  conn,
		},
		reqTimeout:    time.Second * 30,
		retProcMap:    &sync.Map{},
		onQRURL:       func(string) {},
		onScan:        func(ScanResp) {},
		onMsg:         func(Msg) {},
		onLogin:       func() {},
		onLoaded:      func() {},
		onContactSync: func(Contact) {},
	}
}

// OnQRURL 收到二维码回调, 需在执行二维码登录前配置
func (bot *Bot) OnQRURL(f func(string)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onQRURL = f
}

// OnScan 二维码扫描回调
func (bot *Bot) OnScan(f func(resp ScanResp)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onScan = f
}

// OnMsg 微信接收消息回调
func (bot *Bot) OnMsg(f func(msg Msg)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onMsg = f
}

// OnLogin 登录成功回调
func (bot *Bot) OnLogin(f func()) {
	bot.Lock()
	defer bot.Unlock()
	bot.onLogin = f
}

// OnLoaded 联系人加载完成回调
func (bot *Bot) OnLoaded(f func()) {
	bot.Lock()
	defer bot.Unlock()
	bot.onLoaded = f
}

// OnContactSync 联系人同步回调
func (bot *Bot) OnContactSync(f func(contact Contact)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onContactSync = f
}

// SetCommandTimeout 设置微信指令超时时间, 默认为 30 秒
func (bot *Bot) SetCommandTimeout(t time.Duration) {
	bot.reqTimeout = t
}
