package ws

import (
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/v2/util/log"
)

//ServeStream Server Stream fo websocket
func serveStream(ws *websocket.Conn) error {

	// Even if we aren't expecting further requests from the websocket, we still need to read from it to ensure we
	// get close signals
	go func() {
		for {
			if _, _, err := ws.NextReader(); err != nil {
				break
			}
		}
	}()
	log.Info("Received Request")
	//	log.Infof("Received Request: %v", req)

	for {

		var rsp interface{}
		// Write server response to the websocket
		err := ws.WriteJSON(rsp)
		if err != nil {
			// End request if socket is closed
			if isExpectedClose(err) {
				log.Infof("Expected Close on socket", err)
				break
			} else {
				return err
			}
		}
	}

	return nil
}

func isExpectedClose(err error) bool {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
		log.Infof("Unexpected websocket close: ", err)
		return false
	}

	return true
}
