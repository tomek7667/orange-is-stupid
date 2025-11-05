package cloudflare

type Client struct {
	Url      string
	ApiToken string
}

func New(apiToken string, endpoint ...string) *Client {
	var url string
	if len(endpoint) > 0 {
		url = endpoint[0]
	} else {
		url = "https://api.cloudflare.com/client/v4"
	}
	c := &Client{
		Url:      url,
		ApiToken: apiToken,
	}
	return c
}
