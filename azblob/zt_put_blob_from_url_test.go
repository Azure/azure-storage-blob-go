package azblob

import (
	"context"
	"crypto/md5"
	chk "gopkg.in/check.v1"
	"io/ioutil"
	"net/url"
	"time"
)

func (s *aztestsSuite) TestPutBlockBlobFromURLWithTags(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	testSize := 1 * 1024 * 1024 // 1MB
	r, sourceData := getRandomDataAndReader(testSize)
	sourceDataMD5Value := md5.Sum(sourceData)
	srcBlob := container.NewBlockBlobURL("srcBlob")
	destBlob := container.NewBlockBlobURL("destBlob")

	blobTagsMap := BlobTagsMap{
		"Go":         "CPlusPlus",
		"Python":     "CSharp",
		"Javascript": "Android",
	}

	uploadSrcResp, err := srcBlob.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, blobTagsMap, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)

	// Get source blob URL with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Invoke put blob from URL
	resp, err := destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], sourceDataMD5Value[:], DefaultAccessTier, blobTagsMap, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), "")
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ContentMD5(), chk.DeepEquals, sourceDataMD5Value[:])

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)
	c.Assert(len(downloadResp.NewMetadata()), chk.Equals, 1)

	// Edge case 1: Provide bad MD5 and make sure the put fails
	_, badMD5 := getRandomDataAndReader(16)
	_, err = destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, badMD5, badMD5, DefaultAccessTier, blobTagsMap, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)

	// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
	resp, err = destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, nil, DefaultAccessTier, blobTagsMap, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
}

func (s *aztestsSuite) TestPutBlobFromURLWithSASReturnsVID(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	testSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(testSize)
	sourceDataMD5Value := md5.Sum(sourceData)
	ctx := context.Background()
	srcBlob := container.NewBlockBlobURL(generateBlobName())
	destBlob := container.NewBlockBlobURL(generateBlobName())

	uploadSrcResp, err := srcBlob.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)
	c.Assert(uploadSrcResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)

	// Get source blob URL with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Invoke put blob from URL
	resp, err := destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], sourceDataMD5Value[:], DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
	c.Assert(resp.VersionID(), chk.NotNil)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)
	c.Assert(downloadResp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
	c.Assert(len(downloadResp.NewMetadata()), chk.Equals, 1)

	// Edge case 1: Provide bad MD5 and make sure the put fails
	_, badMD5 := getRandomDataAndReader(16)
	_, err = destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, badMD5, badMD5, DefaultAccessTier, BlobTagsMap{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)

	// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
	resp, err = destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, nil, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
	c.Assert(resp.Response().Header.Get("x-ms-version"), chk.Equals, ServiceVersion)
	c.Assert(resp.Response().Header.Get("x-ms-version-id"), chk.NotNil)
}

func (s *aztestsSuite) TestPutBlockBlobFromURL(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	testSize := 8 * 1024 * 1024 // 8MB
	r, sourceData := getRandomDataAndReader(testSize)
	sourceDataMD5Value := md5.Sum(sourceData)
	ctx := context.Background() // Use default Background context
	srcBlob := container.NewBlockBlobURL(generateBlobName())
	destBlob := container.NewBlockBlobURL(generateBlobName())

	// Prepare source blob for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)

	// Get source blob URL with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Invoke put blob from URL.
	resp, err := destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], sourceDataMD5Value[:], DefaultAccessTier, BlobTagsMap{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), "")
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ContentMD5(), chk.DeepEquals, sourceDataMD5Value[:])

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobURL.Download(ctx, 0, CountToEnd, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)

	// Make sure the metadata got copied over
	c.Assert(len(downloadResp.NewMetadata()), chk.Equals, 1)

	// Edge case 1: Provide bad MD5 and make sure the put fails
	_, badMD5 := getRandomDataAndReader(16)
	_, err = destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, badMD5, badMD5, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)

	// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
	resp, err = destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{}, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, nil, DefaultAccessTier, BlobTagsMap{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
}

func (s *aztestsSuite) TestSetTierOnPutBlockBlobFromURL(c *chk.C) {
	bsu := getBSU()

	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	testSize := 1 * 1024 * 1024
	r, sourceData := getRandomDataAndReader(testSize)
	sourceDataMD5Value := md5.Sum(sourceData)
	ctx := context.Background()
	srcBlob := container.NewBlockBlobURL(generateBlobName())

	// Setting blob tier as "cool"
	uploadSrcResp, err := srcBlob.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, AccessTierCool, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)

	// Get source blob URL with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(2 * time.Hour),
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()
	for _, tier := range []AccessTierType{AccessTierArchive, AccessTierCool, AccessTierHot} {
		destBlob := container.NewBlockBlobURL(generateBlobName())
		resp, err := destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlobURLWithSAS, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], sourceDataMD5Value[:], tier, BlobTagsMap{}, ClientProvidedKeyOptions{})
		c.Assert(err, chk.IsNil)
		c.Assert(resp.Response().StatusCode, chk.Equals, 201)

		destBlobPropResp, err := destBlob.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
		c.Assert(err, chk.IsNil)
		c.Assert(destBlobPropResp.AccessTier(), chk.Equals, string(tier))
	}
}

func (s *aztestsSuite) TestPutBlobFromURLWithMissingSAS(c *chk.C) {
	bsu := getBSU()
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	testSize := 8 * 1024 * 1024 // 8MB
	r, sourceData := getRandomDataAndReader(testSize)
	sourceDataMD5Value := md5.Sum(sourceData)
	ctx := context.Background() // Use default Background context
	srcBlob := container.NewBlockBlobURL(generateBlobName())
	destBlob := container.NewBlockBlobURL(generateBlobName())

	// Prepare source blob for put.
	uploadSrcResp, err := srcBlob.Upload(ctx, r, BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.Response().StatusCode, chk.Equals, 201)

	// Invoke put blob from URL with URL without SAS
	resp, err := destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, srcBlob.URL(), Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], sourceDataMD5Value[:], DefaultAccessTier, BlobTagsMap{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)
	c.Assert(resp, chk.IsNil)
}

func (s *aztestsSuite) TestPutBlobFromURLWithIncorrectURL(c *chk.C) {
	bsu := getBSU()
	container, _ := createNewContainer(c, bsu)
	defer delContainer(c, container)

	testSize := 8 * 1024 * 1024 // 8MB
	_, sourceData := getRandomDataAndReader(testSize)
	sourceDataMD5Value := md5.Sum(sourceData)
	ctx := context.Background() // Use default Background context
	destBlob := container.NewBlockBlobURL(generateBlobName())

	// Invoke put blob from URL with incorrect URL
	resp, err := destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, url.URL{}, Metadata{"foo": "bar"}, ModifiedAccessConditions{}, BlobAccessConditions{}, sourceDataMD5Value[:], sourceDataMD5Value[:], DefaultAccessTier, BlobTagsMap{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)
	c.Assert(resp, chk.IsNil)
}
