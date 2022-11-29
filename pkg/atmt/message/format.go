package message

import "fmt"

func Format(a ...any) (msg Message) {
	for _, item := range a {
		switch t := item.(type) {
		case Message:
			msg = msg.Extend(t)
		case Segment:
			msg = msg.Append(t)
		default:
			msg = msg.Append(Text(fmt.Sprint(t)))
		}
	}
	return
}
