package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"proxmox-dns-sync/internal/sync"
	"strings"
	"time"
)

// Client is a client for the Proxmox VE API.
type Client struct {
	PveURL     string
	PveToken   string
	httpClient *http.Client
}

// NewClient creates a new Proxmox client.
func NewClient(pveURL, token string) *Client {
	// Disable TLS verification for self-signed certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		PveURL:   strings.TrimSuffix(pveURL, "/"),
		PveToken: token,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		},
	}
}

// GetResources fetches active VMs, LXCs, and Nodes from Proxmox VE API.
func (c *Client) GetResources() ([]sync.ProxmoxResource, error) {
	reqURL := fmt.Sprintf("%s/cluster/resources", c.PveURL)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s", c.PveToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proxmox returned HTTP %d", resp.StatusCode)
	}

	var wrapper struct {
		Data []sync.ProxmoxResource `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}

	var results []sync.ProxmoxResource
	for _, res := range wrapper.Data {
		if res.Type == "node" && res.Status == "online" {
			results = append(results, res)
		} else if (res.Type == "qemu" || res.Type == "lxc") && res.Status == "running" {
			if res.VMID != nil {
				ip, err := c.getVMIP(res.Node, res.Type, *res.VMID)
				if err == nil && ip != "" {
					res.IP = ip
					results = append(results, res)
				}
			}
		}
	}

	return results, nil
}

func (c *Client) getVMIP(node, resType string, vmid int) (string, error) {
	reqURL := ""
	if resType == "lxc" {
		reqURL = fmt.Sprintf("%s/nodes/%s/lxc/%d/interfaces", c.PveURL, node, vmid)
	} else if resType == "qemu" {
		reqURL = fmt.Sprintf("%s/nodes/%s/qemu/%d/agent/network-get-interfaces", c.PveURL, node, vmid)
	} else {
		return "", fmt.Errorf("unknown resource type %s", resType)
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s", c.PveToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	if resType == "lxc" {
		var wrapper struct {
			Data []struct {
				Name string `json:"name"`
				Inet string `json:"inet"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
			return "", err
		}
		for _, iface := range wrapper.Data {
			if iface.Name == "lo" || iface.Inet == "" {
				continue
			}
			parts := strings.Split(iface.Inet, "/")
			if len(parts) > 0 {
				ip := parts[0]
				if ip != "" && !strings.HasPrefix(ip, "127.") {
					return ip, nil
				}
			}
		}
	} else if resType == "qemu" {
		var wrapper struct {
			Data struct {
				Result []struct {
					Name         string `json:"name"`
					IPAddresses []struct {
						Type string `json:"ip-address-type"`
						Addr string `json:"ip-address"`
					} `json:"ip-addresses"`
				} `json:"result"`
			} `json:"data"`
		}
		// Try parsing wrapped format as in test, or check if direct Result wrapper
		dec := json.NewDecoder(resp.Body)
		var raw map[string]json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return "", err
		}

		// Try parsing inside "data"
		if dataVal, ok := raw["data"]; ok {
			_ = json.Unmarshal(dataVal, &wrapper.Data)
		} else {
			// Fallback: parse directly into Result
			_ = json.Unmarshal(raw["result"], &wrapper.Data.Result)
		}

		for _, iface := range wrapper.Data.Result {
			if iface.Name == "lo" {
				continue
			}
			for _, ipAddr := range iface.IPAddresses {
				if ipAddr.Type == "ipv4" {
					ip := ipAddr.Addr
					if ip != "" && !strings.HasPrefix(ip, "127.") {
						return ip, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("no valid IP address found")
}

// URLEncodeHelper helper to URL-encode strings
func URLEncodeHelper(s string) string {
	return url.PathEscape(s)
}
