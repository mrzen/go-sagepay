package sagepay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

// Client represents a SagePay API client instance
type Client struct {
	// HTTP Client - This can be overridden to change the
	// HTTP transport behaviour.
	HTTP *http.Client

	// A writer which will be used to write out raw requests
	// for diganostics.
	DebugWriter io.Writer

	provider            CredentialsProvider
	testMode            bool
	sessionKey          string
	sessionKeyExpiresAt time.Time
}

const (
	// TestHost is the address of the test API
	TestHost = "https://pi-test.sagepay.com/api/v1"

	// ProductionHost is the address of the production API
	ProductionHost = "https://pi-live.sagepay.com/api/v1"
)

// New creates a new Sagepay API Client with the given CredentialsProvider
func New(ctx context.Context, credentials CredentialsProvider) *Client {
	hc := &http.Client{
		Transport: http.DefaultTransport,
	}

	return &Client{
		HTTP:                hc,
		provider:            credentials,
		testMode:            false,
		sessionKey:          "",
		sessionKeyExpiresAt: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

// Do performs the given HTTP Request
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {

	req.Header.Set("User-Agent", "Sagepay-go +https://github.com/mrzen/go-sagepay")

	credentials, err := c.provider.GetCredentials(ctx)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(credentials.Username, credentials.Password)

	// If a DebugWriter is given and the API is in test mode
	// Write out the raw HTTP Request to the given writer.
	//
	// * Note: This feature is available only in test mode
	//  	   to prevent accidental leakage of sensitive data
	//		   within these logs.
	if c.testMode && c.DebugWriter != nil {
		if req.Body != nil {
			fmt.Fprintln(c.DebugWriter, "--------- REQUEST --------")
			cb := new(bytes.Buffer)
			tr := io.TeeReader(req.Body, cb)
			req.Body = ioutil.NopCloser(tr)
			req.Write(os.Stdout)
			req.Body = ioutil.NopCloser(bytes.NewReader(cb.Bytes()))
		}
	}

	return ctxhttp.Do(ctx, c.HTTP, req)
}

func (c *Client) getEndpoint() string {
	if c.testMode {
		return TestHost
	}

	return ProductionHost
}

// JSON performs an HTTP request with a given request body value encoded as JSON
// and decodes the response as JSON into the given response pointer.
func (c *Client) JSON(ctx context.Context, method, path string, body, into interface{}) error {

	buffer := new(bytes.Buffer)

	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequest(method, c.getEndpoint()+path, bytes.NewReader(buffer.Bytes()))

	// Content-Type and Accept headers must be set for Sage to recognize the input.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if err != nil {
		return err
	}

	res, err := c.Do(ctx, req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	// If we get an error response
	// attempt to decode the response body as an error document
	// and return that as an error.
	if res.StatusCode >= 400 {
		errorBody := ErrorResponse{}

		cb := new(bytes.Buffer)
		tr := io.TeeReader(res.Body, cb)

		// If we can’t decode the body, or the body doesn't appear to
		// contain a list of errors, then return the response body itself
		// as an error.
		if err := json.NewDecoder(tr).Decode(&errorBody); err != nil {
			return err
		} else if len(errorBody.Errors) == 0 {
			return errors.New(cb.String())
		}

		return errorBody
	}

	// Decode the response into the given `into` object,
	// returning any error encountered.
	return json.NewDecoder(res.Body).Decode(&into)
}

// SetTestMode determines if the API client will communite with a test
// or production endpoint
func (c *Client) SetTestMode(testMode bool) {
	c.testMode = testMode
}
