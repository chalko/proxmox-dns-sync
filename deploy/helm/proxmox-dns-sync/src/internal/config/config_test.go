package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("PM_API_URL", "https://pve-test:8006/api2/json")
	os.Setenv("PROXMOX_VE_API_TOKEN", "test-token")
	os.Setenv("TF_VAR_pihole_password", "test-password")
	os.Setenv("DOMAIN_SUFFIX", "test.suffix")
	os.Setenv("STATE_FILE_PATH", "/tmp/test-state.json")

	cfg := LoadConfig()

	if cfg.PveURL != "https://pve-test:8006/api2/json" {
		t.Errorf("Expected PveURL to be https://pve-test:8006/api2/json, got %q", cfg.PveURL)
	}
	if cfg.PveToken != "test-token" {
		t.Errorf("Expected PveToken to be test-token, got %q", cfg.PveToken)
	}
	if cfg.PiholePassword != "test-password" {
		t.Errorf("Expected PiholePassword to be test-password, got %q", cfg.PiholePassword)
	}
	if cfg.DomainSuffix != "test.suffix" {
		t.Errorf("Expected DomainSuffix to be test.suffix, got %q", cfg.DomainSuffix)
	}
	if cfg.StateFilePath != "/tmp/test-state.json" {
		t.Errorf("Expected StateFilePath to be /tmp/test-state.json, got %q", cfg.StateFilePath)
	}
}
