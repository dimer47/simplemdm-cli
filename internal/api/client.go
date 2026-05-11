package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const defaultBaseURL = "https://a.simplemdm.com/api/v1"

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	debug      bool
}

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c *Client) Get(path string) ([]byte, error) {
	return c.Do("GET", path, nil)
}

func (c *Client) Post(path string, body io.Reader) ([]byte, error) {
	return c.Do("POST", path, body)
}

func (c *Client) Put(path string, body io.Reader) ([]byte, error) {
	return c.Do("PUT", path, body)
}

func (c *Client) Patch(path string, body io.Reader) ([]byte, error) {
	return c.Do("PATCH", path, body)
}

func (c *Client) Delete(path string) ([]byte, error) {
	return c.Do("DELETE", path, nil)
}

func (c *Client) DoJSON(method, path string, body interface{}) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reader = bytes.NewReader(data)
	}
	return c.Do(method, path, reader)
}

func (c *Client) DoForm(method, path string, values map[string]string) ([]byte, error) {
	form := url.Values{}
	for k, v := range values {
		form.Set(k, v)
	}
	return c.doWithContentType(method, path, bytes.NewBufferString(form.Encode()), "application/x-www-form-urlencoded")
}

func (c *Client) DoMultipart(method, path string, fields map[string]string, fileField, filePath string) ([]byte, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("failed to write field %s: %w", k, err)
		}
	}

	if fileField != "" && filePath != "" {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		defer f.Close()

		part, err := writer.CreateFormFile(fileField, filepath.Base(filePath))
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		if _, err := io.Copy(part, f); err != nil {
			return nil, fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return c.doWithContentType(method, path, &buf, writer.FormDataContentType())
}

func (c *Client) DoDownload(path string) ([]byte, string, error) {
	u := c.baseURL + path

	if c.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] GET %s\n", u)
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response: %w", err)
	}

	filename := ""
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		if idx := len("filename="); len(cd) > idx {
			for _, part := range bytes.Split([]byte(cd), []byte(";")) {
				p := bytes.TrimSpace(part)
				if bytes.HasPrefix(p, []byte("filename=")) {
					filename = string(bytes.Trim(p[9:], "\" "))
				}
			}
		}
	}

	return body, filename, nil
}

func (c *Client) Do(method, path string, body io.Reader) ([]byte, error) {
	return c.doWithContentType(method, path, body, "application/json")
}

func (c *Client) doWithContentType(method, path string, body io.Reader, contentType string) ([]byte, error) {
	url := c.baseURL + path

	if c.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s %s\n", method, url)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// SimpleMDM uses Basic Auth with API key as username and empty password
	req.SetBasicAuth(c.apiKey, "")
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if c.debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] Status: %d\n", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == 204 {
		return []byte(`{"status":"ok"}`), nil
	}

	if resp.StatusCode >= 400 {
		var apiErr struct {
			Errors []struct {
				Title string `json:"title"`
			} `json:"errors"`
		}
		if json.Unmarshal(respBody, &apiErr) == nil && len(apiErr.Errors) > 0 {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, apiErr.Errors[0].Title)
		}
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
