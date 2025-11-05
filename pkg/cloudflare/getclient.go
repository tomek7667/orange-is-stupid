package cloudflare

import (
	"net/http"
	"time"
)

type headerRoundTripper struct {
	base   http.RoundTripper
	header http.Header
}

func (hrt headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	for k, vs := range hrt.header {
		for _, v := range vs {
			r.Header.Add(k, v)
		}
	}
	return hrt.base.RoundTrip(r)
}

func (c *Client) getHttpClient() *http.Client {
	base := http.DefaultTransport.(*http.Transport).Clone()

	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: headerRoundTripper{
			base: base,
			header: http.Header{
				"Authorization": []string{"Bearer " + c.ApiToken},
				"User-Agent":    []string{"cloudflare-client/1.0"},
				"Accept":        []string{"application/json"},
			},
		},
	}
}
