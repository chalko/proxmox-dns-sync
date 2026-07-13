package sync

import "testing"

func TestResolveCollisions(t *testing.T) {
	vmid101 := 101
	vmid102 := 102
	vmid103 := 103

	resources := []ProxmoxResource{
		{VMID: &vmid102, Name: "grafana", Type: "qemu", Status: "running", IP: "10.7.82.102"},
		{VMID: &vmid101, Name: "grafana", Type: "qemu", Status: "running", IP: "10.7.82.101"},
		{VMID: &vmid103, Name: "unique-node", Type: "lxc", Status: "running", IP: "10.7.82.103"},
	}

	mappings := ResolveCollisions(resources, "fog.lodge.chalko.com")

	// We expect:
	// - unique-node-103.fog.lodge.chalko.com -> 10.7.82.103
	// - unique-node.fog.lodge.chalko.com -> 10.7.82.103 (alias, since it's the only/first)
	// - grafana-101.fog.lodge.chalko.com -> 10.7.82.101
	// - grafana.fog.lodge.chalko.com -> 10.7.82.101 (alias for lowest ID)
	// - grafana-102.fog.lodge.chalko.com -> 10.7.82.102

	expectedCount := 5
	if len(mappings) != expectedCount {
		t.Fatalf("Expected %d mappings, got %d: %v", expectedCount, len(mappings), mappings)
	}

	mappingMap := make(map[string]string)
	for _, m := range mappings {
		mappingMap[m.Hostname] = m.IP
	}

	checks := map[string]string{
		"unique-node-103.fog.lodge.chalko.com": "10.7.82.103",
		"unique-node.fog.lodge.chalko.com":     "10.7.82.103",
		"grafana-101.fog.lodge.chalko.com":     "10.7.82.101",
		"grafana.fog.lodge.chalko.com":         "10.7.82.101",
		"grafana-102.fog.lodge.chalko.com":     "10.7.82.102",
	}

	for host, ip := range checks {
		resIP, ok := mappingMap[host]
		if !ok {
			t.Errorf("Expected host %q to be registered", host)
		} else if resIP != ip {
			t.Errorf("Expected host %q to have IP %q, got %q", host, ip, resIP)
		}
	}
}

func TestCalculateSyncDiff(t *testing.T) {
	desired := []DNSMapping{
		{Hostname: "grafana.fog.lodge.chalko.com", IP: "10.7.82.101"},
		{Hostname: "influxdb.fog.lodge.chalko.com", IP: "10.7.82.102"},
	}

	state := &LocalState{
		RegisteredRecords: map[string]string{
			"grafana.fog.lodge.chalko.com": "10.7.82.101",
			"stale-node.fog.lodge.chalko.com": "10.7.82.199",
		},
	}

	// Active hosts in Pi-hole
	// - grafana.fog.lodge.chalko.com (needs updating if IP changes, but here matches)
	// - stale-node.fog.lodge.chalko.com (needs deletion because it's in our state but not desired)
	// - external-dns-record.fog.lodge.chalko.com (MUST NOT be deleted because it is NOT in our local state!)
	currentPihole := []string{
		"10.7.82.101 grafana.fog.lodge.chalko.com",
		"10.7.82.199 stale-node.fog.lodge.chalko.com",
		"10.7.82.250 external-dns-record.fog.lodge.chalko.com",
	}

	toAdd, toDelete := CalculateSyncDiff(desired, state, currentPihole)

	// We expect:
	// - toAdd: influxdb.fog.lodge.chalko.com -> 10.7.82.102
	// - toDelete: stale-node.fog.lodge.chalko.com -> 10.7.82.199
	// - external-dns-record.fog.lodge.chalko.com is protected and NOT in toDelete.

	if len(toAdd) != 1 || toAdd[0].Hostname != "influxdb.fog.lodge.chalko.com" {
		t.Errorf("Unexpected toAdd: %v", toAdd)
	}

	if len(toDelete) != 1 || toDelete[0].Hostname != "stale-node.fog.lodge.chalko.com" {
		t.Errorf("Unexpected toDelete: %v", toDelete)
	}
}
