package azblob_test

import (
	"bytes"
	"github.com/Azure/azure-storage-blob-go/azblob"
	chk "gopkg.in/check.v1"
	"strings"
	"time"
)

//Creates a container and tests permissions by listing blobs
func (s *aztestsSuite) TestUserDelegationSASContainer(c *chk.C) {
	bsu := getBSU()
	containerURL, containerName := getContainerURL(c, bsu)
	currentTime := time.Now().UTC()
	accountName, _ := accountInfo()
	ocred, err := getOAuthCredential("")
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline w/ OAuth to handle user delegation key obtaining
	p := azblob.NewPipeline(*ocred, azblob.PipelineOptions{})

	bsu = bsu.WithPipeline(p)
	keyInfo := azblob.NewKeyInfo(currentTime, currentTime.Add(48*time.Hour))
	cudk, err := bsu.GetUserDelegationKey(ctx, keyInfo, nil, nil)
	if err != nil {
		c.Fatal(err)
	}

	cSAS, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "racwdl",
		ContainerName: containerName,
	}.NewSASQueryParameters(nil, accountName, &cudk)

	// Create anonymous pipeline
	p = azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})

	// Create the container
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	defer containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	// Craft a container URL w/ container UDK SAS
	cURL := containerURL.URL()
	cURL.RawQuery += cSAS.Encode()
	cSASURL := azblob.NewContainerURL(cURL, p)

	// Lazily test the container by listing blobs
	_, err = cSASURL.ListBlobsFlatSegment(ctx, azblob.Marker{}, azblob.ListBlobsSegmentOptions{})
	if err != nil {
		c.Fatal(err)
	}
}

// Creates a blob, takes a snapshot, downloads from snapshot, and deletes from the snapshot w/ the token
func (s *aztestsSuite) TestUserDelegationSAS(c *chk.C) {
	// Accumulate prerequisite details to create storage etc.
	bsu := getBSU()
	containerURL, containerName := getContainerURL(c, bsu)
	blobURL, blobName := getBlockBlobURL(c, containerURL)
	currentTime := time.Now().UTC()
	accountName, _ := accountInfo()
	ocred, err := getOAuthCredential("")
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline to handle requests
	p := azblob.NewPipeline(*ocred, azblob.PipelineOptions{})

	// Prepare user delegation key
	bsu = bsu.WithPipeline(p)
	keyInfo := azblob.NewKeyInfo(currentTime, currentTime.Add(48*time.Hour))
	budk, err := bsu.GetUserDelegationKey(ctx, keyInfo, nil, nil) //MUST have TokenCredential
	if err != nil {
		c.Fatal(err)
	}

	// Prepare User Delegation SAS query
	bSAS, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "rd",
		ContainerName: containerName,
		BlobName:      blobName,
	}.NewSASQueryParameters(nil, accountName, &budk)
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline
	p = azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})

	// Append User Delegation SAS token to URL
	bSASParts := azblob.NewBlobURLParts(blobURL.URL())
	bSASParts.SAS = bSAS
	bSASURL := azblob.NewBlockBlobURL(bSASParts.URL(), p)

	// Create container & upload sample data
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	defer containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
	data := "Hello World!"
	_, err = blobURL.Upload(ctx, strings.NewReader(data), azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	// Download data via User Delegation SAS URL; must succeed
	downloadResponse, err := bSASURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
		c.Fatal(err)
	}
	downloadedData := &bytes.Buffer{}
	reader := downloadResponse.Body(azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		c.Fatal(err)
	}
	err = reader.Close()
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(data, chk.Equals, downloadedData.String())

	// Delete the item using the User Delegation SAS URL; must succeed
	_, err = bSASURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
}
