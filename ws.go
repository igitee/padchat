package padchat

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/json-iterator/go"
)

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
