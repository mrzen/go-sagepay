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
	// HTTP Client
	HTTP *http.Client

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

// New creates a new Sagepay API Client
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

	if res.StatusCode >= 400 {
		errorBody := ErrorResponse{}

		cb := new(bytes.Buffer)
		tr := io.TeeReader(res.Body, cb)

		if err := json.NewDecoder(tr).Decode(&errorBody); err != nil {
			return err
		}

		if len(errorBody.Errors) == 0 {
			return errors.New(cb.String())
		}

		return errorBody
	}

	return json.NewDecoder(res.Body).Decode(&into)
}

// SetTestMode determines if the API client will communite with a test
// or production endpoint
func (c *Client) SetTestMode(testMode bool) {
	c.testMode = testMode
}
