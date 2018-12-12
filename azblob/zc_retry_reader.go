package azblob

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"
)

const CountToEnd = 0

// HTTPGetter is a function type that refers to a method that performs an HTTP GET operation.
type HTTPGetter func(ctx context.Context, i HTTPGetterInfo) (*http.Response, error)

// HTTPGetterInfo is passed to an HTTPGetter function passing it parameters
// that should be used to make an HTTP GET request.
type HTTPGetterInfo struct {
	// Offset specifies the start offset that should be used when
	// creating the HTTP GET request's Range header
	Offset int64

	// Count specifies the count of bytes that should be used to calculate
	// the end offset when creating the HTTP GET request's Range header
	Count int64

	// ETag specifies the resource's etag that should be used when creating
	// the HTTP GET request's If-Match header
	ETag ETag
}

// FailedReadNotifier is a function type that represents the notification function called when a read fails
type FailedReadNotifier func(failureCount int, lastError error, offset int64, count int64, willRetry bool)

// RetryReaderOptions contains properties which can help to decide when to do retry.
type RetryReaderOptions struct {
	// MaxRetryRequests specifies the maximum number of HTTP GET requests that will be made
	// while reading from a RetryReader. A value of zero means that no additional HTTP
	// GET requests will be made.
	MaxRetryRequests   int
	doInjectError      bool
	doInjectErrorRound int

	// Is called, if non-nil, after any failure to read. Expected usage is diagnostic logging.
	NotifyFailedRead  FailedReadNotifier
}

// retryReader implements io.ReaderCloser methods.
// retryReader tries to read from response, and if there is retriable network error
// returned during reading, it will retry according to retry reader option through executing
// user defined action with provided data to get a new response, and continue the overall reading process
// through reading from the new response.
type retryReader struct {
	ctx             context.Context
	response        *http.Response
	info            HTTPGetterInfo
	countWasBounded bool
	o               RetryReaderOptions
	getter          HTTPGetter

	// forced re-read handling (which is the only thread-safe part of this type)
	rrLock      *sync.Mutex
	rrCanceller RequestCanceller
	rrForced    bool               // necessary because forced retries are done with cancellation, and cancellations would not normally be classified as retryable
}

type RequestCanceller interface {
	CancelRequest()
}

// NewRetryReader creates a retry reader.
func NewRetryReader(ctx context.Context, initialResponse *http.Response,
	info HTTPGetterInfo, o RetryReaderOptions, getter HTTPGetter) io.ReadCloser {
	return &retryReader{
		ctx: ctx,
		getter: getter,
		info: info,
		countWasBounded: info.Count != CountToEnd,
		response:    initialResponse,
		rrCanceller: getCancellableRequestBody(initialResponse),
		rrLock:      &sync.Mutex{},
		o:           o}
}


func (s *retryReader) Read(p []byte) (n int, err error) {
	for try := 0; ; try++ {
		//fmt.Println(try)       // Comment out for debugging.
		if s.countWasBounded && s.info.Count == CountToEnd {
			// User specified an original count and the remaining bytes are 0, return 0, EOF
			return 0, io.EOF
		}

		if s.response == nil { // We don't have a response stream to read from, try to get one.
			response, err := s.getter(s.ctx, s.info)
			if err != nil {
				return 0, err
			}
			// Successful GET; this is the network stream we'll read from.
			s.response = response
			s.setCanceller(getCancellableRequestBody(response))
		}
		n, err := s.response.Body.Read(p) // Read from the stream (this will return non-nil err if forceRetry is called, from another goroutine, while it is running)

		// Injection mechanism for testing.
		if s.o.doInjectError && try == s.o.doInjectErrorRound {
			err = &net.DNSError{IsTemporary: true}
		}

		// We successfully read data or end EOF.
		if err == nil || err == io.EOF {
			s.info.Offset += int64(n) // Increments the start offset in case we need to make a new HTTP request in the future
			if s.info.Count != CountToEnd {
				s.info.Count -= int64(n) // Decrement the count in case we need to make a new HTTP request in the future
			}
			return n, err // Return the return to the caller
		}
		s.Close()        // Error, close stream
		s.response = nil // Our stream is no longer good

		// Check the retry count and error code, and decide whether to retry.
		retriesExhausted := try >= s.o.MaxRetryRequests
		_, isNetError := err.(net.Error)
		willRetry := (isNetError || s.wasForcedRetry()) && !retriesExhausted

		// Notify, for logging purposes, of any failures
		if s.o.NotifyFailedRead != nil {
			failureCount := try + 1 // because try is zero-based
			s.o.NotifyFailedRead(failureCount, err, s.info.Offset, s.info.Count, willRetry)
		}

		if willRetry {
			continue
			// Loop around and try to get and read from new stream.
		}
		return n, err // Not retryable, or retries exhausted, so just return
	}
}

func (s *retryReader) Close() error {
	if s.response != nil && s.response.Body != nil {
		return s.response.Body.Close()
	}
	return nil
}

// Returns a function that can be used to force a retry within (and transparently to) an call to Read
func (s *retryReader) getForceRetryFuncOrNil() func(){
	s.rrLock.Lock()
	defer s.rrLock.Unlock()

	if s.rrCanceller == nil {
		return nil            // cancellation is currently impossible, and we'll assume it will remain so
	} else {
		return s.forceRetry   // cancellation is possible at present, and well assume it will remain so (even after any use of s.getter)
	}
}

// Allows forced triggering of the re-read behavior, by cancelling the existing request.
// Can be called mid-read by another go-routine, to cancel and force retryReader to commence a retry cycle.
// Only works with request pipelines where the request offers us a cancellation method
// in the form of a body that implements BodyReadCanceller.
// Why do it that way? Because the alternative would be to generate a new per-request cancellable context
// on each try. That's not impossible.  But... it's complicated by the fact that sometimes we make our
// own HTTP requests (on retries) and sometimes we don't (on the initial request) AND it does have messier handling
// of needing to clean up (cancel) those contexts after use, whereas exactly that functionality is already built
// into our request pipeline using NewRetryPolicyFactory. It's already creating a per-request context, saving the
// CancelFunc, and cleaning it up automatically when the body is closed. So here we are leveraging that existing
// functionality
func (s *retryReader) forceRetry() {
	s.rrLock.Lock()
	defer s.rrLock.Unlock()
	if s.rrCanceller != nil {
		s.rrForced = true
		s.rrCanceller.CancelRequest()
	}
}

func (s *retryReader) setCanceller(canceller RequestCanceller){
	s.rrLock.Lock()
	defer s.rrLock.Unlock()
	s.rrCanceller = canceller
	s.rrForced = false
}

func (s *retryReader) wasForcedRetry() bool {
	s.rrLock.Lock()
	defer s.rrLock.Unlock()
	return s.rrForced
}

func getCancellableRequestBody(r *http.Response) RequestCanceller {
	if c, ok := r.Body.(RequestCanceller); ok {
		return c
	}
	return nil
}


