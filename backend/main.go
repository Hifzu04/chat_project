package main

import (
	config "chat-backend/Config"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	///router := routes.RegisterRoutes()

	//port
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	fmt.Printf("server is listening on : %s\n", port)
	
	//log.Fatal(http.ListenAndServe(":"+port, router))
}
