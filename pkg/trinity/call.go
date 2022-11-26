package trinity

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"

	"github.com/chiyoi/trinity/pkg/atmt/message"
)

type ApiCallError struct {
	Action Action
	Status string
}

func (e *ApiCallError) Error() string {
	return fmt.Sprintf("api call [%s] error(%s)", e.Action, e.Status)
}

func PostMessageCtx(ctx context.Context, url string, auth string, a ...any) (messageId string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("post message: %w", err)
		}
	}()
	msg := message.Format(a...)

	req := Request{
		Action: ActionPostMessage,
		Data: ReqDataPostMessage{
			Message: msg,
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", auth)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = &ApiCallError{
			Action: ActionPostMessage,
			Status: httpResp.Status,
		}
		return
	}

	if b, err = io.ReadAll(httpReq.Body); err != nil {
		return
	}
	var resp Response[RespDataPostMessage]
	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	messageId = resp.Data.MessageId
	return
}
func PostMessage(url string, auth string, a ...any) (messageId string, err error) {
	return PostMessageCtx(context.Background(), url, auth, a...)
}

func GetMessageCtx(ctx context.Context, url string, auth string, id string) (data RespDataGetMessage, err error) {
	req := Request{
		Action: ActionQueryMessageIdsTimeRange,
		Data: ReqDataGetMessage{
			Id: id,
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", auth)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = &ApiCallError{
			Action: ActionPostMessage,
			Status: httpResp.Status,
		}
		return
	}

	b, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	var resp Response[RespDataGetMessage]
	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	data = resp.Data
	return
}
func GetMessage(url string, auth string, id string) (data RespDataGetMessage, err error) {
	return GetMessageCtx(context.Background(), url, auth, id)
}

func QueryMessageIdsTimeRangeCtx(ctx context.Context, url string, auth string, from, to int64) (ids []string, err error) {
	req := Request{
		Action: ActionQueryMessageIdsTimeRange,
		Data: ReqDataQueryMessageTimeRange{
			From: from,
			To:   to,
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", auth)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = &ApiCallError{
			Action: ActionPostMessage,
			Status: httpResp.Status,
		}
		return
	}

	b, err = io.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	var resp Response[RespDataQueryMessageTimeRange]
	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	ids = resp.Data.Ids
	return
}
func QueryMessageIdsTimeRange(url string, auth string, from, to int64) (ids []string, err error) {
	return QueryMessageIdsTimeRangeCtx(context.Background(), url, auth, from, to)
}

func CacheFileCtx(ctx context.Context, url string, auth string, data []byte) (sasUrl string, err error) {
	md5Sum := md5.Sum(data)
	req := Request{
		Action: ActionCacheFile,
		Data: ReqDataCacheFile{
			Md5SumHex: fmt.Sprintf("%x", md5Sum),
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", auth)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		err = &ApiCallError{
			Action: ActionCacheFile,
			Status: httpResp.Status,
		}
		return
	}

	if b, err = io.ReadAll(httpReq.Body); err != nil {
		return
	}
	var resp Response[RespDataCacheFile]
	if err = json.Unmarshal(b, &resp); err != nil {
		return
	}
	sasUrl = resp.Data.SasURL

	u, err := urlpkg.Parse(sasUrl)
	if err != nil {
		return
	}
	credential := azblob.NewAnonymousCredential()
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	blockBlobUrl := azblob.NewBlockBlobURL(*u, pipeline)

	properties, err := blockBlobUrl.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if respErr, ok := err.(azblob.ResponseError); ok && respErr.Response().StatusCode == http.StatusNotFound {
		err = nil
	}
	if err != nil {
		return
	}

	identicalMd5 := func() bool {
		for i, b := range properties.ContentMD5() {
			if b != md5Sum[i] {
				return false
			}

		}
		return true
	}
	if identicalMd5() {
		return
	}

	if _, err = azblob.UploadBufferToBlockBlob(ctx, data, blockBlobUrl, azblob.UploadToBlockBlobOptions{}); err != nil {
		return
	}
	return
}
func CacheFile(url string, auth string, data []byte) (sasUrl string, err error) {
	return CacheFileCtx(context.Background(), url, auth, data)
}

func CreateAuthorization(user, passwd string) string {
	sum := sha256.Sum256([]byte(passwd))
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(
		"%s:%s",
		user,
		base64.StdEncoding.EncodeToString(sum[:]),
	)))
}
