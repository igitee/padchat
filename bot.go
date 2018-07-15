package padchat

import (
	"encoding/json"
	"fmt"
	"math"
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
	onQRURL       func(string)
	onScan        func(ScanResp)
	onMsg         func(Msg)
	onLogin       func()
	onLoaded      func()
	onContactSync func(Contact)
}

type WSConn struct {
	sync.Mutex
	*websocket.Conn
}

func (c *WSConn) WriteJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	encoder := jsoniter.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	err1 := encoder.Encode(v)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

//NewBot create new Bot instance
func NewBot(url string) (*Bot, error) {
	dialer := websocket.DefaultDialer
	dialer.WriteBufferSize = math.MaxInt32
	conn, _, err := dialer.Dial(url, nil)
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
		retProcMap:    &sync.Map{},
		onQRURL:       func(string) {},
		onScan:        func(ScanResp) {},
		onMsg:         func(Msg) {},
		onLogin:       func() {},
		onLoaded:      func() {},
		onContactSync: func(Contact) {},
	}
}

func (bot *Bot) OnQRURL(f func(string)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onQRURL = f
}

func (bot *Bot) OnScan(f func(resp ScanResp)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onScan = f
}

func (bot *Bot) OnMsg(f func(msg Msg)) {
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

func (bot *Bot) OnContactSync(f func(contact Contact)) {
	bot.Lock()
	defer bot.Unlock()
	bot.onContactSync = f
}
