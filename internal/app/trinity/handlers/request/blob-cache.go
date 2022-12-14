package request

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/chiyoi/neko03/pkg/logs"
	"github.com/chiyoi/trinity/internal/app/trinity/config"
	"github.com/chiyoi/trinity/pkg/atmt"
	"github.com/chiyoi/trinity/pkg/sdk/trinity"
)

const (
	sasAlive = time.Hour * 24 * 365 * 10
)

func getBlobCacheURLHandler() (h reqHandler) {
	accountName, accountKey, containerName :=
		config.Get[string]("AzureStorageAccount"),
		config.Get[string]("AzureStorageKey"),
		config.Get[string]("ContainerBlobCache")

	u, err := url.Parse(fmt.Sprintf(
		"https://%s.blob.core.windows.net",
		accountName,
	))
	if err != nil {
		logs.Panic(err)
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		logs.Panic(err)
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	serviceUrl := azblob.NewServiceURL(*u, pipeline)
	containerUrl := serviceUrl.NewContainerURL(containerName)

	return func(resp *atmt.Message, req atmt.DataRequest[trinity.Action]) {
		var args trinity.ArgsGetBlobCacheURL
		if err := json.Unmarshal(req.Args, &args); err != nil {
			logs.Warning(err)
			atmt.Error(resp, atmt.StatusBadRequest)
			return
		}
		_, pass, err := verifyAuth(resp, args.Auth)
		if err != nil {
			return
		}
		if !pass {
			atmt.Error(resp, atmt.StatusUnauthorized)
			return
		}

		sas := &azblob.BlobSASSignatureValues{
			StartTime:  time.Now().UTC(),
			ExpiryTime: time.Now().UTC().Add(sasAlive),
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
			logs.Error(err)
			atmt.InternalServerError(resp)
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
			logs.Error(err)
			atmt.InternalServerError(resp)
			return
		}
	}
}
