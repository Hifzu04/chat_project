// File: models/message.go

package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a chat message between two users.
// - SenderID and ReceiverID are ObjectIDs of the users collection.
// - Text is the message content.
// - Images is an optional slice of image URLs/filenames.
// - CreatedAt is a timestamp of when the message was sent.
type Message struct {
    ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
    SenderID   primitive.ObjectID   `bson:"sender_id" json:"sender_id"`
    ReceiverID primitive.ObjectID   `bson:"receiver_id" json:"receiver_id"`
    Text       string               `bson:"text,omitempty" json:"text,omitempty"`
    Images     []string             `bson:"images,omitempty" json:"images,omitempty"`
    CreatedAt  time.Time            `bson:"created_at" json:"created_at"`
}

// CollectionNameMessage is the MongoDB collection name for chat messages.
const CollectionNameMessage = "messages"
