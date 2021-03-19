package azblob

import (
	chk "gopkg.in/check.v1"
	"net/url"
)

func (s *aztestSuite) TestParseSASQueryParams(c *chk.C) {
	const blobURL = "https://myaccount.blob.core.windows.net/mycontainer/testfile?sp=r&st=2021-03-19T03:48:02Z&se=2021-03-19T11:48:02Z&spr=https&sv=2020-02-10&sr=d&sdd=10&sig=invalidsignature"

	testURL, _ := url.Parse(blobURL)

	bURLParts := NewBlobURLParts(*testURL)
	sas := bURLParts.SAS

	c.Assert(sas.resource, chk.Equals, "d")
	c.Assert(sas.SignedDirectoryDepth(), chk.Equals, "10")
	c.Assert(sas.protocol, chk.Equals, SASProtocolHTTPS)
}