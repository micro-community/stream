package jessica

import (
	"encoding/binary"
	"net/http"
	"strings"

	"github.com/gobwas/ws"
	"github.com/micro-community/streaming/engine"
	"github.com/micro-community/streaming/engine/avformat"
	"github.com/micro-community/streaming/engine/pool"
)

//WsHandler for websocket
func WsHandler(w http.ResponseWriter, r *http.Request) {
	sign := r.URL.Query().Get("sign")
	isFlv := false
	if err := engine.AuthHooks.Trigger(sign); err != nil {
		w.WriteHeader(403)
		return
	}
	streamPath := strings.TrimLeft(r.RequestURI, "/")
	if strings.HasSuffix(streamPath, ".flv") {
		streamPath = strings.TrimRight(streamPath, ".flv")
		isFlv = true
	}
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	baseStream := engine.Subscriber{Sign: sign}
	baseStream.ID = conn.RemoteAddr().String()
	defer conn.Close()
	if isFlv {
		baseStream.Type = "JessicaFlv"
		baseStream.OnData = func(packet *avformat.SendPacket) error {
			return avformat.WriteFLVTag(conn, packet)
		}
		if err := ws.WriteHeader(conn, ws.Header{
			Fin:    true,
			OpCode: ws.OpBinary,
			Length: int64(13),
		}); err != nil {
			return
		}
		if _, err = conn.Write(avformat.FLVHeader); err != nil {
			return
		}
	} else {
		baseStream.Type = "Jessica"
		baseStream.OnData = func(packet *avformat.SendPacket) error {
			err := ws.WriteHeader(conn, ws.Header{
				Fin:    true,
				OpCode: ws.OpBinary,
				Length: int64(len(packet.Payload) + 5),
			})
			if err != nil {
				return err
			}
			head := pool.GetSlice(5)
			head[0] = packet.Type - 7
			binary.BigEndian.PutUint32(head[1:5], packet.Timestamp)
			if _, err = conn.Write(head); err != nil {
				return err
			}
			pool.RecycleSlice(head)
			//if p.Header[0] == 2 {
			//	fmt.Printf("%6d %X\n", (uint32(p.Packet.Payload[5])<<24)|(uint32(p.Packet.Payload[6])<<16)|(uint32(p.Packet.Payload[7])<<8)|uint32(p.Packet.Payload[8]), p.Packet.Payload[9])
			//}
			if _, err = conn.Write(packet.Payload); err != nil {
				return err
			}
			return nil
		}
	}
	baseStream.Publish(streamPath)
}
