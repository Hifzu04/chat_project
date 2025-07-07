package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	config "chat-backend/Config"
	routes "chat-backend/Routes"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// allow requests from your React dev server
			return r.Header.Get("Origin") == "http://localhost:5173"
		},
	}
	clients   = map[string]*websocket.Conn{} // userID → WS conn
	clientsMu sync.Mutex
)

func main() {
	// 1) Load env & connect to MongoDB
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file; using environment vars")
	}
	config.ConnectDB()

	// 2) Register REST API under /
	apiMux := http.NewServeMux()
	apiMux.Handle("/", routes.RegisterRoutes())

	// 3) Register WebSocket endpoint under /ws
	apiMux.HandleFunc("/ws", wsHandler)

	// 4) Wrap everything in CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		// EnableDebug:    true, // uncomment to log CORS decisions
	}).Handler(apiMux)

	// 5) Start server
	fmt.Println("Server listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}

	// Register and broadcast online users
	clientsMu.Lock()
	clients[userID] = conn
	broadcast("getOnlineUsers", currentUsers())
	clientsMu.Unlock()

	// Clean up when they disconnect
	defer func() {
		clientsMu.Lock()
		delete(clients, userID)
		broadcast("getOnlineUsers", currentUsers())
		clientsMu.Unlock()
		conn.Close()
	}()

	// Read incoming messages
	for {
		var msg struct {
			Event string                 `json:"event"`
			Data  map[string]interface{} `json:"data"`
		}
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		if msg.Event == "sendMessage" {
			to := msg.Data["receiver_id"].(string)
			sendTo(to, "newMessage", msg.Data)
		}
	}
}

func currentUsers() []string {
	ks := make([]string, 0, len(clients))
	for id := range clients {
		ks = append(ks, id)
	}
	return ks
}

func broadcast(event string, payload interface{}) {
	for _, c := range clients {
		_ = c.WriteJSON(map[string]interface{}{
			"event": event,
			"data":  payload,
		})
	}
}

func sendTo(userID, event string, payload interface{}) {
	if c, ok := clients[userID]; ok {
		_ = c.WriteJSON(map[string]interface{}{
			"event": event,
			"data":  payload,
		})
	}
}
