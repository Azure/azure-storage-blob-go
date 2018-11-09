package azblob_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/Azure/azure-storage-blob-go/azblob"
	chk "gopkg.in/check.v1"
)

// Testings for RetryReader
// This reader return one byte through each Read call
type perByteReader struct {
	RandomBytes []byte // Random generated bytes

	byteCount              int // Bytes can be returned before EOF
	currentByteIndex       int // Bytes that have already been returned.
	doInjectError          bool
	doInjectErrorByteIndex int
	doInjectTimes          int
	injectedError          error
}

func newPerByteReader(byteCount int) *perByteReader {
	perByteReader := perByteReader{
		byteCount: byteCount,
	}

	perByteReader.RandomBytes = make([]byte, byteCount)
	rand.Read(perByteReader.RandomBytes)

	return &perByteReader
}

func (r *perByteReader) Read(b []byte) (n int, err error) {
	if r.doInjectError && r.doInjectErrorByteIndex == r.currentByteIndex && r.doInjectTimes > 0 {
		r.doInjectTimes--
		return 0, r.injectedError
	}

	if r.currentByteIndex < r.byteCount {
		n = copy(b, r.RandomBytes[r.currentByteIndex:r.currentByteIndex+1])
		r.currentByteIndex += n
		return
	}

	return 0, io.EOF
}

func (r *perByteReader) Close() error {
	return nil
}

// Test normal retry succeed, note initial response not provided.
// Tests both with and without notification of failures
func (r *aztestsSuite) TestRetryReaderReadWithRetry(c *chk.C) {
	// Test twice, the second time using the optional "logging"/notification callback for failed tries
	// We must test both with and without the callback, since be testing without
	// we are testing that it is, indeed, optional to provide the callback
	for _, logThisRun := range []bool { false, true} {

		// Extra setup for testing notification of failures (i.e. of unsuccessful tries)
		failureMethodNumCalls := 0
		failureWillRetryCount := 0
		failureLastReportedFailureCount := -1
		var failureLastReportedError error = nil
		failureMethod := func(failureCount int, lastError error, willRetry bool) {
			failureMethodNumCalls++;
			if willRetry {
				failureWillRetryCount++
			}
			failureLastReportedFailureCount = failureCount;
			failureLastReportedError = lastError;
		}

		// Main test setup
		byteCount := 1
		body := newPerByteReader(byteCount)
		body.doInjectError = true
		body.doInjectErrorByteIndex = 0
		body.doInjectTimes = 1
		body.injectedError = &net.DNSError{IsTemporary: true}

		getter := func(ctx context.Context, info azblob.HTTPGetterInfo) (*http.Response, error) {
			r := http.Response{}
			body.currentByteIndex = int(info.Offset)
			r.Body = body

			return &r, nil
		}

		httpGetterInfo := azblob.HTTPGetterInfo{Offset: 0, Count: int64(byteCount)}
		initResponse, err := getter(context.Background(), httpGetterInfo)
		c.Assert(err, chk.IsNil)

		rrOptions := azblob.RetryReaderOptions{MaxRetryRequests: 1}
		if logThisRun {
			rrOptions.NotifyFailedRead = failureMethod;
		}
		retryReader := azblob.NewRetryReader(context.Background(), initResponse, httpGetterInfo, rrOptions, getter)

		// should fail and succeed through retry
		can := make([]byte, 1)
		n, err := retryReader.Read(can)
		c.Assert(n, chk.Equals, 1)
		c.Assert(err, chk.IsNil)

		// check "logging", if it was enabled
		if logThisRun {
			// We only expect one failed try in this test
			// And the notification method is not called for successes
			c.Assert(failureMethodNumCalls, chk.Equals, 1)             // this is the number of calls we counted
			c.Assert(failureWillRetryCount, chk.Equals, 1)             // the sole failure was retried
			c.Assert(failureLastReportedFailureCount, chk.Equals, 1)   // this is the number of failures reported by the notification method
			c.Assert(failureLastReportedError, chk.NotNil)
		}
		// should return EOF
		n, err = retryReader.Read(can)
		c.Assert(n, chk.Equals, 0)
		c.Assert(err, chk.Equals, io.EOF)
	}
}

// Test normal retry fail as retry Count not enough.
func (r *aztestsSuite) TestRetryReaderReadNegativeNormalFail(c *chk.C) {
	// Extra setup for testing notification of failures (i.e. of unsuccessful tries)
	failureMethodNumCalls := 0
	failureWillRetryCount := 0
	failureLastReportedFailureCount := -1
	var failureLastReportedError error = nil
	failureMethod := func(failureCount int, lastError error, willRetry bool) {
		failureMethodNumCalls++;
		if willRetry {
			failureWillRetryCount++
		}
		failureLastReportedFailureCount = failureCount;
		failureLastReportedError = lastError;
	}

	// Main test setup
	byteCount := 1
	body := newPerByteReader(byteCount)
	body.doInjectError = true
	body.doInjectErrorByteIndex = 0
	body.doInjectTimes = 2
	body.injectedError = &net.DNSError{IsTemporary: true}

	startResponse := http.Response{}
	startResponse.Body = body

	getter := func(ctx context.Context, info azblob.HTTPGetterInfo) (*http.Response, error) {
		r := http.Response{}
		body.currentByteIndex = int(info.Offset)
		r.Body = body

		return &r, nil
	}

	rrOptions := azblob.RetryReaderOptions{
		MaxRetryRequests: 1,
		NotifyFailedRead: failureMethod}
	retryReader := azblob.NewRetryReader(context.Background(), &startResponse, azblob.HTTPGetterInfo{Offset: 0, Count: int64(byteCount)}, rrOptions, getter)

	// should fail
	can := make([]byte, 1)
	n, err := retryReader.Read(can)
	c.Assert(n, chk.Equals, 0)
	c.Assert(err, chk.Equals, body.injectedError)

	// Check that we recieved the right notification callbacks
	// We only expect two failed tries in this test, but only one
	// of the would have had willRetry = true
	c.Assert(failureMethodNumCalls, chk.Equals, 2)             // this is the number of calls we counted
	c.Assert(failureWillRetryCount, chk.Equals, 1)             // only the first failure was retried
	c.Assert(failureLastReportedFailureCount, chk.Equals, 2)   // this is the number of failures reported by the notification method
	c.Assert(failureLastReportedError, chk.NotNil)
}

// Test boundary case when Count equals to 0 and fail.
func (r *aztestsSuite) TestRetryReaderReadCount0(c *chk.C) {
	byteCount := 1
	body := newPerByteReader(byteCount)
	body.doInjectError = true
	body.doInjectErrorByteIndex = 1
	body.doInjectTimes = 1
	body.injectedError = &net.DNSError{IsTemporary: true}

	startResponse := http.Response{}
	startResponse.Body = body

	getter := func(ctx context.Context, info azblob.HTTPGetterInfo) (*http.Response, error) {
		r := http.Response{}
		body.currentByteIndex = int(info.Offset)
		r.Body = body

		return &r, nil
	}

	retryReader := azblob.NewRetryReader(context.Background(), &startResponse, azblob.HTTPGetterInfo{Offset: 0, Count: int64(byteCount)}, azblob.RetryReaderOptions{MaxRetryRequests: 1}, getter)

	// should consume the only byte
	can := make([]byte, 1)
	n, err := retryReader.Read(can)
	c.Assert(n, chk.Equals, 1)
	c.Assert(err, chk.IsNil)

	// should not read when Count=0, and should return EOF
	n, err = retryReader.Read(can)
	c.Assert(n, chk.Equals, 0)
	c.Assert(err, chk.Equals, io.EOF)
}

func (r *aztestsSuite) TestRetryReaderReadNegativeNonRetriableError(c *chk.C) {
	byteCount := 1
	body := newPerByteReader(byteCount)
	body.doInjectError = true
	body.doInjectErrorByteIndex = 0
	body.doInjectTimes = 1
	body.injectedError = fmt.Errorf("not retriable error")

	startResponse := http.Response{}
	startResponse.Body = body

	getter := func(ctx context.Context, info azblob.HTTPGetterInfo) (*http.Response, error) {
		r := http.Response{}
		body.currentByteIndex = int(info.Offset)
		r.Body = body

		return &r, nil
	}

	retryReader := azblob.NewRetryReader(context.Background(), &startResponse, azblob.HTTPGetterInfo{Offset: 0, Count: int64(byteCount)}, azblob.RetryReaderOptions{MaxRetryRequests: 2}, getter)

	dest := make([]byte, 1)
	_, err := retryReader.Read(dest)
	c.Assert(err, chk.Equals, body.injectedError)
}

// End testings for RetryReader
