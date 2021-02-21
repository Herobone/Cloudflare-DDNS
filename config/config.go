package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

func (c *DNSConfig) getZoneID(api *cloudflare.API) error {
	fromENV := os.Getenv("CLOUDFLARE_ZONE_ID")
	if len(fromENV) > 0 {
		c.ZoneID = fromENV
		zone, err := api.ZoneDetails(fromENV)
		c.ZoneName = zone.Name
		return err
	}
	nameFromENV := os.Getenv("CLOUDFLARE_ZONE_NAME")
	if len(nameFromENV) > 0 {
		zone, err := api.ZoneIDByName(nameFromENV)
		c.ZoneID = zone
		c.ZoneName = nameFromENV
		return err
	}
	return errors.New("neither CLOUDFLARE_ZONE_ID nor CLOUDFLARE_ZONE_NAME were provided")
}

func (c *DNSConfig) getExternalIP() error {
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

func (c *DNSConfig) getDNSName() {
	fromENV := os.Getenv("DDNS_NAME")
	if len(fromENV) > 0 {
		c.DNSName = fromENV
		return
	}
	c.DNSName = "ddns"
}

func (c *DNSConfig) getDNSProxied() {
	fromENV := strings.ToLower(os.Getenv("DDNS_PROXIED"))
	c.Proxied = len(fromENV) > 0 && fromENV != "false"
}

func (c *DNSConfig) getDNSTTL() {
	fromENV := os.Getenv("DDNS_TTL")
	if len(fromENV) > 0 {
		ttl, err := strconv.ParseInt(fromENV, 10, 64)
		if err == nil && (ttl > 120 || ttl == 1) {
			c.TTL = int(ttl)
		}
	} else {
		c.TTL = 120
	}
}

// New creates a new DNSConfig instance and initializes it
func New(api *cloudflare.API) (DNSConfig, error) {
	config := DNSConfig{}
	config.getDNSTTL()
	config.getDNSName()
	config.getDNSProxied()
	err := config.getExternalIP()
	if err != nil {
		return DNSConfig{}, err
	}

	err = config.getZoneID(api)
	if err != nil {
		return DNSConfig{}, err
	}

	return config, nil
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
