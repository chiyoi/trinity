package request

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/chiyoi/trinity/internal/pkg/logs"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

const (
	sasExpire = time.Hour * 24 * 365 * 10
)

func getBlobCacheURLHandler() (h reqHandler, err error) {
	accountName, err := config.GetErr[string]("AzureStorageAccount")
	if err != nil {
		return
	}
	accountKey, err := config.GetErr[string]("AzureStorageKey")
	if err != nil {
		return
	}
	containerName, err := config.GetErr[string]("FileCacheContainer")
	if err != nil {
		return
	}

	u, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	if err != nil {
		return
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	serviceUrl := azblob.NewServiceURL(*u, pipeline)
	containerUrl := serviceUrl.NewContainerURL(containerName)

	h = func(resp *atmt.Message, req atmt.DataRequest[trinity.Action]) {
		logPrefix := "handle cache file:"
		var args trinity.ArgsGetBlobCacheURL
		if err := json.Unmarshal(req.Args, &args); err != nil {
			logs.Warning(logPrefix, err)
			atmt.Error(resp, atmt.StatusBadRequest)
			return
		}
		_, pass, err := verifyAuth(args.Auth)
		if err != nil {
			logs.Error(logPrefix, err)
			atmt.Error(resp, atmt.StatusInternalServerError)
			return
		}
		if !pass {
			logs.Error(logPrefix, err)
			atmt.Error(resp, atmt.StatusUnauthorized)
			return
		}

		sas := &azblob.BlobSASSignatureValues{
			StartTime:  time.Now().UTC(),
			ExpiryTime: time.Now().UTC().Add(sasExpire),
			Permissions: azblob.BlobSASPermissions{
				Read:   true,
				Create: true,
				Write:  true,
				Delete: true,
			}.String(),
			ContainerName: containerName,
			BlobName:      args.BlobName,
		}
		sasQuery, err := sas.NewSASQueryParameters(credential)
		if err != nil {
			logs.Error(logPrefix, err)
			atmt.Error(resp, atmt.StatusInternalServerError)
			return
		}
		sasURL := containerUrl.NewBlobURL(args.BlobName).URL()
		sasURL.RawQuery += sasQuery.Encode()
		b := atmt.MessageBuilder[atmt.DataResponseBuilder[trinity.ValuesGetBlobCacheURL]]{
			Type: atmt.MessageResponse,
			Data: atmt.DataResponseBuilder[trinity.ValuesGetBlobCacheURL]{
				StatusCode: atmt.StatusOK,
				Values: trinity.ValuesGetBlobCacheURL{
					SasURL: sasURL.String(),
				},
			},
			Content: []atmt.Paragraph{},
		}
		if err = b.Write(resp); err != nil {
			logs.Error(logPrefix, err)
			atmt.Error(resp, atmt.StatusInternalServerError)
			return
		}
	}
	return
}
