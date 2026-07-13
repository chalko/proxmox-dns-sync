# Quickstart Validation Guide: Proxmox DNS Sync

This guide details the commands and steps to validate that the synchronization tool functions correctly before deployment.

## Prerequisites
- Go 1.22+

## 1. Unit & Regression Tests (TDD validation)
Run the test suite to verify collision logic, API contract simulations, and mapping logic:
```bash
go test -v ./...
```

## 2. Dry-Run Verification (Manual validation)
Verify changes without sending modifications to Pi-hole:
```bash
# Set dummy credentials
export PM_API_URL="https://pve-mock-host:8006/api2/json"
export PROXMOX_VE_API_TOKEN="dummy-token"
export TF_VAR_pihole_password="dummy-password"

# Build and run in dry-run mode (if implemented)
go build -o proxmox-dns-sync ./cmd/proxmox-dns-sync
./proxmox-dns-sync --dry-run
```

## 3. Deployment Validation
Once installed as a systemd service, verify execution logs using:
```bash
journalctl -u sync-proxmox-dns.service -n 50 --no-pager
```
