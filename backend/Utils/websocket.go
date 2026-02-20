package utils

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// We moved these here from main.go and capitalized them to make them public
	//Client serve as db for currently  conneected user
	Clients   = make(map[string]*websocket.Conn)
	ClientsMu sync.Mutex
)

// SendTo sends a WebSocket event to a specific user
func SendTo(userID string, event string, payload interface{}) {
	ClientsMu.Lock()
	c, ok := Clients[userID]
	ClientsMu.Unlock()

	if !ok {
		// User is offline, they will get the message when they refresh
		return
	}

	if err := c.WriteJSON(map[string]interface{}{
		"event": event,
		"data":  payload,
	}); err != nil {
		log.Printf("SendTo write error for %s: %v\n", userID, err)
	}
}

// Broadcast sends a WebSocket event to all connected users
func Broadcast(event string, payload interface{}) {
	ClientsMu.Lock()
	defer ClientsMu.Unlock()
	for uid, c := range Clients {
		if err := c.WriteJSON(map[string]interface{}{
			"event": event,
			"data":  payload,
		}); err != nil {
			log.Printf("write to %s error: %v\n", uid, err)
		}
	}
}

// CurrentUsers returns a list of all online user IDs
func CurrentUsers() []string {
	ClientsMu.Lock()
	defer ClientsMu.Unlock()
	ids := make([]string, 0, len(Clients))
	for id := range Clients {
		ids = append(ids, id)
	}
	return ids
}
