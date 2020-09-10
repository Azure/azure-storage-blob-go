package azblob

import (
	chk "gopkg.in/check.v1"
)

func (s *aztestsSuite) TestSetBlobTags(c *chk.C) {
	//bsu := getBSU()
	//containerURL, _ := createNewContainer(c, bsu)
	//defer deleteContainer(c, containerURL)
	//blobURL, _ := getBlockBlobURL(c, containerURL)
	//blobTags := map[string]string {
	//	"azure":"blob",
	//	"blob":"sdk",
	//	"sdk":"go",
	//}
	//blockBlobUploadResp, err := blobURL.Upload(ctx, bytes.NewReader([]byte("data")), BlobHTTPHeaders{}, basicMetadata, BlobAccessConditions{}, DefaultAccessTier, nil)
	//c.Assert(err, chk.IsNil)
	//c.Assert(blockBlobUploadResp.ETag(), chk.Equals, blobTags)
	//blobURL.SetTags(ctx, nil, nil,nil, nil, nil, nil, blobTags)

}
