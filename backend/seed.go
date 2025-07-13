// File: seed.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	config "chat-backend/Config"
	models "chat-backend/Models"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// seedUsers wipes & repopulates the users collection with test data.
func seedUsers() {
	col := config.GetCollection(models.CollectionNameUser)

	// 1) Clear existing docs
	if _, err := col.DeleteMany(context.Background(), bson.M{}); err != nil {
		log.Fatalf("seeder: could not clear users: %v", err)
	}

	// 2) Prepare raw data
	raw := []models.User{
		{Email: "emma.thompson@example.com", FullName: "Emma Thompson", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/women/1.jpg"},
		{Email: "james.anderson@example.com", FullName: "James Anderson", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/men/1.jpg"},
		{Email: "olivia.miller@example.com", FullName: "Olivia Miller", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/women/2.jpg"},
		{Email: "william.clark@example.com", FullName: "William Clark", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/men/2.jpg"},
		{Email: "sophia.davis@example.com", FullName: "Sophia Davis", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/women/3.jpg"},
		{Email: "benjamin.taylor@example.com", FullName: "Benjamin Taylor", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/men/3.jpg"},
		{Email: "ava.wilson@example.com", FullName: "Ava Wilson", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/women/4.jpg"},
		{Email: "lucas.moore@example.com", FullName: "Lucas Moore", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/men/4.jpg"},
		{Email: "isabella.brown@example.com", FullName: "Isabella Brown", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/women/5.jpg"},
		{Email: "henry.jackson@example.com", FullName: "Henry Jackson", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/men/5.jpg"},
		{Email: "mia.johnson@example.com", FullName: "Mia Johnson", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/women/6.jpg"},
		{Email: "alexander.martin@example.com", FullName: "Alexander Martin", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/men/6.jpg"},
		{Email: "charlotte.williams@example.com", FullName: "Charlotte Williams", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/women/7.jpg"},
		{Email: "daniel.rodriguez@example.com", FullName: "Daniel Rodriguez", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/men/7.jpg"},
		{Email: "amelia.garcia@example.com", FullName: "Amelia Garcia", Password: "123456", ProfilePic: "https://randomuser.me/api/portraits/women/8.jpg"},
	}

	// 3) Hash + timestamp + collect for insert
	var docs []interface{}
	for _, u := range raw {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("seeder: bcrypt error for %s: %v", u.Email, err)
		}
		u.Password = string(hash)
		u.CreatedAt = time.Now()
		docs = append(docs, u)
	}

	// 4) Insert all at once
	if _, err := col.InsertMany(context.Background(), docs); err != nil {
		log.Fatalf("seeder: insert error: %v", err)
	}

	fmt.Println("✅ User seed complete")
}

// init runs before main; if SEED_DB=true, seed and exit.
func init() {
	if os.Getenv("SEED_DB") == "true" {
		config.ConnectDB()
		seedUsers()
		os.Exit(0)
	}
}
