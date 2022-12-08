package onebot

import "fmt"

type Message []MessageSegment

type MessageSegment struct {
	Type SegmentType       `json:"type"`
	Data map[string]string `json:"data"`
}

func (msg Message) Append(seg MessageSegment) Message { return append(msg, seg) }
func (msg Message) Extend(msg1 Message) Message       { return append(msg, msg1...) }

func (seg MessageSegment) Chain(seg1 MessageSegment) Message { return append(Message{}, seg, seg1) }

type SegmentType string

const (
	SegmentText      SegmentType = "text"
	SegmentFace      SegmentType = "face"
	SegmentRecord    SegmentType = "record"
	SegmentVideo     SegmentType = "video"
	SegmentAt        SegmentType = "at"
	SegmentMusic     SegmentType = "music"
	SegmentImage     SegmentType = "image"
	SegmentReply     SegmentType = "reply"
	SegmentRedbag    SegmentType = "redbag"
	SegmentPoke      SegmentType = "poke"
	SegmentGift      SegmentType = "gift"
	SegmentForward   SegmentType = "forward"
	SegmentNode      SegmentType = "node"
	SegmentXML       SegmentType = "xml"
	SegmentJSON      SegmentType = "json"
	SegmentCardimage SegmentType = "cardimage"
	SegmentTTS       SegmentType = "tts"
)

func FormatMessage(a ...any) (msg Message) {
	for _, aa := range a {
		switch ta := aa.(type) {
		case Message:
			msg.Extend(ta)
		case MessageSegment:
			msg.Append(ta)
		default:
			msg.Append(Text(fmt.Sprint(aa)))
		}
	}
	return
}

func Text(txt string) MessageSegment {
	return MessageSegment{
		Type: SegmentText,
		Data: map[string]string{
			"text": txt,
		},
	}
}

func Image(url string) MessageSegment {
	return MessageSegment{
		Type: SegmentImage,
		Data: map[string]string{
			"file": url,
		},
	}
}

func Record(url string) MessageSegment {
	return MessageSegment{
		Type: SegmentRecord,
		Data: map[string]string{
			"file": url,
		},
	}
}
