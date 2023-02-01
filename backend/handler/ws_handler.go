package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	localOffer         = SDPDescription{}
	remoteIceCandidate = SDPDescription{}
)

type WSReadMessage struct {
	MsgType        int `json:"msg_type"` // 1 for local addd, 2 for remote addd
	SDPDescription `json:"sdp_description"`
}

type SDPDescription struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

func ServeWS(c *gin.Context) {
	conn, err := upgradeConnection.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			msgType, data, err := conn.ReadMessage()
			if err != nil {
				// someone let the page
				log.Println(err.Error(), msgType)
				// if someone leave without proper close, second condition ("IsUnexpectedCloseError") will be true
				// Or timeout
				if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) || strings.Contains(err.Error(), "i/o timeout") {
					conn.Close()
					return
				}
				continue
			}

			// log.Println(string(data))

			wsReadMsg := WSReadMessage{}
			if err := json.Unmarshal(data, &wsReadMsg); err != nil {
				log.Println(err)
				return
			}

			// log.Println(wsReadMsg)

			if wsReadMsg.MsgType == 1 {
				localOffer = wsReadMsg.SDPDescription
			}

			// if wsReadMsg.MsgType == 2 {
			// 	remoteIceCandidate = wsReadMsg.IceCandidate
			// }

			// log.Println("LOCAL:", localOffer, "REMOTE:", remoteIceCandidate)

		}
	}()

}
