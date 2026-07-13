package pihole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is a client for the Pi-hole v6 API.
type Client struct {
	PiholeURL      string
	PiholePassword string
	httpClient     *http.Client
}

// NewClient creates a new Pi-hole client.
func NewClient(url, password string) *Client {
	return &Client{
		PiholeURL:      strings.TrimSuffix(url, "/"),
		PiholePassword: password,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetCustomHosts retrieves custom A/AAAA records.
func (c *Client) GetCustomHosts() ([]string, error) {
	reqURL := fmt.Sprintf("%s/api/config/dns/hosts", c.PiholeURL)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-FTL-SID", c.PiholePassword)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pi-hole returned HTTP %d", resp.StatusCode)
	}

	var hosts []string
	if err := json.NewDecoder(resp.Body).Decode(&hosts); err != nil {
		return nil, err
	}

	return hosts, nil
}

// AddCustomHost inserts an A record.
func (c *Client) AddCustomHost(ip, host string) error {
	reqURL := fmt.Sprintf("%s/api/config/dns/hosts", c.PiholeURL)
	payload, err := json.Marshal(fmt.Sprintf("%s %s", ip, host))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-FTL-SID", c.PiholePassword)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pi-hole returned HTTP %d on add", resp.StatusCode)
	}

	return nil
}

// DeleteCustomHost removes an A record.
func (c *Client) DeleteCustomHost(ip, host string) error {
	encodedParam := url.PathEscape(fmt.Sprintf("%s %s", ip, host))
	// Replace %20 with %20 explicitly in path escape
	encodedParam = strings.ReplaceAll(encodedParam, "+", "%20")
	
	reqURL := fmt.Sprintf("%s/api/config/dns/hosts/%s", c.PiholeURL, encodedParam)
	req, err := http.NewRequest("DELETE", reqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-FTL-SID", c.PiholePassword)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pi-hole returned HTTP %d on delete", resp.StatusCode)
	}

	return nil
}
