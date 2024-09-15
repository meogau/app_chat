package server

import (
	"app_chat/pkg/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatHandler struct {
}

func (handler *ChatHandler) HandleSocketConnection(ws *websocket.Conn, r *http.Request, readTimeout int, writeTimeout int) error {
	go utils.RunWithRecovery(func() {
		handler.readWsMessage(ws)
	})
	return nil
}

func (handler *ChatHandler) readWsMessage(conn *websocket.Conn) {
	for {
		_, mess, err := conn.ReadMessage()
		if err != nil {
			expectedCloseErrorCodes := []int{websocket.CloseGoingAway, websocket.CloseNormalClosure}
			if websocket.IsUnexpectedCloseError(err, expectedCloseErrorCodes...) {
				fmt.Printf("Error read message IsUnexpectedCloseError %v", err)
			} else {
				fmt.Printf("Error read message %v", err)
			}
		}
		messString := string(mess)
		messString += " replied"
		err = conn.WriteMessage(1, []byte(messString))
		if err != nil {
			return
		}
	}
}
