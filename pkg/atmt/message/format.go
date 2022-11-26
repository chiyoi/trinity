package message

import "fmt"

func Format(a ...any) (msg Message) {
	for _, item := range a {
		switch t := item.(type) {
		case Message:
			msg.Extend(t)
		case Segment:
			msg.Append(t)
		default:
			msg.Append(Text(fmt.Sprint(t)))
		}
	}
	return
}
