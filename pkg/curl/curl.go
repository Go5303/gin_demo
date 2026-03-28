package curl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

// SetTimeout sets the global HTTP client timeout
func SetTimeout(d time.Duration) {
	client.Timeout = d
}

// Response wraps the HTTP response
type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// String returns the response body as string
func (r *Response) String() string {
	return string(r.Body)
}

// JSON unmarshals the response body into the given target
func (r *Response) JSON(v any) error {
	return json.Unmarshal(r.Body, v)
}

// Get sends a GET request (uses context.Background)
func Get(rawURL string, params map[string]string, headers map[string]string) (*Response, error) {
	return GetWithCtx(context.Background(), rawURL, params, headers)
}

// GetWithCtx sends a GET request with context
func GetWithCtx(ctx context.Context, rawURL string, params map[string]string, headers map[string]string) (*Response, error) {
	if len(params) > 0 {
		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, fmt.Errorf("invalid url: %w", err)
		}
		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		rawURL = u.String()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	setHeaders(req, headers)
	return doRequest(req)
}

// PostJSON sends a POST request with JSON body (uses context.Background)
func PostJSON(rawURL string, body any, headers map[string]string) (*Response, error) {
	return PostJSONWithCtx(context.Background(), rawURL, body, headers)
}

// PostJSONWithCtx sends a POST request with JSON body and context
func PostJSONWithCtx(ctx context.Context, rawURL string, body any, headers map[string]string) (*Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal body failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	setHeaders(req, headers)
	return doRequest(req)
}

// PostForm sends a POST request with form-urlencoded body (uses context.Background)
func PostForm(rawURL string, data map[string]string, headers map[string]string) (*Response, error) {
	return PostFormWithCtx(context.Background(), rawURL, data, headers)
}

// PostFormWithCtx sends a POST request with form-urlencoded body and context
func PostFormWithCtx(ctx context.Context, rawURL string, data map[string]string, headers map[string]string) (*Response, error) {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setHeaders(req, headers)
	return doRequest(req)
}

func setHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func doRequest(req *http.Request) (*Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
		Header:     resp.Header,
	}, nil
}
