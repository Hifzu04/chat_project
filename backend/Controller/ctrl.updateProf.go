// File: controllers/user.go

package controllers

import (
	config "chat-backend/Config"
	models "chat-backend/Models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateProfileRequest defines the JSON payload for updating user profile.
type UpdateProfileRequest struct {
	FullName   string `json:"fullname,omitempty"`
	ProfilePic string `json:"profile_pic,omitempty"`
}


//todo
//In cloudnary there will be a bucket to store all images  1:05:00





// UpdateProfile allows a logged-in user to update their full name or profile picture.
// The user ID is retrieved from the request context (set by JWT middleware).
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Extract userID from context
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "User ID missing in context", http.StatusUnauthorized)
		return
	}
	userID, ok := userIDValue.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Invalid user ID type", http.StatusInternalServerError)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	usersColl := config.GetCollection(models.CollectionNameUser)

	// Build update document
	update := bson.M{}
	if req.FullName != "" {
		update["fullname"] = req.FullName
	}
	if req.ProfilePic != "" {
		update["profile_pic"] = req.ProfilePic
	}
	if len(update) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": userID}
	updateDoc := bson.M{"$set": update}

	// Update the user document
	res := usersColl.FindOneAndUpdate(ctx, filter, updateDoc, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Return the updated user (excluding password)
	var updatedUser models.User
	err := res.Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Failed to decode updated user", http.StatusInternalServerError)
		return
	}
	updatedUser.Password = "" // ensure password is not sent back

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}
