package azblob_test

import (
	"bytes"
	"github.com/Azure/azure-storage-blob-go/azblob"
	chk "gopkg.in/check.v1"
	"strings"
	"time"
)

func (s *aztestsSuite) TestSnapshotSASUsage(c *chk.C) {
	//Generate URLs ----------------------------------------------------------------------------------------------------
	bsu := getBSU()
	containerURL, containerName := getContainerURL(c, bsu)
	blobURL, blobName := getBlockBlobURL(c, containerURL)

	currentTime := time.Now().UTC()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}

	genericReadSAS, err := azblob.AccountSASSignatureValues{
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(48 * time.Hour),
		Permissions:   "racwdl",
		ResourceTypes: "sco",
		Services:      "bfqt",
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	genericParts := azblob.NewBlobURLParts(containerURL.URL())
	genericParts.SAS = genericReadSAS
	genericURL := genericParts.URL()

	//Create pipeline
	p := azblob.NewPipeline(azblob.NewAnonymousCredential(), azblob.PipelineOptions{})
	genericContainer := azblob.NewContainerURL(genericURL, p)

	//Create container
	_, err = genericContainer.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)

	if err != nil {
		c.Fatal(err)
	}

	//Create file in container, download from snapshot to test. --------------------------------------------------------
	burl := genericContainer.NewBlockBlobURL(blobName)
	data := "Hello world!"

	_, err = burl.Upload(ctx, strings.NewReader(data), azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	//Create a snapshot & URL
	createSnapshot, err := burl.CreateSnapshot(ctx, azblob.Metadata{}, azblob.BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	snapTime, err := time.Parse(azblob.SnapshotTimeFormat, createSnapshot.Snapshot())
	if err != nil {
		c.Fatal(err)
	}

	snapSASQueryParams, err := azblob.BlobSASSignatureValues{
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(48 * time.Hour),
		SnapshotTime:  snapTime,
		Permissions:   "racwd",
		ContainerName: containerName,
		BlobName:      blobName,
		Protocol:      azblob.SASProtocolHTTPS,
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	snapParts := azblob.NewBlobURLParts(blobURL.URL())
	snapParts.SAS = snapSASQueryParams
	sburl := azblob.NewBlockBlobURL(snapParts.URL(), p)

	//Test the snapshot
	downloadResponse, err := sburl.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
		c.Fatal(err)
	}

	downloadedData := &bytes.Buffer{}
	reader := downloadResponse.Body(azblob.RetryReaderOptions{})
	downloadedData.ReadFrom(reader)
	reader.Close()

	c.Assert(data, chk.Equals, downloadedData.String())

	//Try (and fail to) delete from snapshot ---------------------------------------------------------------------------
	deleteResponse, err := sburl.Delete(ctx, azblob.DeleteSnapshotsOptionOnly, azblob.BlobAccessConditions{})
	if err == nil { //This absolutely SHOULD error out.
		c.Fatal(deleteResponse)
	}
}
