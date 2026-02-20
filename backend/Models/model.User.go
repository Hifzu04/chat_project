// File: models/user.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account in the chat application.
// The `bson` Binary JSON tags define how fields are stored in MongoDB.
//The Database World (MongoDB) speaks BSON. /"When you save this struct to the database, map the ID field to the database column _id."
//The Web World (Frontend/React) speaks JSON. :eg json:"id"`  This tag tells Go: "When you send this data to the browser, rename this field to id."
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`     
	Email      string             `bson:"email" json:"email"`
	FullName   string             `bson:"fullname" json:"fullname"`
	Password   string             `bson:"password" json:"-"` // - ensures even if you convert the whole User struct to JSON to send to frontend, the password field is stripped out
	ProfilePic string             `bson:"profile_pic" json:"profile_pic"`
	CreatedAt  time.Time          `bson:"created_at"      json:"createdAt"`
}

// CollectionNameUser is the MongoDB collection(table) name for users.
const CollectionNameUser = "users"
