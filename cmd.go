package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Herobone/cloudflare-ddns/config"
	"github.com/cloudflare/cloudflare-go"
	"github.com/google/subcommands"
)

type updateCMD struct {
	config.DNSConfig
}

func (*updateCMD) Name() string     { return "update" }
func (*updateCMD) Synopsis() string { return "Update the DNS Record from public IP" }
func (*updateCMD) Usage() string {
	return `print [-capitalize] <some text>:
  Print args to stdout.
`
}

func (p *updateCMD) SetFlags(f *flag.FlagSet) {

	proxiedFromENV := strings.ToLower(os.Getenv("DDNS_PROXIED"))
	proxied := len(proxiedFromENV) > 0 && proxiedFromENV != "false"

	f.BoolVar(&p.Proxied, "proxied", proxied, "Proxy IP over Cloudflare CDN")
	f.BoolVar(&p.Proxied, "p", proxied, "Proxy IP over Cloudflare CDN")

	zoneIDFromENV := os.Getenv("CLOUDFLARE_ZONE_ID")

	f.StringVar(&p.ZoneID, "zone", zoneIDFromENV, "Cloudflare Zone ID")
	f.StringVar(&p.ZoneID, "z", zoneIDFromENV, "Cloudflare Zone ID")

	zoneNameFromENV := os.Getenv("CLOUDFLARE_ZONE_NAME")

	f.StringVar(&p.ZoneName, "zoneName", zoneNameFromENV, "Cloudflare Zone name (i.e. 'example.com')")
	f.StringVar(&p.ZoneName, "zn", zoneNameFromENV, "Cloudflare Zone name (i.e. 'example.com')")

	ttlFromEnv := os.Getenv("DDNS_TTL")
	ttl := 120
	if len(ttlFromEnv) > 0 {
		_ttl, err := strconv.ParseInt(ttlFromEnv, 10, 64)
		if err == nil && (_ttl > 120 || _ttl == 1) {
			ttl = int(_ttl)
		}
	}

	f.IntVar(&p.TTL, "ttl", ttl, "Time To Live for the DNS record")

	nameFromENV := os.Getenv("DDNS_NAME")
	if len(nameFromENV) < 1 {
		nameFromENV = "ddns"
	}

	f.StringVar(&p.DNSName, "name", nameFromENV, "The Name of the DNS Record")
	f.StringVar(&p.DNSName, "n", nameFromENV, "The Name of the DNS Record")
}

func (p *updateCMD) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return subcommands.ExitFailure
	}

	if len(p.ZoneID) > 0 && len(p.ZoneName) == 0 {
		zone, err := api.ZoneDetails(p.ZoneID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return subcommands.ExitFailure
		}
		p.ZoneName = zone.Name
	} else if len(p.ZoneID) == 0 && len(p.ZoneName) > 0 {
		zone, err := api.ZoneIDByName(p.ZoneName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return subcommands.ExitFailure
		}
		p.ZoneID = zone
	} else if len(p.ZoneID) > 0 && len(p.ZoneName) > 0 {
		zone, err := api.ZoneDetails(p.ZoneID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return subcommands.ExitFailure
		}
		if p.ZoneName != zone.Name {
			fmt.Fprintln(os.Stderr, "Zone Name and Zone ID do not reference the same Cloudflare Zone")
			return subcommands.ExitFailure
		}
	}

	update(api, &p.DNSConfig)

	return subcommands.ExitSuccess
}
