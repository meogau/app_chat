package socket_server

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler interface {
	HandleSocketConnection(ws *websocket.Conn, r *http.Request, readTimeout int, writeTimeout int) error
}
