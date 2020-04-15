package ws

import (
	"context"
	"io"

	"github.com/gorilla/websocket"
	proto "github.com/micro-community/x-push/proto/stream"
	"github.com/micro/go-micro/v2/util/log"
)

//ServeStream Server Stream fo websocket
func serveStream(cli proto.StreamerService, ws *websocket.Conn) error {
	//Read initial request from websocket

	var req proto.Request
	err := ws.ReadJSON(&req)
	if err != nil {
		return err
	}

	// Even if we aren't expecting further requests from the websocket, we still need to read from it to ensure we
	// get close signals
	go func() {
		for {
			if _, _, err := ws.NextReader(); err != nil {
				break
			}
		}
	}()

	log.Infof("Received Request: %v", req)
	// Send request to stream server
	stream, err := cli.ServerStream(context.Background(), &req)
	if err != nil {
		return err
	}
	defer stream.Close()

	// Read from the stream server and pass responses on to websocket
	for {
		// Read from stream, end request once the stream is closed
		rsp, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return err
			}

			break
		}

		// Write server response to the websocket
		err = ws.WriteJSON(rsp)
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
