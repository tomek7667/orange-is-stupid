package inet

import (
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"strings"
)

func GetCurrentAddr() (*netip.Addr, error) {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return nil, fmt.Errorf("failed to check your ip - got %d status code from 'https://checkip.amazonaws.com': %w", resp.StatusCode, err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read all body from aws: %w", err)
	}
	s := strings.TrimSpace(string(b))
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ip address '%s': %w", s, err)
	}
	return &addr, nil
}
