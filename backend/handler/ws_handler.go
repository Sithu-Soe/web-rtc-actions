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
	localSDP  = ""
	remoteSDP = ""
)

type WSReadMessage struct {
	MsgType int    `json:"msg_type"` // 1 for local addd, 2 for remote addd
	SDP     string `json:"sdp"`
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

			wsReadMsg := WSReadMessage{}
			if err := json.Unmarshal(data, &wsReadMsg); err != nil {
				log.Println(err)
				return
			}

			if wsReadMsg.MsgType == 1 {
				localSDP = wsReadMsg.SDP
			}

			if wsReadMsg.MsgType == 2 {
				remoteSDP = wsReadMsg.SDP
			}

			log.Println("LOCAL:", localSDP, "REMOTE:", remoteSDP)

		}
	}()

}
