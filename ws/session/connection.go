package session

import (
	"github.com/gorilla/websocket"
	"sync"
)

//Connection Connection
type Connection struct {
	conn *websocket.Conn
	lock sync.Mutex
}

//WriteJSON ...
func (c *Connection) WriteJSON(v interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.conn.WriteJSON(v)
}
