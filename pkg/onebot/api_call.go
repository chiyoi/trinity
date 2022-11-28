package onebot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chiyoi/trinity/pkg/onebot/message"
)

type ApiCallError struct {
	Action     Action
	StatusCode int
}

func (e *ApiCallError) Error() string {
	return fmt.Sprintf("api call error[%s](%d %s)", e.Action, e.StatusCode, http.StatusText(e.StatusCode))
}

func CallApiCtx[Data RespData](ctx context.Context, url string, req Request) (resp Response[Data], err error) {
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
				Action:     req.Action,
				StatusCode: httpResp.StatusCode,
			}
		}
		return
	}
	if b, err = io.ReadAll(httpReq.Body); err != nil {
		return
	}

	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	return
}
func CallApi[Data RespData](url string, req Request) (resp Response[Data], err error) {
	return CallApiCtx[Data](context.Background(), url, req)
}

func SendMsg(url string, id UserId, a ...any) (messageId MessageId, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("send msg: %w", err)
		}
	}()

	resp, err := CallApi(url, Request{
		Action: ActionSendMsg,
		Params: ReqParamsSendMsg{
			MessageType: MessagePrivate,
			UserId:      id,
			Message:     message.Format(a...),
		},
	})
	if err != nil {
		return
	}
	messageId = resp.Data.MessageId
	return
}
