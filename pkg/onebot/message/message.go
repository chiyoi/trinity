package message

import "github.com/chiyoi/trinity/pkg/onebot"

func Text(txt string) onebot.MessageSegment {
	return onebot.MessageSegment{
		Type: onebot.MessageText,
		Data: map[string]string{
			"text": txt,
		},
	}
}

func Image(url string) onebot.MessageSegment {
	return onebot.MessageSegment{
		Type: onebot.MessageImage,
		Data: map[string]string{
			"file": url,
		},
	}
}

func Record(url string) onebot.MessageSegment {
	return onebot.MessageSegment{
		Type: onebot.MessageRecord,
		Data: map[string]string{
			"file": url,
		},
	}
}
