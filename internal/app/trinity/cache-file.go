package trinity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/chiyoi/trinity/internal/app/trinity/config"
)

var (
	containerName string
	containerUrl  azblob.ContainerURL

	credential *azblob.SharedKeyCredential
)

const (
	fileCacheSasExpireDelay = time.Hour * 24 * 365 * 10
)

func init() {
	accountName, accountKey :=
		config.Get[string]("AzureStorageAccount"),
		config.Get[string]("AzureStorageKey")
	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	if err != nil {
		return
	}
	if credential, err = azblob.NewSharedKeyCredential(accountName, accountKey); err != nil {
		return
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	serviceUrl := azblob.NewServiceURL(*u, pipeline)

	containerName = config.Get[string]("FileCacheContainer")
	containerUrl = serviceUrl.NewContainerURL(containerName)
}

func handleCacheFile(baseCtx context.Context, w http.ResponseWriter, req Request) {
	var reqData ReqDataCacheFile
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		badRequestCallback(w, err)
	}
	blobName := reqData.Sha256SumHex
	blobUrl := containerUrl.NewBlobURL(blobName)

	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()
	if _, err := blobUrl.GetProperties(ctx, azblob.BlobAccessConditions{}, azblob.ClientProvidedKeyOptions{}); err != nil {
		if responseError, ok := err.(azblob.ResponseError); !ok || responseError.Response().StatusCode == http.StatusNotFound {
			internalServerErrorCallback(w, err)
			return
		}
	}

	sas := &azblob.BlobSASSignatureValues{
		StartTime:  time.Now().UTC(),
		ExpiryTime: time.Now().UTC().Add(fileCacheSasExpireDelay),
		Permissions: azblob.BlobSASPermissions{
			Read:   true,
			Create: true,
			Write:  true,
			Delete: true,
		}.String(),
		ContainerName: containerName,
		BlobName:      blobName,
	}
	sasQuery, err := sas.NewSASQueryParameters(credential)
	if err != nil {
		return
	}

	u := blobUrl.URL()
	u.RawQuery += sasQuery.Encode()

	respData := RespDataCacheFile{
		SasURL: u.String(),
	}
	resp := Response{
		StatusCode: StatusOK,
		Data:       respData,
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		internalServerErrorCallback(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
