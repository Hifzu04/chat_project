package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	config "chat-backend/Config"
	routes "chat-backend/Routes"
	utils "chat-backend/Utils"
)

var (
	//upgrader: This object is responsible for "upgrading" a standard HTTP request to a WebSocket connection.
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func init() {
	_ = godotenv.Load()
}

func main() {
	//config.ConnectDB(): Calls the function in your local config package to connect to MongoDB.
	config.ConnectDB()
	//http.NewServeMux(): Creates a new HTTP Request Multiplexer (router).
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is awake!"))
	})
	//apiMux.Handle("/", ...): Any request starting with / (standard API calls like Login/Signup) is passed to the routes package.
	apiMux.Handle("/", routes.RegisterRoutes())
	//apiMux.HandleFunc("/ws", wsHandler): Any request to the specific endpoint /ws is handled by
	// the wsHandler function defined further down in this file.
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

	//log.Fatal(...): This wraps the server start.

	//  If the server crashes (e.g., port already in use), it logs the error and exits the program immediately.

	log.Fatal(http.ListenAndServe(":"+port, handler))

}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		//	log.Println(" wsHandler: missing userId in query")
		http.Error(w, "userId required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("wsHandler: Upgrade error:", err)
		return
	}
	//log.Printf(" [WS] connection established for userID=%s\n", userID)

	//pause other goroutine
	utils.ClientsMu.Lock()
	//Register: We map the userID to this specific conn.
	utils.Clients[userID] = conn
	utils.ClientsMu.Unlock()
	//log.Printf(" [WS] clients after register: %v\n", keys(utils.Clients))

	//currentUsers(): Gets a list of all User IDs currently in the map.
	//broadcast(...): Sends this list to everyone connected. This allows the frontend to show a green dot next to online friends immediately.
	users := utils.CurrentUsers()
	//log.Printf("[WS] broadcasting getOnlineUsers: %v\n", users)
	utils.Broadcast("getOnlineUsers", users)

	// Clean up on disconnect
	defer func() {
		utils.ClientsMu.Lock()
		delete(utils.Clients, userID)
		utils.ClientsMu.Unlock()
		after := utils.CurrentUsers()
		utils.Broadcast("getOnlineUsers", after)
		conn.Close()
	}()

	// Read loop
	for {
		var msg struct {
			Event string                 `json:"event"`
			Data  map[string]interface{} `json:"data"`
		}
		if err := conn.ReadJSON(&msg); err != nil {
			//	log.Printf("ReadJSON closed for %s: %v\n", userID, err)
			break
		}
		//log.Printf(" received WS event=%s data=%v\n", msg.Event, msg.Data)

		if msg.Event == "sendMessage" {
			to, ok := msg.Data["receiver_id"].(string)
			if !ok {
				log.Printf(" malformed sendMessage payload: %v\n", msg.Data)
				continue
			}
			//log.Printf(" sendMessage → to=%s payload=%v\n", to, msg.Data)
			//Routing Logic:
			//If the event is "sendMessage", it extracts the receiver_id.
			//It then calls sendTo to forward that message to the specific target user.
			utils.SendTo(to, "newMessage", msg.Data)
		}
	}
}
