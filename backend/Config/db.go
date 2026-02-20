//config database to read mongodb uri

package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DBClient *mongo.Client

	// Default database name if not set in .env or environment variables
	DBName = "chatdb"
)

func init() {
	// loads .env from project root (optional log if missing)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  no .env file found, relying on real ENV vars")
	}
} 
func ConnectDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		//log.Fatal(...): If the URI is missing, the program crashes immediately. Fatal prints the message and calls os.Exit(1) no further proceed.
		log.Fatal("mongodb uri is not set in the environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		//Fatalf is equivalent to [Printf] followed by a call to os.Exit(1).
       log.Fatalf("Mongo connect %v", err)
	}
    // Pinging (Verification)
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Mongo ping %v", err)
	}
	fmt.Println("Yeaa, connected to mongodb")
     //DBClient = client: most important.We take the successfully connected local client variable and assign it to the global DBClient variable we defined at the top. 
	 // Now the rest of the app can use the database.
	DBClient = client

	fmt.Printf("Using database: %s\n", DBName)

	
}


//Purpose: This is a utility function to make your code cleaner in other files.
//Instead of writing config.DBClient.Database("chatdb").Collection("users") every time you want to query users, you can just call: config.GetCollection("users").
//DBClient.Database(DBName): Selects the "chatdb" database.
//.Collection(name): Selects the specific collection (table) inside that database (e.g., "users", "messages").

func GetCollection(name string) *mongo.Collection {
	db := DBClient.Database(DBName)
	return db.Collection(name)

}
