package main

import (
	config "chat-backend/Config"
	routes "chat-backend/Routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	fmt.Println("welcome to backend")

	//load from env
	if err := godotenv.Load(); err != nil {
		log.Println("no env found")
	}

	//connect to mongodb , read mongodb uri
	config.ConnectDB()

	//router
	router := routes.RegisterRoutes()

	//handling cors error , when connecting with FE.

	c := cors.New(cors.Options{
		// Change this to your React dev origin(s)
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           600,
	})
	handler := c.Handler(router)

	//port
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	fmt.Printf("server is listening on : %s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
