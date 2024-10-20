package main

import (
	"app_chat/internal/chat/server"
	"app_chat/internal/socket_server"
	"github.com/gorilla/websocket"
)

func main() {
	chatHandler := &server.ChatHandler{
		Upgrade: websocket.Upgrader{},
	}
	chatServer := socket_server.Server{
		Port:            8989,
		Name:            "Chat",
		ReadTimeOut:     1000,
		WriteTimeOut:    1000,
		MaxMessageSize:  1000,
		ReadBufferSize:  1000,
		WriteBufferSize: 1000,
		Handler:         chatHandler,
	}
	err := chatServer.Start()
	if err != nil {
		return
	}
}
