package atmt

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func SendMessageCtx(ctx context.Context, url string, msg Message) (resp Message, err error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return
	}
	if data, err = io.ReadAll(httpResp.Body); err != nil {
		return
	}
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	return
}
func SendMessage(url string, msg Message) (resp Message, err error) {
	return SendMessageCtx(context.Background(), url, msg)
}

func PostCtx(ctx context.Context, url string, msg Message) (err error) {
	if msg.Type != MessagePost {
		err = &messageTypeError{
			typ: msg.Type,
			exp: MessagePost,
		}
	}
	resp, err := SendMessageCtx(ctx, url, msg)
	if err != nil {
		return
	}
	if resp.Type != MessageResponse {
		err = &messageTypeError{
			typ: resp.Type,
			exp: MessageResponse,
		}
	}
	var respData DataResponse
	if err = json.Unmarshal(resp.Data, &respData); err != nil {
		return
	}
	if respData.StatusCode != StatusOK {
		err = &postError{
			code: respData.StatusCode,
		}
		return
	}
	return
}
func Post(url string, msg Message) (err error) {
	return PostCtx(context.Background(), url, msg)
}
