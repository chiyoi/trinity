package message

import "fmt"

func Format(a ...any) (msg Message) {
	for _, aa := range a {
		switch ta := aa.(type) {
		case Message:
			msg.Extend(ta)
		case Segment:
			msg.Append(ta)
		default:
			msg.Append(Text(fmt.Sprint(aa)))
		}
	}
	return
}

func Text(txt string) Segment {
	return Segment{
		Type: TypeText,
		Data: map[string]string{
			"text": txt,
		},
	}
}

func Image(url string) Segment {
	return Segment{
		Type: TypeImage,
		Data: map[string]string{
			"file": url,
		},
	}
}

func Record(url string) Segment {
	return Segment{
		Type: TypeRecord,
		Data: map[string]string{
			"file": url,
		},
	}
}
