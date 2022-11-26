package message

import "strings"

type Message []Segment
type Segment struct {
	Type Type   `json:"type" bson:"type"`
	Ref  Ref    `json:"ref" bson:"ref"`
	Data string `json:"data" bson:"data"`
}

type Type uint8

const (
	TypeText Type = iota
	TypeImage
	TypeRecord
	TypeVideo
	TypeFile
)

type Ref struct {
	Name string `json:"name"`
	Url  string `json:"url"`
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
