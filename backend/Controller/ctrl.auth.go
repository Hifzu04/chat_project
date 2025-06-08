package controller

import (
	config "chat-backend/Config"
	models "chat-backend/Models"
	utils "chat-backend/Utils"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// SignupRequest defines the JSON payload expected when a new user registers.
type SignupRequest struct {
	Email      string `json:"email"`
	FullName   string `json:"fullname"`
	Password   string `json:"password"`
	ProfilePic string `json:"profile_pic"`
}

// SignupResponse defines the JSON payload returned after successful signup.
type SignupResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse defines the JSON payload returned after successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// Steps:
// 1. Parse incoming JSON.
// 2. Check if email already exists.
// 3. Hash password.
// 4. Insert new user into MongoDB.
// 5. Return success message.

func Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	usersColl := config.GetCollection(models.CollectionNameUser)

	// Check if email already in use
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := usersColl.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	// Create the user document
	newUser := models.User{
		Email:      req.Email,
		FullName:   req.FullName,
		Password:   hashedPwd,
		ProfilePic: req.ProfilePic,
	}

	// Insert into MongoDB
	_, err = usersColl.InsertOne(ctx, newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SignupResponse{Message: "User registered successfully"})

}
