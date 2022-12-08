package atmt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Data    json.RawMessage `json:"data"`
	Content []Paragraph     `json:"content"`
}
type MessageType uint8

const (
	MessagePost MessageType = iota + 1
	MessageRequest
	MessageResponse
)

func (typ MessageType) String() (str string) {
	defer func() {
		str = fmt.Sprintf("MessageType(%s)", str)
	}()
	switch typ {
	case MessagePost:
		return "post"
	case MessageRequest:
		return "request"
	case MessageResponse:
		return "response"
	default:
		return "invalid type"
	}
}

type Paragraph struct {
	Type ParagraphType `json:"type" bson:"type,omitempty"`
	Text string        `json:"text" bson:"text,omitempty"`
	Name string        `json:"name" bson:"name,omitempty"`
	Ref  string        `json:"ref" bson:"text,omitempty"`
}

type ParagraphType uint8

const (
	ParagraphText ParagraphType = iota + 1
	ParagraphImage
	ParagraphRecord
)

type DataPost struct {
	MessageId primitive.ObjectID `json:"message_id"`
	SenderId  string             `json:"sender_id"`
}

type DataRequest[Action ~uint8] struct {
	Action Action          `json:"action"`
	Args   json.RawMessage `json:"args"`
}

type DataResponse struct {
	StatusCode StatusCode      `json:"status_code"`
	Args       json.RawMessage `json:"args"`
}

type StatusCode uint8

const (
	StatusOK StatusCode = iota + 1
	StatusBadRequest
	StatusUnauthorized
	StatusInternalServerError
)

func (c StatusCode) String() string {
	return fmt.Sprintf("StatusCode(%s)", c.Text())
}

func (c StatusCode) Text() string {
	switch c {
	case StatusOK:
		return "ok"
	case StatusBadRequest:
		return "bad request"
	case StatusUnauthorized:
		return "unauthorized"
	case StatusInternalServerError:
		return "internal server error"
	default:
		return "invalid code"
	}
}

func (msg *Message) Plaintext() string {
	if reflect.TypeOf(msg) != reflect.TypeOf(&Message{}) {
		return ""
	}
	var buf strings.Builder
	for _, para := range msg.Content {
		if para.Type == ParagraphText {
			buf.WriteString(para.Text)
		}
	}
	return buf.String()
}

type MessageBuilder[Data any] struct {
	Type    MessageType `json:"type"`
	Data    Data        `json:"data"`
	Content []Paragraph `json:"content"`
}

func (b *MessageBuilder[Data]) Message() (msg Message, err error) {
	data, err := json.Marshal(b.Data)
	if err != nil {
		return
	}
	return Message{
		Type:    b.Type,
		Data:    data,
		Content: b.Content,
	}, err
}
