package azblob_test

import (
	"bytes"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	chk "gopkg.in/check.v1"
	"strings"
	"time"
)

func (s *aztestsSuite) TestIdentitySASUsage(c *chk.C) {
	bsu := getBSU()
	containerURL, containerName := getContainerURL(c, bsu)
	blobURL, blobName := getBlockBlobURL(c, containerURL)

	currentTime := time.Now().UTC()
	accountName, _ := accountInfo()

	ocred, err := getOAuthCredential("")
	if err != nil {
		c.Fatal(err)
	}

	p := azblob.NewPipeline(*ocred, azblob.PipelineOptions{})
	bsu = azblob.NewServiceURL(bsu.URL(), p)

	keyInfo := azblob.NewKeyInfo(currentTime, currentTime.Add(48*time.Hour))

	budk, err := bsu.GetUserDelegationKey(ctx, keyInfo, nil, nil) //MUST have TokenCredential
	if err != nil {
		c.Fatal(err)
	}

	bSAS, err := azblob.BlobSASSignatureValues{
		Protocol:          azblob.SASProtocolHTTPS,
		StartTime:         currentTime,
		ExpiryTime:        currentTime.Add(24 * time.Hour),
		Permissions:       "r",
		ContainerName:     containerName,
		BlobName:          blobName,
		UserDelegationKey: budk,
	}.NewIdentitySASQueryParameters(accountName)
	if err != nil {
		c.Fatal(err)
	}

	//Create pipeline
	p = azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})

	//Create blob SAS Identity URL
	bSASParts := azblob.NewBlobURLParts(blobURL.URL())
	bSASParts.SAS = bSAS
	bsurl := bSASParts.URL()
	fmt.Println(bsurl.String())
	bSASURL := azblob.NewBlockBlobURL(bSASParts.URL(), p)
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		c.Fatal(err)
	}

	data := "Hello World!"

	_, err = blobURL.Upload(ctx, strings.NewReader(data), azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	downloadResponse, err := bSASURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
		c.Fatal(err)
	}

	downloadedData := &bytes.Buffer{}
	reader := downloadResponse.Body(azblob.RetryReaderOptions{})
	downloadedData.ReadFrom(reader)
	reader.Close()

	c.Assert(data, chk.Equals, downloadedData.String())

	_, err = bSASURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
}
