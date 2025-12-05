package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"moul.io/http2curl"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}

func (c *Client) Request(ctx context.Context, payload io.Reader, method, endpoint string, header map[string]string) Response {
	var curl string

	req, err := createRequest(ctx, endpoint, method, payload)
	if err != nil {
		return responseError(err, endpoint, method, curl)
	}

	addCustomHeader(req, header)

	curl, err = getCurl(req)
	if err != nil {
		return responseError(err, req.RequestURI, req.Method, curl)
	}

	response, err := c.sendRequest(req)
	if err != nil {
		return responseError(err, req.RequestURI, req.Method, curl)
	}

	responseBody, err := readResponse(response)
	if err != nil {
		return responseError(err, req.RequestURI, req.Method, curl)
	}

	return NewResponse(response.StatusCode, responseBody, response.Request.URL.String(), response.Request.Method, curl, nil)
}

func createRequest(ctx context.Context, url, method string, payload io.Reader) (*http.Request, error) {
	var req *http.Request
	var err error

	req, err = http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return req, err
}

func readResponse(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return body, fmt.Errorf("failed to read response body: %w", err)
	}

	defer res.Body.Close()

	return body, nil
}

func addCustomHeader(request *http.Request, headers map[string]string) {
	for k, v := range headers {
		request.Header.Add(k, v)
	}
}

func (c *Client) sendRequest(request *http.Request) (*http.Response, error) {
	res, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return res, nil
}

func getCurl(request *http.Request) (string, error) {
	curl, err := http2curl.GetCurlCommand(request)
	if err != nil {
		return "", fmt.Errorf("failed to get curl command: %w", err)
	}

	return curl.String(), nil
}
