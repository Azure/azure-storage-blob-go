package azblob

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// UploadStreamToBlockBlobOptions identifies options used by the UploadStreamToBlockBlob function. Note that the
// BlockSize field is mandatory and must be set; other fields are optional.
type UploadStreamToBlockBlobOptions struct {
	// BlockSize is mandatory. It specifies the block size to use; the maximum size is BlockBlobMaxPutBlockBytes.
	BlockSize int64

	// Progress is a function that is invoked periodically as bytes are send in a PutBlock call to the BlockBlobURL.
	Progress pipeline.ProgressReceiver

	// BlobHTTPHeaders indicates the HTTP headers to be associated with the blob when PutBlockList is called.
	BlobHTTPHeaders BlobHTTPHeaders

	// Metadata indicates the metadata to be associated with the blob when PutBlockList is called.
	Metadata Metadata

	// AccessConditions indicates the access conditions for the block blob.
	AccessConditions BlobAccessConditions
}

// UploadStreamToBlockBlob uploads a stream of data in blocks to a block blob.
func UploadStreamToBlockBlob(ctx context.Context, stream io.ReaderAt, streamSize int64,
	blockBlobURL BlockBlobURL, o UploadStreamToBlockBlobOptions) (*BlockBlobsPutBlockListResponse, error) {

	if o.BlockSize <= 0 || o.BlockSize > BlockBlobMaxPutBlockBytes {
		panic(fmt.Sprintf("BlockSize option must be > 0 and <= %d", BlockBlobMaxPutBlockBytes))
	}

	numBlocks := ((streamSize - int64(1)) / o.BlockSize) + 1
	if numBlocks > BlockBlobMaxBlocks {
		panic(fmt.Sprintf("The streamSize is too big or the BlockSize is too small; the number of blocks must be <= %d", BlockBlobMaxBlocks))
	}
	blockIDList := make([]string, numBlocks) // Base 64 encoded block IDs
	blockSize := o.BlockSize

	for blockNum := int64(0); blockNum < numBlocks; blockNum++ {
		if blockNum == numBlocks-1 { // Last block
			blockSize = streamSize - (blockNum * o.BlockSize) // Remove size of all uploaded blocks from total
		}

		streamOffset := blockNum * o.BlockSize
		// Prepare to read the proper block/section of the file
		var body io.ReadSeeker = io.NewSectionReader(stream, streamOffset, blockSize)
		if o.Progress != nil {
			body = pipeline.NewRequestBodyProgress(body,
				func(bytesTransferred int64) { o.Progress(streamOffset + bytesTransferred) })
		}

		// Block IDs are unique values to avoid issue if 2+ clients are uploading blocks
		// at the same time causeing PutBlockList to get a mix of blocks from all the clients.
		blockIDList[blockNum] = base64.StdEncoding.EncodeToString(newUUID().bytes())
		_, err := blockBlobURL.PutBlock(ctx, blockIDList[blockNum], body, o.AccessConditions.LeaseAccessConditions)
		if err != nil {
			return nil, err
		}
	}
	return blockBlobURL.PutBlockList(ctx, blockIDList, o.Metadata, o.BlobHTTPHeaders, o.AccessConditions)
}

// DownloadStreamOptions is used to configure a call to NewDownloadBlobToStream to download a large stream with intelligent retries.
type DownloadStreamOptions struct {
	// Range indicates the starting offset and count of bytes within the blob to download.
	Range BlobRange

	// AccessConditions indicates the BlobAccessConditions to use when accessing the blob.
	AccessConditions BlobAccessConditions
}

type retryStream struct {
	ctx      context.Context
	getBlob  func(ctx context.Context, blobRange BlobRange, ac BlobAccessConditions, rangeGetContentMD5 bool) (*GetResponse, error)
	o        DownloadStreamOptions
	response *http.Response
}

// NewDownloadStream creates a stream over a blob allowing you download the blob's contents.
// When network errors occur, the retry stream internally issues new HTTP GET requests for
// the remaining range of the blob's contents. The GetBlob argument identifies the function
// to invoke when the GetRetryStream needs to make an HTTP GET request as Read methods are called.
// The callback can wrap the response body (with progress reporting, for example) before returning.
func NewDownloadStream(ctx context.Context,
	getBlob func(ctx context.Context, blobRange BlobRange, ac BlobAccessConditions, rangeGetContentMD5 bool) (*GetResponse, error),
	o DownloadStreamOptions) io.ReadCloser {

	// BlobAccessConditions may already have an If-Match:etag header
	if getBlob == nil {
		panic("getBlob must not be nil")
	}
	return &retryStream{ctx: ctx, getBlob: getBlob, o: o, response: nil}
}

func (s *retryStream) Read(p []byte) (n int, err error) {
	for {
		if s.response != nil { // We working with a successful response
			n, err := s.response.Body.Read(p) // Read from the stream
			if err == nil || err == io.EOF {  // We successfully read data or end EOF
				s.o.Range.Offset += int64(n) // Increments the start offset in case we need to make a new HTTP request in the future
				if s.o.Range.Count != 0 {
					s.o.Range.Count -= int64(n) // Decrement the count in case we need to make a new HTTP request in the future
				}
				return n, err // Return the return to the caller
			}
			s.Close()
			s.response = nil // Something went wrong; our stream is no longer good
			if nerr, ok := err.(net.Error); ok {
				if !nerr.Timeout() && !nerr.Temporary() {
					return n, err // Not retryable
				}
			} else {
				return n, err // Not retryable, just return
			}
		}

		// We don't have a response stream to read from, try to get one
		response, err := s.getBlob(s.ctx, s.o.Range, s.o.AccessConditions, false)
		if err != nil {
			return 0, err
		}
		// Successful GET; this is the network stream we'll read from
		s.response = response.Response()

		// Ensure that future requests are from the same version of the source
		s.o.AccessConditions.IfMatch = response.ETag()

		// Loop around and try to read from this stream
	}
}

func (s *retryStream) Close() error {
	if s.response != nil && s.response.Body != nil {
		return s.response.Body.Close()
	}
	return nil
}
