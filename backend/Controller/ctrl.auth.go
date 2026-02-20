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

// Signup handles user registration.
// Steps:
// 1. Parse incoming JSON.
// 2. Check if email already exists.
// 3. Hash password.
// 4. Insert new user into MongoDB.
// 5. Return success message.
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

func Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	//const CollectionNameUser = "users" defined in folder - model
	//config.GetCollection defined in config/db.go: connects your Go code to a specific table (collection) in your MongoDB database.
	usersColl := config.GetCollection(models.CollectionNameUser)

	if req.Email == "" || req.FullName == "" || req.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}
	// Check if email already in use(Duplicates)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := usersColl.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		//w : response back
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		//http.StatusConflict (409): The specific HTTP code for "I can't do this because the resource already exists."
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash the password
	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Could not SignupResponsehash password", http.StatusInternalServerError)
		return
	}
	//take the case : all field are necessary and password must be atleast 6 char.

	// Create the user document
	newUser := models.User{
		Email:      req.Email,
		FullName:   req.FullName,
		Password:   hashedPwd,
		ProfilePic: req.ProfilePic,
		CreatedAt:  time.Now(), // Set the current time as CreatedAt
	}

	// Insert into MongoDB
	_, err = usersColl.InsertOne(ctx, newUser)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	//These three lines form the Final Handshake of your server's response. They package up the data and send it back to the user (the frontend).




//  w.Header() It adds a metadata "header" to the HTTP response.
//  Why: Computers can send text, images, videos, or HTML. 
//  This line tells the frontend (React/Browser): "Hey, what I am about to send you is JSON text.
//  Please treat it as an object, not just a random string."

	w.Header().Set("Content-Type", "application/json")
	//status code 
	//200 OK: Generic "It worked."
    //201 Created: Specific "It worked, and I successfully created a new resource (the user) in the database."
	//if u dont write this line go automatically sets the code to 200 ok
	w.WriteHeader(http.StatusCreated)
     //take signupResponse struct convert it to json and send response back to client
	json.NewEncoder(w).Encode(SignupResponse{
		Message: "User registered successfully"})

}











// Login handles user authentication.
// Steps:
// 1. Parse incoming JSON.
// 2. Find user by email.
// 3. Compare stored hashed password with provided password.
// 4. If valid, generate JWT and set it as an HTTP-only cookie.
// 5. Return token in JSON as well.

// LoginRequest defines the JSON payload expected when a user logs in.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse defines the JSON payload returned after successful login.

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

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

	//Generate Token: You create a digital ID card (JWT) so the user doesn't have to enter their password for every single chat message.
    //Set Cookie: You put that ID card inside a locked, armored box (HttpOnly Cookie) so that malicious scripts in the browser cannot steal

	// Set the token as an HTTP-only cookie It tells the browser: "Save this cookie, but DO NOT let any JavaScript code read it."
	http.SetCookie(w, &http.Cookie{
		Name:        "token",
		Value:       tokenString,
		Path:        "/",                // Cookie works on all pages of the site
		HttpOnly:    true,	             //// 1. SECURITY: JavaScript cannot steal this. // Prevents XSS attacks.
		Secure:      true,                  // only sent over HTTPS 
		SameSite:    http.SameSiteNoneMode, // allow cross‑site // 3. CROSS-ORIGIN: Allows the cookie to be sent even if // Frontend  and Backend (localhost:8000) are different.
		Partitioned: true,
		Expires:     time.Now().Add(48 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}{
		Token: tokenString,
		User:  user,
	})

}





// Logout clears the JWT cookie on the client side.
func Logout(w http.ResponseWriter, r *http.Request) {
	// Overwrite cookie with expired value
	http.SetCookie(w, &http.Cookie{
		Name:        "token",
		Value:       "", // Clear the token value
		Path:        "/",                // Cookie works on all pages of the site
		HttpOnly:    true,	             //// 1. SECURITY: JavaScript cannot steal this. // Prevents XSS attacks.
		Secure:      true,                  // only sent over HTTPS 
		SameSite:    http.SameSiteNoneMode, // allow cross‑site // 3. CROSS-ORIGIN: Allows the cookie to be sent even if // Frontend  and Backend (localhost:8000) are different.
		Partitioned: true,
		MaxAge:      -1, // Expire the cookie immediately
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

//check auth


// CheckAuth simply verifies that the middleware (jwt generated) ran, and echoes back the userID.
//Middleware: Checks Cookie -> Finds ID "123" -> Puts "123" in Context.
//CheckAuth (Start): Looks in Context -> Finds "123".
//CheckAuth (Middle): Asks DB for "User 123".
//CheckAuth (End): Sends "User 123" data to client,

// Authenticate middleware validates the JWT from the "Authorization" header or cookie,
// and attaches the user ID to the request context if valid.
func CheckAuth(w http.ResponseWriter, r *http.Request) {
	// 1) Pull the userID from context (set by Authenticate middleware) commig from user request
	val := r.Context().Value("userID")
	userID, ok := val.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//fetch the full user from database
	usersColl := config.GetCollection(models.CollectionNameUser)
	var user models.User
	err := usersColl.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	// 2) Return a simple JSON with the userID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}
