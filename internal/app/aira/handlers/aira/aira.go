package aira

import (
	"os"
	"path/filepath"

	"github.com/chiyoi/trinity/internal/app/aira/client"
	"github.com/chiyoi/trinity/internal/app/aira/handlers"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/rules"
)

func Aira() (atmt.Matcher, atmt.Handler) {
	return matcher, handler()
}

var matcher = atmt.Matcher{
	Match:    rules.ExactlyOneOf("aira", "アイラ"),
	Priority: 10,
}

func handler() atmt.HandlerFunc {
	var voicePath = filepath.Join("assets", "aira", "aira.mp3")
	const voiceText = "大切な人と、いつかまた巡り合えますように。"
	return func(resp *atmt.Message, post atmt.Message) {
		logPrefix := "aira:"
		voice, err := os.ReadFile(voicePath)
		if err != nil {
			logs.Error(logPrefix, err)
			handlers.Error(err)
			return
		}

		url, err := client.CacheBlob(voice)
		if err != nil {
			logs.Error(logPrefix, err)
			handlers.Error(err)
			return
		}
		client.PostMessage(atmt.Record(filepath.Base(voicePath), url), voiceText)
	}
}
