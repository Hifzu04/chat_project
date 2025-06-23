// File: models/user.go

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account in the chat application.
// The `bson` tags define how fields are stored in MongoDB.
type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email      string             `bson:"email" json:"email"`
	FullName   string             `bson:"fullname" json:"fullname"`
	Password   string             `bson:"password" json:"-"` // omit password in JSON responses
	ProfilePic string             `bson:"profile_pic" json:"profile_pic"`
	CreatedAt  time.Time          `bson:"created_at"      json:"createdAt"`
}

// CollectionNameUser is the MongoDB collection name for users.
const CollectionNameUser = "users"
