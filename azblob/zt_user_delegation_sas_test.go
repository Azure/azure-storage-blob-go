package azblob

import (
	"bytes"
	"github.com/Azure/azure-pipeline-go/pipeline"
	chk "gopkg.in/check.v1"
	"strings"
	"time"
)

func CreateUserDelegationKey(c *chk.C) (ContainerURL, string, BlockBlobURL, string, UserDelegationCredential, time.Time, pipeline.Pipeline) {
	// Accumulate prerequisite details to create storage etc.
	bsu := getBSU()
	containerURL, containerName := getContainerURL(c, bsu)
	blobURL, blobName := getBlockBlobURL(c, containerURL)
	currentTime := time.Now().UTC()
	ocred, err := getOAuthCredential("")
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline to handle requests
	p := NewPipeline(*ocred, PipelineOptions{})

	// Prepare user delegation key
	bsu = bsu.WithPipeline(p)
	keyInfo := NewKeyInfo(currentTime, currentTime.Add(48*time.Hour))
	budk, err := bsu.GetUserDelegationCredential(ctx, keyInfo, nil, nil) //MUST have TokenCredential
	if err != nil {
		c.Fatal(err)
	}

	return containerURL, containerName, blobURL, blobName, budk, currentTime, p
}

// Attempting to create User Delegation Key SAS with Incorrect Permissions, should return err
func (s *aztestsSuite) TestUserDelegationSASIncorrectPermissions(c *chk.C) {
	_, containerName, _, blobName, cudk, currentTime, _ := CreateUserDelegationKey(c)
	// Prepare User Delegation SAS query for Container
	_, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "rdq",
		ContainerName: containerName,
	}.NewSASQueryParameters(cudk)
	c.Assert(err, chk.NotNil)

	// Prepare User Delegation SAS query for Blob; returns err due to wrong permission
	_, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "rdq",
		ContainerName: containerName,
		BlobName:      blobName,
	}.NewSASQueryParameters(cudk)
	c.Assert(err, chk.NotNil)
}

// Creates a container with no permissions, upload fails due to lack of permissions
func (s *aztestsSuite) TestUserDelegationSASContainerNoPermissions(c *chk.C) {
	containerURL, containerName, _, _, cudk, currentTime, p := CreateUserDelegationKey(c)
	// Prepare User Delegation SAS query
	cSAS, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		ContainerName: containerName,
	}.NewSASQueryParameters(cudk)
	if err != nil {
		c.Fatal(err)
	}

	// Create anonymous pipeline
	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Create the container
	_, err = containerURL.Create(ctx, Metadata{}, PublicAccessNone)
	defer containerURL.Delete(ctx, ContainerAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	// Craft a container URL w/ container UDK SAS
	cURL := containerURL.URL()
	cURL.RawQuery += cSAS.Encode()
	cSASURL := NewContainerURL(cURL, p)

	// Create blob; upload returns err due to lack of permissions
	bblob := cSASURL.NewBlockBlobURL("test")
	_, err = bblob.Upload(ctx, strings.NewReader("hello world!"), BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)
}

// Creates a container with all permissions
func (s *aztestsSuite) TestUserDelegationSASContainerAllPermissions(c *chk.C) {
	containerURL, containerName, _, _, cudk, currentTime, p := CreateUserDelegationKey(c)
	// Prepare User Delegation SAS query
	cSAS, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "racwdxlt",
		ContainerName: containerName,
	}.NewSASQueryParameters(cudk)
	if err != nil {
		c.Fatal(err)
	}

	// Create anonymous pipeline
	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Create the container
	_, err = containerURL.Create(ctx, Metadata{}, PublicAccessNone)
	defer containerURL.Delete(ctx, ContainerAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	// Craft a container URL w/ container UDK SAS
	cURL := containerURL.URL()
	cURL.RawQuery += cSAS.Encode()
	cSASURL := NewContainerURL(cURL, p)

	bblob := cSASURL.NewBlockBlobURL("test")
	_, err = bblob.Upload(ctx, strings.NewReader("hello world!"), BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	if err != nil {
		c.Fatal(err)
	}

	resp, err := bblob.Download(ctx, 0, 0, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	data := &bytes.Buffer{}
	body := resp.Body(RetryReaderOptions{})
	if body == nil {
		c.Fatal("download body was nil")
	}
	_, err = data.ReadFrom(body)
	if err != nil {
		c.Fatal(err)
	}
	err = body.Close()
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(data.String(), chk.Equals, "hello world!")
	_, err = bblob.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
}

//Creates a container and tests permissions by listing blobs
func (s *aztestsSuite) TestUserDelegationSASContainer(c *chk.C) {
	containerURL, containerName, _, _, cudk, currentTime, p := CreateUserDelegationKey(c)
	// Prepare User Delegation SAS query
	cSAS, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "racwdl",
		ContainerName: containerName,
	}.NewSASQueryParameters(cudk)
	if err != nil {
		c.Fatal(err)
	}

	// Create anonymous pipeline
	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Create the container
	_, err = containerURL.Create(ctx, Metadata{}, PublicAccessNone)
	defer containerURL.Delete(ctx, ContainerAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	// Craft a container URL w/ container UDK SAS
	cURL := containerURL.URL()
	cURL.RawQuery += cSAS.Encode()
	cSASURL := NewContainerURL(cURL, p)

	bblob := cSASURL.NewBlockBlobURL("test")
	_, err = bblob.Upload(ctx, strings.NewReader("hello world!"), BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	if err != nil {
		c.Fatal(err)
	}

	resp, err := bblob.Download(ctx, 0, 0, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	data := &bytes.Buffer{}
	body := resp.Body(RetryReaderOptions{})
	if body == nil {
		c.Fatal("download body was nil")
	}
	_, err = data.ReadFrom(body)
	if err != nil {
		c.Fatal(err)
	}
	err = body.Close()
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(data.String(), chk.Equals, "hello world!")
	_, err = bblob.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
}

// Creates a blob with all permissions, takes a snapshot, downloads from snapshot, and deletes from the snapshot w/ the token
func (s *aztestsSuite) TestUserDelegationSASBlobAllPermissions(c *chk.C) {
	containerURL, containerName, blobURL, blobName, budk, currentTime, p := CreateUserDelegationKey(c)

	// Prepare User Delegation SAS query
	bSAS, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "racwdxtmeop",
		ContainerName: containerName,
		BlobName:      blobName,
	}.NewSASQueryParameters(budk)
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline
	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Append User Delegation SAS token to URL
	bSASParts := NewBlobURLParts(blobURL.URL())
	bSASParts.SAS = bSAS
	bSASURL := NewBlockBlobURL(bSASParts.URL(), p)

	// Create container & upload sample data
	_, err = containerURL.Create(ctx, Metadata{}, PublicAccessNone)
	defer containerURL.Delete(ctx, ContainerAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
	data := "Hello World!"
	_, err = blobURL.Upload(ctx, strings.NewReader(data), BlobHTTPHeaders{ContentType: "text/plain"}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	if err != nil {
		c.Fatal(err)
	}

	// Download data via User Delegation SAS URL; must succeed
	downloadResponse, err := bSASURL.Download(ctx, 0, 0, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	if err != nil {
		c.Fatal(err)
	}
	downloadedData := &bytes.Buffer{}
	reader := downloadResponse.Body(RetryReaderOptions{})
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
	_, err = bSASURL.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
}

// Creates a blob, takes a snapshot, downloads from snapshot, and deletes from the snapshot w/ the token
func (s *aztestsSuite) TestUserDelegationSASBlob(c *chk.C) {
	containerURL, containerName, blobURL, blobName, budk, currentTime, p := CreateUserDelegationKey(c)

	// Prepare User Delegation SAS query
	bSAS, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "rd",
		ContainerName: containerName,
		BlobName:      blobName,
	}.NewSASQueryParameters(budk)
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline
	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Append User Delegation SAS token to URL
	bSASParts := NewBlobURLParts(blobURL.URL())
	bSASParts.SAS = bSAS
	bSASURL := NewBlockBlobURL(bSASParts.URL(), p)

	// Create container & upload sample data
	_, err = containerURL.Create(ctx, Metadata{}, PublicAccessNone)
	defer containerURL.Delete(ctx, ContainerAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
	data := "Hello World!"
	_, err = blobURL.Upload(ctx, strings.NewReader(data), BlobHTTPHeaders{ContentType: "text/plain"}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	if err != nil {
		c.Fatal(err)
	}

	// Download data via User Delegation SAS URL; must succeed
	downloadResponse, err := bSASURL.Download(ctx, 0, 0, BlobAccessConditions{}, false, ClientProvidedKeyOptions{})
	if err != nil {
		c.Fatal(err)
	}
	downloadedData := &bytes.Buffer{}
	reader := downloadResponse.Body(RetryReaderOptions{})
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
	_, err = bSASURL.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}
}

// Creates User Delegation SAS with saoid and checks if URL is correctly formed
func (s *aztestsSuite) TestUserDelegationSaoid(c *chk.C) {
	_, containerName, blobURL, blobName, budk, currentTime, p := CreateUserDelegationKey(c)
	saoid := newUUID().String()
	// Prepare User Delegation SAS query
	bSAS, err := BlobSASSignatureValues{
		Protocol:                   SASProtocolHTTPS,
		StartTime:                  currentTime,
		ExpiryTime:                 currentTime.Add(24 * time.Hour),
		Permissions:                "rd",
		ContainerName:              containerName,
		BlobName:                   blobName,
		PreauthorizedAgentObjectId: saoid,
	}.NewSASQueryParameters(budk)
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline
	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Append User Delegation SAS token to URL
	bSASParts := NewBlobURLParts(blobURL.URL())
	bSASParts.SAS = bSAS
	bSASURL := NewBlockBlobURL(bSASParts.URL(), p)

	c.Assert(strings.Contains(bSASURL.blobClient.url.RawQuery, "saoid="+saoid), chk.Equals, true)
}

// Creates User Delegation SAS with suoid and checks if URL is correctly formed
func (s *aztestsSuite) TestUserDelegationSuoid(c *chk.C) {
	_, containerName, blobURL, blobName, budk, currentTime, p := CreateUserDelegationKey(c)
	suoid := newUUID().String()
	// Prepare User Delegation SAS query
	bSAS, err := BlobSASSignatureValues{
		Protocol:        SASProtocolHTTPS,
		StartTime:       currentTime,
		ExpiryTime:      currentTime.Add(24 * time.Hour),
		Permissions:     "rd",
		ContainerName:   containerName,
		BlobName:        blobName,
		SignedUnauthOid: suoid,
	}.NewSASQueryParameters(budk)
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline
	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Append User Delegation SAS token to URL
	bSASParts := NewBlobURLParts(blobURL.URL())
	bSASParts.SAS = bSAS
	bSASURL := NewBlockBlobURL(bSASParts.URL(), p)

	c.Assert(strings.Contains(bSASURL.blobClient.url.RawQuery, "suoid="+suoid), chk.Equals, true)
}

// Creates User Delegation SAS with correlation id and checks if URL is correctly formed
func (s *aztestsSuite) TestUserDelegationCid(c *chk.C) {
	_, containerName, blobURL, blobName, budk, currentTime, p := CreateUserDelegationKey(c)
	cid := newUUID().String()
	// Prepare User Delegation SAS query
	bSAS, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "rd",
		ContainerName: containerName,
		BlobName:      blobName,
		CorrelationId: cid,
	}.NewSASQueryParameters(budk)
	if err != nil {
		c.Fatal(err)
	}

	// Create pipeline
	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Append User Delegation SAS token to URL
	bSASParts := NewBlobURLParts(blobURL.URL())
	bSASParts.SAS = bSAS
	bSASURL := NewBlockBlobURL(bSASParts.URL(), p)

	c.Assert(strings.Contains(bSASURL.blobClient.url.RawQuery, "cid="+cid), chk.Equals, true)
}
