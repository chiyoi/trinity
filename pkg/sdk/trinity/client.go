package trinity

import (
	"encoding/json"
	"fmt"

	"github.com/chiyoi/trinity/pkg/atmt"
)

type RequestError interface {
	error
	Action() Action
	StatusCode() atmt.StatusCode
}

type requestError struct {
	act  Action
	code atmt.StatusCode
}

var _ RequestError = (*requestError)(nil)

func (err *requestError) Error() string {
	return fmt.Sprintf("post error(%d %s)", err.code, err.code)
}

func (err *requestError) Action() Action              { return err.act }
func (err *requestError) StatusCode() atmt.StatusCode { return err.code }

func PostMessage(url string, auth string, content []atmt.Paragraph) (err error) {
	sender, _, err := ParseAuthToken(auth)
	if err != nil {
		return
	}
	req, err := (&atmt.MessageBuilder[RequestBuilder[ArgsPostMessage]]{
		Type: atmt.MessageRequest,
		Data: RequestBuilder[ArgsPostMessage]{
			Action: ActionPostMessage,
			Args: ArgsPostMessage{
				Sender: sender,
				Auth:   auth,
			},
		},
		Content: content,
	}).Message()
	resp, err := atmt.SendMessage(url, req)
	if err != nil {
		return
	}
	if resp.Type != atmt.MessageResponse {
		err = fmt.Errorf("unexpected non-response message")
		return
	}
	var data atmt.DataResponse
	if err = json.Unmarshal(resp.Data, &data); err != nil {
		return
	}
	if data.StatusCode != atmt.StatusOK {
		err = &requestError{
			act:  ActionPostMessage,
			code: data.StatusCode,
		}
	}
	return
}

// func CallApiCtx[Data RespData](ctx context.Context, url string, auth string, req Request) (resp Response[Data], err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("call api: %w", err)
// 		}
// 	}()
// 	b, err := json.Marshal(req)
// 	if err != nil {
// 		return
// 	}

// 	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
// 	if err != nil {
// 		return
// 	}
// 	httpReq.Header.Set("Content-Type", "application/json")
// 	httpReq.Header.Set("Authorization", auth)

// 	httpResp, err := http.DefaultClient.Do(httpReq)
// 	if err != nil || httpResp.StatusCode != http.StatusOK {
// 		if err == nil {
// 			err = &apiCallError{
// 				action:     req.Action,
// 				statusCode: httpResp.StatusCode,
// 			}
// 		}
// 		return
// 	}

// 	if b, err = io.ReadAll(httpResp.Body); err != nil {
// 		return
// 	}
// 	if err = json.Unmarshal(b, &resp); err != nil {
// 		logs.Debug(string(b))
// 		return
// 	}
// 	return
// }
// func CallApi[Data RespData](url string, auth string, req Request) (resp Response[Data], err error) {
// 	return CallApiCtx[Data](context.Background(), url, auth, req)
// }

// func PostMessage(url string, auth string, a ...any) (messageId string, err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("post message: %w", err)
// 		}
// 	}()
// 	resp, err := CallApi[RespDataPostMessage](url, auth, Request{
// 		Action: ActionPostMessage,
// 		Data: ReqDataPostMessage{
// 			Message: message.Format(a...),
// 		},
// 	})
// 	if err != nil {
// 		return
// 	}
// 	messageId = resp.Data.MessageId
// 	return
// }

// func GetMessage(url string, auth string, id string) (data RespDataGetMessage, err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("get message: %w", err)
// 		}
// 	}()
// 	resp, err := CallApi[RespDataGetMessage](url, auth, Request{
// 		Action: ActionQueryMessageIdsTimeRange,
// 		Data: ReqDataGetMessage{
// 			Id: id,
// 		},
// 	})
// 	if err != nil {
// 		return
// 	}
// 	data = resp.Data
// 	return
// }

// func QueryMessageIdsTimeRange(url string, auth string, from, to int64) (ids []string, err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("query message ids time range: %w", err)
// 		}
// 	}()
// 	resp, err := CallApi[RespDataQueryMessageTimeRange](url, auth, Request{
// 		Action: ActionQueryMessageIdsTimeRange,
// 		Data: ReqDataQueryMessageTimeRange{
// 			From: from,
// 			To:   to,
// 		},
// 	})
// 	if err != nil {
// 		return
// 	}
// 	ids = resp.Data.Ids
// 	return
// }

// func CacheFileCtx(ctx context.Context, url string, auth string, data []byte) (sasUrl string, err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("cache file: %w", err)
// 		}
// 	}()
// 	md5Sum := md5.Sum(data)
// 	resp, err := CallApiCtx[RespDataCacheFile](ctx, url, auth, Request{
// 		Action: ActionQueryMessageIdsTimeRange,
// 		Data: Request{
// 			Action: ActionCacheFile,
// 			Data: ReqDataCacheFile{
// 				Md5SumHex: fmt.Sprintf("%x", md5Sum),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		return
// 	}
// 	sasUrl = resp.Data.SasURL

// 	u, err := urlpkg.Parse(sasUrl)
// 	if err != nil {
// 		return
// 	}
// 	credential := azblob.NewAnonymousCredential()
// 	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
// 	blockBlobUrl := azblob.NewBlockBlobURL(*u, pipeline)

// 	properties, err := blockBlobUrl.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
// 	if respErr, ok := err.(azblob.ResponseError); ok && respErr.Response().StatusCode == http.StatusNotFound {
// 		err = nil
// 	}
// 	if err != nil {
// 		return
// 	}

// 	if !reflect.DeepEqual(md5Sum, properties.ContentMD5()) {
// 		if _, err = azblob.UploadBufferToBlockBlob(ctx, data, blockBlobUrl, azblob.UploadToBlockBlobOptions{}); err != nil {
// 			return
// 		}
// 	}
// 	return
// }
// func CacheFile(url string, auth string, data []byte) (sasUrl string, err error) {
// 	return CacheFileCtx(context.Background(), url, auth, data)
// }

// func VerifyAuthorization(url string, auth string) (pass bool, err error) {
// 	defer func() {
// 		if err != nil {
// 			err = fmt.Errorf("verify authorization: %w", err)
// 		}
// 	}()
// 	if _, err = CallApi[RespDataVerifyAuthorization](url, auth, Request{
// 		Action: ActionVerifyAuthorization,
// 		Data:   ReqDataVerifyAuthorization{},
// 	}); err != nil {
// 		if acErr, ok := err.(ApiCallError); ok && (acErr.StatusCode() == http.StatusUnauthorized || acErr.StatusCode() == http.StatusForbidden) {
// 			pass = false
// 			return
// 		}
// 		return
// 	}
// 	pass = true
// 	return
// }
