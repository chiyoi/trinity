package aira

import (
	"github.com/chiyoi/trinity/internal/app/aira/client"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/rules"
)

const (
	voicePath = "assets/aira/aira.mp3"
	voiceText = "大切な人と、いつかまた巡り合えますように。"
)

func Aira() (atmt.Matcher, atmt.Handler) { return matcher, atmt.HandlerFunc(handler) }

var matcher = atmt.Matcher{
	Match: rules.And(
		rules.MessageType(atmt.MessagePost),
		rules.ExactlyOneOf("aira", "アイラ"),
	),
	Priority: 10,
}

func handler(resp *atmt.Message, post atmt.Message) {
	// voice, err := os.ReadFile(voicePath)
	// if err != nil {
	// 	logs.Error(err)
	// 	return
	// }

	// url, err := client.CacheFile(voice)
	// if err != nil {
	// 	handlers.ErrorCallback(err)
	// 	return
	// }
	client.PostMessage(voiceText)
}
