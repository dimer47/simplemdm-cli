package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "fixtures", name))
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", name, err)
	}
	return data
}

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("test-key")
	if c.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL %s, got %s", defaultBaseURL, c.baseURL)
	}
	if c.apiKey != "test-key" {
		t.Errorf("expected apiKey test-key, got %s", c.apiKey)
	}
	if c.debug {
		t.Error("expected debug to be false")
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	hc := &http.Client{}
	c := NewClient("key",
		WithBaseURL("https://custom.api.com"),
		WithHTTPClient(hc),
		WithDebug(true),
	)
	if c.baseURL != "https://custom.api.com" {
		t.Errorf("expected custom baseURL, got %s", c.baseURL)
	}
	if c.httpClient != hc {
		t.Error("expected custom http client")
	}
	if !c.debug {
		t.Error("expected debug to be true")
	}
}

func TestClient_Get(t *testing.T) {
	fixture := loadFixture(t, "devices.json")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/devices" {
			t.Errorf("expected /devices, got %s", r.URL.Path)
		}
		// Check Basic Auth
		user, _, ok := r.BasicAuth()
		if !ok || user != "test-key" {
			t.Errorf("expected basic auth with user test-key")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	data, err := c.Get("/devices")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if _, ok := result["data"]; !ok {
		t.Error("expected 'data' key in response")
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	data, err := c.Post("/devices/121/refresh", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"status":"ok"}` {
		t.Errorf("expected ok response for 204, got %s", string(data))
	}
}

func TestClient_DoForm(t *testing.T) {
	fixture := loadFixture(t, "device.json")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/x-www-form-urlencoded" {
			t.Errorf("expected form content type, got %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		if len(body) == 0 {
			t.Error("expected non-empty body")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	data, err := c.DoForm("PATCH", "/devices/121", map[string]string{"name": "New Name"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty response")
	}
}

func TestClient_DoJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("expected json content type, got %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		var parsed map[string]interface{}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Errorf("expected valid JSON body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	data, err := c.DoJSON("POST", "/test", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"status":"ok"}` {
		t.Errorf("unexpected response: %s", string(data))
	}
}

func TestClient_ErrorResponse(t *testing.T) {
	fixture := loadFixture(t, "error_401.json")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(fixture)
	}))
	defer server.Close()

	c := NewClient("bad-key", WithBaseURL(server.URL))
	_, err := c.Get("/account")
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if got := err.Error(); got != "API error (status 401): Unauthorized" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestClient_Error404(t *testing.T) {
	fixture := loadFixture(t, "error_404.json")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(fixture)
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	_, err := c.Get("/devices/999")
	if err == nil {
		t.Fatal("expected error for 404")
	}
	if got := err.Error(); got != "API error (status 404): Record not found" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := NewClient("test-key", WithBaseURL(server.URL))
	data, err := c.Delete("/devices/121")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"status":"ok"}` {
		t.Errorf("expected ok response, got %s", string(data))
	}
}

func TestClient_SetDebug(t *testing.T) {
	c := NewClient("key")
	if c.debug {
		t.Error("expected debug false initially")
	}
	c.SetDebug(true)
	if !c.debug {
		t.Error("expected debug true after SetDebug")
	}
}
