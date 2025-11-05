package cloudflare

import (
	"fmt"
	"net/netip"

	"orange-is-stupid/pkg/utils"
)

func (c *Client) UpsertUrl(urlS string, inet string, cloudflareProxy bool) error {
	url, err := utils.UnpackSubdomainUrl(urlS)
	if err != nil {
		return fmt.Errorf("failed to unpack subdomain url '%s': %w", urlS, err)
	}
	addr, err := netip.ParseAddr(inet)
	if err != nil {
		return fmt.Errorf("failed to parse ip='%s': %w", inet, err)
	}

	// getting the zone id
	zones, err := c.GetZones()
	if err != nil {
		return fmt.Errorf("failed to get zones: %w", err)
	}
	var zoneID string
	for _, zone := range zones {
		if zone.Name == url.ZoneName {
			zoneID = zone.ID
			break
		}
	}
	if zoneID == "" {
		return fmt.Errorf("zone for '%s' not found; zones found: %d", urlS, len(zones))
	}

	// get the record; if it exists update if not create
	dnsRecords, err := c.GetDNSRecords(zoneID)
	if err != nil {
		return fmt.Errorf("failed to get DNS records for zone ID '%s': %w", zoneID, err)
	}
	fullHostname := url.GetFullHostname()
	var dnsRecordID string
	for _, dnsRecord := range dnsRecords {
		if dnsRecord.Name == fullHostname {
			dnsRecordID = dnsRecord.ID
			break
		}
	}

	shouldBeUpdated := dnsRecordID != ""
	if shouldBeUpdated {
		return c.UpdateDNSRecord(dnsRecordID, zoneID, fullHostname, addr, cloudflareProxy)
	} else {
		return c.CreateDNSRecord(zoneID, fullHostname, addr, cloudflareProxy)
	}
}
