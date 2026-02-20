// File: routes/routes.go

package routes

import (
	"net/http"

	controllers "chat-backend/Controller"
	middleware "chat-backend/Middleware"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all application routes and returns the mux.Router instance.
//
// Public routes:

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
	//mux.NewRouter(): Creates a new, empty router. This is the main object that will hold all your rules.
	router := mux.NewRouter()

	// Public auth routes
	//HandleFunc("/signup", ...): Tells the router: "If someone visits /signup, run the controllers.Signup function."
	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/logout", controllers.Logout).Methods("POST")

	// Create a subrouter for protected endpoints . when these routes pass we move next.
	authRouter := router.PathPrefix("/").Subrouter()
	//Use(...): This applies your Authenticate middleware to every single route attached to authRouter
	authRouter.Use(middleware.Authenticate)

	// User-related

	//update prof
	authRouter.HandleFunc("/user/update", controllers.UpdateProfile).Methods("PUT")
	//check auth
	authRouter.HandleFunc("/auth/check", controllers.CheckAuth).Methods("GET")

	// // Message-related

	//user for sidebar
	authRouter.HandleFunc("/users", controllers.GetAllUsersForSidebar).Methods("GET")

	authRouter.HandleFunc("/messages/send", controllers.SendMessage).Methods("POST")

	//The router extracts the value (e.g., "123") and makes it available to your controller via mux.Vars(r).
	//This is how your GetMessages controller knows which friend's chat history to load.
	authRouter.HandleFunc("/messages/{userID}", controllers.GetMessages).Methods("GET")

	// A simple health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"up"}`))
	}).Methods("GET")

	return router
}



/*   Operation   SQL         http method     MongoDB
--------------------------------------------------------------------
C    Create      INSERT,     POST,           InsertOne/InsertMany
R    Read        SELECT,     GET,            Find
U    Update      UPDATE,     PUT/PATCH,      UpdateOne/updateMany
D    Delete      DELETE,     DELETE,         DeleteOne/DeleteMany
*/
