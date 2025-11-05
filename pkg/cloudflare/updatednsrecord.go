package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"time"
)

type UpdateDNSRecordRequestBody struct {
	Name    string `json:"name"`
	Ip      string `json:"content"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
	Proxied bool   `json:"proxied"`
}

type UpdateDNSRecordResponseBody struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func (c *Client) UpdateDNSRecord(dnsRecordID, zoneID, fullHostname string, addr netip.Addr, cloudflareProxy bool) error {
	var recordType string
	if addr.Is4() {
		recordType = "A"
	} else if addr.Is6() {
		recordType = "AAAA"
	} else {
		return fmt.Errorf("invalid ip address '%s' for hostname '%s'", addr.String(), fullHostname)
	}

	body := &UpdateDNSRecordRequestBody{
		Name:    fullHostname,
		Ip:      addr.String(),
		Type:    recordType,
		Comment: fmt.Sprintf("updated by https://github.com/tomek7667/orange-is-stupid at %s", time.Now().Format(time.RFC3339Nano)),
		Proxied: cloudflareProxy,
	}
	u := fmt.Sprintf("%s/zones/%s/dns_records/%s", c.Url, zoneID, dnsRecordID)
	requestBodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal the body for %s: %w", fullHostname, err)
	}
	req, err := http.NewRequest(http.MethodPut, u, bytes.NewReader(requestBodyBytes))
	if err != nil {
		return fmt.Errorf("faield to create put request to '%s': %w", u, err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.getHttpClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to update dns record on '%s' and body '%s': %w", u, string(requestBodyBytes), err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read bytes from the response to '%s': %w", u, err)
	}
	var response UpdateDNSRecordResponseBody
	err = json.Unmarshal(b, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal update dns record response '%s' from url '%s': %w", string(b), u, err)
	}
	if response.Success {
		return nil
	}
	errS := ""
	for _, err := range response.Errors {
		errS += err.Message + "; "
	}
	return fmt.Errorf("failed to update DNS record for %s: %s", fullHostname, errS)
}
