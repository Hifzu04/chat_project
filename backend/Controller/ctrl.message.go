// File: controllers/user.go

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	config "chat-backend/Config"
	models "chat-backend/Models"
	utils "chat-backend/Utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"
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

	//Find does not return a list (Slice). It returns a Cursor.
	cursor, err := coll.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{
		"password": 0,
	}))
	if err != nil {
		log.Println(" GetAllUsers Find error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// 3) Decode into slice of UserListItem
	var list []UserListItem
	//cursor.Next(ctx): Moves the pointer to the next user. Returns true if there is a user, false if we reached the end.
	for cursor.Next(ctx) {
		var u models.User
		if err := cursor.Decode(&u); err != nil {
			log.Println(" GetAllUsers Decode error:", err)
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
		log.Println(" GetAllUsers Cursor error:", err)
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
	//r.Context().Value("userID"): Pulls the "who is logged in?" info that your Middleware put there.
	senderID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form-data for text + images
	//Why ParseMultipartForm? Normal JSON requests cannot easily handle file uploads (images). We use a format called Multipart/Form-Data.
	// //It splits the request into parts: Part 1 is the text ("Hello"), Part 2 is the image file data.
	// //10 << 20: This is a bitwise operation that equals 10 MB (10 * 1024 * 1024).

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

	//slice of strings to hold urls of uploaded images
	var imageURLs []string

	//r.MultipartForm.File["images"]: Since a user can send multiple images at once, this returns a list (slice) of file headers.
	files := r.MultipartForm.File["images"]
	for _, fh := range files {
		//.Open(): Actually opens the data stream so we can read the pixels.
		file, err := fh.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Convert multipart.File to io.Reader
		//We cant send actual image to db (heavy and slow ) , so we pass it to cloudinary and get back a url
		//UploadToCloudinary reads the file data and uploads it to Cloudinary, returning the URL of the uploaded image.
		uploadedURL, err := utils.UploadToCloudinary(file, fh)
		if err != nil {
			log.Println("Cloudinary error:", err)
			continue
		}
		//keep adding uploaded image urls to slice(like in map keep pushing back values)
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

	//Updating the ID: InsertOne generates a new unique ID (_id) for the message. We grab this ID from the result (insertRes) and put it back into our msg struct.

	//Why? The frontend needs this new ID immediately (e.g., to confirm "Message Sent" or to let the user delete it later).

	///Encode: Sends the full saved message object back to the client.

	msg.ID = insertRes.InsertedID.(primitive.ObjectID)

	nestBotIDHex := os.Getenv("NESTBOT_ID")
	if receiverIDHex == nestBotIDHex {
		// Use 'go' to run this in the background so the user doesn't have to wait!
		go replyWithAI(senderID, receiverID, text, imageURLs)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// Get messages between two users
// This function is responsible for loading the chat history when you click on a friend's name in the sidebar. I
// and the user whose ID is in the path: GET /messages/{userID}
func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1) Extract the receiverID from the path
	//gorilla/mux router, this function grabs variables defined in CURRENT route URL.  //FRIENDS ID
	vars := mux.Vars(r)

	receiverIDHex, ok := vars["userID"]
	if !ok {
		http.Error(w, "Missing userID in path", http.StatusBadRequest)
		return
	}
	//Why: The URL contains a string ("abc...").
	//  MongoDB requires a binary ObjectID. This converts the string into the database format.
	receiverID, err := primitive.ObjectIDFromHex(receiverIDHex)
	if err != nil {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
		return
	}

	// 2) Get the authenticated userID from the request context . (logged in user)
	userIDVal := r.Context().Value("userID")
	senderID, ok := userIDVal.(primitive.ObjectID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3) Build the two‑way filter USES OR TO GRAB TWO WAY MSG
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

	//options.Find().SetSort(...):   Sorts the results before returning them.
	//"created_at": 1: The 1 stands for Ascending Order (Oldest -> Newest). Show chat history form oldest(top) to newest(bottom).
	//If you used -1 (Descending), the newest messages would appear at the top, which is useful for "News Feeds" but usually wrong for Chat.
	//cursor: A pointer to the list of messages matching our filter. finds returns a cursor
	cursor, err := coll.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": 1}))
	if err != nil {
		log.Println(" GetMessages Find error:", err)
		http.Error(w, "Database error while fetching messages betweeen u and ur frnd", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	// 5) Decode all messages
	var messages []models.Message
	//Difference: In your previous GetAllUsers code, you used a manual loop (for cursor.Next()).
	// Here, cursor.All does the loop for you automatically and dumps everything into the messages slice.

	if err := cursor.All(ctx, &messages); err != nil {
		log.Println(" GetMessages Decode error:", err)
		http.Error(w, "Data error while decoding messages", http.StatusInternalServerError)
		return
	}

	// 6) Return as JSON
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Println(" GetMessages Encode error:", err)
	}
}

func replyWithAI(user primitive.ObjectID, botID primitive.ObjectID, usertext string, imageURLs []string) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		//In Go, log.Fatal() means: "Print this error, and then instantly CRASH the entire server." If Google's API has a hiccup,
		// or your internet drops for a second, your whole backend will die, kicking every user offline.
		fmt.Println("maybe something is wrong in getting genAI API")
		return
	}

	configofgenAI := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText("You are NestBot, a friendly and helpful AI assistant inside a chat application. Keep your answers concise, conversational, and under 3 sentences. Do not use complex markdown formatting.", genai.RoleUser),
	}

	// 1. Create an empty list to hold our text and image parts
	var parts []*genai.Part

	// 2. Add the text (if the user typed anything)
	if usertext != "" {
		parts = append(parts, genai.NewPartFromText(usertext))
	}

	for _, imgURL := range imageURLs {
		imageResp, err := http.Get(imgURL)
		if err != nil {
			log.Println("Failed to download image from Cloudinary:", err)
			continue // Skip this one and move to the next image
		}

		//// Read the pixels
		imageBytes, err := io.ReadAll(imageResp.Body)
		imageResp.Body.Close() // ALWAYS close the body to prevent memory leaks
		
		if err == nil {
			// Append the image part
			parts = append(parts, genai.NewPartFromBytes(imageBytes, "image/jpeg"))
		}

	}


	if len(parts) == 0 {
        return
    }

	contents := []*genai.Content{
        genai.NewContentFromParts(parts, genai.RoleUser),
    }

	
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		contents,
		configofgenAI,
	)
	if err != nil {
		//log.Fatal(err)
		log.Println("Error generating AI response:", err)
		return
	}

	aiReplyText := (result.Text())
	fmt.Println(aiReplyText)

	// 2. SAVE TO MONGODB(The response of AI)
	// We flip the sender/receiver. The Bot is now the sender!
	msg := models.Message{
		SenderID:   botID,
		ReceiverID: user,
		Text:       aiReplyText,
		CreatedAt:  time.Now(),
	}

	collection := config.GetCollection(models.CollectionNameMessage)
	insertRes, err := collection.InsertOne(ctx, msg)
	if err != nil {
		log.Println("Error saving AI message to DB:", err)
		return
	}

	//broadcast to websocket?
	// Grab the new Mongo ID for the WebSocket payload
	msg.ID = insertRes.InsertedID.(primitive.ObjectID)

	// Tell the specific user (userWhoSentMsg) that they have a "newMessage"
	// We must pass the `msg` object so React knows the ID, Sender, and Text.
	//React frontend is currently programmed (in useAuthStore.js) to listen for a WebSocket event called "newMessage".
	//When the AI finishes thinking and saving to the database, it triggers utils.SendTo.
	//This blasts the message straight through the open WebSocket to the user's browser, and the chat bubble instantly appears!
	utils.SendTo(user.Hex(), "newMessage", msg)

	//TODO //Try for Streaming responses in chunks

}
