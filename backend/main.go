package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	config "chat-backend/Config"
	routes "chat-backend/Routes"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	clients   = map[string]*websocket.Conn{} // userID → WS conn
	clientsMu sync.Mutex
)

func init() {
	_ = godotenv.Load()
}

func main() {
	config.ConnectDB()

	apiMux := http.NewServeMux()
	apiMux.Handle("/", routes.RegisterRoutes())
	apiMux.HandleFunc("/ws", wsHandler)

	frontendURL := os.Getenv("FRONTEND_URL")
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(apiMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Printf("Server listening on port %s\n", port)
	fmt.Println("FRONTEND_URL from env:", frontendURL)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		log.Println("❌ wsHandler: missing userId in query")
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ wsHandler: Upgrade error:", err)
		return
	}
	log.Printf("🔴 [WS] connection established for userID=%s\n", userID)

	// Register
	clientsMu.Lock()
	clients[userID] = conn
	clientsMu.Unlock()
	log.Printf("🔴 [WS] clients after register: %v\n", keys(clients))

	// Broadcast getOnlineUsers
	users := currentUsers()
	log.Printf("🔴 [WS] broadcasting getOnlineUsers: %v\n", users)
	broadcast("getOnlineUsers", users)

	// Clean up on disconnect
	defer func() {
		clientsMu.Lock()
		delete(clients, userID)
		clientsMu.Unlock()
		after := currentUsers()
		log.Printf("🔴 [WS] after unregister, currentUsers: %v\n", after)
		broadcast("getOnlineUsers", after)
		conn.Close()
	}()

	// Read loop
	for {
		var msg struct {
			Event string                 `json:"event"`
			Data  map[string]interface{} `json:"data"`
		}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("⚠️ ReadJSON closed for %s: %v\n", userID, err)
			break
		}
		log.Printf("📨 received WS event=%s data=%v\n", msg.Event, msg.Data)

		if msg.Event == "sendMessage" {
			to, ok := msg.Data["receiver_id"].(string)
			if !ok {
				log.Printf("❌ malformed sendMessage payload: %v\n", msg.Data)
				continue
			}
			log.Printf("🔴 sendMessage → to=%s payload=%v\n", to, msg.Data)
			sendTo(to, "newMessage", msg.Data)
		}
	}
}

func currentUsers() []string {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	ids := make([]string, 0, len(clients))
	for id := range clients {
		ids = append(ids, id)
	}
	return ids
}

func broadcast(event string, payload interface{}) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for uid, c := range clients {
		if err := c.WriteJSON(map[string]interface{}{
			"event": event,
			"data":  payload,
		}); err != nil {
			log.Printf("❌ write to %s error: %v\n", uid, err)
		}
	}
}

func sendTo(userID, event string, payload interface{}) {
	clientsMu.Lock()
	c, ok := clients[userID]
	clientsMu.Unlock()
	if !ok {
		log.Printf("⚠️ sendTo: no client for userID=%s\n", userID)
		return
	}
	if err := c.WriteJSON(map[string]interface{}{
		"event": event,
		"data":  payload,
	}); err != nil {
		log.Printf("❌ sendTo write error for %s: %v\n", userID, err)
	}
}

func keys(m map[string]*websocket.Conn) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
