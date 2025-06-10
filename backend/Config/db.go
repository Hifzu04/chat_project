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
		log.Fatal("mongodb uri is not set in the environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOpts)

	if err != nil {
		log.Fatalf("Mongo connect %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Mongo ping %v", err)
	}
	fmt.Println("Yeaa, connected to mongodb")

	DBClient = client
	fmt.Printf("Using database: %s\n", DBName)
}

func GetCollection(name string) *mongo.Collection {
	db := DBClient.Database(DBName)
	return db.Collection(name)

}
