# Feature Specification: Proxmox DNS Sync

**Feature Branch**: `001-proxmox-dns-sync`

**Created**: 2026-07-13

**Status**: Draft

**Input**: User description: "read /home/nick/src/lodge/plans/dynamic_dns_api_push.md"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Sync VM/LXC states to Pi-hole (Priority: P1)

As a system administrator, I want running VMs and LXCs in Proxmox to automatically register their hostnames in Pi-hole DNS under `*.fog.lodge.chalko.com`, so that I can resolve local node hostnames directly.

**Why this priority**: Core sync functionality. Without registering running nodes, the DNS sync service serves no purpose.

**Independent Test**: Create/start a VM/LXC node, run the sync script manually, and confirm the hostname resolves to the VM's IP address.

**Acceptance Scenarios**:

1. **Given** a running VM named `test-vm-1` with IP `10.7.82.100` on Proxmox, **When** the sync script runs, **Then** a DNS A record for `test-vm-1.fog.lodge.chalko.com` pointing to `10.7.82.100` is added to Pi-hole.
2. **Given** a running LXC container named `test-container-1` with IP `10.7.82.101` on Proxmox, **When** the sync script runs, **Then** a DNS A record for `test-container-1.fog.lodge.chalko.com` pointing to `10.7.82.101` is added to Pi-hole.

---

### User Story 2 - Prune stale DNS records (Priority: P2)

As a system administrator, I want stopped or deleted Proxmox nodes to have their DNS records removed from Pi-hole automatically.

**Why this priority**: Essential for maintaining clean and correct DNS state, avoiding routing to offline resources.

**Independent Test**: Stop a running VM/LXC node, run the sync script, and verify that its DNS record is deleted from Pi-hole.

**Acceptance Scenarios**:

1. **Given** a stopped VM `test-vm-1` whose DNS A record `test-vm-1.fog.lodge.chalko.com` exists in Pi-hole, **When** the sync script runs, **Then** the record `test-vm-1.fog.lodge.chalko.com` is deleted from Pi-hole.

---

### User Story 3 - Periodic automated sync (Priority: P3)

As a system administrator, I want the synchronization to run automatically every 5 minutes in the background.

**Why this priority**: Automates the process to keep DNS up to date without manual trigger.

**Independent Test**: Start a new VM, wait 5 minutes, and verify that the hostname resolves.

**Acceptance Scenarios**:

1. **Given** the systemd timer is active, **When** a VM starts, **Then** within 5 minutes the VM's hostname is resolvable via Pi-hole.

---

### Edge Cases

- **Multiple IP Addresses**: A VM may have multiple IPv4 addresses or both IPv4 and IPv6. The system needs a rule for which IP to sync.
- **Unreachable APIs**: If either Proxmox VE or Pi-hole API is temporarily unreachable, the sync must not crash, should log a warning/error, and retry on the next interval.
- **Name Collisions**: If multiple VMs/LXCs have the same hostname, the system needs to handle the collision gracefully.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST retrieve VM and LXC lists and statuses from the Proxmox VE API.
- **FR-002**: System MUST retrieve network interface IP addresses for each running VM and LXC.
- **FR-003**: System MUST identify active VMs/LXCs (status: running) and compile target DNS mappings.
- **FR-004**: System MUST push new DNS records and delete stale/non-running records via the Pi-hole REST API. The system MUST NOT delete any DNS records created by other services (e.g. Kubernetes ExternalDNS) or users.
- **FR-005**: System MUST maintain a local state file registry (e.g. `/var/lib/proxmox-dns-sync/state.json`) containing records created by this service. The system MUST only prune records that are present in both the local state file and the Pi-hole active host list, but are no longer active on Proxmox.
- **FR-006**: System MUST run securely, sourcing credentials from `/etc/default/sync-proxmox-dns` and executing via systemd service/timer.
- **FR-007**: System MUST log all DNS additions, deletions, and failures using structured logging.
- **FR-008**: System MUST handle duplicate VM hostnames by appending the VM/LXC ID to the hostname to guarantee uniqueness, with the oldest VM/LXC (by ID or creation date) additionally receiving the raw hostname as an alias.
- **FR-009**: System MUST authenticate with Pi-hole API using the Pi-hole v6 REST API schema (using HTTP API endpoints and passwords/tokens).

### Key Entities *(include if feature involves data)*

- **Proxmox Node**: Represents a VM or LXC container with ID, name, status, and network interfaces.
- **DNS Record**: Represents a mapping between an FQDN (`<name>.fog.lodge.chalko.com`) and an IP address in Pi-hole.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of running Proxmox VMs and LXCs with valid IPs have correct DNS records in Pi-hole.
- **SC-002**: DNS records for stopped or deleted VMs/LXCs are cleaned up within 5 minutes.
- **SC-003**: Synchronization completes execution in under 2 seconds under normal load (e.g., <100 nodes).
- **SC-004**: Synchronization operates continuously without manual restart or intervention.

## Assumptions

- Proxmox VMs/LXCs have QEMU guest agent installed/configured or network interfaces are queryable via Proxmox API.
- Pi-hole is accessible over HTTP/HTTPS from the execution host.
- The domain suffix `fog.lodge.chalko.com` is used for all synced hosts.
