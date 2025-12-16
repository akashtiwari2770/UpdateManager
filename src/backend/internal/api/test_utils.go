package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"updatemanager/internal/api/router"
	"updatemanager/internal/service"
	"updatemanager/pkg/database"
)

// setupTestServer creates a test HTTP server with all handlers
func setupTestServer(t *testing.T) (*httptest.Server, *service.ServiceFactory, func()) {
	ctx := context.Background()

	// Connect to test MongoDB
	cfg := database.DefaultConfig()
	cfg.Database = "updatemanager_test"

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Initialize services
	services := service.NewServiceFactory(db.Database)

	// Setup router
	r := router.NewRouter(services)
	handler := r.Handler()

	// Create test server
	server := httptest.NewServer(handler)

	cleanup := func() {
		server.Close()
		db.Disconnect(ctx)
	}

	return server, services, cleanup
}

// makeRequest makes an HTTP request to the test server
func makeRequest(t *testing.T, server *httptest.Server, method, path string, body interface{}, headers map[string]string) *http.Response {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, server.URL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	return resp
}

// parseResponse parses JSON response
func parseResponse(t *testing.T, resp *http.Response, target interface{}) {
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
}
