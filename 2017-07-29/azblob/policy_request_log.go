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
	// LogWarningIfTryOverThreshold logs a warning if a tried operation takes longer than the specified
	// duration (-1=no logging; 0=default threshold).
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

func (f *requestLogPolicyFactory) New(next pipeline.Policy, config *pipeline.Configuration) pipeline.Policy {
	return &requestLogPolicy{o: f.o, next: next, config: config}
}

type requestLogPolicy struct {
	next           pipeline.Policy
	config         *pipeline.Configuration
	o              RequestLogOptions
	try            int32
	operationStart time.Time
}

func redactSigQueryParam(rawQuery string) (bool, string) {
	rawQuery = strings.ToLower(rawQuery) // lowercase the string so we can look for ?sig= and &sig=
	sigFound := strings.Contains(rawQuery, "?sig=")
	if !sigFound {
		sigFound = strings.Contains(rawQuery, "&sig=")
		if !sigFound {
			return sigFound, rawQuery // [?|&]sig= not found; return same rawQuery passed in (no memory allocation)
		}
	}
	// [?|&]sig= found, redact its value
	values, _ := url.ParseQuery(rawQuery)
	for name := range values {
		if strings.EqualFold(name, "sig") {
			values[name] = []string{"REDACTED"}
		}
	}
	return sigFound, values.Encode()
}

func prepareRequestForLogging(request pipeline.Request) *http.Request {
	req := request
	if sigFound, rawQuery := redactSigQueryParam(req.URL.RawQuery); sigFound {
		// Make copy so we don't destroy the query parameters we actually need to send in the request
		req = request.Copy()
		req.Request.URL.RawQuery = rawQuery
	}
	return req.Request
}

func (p *requestLogPolicy) Do(ctx context.Context, request pipeline.Request) (response pipeline.Response, err error) {
	p.try++ // The first try is #1 (not #0)
	if p.try == 1 {
		p.operationStart = time.Now() // If this is the 1st try, record the operation state time
	}

	// Log the outgoing request as informational
	if p.config.ShouldLog(pipeline.LogInfo) {
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "==> OUTGOING REQUEST (Try=%d)\n", p.try)
		pipeline.WriteRequestWithResponse(b, prepareRequestForLogging(request), nil, nil)
		p.config.Log(pipeline.LogInfo, b.String())
	}

	// Set the time for this particular retry operation and then Do the operation.
	tryStart := time.Now()
	response, err = p.next.Do(ctx, request) // Make the request
	tryEnd := time.Now()
	tryDuration := tryEnd.Sub(tryStart)
	opDuration := tryEnd.Sub(p.operationStart)

	logLevel, forceLog := pipeline.LogInfo, false // Default logging information

	// If the response took too long, we'll upgrade to warning.
	if p.o.LogWarningIfTryOverThreshold > 0 && tryDuration > p.o.LogWarningIfTryOverThreshold {
		// Log a warning if the try duration exceeded the specified threshold
		logLevel, forceLog = pipeline.LogWarning, true
	}

	if err == nil { // We got a response from the service
		sc := response.Response().StatusCode
		if ((sc >= 400 && sc <= 499) && sc != http.StatusNotFound && sc != http.StatusConflict && sc != http.StatusPreconditionFailed && sc != http.StatusRequestedRangeNotSatisfiable) || (sc >= 500 && sc <= 599) {
			logLevel, forceLog = pipeline.LogError, true // Promote to Error any 4xx (except those listed is an error) or any 5xx
		} else {
			// For other status codes, we leave the level as is.
		}
	} else { // This error did not get an HTTP response from the service; upgrade the severity to Error
		logLevel, forceLog = pipeline.LogError, true
	}

	if shouldLog := p.config.ShouldLog(logLevel); forceLog || shouldLog {
		// We're going to log this; build the string to log
		b := &bytes.Buffer{}
		slow := ""
		if p.o.LogWarningIfTryOverThreshold > 0 && tryDuration > p.o.LogWarningIfTryOverThreshold {
			slow = fmt.Sprintf("[SLOW >%v]", p.o.LogWarningIfTryOverThreshold)
		}
		fmt.Fprintf(b, "==> REQUEST/RESPONSE (Try=%d/%v%s, OpTime=%v) -- ", p.try, tryDuration, slow, opDuration)
		if err != nil { // This HTTP request did not get a response from the service
			fmt.Fprintf(b, "REQUEST ERROR\n")
		} else {
			if logLevel == pipeline.LogError {
				fmt.Fprintf(b, "RESPONSE STATUS CODE ERROR\n")
			} else {
				fmt.Fprintf(b, "RESPONSE SUCCESSFULLY RECEIVED\n")
			}
		}

		pipeline.WriteRequestWithResponse(b, prepareRequestForLogging(request), response.Response(), err)
		if logLevel <= pipeline.LogError {
			b.Write(stack()) // For errors (or lower levels), we append the stack trace (an expensive operation)
		}
		msg := b.String()

		if forceLog {
			pipeline.ForceLog(logLevel, msg)
		}
		if shouldLog {
			p.config.Log(logLevel, msg)
		}
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
