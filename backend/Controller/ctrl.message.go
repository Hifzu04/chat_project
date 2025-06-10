// File: controllers/user.go

package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	config "chat-backend/Config"
	models "chat-backend/Models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserListItem is the subset of User fields we expose in the sidebar.
type UserListItem struct {
	ID         string `json:"id"`
	FullName   string `json:"fullname"`
	ProfilePic string `json:"profile_pic"`
	Email      string `json:"email"`
}

// GetAllUsers returns a list of all registered users (excluding passwords).
//only if the user is logged in . in routes GET /users (protected by your Authenticate middleware) returns:
// Protected by JWT middleware.
func GetAllUsersForSidebar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1) Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2) Query the users collection, projecting out the password field
	coll := config.GetCollection(models.CollectionNameUser)
	cursor, err := coll.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{
		"password": 0,
	}))
	if err != nil {
		log.Println("❌ GetAllUsers Find error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// 3) Decode into slice of UserListItem
	var list []UserListItem
	for cursor.Next(ctx) {
		var u models.User
		if err := cursor.Decode(&u); err != nil {
			log.Println("❌ GetAllUsers Decode error:", err)
			http.Error(w, "Data error", http.StatusInternalServerError)
			return
		}
		list = append(list, UserListItem{
			ID:         u.ID.Hex(),
			FullName:   u.FullName,
			ProfilePic: u.ProfilePic,
			Email:      u.Email,
		})
	}
	if err := cursor.Err(); err != nil {
		log.Println("❌ GetAllUsers Cursor error:", err)
		http.Error(w, "Cursor error", http.StatusInternalServerError)
		return
	}

	// 4) Return the JSON list
	json.NewEncoder(w).Encode(list)
}
