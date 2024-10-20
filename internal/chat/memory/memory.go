package memory

import (
	"app_chat/pkg/model/socket"
	"sync"
)

type InMemory struct {
	data map[string]*socket.User
	sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		data: make(map[string]*socket.User),
	}
}

func (m *InMemory) AddUSer(user *socket.User) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.data[user.Username] = user
}

func (m *InMemory) GetUser(username string) *socket.User {
	m.RWMutex.RLock()
	defer m.RWMutex.Unlock()
	user, ok := m.data[username]
	if !ok {
		return nil
	}
	return user
}
