package atri

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"path"
	"sync"
	"time"

	"github.com/chiyoi/trinity/internal/app/aira/client"
	atmt2 "github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/message"
	"github.com/chiyoi/trinity/pkg/atmt/rules"
)

func initiate() (err error) {
	f, err := os.Open(path.Join(assetsBase, "text", "atri.json"))
	if err != nil {
		return
	}

	var data []byte
	data, err = io.ReadAll(f)
	if err != nil {
		return
	}
	if err = json.Unmarshal(data, &voiceList); err != nil {
		return
	}
	rand.Seed(time.Now().UnixMicro())
	return
}

func Atri() (atmt2.Matcher, atmt2.HandlerFunc) {
	var init sync.Once
	var err error
	return rules.ExactMessageOneOf("aira", "アトリ"), func(ev atmt2.Event) {
		if init.Do(func() {
			err = initiate()
		}); err != nil {
			client.ErrorCallback(err)
			return
		}
		handler(ev)
	}
}

var assetsBase = "assets/atri"

type voiceDisp struct {
	Path string `json:"o"`
	Text string `json:"s"`
}

var voiceList []voiceDisp

func handler(ev atmt2.Event) {
	idx := rand.Intn(len(voiceList))
	fp, txt := path.Join(assetsBase, "voice", voiceList[idx].Path), voiceList[idx].Text
	voice, err := os.ReadFile(fp)
	if err != nil {
		client.ErrorCallback(err)
		return
	}
	url, err := client.CacheFile(voice)
	if err != nil {
		client.ErrorCallback(err)
		return
	}
	client.PostMessage(txt, message.Record("atri.mp3", url))
}
