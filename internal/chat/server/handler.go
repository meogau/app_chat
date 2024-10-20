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

func (handler *ChatHandler) HandleSocketConnection(ws *websocket.Conn, r *http.Request, _ int, _ int) error {
	// todo: create user
	queryParams := r.URL.Query()
	username := queryParams.Get("username")
	user := &socket.User{
		Username: username,
		Conn:     ws,
	}
	loginHandle(handler, user)
	go utils.RunWithRecovery(func() {
		handler.readWsMessage(user)
	})
	return nil
}

func (handler *ChatHandler) readWsMessage(user *socket.User) {
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Println("websocket closed")
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
			case "select_user":
				handleSelectUserMess(handler, user, data)
			case "message":
				handleSendMess(handler, data)
			}
		}
	}
}

func loginHandle(handler *ChatHandler, user *socket.User) {
	handler.InMemory.AddUSer(user)
	log.Printf("User %s logged in", user.Username)
	response := map[string]interface{}{
		"action": "login",
		"status": "success",
		"data": map[string]string{
			"message": "Login successful!",
		},
	}
	// todo: use send function
	err := user.SendMessage(response)
	if err != nil {
		log.Println("Error sending response:", err)
	}
}

func handleSelectUserMess(handler *ChatHandler, user *socket.User, data map[string]interface{}) {
	receiverUsername := data["receiver"].(string)
	// todo: get user
	receiverUser := handler.InMemory.GetUser(receiverUsername)
	if receiverUser == nil {
		response := map[string]interface{}{
			"action": "select_user",
			"status": "error",
			"data": map[string]string{
				"message": "User " + receiverUsername + " does not exist.",
			},
		}
		err := user.SendMessage(response)
		if err != nil {
			log.Println("Error sending response:", err)
		}
	}
	response := map[string]interface{}{
		"action": "select_user",
		"status": "success",
		"data": map[string]string{
			"message": "Selected user: " + receiverUsername,
		},
	}
	err := user.SendMessage(response)
	if err != nil {
		log.Println("Error sending response:", err)
	}
	log.Printf("%s has selected %s to chat", user.Username, receiverUsername)
}

func handleSendMess(handler *ChatHandler, data map[string]interface{}) {
	// todo: separate function
	senderUsername := data["username"].(string)
	receiverUsername := data["receiver"].(string)
	content := data["message"].(string)

	receiverUser := handler.InMemory.GetUser(receiverUsername)
	if receiverUser != nil {
		//todo: use send function instead
		err := receiverUser.Conn.WriteJSON(map[string]interface{}{
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
