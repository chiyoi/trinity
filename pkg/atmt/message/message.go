package message

import "strings"

type Message []Segment
type Segment struct {
	Type Type   `json:"type" bson:"type,omitempty"`
	Ref  Ref    `json:"ref" bson:"ref,omitempty"`
	Data string `json:"data" bson:"data,omitempty"`
}

type Type uint8

const (
	invalid Type = iota
	TypeText
	TypeImage
	TypeRecord
	TypeVideo
	TypeFile
)

func (t Type) String() string {
	switch t {
	case TypeText:
		return "text"
	case TypeImage:
		return "image"
	case TypeRecord:
		return "record"
	case TypeVideo:
		return "video"
	case TypeFile:
		return "file"
	default:
		return "invalid"
	}
}

type Ref struct {
	Name string `json:"name" bson:"name,omitempty"`
	Url  string `json:"url" bson:"url,omitempty"`
}

func (msg Message) Append(seg Segment) Message  { return append(msg, seg) }
func (msg Message) Extend(msg1 Message) Message { return append(msg, msg1...) }

func (seg Segment) Chain(seg1 Segment) Message { return append(Message{}, seg, seg1) }

func (msg Message) Plaintext() string {
	var buf strings.Builder
	for _, seg := range msg {
		if seg.Type == TypeText {
			buf.WriteString(seg.Data)
		}
	}
	return buf.String()
}
