package azblob

import (
	"context"
	"errors"
	"github.com/Azure/azure-pipeline-go/pipeline"
	chk "gopkg.in/check.v1"
	"net/http"
	"net/url"
)

type requestIDTestScenario int

const (
	// Testing scenarios for echoing Client Request ID
	clientRequestIDMissing requestIDTestScenario = 1
	clientRequestIDError   requestIDTestScenario = 2
	clientRequestIDMatch   requestIDTestScenario = 3
	clientRequestIDNoMatch requestIDTestScenario = 4
)

type dummyPolicy struct {
	matchID  string
	scenario requestIDTestScenario
}

func (p dummyPolicy) Do(ctx context.Context, request pipeline.Request) (pipeline.Response, error) {
	var header http.Header = make(map[string][]string)
	var err error

	// Set headers and errors according to each scenario
	switch p.scenario {
	case clientRequestIDMissing:
		header.Add("x-ms-client-request-id", "")
	case clientRequestIDError:
		err = errors.New("error is not nil")
	case clientRequestIDMatch:
		header.Add("x-ms-client-request-id", p.matchID)
	case clientRequestIDNoMatch:
		header.Add("x-ms-client-request-id", "fake-client-request-id")
	default:
		header.Add("x-ms-client-request-id", newUUID().String())
	}

	response := http.Response{Header: header}

	return pipeline.NewHTTPResponse(&response), err
}

func (s *aztestsSuite) TestEchoClientRequestIDMissing(c *chk.C) {
	factory := NewUniqueRequestIDPolicyFactory()

	// Scenario 1: Client Request ID is missing
	policy := factory.New(dummyPolicy{scenario: requestIDTestScenario(1)}, nil)
	request, _ := pipeline.NewRequest("GET", url.URL{}, nil)
	resp, err := policy.Do(context.Background(), request)

	c.Assert(err, chk.IsNil)
	c.Assert(resp, chk.NotNil)
}

func (s *aztestsSuite) TestEchoClientRequestIDError(c *chk.C) {
	factory := NewUniqueRequestIDPolicyFactory()

	// Scenario 2: Do method returns an error
	policy := factory.New(dummyPolicy{scenario: requestIDTestScenario(2)}, nil)
	request, _ := pipeline.NewRequest("GET", url.URL{}, nil)
	resp, err := policy.Do(context.Background(), request)

	c.Assert(err, chk.NotNil)
	c.Assert(resp, chk.NotNil)
}

func (s *aztestsSuite) TestEchoClientRequestIDMatch(c *chk.C) {
	factory := NewUniqueRequestIDPolicyFactory()

	// Scenario 3: Client Request ID matches
	matchRequestID := newUUID().String()
	policy := factory.New(dummyPolicy{matchID: matchRequestID, scenario: requestIDTestScenario(3)}, nil)
	request, _ := pipeline.NewRequest("GET", url.URL{}, nil)
	request.Header.Set(xMsClientRequestID, matchRequestID)
	resp, err := policy.Do(context.Background(), request)

	c.Assert(err, chk.IsNil)
	c.Assert(resp, chk.NotNil)
}

func (s *aztestsSuite) TestEchoClientRequestIDNoMatch(c *chk.C) {
	factory := NewUniqueRequestIDPolicyFactory()

	// Scenario 4: Client Request ID does not match
	matchRequestID := newUUID().String()
	policy := factory.New(dummyPolicy{matchID: matchRequestID, scenario: requestIDTestScenario(4)}, nil)
	request, _ := pipeline.NewRequest("GET", url.URL{}, nil)
	request.Header.Set(xMsClientRequestID, matchRequestID)
	resp, err := policy.Do(context.Background(), request)

	c.Assert(err, chk.NotNil)
	c.Assert(err.Error(), chk.Equals, "client Request ID from request and response does not match")
	c.Assert(resp, chk.NotNil)
}
