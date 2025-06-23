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

	// 2) parse multipart
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Error(w, "could not parse form", http.StatusBadRequest)
		return
	}
	file, hdr, err := r.FormFile("profilePic")
	if err != nil {
		http.Error(w, "profilePic is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 3) upload
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
