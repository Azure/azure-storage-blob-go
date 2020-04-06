package azblob

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

// blockWriter provides methods to upload blocks that represent a file to a server and commit them.
// This allows us to provide a local implementation that fakes the server for hermetic testing.
type blockWriter interface {
	StageBlock(context.Context, string, io.ReadSeeker, LeaseAccessConditions, []byte) (*BlockBlobStageBlockResponse, error)
	CommitBlockList(context.Context, []string, BlobHTTPHeaders, Metadata, BlobAccessConditions) (*BlockBlobCommitBlockListResponse, error)
}

// copyFromReader copies a source io.Reader to blob storage using concurrent uploads.
// TODO(someone): The existing model provides a buffer size and buffer limit as limiting factors.  The buffer size is probably
// useless other than needing to be above some number, as the network stack is going to hack up the buffer over some size. The
// max buffers is providing a cap on how much memory we use (by multiplying it times the buffer size) and how many go routines can upload
// at a time.  I think having a single max memory dial would be more efficient.  We can choose an internal buffer size that works
// well, 4 MiB or 8 MiB, and autoscale to as many goroutines within the memory limit. This gives a single dial to tweak and we can
// choose a max value for the memory setting based on internal transfers within Azure (which will give us the maximum throughput model).
// We can even provide a utility to dial this number in for customer networks to optimize their copies.
func copyFromReader(ctx context.Context, from io.Reader, to blockWriter, o UploadStreamToBlockBlobOptions) (*BlockBlobCommitBlockListResponse, error) {
	o.defaults()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	copy := &copier{
		ctx:    ctx,
		cancel: cancel,
		reader: from,
		to:     to,
		prefix: newUUID(),
		o:      o,
		ch:     make(chan *chunk, 1),
		errCh:  make(chan error, 1),
		buffers: sync.Pool{
			New: func() interface{} {
				return &chunk{payload: make([]byte, o.BufferSize)}
			},
		},
	}

	// Starts the pools of concurrent writers.
	copy.wg.Add(o.MaxBuffers)
	for i := 0; i < o.MaxBuffers; i++ {
		go copy.writer()
	}

	// Send all our chunks until we get an error.
	var err error
	for {
		if err = copy.sendChunk(); err != nil {
			break
		}
	}
	// If the error is not EOF, then we have a problem.
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	// Close out our upload.
	if err := copy.close(); err != nil {
		return nil, err
	}

	return copy.result, nil
}

// chunk represents a chunk of data to be sent to blob storage.
type chunk struct {
	// num is the number of the chunk. Must be unique per chunk upload for a single file.
	num int32
	// payload is the chunk data.
	payload []byte
}

// copier streams a file via chunks in parallel from a reader representing a file.
// Do not use directly, instead use copyFromReader().
type copier struct {
	// ctx holds the context of a copier. This is normally a faux paux to store a Context in a struct. In this case,
	// the copier has the lifetime of a function call, so its fine.
	ctx    context.Context
	cancel context.CancelFunc

	// reader is the source to be written to disk.
	reader io.Reader
	// to is the location we are writing our chunks to.
	to blockWriter

	prefix uuid
	o      UploadStreamToBlockBlobOptions

	// num is the current chunk we are on.
	num int32
	// ch is used to pass the next chunk of data from our reader to one of the writers.
	ch chan *chunk
	// errCh is used to hold the first error from our concurrent writers.
	errCh chan error
	// wg provides a count of how many writers we are waiting to finish.
	wg sync.WaitGroup
	// buffers provides a pool of chunks that can be reused.
	buffers sync.Pool

	// result holds the final result from blob storage after we have submitted all chunks.
	result *BlockBlobCommitBlockListResponse
}

// getErr returns an error by priority. First, if a function set an error, it returns that error. Next, if the Context has an error
// it returns that error. Otherwise it is nil. getErr supports only returning an error once per copier.
func (c *copier) getErr() error {
	select {
	case err := <-c.errCh:
		return err
	default:
	}
	return c.ctx.Err()
}

// sendChunk reads data from out internal reader, creates a chunk, and sends it to be written via a channel.
// sendChunk returns io.EOF when the reader returns an io.EOF or io.ErrUnexpectedEOF.
func (c *copier) sendChunk() error {
	if err := c.getErr(); err != nil {
		return err
	}

	chunk := c.buffers.Get().(*chunk)
	n, err := io.ReadFull(c.reader, chunk.payload)
	chunk.payload = chunk.payload[0:n]
	switch {
	case err == nil && n == 0:
		return nil
	case err == nil:
		chunk.num = c.num
		c.num++
		c.ch <- chunk
		return nil
	case err != nil && (err == io.EOF || err == io.ErrUnexpectedEOF) && n == 0:
		return io.EOF
	}

	if err == io.EOF || err == io.ErrUnexpectedEOF {
		chunk.num = c.num
		c.num++
		c.ch <- chunk
		return io.EOF
	}
	if err := c.getErr(); err != nil {
		return err
	}
	return err
}

// writer writes chunks sent on a channel.
func (c *copier) writer() {
	defer c.wg.Done()

	for chunk := range c.ch {
		if err := c.write(chunk); err != nil {
			if !errors.Is(err, context.Canceled) {
				select {
				case c.errCh <- err:
					c.cancel()
				default:
				}
				return
			}
		}
	}
}

// write uploads a chunk to blob storage.
// TODO(someone): Might be worth having StageBlock() retry with some exponential delays before giving up.
// Sucks to have a 100GiB upload die because of a single write that could be corrected. Right now this
// is just mimicking the previous behavior.
func (c *copier) write(chunk *chunk) error {
	defer c.buffers.Put(chunk)

	if err := c.ctx.Err(); err != nil {
		return err
	}

	blockID := newUuidBlockID(c.prefix).WithBlockNumber(uint32(chunk.num)).ToBase64()
	_, err := c.to.StageBlock(c.ctx, blockID, bytes.NewReader(chunk.payload), LeaseAccessConditions{}, nil)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

// close commits our blocks to blob storage and closes our writer.
func (c *copier) close() error {
	close(c.ch)
	c.wg.Wait()

	if err := c.getErr(); err != nil {
		return err
	}

	blockID := newUuidBlockID(c.prefix)
	var blockIDs []string
	if c.num > 0 {
		blockIDs = make([]string, 0, c.num-1)
		for i := 0; i < int(c.num); i++ {
			blockIDs = append(blockIDs, blockID.WithBlockNumber(uint32(i)).ToBase64())
		}
	}

	var err error
	c.result, err = c.to.CommitBlockList(c.ctx, blockIDs, c.o.BlobHTTPHeaders, c.o.Metadata, c.o.AccessConditions)
	return err
}
