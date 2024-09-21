package main

import (
	"app_chat/internal/chat/server"
	"app_chat/internal/socket_server"
)

func main() {
	chatHandler := &server.ChatHandler{}
	chatServer := socket_server.Server{
		Port:            8888,
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
