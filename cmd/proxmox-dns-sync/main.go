package main

import (
	"flag"
	"log/slog"
	"os"
	"proxmox-dns-sync/internal/config"
	"proxmox-dns-sync/internal/logger"
	"proxmox-dns-sync/internal/pihole"
	"proxmox-dns-sync/internal/proxmox"
	"proxmox-dns-sync/internal/sync"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Calculate and log changes without applying them to Pi-hole or saving state")
	flag.Parse()

	// Initialize structured logger
	logger.InitLogger()

	slog.Info("Starting Proxmox DNS Sync service")

	// Load configuration
	cfg := config.LoadConfig()
	if cfg.PveURL == "" || cfg.PveToken == "" {
		slog.Error("Missing required Proxmox credentials. Ensure PM_API_URL and PROXMOX_VE_API_TOKEN are set.")
		os.Exit(1)
	}
	if cfg.PiholePassword == "" {
		slog.Error("Missing required Pi-hole credentials. Ensure TF_VAR_pihole_password is set.")
		os.Exit(1)
	}

	// Load local state registry
	slog.Info("Loading local state registry", "path", cfg.StateFilePath)
	state, err := sync.LoadState(cfg.StateFilePath)
	if err != nil {
		slog.Error("Failed to load local state registry", "error", err)
		os.Exit(1)
	}

	// Instantiate clients
	pveClient := proxmox.NewClient(cfg.PveURL, cfg.PveToken)
	piholeClient := pihole.NewClient(cfg.PiholeURL, cfg.PiholePassword)

	// Fetch resources from Proxmox
	slog.Info("Fetching resources from Proxmox VE API")
	resources, err := pveClient.GetResources()
	if err != nil {
		slog.Error("Failed to fetch Proxmox resources", "error", err)
		os.Exit(1)
	}
	slog.Info("Retrieved active Proxmox resources", "count", len(resources))

	// Resolve hostname collisions and compile desired mappings
	desiredMappings := sync.ResolveCollisions(resources, cfg.DomainSuffix)
	slog.Info("Compiled desired DNS mappings", "count", len(desiredMappings))

	// Fetch current DNS mappings from Pi-hole
	slog.Info("Fetching current custom hosts from Pi-hole")
	currentHosts, err := piholeClient.GetCustomHosts()
	if err != nil {
		slog.Error("Failed to fetch custom hosts from Pi-hole", "error", err)
		os.Exit(1)
	}

	// Calculate changes (additions and safe deletions)
	toAdd, toDelete := sync.CalculateSyncDiff(desiredMappings, state, currentHosts)

	slog.Info("Calculated synchronization delta", "additions", len(toAdd), "deletions", len(toDelete))

	if *dryRun {
		slog.Info("[DRY RUN] Execution complete. No modifications applied.")
		for _, add := range toAdd {
			slog.Info("[DRY RUN] Add record", "host", add.Hostname, "ip", add.IP)
		}
		for _, del := range toDelete {
			slog.Info("[DRY RUN] Delete record", "host", del.Hostname, "ip", del.IP)
		}
		return
	}

	stateChanged := false

	// Apply deletions
	for _, del := range toDelete {
		slog.Info("Deleting stale DNS record from Pi-hole", "host", del.Hostname, "ip", del.IP)
		if err := piholeClient.DeleteCustomHost(del.IP, del.Hostname); err != nil {
			slog.Error("Failed to delete record", "host", del.Hostname, "error", err)
		} else {
			delete(state.RegisteredRecords, del.Hostname)
			stateChanged = true
		}
	}

	// Apply additions
	for _, add := range toAdd {
		slog.Info("Adding/Updating DNS record in Pi-hole", "host", add.Hostname, "ip", add.IP)
		if err := piholeClient.AddCustomHost(add.IP, add.Hostname); err != nil {
			slog.Error("Failed to add record", "host", add.Hostname, "error", err)
		} else {
			state.RegisteredRecords[add.Hostname] = add.IP
			stateChanged = true
		}
	}

	// Persist updated state registry
	if stateChanged {
		slog.Info("Saving updated local state registry", "path", cfg.StateFilePath)
		if err := sync.SaveState(cfg.StateFilePath, state); err != nil {
			slog.Error("Failed to save state registry", "error", err)
			os.Exit(1)
		}
	}

	slog.Info("Proxmox DNS Sync execution completed successfully")
}
