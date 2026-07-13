# Implementation Plan: Proxmox DNS Sync

**Branch**: `001-proxmox-dns-sync` | **Date**: 2026-07-13 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/001-proxmox-dns-sync/spec.md`

## Summary
Implement a synchronization service that queries Proxmox VE API for running VMs and LXC containers, maps their hostnames to their first active lab IPv4 addresses, and synchronizes these records to Pi-hole local DNS using the Pi-hole v6 REST API. The sync runs automatically every 5 minutes using systemd timer.

## Technical Context

**Language/Version**: Go 1.22+

**Primary Dependencies**: Standard library (`net/http`, etc.)

**Storage**: Local state file (`/var/lib/proxmox-dns-sync/state.json`) to track registered hostnames.

**Testing**: `go test` with standard library testing framework

**Target Platform**: Linux helper node/server

**Project Type**: CLI / Daemon

**Performance Goals**: DNS Sync execution completed in < 2 seconds.

**Constraints**: API failure tolerance (log error, do not fail/crash), hostname collision resolution via VM ID suffix, first VM (lowest ID) receives hostname alias. Must NOT prune or modify records created by other jobs (e.g. Kubernetes ExternalDNS); only prune records registered in the local state file.

**Scale/Scope**: ~100 VM/LXC nodes on local lab network.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **I. Code Quality & Standards**: Must use strict typing, pass `go vet` / golangci-lint, and achieve >90% test coverage with `go test`.
- **II. Focused Microservice Design**: Project structured as a single-responsibility Go command CLI utility with distinct packages for Proxmox VE clients, Pi-hole API clients, and the synchronization engine.
- **III. Test-Driven Development (TDD)**: Mock HTTP handler tests must be written first to verify DNS target generation, collision logic, and HTTP query patterns.
- **IV. Integration & Contract Testing**: Implement mock-based unit tests for Pi-hole v6 REST API contracts.
- **V. Observability & Simplicity**: Implement structured JSON logging (using Go `slog` library) mapping to stdout for systemd-journald consumption.

## Project Structure

### Documentation (this feature)

```text
specs/001-proxmox-dns-sync/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 data modeling
├── quickstart.md        # Run/Validation guide
├── checklists/
│   └── requirements.md  # Spec checklist
└── contracts/
    └── pihole-api.md    # Pi-hole API contract
```

### Source Code (repository root)

```text
cmd/
└── proxmox-dns-sync/
    └── main.go          # Main entrypoint

internal/
├── proxmox/
│   ├── client.go        # Proxmox API client
│   └── client_test.go
├── pihole/
│   ├── client.go        # Pi-hole API client
│   └── client_test.go
└── sync/
    ├── sync.go          # Sync logic & collision resolution
    └── sync_test.go

deploy/
├── helm/
│   └── external-dns-values.yaml  # Helm configuration for ExternalDNS
└── k8s/
    ├── secret-store.yaml         # ESO ClusterSecretStore targeting Vault
    └── external-secret.yaml      # ESO ExternalSecret mapping Vault credentials

go.mod
```

**Structure Decision**: Go project structure with sub-packages inside `internal/` to encapsulate concerns.

## Complexity Tracking

No violations of constitution detected.
