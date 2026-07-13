package sync

// ProxmoxResource represents raw node metadata fetched from Proxmox cluster resources.
type ProxmoxResource struct {
	VMID   *int   `json:"vmid"`   // VM/LXC ID (nil for hypervisor nodes)
	Name   string `json:"name"`   // Hostname / Name of the resource
	Type   string `json:"type"`   // "qemu" (VM), "lxc" (Container), or "node" (Hypervisor)
	Status string `json:"status"` // "running", "stopped", "online", etc.
	Node   string `json:"node"`   // Name of the physical hypervisor node hosting this resource
	IP     string `json:"ip"`     // Resolved IP address
}

// DNSMapping represents a calculated DNS host record targeted for Pi-hole.
type DNSMapping struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

// LocalState represents the persisted list of hosts registered by this service.
type LocalState struct {
	RegisteredRecords map[string]string `json:"registered_records"` // maps FQDN -> IP address
}
