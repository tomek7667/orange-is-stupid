package cloudflare

import (
	"encoding/json"
	"fmt"
	"io"
)

type GetZonesResponseBody struct {
	Result     []Zone     `json:"result"`
	ResultInfo ResultInfo `json:"result_info"`
}

type Zone struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	Type                string   `json:"type"`
	DevelopmentMode     int      `json:"development_mode"`
	NameServers         []string `json:"name_servers"`
	OriginalNameServers []string `json:"original_name_servers"`
}

func (c *Client) GetZones() ([]Zone, error) {
	u := fmt.Sprintf("%s/zones", c.Url)
	resp, err := c.getHttpClient().Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to get '%s': %w", u, err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes from the response to '%s': %w", u, err)
	}
	var response GetZonesResponseBody
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal zones response '%s': %w", string(b), err)
	}
	// TODO: improve by paginating to get all zones
	return response.Result, nil
}
