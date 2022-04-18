package azblob

import (
	"encoding/base64"
	chk "gopkg.in/check.v1"
	"math/rand"
	"strings"
)

func (s *aztestsSuite) TestBlockBlobS2SOAuth(c *chk.C) {
	bsu := getBSU()
	ocred, err := getOAuthCredential("", "")
	c.Assert(err, chk.IsNil)

	bsu = bsu.WithPipeline(NewPipeline(ocred, PipelineOptions{})) // get based on the OAuth token

	source := bsu.NewContainerURL(newUUID().String())
	dest := bsu.NewContainerURL(newUUID().String())

	_, err = source.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	_, err = dest.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)

	sourceBlob := source.NewBlockBlobURL("SourceBlob")

	_, err = sourceBlob.Upload(ctx, strings.NewReader("Hello, World!"), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, AccessTierHot, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)

	destBlob := dest.NewBlockBlobURL("DestBlob")

	_, err = destBlob.PutBlobFromURL(ctx, BlobHTTPHeaders{}, sourceBlob.URL(), nil, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, nil, AccessTierHot, nil, ClientProvidedKeyOptions{}, ocred)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlockBlobS2SOAuthByBlock(c *chk.C) {
	bsu := getBSU()
	ocred, err := getOAuthCredential("", "")
	c.Assert(err, chk.IsNil)

	bsu = bsu.WithPipeline(NewPipeline(ocred, PipelineOptions{})) // get based on the OAuth token

	source := bsu.NewContainerURL(newUUID().String())
	dest := bsu.NewContainerURL(newUUID().String())

	_, err = source.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	_, err = dest.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)

	sourceBlob := source.NewBlockBlobURL("SourceBlob")

	_, err = sourceBlob.Upload(ctx, strings.NewReader("Hello, World!"), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, AccessTierHot, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)

	destBlob := dest.NewBlockBlobURL("DestBlob")

	_, err = destBlob.StageBlockFromURL(ctx, base64.StdEncoding.EncodeToString([]byte(newUUID().String())), sourceBlob.URL(), 0, int64(len("Hello, World!")), LeaseAccessConditions{}, ModifiedAccessConditions{}, ClientProvidedKeyOptions{}, ocred)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlockBlobS2SOAuthCopyFromURL(c *chk.C) {
	bsu := getBSU()
	ocred, err := getOAuthCredential("", "")
	c.Assert(err, chk.IsNil)

	bsu = bsu.WithPipeline(NewPipeline(ocred, PipelineOptions{})) // get based on the OAuth token

	source := bsu.NewContainerURL(newUUID().String())
	dest := bsu.NewContainerURL(newUUID().String())

	_, err = source.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	_, err = dest.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)

	sourceBlob := source.NewBlockBlobURL("SourceBlob")

	_, err = sourceBlob.Upload(ctx, strings.NewReader("Hello, World!"), BlobHTTPHeaders{}, nil, BlobAccessConditions{}, AccessTierHot, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)

	destBlob := dest.NewBlockBlobURL("DestBlob")

	_, err = destBlob.CopyFromURL(ctx, sourceBlob.URL(), nil, ModifiedAccessConditions{}, BlobAccessConditions{}, nil, AccessTierHot, nil, ImmutabilityPolicyOptions{}, ocred)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestPageBlobS2SOAuth(c *chk.C) {
	bsu := getBSU()
	ocred, err := getOAuthCredential("", "")
	c.Assert(err, chk.IsNil)

	bsu = bsu.WithPipeline(NewPipeline(ocred, PipelineOptions{})) // get based on the OAuth token

	source := bsu.NewContainerURL(newUUID().String())
	dest := bsu.NewContainerURL(newUUID().String())

	_, err = source.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	_, err = dest.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)

	sourceBlob := source.NewPageBlobURL("SourceBlob")

	_, err = sourceBlob.Create(ctx, 512, 0, BlobHTTPHeaders{}, nil, BlobAccessConditions{}, PremiumPageBlobAccessTierNone, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)

	page := make([]byte, 512)
	for k := range page {
		page[k] = byte(rand.Intn(256))
	}

	// bytes.NewBuffer does not work, because bytes.Buffer does not satisfy Seeker.
	_, err = sourceBlob.UploadPages(ctx, 0, strings.NewReader(string(page)), PageBlobAccessConditions{}, nil, ClientProvidedKeyOptions{})

	destBlob := dest.NewPageBlobURL("DestBlob")

	_, err = destBlob.Create(ctx, 512, 0, BlobHTTPHeaders{}, nil, BlobAccessConditions{}, PremiumPageBlobAccessTierNone, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)
	_, err = destBlob.UploadPagesFromURL(ctx, sourceBlob.URL(), 0, 0, 512, nil, PageBlobAccessConditions{}, ModifiedAccessConditions{}, ClientProvidedKeyOptions{}, ocred)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestAppendBlobS2SOAuth(c *chk.C) {
	bsu := getBSU()
	ocred, err := getOAuthCredential("", "")
	c.Assert(err, chk.IsNil)

	bsu = bsu.WithPipeline(NewPipeline(ocred, PipelineOptions{})) // get based on the OAuth token

	source := bsu.NewContainerURL(newUUID().String())
	dest := bsu.NewContainerURL(newUUID().String())

	_, err = source.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	_, err = dest.Create(ctx, nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)

	sourceBlob := source.NewAppendBlobURL("SourceBlob")

	_, err = sourceBlob.Create(ctx, BlobHTTPHeaders{}, nil, BlobAccessConditions{}, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)
	_, err = sourceBlob.AppendBlock(ctx, strings.NewReader("Hello, World!"), AppendBlobAccessConditions{}, nil, ClientProvidedKeyOptions{})
	c.Assert(err, chk.IsNil)

	destBlob := dest.NewAppendBlobURL("DestBlob")

	_, err = destBlob.Create(ctx, BlobHTTPHeaders{}, nil, BlobAccessConditions{}, nil, ClientProvidedKeyOptions{}, ImmutabilityPolicyOptions{})
	c.Assert(err, chk.IsNil)
	_, err = destBlob.AppendBlockFromURL(ctx, sourceBlob.URL(), 0, int64(len("Hello, World!")), AppendBlobAccessConditions{}, ModifiedAccessConditions{}, nil, ClientProvidedKeyOptions{}, ocred)
	c.Assert(err, chk.IsNil)
}
