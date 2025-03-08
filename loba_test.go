package loba

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ssenthilnathan3/loba/internal/config"
)

func TestBalanceLoads(t *testing.T) {
	tests := []struct {
		name     string
		targets  []config.Target
		wantCode int
		wantURL  string
	}{
		{
			name: "Valid server available",
			targets: []config.Target{
				{Host_name: "server1", Host_url: "http://localhost:9001", Capacity: "100", Consumption_rate: "50"},
			},
			wantCode: http.StatusTemporaryRedirect,
			wantURL:  "http://localhost:9001",
		},
		{
			name: "No available servers",
			targets: []config.Target{
				{Host_name: "server1", Host_url: "http://localhost:9001", Capacity: "50", Consumption_rate: "50"},
			},
			wantCode: http.StatusServiceUnavailable,
			wantURL:  "",
		},
		{
			name: "Multiple servers - pick least loaded",
			targets: []config.Target{
				{Host_name: "server1", Host_url: "http://localhost:9001", Capacity: "100", Consumption_rate: "70"},
				{Host_name: "server2", Host_url: "http://localhost:9002", Capacity: "150", Consumption_rate: "50"},
			},
			wantCode: http.StatusTemporaryRedirect,
			wantURL:  "http://localhost:9001",
		},
		{
			name: "All servers overloaded",
			targets: []config.Target{
				{Host_name: "server1", Host_url: "http://localhost:9001", Capacity: "100", Consumption_rate: "100"},
				{Host_name: "server2", Host_url: "http://localhost:9002", Capacity: "150", Consumption_rate: "150"},
			},
			wantCode: http.StatusServiceUnavailable,
			wantURL:  "",
		},
		{
			name: "Invalid capacity values",
			targets: []config.Target{
				{Host_name: "server1", Host_url: "http://localhost:9001", Capacity: "invalid", Consumption_rate: "50"},
			},
			wantCode: http.StatusInternalServerError,
			wantURL:  "",
		},
		{
			name: "Invalid consumption values",
			targets: []config.Target{
				{Host_name: "server1", Host_url: "http://localhost:9001", Capacity: "100", Consumption_rate: "invalid"},
			},
			wantCode: http.StatusInternalServerError,
			wantURL:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Configuration{Targets: tt.targets}
			handler := BalanceLoads(cfg)

			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("Expected status %d, got %d", tt.wantCode, rec.Code)
			}

			if tt.wantURL != "" {
				gotURL := rec.Header().Get("Location")
				if gotURL != tt.wantURL {
					t.Errorf("Expected redirect to %s, got %s", tt.wantURL, gotURL)
				}
			}
		})
	}
}
