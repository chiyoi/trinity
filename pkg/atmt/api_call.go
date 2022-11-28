package atmt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiCallError struct {
	StatusCode int
}

func (e *ApiCallError) Error() string {
	return fmt.Sprintf("api call error(%d %s)", e.StatusCode, http.StatusText(e.StatusCode))
}

func CallApiCtx(ctx context.Context, url string, req Request) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("call api: %w", err)
		}
	}()

	b, err := json.Marshal(req)
	if err != nil {
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil || httpResp.StatusCode != http.StatusOK {
		if err == nil {
			err = &ApiCallError{
				StatusCode: httpResp.StatusCode,
			}
		}
		return
	}
	return
}
func CallApi(url string, req Request) (err error) {
	return CallApiCtx(context.Background(), url, req)
}

func SendEvent(url string, ev Event) (err error) {
	return CallApi(url, Request{
		Time:      ev.Time.Unix(),
		User:      ev.User,
		MessageId: ev.MessageId,
		Message:   ev.Message,
	})
}
