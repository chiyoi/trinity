package eroira

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/chiyoi/trinity/internal/app/aira/client"
	atmt2 "github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/atmt/message"
	"github.com/chiyoi/trinity/pkg/atmt/rules"
)

func Eroira() (atmt2.Matcher, atmt2.HandlerFunc) {
	return rules.Keywords("色图"), handler
}

func handler(ev atmt2.Event) {
	client.PostMessage("アトリ、検索中ーー")
	img, name, disp, err := GetImage()
	if err != nil {
		client.ErrorCallback(err)
		return
	}

	url, err := client.CacheFile(img)
	if err != nil {
		client.ErrorCallback(err)
		return
	}
	client.PostMessage(message.Image(name, url), disp)
}

var lolicon = "https://api.lolicon.app/setu/v2"

var dispTmpl = `
title: %v[%v]
author: %v[%v]`

func GetImage() (img []byte, name string, disp string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get image: %w", err)
		}
	}()

	var u *url.URL
	if u, err = url.Parse(lolicon); err != nil {
		return
	}

	var query = url.Values{}
	query.Add("r18", "2")
	query.Add("proxy", "0")
	u.RawQuery += query.Encode()

	resp, err := http.Get(u.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("lolicon api call error(%d %s)", resp.StatusCode, resp.Status)
		}
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var loliconResp LoliconResponse
	if err = json.Unmarshal(body, &loliconResp); err != nil {
		return
	}
	if len(loliconResp.Data) <= 0 {
		err = errors.New("api responded empty data")
		return
	}

	data := loliconResp.Data[0]
	req, err := http.NewRequest("GET", data.Urls.Original, nil)
	if err != nil {
		return
	}
	req.Header.Add("Referer", "https://www.pixiv.net/")

	if resp, err = http.DefaultClient.Do(req); err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = errors.New(resp.Status)
		}
		return
	}

	if img, err = io.ReadAll(resp.Body); err != nil {
		return
	}
	name = path.Base(data.Urls.Original)
	disp = fmt.Sprintf(dispTmpl, data.Title, data.Pid, data.Author, data.Uid)
	return
}
