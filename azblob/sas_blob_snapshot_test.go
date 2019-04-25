package azblob_test

import(
	"github.com/Azure/azure-storage-blob-go/azblob"
	chk "gopkg.in/check.v1"
	"os"
	"time"
)

func (s *aztestsSuite) TestNewSASQueryParametersSnapshot(c *chk.C) {
	bsu := getBSU()
	containerURL, containerName := getContainerURL(c, bsu)
	blobURL, blobName := getBlockBlobURL(c, containerURL)

	currentTime := time.Now().UTC()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		StartTime: currentTime,
		ExpiryTime: currentTime.Add(48 * time.Hour),
		SnapshotTime: currentTime,
		Permissions: "racwd",
		ContainerName: containerName,
		BlobName: blobName,
		Protocol: azblob.SASProtocolHTTPS,
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	parts := azblob.NewBlobURLParts(blobURL.URL())
	parts.SAS = sasQueryParams
	//parts.Snapshot = sasQueryParams.SnapshotTime().Format(azblob.SASTimeFormat) No need, if it's present in the sasQueryParams, it'll be appended.
	testURL := parts.URL()

	correctURL := "https://" + os.Getenv("ACCOUNT_NAME") + ".blob.core.windows.net/" + containerName + "/" + blobName +
		"?snapshot=" + sasQueryParams.SnapshotTime().Format(azblob.SASTimeFormat) + "&" + sasQueryParams.Encode()
	c.Assert(testURL.String(), chk.Equals, correctURL)
}