package azblob

import (
	"context"
	"strings"
	"time"

	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

func (s *aztestsSuite) TestGetAccountInfo(c *chk.C) {
	sa := getBSU()

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sAccInfo, err := sa.GetAccountInfo(context.Background())
	c.Assert(err, chk.IsNil)
	c.Assert(*sAccInfo, chk.Not(chk.DeepEquals), ServiceGetAccountInfoResponse{})

	// Test on a container
	cURL := sa.NewContainerURL(generateContainerName())
	defer deleteContainer(c, cURL, false)
	_, err = cURL.Create(ctx, Metadata{}, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	cAccInfo, err := cURL.GetAccountInfo(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(*cAccInfo, chk.Not(chk.DeepEquals), ContainerGetAccountInfoResponse{})

	// test on a block blob URL. They all call the same thing on the base URL, so only one test is needed for that.
	bbURL := cURL.NewBlockBlobURL(generateBlobName())
	_, err = bbURL.Upload(ctx, strings.NewReader("blah"), BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)
	bAccInfo, err := bbURL.GetAccountInfo(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(*bAccInfo, chk.Not(chk.DeepEquals), BlobGetAccountInfoResponse{})
}

func (s *aztestsSuite) TestListContainers(c *chk.C) {
	sa := getBSU()
	resp, err := sa.ListContainersSegment(context.Background(), Marker{}, ListContainersSegmentOptions{Prefix: containerPrefix})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
	c.Assert(len(resp.ContainerItems) >= 0, chk.Equals, true)
	c.Assert(resp.ServiceEndpoint, chk.NotNil)

	container, name := createNewContainer(c, sa)
	defer deleteContainer(c, container, false)

	md := Metadata{
		"foo": "foovalue",
		"bar": "barvalue",
	}
	_, err = container.SetMetadata(context.Background(), md, ContainerAccessConditions{})
	c.Assert(err, chk.IsNil)

	resp, err = sa.ListContainersSegment(context.Background(), Marker{}, ListContainersSegmentOptions{Detail: ListContainersDetail{Metadata: true}, Prefix: name})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.ContainerItems, chk.HasLen, 1)
	c.Assert(resp.ContainerItems[0].Name, chk.NotNil)
	c.Assert(resp.ContainerItems[0].Properties, chk.NotNil)
	c.Assert(resp.ContainerItems[0].Properties.LastModified, chk.NotNil)
	c.Assert(resp.ContainerItems[0].Properties.Etag, chk.NotNil)
	c.Assert(resp.ContainerItems[0].Properties.LeaseStatus, chk.Equals, LeaseStatusUnlocked)
	c.Assert(resp.ContainerItems[0].Properties.LeaseState, chk.Equals, LeaseStateAvailable)
	c.Assert(string(resp.ContainerItems[0].Properties.LeaseDuration), chk.Equals, "")
	c.Assert(string(resp.ContainerItems[0].Properties.PublicAccess), chk.Equals, string(PublicAccessNone))
	c.Assert(resp.ContainerItems[0].Metadata, chk.DeepEquals, md)
}

func (s *aztestsSuite) TestListContainersPaged(c *chk.C) {
	sa := getBSU()

	const numContainers = 4
	const maxResultsPerPage = 2
	const pagedContainersPrefix = "azblobspagedtest"

	containers := make([]ContainerURL, numContainers)
	for i := 0; i < numContainers; i++ {
		containers[i], _ = createNewContainerWithSuffix(c, sa, pagedContainersPrefix)
	}

	defer func() {
		for i := range containers {
			deleteContainer(c, containers[i], false)
		}
	}()

	marker := Marker{}
	iterations := numContainers / maxResultsPerPage

	for i := 0; i < iterations; i++ {
		resp, err := sa.ListContainersSegment(context.Background(), marker, ListContainersSegmentOptions{MaxResults: maxResultsPerPage, Prefix: containerPrefix + pagedContainersPrefix})
		c.Assert(err, chk.IsNil)
		c.Assert(resp.ContainerItems, chk.HasLen, maxResultsPerPage)

		hasMore := i < iterations-1
		c.Assert(resp.NextMarker.NotDone(), chk.Equals, hasMore)
		marker = resp.NextMarker
	}
}

func (s *aztestsSuite) TestAccountListContainersEmptyPrefix(c *chk.C) {
	bsu := getBSU()
	containerURL1, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL1, false)
	containerURL2, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL2, false)

	response, err := bsu.ListContainersSegment(ctx, Marker{}, ListContainersSegmentOptions{})

	c.Assert(err, chk.IsNil)
	c.Assert(len(response.ContainerItems) >= 2, chk.Equals, true) // The response should contain at least the two created containers. Probably many more
}

func (s *aztestsSuite) TestAccountListContainersIncludeTypeMetadata(c *chk.C) {
	bsu := getBSU()
	containerURLNoMetadata, nameNoMetadata := createNewContainerWithSuffix(c, bsu, "a")
	defer deleteContainer(c, containerURLNoMetadata, false)
	containerURLMetadata, nameMetadata := createNewContainerWithSuffix(c, bsu, "b")
	defer deleteContainer(c, containerURLMetadata, false)

	// Test on containers with and without metadata
	_, err := containerURLMetadata.SetMetadata(ctx, basicMetadata, ContainerAccessConditions{})
	c.Assert(err, chk.IsNil)

	// Also validates not specifying MaxResults
	response, err := bsu.ListContainersSegment(ctx, Marker{},
		ListContainersSegmentOptions{Prefix: containerPrefix, Detail: ListContainersDetail{Metadata: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(response.ContainerItems[0].Name, chk.Equals, nameNoMetadata)
	c.Assert(response.ContainerItems[0].Metadata, chk.HasLen, 0)
	c.Assert(response.ContainerItems[1].Name, chk.Equals, nameMetadata)
	c.Assert(response.ContainerItems[1].Metadata, chk.DeepEquals, basicMetadata)
}

func (s *aztestsSuite) TestAccountListContainersMaxResultsNegative(c *chk.C) {
	bsu := getBSU()
	containerURL, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerURL, false)
	_, err := bsu.ListContainersSegment(ctx,
		Marker{}, *(&ListContainersSegmentOptions{Prefix: containerPrefix, MaxResults: -2}))
	c.Assert(err, chk.Not(chk.IsNil))
}

func (s *aztestsSuite) TestAccountListContainersMaxResultsZero(c *chk.C) {
	bsu := getBSU()
	containerURL, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerURL, false)

	// Max Results = 0 means the value will be ignored, the header not set, and the server default used
	resp, err := bsu.ListContainersSegment(ctx,
		Marker{}, *(&ListContainersSegmentOptions{Prefix: containerPrefix, MaxResults: 0}))

	c.Assert(err, chk.IsNil)
	// There could be existing container
	c.Assert(len(resp.ContainerItems) >= 1, chk.Equals, true)
}

func (s *aztestsSuite) TestAccountListContainersMaxResultsExact(c *chk.C) {
	// If this test fails, ensure there are no extra containers prefixed with go in the account. These may be left over if a test is interrupted.
	bsu := getBSU()
	containerURL1, containerName1 := createNewContainerWithSuffix(c, bsu, "a")
	defer deleteContainer(c, containerURL1, false)
	containerURL2, containerName2 := createNewContainerWithSuffix(c, bsu, "b")
	defer deleteContainer(c, containerURL2, false)

	response, err := bsu.ListContainersSegment(ctx,
		Marker{}, *(&ListContainersSegmentOptions{Prefix: containerPrefix, MaxResults: 2}))

	c.Assert(err, chk.IsNil)
	c.Assert(response.ContainerItems, chk.HasLen, 2)
	c.Assert(response.ContainerItems[0].Name, chk.Equals, containerName1)
	c.Assert(response.ContainerItems[1].Name, chk.Equals, containerName2)
}

func (s *aztestsSuite) TestAccountListContainersMaxResultsInsufficient(c *chk.C) {
	bsu := getBSU()
	containerURL1, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL1, false)
	containerURL2, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL2, false)

	response, err := bsu.ListContainersSegment(ctx, Marker{},
		*(&ListContainersSegmentOptions{Prefix: containerPrefix, MaxResults: 1}))

	c.Assert(err, chk.IsNil)
	c.Assert(len(response.ContainerItems), chk.Equals, 1)
}

func (s *aztestsSuite) TestAccountListContainersMaxResultsSufficient(c *chk.C) {
	bsu := getBSU()
	containerURL1, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL1, false)
	containerURL2, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerURL2, false)

	response, err := bsu.ListContainersSegment(ctx, Marker{},
		*(&ListContainersSegmentOptions{Prefix: containerPrefix, MaxResults: 3}))

	c.Assert(err, chk.IsNil)

	// This case could be instable, there could be existing containers, so the count should be >= 2
	c.Assert(len(response.ContainerItems) >= 2, chk.Equals, true)
}

func CreateBlobWithRetentionPolicy(c *chk.C) (BlockBlobURL, ContainerURL) {
	bsu, err := getAlternateBSU()
	c.Assert(err, chk.IsNil)

	days := int32(5)
	allowDelete := true
	_, err = bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true, Days: &days, AllowPermanentDelete: &allowDelete}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	// create container and blobs
	containerURL, _ := createNewContainer(c, bsu)
	testSize := 8 * 1024
	r, _ := getRandomDataAndReader(testSize)
	blobURL, _ := getBlockBlobURL(c, containerURL)

	cResp, err := blobURL.Upload(ctx, r, BlobHTTPHeaders{}, nil, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.StatusCode(), chk.Equals, 201)

	return blobURL, containerURL
}

func (s *aztestsSuite) TestUndelete(c *chk.C) {
	blobURL, containerURL := CreateBlobWithRetentionPolicy(c)
	defer deleteContainer(c, containerURL, false)

	// Soft delete blob
	_, err := blobURL.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{}) //soft delete
	c.Assert(err, chk.IsNil)

	// Check that blob has been soft deleted
	_, err = blobURL.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)

	// Undelete soft deleted blob
	undelResp, err := blobURL.Undelete(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(undelResp, chk.NotNil)
	c.Assert(undelResp.StatusCode(), chk.Equals, 200)

	blobProp, err := blobURL.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
	c.Assert(blobProp.StatusCode(), chk.Equals, 200)
}

/*func (s *aztestsSuite) TestPermanentDelete(c *chk.C) {
	blobURL, containerURL := CreateBlobWithRetentionPolicy(c)
	defer deleteContainer(c, containerURL, false)

	// Create snapshot for blob
	snapResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
	c.Assert(snapResp, chk.NotNil)
	c.Assert(err, chk.IsNil)

	// Check snapshot and blob exist
	listBlobResp1, err := containerURL.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(listBlobResp1.Segment.BlobItems, chk.HasLen, 2)

	// Soft delete snapshot
	snapshotBlob := blobURL.WithSnapshot(snapResp.Snapshot())
	_, err = snapshotBlob.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	// Check that both blobs and snapshot exist
	listBlobResp2, err := containerURL.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Deleted: true, Snapshots: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(listBlobResp2.Segment.BlobItems, chk.HasLen, 2)

	// Permanent delete snapshot
	delResp, err := snapshotBlob.PermanentDelete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(delResp, chk.NotNil)
	c.Assert(delResp.StatusCode(), chk.Equals, 202)

	// Check that snapshot has been deleted
	spBlobResp, err := snapshotBlob.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)
	c.Assert(spBlobResp, chk.IsNil)

	// Check that only blob exists
	listBlobResp3, err := containerURL.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Deleted: true, Snapshots: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(listBlobResp3.Segment.BlobItems, chk.HasLen, 1)
	c.Assert(listBlobResp3.Segment.BlobItems[0].Deleted, chk.Equals, false)
}*/

/*func (s *aztestsSuite) TestPermanentDeleteAccountSAS(c *chk.C) {
	accountName, accountKey := accountInfo()
	credential, err := NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		c.Fail()
	}

	sasQueryParams, err := AccountSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour),
		Permissions:   AccountSASPermissions{Read: true, List: true, Write: true, Create: true, PermanentDelete: true, Delete: true}.String(),
		Services:      AccountSASServices{Blob: true}.String(),
		ResourceTypes: AccountSASResourceTypes{Service: true, Container: true, Object: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		log.Fatal(err)
	}

	qp := sasQueryParams.Encode()
	urlToSendToSomeone := fmt.Sprintf("https://%s.blob.core.windows.net?%s", accountName, qp)
	u, _ := url.Parse(urlToSendToSomeone)
	serviceURL := NewServiceURL(*u, NewPipeline(NewAnonymousCredential(), PipelineOptions{}))
	days := int32(5)
	allowDelete := true
	_, err = serviceURL.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true, Days: &days, AllowPermanentDelete: &allowDelete}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed to be enough.
	time.Sleep(time.Second * 30)

	containerName := generateContainerName()
	containerURL := serviceURL.NewContainerURL(containerName)
	_, err = containerURL.Create(ctx, Metadata{}, PublicAccessNone)
	defer containerURL.Delete(ctx, ContainerAccessConditions{})
	if err != nil {
		c.Fatal(err)
	}

	blobURL := containerURL.NewBlockBlobURL("temp")
	_, err = blobURL.Upload(ctx, bytes.NewReader([]byte("random data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil, ClientProvidedKeyOptions{})
	if err != nil {
		c.Fail()
	}

	// Create snapshot for blob
	snapResp, err := blobURL.CreateSnapshot(ctx, Metadata{}, BlobAccessConditions{}, ClientProvidedKeyOptions{})
	c.Assert(snapResp, chk.NotNil)
	c.Assert(err, chk.IsNil)

	// Check snapshot and blob exist
	listBlobResp1, err := containerURL.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(listBlobResp1.Segment.BlobItems, chk.HasLen, 2)

	// Soft delete snapshot
	snapshotBlob := blobURL.WithSnapshot(snapResp.Snapshot())
	_, err = snapshotBlob.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	// Check that both blobs and snapshot exist
	listBlobResp2, err := containerURL.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Deleted: true, Snapshots: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(listBlobResp2.Segment.BlobItems, chk.HasLen, 2)

	// Permanent delete snapshot
	delResp, err := snapshotBlob.PermanentDelete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(delResp, chk.NotNil)
	c.Assert(delResp.StatusCode(), chk.Equals, 202)

	// Check that snapshot has been deleted
	spBlobResp, err := snapshotBlob.GetProperties(ctx, BlobAccessConditions{}, ClientProvidedKeyOptions{})
	c.Assert(err, chk.NotNil)
	c.Assert(spBlobResp, chk.IsNil)

	// Check that only blob exists
	listBlobResp3, err := containerURL.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{Details: BlobListingDetails{Deleted: true, Snapshots: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(listBlobResp3.Segment.BlobItems, chk.HasLen, 1)
	c.Assert(listBlobResp3.Segment.BlobItems[0].Deleted, chk.Equals, false)
}*/

func (s *aztestsSuite) TestAccountDeleteRetentionPolicy(c *chk.C) {
	bsu := getBSU()

	days := int32(5)
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true, Days: &days}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed t be enough.
	time.Sleep(time.Second * 30)

	resp, err := bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.DeleteRetentionPolicy.Enabled, chk.Equals, true)
	c.Assert(*resp.DeleteRetentionPolicy.Days, chk.Equals, int32(5))

	_, err = bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: false}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed t be enough.
	time.Sleep(time.Second * 30)

	resp, err = bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.DeleteRetentionPolicy.Enabled, chk.Equals, false)
	c.Assert(resp.DeleteRetentionPolicy.Days, chk.IsNil)
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyEmpty(c *chk.C) {
	bsu := getBSU()

	days := int32(5)
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true, Days: &days}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed t be enough.
	time.Sleep(time.Second * 30)

	resp, err := bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.DeleteRetentionPolicy.Enabled, chk.Equals, true)
	c.Assert(*resp.DeleteRetentionPolicy.Days, chk.Equals, int32(5))

	// Enabled should default to false and therefore disable the policy
	_, err = bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed t be enough.
	time.Sleep(time.Second * 30)

	resp, err = bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.DeleteRetentionPolicy.Enabled, chk.Equals, false)
	c.Assert(resp.DeleteRetentionPolicy.Days, chk.IsNil)
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyNil(c *chk.C) {
	bsu := getBSU()

	days := int32(5)
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true, Days: &days}})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed t be enough.
	time.Sleep(time.Second * 30)

	resp, err := bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.DeleteRetentionPolicy.Enabled, chk.Equals, true)
	c.Assert(*resp.DeleteRetentionPolicy.Days, chk.Equals, int32(5))

	_, err = bsu.SetProperties(ctx, StorageServiceProperties{})
	c.Assert(err, chk.IsNil)

	// From FE, 30 seconds is guaranteed t be enough.
	time.Sleep(time.Second * 30)

	// If an element of service properties is not passed, the service keeps the current settings.
	resp, err = bsu.GetProperties(ctx)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.DeleteRetentionPolicy.Enabled, chk.Equals, true)
	c.Assert(*resp.DeleteRetentionPolicy.Days, chk.Equals, int32(5))

	// Disable for other tests
	bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: false}})
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyDaysTooSmall(c *chk.C) {
	bsu := getBSU()

	days := int32(0) // Minimum days is 1. Validated on the client.
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true, Days: &days}})
	c.Assert(strings.Contains(err.Error(), validationErrorSubstring), chk.Equals, true)
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyDaysTooLarge(c *chk.C) {
	bsu := getBSU()

	days := int32(366) // Max days is 365. Left to the service for validation.
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true, Days: &days}})
	validateStorageError(c, err, ServiceCodeInvalidXMLDocument)
}

func (s *aztestsSuite) TestAccountDeleteRetentionPolicyDaysOmitted(c *chk.C) {
	bsu := getBSU()

	// Days is required if enabled is true.
	_, err := bsu.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: true}})
	validateStorageError(c, err, ServiceCodeInvalidXMLDocument)
}
