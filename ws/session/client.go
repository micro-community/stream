package session

import (
	"sync"

	"log"

	"github.com/gorilla/websocket"
)

// Client Connection
type Client struct {
	userID string
	conns  map[*websocket.Conn]*Connection
}

func (c *Client) add(conn *websocket.Conn) {

	_, ok := c.conns[conn]
	if ok {
		return
	}
	c.conns[conn] = &Connection{
		conn: conn,
	}
}

func (c *Client) delete(conn *websocket.Conn) {
	delete(c.conns, conn)
}

func (c *Client) sendMessage(msg interface{}) {
	for _, conn := range c.conns {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("conn.WriteJSON error:%v", err)
		}
	}
}

// parameters
var (
	clients = make(map[string]*Client)
	lock    sync.RWMutex
)

// SendMessage To User..
func SendMessage(msg interface{}, userIDs []string) {
	lock.RLock()
	for _, userID := range userIDs {
		client, exist := clients[userID]
		if exist {
			log.Printf("streaming log (data:%s), user(%s)", msg, userID)
			client.sendMessage(msg)
		}
	}
	lock.RUnlock()
}
