package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type User struct {
	Username string
	Conn     *websocket.Conn
	wsLock   sync.Mutex
}

func (u *User) SendMessage(msg interface{}) error {
	u.wsLock.Lock()
	defer u.wsLock.Unlock()
	err := u.Conn.WriteJSON(msg)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("Send message to streamer failed %+v, err: %+v", u.Username, err)
		}
		return err
	}
	return nil
}
