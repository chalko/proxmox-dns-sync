# Technical Research: Proxmox DNS Sync

## 1. Programming Language & Runtime
- **Decision**: Go 1.22+ using standard library (specifically `net/http` and `encoding/json`).
- **Rationale**: Go produces a single, statically linked binary with zero dependencies. This provides extremely fast startup time, minimal CPU/memory footprint, and ease of delivery (just copy the compiled binary to `/usr/local/bin/` on the target host).
- **Alternatives Considered**: Python (rejected because python script delivery requires managing python virtual environments or installing dependencies like `requests` system-wide, which adds packaging complexity and increases execution overhead/startup latency).

## 2. Pi-hole Integration
- **Decision**: Authenticate and configure using the Pi-hole v6 REST API.
- **Rationale**: Pi-hole v6 introduces a robust REST API schema with proper HTTP methods and headers (using `X-FTL-SID` or `sid` cookies/headers). The legacy `/admin/api.php` is deprecated.
- **API Target**:
  - `GET /api/config/dns/hosts` - List custom DNS hosts
  - `POST /api/config/dns/hosts` - Add a custom host mapping
  - `DELETE /api/config/dns/hosts/<ip>%20<host>` - Delete a host mapping

## 3. Collision Handling Strategy
- **Decision**: Sort VMs/LXCs by creation time or VM ID. The lowest VM ID gets the base hostname alias (e.g., `grafana.fog.lodge.chalko.com`) and its own unique ID-appended record (`grafana-101.fog.lodge.chalko.com`). Subsequent VMs with duplicate hostnames only get their unique ID-appended records (e.g. `grafana-102.fog.lodge.chalko.com`).
- **Rationale**: Guarantees resolution for all running nodes while preserving expected default access for the primary/original service node.
