package main

import (
	"errors"
	"log"

	"github.com/Herobone/cloudflare-ddns/config"
	"github.com/cloudflare/cloudflare-go"
)

func update(api *cloudflare.API, dnsConfig *config.DNSConfig) error {

	log.Println("External IP:", dnsConfig.ExternalIP)

	log.Println("Zone ID is:", dnsConfig.ZoneID)

	log.Println("DNS Name is", dnsConfig.DNSName)

	log.Println("IP will be proxied", dnsConfig.Proxied)

	log.Println("DNS TTL is", dnsConfig.TTL)

	results, err := api.DNSRecords(dnsConfig.ZoneID, cloudflare.DNSRecord{Name: dnsConfig.GetName()})
	if err != nil {
		return err
	}

	if len(results) > 1 {
		return errors.New("too many records")
	} else if len(results) == 1 {
		log.Println("Updating record")
		err := api.UpdateDNSRecord(dnsConfig.ZoneID, results[0].ID, dnsConfig.ToDNSRecord())
		if err != nil {
			return err
		}
	} else {
		log.Println("Creating record")
		result, err := api.CreateDNSRecord(dnsConfig.ZoneID, dnsConfig.ToDNSRecord())
		if err != nil {
			return err
		}
		log.Println(result)
	}
	return nil
}
