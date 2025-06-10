// File: routes/routes.go

package routes

import (
	"net/http"

	middleware "chat-backend/Middleware"
	controllers "chat-backend/Controller"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all application routes and returns the mux.Router instance.
//
// Public routes:
//
//	POST /signup      → Signup
//	POST /login       → Login
//	POST /logout      → Logout
//
// Protected routes (require valid JWT):
//
//	PUT  /user/update         → UpdateProfile
//	POST /messages/send       → SendMessage
//	GET  /messages/{userID}   → GetMessages
func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	// Public auth routes
	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/logout", controllers.Logout).Methods("POST")

	// Create a subrouter for protected endpoints
	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.Authenticate)

	// // User-related
	// authRouter.HandleFunc("/user/update", controllers.UpdateProfile).Methods("PUT")

	// // Message-related
	// authRouter.HandleFunc("/messages/send", controllers.SendMessage).Methods("POST")
	// authRouter.HandleFunc("/messages/{userID}", controllers.GetMessages).Methods("GET")

	// A simple health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"up"}`))
	}).Methods("GET")

	return router
}
