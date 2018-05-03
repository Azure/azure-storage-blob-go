package azblob_test

import (
	chk "gopkg.in/check.v1"
	"github.com/Azure/azure-storage-blob-go/2017-07-29/azblob"
	"bytes"
	"io/ioutil"
)

func (s *aztestsSuite) TestUploadStreamToBlockBlobInChunks(c *chk.C) {
	// Set up test container
	bsu := getBSU()
	containerURL, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL)

	// Set up test blob
	blobURL, _ := getBlockBlobURL(c, containerURL)

	// Create a some data to test the upload stream
	blobSize := 8 * 1024
	blobData := make([]byte, blobSize, blobSize)
	for i := range blobData {
		blobData[i] = byte('a' + i%26)
	}

	// Perform UploadStreamToBlockBlob
	uploadResp, err := azblob.UploadStreamToBlockBlob(ctx, bytes.NewReader(blobData), blobURL,
		azblob.UploadStreamToBlockBlobOptions{BufferSize: 1024, MaxBuffers: 3})

	// Assert that upload was successful
	c.Assert(err, chk.Equals, nil)
	c.Assert(uploadResp.StatusCode(), chk.Equals, 201)

	// Download the blob and assert that the content is correct
	downloadResponse, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	c.Assert(err, chk.IsNil)
	actualBlobData, err := ioutil.ReadAll(downloadResponse.Response().Body)
	c.Assert(len(actualBlobData), chk.Equals, blobSize)
	c.Assert(actualBlobData, chk.DeepEquals, blobData)
}

func (s *aztestsSuite) TestUploadStreamToBlockBlobSingleIO(c *chk.C) {
	// Set up test container
	bsu := getBSU()
	containerURL, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL)

	// Set up test blob
	blobURL, _ := getBlockBlobURL(c, containerURL)

	// Create a some data to test the upload stream
	blobSize := 8 * 1024
	blobData := make([]byte, blobSize, blobSize)
	for i := range blobData {
		blobData[i] = byte('a' + i%26)
	}

	// Perform UploadStreamToBlockBlob
	uploadResp, err := azblob.UploadStreamToBlockBlob(ctx, bytes.NewReader(blobData), blobURL,
		azblob.UploadStreamToBlockBlobOptions{BufferSize: 10 * 1024, MaxBuffers: 3})

	// Assert that upload was successful
	c.Assert(err, chk.Equals, nil)
	c.Assert(uploadResp.StatusCode(), chk.Equals, 201)

	// Download the blob and assert that the content is correct
	downloadResponse, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	c.Assert(err, chk.IsNil)
	actualBlobData, err := ioutil.ReadAll(downloadResponse.Response().Body)
	c.Assert(len(actualBlobData), chk.Equals, blobSize)
	c.Assert(actualBlobData, chk.DeepEquals, blobData)
}