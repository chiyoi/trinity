package db

import (
	"github.com/chiyoi/trinity/pkg/atmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Name  string `bson:"name,omitempty"`
	Token string `bson:"token,omitempty"`
}

type Message struct {
	MessageId primitive.ObjectID `bson:"_id,omitempty"`
	SenderId  string             `bson:"user,omitempty"`
	Content   []atmt.Paragraph   `bson:"content,omitempty"`
}
