package server

import (
	"app_chat/internal/chat/memory"
	"app_chat/pkg/model/socket"
	"app_chat/pkg/utils"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type ChatHandler struct {
	Upgrade  websocket.Upgrader
	InMemory memory.InMemory
}

func (handler *ChatHandler) HandleSocketConnection(ws *websocket.Conn, _ *http.Request, _ int, _ int) error {
	// todo: create user
	user := &socket.User{
		Username: "",
		Conn:     ws,
	}
	go utils.RunWithRecovery(func() {
		handler.readWsMessage(user)
	})
	return nil
}

func (handler *ChatHandler) readWsMessage(user *socket.User) {
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {

		}
	}(user.Conn)

	var message map[string]interface{}
	for {
		if err := user.Conn.ReadJSON(&message); err != nil {
			log.Println("Error reading json:", err)
			break
		}

		if action, ok := message["action"].(string); ok {
			data := message["data"].(map[string]interface{})

			switch action {
			case "login":
				username := data["username"].(string)
				// todo: check user name
				user := &socket.User{
					Username: username,
					Conn:     ws,
				}
				handler.InMemory.AddUSer(user)
				log.Printf("User %s logged in", username)
				response := map[string]interface{}{
					"action": "login",
					"status": "success",
					"data": map[string]string{
						"message": "Login successful!",
					},
				}

				// todo: use send function
				err := user.SendMessage(response)
				//err := ws.WriteJSON(response)
				if err != nil {
					log.Println("Error sending response:", err)
				}

			case "select_user":
				senderUsername := data["username"].(string)
				receiverUsername := data["receiver"].(string)
				// todo: get user
				if _, exists := srv.Clients[receiverUsername]; exists {
					response := map[string]interface{}{
						"action": "select_user",
						"status": "success",
						"data": map[string]string{
							"message": "Selected user: " + receiverUsername,
						},
					}
					err := ws.WriteJSON(response)
					if err != nil {
						log.Println("Error sending response:", err)
					}
					log.Printf("%s has selected %s to chat", senderUsername, receiverUsername)
				} else {
					response := map[string]interface{}{
						"action": "select_user",
						"status": "error",
						"data": map[string]string{
							"message": "User " + receiverUsername + " does not exist.",
						},
					}
					err := ws.WriteJSON(response)
					if err != nil {
						log.Println("Error sending response:", err)
					}
				}

			case "message":
				// todo: separate function
				senderUsername := data["username"].(string)
				receiverUsername := data["receiver"].(string)
				content := data["message"].(string)

				receiver, exists := srv.Clients[receiverUsername]
				if exists {
					//todo: use send function instead
					err := receiver.Conn.WriteJSON(map[string]interface{}{
						"action":  "message",
						"from":    senderUsername,
						"content": content,
					})
					if err != nil {
						log.Println("Error sending message to receiver:", err)
					} else {
						log.Printf("Message from %s to %s: %s", senderUsername, receiverUsername, content)
					}
				} else {
					log.Printf("Receiver %s does not exist or is not connected", receiverUsername)
				}
			}
		}
	}
}
