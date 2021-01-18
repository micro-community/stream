package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/micro-community/streaming/ws/session"
	"github.com/micro/go-micro/v2/util/log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

//verifyAuth token from head
func verifyAuth(tokenStr string) (string, error) {
	//	log.Info(token,check)
	// todo ..
	//send token server to verify auth .

	return "userID", nil
}

//HandleConn of websocket
func HandleConn(w http.ResponseWriter, r *http.Request) {
	// Upgrade request to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal("Upgrade: ", err)
		return
	}
	defer conn.Close()

	var userID string
	////to do
	//	userID = verifyAuth()
	///

	session.AddClient(userID, conn)

	session.RemoveClient(userID, conn)
	log.Infof("Stream complete")
}
