package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpRequest struct {
	Endpoint string
	Params   *url.Values
	Headers  *map[string]string
}

type Client struct {
	client         *http.Client
	baseUrl        string
	defaultHeaders *map[string]string
}

func NewClient(baseUrl string, defaultHeaders *map[string]string) *Client {
	return &Client{
		client:         &http.Client{},
		baseUrl:        baseUrl,
		defaultHeaders: defaultHeaders,
	}
}

func (c *Client) GetBaseUrl() string {
	return c.baseUrl
}

// Helper method to build HTTP request
func (c *Client) buildRequest(method string, req *HttpRequest, bodyReader io.Reader) (*http.Request, error) {
	url := c.baseUrl + req.Endpoint

	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	c.setQueryParams(request, req.Params)
	c.setHeaders(request, req.Headers)

	return request, nil
}

// Helper to set query parameters
func (c *Client) setQueryParams(request *http.Request, params *url.Values) {
	if params == nil {
		return
	}
	request.URL.RawQuery = params.Encode()
}

// Helper to set headers (including default headers)
func (c *Client) setHeaders(request *http.Request, headers *map[string]string) {
	// Set default headers first
	if c.defaultHeaders != nil {
		for key, value := range *c.defaultHeaders {
			request.Header.Set(key, value)
		}
	}

	// Override with request-specific headers
	if headers != nil {
		for key, value := range *headers {
			request.Header.Set(key, value)
		}
	}
}

// Helper to execute request and handle response
func (c *Client) executeRequest(request *http.Request, result any) (*http.Response, error) {
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if result != nil {
		if err := json.NewDecoder(response.Body).Decode(result); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return response, nil
}

// Simplified HTTP method implementations
func (c *Client) Get(req *HttpRequest, result any) (*http.Response, error) {
	request, err := c.buildRequest("GET", req, nil)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(request, result)
}

func (c *Client) Post(req *HttpRequest, body io.Reader, result any) (*http.Response, error) {
	request, err := c.buildRequest("POST", req, body)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(request, result)
}

func (c *Client) Put(req *HttpRequest, body io.Reader, result any) (*http.Response, error) {
	request, err := c.buildRequest("PUT", req, body)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(request, result)
}

func (c *Client) Delete(req *HttpRequest, result any) (*http.Response, error) {
	request, err := c.buildRequest("DELETE", req, nil)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(request, result)
}
