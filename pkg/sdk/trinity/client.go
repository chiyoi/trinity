package trinity

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/chiyoi/trinity/pkg/atmt"
)

func Request[Args any, Values any](url string, action Action, args Args, content atmt.Content) (vals Values, ret atmt.Content, err error) {
	return RequestCtx[Args, Values](context.Background(), url, action, args, content)
}
func RequestCtx[Args any, Values any](ctx context.Context, url string, action Action, args Args, content atmt.Content) (vals Values, ret atmt.Content, err error) {
	b := atmt.MessageBuilder[atmt.DataRequestBuilder[Action, Args]]{
		Type: atmt.MessageRequest,
		Data: atmt.DataRequestBuilder[Action, Args]{
			Action: action,
			Args:   args,
		},
		Content: content,
	}
	req, err := b.Message()
	if err != nil {
		return
	}

	resp, err := atmt.SendMessageCtx(ctx, url, req)
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
		return
	}
	if err = json.Unmarshal(data.Values, &vals); err != nil {
		return
	}
	ret = resp.Content
	return
}

const (
	reqTimeout = time.Second * 10
)

func CacheBlob(url string, auth string, b []byte) (sasURL string, err error) {
	md5Sum := md5.Sum(b)
	vals, _, err := Request[ArgsGetBlobCacheURL, ValuesGetBlobCacheURL](
		url,
		ActionGetBlobCacheURL,
		ArgsGetBlobCacheURL{
			Auth:     auth,
			BlobName: fmt.Sprintf("%x", md5Sum),
		},
		nil,
	)
	if err != nil {
		return
	}
	sasURL = vals.SasURL

	bg := context.Background()
	credential := azblob.NewAnonymousCredential()
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	u, err := urlpkg.Parse(sasURL)
	if err != nil {
		return
	}
	blockBlobURL := azblob.NewBlockBlobURL(*u, pipeline)
	ctx, cancel := context.WithTimeout(bg, reqTimeout)
	defer cancel()
	properties, err := blockBlobURL.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{})
	if err == nil && *(*[md5.Size]byte)(properties.ContentMD5()) == md5Sum {
		return
	}
	if responseErr, ok := err.(azblob.ResponseError); err != nil && (!ok || responseErr.Response().StatusCode != http.StatusNotFound) {
		return
	}
	_, err = azblob.UploadBufferToBlockBlob(bg, b, blockBlobURL, azblob.UploadToBlockBlobOptions{})
	return
}
