#!/bin/bash
set -e

# Compile the binary
echo "Building proxmox-dns-sync..."
go build -o proxmox-dns-sync ./cmd/proxmox-dns-sync

# Copy binary to destination
echo "Installing binary to /usr/local/bin..."
sudo cp proxmox-dns-sync /usr/local/bin/

# Copy systemd units
echo "Installing systemd files..."
sudo cp deploy/systemd/sync-proxmox-dns.service /etc/systemd/system/
sudo cp deploy/systemd/sync-proxmox-dns.timer /etc/systemd/system/

# Reload systemd and enable timer
echo "Reloading systemd and enabling timer..."
sudo systemctl daemon-reload
sudo systemctl enable --now sync-proxmox-dns.timer

echo "Installation complete!"
