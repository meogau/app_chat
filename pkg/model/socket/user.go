package socket

import "github.com/gorilla/websocket"

type User struct {
	Username string
	Conn     *websocket.Conn
}

type Server struct {
	Clients map[string]*User
}

func NewServer() *Server {
	return &Server{
		Clients: make(map[string]*User),
	}
}
