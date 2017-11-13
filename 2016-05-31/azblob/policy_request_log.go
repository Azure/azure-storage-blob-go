package azblob

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// RequestLogOptions configures the retry policy's behavior.
type RequestLogOptions struct {
	// LogWarningIfTryOverThreshold logs a warning if a tried operation takes longer than the specified duration (0=no logging).
	LogWarningIfTryOverThreshold time.Duration
}

// NewRequestLogPolicyFactory creates a RequestLogPolicyFactory object configured using the specified options.
func NewRequestLogPolicyFactory(o RequestLogOptions) pipeline.Factory {
	if o.LogWarningIfTryOverThreshold == 0 {
		// It would be good to relate this to https://azure.microsoft.com/en-us/support/legal/sla/storage/v1_2/
		// But this monitors the time to get the HTTP response; NOT the time to download the response body.
		o.LogWarningIfTryOverThreshold = 3 * time.Second // Default to 3 seconds
	}
	return &requestLogPolicyFactory{o: o}
}

type requestLogPolicyFactory struct {
	o RequestLogOptions
}

func (f *requestLogPolicyFactory) New(node pipeline.Node) pipeline.Policy {
	return &requestLogPolicy{node: node, o: f.o}
}

type requestLogPolicy struct {
	node           pipeline.Node
	o              RequestLogOptions
	try            int
	operationStart time.Time
}

func redactSigQueryParam(rawQuery string) (bool, string) {
	sigFound := strings.EqualFold(rawQuery, "?sig=")
	if !sigFound {
		sigFound = strings.EqualFold(rawQuery, "&sig=")
		if !sigFound {
			return sigFound, rawQuery // [?|&]sig= not found; return same rawQuery passed in (no memory allocation)
		}
	}
	// [?|&]sig= was found, redact its value
	values, _ := url.ParseQuery(rawQuery)
	for name := range values {
		if strings.EqualFold(name, "sig") {
			values[name] = []string{"(redacted)"}
		}
	}
	return sigFound, values.Encode()
}

func (p *requestLogPolicy) Do(ctx context.Context, request pipeline.Request) (response pipeline.Response, err error) {
	p.try++ // The first try is #1 (not #0)
	if p.try == 1 {
		p.operationStart = time.Now() // If this is the 1st try, record the operation state time
	}

	// Log the outgoing request as informational
	if p.node.WouldLog(pipeline.LogInfo) {
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "==> OUTGOING REQUEST (Try=%d)\n", p.try)
		req := request
		if sigFound, rawQuery := redactSigQueryParam(req.URL.RawQuery); sigFound {
			// TODO: Make copy so we dont' destroy the query parameters we actually need to send in the request
			req = request.Copy()
			req.Request.URL.RawQuery = rawQuery
		}
		pipeline.WriteRequest(b, req.Request)
		p.node.Log(pipeline.LogInfo, b.String())
	}

	// Set the time for this particular retry operation and then Do the operation.
	tryStart := time.Now()
	response, err = p.node.Do(ctx, request) // Make the request
	tryEnd := time.Now()
	tryDuration := tryEnd.Sub(tryStart)
	opDuration := tryEnd.Sub(p.operationStart)

	severity := pipeline.LogInfo // Assume success and default to informational logging
	logMsg := func(b *bytes.Buffer) {
		b.WriteString("SUCCESSFUL OPERATION\n")
		pipeline.WriteResponseWithRequest(b, response.Response())
	}

	// If the response took too long, we'll upgrade to warning.
	if p.o.LogWarningIfTryOverThreshold > 0 && tryDuration > p.o.LogWarningIfTryOverThreshold {
		// Log a warning if the try duration exceeded the specified threshold
		severity = pipeline.LogWarning
		logMsg = func(b *bytes.Buffer) {
			fmt.Fprintf(b, "SLOW OPERATION [tryDuration > %v]\n", p.o.LogWarningIfTryOverThreshold)
			pipeline.WriteResponseWithRequest(b, response.Response())
		}
	}

	if err == nil { // We got a response from the service
		sc := response.Response().StatusCode
		if ((sc >= 400 && sc <= 499) && sc != http.StatusNotFound && sc != http.StatusConflict && sc != http.StatusPreconditionFailed && sc != http.StatusRequestedRangeNotSatisfiable) || (sc >= 500 && sc <= 599) {
			severity = pipeline.LogError // Promote to Error any 4xx (except those listed is an error) or any 5xx
			logMsg = func(b *bytes.Buffer) {
				// Write the error, the originating request and the stack
				fmt.Fprintf(b, "OPERATION ERROR:\n%v\n", err)
				pipeline.WriteResponseWithRequest(b, response.Response())
				b.Write(stack()) // For errors, we append the stack trace (an expensive operation)
			}
		} else {
			// For other status codes, we leave the severity as is.
		}
	} else { // This error did not get an HTTP response from the service; upgrade the severity to Error
		severity = pipeline.LogError
		logMsg = func(b *bytes.Buffer) {
			// Write the error, the originating request and the stack
			fmt.Fprintf(b, "NETWORK ERROR:\n%v\n", err)
			pipeline.WriteRequest(b, request.Request)
			b.Write(stack()) // For errors, we append the stack trace (an expensive operation)
		}
	}

	if p.node.WouldLog(severity) || false { // Change false to true for testing
		// We're going to log this; build the string to log
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "==> REQUEST/RESPONSE (Try=%d, TryDuration=%v, OpDuration=%v) -- ", p.try, tryDuration, opDuration)
		logMsg(b)
		p.node.Log(severity, b.String())
	}
	return response, err
}

func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
