// File: controllers/auth.go

package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	config "chat-backend/Config"
	models "chat-backend/Models"
	utils "chat-backend/Utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	middleware "chat-backend/Middleware"
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

// LoginRequest defines the JSON payload expected when a user logs in.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse defines the JSON payload returned after successful login.
type LoginResponse struct {
	Token string `json:"token"`
}

// Signup handles user registration.
// Steps:
// 1. Parse incoming JSON.
// 2. Check if email already exists.
// 3. Hash password.
// 4. Insert new user into MongoDB.
// 5. Return success message.
func Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
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
	//take the case : all field are necessary and password must be atleast 6 char.

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SignupResponse{Message: "User registered successfully"})
}

// Login handles user authentication.
// Steps:
// 1. Parse incoming JSON.
// 2. Find user by email.
// 3. Compare stored hashed password with provided password.
// 4. If valid, generate JWT and set it as an HTTP-only cookie.
// 5. Return token in JSON as well.
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	usersColl := config.GetCollection(models.CollectionNameUser)

	// Find user document by email
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := usersColl.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token (expires in 24 hours)
	tokenString, err := middleware.GenerateToken(user.ID, time.Hour*48)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set the token as an HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}

// Logout clears the JWT cookie on the client side.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Overwrite cookie with expired value
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

//check auth

// CheckAuthResponse is the payload returned when auth succeeds.
type CheckAuthResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

// CheckAuth simply verifies that the middleware ran, and echoes back the userID.
func CheckAuth(w http.ResponseWriter, r *http.Request) {
	// 1) Pull the userID from context (set by Authenticate middleware)
	val := r.Context().Value("userID")
	id, ok := val.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2) Return a simple JSON with the userID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CheckAuthResponse{
		UserID:  id.Hex(),
		Message: "You are authenticated",
	})
}
