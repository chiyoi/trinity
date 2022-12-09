package db

import (
	"github.com/chiyoi/trinity/pkg/atmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Name  string `bson:"name,omitempty"`
	Token string `bson:"token,omitempty"`
}

type MessageID = primitive.ObjectID
type Message struct {
	ID      MessageID        `bson:"_id,omitempty"`
	Sender  string           `bson:"sender,omitempty"`
	Content []atmt.Paragraph `bson:"content,omitempty"`
}
