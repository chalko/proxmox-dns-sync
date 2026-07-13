package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStateLoadSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "proxmox-dns-sync-test-*")
	if err != nil {
		t.Fatalf("Failed to create tmp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	statePath := filepath.Join(tmpDir, "state.json")

	// 1. Loading non-existent file should return empty state instead of failing
	state, err := LoadState(statePath)
	if err != nil {
		t.Fatalf("Expected no error for missing file, got: %v", err)
	}
	if state == nil || len(state.RegisteredRecords) != 0 {
		t.Errorf("Expected empty state registry, got: %v", state)
	}

	// 2. Add records and save
	state.RegisteredRecords["grafana.fog.lodge.chalko.com"] = "10.7.82.101"
	err = SaveState(statePath, state)
	if err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// 3. Load again and verify
	loadedState, err := LoadState(statePath)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	ip, ok := loadedState.RegisteredRecords["grafana.fog.lodge.chalko.com"]
	if !ok || ip != "10.7.82.101" {
		t.Errorf("Loaded state incorrect: %v", loadedState)
	}
}
