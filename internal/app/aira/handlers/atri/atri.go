package atri

// func initiate() (err error) {
// 	f, err := os.Open(path.Join(assetsBase, "text", "atri.json"))
// 	if err != nil {
// 		return
// 	}

// 	var data []byte
// 	data, err = io.ReadAll(f)
// 	if err != nil {
// 		return
// 	}
// 	if err = json.Unmarshal(data, &voiceList); err != nil {
// 		return
// 	}
// 	rand.Seed(time.Now().UnixMicro())
// 	return
// }

// func Atri() (atmt.Matcher, atmt.HandlerFunc) {
// 	var init sync.Once
// 	var err error
// 	return rules.ExactMessageOneOf("aira", "アトリ"), func(ev atmt.Event) {
// 		if init.Do(func() {
// 			err = initiate()
// 		}); err != nil {
// 			handlers.ErrorCallback(err)
// 			return
// 		}
// 		handler(ev)
// 	}
// }

// var assetsBase = "assets/atri"

// type voiceDisp struct {
// 	Path string `json:"o"`
// 	Text string `json:"s"`
// }

// var voiceList []voiceDisp

// func handler(ev atmt.Event) {
// 	idx := rand.Intn(len(voiceList))
// 	fp, txt := path.Join(assetsBase, "voice", voiceList[idx].Path), voiceList[idx].Text
// 	voice, err := os.ReadFile(fp)
// 	if err != nil {
// 		handlers.ErrorCallback(err)
// 		return
// 	}
// 	url, err := client.CacheFile(voice)
// 	if err != nil {
// 		handlers.ErrorCallback(err)
// 		return
// 	}
// 	client.PostMessage(txt, message.Record("atri.mp3", url))
// }
