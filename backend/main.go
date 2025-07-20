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
	upgrader  websocket.Upgrader
	clients   = map[string]*websocket.Conn{} // userID → WS conn
	clientsMu sync.Mutex
)

func init() {
	// Load .env in local development; in production Render injects real env vars
	_ = godotenv.Load()

	frontendURL := os.Getenv("FRONTEND_URL") // e.g. "https://chatnest.onrender.com"
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Only allow our frontend origin
			return r.Header.Get("Origin") == frontendURL
		},
	}
}

func main() {
	// Connect to MongoDB
	config.ConnectDB()

	// Register REST API
	apiMux := http.NewServeMux()
	apiMux.Handle("/", routes.RegisterRoutes())

	// Register WebSocket endpoint under /ws
	apiMux.HandleFunc("/ws", wsHandler)

	// Wrap everything in CORS
	frontendURL := os.Getenv("FRONTEND_URL")
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(apiMux)

	// Determine port (Render provides $PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
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

	// Register user and broadcast updated online list
	clientsMu.Lock()
	clients[userID] = conn
	broadcast("getOnlineUsers", currentUsers())
	clientsMu.Unlock()

	// Clean up on disconnect
	defer func() {
		clientsMu.Lock()
		delete(clients, userID)
		broadcast("getOnlineUsers", currentUsers())
		clientsMu.Unlock()
		conn.Close()
	}()

	// Listen for incoming WS messages
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

	for _, c := range clients {
		_ = c.WriteJSON(map[string]interface{}{
			"event": event,
			"data":  payload,
		})
	}
}

func sendTo(userID, event string, payload interface{}) {
	clientsMu.Lock()
	c, ok := clients[userID]
	clientsMu.Unlock()

	if !ok {
		return
	}
	_ = c.WriteJSON(map[string]interface{}{
		"event": event,
		"data":  payload,
	})
}
