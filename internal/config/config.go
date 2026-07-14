package config

import (
	"os"
	"strings"
)

// Config holds the application configuration.
type Config struct {
	PveURL         string
	PveToken       string
	PiholeURL      string
	PiholePassword string
	DomainSuffix   string
	StateFilePath  string
}

// LoadConfig retrieves configuration from environment variables.
func LoadConfig() *Config {
	pveURL := os.Getenv("PM_API_URL")
	pveToken := os.Getenv("PROXMOX_VE_API_TOKEN")
	piholeURL := os.Getenv("PIHOLE_URL")
	piholePassword := os.Getenv("TF_VAR_pihole_password")
	domainSuffix := os.Getenv("DOMAIN_SUFFIX")
	if domainSuffix == "" {
		domainSuffix = "fog.lodge.chalko.com"
	}
	stateFilePath := os.Getenv("STATE_FILE_PATH")
	if stateFilePath == "" {
		stateFilePath = "/var/lib/proxmox-dns-sync/state.json"
	}

	return &Config{
		PveURL:         strings.TrimSuffix(pveURL, "/"),
		PveToken:       pveToken,
		PiholeURL:      strings.TrimSuffix(piholeURL, "/"),
		PiholePassword: piholePassword,
		DomainSuffix:   domainSuffix,
		StateFilePath:  stateFilePath,
	}
}
