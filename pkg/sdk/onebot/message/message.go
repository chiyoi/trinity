package message

type Message []Segment

type Segment struct {
	Type MessageType       `json:"type"`
	Data map[string]string `json:"data"`
}

func (msg Message) Append(seg Segment) Message  { return append(msg, seg) }
func (msg Message) Extend(msg1 Message) Message { return append(msg, msg1...) }

func (seg Segment) Chain(seg1 Segment) Message { return append(Message{}, seg, seg1) }

type MessageType string

const (
	TypeText      MessageType = "text"
	TypeFace      MessageType = "face"
	TypeRecord    MessageType = "record"
	TypeVideo     MessageType = "video"
	TypeAt        MessageType = "at"
	TypeMusic     MessageType = "music"
	TypeImage     MessageType = "image"
	TypeReply     MessageType = "reply"
	TypeRedbag    MessageType = "redbag"
	TypePoke      MessageType = "poke"
	TypeGift      MessageType = "gift"
	TypeForward   MessageType = "forward"
	TypeNode      MessageType = "node"
	TypeXML       MessageType = "xml"
	TypeJSON      MessageType = "json"
	TypeCardimage MessageType = "cardimage"
	TypeTTS       MessageType = "tts"
)
