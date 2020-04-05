package azblob

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

type blockWriter interface {
	StageBlock(context.Context, string, io.ReadSeeker, LeaseAccessConditions, []byte) (*BlockBlobStageBlockResponse, error)
	CommitBlockList(context.Context, []string, BlobHTTPHeaders, Metadata, BlobAccessConditions) (*BlockBlobCommitBlockListResponse, error)
}

// copyFromReader copies a source io.Reader to blob storage using concurrent uploads.
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

	prefix uuid
	to     blockWriter
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

// getErr returns an error by priority. First, if a function set an error, it returns that error. Next, if the Context has an error it returns
// that error. Otherwise it is nil. getErr supports only a single call.
func (c *copier) getErr() error {
	select {
	case err := <-c.errCh:
		return err
	default:
	}
	return c.ctx.Err()
}

// sendChunk reads data from out internal reader, creates a chunk, and sends it to be written via a channel.
// sendChunk returns io.EOF when the reader returns an io.EOF.
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

// writer writes chunks that come on a channel via write(). This is mean to be started by a goroutine.
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

	select {
	case err := <-c.errCh:
		return err
	default:
	}

	if c.ctx.Err() != nil {
		return c.ctx.Err()
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
