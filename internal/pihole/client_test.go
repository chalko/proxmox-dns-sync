package pihole

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPiholeClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header/token
		token := r.Header.Get("X-FTL-SID")
		if token != "secret-sid" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" && r.URL.Path == "/api/config/dns/hosts" {
			hosts := []string{
				"10.7.82.10 misty.fog.lodge.chalko.com",
				"10.7.82.100 grafana-101.fog.lodge.chalko.com",
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(hosts)
			return
		}

		if r.Method == "POST" && r.URL.Path == "/api/config/dns/hosts" {
			var body string
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body != "10.7.82.101 influxdb-102.fog.lodge.chalko.com" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == "DELETE" && (r.URL.Path == "/api/config/dns/hosts/10.7.82.100 grafana-101.fog.lodge.chalko.com" || r.URL.Path == "/api/config/dns/hosts/10.7.82.100%20grafana-101.fog.lodge.chalko.com") {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "secret-sid")

	// 1. Test GetCustomHosts
	hosts, err := client.GetCustomHosts()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(hosts) != 2 || hosts[0] != "10.7.82.10 misty.fog.lodge.chalko.com" {
		t.Errorf("Unexpected hosts response: %v", hosts)
	}

	// 2. Test AddCustomHost
	err = client.AddCustomHost("10.7.82.101", "influxdb-102.fog.lodge.chalko.com")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// 3. Test DeleteCustomHost
	err = client.DeleteCustomHost("10.7.82.100", "grafana-101.fog.lodge.chalko.com")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}
