package azblob

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// RetryPolicy tells the pipeline what kind of retry policy to use. See the RetryPolicy* constants.
type RetryPolicy int32

const (
	// RetryPolicyExponential tells the pipeline to use an exponential back-off retry policy
	RetryPolicyExponential RetryPolicy = 0

	// RetryPolicyFixed tells the pipeline to use a fixed back-off retry policy
	RetryPolicyFixed RetryPolicy = 1
)

// RetryOptions configures the retry policy's behavior.
type RetryOptions struct {
	// Policy tells the pipeline what kind of retry policy to use. See the RetryPolicy* constants.\
	// A value of zero means that you accept our default policy.
	Policy RetryPolicy

	// MaxTries specifies the maximum number of attempts an operation will be tried before producing an error (0=default).
	// A value of zero means that you accept our default policy. A value of 1 means 1 try and no retries.
	MaxTries int32

	// TryTimeout indicates the maximum time allowed for any single try of an HTTP request.
	// A value of zero means that you accept our default timeout. NOTE: When transferring large amounts
	// of data, the default TryTimeout will probably not be sufficient. You should override this value
	// based on the bandwidth available to the host machine and proximity to the Storage service. A good
	// starting point may be something like (60 seconds per MB of anticipated-payload-size).
	TryTimeout time.Duration

	// RetryDelay specifies the amount of delay to use before retrying an operation (0=default).
	// The delay increases (exponentially or linearly) with each retry up to a maximum specified by
	// MaxRetryDelay. If you specify 0, then you must also specify 0 for MaxRetryDelay.
	RetryDelay time.Duration

	// MaxRetryDelay specifies the maximum delay allowed before retrying an operation (0=default).
	// If you specify 0, then you must also specify 0 for RetryDelay.
	MaxRetryDelay time.Duration

	// RetryReadsFromSecondaryHost specifies whether the retry policy should retry a read operation against another host.
	// If RetryReadsFromSecondaryHost is "" (the default) then operations are not retried against another host.
	// NOTE: Before setting this field, make sure you understand the issues around reading stale & potentially-inconsistent
	// data at this webpage: https://docs.microsoft.com/en-us/azure/storage/common/storage-designing-ha-apps-with-ragrs
	RetryReadsFromSecondaryHost string
}

func (o RetryOptions) defaults() RetryOptions {
	if o.Policy != RetryPolicyExponential && o.Policy != RetryPolicyFixed {
		panic(errors.New("RetryPolicy must be RetryPolicyExponential or RetryPolicyFixed"))
	}
	if o.MaxTries < 0 {
		panic("MaxTries must be >= 0")
	}
	if o.TryTimeout < 0 || o.RetryDelay < 0 || o.MaxRetryDelay < 0 {
		panic("TryTimeout, RetryDelay, and MaxRetryDelay must all be >= 0")
	}
	if o.RetryDelay > o.MaxRetryDelay {
		panic("RetryDelay must be <= MaxRetryDelay")
	}
	if (o.RetryDelay == 0 && o.MaxRetryDelay != 0) || (o.RetryDelay != 0 && o.MaxRetryDelay == 0) {
		panic(errors.New("Both RetryDelay and MaxRetryDelay must be 0 or neither can be 0"))
	}

	IfDefault := func(current *time.Duration, desired time.Duration) {
		if *current == time.Duration(0) {
			*current = desired
		}
	}

	// Set defaults if unspecified
	if o.MaxTries == 0 {
		o.MaxTries = 4
	}
	switch o.Policy {
	case RetryPolicyExponential:
		IfDefault(&o.TryTimeout, 30*time.Second)
		IfDefault(&o.RetryDelay, 4*time.Second)
		IfDefault(&o.MaxRetryDelay, 120*time.Second)

	case RetryPolicyFixed:
		IfDefault(&o.TryTimeout, 30*time.Second)
		IfDefault(&o.RetryDelay, 30*time.Second)
		IfDefault(&o.MaxRetryDelay, 120*time.Second)
	}
	return o
}

func (o RetryOptions) calcDelay(try int32) time.Duration { // try is >=1; never 0
	pow := func(number int64, exponent int32) int64 { // pow is nested helper function
		var result int64 = 1
		for n := int32(0); n < exponent; n++ {
			result *= number
		}
		return result
	}

	delay := time.Duration(0)
	switch o.Policy {
	case RetryPolicyExponential:
		delay = time.Duration(pow(2, try-1)-1) * o.RetryDelay

	case RetryPolicyFixed:
		if try > 1 { // Any try after the 1st uses the fixed delay
			delay = o.RetryDelay
		}
	}

	// Introduce some jitter:  [0.0, 1.0) / 2 = [0.0, 0.5) + 0.8 = [0.8, 1.3)
	delay *= time.Duration(rand.Float32()/2 + 0.8) // NOTE: We want math/rand; not crypto/rand
	if delay > o.MaxRetryDelay {
		delay = o.MaxRetryDelay
	}
	return delay
}

// NewRetryPolicyFactory creates a RetryPolicyFactory object configured using the specified options.
func NewRetryPolicyFactory(o RetryOptions) pipeline.Factory {
	return &retryPolicyFactory{o: o.defaults()}
}

type retryPolicyFactory struct {
	o RetryOptions
}

func (f *retryPolicyFactory) New(next pipeline.Policy, po *pipeline.PolicyOptions) pipeline.Policy {
	return &retryPolicy{o: f.o, next: next}
}

type retryPolicy struct {
	next pipeline.Policy
	o    RetryOptions
}

// According to https://github.com/golang/go/wiki/CompilerOptimizations, the compiler will inline this method and hopefully optimize all calls to it away
var logf = func(format string, a ...interface{}) {}

// Use this version to see the retry method's code path (import "fmt")
//var logf = fmt.Printf

func (p *retryPolicy) Do(ctx context.Context, request pipeline.Request) (response pipeline.Response, err error) {
	// Before each try, we'll select either the primary or secondary URL.
	secondaryHost := ""
	primaryTry := int32(0) // This indicates how many tries we've attempted against the primary DC

	// We only consider retring against a secondary if we have a read request (GET/HEAD) AND this policy has a Secondary URL it can use
	considerSecondary := (request.Method == http.MethodGet || request.Method == http.MethodHead) && p.o.RetryReadsFromSecondaryHost != ""
	if considerSecondary {
		secondaryHost = p.o.RetryReadsFromSecondaryHost
	}

	// Exponential retry algorithm: ((2 ^ attempt) - 1) * delay * random(0.8, 1.2)
	// When to retry: connection failure or an HTTP status code of 500 or greater, except 501 and 505
	// If using a secondary:
	//    Even tries go against primary; odd tries go against the secondary
	//    For a primary wait ((2 ^ primaryTries - 1) * delay * random(0.8, 1.2)
	//    If secondary gets a 404, don't fail, retry but future retries are only against the primary
	//    When retrying against a secondary, ignore the retry count and wait (.1 second * random(0.8, 1.2))
	for try := int32(1); try <= p.o.MaxTries; try++ {
		logf("\n=====> Try=%d\n", try)

		// Determine which endpoint to try. It's primary if there is no secondary or if it is an add # attempt.
		tryingPrimary := !considerSecondary || (try%2 == 1)
		// Select the correct host and delay
		if tryingPrimary {
			primaryTry++
			delay := p.o.calcDelay(primaryTry)
			logf("Primary try=%d, Delay=%v\n", primaryTry, delay)
			time.Sleep(delay) // The 1st try returns 0 delay
		} else {
			delay := time.Second * time.Duration(rand.Float32()/2+0.8)
			logf("Secondary try=%d, Delay=%v\n", try-primaryTry, delay)
			time.Sleep(delay) // Delay with some jitter before trying secondary
		}

		// Clone the original request to ensure that each try starts with the original (unmutated) request.
		requestCopy := request.Copy()
		if try > 1 {
			// For a retry, seek to the beginning of the Body stream.
			if err = requestCopy.RewindBody(); err != nil {
				panic(err)
			}
		}
		if !tryingPrimary {
			requestCopy.Request.URL.Host = secondaryHost
		}

		// Set the server-side timeout query parameter "timeout=[seconds]"
		timeout := int32(p.o.TryTimeout.Seconds()) // Max seconds per try
		if deadline, ok := ctx.Deadline(); ok {    // If user's ctx has a deadline, make the timeout the smaller of the two
			t := int32(deadline.Sub(time.Now()).Seconds()) // Duration from now until user's ctx reaches its deadline
			logf("MaxTryTimeout=%d secs, TimeTilDeadline=%d sec\n", timeout, t)
			if t < timeout {
				timeout = t
			}
			if timeout < 0 {
				timeout = 0 // If timeout ever goes negative, set it to zero; this happen while debugging
			}
			logf("TryTimeout adjusted to=%d sec\n", timeout)
		}
		q := requestCopy.Request.URL.Query()
		q.Set("timeout", strconv.Itoa(int(timeout)))
		requestCopy.Request.URL.RawQuery = q.Encode()
		logf("Url=%s\n", requestCopy.Request.URL.String())

		// Set the time for this particular retry operation and then Do the operation.
		tryCtx, tryCancel := context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		response, err = p.next.Do(tryCtx, requestCopy) // Make the request
		logf("Err=%v, response=%v\n", err, response)

		action := "" // This MUST get changed within the switch code below
		switch {
		case ctx.Err() != nil:
			action = "NoRetry: Op timeout"
		case !tryingPrimary && response != nil && response.Response().StatusCode == http.StatusNotFound:
			// If attempt was against the secondary & it returned a StatusNotFound (404), then
			// the resource was not found. This may be due to replication delay. So, in this
			// case, we'll never try the secondary again for this operation.
			considerSecondary = false
			action = "Retry: Secondary URL returned 404"
		case err == context.DeadlineExceeded: // tryCtx.Err should also return context.DeadlineExceeded
			action = "Retry: timeout"
		case err != nil:
			// NOTE: Protocol Responder returns non-nil if REST API returns invalid status code for the invoked operation
			if nerr, ok := err.(net.Error); ok && (nerr.Temporary() || nerr.Timeout()) { // We have a network or StorageError
				action = "Retry: net.Error and Temporary() or Timeout()"
			} else {
				action = "NoRetry: unrecognized error"
			}
		default:
			action = "NoRetry: successful HTTP request" // no error
		}

		logf("Action=%s\n", action)
		// fmt.Println(action + "\n") // This is where we could log the retry operation; action is why we're retrying
		if action[0] != 'R' { // Retry only if action starts with 'R'
			if err != nil {
				tryCancel() // If we're returning an error, cancel this current/last per-retry timeout context
			} else {
				// TODO: Right now, we've decided to leak the per-try Context until the user's Context is canceled.
				// Another option is that we wrap the last per-try context in a body and overwrite the Response's Body field with our wrapper.
				// So, when the user closes the Body, the our per-try context gets closed too.
				// Another option, is that the Last Policy do this wrapping for a per-retry context (not for the user's context)
				_ = tryCancel // So, for now, we don't call cancel: cancel()
			}
			break // Don't retry
		}
		// If retrying, cancel the current per-try timeout context
		tryCancel()
	}
	return response, err // Not retryable or too many retries; return the last response/error
}
