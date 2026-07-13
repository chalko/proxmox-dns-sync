package proxmox

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetResources(t *testing.T) {
	// Mock server mimicking Proxmox API behavior
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header
		token := r.Header.Get("Authorization")
		if token != "PVEAPIToken=root@pam!terraform=SECRET" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.URL.Path == "/cluster/resources" {
			data := map[string]interface{}{
				"data": []map[string]interface{}{
					{"type": "node", "name": "pve", "status": "online"},
					{"type": "qemu", "vmid": 101, "name": "grafana", "status": "running", "node": "pve"},
					{"type": "lxc", "vmid": 102, "name": "influxdb", "status": "running", "node": "pve"},
					{"type": "qemu", "vmid": 103, "name": "stopped-vm", "status": "stopped", "node": "pve"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(data)
			return
		}

		// Mock network interface calls
		if r.URL.Path == "/nodes/pve/qemu/101/agent/network-get-interfaces" {
			data := map[string]interface{}{
				"result": []map[string]interface{}{
					{
						"name": "eth0",
						"ip-addresses": []map[string]interface{}{
							{"ip-address-type": "ipv4", "ip-address": "10.7.82.100"},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
			return
		}

		if r.URL.Path == "/nodes/pve/lxc/102/interfaces" {
			data := []map[string]interface{}{
				{
					"name": "eth0",
					"inet": "10.7.82.101/24",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "root@pam!terraform=SECRET")
	resources, err := client.GetResources()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(resources) != 3 { // pve node + vmid 101 + vmid 102
		t.Errorf("Expected 3 resources, got %d", len(resources))
	}

	var hasGrafana, hasNode bool
	for _, res := range resources {
		if res.Name == "grafana" && res.IP == "10.7.82.100" {
			hasGrafana = true
		}
		if res.Type == "node" && res.Name == "pve" {
			hasNode = true
		}
	}

	if !hasGrafana {
		t.Error("Expected resources to include grafana VM with IP 10.7.82.100")
	}
	if !hasNode {
		t.Error("Expected resources to include pve hypervisor node")
	}
}
