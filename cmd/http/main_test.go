package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMCPRequestProtection(t *testing.T) {
	const initialize = `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1"}}}`

	tests := []struct {
		name        string
		contentType string
		origin      string
		fetchSite   string
		body        string
		wantStatus  int
	}{
		{name: "native client", contentType: "application/json", body: initialize, wantStatus: http.StatusOK},
		{name: "same origin", contentType: "application/json", origin: "same", body: initialize, wantStatus: http.StatusOK},
		{name: "cross origin", contentType: "application/json", origin: "https://untrusted.invalid", body: initialize, wantStatus: http.StatusForbidden},
		{name: "cross-site fetch", contentType: "application/json", fetchSite: "cross-site", body: initialize, wantStatus: http.StatusForbidden},
		{name: "non-JSON request", contentType: "text/plain", body: initialize, wantStatus: http.StatusUnsupportedMediaType},
		{name: "cross-site simple request", contentType: "text/plain", origin: "https://untrusted.invalid", body: initialize, wantStatus: http.StatusForbidden},
		{
			name:        "null-suffixed method does not override",
			contentType: "application/json",
			body:        strings.Replace(initialize, `"method":"initialize"`, `"method":"initialize","method\u0000":"not/a/method"`, 1),
			wantStatus:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(newHandler())
			defer srv.Close()

			req, err := http.NewRequest(http.MethodPost, srv.URL+"/mcp", strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", tt.contentType)
			req.Header.Set("Accept", "application/json, text/event-stream")
			origin := tt.origin
			if origin == "same" {
				origin = srv.URL
			}
			if origin != "" {
				req.Header.Set("Origin", origin)
			}
			if tt.fetchSite != "" {
				req.Header.Set("Sec-Fetch-Site", tt.fetchSite)
			}

			resp, err := srv.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := resp.Body.Close(); err != nil {
					t.Error(err)
				}
			}()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", resp.StatusCode, tt.wantStatus, body)
			}
			if tt.wantStatus == http.StatusOK && !strings.Contains(string(body), `"serverInfo"`) {
				t.Fatalf("initialize response missing serverInfo: %s", body)
			}
		})
	}
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	newHandler().ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Status != "ok" {
		t.Fatalf("status = %q, want ok", body.Status)
	}
}
