package aira

import (
	"os"

	"github.com/chiyoi/trinity/internal/app/aira/client"
	atmt2 "github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/message"
	"github.com/chiyoi/trinity/pkg/atmt/rules"
)

func Aira() (atmt2.Matcher, atmt2.HandlerFunc) {
	return rules.ExactMessageOneOf("aira", "アイラ"), handler
}

var (
	voicePath = "assets/aira/aira.mp3"
	voiceText = "大切な人と、いつかまた巡り合えますように。"
)

func handler(ev atmt2.Event) {
	voice, err := os.ReadFile(voicePath)
	if err != nil {
		client.ErrorCallback(err)
		return
	}

	url, err := client.CacheFile(voice)
	if err != nil {
		client.ErrorCallback(err)
		return
	}
	client.PostMessage(voiceText, message.Record("aira.mp3", url))
}
