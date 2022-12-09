package atmt

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Data    json.RawMessage `json:"data"`
	Content Content         `json:"content"`
}

type Content []Paragraph

func (c Content) String() string {
	return c.Plaintext()
}

type MessageType uint8

const (
	MessagePush MessageType = iota + 1
	MessageRequest
	MessageResponse
)

func (typ MessageType) String() (str string) {
	defer func() {
		str = fmt.Sprintf("MessageType(%s)", str)
	}()
	switch typ {
	case MessagePush:
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
	Ref  string        `json:"ref" bson:"ref,omitempty"`
}

type ParagraphType uint8

const (
	ParagraphText ParagraphType = iota + 1
	ParagraphImage
	ParagraphRecord
	ParagraphVideo
)

type MessageID = primitive.ObjectID

type DataPush struct {
	MessageID MessageID `json:"message_id"`
	Sender    string    `json:"sender"`
}

type DataRequest[Action ~uint8] struct {
	Action Action          `json:"action"`
	Args   json.RawMessage `json:"args"`
}

type DataResponse struct {
	StatusCode StatusCode      `json:"status_code"`
	Values     json.RawMessage `json:"values"`
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

func (c *Content) Plaintext() string {
	var buf strings.Builder
	for _, para := range *c {
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

func (b *MessageBuilder[Data]) Write(w *Message) (err error) {
	*w, err = b.Message()
	return
}

type DataRequestBuilder[Action ~uint8, Args any] struct {
	Action Action `json:"action"`
	Args   Args   `json:"args"`
}

type DataResponseBuilder[Values any] struct {
	StatusCode StatusCode `json:"status_code"`
	Values     Values     `json:"values"`
}

func FormatContent(v ...any) (content Content) {
	for _, a := range v {
		switch t := a.(type) {
		case Paragraph:
			content = append(content, t)
		default:
			content = append(content, Text(fmt.Sprint(a)))
		}
	}
	return
}

func Text(txt string) Paragraph {
	return Paragraph{
		Type: ParagraphText,
		Text: txt,
	}
}

func Image(name, ref string) Paragraph {
	return Paragraph{
		Type: ParagraphImage,
		Name: name,
		Ref:  ref,
	}
}

func Record(name, ref string) Paragraph {
	return Paragraph{
		Type: ParagraphRecord,
		Name: name,
		Ref:  ref,
	}
}

func Video(name, ref string) Paragraph {
	return Paragraph{
		Type: ParagraphVideo,
		Name: name,
		Ref:  ref,
	}
}
