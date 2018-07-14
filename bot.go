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
	ws         *websocket.Conn
	retProcMap *sync.Map
	onQRURL    func(string)
	onScan     func(*ScanResp)
	onMsg      func(Msg)
	onLogin    func()
	onLoaded   func()
}

//NewBot create new Bot instance
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
		jsoniter.Unmarshal(data.Data, push)
		for _, v := range push.List {
			msgType := jsoniter.Get(v, "msg_type").ToInt()
			switch msgType {
			case 5:
				msg := &Msg{}
				json.Unmarshal(v, msg)
				msg.MType = msg.SubType
				go func() {
					bot.RLock()
					defer bot.RUnlock()
					bot.onMsg(*msg)
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
		RWMutex:    sync.RWMutex{},
		ws:         conn,
		retProcMap: &sync.Map{},
		onQRURL:    func(string) {},
		onScan:     func(*ScanResp) {},
		onMsg:      func(Msg) {},
		onLogin:    func() {},
		onLoaded:   func() {},
	}
}

func (bot *Bot) OnQRURL(f func(string)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onQRURL = f
}

func (bot *Bot) OnScan(f func(resp *ScanResp)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onScan = f
}

func (bot *Bot) OnMsg(f func(msgList Msg)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onMsg = f
}

func (bot *Bot) OnLogin(f func()) {
	bot.Lock()
	defer bot.Unlock()
	bot.onLogin = f
}

func (bot *Bot) OnLoaded(f func()) {
	bot.Lock()
	defer bot.Unlock()
	bot.onLoaded = f
}
