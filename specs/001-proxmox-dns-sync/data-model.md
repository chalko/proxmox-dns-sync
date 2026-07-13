# Data Models: Proxmox DNS Sync

This document defines the key entities and in-memory structures used by the synchronization service.

## 1. ProxmoxResource
Represents a raw virtual machine, LXC container, or hypervisor node retrieved from Proxmox VE `/cluster/resources` API.

```go
type ProxmoxResource struct {
    VMID   *int   `json:"vmid"`   // VM/LXC ID (nil for hypervisor nodes)
    Name   string `json:"name"`   // Hostname / Name of the resource
    Type   string `json:"type"`   // "qemu" (VM), "lxc" (Container), or "node" (Hypervisor)
    Status string `json:"status"` // "running", "stopped", "online", etc.
    Node   string `json:"node"`   // Name of the physical hypervisor node hosting this resource
    IP     string `json:"ip"`     // Resolved IP address
}
```

## 2. DNSMapping
Represents a calculated DNS host record targeted for Pi-hole.

```go
type DNSMapping struct {
    Hostname string // FQDN (e.g. "grafana-101.fog.lodge.chalko.com")
    IP       string // Target IP address (e.g. "10.7.82.100")
}
```

## 3. Local State Registry
Persisted to `/var/lib/proxmox-dns-sync/state.json` to keep track of DNS hostnames registered by this tool and prevent accidental pruning of external entries.

```go
type LocalState struct {
    RegisteredRecords map[string]string `json:"registered_records"` // maps FQDN -> IP address
}
```

## 4. Hostname Collision Resolution
When multiple resources share the same name, they are processed as follows:
1. Sort resources sharing a name by `vmid` in ascending order.
2. The resource with the lowest `vmid` (assumed first/primary) gets two mappings:
   - `<name>-<vmid>.<domain>`
   - `<name>.<domain>` (as an alias)
3. Other resources with the same name get:
   - `<name>-<vmid>.<domain>`
