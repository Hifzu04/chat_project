// File: controllers/user.go

package controllers

import (
	config "chat-backend/Config"
	models "chat-backend/Models"
	utils "chat-backend/Utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// 1) pull ObjectID straight out of context
	oid, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2)  ParseMultipartForm: Standard JSON requests cannot carry image data. We use Multipart Form Data.
	//5 << 20: This is a bitwise operation that equals 5 MB (5 * 1024 * 1024).
	if err := r.ParseMultipartForm(5 << 20); err != nil {

		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}

	// r.FormFile("profilePic"): Looks for the specific file input named "profilePic" coming from the React frontend.
	file, hdr, err := r.FormFile("profilePic")
	if err != nil {
		http.Error(w, "profilePic is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 3) upload to cloudinary
	secureURL, err := utils.UploadToCloudinary(file, hdr)
	if err != nil {
		http.Error(w, "upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4) update Mongo using the same "profile_pic" key
	usersColl := config.GetCollection(models.CollectionNameUser)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()


	
	filter := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"profile_pic": secureURL}}

	//SetReturnDocument(options.After): This is the most critical line in this block.
	//Default Behavior: MongoDB updates the document but returns the OLD version (the one before the update).
	//options.After: Tells MongoDB: "Update the document, and then give me back the NEW version."
	//Why? We want to send the new profile picture back to the frontend immediately so the UI updates instantly.
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser models.User
	err = usersColl.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, "db update error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 5) return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}
