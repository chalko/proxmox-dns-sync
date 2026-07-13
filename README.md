# Proxmox DNS Sync

`proxmox-dns-sync` is a lightweight, zero-dependency Go utility designed to automatically synchronize active virtual machines and LXC container hostnames from a **Proxmox VE** cluster directly to a **Pi-hole v6** local DNS server. 

It is designed to run periodically in the background (via a systemd timer or cron) on a helper node, ensuring seamless DNS resolution for lab environments.

---

## Features

- **Statically Linked Binary**: Written in Go, compiling down to a single binary with zero external runtime dependencies.
- **Dynamic Session Authentication**: Implements automatic login handshakes with the Pi-hole v6 `/api/auth` endpoint, supporting standard and application passwords.
- **Hostname Collision Resolution**:
  - Automatically standardizes hostnames (lowercasing and replacing underscores with hyphens).
  - Appends the VM/LXC ID to duplicate hostnames (e.g. `grafana-101.fog.lodge.chalko.com`) to guarantee uniqueness.
  - Maps the lowest VM/LXC ID as an alias to the raw hostname (e.g. `grafana.fog.lodge.chalko.com`).
- **Safe Pruning (Coexistence)**:
  - Maintains a persistent local JSON state registry (`/var/lib/proxmox-dns-sync/state.json`) tracking which records it registered.
  - Only prunes DNS records that it owns, ensuring manual entries and records created by other automated jobs (such as Kubernetes `ExternalDNS`) are **never deleted or modified**.
- **Structured Observability**: Logs in JSON format using Go's structured logging library (`slog`), integrating cleanly with `systemd-journald`.

---

## Directory Structure

```text
├── cmd/
│   └── proxmox-dns-sync/
│       └── main.go          # CLI Entrypoint
├── internal/
│   ├── config/              # Configuration Loader
│   ├── logger/              # Structured slog JSON logger
│   ├── pihole/              # Pi-hole v6 API REST client
│   ├── proxmox/             # Proxmox VE API Client
│   └── sync/                # Synchronization & Collision Resolution
├── deploy/
│   ├── helm/
│   │   └── external-dns-values.yaml # Kubernetes ExternalDNS values
│   ├── k8s/
│   │   └── external-secret.yaml      # External Secrets Operator manifest
│   └── systemd/
│       ├── install.sh                # Helper deployment script
│       ├── sync-proxmox-dns.service  # systemd service definition
│       └── sync-proxmox-dns.timer    # systemd timer definition
```

---

## Configuration

The application is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PM_API_URL` | Proxmox VE API endpoint url (e.g. `https://10.7.82.10:8006/api2/json`) | *Required* |
| `PROXMOX_VE_API_TOKEN` | Proxmox API Token ID and Secret (e.g. `user@pam!tokenid=secret-uuid`) | *Required* |
| `TF_VAR_pihole_password` | Pi-hole Web Admin / Application Password | *Required* |
| `PIHOLE_URL` | Pi-hole base URL | `http://10.5.110.3` |
| `DOMAIN_SUFFIX` | Target DNS subdomain suffix | `fog.lodge.chalko.com` |
| `STATE_FILE_PATH` | Path to persistent local state file | `/var/lib/proxmox-dns-sync/state.json` |

---

## Getting Started

### Prerequisites
- Go 1.22+

### Build & Unit Testing
```bash
# Run the test suite
go test -v ./...

# Build the binary
go build -o proxmox-dns-sync ./cmd/proxmox-dns-sync
```

### Dry Run Verification
Execute the sync in dry-run mode to check mapping outputs without applying changes to Pi-hole or writing state:
```bash
# Load environment credentials
export PM_API_URL="https://10.7.82.10:8006/api2/json"
export PROXMOX_VE_API_TOKEN="root@pam!terraform=SECRET-TOKEN"
export TF_VAR_pihole_password="SECRET-PIHOLE-PASSWORD"

# Run dry run
./proxmox-dns-sync --dry-run
```

---

## Deployment

### systemd Installation (Helper Node)
The project includes systemd service files to run synchronization automatically every 5 minutes:

1. Create a secure environment configuration file at `/etc/default/sync-proxmox-dns` containing your credentials:
   ```bash
   PM_API_URL="https://10.7.82.10:8006/api2/json"
   PROXMOX_VE_API_TOKEN="root@pam!terraform=SECRET-TOKEN"
   PIHOLE_URL="http://10.5.110.3"
   TF_VAR_pihole_password="SECRET-PIHOLE-PASSWORD"
   ```
2. Run the deployment script to compile the binary, register systemd units, and start the timer:
   ```bash
   chmod +x deploy/systemd/install.sh
   ./deploy/systemd/install.sh
   ```
3. Monitor logs:
   ```bash
   journalctl -u sync-proxmox-dns.service -f --no-pager
   ```

### Kubernetes Integration (ExternalDNS)
To synchronize Kubernetes Ingresses to the same Pi-hole instance without collision:
1. Apply the [ExternalSecret](deploy/k8s/external-secret.yaml) to sync the Pi-hole password from Vault via the External Secrets Operator:
   ```bash
   kubectl apply -f deploy/k8s/external-secret.yaml
   ```
2. Deploy the `external-dns` controller using Helm and the provided [values file](deploy/helm/external-dns-values.yaml):
   ```bash
   helm repo add external-dns https://kubernetes-sigs.github.io/external-dns/
   helm upgrade --install external-dns external-dns/external-dns \
     -n external-dns \
     -f deploy/helm/external-dns-values.yaml
   ```
