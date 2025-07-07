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
	utils "chat-backend/Utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// only if the user is logged in . in routes GET /users (protected by your Authenticate middleware) returns:
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

func SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//http.request has prebuilt context
	senderID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form-data for text + images
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	receiverIDHex := r.FormValue("receiver_id")
	receiverID, err := primitive.ObjectIDFromHex(receiverIDHex)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")

	// Handle multiple image uploads
	var imageURLs []string
	files := r.MultipartForm.File["images"]
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Convert multipart.File to io.Reader
		uploadedURL, err := utils.UploadToCloudinary(file, fh)
		if err != nil {
			log.Println("Cloudinary error:", err)
			continue
		}
		imageURLs = append(imageURLs, uploadedURL)
	}
	// Save to DB
	msg := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Text:       text,
		Images:     imageURLs,
		CreatedAt:  time.Now(),
	}

	collection := config.GetCollection(models.CollectionNameMessage)
	insertRes, err := collection.InsertOne(ctx, msg)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	msg.ID = insertRes.InsertedID.(primitive.ObjectID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// Get messages between two users
// and the user whose ID is in the path: GET /messages/{userID}
func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1) Extract the receiverID from the path
	vars := mux.Vars(r)
	receiverIDHex, ok := vars["userID"]
	if !ok {
		http.Error(w, "Missing userID in path", http.StatusBadRequest)
		return
	}
	receiverID, err := primitive.ObjectIDFromHex(receiverIDHex)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	// 2) Get the authenticated userID from the request context
	userIDVal := r.Context().Value("userID")
	senderID, ok := userIDVal.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3) Build the two‑way filter
	filter := bson.M{
		"$or": []bson.M{
			{"sender_id": senderID, "receiver_id": receiverID},
			{"sender_id": receiverID, "receiver_id": senderID},
		},
	}

	// 4) Query MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := config.GetCollection(models.CollectionNameMessage)
	cursor, err := coll.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		log.Println("❌ GetMessages Find error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// 5) Decode all messages
	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		log.Println("❌ GetMessages Decode error:", err)
		http.Error(w, "Data error", http.StatusInternalServerError)
		return
	}

	// 6) Return as JSON
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Println("❌ GetMessages Encode error:", err)
	}
}