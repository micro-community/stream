package session

import (
	"github.com/gorilla/websocket"
)

//RemoveClient ...
func RemoveClient(userID string, conn *websocket.Conn) {
	lock.Lock()
	client, ok := clients[userID]

	if !ok {
		lock.Unlock()
		return
	}

	client.delete(conn)
	if len(client.conns) == 0 {
		delete(clients, userID)
	}
	lock.Unlock()
}

//AddClient ...
func AddClient(userID string, conn *websocket.Conn) {
	lock.Lock()
	client, ok := clients[userID]
	if !ok {
		client = &Client{
			userID: userID,
			conns:  make(map[*websocket.Conn]*Connection),
		}
		clients[userID] = client
	}
	client.add(conn)
	lock.Unlock()
}
