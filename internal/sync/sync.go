package sync

import (
	"fmt"
	"sort"
	"strings"
)

// ResolveCollisions implements the hostname collision sorting, suffixing, and alias generation.
func ResolveCollisions(resources []ProxmoxResource, domain string) []DNSMapping {
	// Group resources by name
	groups := make(map[string][]ProxmoxResource)
	for _, res := range resources {
		if res.VMID == nil || res.IP == "" {
			continue
		}
		// Standardize name for DNS: lowercase, replace underscores with hyphens
		dnsName := strings.ToLower(res.Name)
		dnsName = strings.ReplaceAll(dnsName, "_", "-")
		groups[dnsName] = append(groups[dnsName], res)
	}

	var mappings []DNSMapping

	for dnsName, group := range groups {
		// Sort by VMID ascending
		sort.Slice(group, func(i, j int) bool {
			return *group[i].VMID < *group[j].VMID
		})

		for idx, res := range group {
			// Suffix FQDN: <name>-<vmid>.<domain>
			suffixHostname := fmt.Sprintf("%s-%d.%s", dnsName, *res.VMID, domain)
			mappings = append(mappings, DNSMapping{
				Hostname: suffixHostname,
				IP:       res.IP,
			})

			// Oldest VMID (first one in sorted slice) gets the raw hostname as an alias
			if idx == 0 {
				rawHostname := fmt.Sprintf("%s.%s", dnsName, domain)
				mappings = append(mappings, DNSMapping{
					Hostname: rawHostname,
					IP:       res.IP,
				})
			}
		}
	}

	return mappings
}

// CalculateSyncDiff computes the additions and deletions for Pi-hole.
// It ensures that deletions only apply to records present in our local state registry.
func CalculateSyncDiff(desired []DNSMapping, state *LocalState, currentPihole []string) (toAdd []DNSMapping, toDelete []DNSMapping) {
	// Parse current Pi-hole entries into a map for fast lookup: FQDN -> IP
	piholeRecords := make(map[string]string)
	for _, entry := range currentPihole {
		parts := strings.Fields(entry)
		if len(parts) >= 2 {
			ip := parts[0]
			host := strings.ToLower(parts[1])
			piholeRecords[host] = ip
		}
	}

	// Helper maps for desired records
	desiredMap := make(map[string]string)
	for _, m := range desired {
		host := strings.ToLower(m.Hostname)
		desiredMap[host] = m.IP
	}

	// 1. Calculate additions/updates
	for _, m := range desired {
		host := strings.ToLower(m.Hostname)
		currentIP, exists := piholeRecords[host]
		if !exists || currentIP != m.IP {
			toAdd = append(toAdd, m)
		}
	}

	// 2. Calculate safe deletions:
	// Only delete if the record is present in both Pi-hole and our local state, but is NOT desired anymore.
	for host, ip := range state.RegisteredRecords {
		host = strings.ToLower(host)
		_, isDesired := desiredMap[host]
		_, existsInPihole := piholeRecords[host]

		if existsInPihole && !isDesired {
			toDelete = append(toDelete, DNSMapping{
				Hostname: host,
				IP:       ip,
			})
		}
	}

	return toAdd, toDelete
}
