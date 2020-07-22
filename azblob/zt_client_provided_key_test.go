package azblob

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

var testEncryptedKey = "MDEyMzQ1NjcwMTIzNDU2NzAxMjM0NTY3MDEyMzQ1Njc="
var testEncryptedHash = "3QFFFpRA5+XANHqwwbT4yXDmrT/2JaLt/FKHjzhOdoE="
var testEncryptedScope = ""
var testCPK = InitClientProvidedKeyOptions(&testEncryptedKey, &testEncryptedHash, &testEncryptedScope)


func blockIDBinaryToBase64(blockID []byte) string {
	return base64.StdEncoding.EncodeToString(blockID)
}

func blockIDBase64ToBinary(blockID string) []byte {
	binary, _ := base64.StdEncoding.DecodeString(blockID)
	return binary
}

// blockIDIntToBase64 functions convert an int block ID to a base-64 string and vice versa
func blockIDIntToBase64(blockID int) string {
	binaryBlockID := (&[4]byte{})[:] // All block IDs are 4 bytes long
	binary.LittleEndian.PutUint32(binaryBlockID, uint32(blockID))
	return blockIDBinaryToBase64(binaryBlockID)
}
func blockIDBase64ToInt(blockID string) int {
	blockIDBase64ToBinary(blockID)
	return int(binary.LittleEndian.Uint32(blockIDBase64ToBinary(blockID)))
}

func (s *aztestsSuite) TestPutBlockAndPutBlockListWithCPK(c *chk.C) {
	bsu := getBSU()
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	blobURL := container.NewBlockBlobURL(generateBlobName())

	words := []string{"AAA ", "BBB ", "CCC "}
	base64BlockIDs := make([]string, len(words))
	for index, word := range words {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		_, err := blobURL.StageBlock(ctx, base64BlockIDs[index], strings.NewReader(word), LeaseAccessConditions{}, nil, *testCPK)
		if err != nil {
			c.Fail()
		}
	}

	resp, err := blobURL.CommitBlockList(ctx, base64BlockIDs, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, *testCPK)
	if err != nil {
		c.Fail()
	}

	c.Assert(resp.ETag(), chk.NotNil)
	c.Assert(resp.LastModified(), chk.NotNil)
	c.Assert(resp.IsServerEncrypted(), chk.Equals, "true")
	c.Assert(resp.EncryptionKeySha256(), chk.DeepEquals, *(testCPK.EncryptionKeySha256))

	// Get blob content without encryption key should fail the request.
	_, err = blobURL.Download(ctx, 0, 0, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	if err == nil {
		c.Fail()
	}

	getResp, err := blobURL.Download(ctx, 0, 0, BlobAccessConditions{}, false, *testCPK)
	if err != nil {
		c.Fail()
	}
	b := bytes.Buffer{}
	reader := getResp.Body(RetryReaderOptions{})
	b.ReadFrom(reader)
	reader.Close() // The client must close the response body when finished with it
	// fmt.Println(b.String())
	c.Assert(b.String(), chk.Equals, "AAA BBB CCC ")
	c.Assert(getResp.ETag(), chk.Equals, resp.ETag())
	c.Assert(getResp.LastModified(), chk.DeepEquals, resp.LastModified())
	c.Assert(getResp.r.EncryptionKeySha256(), chk.Equals, *(testCPK.EncryptionKeySha256))
}

func (s *aztestsSuite) TestPutBlockFromURLAndCommitWithCPK(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	testSize := 2 * 1024 // 2KB
	r, sourceData := getRandomDataAndReader(testSize)
	ctx := context.Background()
	srcBlob := container.NewBlockBlobURL(generateBlobName())

	uploadSrcResp, err := srcBlob.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)

	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	destBlob := container.NewBlockBlobURL(generateBlobName())
	blockID1, blockID2 := blockIDIntToBase64(0), blockIDIntToBase64(1)
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, 1*1024, LeaseAccessConditions{}, ModifiedAccessConditions{}, *testCPK)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp1.Response().StatusCode, chk.Equals, 201)
	c.Assert(stageResp1.ContentMD5(), chk.Not(chk.Equals), "")
	c.Assert(stageResp1.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Version(), chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Date().IsZero(), chk.Equals, false)
	c.Assert(stageResp1.IsServerEncrypted(), chk.Equals, "true")

	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 1*1024, CountToEnd, LeaseAccessConditions{}, ModifiedAccessConditions{}, *testCPK)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp2.Response().StatusCode, chk.Equals, 201)
	c.Assert(stageResp2.ContentMD5(), chk.Not(chk.Equals), "")
	c.Assert(stageResp2.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Version(), chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Date().IsZero(), chk.Equals, false)
	c.Assert(stageResp2.IsServerEncrypted(), chk.Equals, "true")

	blockList, err := destBlob.GetBlockList(ctx, BlockListAll, LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.Response().StatusCode, chk.Equals, 200)
	c.Assert(blockList.UncommittedBlocks, chk.HasLen, 2)
	c.Assert(blockList.CommittedBlocks, chk.HasLen, 0)

	listResp, err := destBlob.CommitBlockList(ctx, []string{blockID1, blockID2}, BlobHTTPHeaders{}, nil, BlobAccessConditions{}, *testCPK)
	c.Assert(err, chk.IsNil)
	c.Assert(listResp.Response().StatusCode, chk.Equals, 201)
	c.Assert(listResp.IsServerEncrypted(), chk.Equals, "true")

	blockList, err = destBlob.GetBlockList(ctx, BlockListAll, LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.Response().StatusCode, chk.Equals, 200)
	c.Assert(blockList.UncommittedBlocks, chk.HasLen, 0)
	c.Assert(blockList.CommittedBlocks, chk.HasLen, 2)

	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	if err == nil {
		c.Fail()
	}

	downloadResp, err = destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, *testCPK)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)
}

func (s *aztestsSuite) TestAppendBlockWithCPK(c *chk.C) {
	bsu := getBSU()
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	appendBlobURL := container.NewAppendBlobURL(generateBlobName())

	resp, err := appendBlobURL.Create(context.Background(), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, *testCPK)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.StatusCode(), chk.Equals, 201)

	words := []string{"AAA ", "BBB ", "CCC "}
	for index, word := range words {
		resp, err := appendBlobURL.AppendBlock(context.Background(), strings.NewReader(word), AppendBlobAccessConditions{}, nil, *testCPK)
		if err != nil {
			c.Fail()
		}
		c.Assert(err, chk.IsNil)
		c.Assert(resp.Response().StatusCode, chk.Equals, 201)
		c.Assert(resp.BlobAppendOffset(), chk.Equals, strconv.Itoa(index*4))
		c.Assert(resp.BlobCommittedBlockCount(), chk.Equals, int32(index+1))
		c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
		c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
		c.Assert(resp.ContentMD5(), chk.Not(chk.Equals), "")
		c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
		c.Assert(resp.Version(), chk.Not(chk.Equals), "")
		c.Assert(resp.Date().IsZero(), chk.Equals, false)
		c.Assert(resp.IsServerEncrypted(), chk.Equals, "true")
		c.Assert(resp.EncryptionKeySha256(), chk.Equals, *(testCPK.EncryptionKeySha256))
	}

	_, err = appendBlobURL.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	if err == nil {
		c.Fail()
	}

	downloadResp, err := appendBlobURL.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, *testCPK)
	c.Assert(err, chk.IsNil)

	data, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(string(data), chk.DeepEquals, "AAA BBB CCC ")
}

func (s *aztestsSuite) TestAppendBlockFromURLWithCPK(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	testSize := 2 * 1024 * 1024 // 2MB
	r, sourceData := getRandomDataAndReader(testSize)
	ctx := context.Background() // Use default Background context
	srcBlob := container.NewAppendBlobURL(generateName("src"))
	destBlob := container.NewAppendBlobURL(generateName("dest"))

	cResp1, err := srcBlob.Create(context.Background(), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(cResp1.StatusCode(), chk.Equals, 201)
	
	resp, err := srcBlob.AppendBlock(context.Background(), r, AppendBlobAccessConditions{}, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.ContentMD5(), chk.Not(chk.Equals), "")

	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS, 
		ExpiryTime:    time.Now().UTC().Add(1 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	cResp2, err := destBlob.Create(context.Background(), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, *testCPK)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp2.StatusCode(), chk.Equals, 201)
	
	appendResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, 0, int64(testSize), AppendBlobAccessConditions{}, ModifiedAccessConditions{}, nil, *testCPK)
	c.Assert(err, chk.IsNil)
	c.Assert(appendResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(appendResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(appendResp.IsServerEncrypted(), chk.Equals, "true")

	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	if err == nil {
		c.Fail()
	}

	downloadResp, err = destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, *testCPK)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)
}
