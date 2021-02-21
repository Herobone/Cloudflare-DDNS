package config

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudflare/cloudflare-go"
)

// GetExternalIP fetches the external IP from the Internet and saves it
func (c *DNSConfig) GetExternalIP() error {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.ExternalIP = string(body)
	return nil
}

// DNSConfig is used for configuring the DNS Record
type DNSConfig struct {
	TTL        int
	Proxied    bool
	DNSName    string
	ExternalIP string
	ZoneID     string
	ZoneName   string
}

// GetName returns the name of the Record in format DNSName.ZoneName
func (c *DNSConfig) GetName() string {
	return fmt.Sprintf("%s.%s", c.DNSName, c.ZoneName)
}

// ToDNSRecord creates a cloudflare.DNSRecord from the struct
func (c *DNSConfig) ToDNSRecord() cloudflare.DNSRecord {
	return cloudflare.DNSRecord{
		Name:    c.DNSName,
		Content: c.ExternalIP,
		Type:    "A",
		Proxied: c.Proxied,
		TTL:     c.TTL,
	}
}
