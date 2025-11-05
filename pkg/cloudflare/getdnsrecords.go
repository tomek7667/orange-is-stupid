package cloudflare

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type GetDNSRecordsResponseBody struct {
	Result     []DNSRecord `json:"result"`
	ResultInfo ResultInfo  `json:"result_info"`
}

type DNSRecord struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	TTL        int       `json:"ttl"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
	Priority   int       `json:"priority,omitempty"`
}

func (c *Client) GetDNSRecords(zoneID string) ([]DNSRecord, error) {
	u := fmt.Sprintf("%s/zones/%s/dns_records", c.Url, zoneID)
	resp, err := c.getHttpClient().Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to get '%s': %w", u, err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes from the response to '%s': %w", u, err)
	}
	var response GetDNSRecordsResponseBody
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dns records response '%s' from url '%s': %w", string(b), u, err)
	}
	// TODO: improve by paginating to get all DNS records. Cloudflare's default limit is 100.
	return response.Result, nil
}
