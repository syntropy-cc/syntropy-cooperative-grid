// Package mocks provides mock implementations for testing the setup component
package mocks

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// MockAPIClient provides a mock implementation of API client operations
type MockAPIClient struct {
	// Response configuration
	StatusCode    int
	ResponseBody  interface{}
	ResponseError error
	
	// Request tracking
	Requests []APIRequest
	
	// Behavior configuration
	Delay        time.Duration
	RetryCount   int
	ShouldFail   bool
	FailAfter    int // Fail after N requests
	requestCount int
}

// APIRequest tracks API request details
type APIRequest struct {
	Method   string
	URL      string
	Headers  map[string]string
	Body     []byte
	Context  context.Context
	Timestamp time.Time
}

// NewMockAPIClient creates a new mock API client
func NewMockAPIClient() *MockAPIClient {
	return &MockAPIClient{
		StatusCode:   http.StatusOK,
		ResponseBody: map[string]interface{}{"status": "success"},
		Requests:     make([]APIRequest, 0),
	}
}

// Get mocks HTTP GET requests
func (m *MockAPIClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return m.makeRequest(ctx, "GET", url, headers, nil)
}

// Post mocks HTTP POST requests
func (m *MockAPIClient) Post(ctx context.Context, url string, headers map[string]string, body []byte) (*http.Response, error) {
	return m.makeRequest(ctx, "POST", url, headers, body)
}

// Put mocks HTTP PUT requests
func (m *MockAPIClient) Put(ctx context.Context, url string, headers map[string]string, body []byte) (*http.Response, error) {
	return m.makeRequest(ctx, "PUT", url, headers, body)
}

// Delete mocks HTTP DELETE requests
func (m *MockAPIClient) Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return m.makeRequest(ctx, "DELETE", url, headers, nil)
}

// makeRequest handles the common request logic
func (m *MockAPIClient) makeRequest(ctx context.Context, method, url string, headers map[string]string, body []byte) (*http.Response, error) {
	m.requestCount++
	
	// Track the request
	request := APIRequest{
		Method:    method,
		URL:       url,
		Headers:   headers,
		Body:      body,
		Context:   ctx,
		Timestamp: time.Now(),
	}
	m.Requests = append(m.Requests, request)
	
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	
	// Simulate delay
	if m.Delay > 0 {
		time.Sleep(m.Delay)
	}
	
	// Check if should fail
	if m.ShouldFail || (m.FailAfter > 0 && m.requestCount > m.FailAfter) {
		if m.ResponseError != nil {
			return nil, m.ResponseError
		}
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       http.NoBody,
		}, nil
	}
	
	// Create response body
	var responseBody []byte
	if m.ResponseBody != nil {
		var err error
		responseBody, err = json.Marshal(m.ResponseBody)
		if err != nil {
			return nil, err
		}
	}
	
	// Create mock response
	response := &http.Response{
		StatusCode: m.StatusCode,
		Header:     make(http.Header),
		Body:       &MockResponseBody{data: responseBody},
	}
	
	response.Header.Set("Content-Type", "application/json")
	
	return response, nil
}

// SetResponse configures the mock response
func (m *MockAPIClient) SetResponse(statusCode int, body interface{}) {
	m.StatusCode = statusCode
	m.ResponseBody = body
}

// SetError configures the mock to return an error
func (m *MockAPIClient) SetError(err error) {
	m.ResponseError = err
	m.ShouldFail = true
}

// SetDelay configures response delay
func (m *MockAPIClient) SetDelay(delay time.Duration) {
	m.Delay = delay
}

// SetFailAfter configures the mock to fail after N requests
func (m *MockAPIClient) SetFailAfter(count int) {
	m.FailAfter = count
}

// GetRequestCount returns the number of requests made
func (m *MockAPIClient) GetRequestCount() int {
	return len(m.Requests)
}

// GetLastRequest returns the last request made
func (m *MockAPIClient) GetLastRequest() *APIRequest {
	if len(m.Requests) == 0 {
		return nil
	}
	return &m.Requests[len(m.Requests)-1]
}

// GetRequestsByMethod returns requests filtered by HTTP method
func (m *MockAPIClient) GetRequestsByMethod(method string) []APIRequest {
	var filtered []APIRequest
	for _, req := range m.Requests {
		if req.Method == method {
			filtered = append(filtered, req)
		}
	}
	return filtered
}

// GetRequestsByURL returns requests filtered by URL
func (m *MockAPIClient) GetRequestsByURL(url string) []APIRequest {
	var filtered []APIRequest
	for _, req := range m.Requests {
		if req.URL == url {
			filtered = append(filtered, req)
		}
	}
	return filtered
}

// Reset clears all requests and resets configuration
func (m *MockAPIClient) Reset() {
	m.StatusCode = http.StatusOK
	m.ResponseBody = map[string]interface{}{"status": "success"}
	m.ResponseError = nil
	m.Requests = make([]APIRequest, 0)
	m.Delay = 0
	m.RetryCount = 0
	m.ShouldFail = false
	m.FailAfter = 0
	m.requestCount = 0
}

// MockResponseBody implements io.ReadCloser for HTTP response body
type MockResponseBody struct {
	data   []byte
	offset int
}

// Read implements io.Reader
func (m *MockResponseBody) Read(p []byte) (n int, err error) {
	if m.offset >= len(m.data) {
		return 0, nil // EOF
	}
	
	n = copy(p, m.data[m.offset:])
	m.offset += n
	return n, nil
}

// Close implements io.Closer
func (m *MockResponseBody) Close() error {
	return nil
}

// MockHTTPServer provides a mock HTTP server for integration testing
type MockHTTPServer struct {
	Server   *http.Server
	Handlers map[string]http.HandlerFunc
	Requests []APIRequest
}

// NewMockHTTPServer creates a new mock HTTP server
func NewMockHTTPServer() *MockHTTPServer {
	return &MockHTTPServer{
		Handlers: make(map[string]http.HandlerFunc),
		Requests: make([]APIRequest, 0),
	}
}

// AddHandler adds a handler for a specific path
func (m *MockHTTPServer) AddHandler(path string, handler http.HandlerFunc) {
	m.Handlers[path] = handler
}

// Start starts the mock server on the specified address
func (m *MockHTTPServer) Start(addr string) error {
	mux := http.NewServeMux()
	
	// Add all handlers
	for path, handler := range m.Handlers {
		mux.HandleFunc(path, handler)
	}
	
	m.Server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	
	return m.Server.ListenAndServe()
}

// Stop stops the mock server
func (m *MockHTTPServer) Stop(ctx context.Context) error {
	if m.Server != nil {
		return m.Server.Shutdown(ctx)
	}
	return nil
}

// GetRequestCount returns the number of requests received
func (m *MockHTTPServer) GetRequestCount() int {
	return len(m.Requests)
}